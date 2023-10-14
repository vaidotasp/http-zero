package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/rs/cors"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{\"Server Response\": \"OK\"}"))
	})

	mux.HandleFunc("/socket", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "client/ws.html")
	})

	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}
		log.Println("Client Connected")

		for {
			messageType, p, err := conn.ReadMessage()
			if err != nil {
				log.Fatal()
				log.Println(err)
				return
			}
			clientMsg := string(p)
			log.Println(clientMsg)

			// Client can close server connection by sending "exit"
			if clientMsg == "exit" {
				conn.Close()
			}

			uppered := strings.ToUpper(clientMsg)

			if err := conn.WriteMessage(messageType, []byte(uppered)); err != nil {
				log.Println(err)
				return
			}
		}
	})

	handler := cors.Default().Handler(mux)
	http.ListenAndServe("localhost:8080", handler)
}
