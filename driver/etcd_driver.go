package driver

import (
	"time"

	// uhttp "github.com/innotech/hydra/vendors/github.com/coreos/etcd/pkg/http"
	"github.com/innotech/hydra/vendors/github.com/coreos/etcd/server"
	"github.com/innotech/hydra/vendors/github.com/coreos/etcd/third_party/github.com/coreos/raft"
)

type EtcdDriver struct {
	// store *store.Store
	etcd       *server.Server
	peerServer *server.PeerServer
}

func NewEtcdDriver(etcd *server.Server, peerServer *server.PeerServer) *EtcdDriver {
	e := new(EtcdDriver)
	e.etcd = etcd
	e.peerServer = peerServer
	return e
}

// func (ed *EtcdDriver) Create(key string, dir bool, value string, expireTime time.Time, unique bool) error {
// 	c := e.store.CommandFactory().CreateCreateCommand(key, dir, value, expireTime, true)
// 	return e.dispatch(c)
// }

// func (ed *EtcdDriver) Delete(key string, dir, recursive bool) {
// 	c := s.Store().CommandFactory().CreateDeleteCommand(key, dir, recursive)
// 	return e.dispatch(c)
// }

// func (ed *EtcdDriver) Get(key string, recursive, sort bool) {
// 	event, err := s.Store().Get(key, recursive, sort)
// 	if err != nil {
// 		return err
// 	}

// 	b, _ := json.Marshal(event)

// 	return b
// }

func (e *EtcdDriver) Set(key string, dir bool, value string, expireTime time.Time) error {
	c := e.etcd.Store().CommandFactory().CreateSetCommand(key, dir, value, expireTime)
	return e.dispatch(c)
}

// func (ed *EtcdDriver) Update() error {

// }

func (e *EtcdDriver) dispatch(c raft.Command) error {
	ps := e.peerServer
	if ps.RaftServer().State() == raft.Leader {
		result, err := ps.RaftServer().Do(c)
		if err != nil {
			return err
		}

		if result == nil {
			return nil
			// TODO: return Error
			// return etcdErr.NewError(300, "Empty result from raft", s.Store().Index())
		}

		// response for raft related commands[join/remove]
		//TODO:
		// if b, ok := result.([]byte); ok {
		// 	// w.WriteHeader(http.StatusOK)
		// 	// w.Write(b)
		// 	return nil
		// }

		// var b []byte
		// if strings.HasPrefix(req.URL.Path, "/v1") {
		// 	b, _ = json.Marshal(result.(*store.Event).Response(0))
		// 	w.WriteHeader(http.StatusOK)
		// } else {
		// 	e, _ := result.(*store.Event)
		// 	b, _ = json.Marshal(e)

		// 	w.Header().Set("Content-Type", "application/json")
		// 	// etcd index should be the same as the event index
		// 	// which is also the last modified index of the node
		// 	w.Header().Add("X-Etcd-Index", fmt.Sprint(e.Index()))
		// 	w.Header().Add("X-Raft-Index", fmt.Sprint(s.CommitIndex()))
		// 	w.Header().Add("X-Raft-Term", fmt.Sprint(s.Term()))

		// 	if e.IsCreated() {
		// 		w.WriteHeader(http.StatusCreated)
		// 	} else {
		// 		w.WriteHeader(http.StatusOK)
		// 	}
		// }

		// w.Write(b)

		return nil

	} else {
		leader := ps.RaftServer().Leader()

		// No leader available.
		if leader == "" {
			// TODO: return error
			return nil
			// return etcdErr.NewError(300, "", s.Store().Index())
		}

		// var url string
		// TODO:
		// switch c.(type) {
		// case *server.JoinCommand, *server.RemoveCommand:
		// 	url, _ = ps.registry.PeerURL(leader)
		// default:
		// 	url, _ = ps.registry.ClientURL(leader)
		// }
		// uhttp.Redirect(url, w, req)

		return nil
	}
}
