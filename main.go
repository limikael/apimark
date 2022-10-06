package main

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024 * 4,
	WriteBufferSize: 1024 * 4,
}

func main() {
	http.HandleFunc("/", serveRoot)
	http.HandleFunc("/ws", serveWebSocket)
	http.HandleFunc("/api", serveHttpApi)

	err := http.ListenAndServe("0.0.0.0:8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func serveRoot(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "home.html")
}

func serveWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	for {
		mt, _, err := conn.ReadMessage()

		err = conn.WriteMessage(mt, []byte("ok"))
		if err != nil {
			log.Println("write: ", err)
			break
		}
	}
}

func serveHttpApi(w http.ResponseWriter, r *http.Request) {
	_, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Println("read err: ", err)
		http.Error(w, "Cannot read body", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("ok"))
}
