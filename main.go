package main

import (
	"net/http"

	log "github.com/sirupsen/logrus"
)

func main() {
	http.Handle("/queue/", newQueue())
	http.Handle("/topic/", newPubSub())
	log.Fatal(http.ListenAndServe(":8000", nil))
}
