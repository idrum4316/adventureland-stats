# Adventureland Stats

Real-time stats for your AdventureLand characters.

![Imgur](http://i.imgur.com/oJHmZrm.gif)

## Building
Building should be pretty simple. All dependencies are vendored, so there's no need to fetch any external libraries. A simple `go build` should do the trick.

## Running
The `public` and `templates` folders need to be in the same folder as the executable. Other than that, I've tried to make sure that the websocket addresses and stuff are generated dynamically based on the host, so there should be no need to update URLs before running it yourself.