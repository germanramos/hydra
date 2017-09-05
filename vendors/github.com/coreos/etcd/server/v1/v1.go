package v1

import (
	"github.com/innotech/hydra/vendors/github.com/coreos/etcd/store"
	"github.com/innotech/hydra/vendors/github.com/coreos/etcd/third_party/github.com/coreos/raft"
	"net/http"
)

// The Server interface provides all the methods required for the v1 API.
type Server interface {
	CommitIndex() uint64
	Term() uint64
	Store() store.Store
	Dispatch(raft.Command, http.ResponseWriter, *http.Request) error
}
