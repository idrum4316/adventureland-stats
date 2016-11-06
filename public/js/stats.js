var template = $('#handlebars-card').html();
var templateScript = Handlebars.compile(template);

var template = $('#handlebars-welcome').html();
var welcomeScript = Handlebars.compile(template);

// This will keep track of all active cards
var active_cards = [];

// Add a new ID to track
function track() {
	key = localStorage.getItem("key");
	sock = new WebSocket(get_ws_url("report"));

	sock.onopen = function (event) {
		sock.send(key);

		sock.onmessage = function (event) {

			rows = event.data.split("\n");
			names_this_message = [];
			cards_changed = false;

			for (var i = 0; i < rows.length - 1; i++) {
				stats = rows[i].split("|");
				name = stats[0]
				names_this_message.push(name);

				// Check to see if the card exists
				if (!document.getElementById("card_" + name)) {

					// Card doesn't exist, so we'll create it
					var card_context = { "name" : name };
					card_html = templateScript(card_context);
					$("#card_list").append(card_html);
					active_cards.push(name);
					cards_changed = true;

					if(!(typeof(componentHandler) == 'undefined')){
						componentHandler.upgradeAllRegistered();
					}

				}

				// Now that we know it exists, update the values
				if (stats.length === 8) {
					document.querySelector('#health_bar_' + name).MaterialProgress.setProgress((stats[1]/stats[2]) * 100);
					document.querySelector('#mana_bar_' + name).MaterialProgress.setProgress((stats[3]/stats[4]) * 100);
					document.querySelector('#xp_bar_' + name).MaterialProgress.setProgress((stats[5]/stats[6]) * 100);
					$("#health_nums_" + name).html(stats[1] + "/" + stats[2]);
					$("#mana_nums_" + name).html(stats[3] + "/" + stats[4]);
					$("#xp_nums_" + name).html(stats[5] + "/" + stats[6] + " (" + Math.round((stats[5] / stats[6]) * 100) + "%)");
					$("#lvl_" + name).html(stats[7]);
				}

			}

			// Remove cards we're no longer recieving data for
			for (var i = 0; i < active_cards.length; i++) {
				if (names_this_message.indexOf(active_cards[i]) === -1) {
					$("#card_" + active_cards[i]).remove();
					cards_changed = true;
				}
			}
			
			var i = active_cards.length
			while (i--) {
				if (names_this_message.indexOf(active_cards[i]) === -1) {
					$("#card_" + active_cards[i]).remove();
					active_cards.splice(i, 1);
				}
			}

			if (cards_changed) {
				if (active_cards.length === 0) {
					add_welcome_message();
				} else {
					$("#welcome_card").remove()
				}
			}

			// console.log(event.data);
		}
	};

	sock.onclose = function(){
        //try to reconnect
        sock = null;
        track();
    };

}