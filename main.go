package main

import (
    "os"
    "net/http"

    "github.com/gorilla/mux"
    "github.com/gorilla/handlers"
    "github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
    CheckOrigin: func(r *http.Request) bool {
        return true
    },
}

func main() {
    users = map[string]User{}

    router := mux.NewRouter()
    router.HandleFunc("/", handleDisplay).Methods("GET")
    router.HandleFunc("/help", handleHelp).Methods("GET")
    router.HandleFunc("/say", handleSay).Methods("GET")
    router.HandleFunc("/update", handleUpdate)
    router.HandleFunc("/report", handleReport)

    // Serves static files
    router.PathPrefix("/public/").Handler(http.StripPrefix("/public/",
        http.FileServer(http.Dir("public/"))))
    router.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, "public/favicon.ico")
    })

    logged_router := handlers.LoggingHandler(os.Stdout, router)
    http.ListenAndServe(":8080", logged_router)
}
