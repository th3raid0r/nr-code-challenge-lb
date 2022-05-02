package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/th3raid0r/nr-code-challenge-lb/util"
)

const (
	Attempts int = iota
	Retry
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

func isUpstreamAlive(u *url.URL) bool {
	timeout := 2 * time.Second
	conn, err := net.DialTimeout("tcp", u.Host, timeout)
	if err != nil {
		log.Println("Site unreachable, error: ", err)
		return false
	}
	_ = conn.Close()
	return true
}

func (s *ServerPool) HealthCheck() {
	for _, b := range s.backends {
		status := "up"
		alive := isUpstreamAlive(b.URL)
		b.SetAlive(alive)
		if !alive {
			status = "down"
		}
		log.Printf("%s [%s]\n", b.URL, status)
	}
}

func healthCheck() {
	t := time.NewTicker(time.Second * 20)
	for {
		select {
		case <-t.C:
			log.Println("Starting health check...")
			serverPool.HealthCheck()
			log.Println("Health check completed")
		}
	}
}

// GetAttemptsFromContext returns the attempts for request
func GetAttemptsFromContext(r *http.Request) int {
	if attempts, ok := r.Context().Value(Attempts).(int); ok {
		return attempts
	}
	return 1
}

// GetAttemptsFromContext returns the attempts for request
func GetRetryFromContext(r *http.Request) int {
	if retry, ok := r.Context().Value(Retry).(int); ok {
		return retry
	}
	return 0
}

func lb(w http.ResponseWriter, r *http.Request) {
	attempts := GetAttemptsFromContext(r)
	if attempts > 3 {
		log.Printf("%s(%s) Max attempts reached, terminating\n", r.RemoteAddr, r.URL.Path)
		http.Error(w, "Service not available", http.StatusServiceUnavailable)
		return
	}

	upstream := Targets.GetNextUpstream()
	if upstream != nil {
		upstream.ReverseProxy.ServeHTTP(w, r)
		return
	}
	http.Error(w, "Service not available", http.StatusServiceUnavailable)
}

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	var targets string = config.TargetList
	var port int = config.Port

	if len(targets) == 0 {
		log.Fatal("Please provide one or more upstream servers to the load balancer")
	}

	tokens := strings.Split(targets, ",")
	for _, tok := range tokens {
		serverUrl, err := url.Parse(tok)
		if err != nill {
			log.Fatal(err)
		}

		proxy := httputil.NewSingleHostReverseProxy(serverUrl)
		proxy.ErrorHandler = func(writer http.ResponseWriter, request *http.Request, e error) {
			log.Printf("[%s] %s\n", serverUrl.Host, e.Error())
			retries := GetRetryFromContext(request)
			if retries < 3 {
				select {
				case <-time.After(10 * time.Millisecond):
					ctx := context.WithValue(request.Context(), Retry, retries+1)
					proxy.ServeHTTP(writer, request.WithContext(ctx))
				}
				return
			}

			// after 3 retries, mark this backend as down
			serverPool.MarkBackendStatus(serverUrl, false)

			// if the same request routing for few attempts with different backends, increase the count
			attempts := GetAttemptsFromContext(request)
			log.Printf("%s(%s) Attempting retry %d\n", request.RemoteAddr, request.URL.Path, attempts)
			ctx := context.WithValue(request.Context(), Attempts, attempts+1)
			lb(writer, request.WithContext(ctx))
		}

		serverPool.AddBackend(&Upstream{
			URL:          serverUrl,
			Alive:        true,
			ReverseProxy: proxy,
		})
		log.Printf("Configured server: %s\n", serverUrl)
	}

	// create http server
	server := http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: http.HandlerFunc(lb),
	}

	// start health checking
	go healthCheck()

	log.Printf("Load Balancer started at :%d\n", port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}

}
