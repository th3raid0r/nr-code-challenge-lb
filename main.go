package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"sync/atomic"
)

type Upstream struct {
	// need net/url for url
	URL   *url.URL
	Alive bool
	mux   sync.RWMutex
	//need net/http/httputil
	ReverseProxy *httputil.ReverseProxy
}

type Targets struct {
	upstreams []*Upstream
	current   uint64
}

func (t *Targets) NextIndex() int {
	// Important that this is an incremental update - so import atomic above
	return int(atomic.AddUint64(&t.current, uint64(1)) % uint64(len(t.upstreams)))
}

func (u *Upstream) SetAlive(alive bool) {
	//Use mux lock to avoid race conditions. Let's just hope we don't need to write a ton.
	u.mux.Lock()
	u.Alive = alive
	b.mux.Unlock()
}

func (u *Upstream) IsAlive() (alive bool) {
	u.mux.RLock()
	alive = u.Alive
	b.mux.RUnlock()
	return
}

// Return next active upstream
func (t *Targets) GetNextUpstream() *Upstream {
	// standard loop over targets to find an active one.
	next := t.NextIndex
	l := len(t.upstreams) + next
	for i := next; i < l; i++ {
		idx := i % len(t.upstreams)
		if s.backends[idx].IsAlive() {
			if i != next {
				atomic.StoreUint64(&t.current, uint64(idx))
			}
			return t.upstreams[idx]
		}
	}
	return nill

}

func lb(w http.ResponseWriter, r *http.Request) {
	upstream := Targets.GetNextUpstream()
	if upstream != nil {
		upstream.ReverseProxy.ServeHTTP(w, r)
		return
	}
	http.Error(w, "Service not available", http.StatusServiceUnavailable)
}

func main() {
	fmt.Println("vim-go")
}
