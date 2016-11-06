/////////////////////
// START STAT CODE
/////////////////////

var StatSocket = new WebSocket("ws://land.devpad.org/update");
var StatKey = "";

StatSocket.onopen = function (event) {

    // Identify yourself to the stat server
    StatSocket.send(StatKey);
    StatSocket.send(character.name);

    // Start updating your stats every 500ms
    setInterval(function(){
        StatSocket.send(character.hp +
            "|" + character.max_hp +
            "|" + character.mp +
            "|" + character.max_mp +
            "|" + character.xp +
            "|" + character.max_xp +
            "|" + character.level
        );
    },500); // Loops every 0.5 seconds.

};

///////////////////
// END STAT CODE
///////////////////