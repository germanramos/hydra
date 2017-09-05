package v2

import (
	"time"

	"github.com/innotech/hydra/vendors/github.com/coreos/etcd/log"
	"github.com/innotech/hydra/vendors/github.com/coreos/etcd/store"
	"github.com/innotech/hydra/vendors/github.com/coreos/etcd/third_party/github.com/coreos/raft"
)

func init() {
	raft.RegisterCommand(&CreateCommand{})
}

// Create command
type CreateCommand struct {
	Key		string		`json:"key"`
	Value		string		`json:"value"`
	ExpireTime	time.Time	`json:"expireTime"`
	Unique		bool		`json:"unique"`
	Dir		bool		`json:"dir"`
}

// The name of the create command in the log
func (c *CreateCommand) CommandName() string {
	return "etcd:create"
}

// Create node
func (c *CreateCommand) Apply(context raft.Context) (interface{}, error) {
	s, _ := context.Server().StateMachine().(store.Store)

	e, err := s.Create(c.Key, c.Dir, c.Value, c.Unique, c.ExpireTime)

	if err != nil {
		log.Debug(err)
		return nil, err
	}

	return e, nil
}
