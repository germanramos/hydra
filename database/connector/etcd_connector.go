package connector

import (
	"github.com/innotech/hydra/etcd"
	"log"
	"net/http"
	"strconv"

	etcdErr "github.com/innotech/hydra/vendors/github.com/coreos/etcd/error"
	uhttp "github.com/innotech/hydra/vendors/github.com/coreos/etcd/pkg/http"
	. "github.com/innotech/hydra/vendors/github.com/coreos/etcd/server"
	"github.com/innotech/hydra/vendors/github.com/coreos/etcd/store"
	"github.com/innotech/hydra/vendors/github.com/coreos/etcd/third_party/github.com/coreos/raft"
)

type EtcdDriver interface {
	Delete(key string, dir, recursive bool) error
	Get(key string, recursive bool, sort bool) []interface{}
	Set(key string, dir bool, value string, ttl string, w http.ResponseWriter, req *http.Request) error
}

type EtcdConnector struct {
	etcd *etcd.Etcd
}

var e *EtcdConnector

func SetEtcdConnector(etcd *etcd.Etcd) {
	e = new(EtcdConnector)
	e.etcd = etcd
}

func GetEtcdConnector() *EtcdConnector {
	return e
}

func (e EtcdConnector) Delete(key string, dir, recursive bool) error {
	return nil

}

func (e EtcdConnector) Get(key string, recursive bool, sort bool) (*store.Event, error) {
	return e.etcd.EtcdServer.Store().Get(key, recursive, sort)
}

func (e EtcdConnector) CreateSetCommand(key string, dir bool, value string, ttl string) (raft.Command, error) {
	expireTime, err := store.TTL(ttl)
	if err != nil {
		return nil, etcdErr.NewError(etcdErr.EcodeTTLNaN, "Create", e.etcd.EtcdServer.Store().Index())
	}
	c := e.etcd.EtcdServer.Store().CommandFactory().CreateSetCommand(key, dir, value, expireTime)
	return c, nil
}

func (e EtcdConnector) Set(key string, dir bool, value string, ttl string, w http.ResponseWriter, req *http.Request) error {
	ps := e.etcd.PeerServer
	registry := e.etcd.Registry
	if ps.RaftServer().State() == raft.Leader {
		expireTime, err := store.TTL(ttl)
		if err != nil {
			return etcdErr.NewError(etcdErr.EcodeTTLNaN, "Create", e.etcd.EtcdServer.Store().Index())
		}
		c := e.etcd.EtcdServer.Store().CommandFactory().CreateSetCommand(key, dir, value, expireTime)
		result, err := ps.RaftServer().Do(c)
		if err != nil {
			return err
		}

		if result == nil {
			return etcdErr.NewError(300, "Empty result from raft", e.etcd.EtcdServer.Store().Index())
		}

		if w != nil {
			w.WriteHeader(http.StatusOK)
		}

		return nil
	} else {
		leader := ps.RaftServer().Leader()

		// No leader available.
		if leader == "" {
			return etcdErr.NewError(300, "", e.etcd.EtcdServer.Store().Index())
		}

		leaderUrl, _ := registry.ClientURL(leader)

		client := http.DefaultClient
		v := "?"
		if !dir {
			v += "value=" + value + "&ttl=" + ttl
		} else {
			v += "dir=" + strconv.FormatBool(dir) + "&ttl=" + ttl
		}

		req, _ := http.NewRequest("PUT", leaderUrl+"/v2/keys"+key+v, nil)
		resp, err := client.Do(req)
		resp.Body.Close()

		if err != nil {
			log.Println(err)
			return err
		} else {
			if w != nil {
				w.WriteHeader(http.StatusOK)
			}

			return nil
		}
	}

	return nil
}

func (e EtcdConnector) Dispatch(c raft.Command, w http.ResponseWriter, req *http.Request) error {
	ps := e.etcd.PeerServer
	registry := e.etcd.Registry
	if ps.RaftServer().State() == raft.Leader {
		result, err := ps.RaftServer().Do(c)
		if err != nil {
			return err
		}

		if result == nil {
			return etcdErr.NewError(300, "Empty result from raft", e.etcd.EtcdServer.Store().Index())
		}

		if w != nil {
			w.WriteHeader(http.StatusOK)
		}

		return nil

	} else {
		leader := ps.RaftServer().Leader()

		// No leader available.
		if leader == "" {
			return etcdErr.NewError(300, "", e.etcd.EtcdServer.Store().Index())
		}

		var url string
		switch c.(type) {
		case *JoinCommand, *RemoveCommand:
			url, _ = registry.PeerURL(leader)
		default:
			url, _ = registry.ClientURL(leader)
		}
		uhttp.Redirect(url, w, req)

		return nil
	}
}
