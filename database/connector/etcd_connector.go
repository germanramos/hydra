package connector

import (
	// "encoding/json"
	// "fmt"
	"github.com/innotech/hydra/etcd"
	"time"

	// . "github.com/innotech/hydra/vendors/github.com/coreos/etcd/server"
	"github.com/innotech/hydra/vendors/github.com/coreos/etcd/store"
	"github.com/innotech/hydra/vendors/github.com/coreos/etcd/third_party/github.com/coreos/raft"
)

type EtcdDriver interface {
	Delete(key string, dir, recursive bool) error
	Get(key string, recursive bool, sort bool) []interface{}
	Set(key string, dir bool, value string, expireTime time.Time) error
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

func (e EtcdConnector) Set(key string, dir bool, value string, expireTime time.Time) error {
	c := e.etcd.EtcdServer.Store().CommandFactory().CreateSetCommand(key, dir, value, expireTime)
	// c := e.etcd.Store().CommandFactory().CreateSetCommand(key, dir, value, expireTime)
	// return e.dispatch(c)
	err := e.dispatch(c)
	if err != nil {
		return err
	}
	// event, err := e.etcd.EtcdServer.Store().Get("/v2/keys/testapp/Number", true, false)
	// event, err := e.etcd.EtcdServer.Store().Get("/", true, false)
	// if err != nil {
	// 	return err
	// }
	// fmt.Printf("%+v\n", event)
	// fmt.Printf("%#v\n", event)
	// js, _ := json.Marshal(event)
	// fmt.Println(string(js))
	// inter := event.Response(e.etcd.EtcdServer.Store().Index())
	// fmt.Printf("%+v\n", inter)
	// fmt.Printf("%#v\n", inter)
	// js2, err := json.Marshal(inter)
	// fmt.Println(string(js2))
	// return err

	return nil
}

func (e EtcdConnector) dispatch(c raft.Command) error {
	ps := e.etcd.PeerServer
	if ps.RaftServer().State() == raft.Leader {
		result, err := ps.RaftServer().Do(c)
		if err != nil {
			return err
		}

		if result == nil {
			return nil
		}

		return nil

	} else {
		leader := ps.RaftServer().Leader()

		if leader == "" {
			return nil
		}

		// var url string
		// switch c.(type) {
		// case *JoinCommand, *RemoveCommand:
		// 	url, _ = ps.registry.PeerURL(leader)
		// default:
		// 	url, _ = ps.registry.ClientURL(leader)
		// }
		// uhttp.Redirect(url, w, req)

		return nil
	}
}
