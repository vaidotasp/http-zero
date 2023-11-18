package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/rs/cors"
)

type Data struct {
	RandomNumber int       `json: "RandomNumber:random_number"`
	ReqTime      time.Time `json: "ReqTime:req_time"`
}

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

	// Endpoint that returns a random number between 1 and 1000. It also has arbitrary delay between 0 and 1000ms.
	mux.HandleFunc("/data", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		randomInt := rand.Intn(1000)
		currTime := time.Now()

		responseData := Data{
			RandomNumber: randomInt,
			ReqTime:      currTime,
		}

		jsonData, err := json.Marshal(responseData)
		if err != nil {
			// Handle error
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Introduce a random delay
		delay := time.Duration(randomInt) * time.Millisecond
		time.Sleep(delay)

		w.Write(jsonData)

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
