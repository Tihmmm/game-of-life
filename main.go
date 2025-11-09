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
	IsDeterministic bool
	GridSizeX       int
	GridSizeY       int
	InitialPoints   []set.Point
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		http.ServeFile(w, r, "./index.html")
	})

	mux.HandleFunc("/api/points", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("websocket upgrade failed: %v", err)
			http.Error(w, "websocket upgrade failed", http.StatusBadRequest)
			return
		}
		//defer conn.Close()

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

		g := game.NewGame(cfg.InitialPoints, cfg.GridSizeX, cfg.GridSizeY, cfg.IsDeterministic)
		log.Printf("board:%dx%d\ninitial state:\n%v\nIs determenistic: %t", g.GridSizeX, g.GridSizeY, g.Generations[0], g.IsDeterministic)

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

	log.Fatal(http.ListenAndServe(":8080", mux))
}
