package connector

import (
	// "encoding/json"
	// "fmt"
	"github.com/innotech/hydra/etcd"
	"log"
	"net/http"
	// "net/url"
	"strconv"
	// "strings"
	// "time"

	// . "github.com/innotech/hydra/vendors/github.com/coreos/etcd/server"
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

// func (e EtcdConnector) Set(key string, dir bool, value string, ttl string, w http.ResponseWriter, req *http.Request) error {
// 	expireTime, err := store.TTL(ttl)
// 	if err != nil {
// 		return etcdErr.NewError(etcdErr.EcodeTTLNaN, "Create", e.etcd.EtcdServer.Store().Index())
// 	}
// 	c := e.etcd.EtcdServer.Store().CommandFactory().CreateSetCommand(key, dir, value, expireTime)
// 	// c := e.etcd.Store().CommandFactory().CreateSetCommand(key, dir, value, expireTime)
// 	// return e.dispatch(c)
// 	dispatchErr := e.Dispatch(c, w, req)
// 	if dispatchErr != nil {
// 		log.Println("DISPATCH ERROR")
// 		log.Println(dispatchErr)
// 		return dispatchErr
// 	}
// 	// event, err := e.etcd.EtcdServer.Store().Get("/v2/keys/testapp/Number", true, false)
// 	// event, err := e.etcd.EtcdServer.Store().Get("/", true, false)
// 	// if err != nil {
// 	// 	return err
// 	// }
// 	// fmt.Printf("%+v\n", event)
// 	// fmt.Printf("%#v\n", event)
// 	// js, _ := json.Marshal(event)
// 	// fmt.Println(string(js))
// 	// inter := event.Response(e.etcd.EtcdServer.Store().Index())
// 	// fmt.Printf("%+v\n", inter)
// 	// fmt.Printf("%#v\n", inter)
// 	// js2, err := json.Marshal(inter)
// 	// fmt.Println(string(js2))
// 	// return err

// 	return nil
// }

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
			// return nil
		}

		if w != nil {
			// log.Println("Pre WriteHeader")
			w.WriteHeader(http.StatusOK)
			// log.Println("Post WriteHeader")
		}

		return nil
	} else {
		leader := ps.RaftServer().Leader()

		// No leader available.
		if leader == "" {
			return etcdErr.NewError(300, "", e.etcd.EtcdServer.Store().Index())
		}

		leaderUrl, _ := registry.ClientURL(leader)
		// log.Println("CLIENT POST to: " + leaderUrl + key)
		// client := &http.Client{
		// 	CheckRedirect: redirectPolicyFunc,
		// }

		client := http.DefaultClient
		// resp, err := client.Get("http://example.com")
		// ...

		// v := url.Values{}
		v := "?"
		if !dir {
			v += "value=" + value + "&ttl=" + ttl
			// v.Set("value", value)
			// v.Add("ttl", ttl)
			// v := url.Values{"value": {value}, "ttl": {ttl}}
		} else {
			v += "dir=" + strconv.FormatBool(dir) + "&ttl=" + ttl
			// v.Set("dir", strconv.FormatBool(dir))
			// v.Add("ttl", ttl)
			// v := url.Values{"dir": {strconv.FormatBool(dir)}, "ttl": {ttl}}
		}

		log.Printf("%#v", v)
		// req, _ := http.NewRequest("PUT", leaderUrl+"/v2/keys"+key, strings.NewReader(v.Encode()))
		req, _ := http.NewRequest("PUT", leaderUrl+"/v2/keys"+key+v, nil)
		// log.Println("CLIENT POST to: " + leaderUrl + "/v2/keys" + key)
		// ...
		// req.Header.Add("If-None-Match", `W/"wyzzy"`)
		resp, err := client.Do(req)
		log.Printf("%#v", resp)
		// _, err := http.PostForm(leaderUrl+key,
		// 	url.Values{"dir": {strconv.FormatBool(dir)}, "value": {value}, "ttl": {ttl}})

		if err != nil {
			log.Println(err)
			return err
		} else {
			// log.Println("Pre WriteHeader 2")
			if w != nil {
				w.WriteHeader(http.StatusOK)
			}
			// log.Println("Post WriteHeader 2")

			return nil
		}
	}
	// dispatchErr := e.Dispatch(c, w, req)
	// if dispatchErr != nil {
	// 	log.Println("DISPATCH ERROR")
	// 	log.Println(dispatchErr)
	// 	return dispatchErr
	// }

	return nil
}

// func (e EtcdConnector) sendPostRequest() {
// 	resp, err := http.PostForm(url+key,
// 		url.Values{"dir": {strconv.FormatBool(dir)}, "value": {value}, "ttl": {ttl}})
// }

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
			// return nil
		}

		if w != nil {
			log.Println("Pre WriteHeader")
			w.WriteHeader(http.StatusOK)
			log.Println("Post WriteHeader")
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
		log.Println("PRE REDIRECT to " + url)
		log.Printf("%#v", w)
		log.Printf("%#v", req)
		uhttp.Redirect(url, w, req)

		return nil
	}
}
