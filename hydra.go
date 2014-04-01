package main

import (
	// "fmt"
	"net"
	"net/http"
	"os"

	"github.com/innotech/hydra/config"
	"github.com/innotech/hydra/database/connector"
	"github.com/innotech/hydra/etcd"
	"github.com/innotech/hydra/load_balancer"
	"github.com/innotech/hydra/log"
	"github.com/innotech/hydra/server"
)

func main() {
	// Load configuration.
	var conf = config.New()
	if err := conf.Load(os.Args[1:]); err != nil {
		log.Fatal(err.Error() + "\n")
	}

	if conf.DataDir == "" {
		log.Fatal("Data dir does't exist")
	}

	// Create data directory if it doesn't already exist.
	if err := os.MkdirAll(conf.DataDir, 0744); err != nil {
		log.Fatalf("Unable to create path: %s", err)
	}

	// Load etcd configuration.
	if err := conf.LoadEtcdConfig(); err != nil {
		log.Fatalf("Unable to load etcd conf: %s", err)
	}

	// connector.SetEtcdConnector(etcd)

	

	// Launch services
	var etcd = etcd.New(conf.EtcdConf)
	etcd.Load()
	hydraEnv := os.Getenv("HYDRA_ENV")
	if hydraEnv == "ETCD_TEST" {
		etcd.Start(hydraEnv)
	} else {
		go func() {
			etcd.Start(hydraEnv)
		}()

		connector.SetEtcdConnector(etcd)

		// Private Server API
		privateHydraListener, err := net.Listen("tcp", conf.PrivateAddr)
		if err != nil {
			log.Fatalf("Failed to create hydra private listener: ", err)
		}
		var privateServer = server.NewPrivateServer(privateHydraListener)
		privateServer.RegisterHandlers()
		go func() {
			log.Infof("hydra private server [name %s, listen on %s, advertised url %s]", conf.Name, conf.PrivateAddr, "http://"+conf.PrivateAddr)
			log.Fatal(http.Serve(privateServer.Listener, privateServer.Router))
		}()

		// Public Server API
		const loadBalancerFrontendEndpoint string = "ipc://frontend.ipc"
		publicHydraListener, err := net.Listen("tcp", conf.PublicAddr)
		if err != nil {
			log.Fatalf("Failed to create hydra public listener: ", err)
		}
		var publicServer = server.NewPublicServer(publicHydraListener, loadBalancerFrontendEndpoint)
		publicServer.RegisterHandlers()
		go func() {
			log.Infof("hydra public server [name %s, listen on %s, advertised url %s]", conf.Name, conf.PublicAddr, "http://"+conf.PublicAddr)
			log.Fatal(http.Serve(publicServer.Listener, publicServer.Router))
		}()

		// Load applications.
		var appsConfig = config.NewApplicationsConfig()
		if _, err := os.Stat(conf.AppsFile); os.IsNotExist(err) {
			log.Warnf("Unable to find apps file: %s", err)
		} else {
			if err := appsConfig.Load(conf.AppsFile); err != nil {
				log.Fatalf("Unable to load applications: %s", err)
			}
		}
		
		// Persist Configured applications
		if err := appsConfig.Persists(); err != nil {
			log.Fatalf("Failed to save configured applications: ", err)
		}
		
		// Load Balancer
		loadBalancer := load_balancer.NewLoadBalancer(loadBalancerFrontendEndpoint, "tcp://"+conf.LoadBalancerAddr)
		defer loadBalancer.Close()
		log.Info("PRE LOAD BALANCER")
		loadBalancer.Run()
		log.Info("POST LOAD BALANCER")
	}
}
