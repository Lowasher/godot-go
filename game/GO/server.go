package main

import (
 "fmt"
 "net/http"
 "sync"
 "github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}

type Player struct {
 ID   int    json:"id"
 X    int    json:"x"
 Y    int    json:"y"
}

var (
 clients   = make(map[*websocket.Conn]int)
 players   = make(map[int]*Player)
 nextID    = 1
 mutex     = &sync.Mutex{}
)

func handleConnection(w http.ResponseWriter, r *http.Request) {
 conn, _ := upgrader.Upgrade(w, r, nil)
 
 mutex.Lock()
 id := nextID
 nextID++
 clients[conn] = id
 players[id] = &Player{ID: id, X: 100, Y: 100}
 mutex.Unlock()

 defer func() {
  mutex.Lock()
  delete(clients, conn)
  delete(players, id)
  mutex.Unlock()
  conn.Close()
 }()

 for {
  var input map[string]string
  if err := conn.ReadJSON(&input); err != nil { break }

  mutex.Lock()
  p := players[id]
  if input["key"] == "right" { p.X += 3 }
  if input["key"] == "left"  { p.X -= 3 }
  if input["key"] == "up"    { p.Y -= 3 }
  if input["key"] == "down"  { p.Y += 3 }

  // Рассылаем всем список ВСЕХ игроков и сообщение чата
  payload := map[string]interface{}{
   "players": players,
   "msg":     input["msg"],
  }
  for client := range clients {
   client.WriteJSON(payload)
  }
  mutex.Unlock()
 }
}

func main() {
 http.HandleFunc("/ws", handleConnection)
 fmt.Println("Сервер для двоих запущен!")
 http.ListenAndServe(":8080", nil)
}
