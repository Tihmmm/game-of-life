package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
	"github.com/tihmmm/game-of-life/game"
	"github.com/tihmmm/game-of-life/set"
)

type Req struct {
	GridSizeX     int
	GridSizeY     int
	InitialPoints []set.Point
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/points", func(w http.ResponseWriter, r *http.Request) {
		//if r.Method != http.MethodGet {
		//	http.Error(w, "Only `Get` method allowed", http.StatusMethodNotAllowed)
		//	return
		//}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("websocket upgrade failed: %v", err)
			http.Error(w, "websocket upgrade failed", http.StatusBadRequest)
			return
		}

		var cfg Req
		if _, msg, err := conn.ReadMessage(); err == nil && len(msg) > 0 && msg[0] == '{' {
			_ = json.Unmarshal(msg, &cfg)
		}

		if cfg.GridSizeX <= 0 {
			cfg.GridSizeX = 10
		}
		if cfg.GridSizeY <= 0 {
			cfg.GridSizeY = 10
		}

		g := game.NewGame(cfg.InitialPoints, cfg.GridSizeX, cfg.GridSizeY)

		if err := conn.WriteJSON(g); err != nil {
			log.Printf("Error writing JSON: %v\n", err)
			return
		}

		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				log.Printf("Error reading message: %v\n", err)
				return
			}

			switch string(msg) {
			case "next":
				g.NextGen()
				if err := conn.WriteJSON(g.GetLastGen()); err != nil {
					log.Printf("Error writing JSON: %v\n", err)
					return
				}
			case "previous":
				genNum := len(g.Generations)
				if err := conn.WriteJSON(g.GetNthGeneration(genNum - 1)); err != nil {
					log.Printf("Error writing JSON: %v\n", err)
				}
			default:
				n, err := strconv.ParseInt(string(msg), 10, 64)
				if err != nil || int(n) > len(g.Generations) {
					return
				}
				_ = conn.WriteJSON(g.GetNthGeneration(int(n)))
			}
		}
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		http.ServeFile(w, r, "./index.html")
	})

	log.Fatal(http.ListenAndServe(":8080", mux))
}
