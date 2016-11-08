package main


// A User just has an array of characters
type User struct {
	chars map[string]*Character
}


// A character has stats and a queue of commands
type Character struct {
	stats	      string
	command_queue chan string
}


// Initialize a command queue with a 5 command buffer
func (c *Character) Init() {
	c.command_queue = make(chan string, 5)
}


// Send a message through this character's channel
func (c *Character) SendMessage(msg string) {
	c.command_queue <- msg
}


// Map Users to their keys
var users map[string]User