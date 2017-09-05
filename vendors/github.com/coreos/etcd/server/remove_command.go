package server

import (
	"encoding/binary"
	"os"

	"github.com/innotech/hydra/vendors/github.com/coreos/etcd/log"
	"github.com/innotech/hydra/vendors/github.com/coreos/etcd/third_party/github.com/coreos/raft"
)

func init() {
	raft.RegisterCommand(&RemoveCommand{})
}

// The RemoveCommand removes a server from the cluster.
type RemoveCommand struct {
	Name string `json:"name"`
}

// The name of the remove command in the log
func (c *RemoveCommand) CommandName() string {
	return "etcd:remove"
}

// Remove a server from the cluster
func (c *RemoveCommand) Apply(context raft.Context) (interface{}, error) {
	ps, _ := context.Server().Context().(*PeerServer)

	// Remove node from the shared registry.
	err := ps.registry.Unregister(c.Name)

	// Delete from stats
	delete(ps.followersStats.Followers, c.Name)

	if err != nil {
		log.Debugf("Error while unregistering: %s (%v)", c.Name, err)
		return []byte{0}, err
	}

	// Remove peer in raft
	err = context.Server().RemovePeer(c.Name)
	if err != nil {
		log.Debugf("Unable to remove peer: %s (%v)", c.Name, err)
		return []byte{0}, err
	}

	if c.Name == context.Server().Name() {
		// the removed node is this node

		// if the node is not replaying the previous logs
		// and the node has sent out a join request in this
		// start. It is sure that this node received a new remove
		// command and need to be removed
		if context.CommitIndex() > ps.joinIndex && ps.joinIndex != 0 {
			log.Debugf("server [%s] is removed", context.Server().Name())
			os.Exit(0)
		} else {
			// else ignore remove
			log.Debugf("ignore previous remove command.")
		}
	}

	b := make([]byte, 8)
	binary.PutUvarint(b, context.CommitIndex())

	return b, err
}
