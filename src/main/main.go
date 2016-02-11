package main

import (
	"fmt"
	"game"
	"html/template"
	"net/http"
	"ws"
)

var homeTempl = template.Must(template.ParseFiles("templates/home.html"))

func serveHome(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	homeTempl.Execute(w, r.Host)
}

func main() {
	fmt.Println("start")
	game := game.NewGame()
	go game.Start()

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", ws.HandlerFactory(game))

	err := http.ListenAndServe(":7101", nil)
	if err != nil {
		panic(err)
	}
	fmt.Println("stop")
}
