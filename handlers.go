package main

import (
    "log"
    "time"
    "net/http"
    "html/template"

    "github.com/gorilla/websocket"
)

// The game side connects to this
func handleUpdate(w http.ResponseWriter, r *http.Request) {

    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Println(err)
        return
    }

    // This is the API key and character name
    var key string
    var name string

    // First message will be the client's API key
    _, p, err := conn.ReadMessage()
    if err != nil {
        return
    }
    key = string(p)

    // Second message will be the client's character name
    _, p, err = conn.ReadMessage()
    if err != nil {
        return
    }
    name = string(p)

    // Make sure the API key exists
    db_lock.Lock()
    if _, ok := db[key]; !ok {
        db[key] = map[string]string{}
    }
    db_lock.Unlock()

    // Start listening for socket messages
    for {

        _, p, err = conn.ReadMessage()
        if err != nil {
            // Delete the cache if the client diconnects
            db_lock.Lock()
            delete(db[key], name)
            db_lock.Unlock()
            return
        }

        // Update the DB
        db_lock.Lock()
        db[key][name] = string(p)
        db_lock.Unlock()

    }

    return

}

// The site connects to this
func handleReport(w http.ResponseWriter, r *http.Request) {

    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Println(err)
        return
    }

    // This is the API key
    var key string

    // First message will be the client's API key
    _, p, err := conn.ReadMessage()
    if err != nil {
        return
    }
    key = string(p)

    // Start listening for socket messages
    for {

        // Get the data from the db
        db_lock.Lock()
        key_data := db[key]
        db_lock.Unlock()

        var encoded_data string

        for k, v := range key_data { 
            encoded_data = encoded_data + k + "|" + v + "\n"
        }

        err = conn.WriteMessage(websocket.TextMessage, []byte(encoded_data))
        if err != nil {
            return
        }

        time.Sleep(500 * time.Millisecond)
        
    }

    return
    
}

// This serves the actual webpage
func handleDisplay(w http.ResponseWriter, r *http.Request) {

    t, err := template.New("base.html").Delims("[[", "]]").ParseFiles("templates/base.html", "templates/home.html")
    if err != nil {
        log.Println(err)
        return
    }
    t.Execute(w, nil)

    return

}

// This serves the actual webpage
func handleHelp(w http.ResponseWriter, r *http.Request) {

    t, err := template.New("base.html").Delims("[[", "]]").ParseFiles("templates/base.html", "templates/help.html")
    if err != nil {
        log.Println(err)
        return
    }
    t.Execute(w, nil)

    return

}