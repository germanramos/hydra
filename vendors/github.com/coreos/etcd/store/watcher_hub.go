package store

import (
	"container/list"
	"path"
	"strings"
	"sync"
	"sync/atomic"

	etcdErr "github.com/innotech/hydra/vendors/github.com/coreos/etcd/error"
)

// A watcherHub contains all subscribed watchers
// watchers is a map with watched path as key and watcher as value
// EventHistory keeps the old events for watcherHub. It is used to help
// watcher to get a continuous event history. Or a watcher might miss the
// event happens between the end of the first watch command and the start
// of the second command.
type watcherHub struct {
	mutex		sync.Mutex	// protect the hash map
	watchers	map[string]*list.List
	count		int64	// current number of watchers.
	EventHistory	*EventHistory
}

// newWatchHub creates a watchHub. The capacity determines how many events we will
// keep in the eventHistory.
// Typically, we only need to keep a small size of history[smaller than 20K].
// Ideally, it should smaller than 20K/s[max throughput] * 2 * 50ms[RTT] = 2000
func newWatchHub(capacity int) *watcherHub {
	return &watcherHub{
		watchers:	make(map[string]*list.List),
		EventHistory:	newEventHistory(capacity),
	}
}

// Watch function returns a watcher.
// If recursive is true, the first change after index under key will be sent to the event channel of the watcher.
// If recursive is false, the first change after index at key will be sent to the event channel of the watcher.
// If index is zero, watch will start from the current index + 1.
func (wh *watcherHub) watch(key string, recursive, stream bool, index uint64) (*Watcher, *etcdErr.Error) {
	event, err := wh.EventHistory.scan(key, recursive, index)

	if err != nil {
		return nil, err
	}

	w := &Watcher{
		EventChan:	make(chan *Event, 1),	// use a buffered channel
		recursive:	recursive,
		stream:		stream,
		sinceIndex:	index,
	}

	if event != nil {
		w.EventChan <- event
		return w, nil
	}

	wh.mutex.Lock()
	defer wh.mutex.Unlock()

	l, ok := wh.watchers[key]

	var elem *list.Element

	if ok {	// add the new watcher to the back of the list
		elem = l.PushBack(w)

	} else {	// create a new list and add the new watcher
		l = list.New()
		elem = l.PushBack(w)
		wh.watchers[key] = l
	}

	w.remove = func() {
		if w.removed {	// avoid remove it twice
			return
		}

		wh.mutex.Lock()
		defer wh.mutex.Unlock()

		w.removed = true
		l.Remove(elem)
		atomic.AddInt64(&wh.count, -1)
		if l.Len() == 0 {
			delete(wh.watchers, key)
		}
	}

	atomic.AddInt64(&wh.count, 1)

	return w, nil
}

// notify function accepts an event and notify to the watchers.
func (wh *watcherHub) notify(e *Event) {
	e = wh.EventHistory.addEvent(e)	// add event into the eventHistory

	segments := strings.Split(e.Node.Key, "/")

	currPath := "/"

	// walk through all the segments of the path and notify the watchers
	// if the path is "/foo/bar", it will notify watchers with path "/",
	// "/foo" and "/foo/bar"

	for _, segment := range segments {
		currPath = path.Join(currPath, segment)
		// notify the watchers who interests in the changes of current path
		wh.notifyWatchers(e, currPath, false)
	}
}

func (wh *watcherHub) notifyWatchers(e *Event, nodePath string, deleted bool) {
	wh.mutex.Lock()
	defer wh.mutex.Unlock()

	l, ok := wh.watchers[nodePath]
	if ok {
		curr := l.Front()

		for curr != nil {
			next := curr.Next()	// save reference to the next one in the list

			w, _ := curr.Value.(*Watcher)

			originalPath := (e.Node.Key == nodePath)
			if (originalPath || !isHidden(nodePath, e.Node.Key)) && w.notify(e, originalPath, deleted) {
				if !w.stream {	// do not remove the stream watcher
					// if we successfully notify a watcher
					// we need to remove the watcher from the list
					// and decrease the counter
					l.Remove(curr)
					atomic.AddInt64(&wh.count, -1)
				}
			}

			curr = next	// update current to the next element in the list
		}

		if l.Len() == 0 {
			// if we have notified all watcher in the list
			// we can delete the list
			delete(wh.watchers, nodePath)
		}
	}
}

// clone function clones the watcherHub and return the cloned one.
// only clone the static content. do not clone the current watchers.
func (wh *watcherHub) clone() *watcherHub {
	clonedHistory := wh.EventHistory.clone()

	return &watcherHub{
		EventHistory: clonedHistory,
	}
}

// isHidden checks to see if key path is considered hidden to watch path i.e. the
// last element is hidden or it's within a hidden directory
func isHidden(watchPath, keyPath string) bool {
	// When deleting a directory, watchPath might be deeper than the actual keyPath
	// For example, when deleting /foo we also need to notify watchers on /foo/bar.
	if len(watchPath) > len(keyPath) {
		return false
	}
	// if watch path is just a "/", after path will start without "/"
	// add a "/" to deal with the special case when watchPath is "/"
	afterPath := path.Clean("/" + keyPath[len(watchPath):])
	return strings.Contains(afterPath, "/_")
}
