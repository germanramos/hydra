package etcd

import (
	// "fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
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
	Store      store.Store
	PeerServer *server.PeerServer
}

func New() *Etcd {
	etcd := new(Etcd)
	return etcd
}

func (etcd *Etcd) Start(config *config.Config) {
	// Enable options.
	if config.VeryVeryVerbose {
		log.Verbose = true
		raft.SetLogLevel(raft.Trace)
	} else if config.VeryVerbose {
		log.Verbose = true
		raft.SetLogLevel(raft.Debug)
	} else if config.Verbose {
		log.Verbose = true
	}
	// if config.CPUProfileFile != "" {
	// 	profile(config.CPUProfileFile)
	// }

	if config.DataDir == "" {
		log.Fatal("The data dir was not set and could not be guessed from machine name")
	}

	// Create data directory if it doesn't already exist.
	if err := os.MkdirAll(config.DataDir, 0744); err != nil {
		log.Fatalf("Unable to create path: %s", err)
	}

	// Warn people if they have an info file
	info := filepath.Join(config.DataDir, "info")
	if _, err := os.Stat(info); err == nil {
		log.Warnf("All cached configuration is now ignored. The file %s can be removed.", info)
	}

	var mbName string
	if config.Trace() {
		mbName = config.MetricsBucketName()
		runtime.SetBlockProfileRate(1)
	}

	mb := metrics.NewBucket(mbName)

	if config.GraphiteHost != "" {
		err := mb.Publish(config.GraphiteHost)
		if err != nil {
			panic(err)
		}
	}

	// Retrieve CORS configuration
	corsInfo, err := ehttp.NewCORSInfo(config.CorsOrigins)
	if err != nil {
		log.Fatal("CORS:", err)
	}

	// Create etcd key-value store and registry.
	etcd.Store = store.New()
	registry := server.NewRegistry(etcd.Store)

	// Create stats objects
	followersStats := server.NewRaftFollowersStats(config.Name)
	serverStats := server.NewRaftServerStats(config.Name)

	// Calculate all of our timeouts
	heartbeatTimeout := time.Duration(config.Peer.HeartbeatTimeout) * time.Millisecond
	electionTimeout := time.Duration(config.Peer.ElectionTimeout) * time.Millisecond
	dialTimeout := (3 * heartbeatTimeout) + electionTimeout
	responseHeaderTimeout := (3 * heartbeatTimeout) + electionTimeout

	// Create peer server
	psConfig := server.PeerServerConfig{
		Name:           config.Name,
		Scheme:         config.PeerTLSInfo().Scheme(),
		URL:            config.Peer.Addr,
		SnapshotCount:  config.SnapshotCount,
		MaxClusterSize: config.MaxClusterSize,
		RetryTimes:     config.MaxRetryAttempts,
		RetryInterval:  config.RetryInterval,
	}
	etcd.PeerServer = server.NewPeerServer(psConfig, registry, etcd.Store, &mb, followersStats, serverStats)

	var psListener net.Listener
	if psConfig.Scheme == "https" {
		peerServerTLSConfig, err := config.PeerTLSInfo().ServerConfig()
		if err != nil {
			log.Fatal("peer server TLS error: ", err)
		}

		psListener, err = server.NewTLSListener(config.Peer.BindAddr, peerServerTLSConfig)
		if err != nil {
			log.Fatal("Failed to create peer listener: ", err)
		}
	} else {
		psListener, err = server.NewListener(config.Peer.BindAddr)
		if err != nil {
			log.Fatal("Failed to create peer listener: ", err)
		}
	}

	// Create raft transporter and server
	raftTransporter := server.NewTransporter(followersStats, serverStats, registry, heartbeatTimeout, dialTimeout, responseHeaderTimeout)
	if psConfig.Scheme == "https" {
		raftClientTLSConfig, err := config.PeerTLSInfo().ClientConfig()
		if err != nil {
			log.Fatal("raft client TLS error: ", err)
		}
		raftTransporter.SetTLSConfig(*raftClientTLSConfig)
	}
	raftServer, err := raft.NewServer(config.Name, config.DataDir, raftTransporter, etcd.Store, etcd.PeerServer, "")
	if err != nil {
		log.Fatal(err)
	}
	raftServer.SetElectionTimeout(electionTimeout)
	raftServer.SetHeartbeatInterval(heartbeatTimeout)
	etcd.PeerServer.SetRaftServer(raftServer)

	// Create etcd server
	s := server.New(config.Name, config.Addr, etcd.PeerServer, registry, etcd.Store, &mb)

	if config.Trace() {
		s.EnableTracing()
	}

	var sListener net.Listener
	if config.EtcdTLSInfo().Scheme() == "https" {
		etcdServerTLSConfig, err := config.EtcdTLSInfo().ServerConfig()
		if err != nil {
			log.Fatal("etcd TLS error: ", err)
		}

		sListener, err = server.NewTLSListener(config.BindAddr, etcdServerTLSConfig)
		if err != nil {
			log.Fatal("Failed to create TLS etcd listener: ", err)
		}
	} else {
		sListener, err = server.NewListener(config.BindAddr)
		if err != nil {
			log.Fatal("Failed to create etcd listener: ", err)
		}
	}

	etcd.PeerServer.SetServer(s)
	etcd.PeerServer.Start(config.Snapshot, config.Peers)

	go func() {
		log.Infof("peer server [name %s, listen on %s, advertised url %s]", etcd.PeerServer.Config.Name, psListener.Addr(), etcd.PeerServer.Config.URL)
		sHTTP := &ehttp.CORSHandler{etcd.PeerServer.HTTPHandler(), corsInfo}
		log.Fatal(http.Serve(psListener, sHTTP))
	}()

	log.Infof("etcd server [name %s, listen on %s, advertised url %s]", s.Name, sListener.Addr(), s.URL())
	sHTTP := &ehttp.CORSHandler{s.HTTPHandler(), corsInfo}
	log.Fatal(http.Serve(sListener, sHTTP))
}

// func (etcd *Etcd) Start() {
// 	etcd.PeerServer.Start(config.Snapshot, config.Peers)

// 	go func() {
// 		log.Infof("peer server [name %s, listen on %s, advertised url %s]", etcd.PeerServer.Config.Name, psListener.Addr(), etcd.PeerServer.Config.URL)
// 		sHTTP := &ehttp.CORSHandler{ps.HTTPHandler(), corsInfo}
// 		log.Fatal(http.Serve(psListener, sHTTP))
// 	}()

// 	log.Infof("etcd server [name %s, listen on %s, advertised url %s]", s.Name, sListener.Addr(), s.URL())
// 	sHTTP := &ehttp.CORSHandler{s.HTTPHandler(), corsInfo}
// 	log.Fatal(http.Serve(sListener, sHTTP))
// }
