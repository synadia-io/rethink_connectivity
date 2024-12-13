# Introduction

Hey everyone, my name is Nate and I like NATS! I'm a software architect who's
been building web applications and infrastructure for about 20 years. And I got a new job!

I recently joined Synadia as a DX engineer, where I'm excited to work on our web
properties and educational content. While I've built and maintained systems with
RabbitMQ, Sidekiq, Kafka, Redis, and plethora AWS infra components, you might be
surprised to learn I'd never heard of NATS prior to meeting Jeremy at what I
would come to learn was my first job interview.

[NATS architecture diagram]

I've come to love the simple and extensible architecture NATS provides for
handling so many of the requirements of distributed systems. Since I've been
learning the ropes so to speak, I thought I'd bring you along as we explore the
basics of what this NATS thing is all about.

Welcome to our 2024 refresh on NATS!

NATS is an open-source messaging system with one goal: connecting everything.
Think of it as the nervous system for your distributed applications. It's fast,
scalable, and designed for modern cloud-native and edge computing environments.

[outline]

In this video, we'll:

1. Cover the basics of NATS architecture
2. Build a real-time chat app
3. Add persistence with JetStream
4. Spice it up with a Twitch integration

We'll keep things snappy, using diagrams and pre-written code to get you up to speed quickly.
Let's dive in and see what NATS can do!

# NATS Architectural Overview

Alright, let's break down how NATS works. We'll cover three key aspects: core 
NATS, clustering, and JetStream.
[pub-sub diagram]
First up, core NATS. At its heart, NATS uses a publish-subscribe model. Here's how it works:

Publishers send messages on specific subjects.
The NATS server acts as a message broker, the courier of the system.
Subscribers listen for messages on subjects they're interested in.

It's that simple. This model allows for flexible, decoupled communication
between components in your system.

[cluster diagram]
Now, let's talk scalability. NATS shines in distributed environments through
clustering:

1. Multiple NATS servers form a cluster.
2. Clients can connect to any server in the cluster.
3. Servers use a Gossip Protocol to share information.

This setup gives you high availability and fault tolerance. If one server goes
down, clients can seamlessly connect to another.

# NATS JetStream Overview
[jetstream diagram]
Finally, let's look at JetStream, NATS' persistence layer:

JetStream adds message streaming and persistence to NATS.
It uses 'Streams' to store messages.
'Consumers' read messages from streams, allowing for things like replay and
filtered consumption.

JetStream is perfect for event-driven architectures and situations where you
need guaranteed message delivery.
And there you have it – NATS in a nutshell. Simple, scalable, and now with
built-in persistence. Enough talk, let's see it in action!


# Project Setup

We'll start by installing NATS and the NATS CLI locally with a package manager.
Now we can run a nats-server locally to show you how this works. Using nats CLI
in another terminal we can verify that NATS is running and we are connected.

Now let's build our first chat client.

```bash
mkdir nats-chat && cd "$_"
go mod init nats-chat
go get github.com/nats-io/nats.go
```

# Basic Chat Application

```go
package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/nats-io/nats.go"
)

func main() {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer nc.Close()

	nc.Subscribe("chat", func(m *nats.Msg) {
		fmt.Printf("Received: %s\n", string(m.Data))
	})

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		msg := scanner.Text()
		if strings.ToLower(msg) == "quit" {
			break
		}
		nc.Publish("chat", []byte(msg))
	}
}
```

Let's walk through this extremely simple chat client:
1. First we connect to the NATS server using the default URL
2. We set up a subscriber that listens for messages on the "chat" subject,
printing anything received to the console.
3. Then we have a simple loop that publishes anything entered on the command
line to the "chat subject"...until receiving the command "quit"

We can run this in two terminals and see the messages being sent and received.
However, if we start a new client, the chat history is empty. A client must be
connected when a message is sent to see it.

# Adding Web Client

First we'll add a simple backend with Go to serve a websocket endpoint and
static HTML. Every connected subscribes to the "chat" subject and is able to
publish messages to it. In this extremely simple front-end, we have a simple
chat interface in HTML that uses JS to manage the websocket connection.

web backend:
```go
package main

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

func main() {
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
```

Web chat front-end:

```html
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>NATS Web Chat</title>
    <style>
        body { font-family: Arial, sans-serif; max-width: 800px; margin: 0 auto; padding: 20px; }
        #chat { height: 300px; overflow-y: scroll; border: 1px solid #ccc; padding: 10px; margin-bottom: 10px; }
        #message { width: 70%; padding: 5px; }
        #send { padding: 5px 10px; }
    </style>
</head>
<body>
    <h1>NATS Web Chat</h1>
    <div id="chat"></div>
    <input type="text" id="message" placeholder="Type your message...">
    <button id="send">Send</button>

    <script>
        const chat = document.getElementById('chat');
        const messageInput = document.getElementById('message');
        const sendButton = document.getElementById('send');

        const socket = new WebSocket('ws://localhost:8080/ws');

        socket.onmessage = function(event) {
            const message = document.createElement('p');
            message.textContent = event.data;
            chat.appendChild(message);
            chat.scrollTop = chat.scrollHeight;
        };

        function sendMessage() {
            if (messageInput.value) {
                socket.send(messageInput.value);
                messageInput.value = '';
            }
        }

        sendButton.onclick = sendMessage;
        messageInput.onkeypress = function(event) {
            if (event.key === 'Enter') {
                sendMessage();
            }
        };
    </script>
</body>
</html>
```

# Enhancing with JetStream

Now the obvious issue with this application so far is that you will only see
messages sent when you have the chat open. Let's fix that by adding persistence
using JetStream.

```go
package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/nats-io/nats.go"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	// Create JetStream context
	js, err := nc.JetStream()
	if err != nil {
		log.Fatal(err)
	}

	// Create a stream for chat messages
	_, err = js.AddStream(&nats.StreamConfig{
		Name:     "CHAT",
		Subjects: []string{"chat.>"},
	})
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}
		defer conn.Close()

		// Subscribe to chat messages
		sub, err := js.Subscribe("chat.messages", func(msg *nats.Msg) {
			conn.WriteMessage(websocket.TextMessage, msg.Data)
		})
		if err != nil {
			log.Println(err)
			return
		}
		defer sub.Unsubscribe()

		// Retrieve chat history
		consumer, err := js.PullSubscribe("chat.messages", "")
		if err != nil {
			log.Println(err)
			return
		}

		msgs, err := consumer.Fetch(100)
		if err != nil {
			log.Println(err)
		} else {
			for _, msg := range msgs {
				conn.WriteMessage(websocket.TextMessage, msg.Data)
			}
		}

		// Handle incoming messages
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Println(err)
				break
			}
			_, err = js.Publish("chat.messages", message)
			if err != nil {
				log.Println(err)
			}
		}
	})

	http.Handle("/", http.FileServer(http.Dir("./static")))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
```

# Twitch Chat Integration


```go
package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gempir/go-twitch-irc/v2"
	"github.com/nats-io/nats.go"
)

func main() {
	// Connect to NATS
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	js, err := nc.JetStream()
	if err != nil {
		log.Fatal(err)
	}

	// Twitch bot setup
	client := twitch.NewClient("YourBotUsername", "oauth:YourOAuthToken")

	client.OnPrivateMessage(func(message twitch.PrivateMessage) {
		chatMessage := fmt.Sprintf("%s: %s", message.User.DisplayName, message.Message)
		
		// Publish to NATS JetStream
		_, err := js.Publish("chat.messages", []byte(chatMessage))
		if err != nil {
			log.Println("Error publishing to NATS:", err)
		}
	})

	client.Join("YourTwitchChannel")

	err = client.Connect()
	if err != nil {
		panic(err)
	}
}
```

# Conclusion and further enhancements

Great job, everyone! Let's recap what we've built and explore some exciting 
ways to take this project even further.

## What We've Accomplished

1. Set up a local NATS environment
2. Created a basic chat application using core NATS
3. Added a web client for broader accessibility
4. Implemented persistence with NATS JetStream
5. Integrated Twitch chat, demonstrating NATS' flexibility

Our chat application now showcases several key NATS features: real-time
messaging, persistence, and system integration.

## Potential Enhancements

1. Multi-room support:
   - Use NATS subjects to create multiple chat rooms
   - Example: "chat.room.general", "chat.room.tech", etc.

2. User authentication:
   - Implement a login system
   - Use NATS authorization to control access to different chat rooms

3. Message filtering:
   - Leverage NATS subject hierarchies for message categorization
   - Implement client-side filtering based on message types or user preferences

4. Scalability demonstration:
   - Set up a NATS cluster to show how the application scales
   - Implement load balancing for web servers

5. More integrations:
   - Add connectors for other platforms (Discord, Slack, etc.)
   - Demonstrate NATS' power in creating a unified chat ecosystem

6. Rich media support:
   - Extend the application to handle file uploads, images, or even video streams
   - Use NATS streaming for larger data transfers

7. Real-time analytics:
   - Implement a separate service that consumes chat messages
   - Generate real-time statistics on chat activity, popular topics, etc.

These enhancements showcase NATS' versatility in building complex, distributed
systems. Each one demonstrates a different aspect of what NATS can do, from
fine-grained message routing to handling high-throughput data streams.

Remember, NATS is all about connecting systems, services, and devices flexibly
and at scale. The possibilities are endless!

Thanks for joining me on this journey through NATS. Happy coding, and I can't
wait to see what you build with NATS!

