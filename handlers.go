package main

import (
    "log"
    "time"
    "strings"
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
    if _, ok := users[key]; !ok {
        users[key] = User{ chars: map[string]*Character{} }
    }
    db_lock.Unlock()

    users[key].chars[name] = &Character{}

    go processUpdates(conn, key, name)

    
    users[key].chars[name].Init()

    for {
        msg := <- users[key].chars[name].command_queue

        err = conn.WriteMessage(websocket.TextMessage, []byte(msg))
        if err != nil {
            return
        }
    }

    return

}

func processUpdates(c *websocket.Conn, key string, name string) {
    for {
        _, p, err := c.ReadMessage()
        if err != nil {
            // Delete the cache if the client diconnects
            db_lock.Lock()
            delete(users[key].chars, name)
            db_lock.Unlock()
            return
        }

        // Update the DB
        db_lock.Lock()
        users[key].chars[name].stats = string(p)
        db_lock.Unlock()
    }
}

func processWebMessages(c *websocket.Conn, key string) {

    for {
        _, p, err := c.ReadMessage()
        if err != nil {
            return
        }

        msg := strings.Split(string(p),"|")


        // Update the DB
        users[key].chars[msg[0]].SendMessage(msg[1])
    }

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

    go processWebMessages(conn, key)

    // Start listening for socket messages
    for {

        // Get the data from the db
        db_lock.Lock()
        key_data := users[key]
        db_lock.Unlock()

        var encoded_data string

        for k, v := range key_data.chars { 
            encoded_data = encoded_data + k + "|" + v.stats + "\n"
        }

        err = conn.WriteMessage(websocket.TextMessage, []byte(encoded_data))
        if err != nil {
            return
        }

        time.Sleep(500 * time.Millisecond)
        
    }

    return
    
}

// START DEBUG //////////////////////////////////
func handleSay(w http.ResponseWriter, r *http.Request) {
    users["password"].chars["Dora2"].SendMessage("Hello")
    w.Write([]byte("Said Hello"))
}
// END DEBUG ////////////////////////////////////

// Serves the main screen
func handleDisplay(w http.ResponseWriter, r *http.Request) {

    t, err := template.New("base.html").Delims("[[", "]]").ParseFiles("templates/base.html", "templates/home.html")
    if err != nil {
        log.Println(err)
        return
    }
    t.Execute(w, "Home")

    return

}

// Serves the help page
func handleHelp(w http.ResponseWriter, r *http.Request) {

    t, err := template.New("base.html").Delims("[[", "]]").ParseFiles("templates/base.html", "templates/help.html")
    if err != nil {
        log.Println(err)
        return
    }
    t.Execute(w, "Help")

    return

}