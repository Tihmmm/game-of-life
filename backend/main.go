package main

import (
	"encoding/json"
	"github.com/rs/cors"
	"github.com/tihmmm/game-of-life/game"
	"github.com/tihmmm/game-of-life/set"
	"log"
	"net/http"
	"strconv"
)

type req struct {
	Points []set.Point
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/points", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Only `POST` method allowed", http.StatusMethodNotAllowed)
			return
		}

		gridSizeX, gridSizeY := 10, 10
		if gridSizeXYParam := r.URL.Query()["gridSize"]; len(gridSizeXYParam) > 0 {
			gridSize, _ := strconv.Atoi(gridSizeXYParam[0])
			gridSizeX, gridSizeY = gridSize, gridSize
		}
		if gridSizeXParam := r.URL.Query()["gridSizeX"]; len(gridSizeXParam) != 0 {
			gridSizeX, _ = strconv.Atoi(gridSizeXParam[0])
		}
		if gridSizeYParam := r.URL.Query()["gridSizeY"]; len(gridSizeYParam) != 0 {
			gridSizeY, _ = strconv.Atoi(gridSizeYParam[0])
		}

		generationsNum := 400
		if generationsNumParam := r.URL.Query()["generationsNum"]; len(generationsNumParam) != 0 {
			generationsNum, _ = strconv.Atoi(generationsNumParam[0])
		}

		reqBody := new(req)
		err := json.NewDecoder(r.Body).Decode(reqBody)
		log.Printf("reqBody:\n%v\n", reqBody)
		if err != nil {
			log.Printf("Error decoding body: %v\n", err)
			http.Error(w, "Error parsing request data", http.StatusBadRequest)
			return
		}

		log.Printf("req:%v\n", *reqBody)
		res := game.NewGameWithNGenerations(reqBody.Points, gridSizeX, gridSizeY, generationsNum).Generations

		w.Header().Set("Content-Type", "application/json")

		err = json.NewEncoder(w).Encode(res)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	})

	handler := cors.Default().Handler(mux)

	log.Fatal(http.ListenAndServe(":8080", handler))
}
