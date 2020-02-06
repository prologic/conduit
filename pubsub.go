package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"

	log "github.com/sirupsen/logrus"
)

type pubsub struct {
	sync.RWMutex
	subscribers map[string][]chan []byte
}

func newPubSub() *pubsub {
	return &pubsub{
		subscribers: make(map[string][]chan []byte),
	}
}

func (p *pubsub) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		topic := r.URL.Path

		ch := make(chan []byte, 1)

		p.Lock()
		p.subscribers[topic] = append(p.subscribers[topic], ch)
		p.Unlock()

		if r.URL.Query().Get("persist") == "true" {
			for {
				w.Write(<-ch)
				w.(http.Flusher).Flush()
			}
		} else {
			w.Write(<-ch)
		}
	} else if r.Method == "POST" {
		topic := r.URL.Path

		p.RLock()
		defer p.RUnlock()

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			msg := fmt.Sprintf("error reading request body: %s", err)
			log.WithError(err).Error(msg)
			fmt.Fprint(w, msg)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		for _, ch := range p.subscribers[topic] {
			go func(ch chan []byte) {
				ch <- body
			}(ch)
		}
	} else {
		http.Error(w, "Unsupported Method", http.StatusMethodNotAllowed)
	}
}
