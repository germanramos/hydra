package store

import (
	"testing"
	"time"

	"github.com/innotech/hydra/vendors/github.com/coreos/etcd/third_party/github.com/stretchr/testify/assert"
)

// Ensure that a successful Get is recorded in the stats.
func TestStoreStatsGetSuccess(t *testing.T) {
	s := newStore()
	s.Create("/foo", false, "bar", false, Permanent)
	s.Get("/foo", false, false)
	assert.Equal(t, uint64(1), s.Stats.GetSuccess, "")
}

// Ensure that a failed Get is recorded in the stats.
func TestStoreStatsGetFail(t *testing.T) {
	s := newStore()
	s.Create("/foo", false, "bar", false, Permanent)
	s.Get("/no_such_key", false, false)
	assert.Equal(t, uint64(1), s.Stats.GetFail, "")
}

// Ensure that a successful Create is recorded in the stats.
func TestStoreStatsCreateSuccess(t *testing.T) {
	s := newStore()
	s.Create("/foo", false, "bar", false, Permanent)
	assert.Equal(t, uint64(1), s.Stats.CreateSuccess, "")
}

// Ensure that a failed Create is recorded in the stats.
func TestStoreStatsCreateFail(t *testing.T) {
	s := newStore()
	s.Create("/foo", true, "", false, Permanent)
	s.Create("/foo", false, "bar", false, Permanent)
	assert.Equal(t, uint64(1), s.Stats.CreateFail, "")
}

// Ensure that a successful Update is recorded in the stats.
func TestStoreStatsUpdateSuccess(t *testing.T) {
	s := newStore()
	s.Create("/foo", false, "bar", false, Permanent)
	s.Update("/foo", "baz", Permanent)
	assert.Equal(t, uint64(1), s.Stats.UpdateSuccess, "")
}

// Ensure that a failed Update is recorded in the stats.
func TestStoreStatsUpdateFail(t *testing.T) {
	s := newStore()
	s.Update("/foo", "bar", Permanent)
	assert.Equal(t, uint64(1), s.Stats.UpdateFail, "")
}

// Ensure that a successful CAS is recorded in the stats.
func TestStoreStatsCompareAndSwapSuccess(t *testing.T) {
	s := newStore()
	s.Create("/foo", false, "bar", false, Permanent)
	s.CompareAndSwap("/foo", "bar", 0, "baz", Permanent)
	assert.Equal(t, uint64(1), s.Stats.CompareAndSwapSuccess, "")
}

// Ensure that a failed CAS is recorded in the stats.
func TestStoreStatsCompareAndSwapFail(t *testing.T) {
	s := newStore()
	s.Create("/foo", false, "bar", false, Permanent)
	s.CompareAndSwap("/foo", "wrong_value", 0, "baz", Permanent)
	assert.Equal(t, uint64(1), s.Stats.CompareAndSwapFail, "")
}

// Ensure that a successful Delete is recorded in the stats.
func TestStoreStatsDeleteSuccess(t *testing.T) {
	s := newStore()
	s.Create("/foo", false, "bar", false, Permanent)
	s.Delete("/foo", false, false)
	assert.Equal(t, uint64(1), s.Stats.DeleteSuccess, "")
}

// Ensure that a failed Delete is recorded in the stats.
func TestStoreStatsDeleteFail(t *testing.T) {
	s := newStore()
	s.Delete("/foo", false, false)
	assert.Equal(t, uint64(1), s.Stats.DeleteFail, "")
}

//Ensure that the number of expirations is recorded in the stats.
func TestStoreStatsExpireCount(t *testing.T) {
	s := newStore()

	c := make(chan bool)
	defer func() {
		c <- true
	}()

	go mockSyncService(s.DeleteExpiredKeys, c)
	s.Create("/foo", false, "bar", false, time.Now().Add(500*time.Millisecond))
	assert.Equal(t, uint64(0), s.Stats.ExpireCount, "")
	time.Sleep(600 * time.Millisecond)
	assert.Equal(t, uint64(1), s.Stats.ExpireCount, "")
}
