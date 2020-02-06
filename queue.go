package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type queue struct {
	channels map[string]chan []byte
}

func newQueue() *queue {
	return &queue{
		channels: make(map[string]chan []byte),
	}
}

func (q *queue) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var (
		ch chan []byte
		ok bool
	)

	ch, ok = q.channels[r.URL.Path]
	if !ok {
		ch = make(chan []byte)
		q.channels[r.URL.Path] = ch
	}

	if r.Method == "GET" {
		if r.URL.Query().Get("persist") == "true" {
			for {
				w.Write(<-ch)
				w.(http.Flusher).Flush()
			}
		} else {
			w.Write(<-ch)
		}
	} else if r.Method == "POST" {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			msg := fmt.Sprintf("error reading request body: %s", err)
			log.WithError(err).Error(msg)
			fmt.Fprint(w, msg)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		ch <- body
	} else {
		http.Error(w, "Unsupported Method", http.StatusMethodNotAllowed)
	}
}
