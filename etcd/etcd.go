package etcd

import (
	// "fmt"
	"net"
	"net/http"
	// "os"
	// "path/filepath"
	// "runtime"
	"time"

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
	// Store      store.Store
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

func (e *Etcd) Load() {
	// IGNORE: not load etcd configuration
	// IGNORE: not set verbose etcd log
	// IGNORE: not set etcd profiling
	// IGNORE: check data dir configuration
	// IGNORE: create data dir
	// IGNORE: info file warning

	// Create metrics bucket
	mb := e.configMetrics()

	// IGNORE: retrieve CORS configuration

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

	// Create peer listener
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
	raftServer, err := raft.NewServer(e.Config.Name, e.Config.DataDir, raftTransporter, store, e.PeerServer, "")
	if err != nil {
		log.Fatal(err)
	}
	raftServer.SetElectionTimeout(electionTimeout)
	raftServer.SetHeartbeatInterval(heartbeatTimeout)
	ps.SetRaftServer(raftServer)

	// Create etcd server
	// s := server.New(etcd.Config.Name, etcd.Config.Addr, etcd.PeerServer, registry, etcd.Store, &mb)

	// if etcd.Config.Trace() {
	// 	s.EnableTracing()
	// }
	// etcd.PeerServer.SetServer(s)

	// IGNORE: etcd server listener

	e.PeerServer = ps
	e.PeerServerListener = psListener
}

func (e *Etcd) Start() {
	e.PeerServer.Start(e.Config.Snapshot, e.Config.Peers)

	// go func() {
	// 	var Permanent time.Time
	// 	log.Infof("Sleeping 5s...")
	// 	time.Sleep(1000 * time.Millisecond)
	// 	log.Infof("Setting foo = bar")
	// 	// _, err := store.Set("/foo", false, "bar", Permanent)
	// 	// s.Store().CommandFactory().CreateSetCommand(key, dir, value, expireTime)
	// 	// _, err := s.Store().Set("/foo", false, "bar", Permanent)
	// 	c := s.Store().CommandFactory().CreateSetCommand("/foo", false, "bar", Permanent)
	// 	result, err := etcd.PeerServer.RaftServer().Do(c)
	// 	if err != nil {
	// 		// return err
	// 		log.Fatal("Failed 1 to set key", err)
	// 	}

	// 	// if result == nil {
	// 	// 	// return etcdErr.NewError(300, "Empty result from raft", s.Store().Index())
	// 	// 	log.Fatal("Failed 2 to set key", err)
	// 	// }
	// 	// // if err != nil {
	// 	// // 	log.Fatal("Failed to set key", err)
	// 	// // }

	// 	log.Infof("Sleeping 2000ms...")
	// 	time.Sleep(2000 * time.Millisecond)
	// 	g, err := s.Store().Get("/foo", false, false)
	// 	if err != nil {
	// 		log.Fatal("Failed to get key", err)
	// 	}
	// 	log.Infof("printing results...")
	// 	log.Infof(g.Node.Key)
	// 	log.Infof(g.Node.Value)
	// }()

	log.Infof("peer server [name %s, listen on %s, advertised url %s]", e.PeerServer.Config.Name, e.PeerServerListener.Addr(), e.PeerServer.Config.URL)

	// Retrieve CORS configuration
	corsInfo, err := ehttp.NewCORSInfo(e.Config.CorsOrigins)
	if err != nil {
		log.Fatal("CORS:", err)
	}
	sHTTP := &ehttp.CORSHandler{e.PeerServer.HTTPHandler(), corsInfo}

	log.Fatal(http.Serve(e.PeerServerListener, sHTTP))
}
