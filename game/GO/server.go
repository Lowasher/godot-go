package main

import (
 "fmt"
 "net/http"
 "sync"
 "github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
 CheckOrigin: func(r *http.Request) bool { return true },
}

type State struct {
 X   int
 Y   int
 Msg string
}

var clients = make(map[*websocket.Conn]*State)
var mutex = &sync.Mutex{}

func handleConnection(w http.ResponseWriter, r *http.Request) {
 conn, err := upgrader.Upgrade(w, r, nil)
 if err != nil {
  return
 }
 
 mutex.Lock()
 clients[conn] = &State{X: 100, Y: 100, Msg: ""}
 mutex.Unlock()

 defer func() {
  mutex.Lock()
  delete(clients, conn)
  mutex.Unlock()
  conn.Close()
 }()

 for {
  var input map[string]string
  if err := conn.ReadJSON(&input); err != nil {
   break
  }

  mutex.Lock()
  st := clients[conn]
  
  if input["key"] == "right" { st.X += 3 }
  if input["key"] == "left"  { st.X -= 3 }
  if input["key"] == "up"    { st.Y -= 3 }
  if input["key"] == "down"  { st.Y += 3 }
  
  msgToSend := ""
  if input["msg"] != "" {
   msgToSend = input["msg"]
  }

  for client := range clients {
   client.WriteJSON(map[string]interface{}{
    "X":   st.X,
    "Y":   st.Y,
    "Msg": msgToSend,
   })
  }
  mutex.Unlock()
 }
}

func main() {
 http.HandleFunc("/ws", handleConnection)
 fmt.Println("СЕРВЕР ЗАПУЩЕН! ПАЦАНЫ, НЕ ПЛАЧЬТЕ!")
 http.ListenAndServe(":8080", nil)
}