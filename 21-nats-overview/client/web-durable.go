package client

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/nats-io/nats.go"
)

var upgraderJS = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func RunWebDurable() {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		fmt.Println(err)
	}
	defer nc.Close()

  js, err := nc.JetStream()
  if err != nil {
    fmt.Println(err)
  }

  _, err = js.AddStream(&nats.StreamConfig{
    Name: "CHAT",
    Subjects: []string{"chat.>"},
  })
  if err != nil {
    fmt.Println(err)
  }

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgraderJS.Upgrade(w, r, nil)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer conn.Close()

		sub, err := nc.Subscribe("chat.messages", func(msg *nats.Msg) {
			conn.WriteMessage(websocket.TextMessage, msg.Data)
		})
		if err != nil {
			fmt.Println(err)
			return
		}
		defer sub.Unsubscribe()

    // Retrieve chat history
    consumer, err := js.PullSubscribe("chat.messages", "")
    if err != nil {
      fmt.Println(err)
      return
    }

    msgs, err := consumer.Fetch(100)
    if err != nil {
      fmt.Println(err)
    } else {
      for _, msg := range msgs {
        conn.WriteMessage(websocket.TextMessage, msg.Data)
      }
    }

    for {
      _, msg, err := conn.ReadMessage()
      if err != nil {
        fmt.Println(err)
        break
      }
      _, err = js.Publish("chat.messages", msg)
      if err != nil {
        fmt.Println(err)
      }
    }
	})

	http.Handle("/", http.FileServer(http.Dir("./static")))
	fmt.Println(http.ListenAndServe(":8080", nil))
}
