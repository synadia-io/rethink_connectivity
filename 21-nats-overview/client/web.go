package client

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/nats-io/nats.go"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func RunWeb() {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}
		defer conn.Close()

		sub, err := nc.Subscribe("chat", func(msg *nats.Msg) {
			conn.WriteMessage(websocket.TextMessage, msg.Data)
		})
		if err != nil {
			log.Println(err)
			return
		}
		defer sub.Unsubscribe()

		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Println(err)
				break
			}
			nc.Publish("chat", message)
		}
	})

	http.Handle("/", http.FileServer(http.Dir("./static")))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
