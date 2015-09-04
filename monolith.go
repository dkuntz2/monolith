package main

import (
	"github.com/dkuntz2/monolith/protocol"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

func main() {
	hub := protocol.NewHub()
	r := mux.NewRouter()
	r.HandleFunc("/socket", func(w http.ResponseWriter, r *http.Request) {
		var upgrader = websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}

		hub.Attach(conn)
	})

	go hub.Run()
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("frontend")))

	log.Println("Starting server on :3000")
	http.ListenAndServe(":3000", r)
}
