package etcd

import (
	// "fmt"
	"net"
	"net/http"
	// "os"
	// "path/filepath"
	// "runtime"
	"time"

	// hlog "github.com/innotech/hydra/log"

	"github.com/innotech/hydra/vendors/github.com/coreos/etcd/third_party/github.com/coreos/raft"

	"github.com/innotech/hydra/vendors/github.com/coreos/etcd/config"
	ehttp "github.com/innotech/hydra/vendors/github.com/coreos/etcd/http"
	"github.com/innotech/hydra/vendors/github.com/coreos/etcd/log"
	"github.com/innotech/hydra/vendors/github.com/coreos/etcd/metrics"
	"github.com/innotech/hydra/vendors/github.com/coreos/etcd/server"
	"github.com/innotech/hydra/vendors/github.com/coreos/etcd/store"
)

type Etcd struct {
	Config *config.Config
	// Store              store.Store
	EtcdServer         *server.Server
	PeerServer         *server.PeerServer
	PeerServerListener net.Listener
}

func New(conf *config.Config) *Etcd {
	etcd := new(Etcd)
	etcd.Config = conf
	return etcd
}

func (e *Etcd) configMetrics() metrics.Bucket {
	// var mbName string
	// if e.Config.Trace() {
	// 	mbName = e.Config.MetricsBucketName()
	// 	runtime.SetBlockProfileRate(1)
	// }

	// mb := metrics.NewBucket(mbName)

	// if e.Config.GraphiteHost != "" {
	// 	err := mb.Publish(e.Config.GraphiteHost)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// }

	// return mb

	return metrics.NewBucket("")
}

func (e *Etcd) configPsListener(psConfig server.PeerServerConfig) net.Listener {
	var psListener net.Listener
	var err error
	// hlog.Info(e.Config.Peer.BindAddr)
	if psConfig.Scheme == "https" {
		peerServerTLSConfig, err := e.Config.PeerTLSInfo().ServerConfig()
		if err != nil {
			log.Fatal("peer server TLS error: ", err)
		}

		psListener, err = server.NewTLSListener(e.Config.Peer.BindAddr, peerServerTLSConfig)
		if err != nil {
			log.Fatal("Failed to create peer listener: ", err)
		}
	} else {
		psListener, err = server.NewListener(e.Config.Peer.BindAddr)
		if err != nil {
			log.Fatal("Failed to create peer listener: ", err)
		}
	}

	return psListener
}

func (e *Etcd) configEtcdListener() net.Listener {
	var sListener net.Listener
	var err error
	if e.Config.EtcdTLSInfo().Scheme() == "https" {
		etcdServerTLSConfig, err := e.Config.EtcdTLSInfo().ServerConfig()
		if err != nil {
			log.Fatal("etcd TLS error: ", err)
		}

		sListener, err = server.NewTLSListener(e.Config.BindAddr, etcdServerTLSConfig)
		if err != nil {
			log.Fatal("Failed to create TLS etcd listener: ", err)
		}
	} else {
		sListener, err = server.NewListener(e.Config.BindAddr)
		if err != nil {
			log.Fatal("Failed to create etcd listener: ", err)
		}
	}

	return sListener
}

func (e *Etcd) Load() {
	mb := e.configMetrics()

	// Create etcd key-value store and registry.
	store := store.New()
	registry := server.NewRegistry(store)

	// Create stats objects
	followersStats := server.NewRaftFollowersStats(e.Config.Name)
	serverStats := server.NewRaftServerStats(e.Config.Name)

	// Calculate all of our timeouts
	heartbeatTimeout := time.Duration(e.Config.Peer.HeartbeatTimeout) * time.Millisecond
	electionTimeout := time.Duration(e.Config.Peer.ElectionTimeout) * time.Millisecond
	dialTimeout := (3 * heartbeatTimeout) + electionTimeout
	responseHeaderTimeout := (3 * heartbeatTimeout) + electionTimeout

	// Create peer server
	psConfig := server.PeerServerConfig{
		Name:           e.Config.Name,
		Scheme:         e.Config.PeerTLSInfo().Scheme(),
		URL:            e.Config.Peer.Addr,
		SnapshotCount:  e.Config.SnapshotCount,
		MaxClusterSize: e.Config.MaxClusterSize,
		RetryTimes:     e.Config.MaxRetryAttempts,
		RetryInterval:  e.Config.RetryInterval,
	}
	ps := server.NewPeerServer(psConfig, registry, store, &mb, followersStats, serverStats)

	psListener := e.configPsListener(psConfig)

	// Create raft transporter and server
	raftTransporter := server.NewTransporter(followersStats, serverStats, registry, heartbeatTimeout, dialTimeout, responseHeaderTimeout)
	if psConfig.Scheme == "https" {
		raftClientTLSConfig, err := e.Config.PeerTLSInfo().ClientConfig()
		if err != nil {
			log.Fatal("raft client TLS error: ", err)
		}
		raftTransporter.SetTLSConfig(*raftClientTLSConfig)
	}
	raftServer, err := raft.NewServer(e.Config.Name, e.Config.DataDir, raftTransporter, store, ps, "")
	if err != nil {
		log.Fatal(err)
	}
	raftServer.SetElectionTimeout(electionTimeout)
	raftServer.SetHeartbeatInterval(heartbeatTimeout)
	ps.SetRaftServer(raftServer)

	// Create etcd server
	s := server.New(e.Config.Name, e.Config.Addr, ps, registry, store, nil)

	ps.SetServer(s)

	e.EtcdServer = s
	e.PeerServer = ps
	e.PeerServerListener = psListener
}

func (e *Etcd) Start(withEtcdServer string) {
	e.PeerServer.Start(e.Config.Snapshot, e.Config.Peers)

	if withEtcdServer == "TEST" {
		sListener := e.configEtcdListener()
		go func() {
			log.Infof("etcd server [name %s, listen on %s, advertised url %s]", e.EtcdServer.Name, sListener.Addr(), e.EtcdServer.URL())
			corsInfo, err := ehttp.NewCORSInfo(e.Config.CorsOrigins)
			if err != nil {
				log.Fatal("CORS:", err)
			}
			sHTTP := &ehttp.CORSHandler{e.EtcdServer.HTTPHandler(), corsInfo}
			log.Fatal(http.Serve(sListener, sHTTP))
		}()
	}

	log.Infof("peer server [name %s, listen on %s, advertised url %s]", e.PeerServer.Config.Name, e.PeerServerListener.Addr(), e.PeerServer.Config.URL)
	// Retrieve CORS configuration
	corsInfo, err := ehttp.NewCORSInfo(e.Config.CorsOrigins)
	if err != nil {
		log.Fatal("CORS:", err)
	}
	sHTTP := &ehttp.CORSHandler{e.PeerServer.HTTPHandler(), corsInfo}
	log.Fatal(http.Serve(e.PeerServerListener, sHTTP))
}
