package etcd

import (
	"fmt"
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
	Config     *config.Config
	Store      store.Store
	PeerServer *server.PeerServer
}

func New(conf *config.Config) *Etcd {
	etcd := new(Etcd)
	etcd.Config = conf
	return etcd
}

func (etcd *Etcd) configMetrics() metrics.Bucket {
	var mbName string
	if etcd.Config.Trace() {
		mbName = etcd.Config.MetricsBucketName()
		runtime.SetBlockProfileRate(1)
	}

	mb := metrics.NewBucket(mbName)

	if etcd.Config.GraphiteHost != "" {
		err := mb.Publish(etcd.Config.GraphiteHost)
		if err != nil {
			panic(err)
		}
	}

	return mb
}

func (etcd *Etcd) configPsListener(psConfig server.PeerServerConfig) net.Listener {
	var psListener net.Listener
	var err error
	if psConfig.Scheme == "https" {
		peerServerTLSConfig, err := etcd.Config.PeerTLSInfo().ServerConfig()
		if err != nil {
			log.Fatal("peer server TLS error: ", err)
		}

		psListener, err = server.NewTLSListener(etcd.Config.Peer.BindAddr, peerServerTLSConfig)
		if err != nil {
			log.Fatal("Failed to create peer listener: ", err)
		}
	} else {
		psListener, err = server.NewListener(etcd.Config.Peer.BindAddr)
		if err != nil {
			log.Fatal("Failed to create peer listener: ", err)
		}
	}

	return psListener
}

func (etcd *Etcd) Start() {
	// Load configuration.
	// etcd.Config = config.New()
	// if err := etcd.Config.Load(os.Args[1:]); err != nil {
	// 	fmt.Println(server.Usage() + "\n")
	// 	fmt.Println(err.Error() + "\n")
	// 	os.Exit(1)
	// } else if etcd.Config.ShowVersion {
	// 	fmt.Println(server.ReleaseVersion)
	// 	os.Exit(0)
	// } else if etcd.Config.ShowHelp {
	// 	fmt.Println(server.Usage() + "\n")
	// 	os.Exit(0)
	// }

	// Enable options.
	if etcd.Config.VeryVeryVerbose {
		log.Verbose = true
		raft.SetLogLevel(raft.Trace)
	} else if etcd.Config.VeryVerbose {
		log.Verbose = true
		raft.SetLogLevel(raft.Debug)
	} else if etcd.Config.Verbose {
		log.Verbose = true
	}
	// if etcd.Config.CPUProfileFile != "" {
	// 	profile(etcd.Config.CPUProfileFile)
	// }

	if etcd.Config.DataDir == "" {
		log.Fatal("The data dir was not set and could not be guessed from machine name")
	}

	// Create data directory if it doesn't already exist.
	if err := os.MkdirAll(etcd.Config.DataDir, 0744); err != nil {
		log.Fatalf("Unable to create path: %s", err)
	}

	// Warn people if they have an info file
	info := filepath.Join(etcd.Config.DataDir, "info")
	if _, err := os.Stat(info); err == nil {
		log.Warnf("All cached configuration is now ignored. The file %s can be removed.", info)
	}

	// var mbName string
	// if etcd.Config.Trace() {
	// 	mbName = etcd.Config.MetricsBucketName()
	// 	runtime.SetBlockProfileRate(1)
	// }

	// mb := metrics.NewBucket(mbName)

	// if etcd.Config.GraphiteHost != "" {
	// 	err := mb.Publish(etcd.Config.GraphiteHost)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// }
	mb := etcd.configMetrics()

	// Retrieve CORS configuration
	corsInfo, err := ehttp.NewCORSInfo(etcd.Config.CorsOrigins)
	if err != nil {
		log.Fatal("CORS:", err)
	}

	// Create etcd key-value store and registry.
	etcd.Store = store.New()
	registry := server.NewRegistry(etcd.Store)

	// Create stats objects
	followersStats := server.NewRaftFollowersStats(etcd.Config.Name)
	serverStats := server.NewRaftServerStats(etcd.Config.Name)

	// Calculate all of our timeouts
	heartbeatTimeout := time.Duration(etcd.Config.Peer.HeartbeatTimeout) * time.Millisecond
	electionTimeout := time.Duration(etcd.Config.Peer.ElectionTimeout) * time.Millisecond
	dialTimeout := (3 * heartbeatTimeout) + electionTimeout
	responseHeaderTimeout := (3 * heartbeatTimeout) + electionTimeout

	// Create peer server
	psConfig := server.PeerServerConfig{
		Name:           etcd.Config.Name,
		Scheme:         etcd.Config.PeerTLSInfo().Scheme(),
		URL:            etcd.Config.Peer.Addr,
		SnapshotCount:  etcd.Config.SnapshotCount,
		MaxClusterSize: etcd.Config.MaxClusterSize,
		RetryTimes:     etcd.Config.MaxRetryAttempts,
		RetryInterval:  etcd.Config.RetryInterval,
	}
	etcd.PeerServer = server.NewPeerServer(psConfig, registry, etcd.Store, &mb, followersStats, serverStats)

	// var psListener net.Listener
	// if psConfig.Scheme == "https" {
	// 	peerServerTLSConfig, err := etcd.Config.PeerTLSInfo().ServerConfig()
	// 	if err != nil {
	// 		log.Fatal("peer server TLS error: ", err)
	// 	}

	// 	psListener, err = server.NewTLSListener(etcd.Config.Peer.BindAddr, peerServerTLSConfig)
	// 	if err != nil {
	// 		log.Fatal("Failed to create peer listener: ", err)
	// 	}
	// } else {
	// 	psListener, err = server.NewListener(etcd.Config.Peer.BindAddr)
	// 	if err != nil {
	// 		log.Fatal("Failed to create peer listener: ", err)
	// 	}
	// }
	var psListener net.Listener
	psListener = etcd.configPsListener(psConfig)

	// Create raft transporter and server
	raftTransporter := server.NewTransporter(followersStats, serverStats, registry, heartbeatTimeout, dialTimeout, responseHeaderTimeout)
	if psConfig.Scheme == "https" {
		raftClientTLSConfig, err := etcd.Config.PeerTLSInfo().ClientConfig()
		if err != nil {
			log.Fatal("raft client TLS error: ", err)
		}
		raftTransporter.SetTLSConfig(*raftClientTLSConfig)
	}
	raftServer, err := raft.NewServer(etcd.Config.Name, etcd.Config.DataDir, raftTransporter, etcd.Store, etcd.PeerServer, "")
	if err != nil {
		log.Fatal(err)
	}
	raftServer.SetElectionTimeout(electionTimeout)
	raftServer.SetHeartbeatInterval(heartbeatTimeout)
	etcd.PeerServer.SetRaftServer(raftServer)

	// Create etcd server
	s := server.New(etcd.Config.Name, etcd.Config.Addr, etcd.PeerServer, registry, etcd.Store, &mb)

	if etcd.Config.Trace() {
		s.EnableTracing()
	}

	// var sListener net.Listener
	// if etcd.Config.EtcdTLSInfo().Scheme() == "https" {
	// 	etcdServerTLSConfig, err := etcd.Config.EtcdTLSInfo().ServerConfig()
	// 	if err != nil {
	// 		log.Fatal("etcd TLS error: ", err)
	// 	}

	// 	sListener, err = server.NewTLSListener(etcd.Config.BindAddr, etcdServerTLSConfig)
	// 	if err != nil {
	// 		log.Fatal("Failed to create TLS etcd listener: ", err)
	// 	}
	// } else {
	// 	sListener, err = server.NewListener(etcd.Config.BindAddr)
	// 	if err != nil {
	// 		log.Fatal("Failed to create etcd listener: ", err)
	// 	}
	// }

	etcd.PeerServer.SetServer(s)
	etcd.PeerServer.Start(etcd.Config.Snapshot, etcd.Config.Peers)

	// go func() {
	// 	var Permanent time.Time
	// 	log.Infof("Sleeping 5s...")
	// 	time.Sleep(1000 * time.Millisecond)
	// 	log.Infof("Setting foo = bar")
	// 	_, err := etcd.Store.Set("/foo", false, "bar", Permanent)
	// 	if err != nil {
	// 		log.Fatal("Failed to set key", err)
	// 	}

	// 	log.Infof("Sleeping 200ms...")
	// 	time.Sleep(1000 * time.Millisecond)
	// 	g, err := etcd.Store.Get("/foo", false, false)
	// 	if err != nil {
	// 		log.Fatal("Failed to get key", err)
	// 	}
	// 	log.Infof("printing results...")
	// 	log.Infof(g.Node.Key)
	// 	log.Infof(g.Node.Value)
	// }()

	go func() {
		// var Permanent time.Time
		// log.Infof("Sleeping 5s...")
		// time.Sleep(1000 * time.Millisecond)
		// log.Infof("Setting foo = bar")
		// // _, err := store.Set("/foo", false, "bar", Permanent)
		// // s.Store().CommandFactory().CreateSetCommand(key, dir, value, expireTime)
		// // _, err := s.Store().Set("/foo", false, "bar", Permanent)
		// c := s.Store().CommandFactory().CreateSetCommand("/foo", false, "bar", Permanent)
		// result, err := etcd.PeerServer.RaftServer().Do(c)
		// if err != nil {
		// 	// return err
		// 	log.Fatal("Failed 1 to set key", err)
		// }

		// if result == nil {
		// 	// return etcdErr.NewError(300, "Empty result from raft", s.Store().Index())
		// 	log.Fatal("Failed 2 to set key", err)
		// }
		// // if err != nil {
		// // 	log.Fatal("Failed to set key", err)
		// // }

		log.Infof("Sleeping 200ms...")
		time.Sleep(10000 * time.Millisecond)
		g, err := s.Store().Get("/foo", false, false)
		if err != nil {
			log.Fatal("Failed to get key", err)
		}
		log.Infof("printing results...")
		log.Infof(g.Node.Key)
		log.Infof(g.Node.Value)
	}()

	// go func() {
	log.Infof("peer server [name %s, listen on %s, advertised url %s]", etcd.PeerServer.Config.Name, psListener.Addr(), etcd.PeerServer.Config.URL)
	sHTTP := &ehttp.CORSHandler{etcd.PeerServer.HTTPHandler(), corsInfo}
	log.Fatal(http.Serve(psListener, sHTTP))
	// h := waitHandler{w, ps.HTTPHandler()}
	// log.Fatal(http.Serve(psListener, &h))
	// }()

	// log.Infof("etcd server [name %s, listen on %s, advertised url %s]", s.Name, sListener.Addr(), s.URL())
	// sHTTP := &ehttp.CORSHandler{s.HTTPHandler(), corsInfo}
	// log.Fatal(http.Serve(sListener, sHTTP))
}

func (etcd *Etcd) Start2() {
	// Load configuration.
	var config = config.New()
	if err := config.Load(os.Args[1:]); err != nil {
		fmt.Println(server.Usage() + "\n")
		fmt.Println(err.Error() + "\n")
		os.Exit(1)
	} else if config.ShowVersion {
		fmt.Println(server.ReleaseVersion)
		os.Exit(0)
	} else if config.ShowHelp {
		fmt.Println(server.Usage() + "\n")
		os.Exit(0)
	}

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
	store := store.New()
	registry := server.NewRegistry(store)

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
	ps := server.NewPeerServer(psConfig, registry, store, &mb, followersStats, serverStats)

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
	raftServer, err := raft.NewServer(config.Name, config.DataDir, raftTransporter, store, ps, "")
	if err != nil {
		log.Fatal(err)
	}
	raftServer.SetElectionTimeout(electionTimeout)
	raftServer.SetHeartbeatInterval(heartbeatTimeout)
	ps.SetRaftServer(raftServer)

	// Create etcd server
	s := server.New(config.Name, config.Addr, ps, registry, store, &mb)

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

	ps.SetServer(s)
	ps.Start(config.Snapshot, config.Peers)

	go func() {
		var Permanent time.Time
		log.Infof("Sleeping 5s...")
		time.Sleep(1000 * time.Millisecond)
		log.Infof("Setting foo = bar")
		// _, err := store.Set("/foo", false, "bar", Permanent)
		// s.Store().CommandFactory().CreateSetCommand(key, dir, value, expireTime)
		// _, err := s.Store().Set("/foo", false, "bar", Permanent)
		c := s.Store().CommandFactory().CreateSetCommand("/foo", false, "bar", Permanent)
		result, err := ps.RaftServer().Do(c)
		if err != nil {
			// return err
			log.Fatal("Failed 1 to set key", err)
		}

		if result == nil {
			// return etcdErr.NewError(300, "Empty result from raft", s.Store().Index())
			log.Fatal("Failed 2 to set key", err)
		}
		// if err != nil {
		// 	log.Fatal("Failed to set key", err)
		// }

		log.Infof("Sleeping 200ms...")
		time.Sleep(10000 * time.Millisecond)
		g, err := store.Get("/foo", false, false)
		if err != nil {
			log.Fatal("Failed to get key", err)
		}
		log.Infof("printing results...")
		log.Infof(g.Node.Key)
		log.Infof(g.Node.Value)
	}()

	go func() {
		log.Infof("peer server [name %s, listen on %s, advertised url %s]", ps.Config.Name, psListener.Addr(), ps.Config.URL)
		sHTTP := &ehttp.CORSHandler{ps.HTTPHandler(), corsInfo}
		log.Fatal(http.Serve(psListener, sHTTP))
	}()

	log.Infof("etcd server [name %s, listen on %s, advertised url %s]", s.Name, sListener.Addr(), s.URL())
	sHTTP := &ehttp.CORSHandler{s.HTTPHandler(), corsInfo}
	log.Fatal(http.Serve(sListener, sHTTP))
}
