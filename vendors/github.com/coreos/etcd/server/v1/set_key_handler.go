package v1

import (
	"net/http"

	etcdErr "github.com/innotech/hydra/vendors/github.com/coreos/etcd/error"
	"github.com/innotech/hydra/vendors/github.com/coreos/etcd/store"
	"github.com/innotech/hydra/vendors/github.com/coreos/etcd/third_party/github.com/coreos/raft"
	"github.com/innotech/hydra/vendors/github.com/coreos/etcd/third_party/github.com/gorilla/mux"
)

// Sets the value for a given key.
func SetKeyHandler(w http.ResponseWriter, req *http.Request, s Server) error {
	vars := mux.Vars(req)
	key := "/" + vars["key"]

	req.ParseForm()

	// Parse non-blank value.
	value := req.Form.Get("value")
	if len(value) == 0 {
		return etcdErr.NewError(200, "Set", s.Store().Index())
	}

	// Convert time-to-live to an expiration time.
	expireTime, err := store.TTL(req.Form.Get("ttl"))
	if err != nil {
		return etcdErr.NewError(202, "Set", s.Store().Index())
	}

	// If the "prevValue" is specified then test-and-set. Otherwise create a new key.
	var c raft.Command
	if prevValueArr, ok := req.Form["prevValue"]; ok {
		if len(prevValueArr[0]) > 0 {
			// test against previous value
			c = s.Store().CommandFactory().CreateCompareAndSwapCommand(key, value, prevValueArr[0], 0, expireTime)
		} else {
			// test against existence
			c = s.Store().CommandFactory().CreateCreateCommand(key, false, value, expireTime, false)
		}

	} else {
		c = s.Store().CommandFactory().CreateSetCommand(key, false, value, expireTime)
	}

	return s.Dispatch(c, w, req)
}
