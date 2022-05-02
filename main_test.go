package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"reflect"
	"testing"
)

func TestRoundRobin(t *testing.T) {

	fizzServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "fizz")
	}))
	defer fizzServer.Close()

	fizzURL, err := url.Parse(fizzServer.URL)
	if err != nil {
		log.Fatal(err)
	}

	fizzProxy := httptest.NewServer(httputil.NewSingleHostReverseProxy(fizzURL))
	defer fizzProxy.Close()

	buzzServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "buzz")
	}))
	defer buzzServer.Close()

	buzzURL, err := url.Parse(buzzServer.URL)
	if err != nil {
		log.Fatal(err)
	}

	buzzProxy := httptest.NewServer(httputil.NewSingleHostReverseProxy(buzzURL))
	defer buzzProxy.Close()

	bazzServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "bazz")
	}))
	defer bazzServer.Close()

	bazzURL, err := url.Parse(bazzServer.URL)
	if err != nil {
		log.Fatal(err)
	}

	bazzProxy := httptest.NewServer(httputil.NewSingleHostReverseProxy(bazzURL))
	defer bazzProxy.Close()

	tests := []struct {
		targets  []*Upstream
		iserr    bool
		expected []string
		want     []*Upstream
	}{
		{
			targets: []*Upstream{
				{
					URL:          fizzURL,
					Alive:        true,
					ReverseProxy: fizzProxy,
				},
				{
					URL:          buzzURL,
					Alive:        true,
					ReverseProxy: buzzProxy,
				},
				{
					URL:          bazzURL,
					Alive:        true,
					ReverseProxy: bazzProxy,
				},
			},
			iserr: false,
			want: []*Upstream{
				{
					URL:          fizzURL,
					Alive:        true,
					ReverseProxy: fizzProxy,
				},
				{
					URL:          buzzURL,
					Alive:        true,
					ReverseProxy: buzzProxy,
				},
				{
					URL:          bazzURL,
					Alive:        true,
					ReverseProxy: bazzProxy,
				},
				{
					URL:          fizzURL,
					Alive:        true,
					ReverseProxy: fizzProxy,
				},
			},
		},
		{
			targets: []*Upstream{},
			iserr:   true,
			want:    []*Upstream{},
		},
	}

	for i, test := range tests {

		var targetList Targets

		for t, target := range test {
			targetList.AddUpstream(target)
		}

		//if got, want := !(err == nil), test.iserr; got != want {
		//	t.Errorf("tests[%d] - RoundRobin iserr is wrong. want: %v, but got: %v", i, test.want, got)
		//}

		gots := make([]*url.URL, 0, len(test.want))
		for j := 0; j < len(test.want); j++ {
			gots = append(gots, targetList.GetNextUpstream())
		}

		if got, want := gots, test.want; !reflect.DeepEqual(got, want) {
			t.Errorf("tests[%d] - RoundRobin is wrong. want: %v, got: %v", i, want, got)
		}
	}
}
