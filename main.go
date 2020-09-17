package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

// need to upgrade an incoming connection from a standard HTTP endpoint to a long-lasting WebSocket connection
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func home(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w,r,"static/index.html")
}

func webSocketEndpoint(w http.ResponseWriter, r *http.Request) {
	// whether or not an incoming request from a different domain is allowed to connect
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	// upgrade the incoming HTTP connection
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Websocket connected")
	err = ws.WriteMessage(1, []byte("Hi Client!"))
	if err != nil {
		log.Println(err)
	}
	reader(ws)

}

// define a reader which will listen for
// new messages being sent to our WebSocket
// endpoint
func reader(conn *websocket.Conn) {
	for {
		// read in a message
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		// print out that message for clarity
		fmt.Println(string(p))

		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Println(err)
			return
		}

	}
}

func main() {
	http.HandleFunc("/", home)
	http.HandleFunc("/ws", webSocketEndpoint)
	log.Fatal(http.ListenAndServe(":3000", nil))
}
