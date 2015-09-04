var monolith = {
    server: "localhost:3000",
    id: null,
    socket: null,

    kickoff: function() {
        var id = localStorage.getItem("id");
        if (id === null) {
            var name = null;
            while (name === null || name === "") {
                name = prompt("What is your name?");
            }
        } else {
            id = parseInt(id)
        }

        monolith.socket = new WebSocket("ws://" + monolith.server + "/socket");
        monolith.socket.onopen = function(event) {
            monolith.send({
                type: "hello",
                text: monolith.name,
                date: Date.now(),
            });

            monolith.send({
                type: "get_messages",
            });
        }

        monolith.socket.onmessage = function(event) {
            var message = JSON.parse(event.data);

            console.log("got a _" + message.type + "_")
            switch(message.type) {
            case "new_user":
                monolith.id = parseInt(message.text);
                localStorage.setItem("id", monolith.id);
                break;
            case "messages":
                var messagesBox = document.getElementById("messages");
                messagesBox.innerHTML = "";
                srv_messages = JSON.parse(message.text);
                console.log(srv_messages);
                _.forEach(srv_messages, function(srv_message) {
                    var html = "<strong>" + srv_message.author_id + ":</strong> ";
                    html += srv_message.payload + "\n\n";

                    messagesBox.innerHTML += html;
                });
                break;
            case "message_broadcast":
                var messagesBox = document.getElementById("messages");
                messagesBox.innerHTML += "<strong>" + message.id + ":</strong> ";
                messagesBox.innerHTML += message.text + "\n\n";
                break;
            }
            //document.getElementById("textarea").innerHTML += event.data + "\n\n";
            console.log(event.data);
        }

        document.getElementById("send_text").onsubmit = function(event) {
            event.preventDefault();

            var message = {
                type: "send_message",
                text: document.getElementById("input").value,
                date: Date.now(),
                id: id,
            }

            monolith.socket.send(JSON.stringify(message));
        }
    },

    cleanup: function() {
        monolith.socket.close();
    },

    send: function(payload) {
        if (payload.date === undefined) {
            payload.date = Date.now()
        }

        monolith.socket.send(JSON.stringify(payload));
    }
}

window.onload = function() {
    monolith.kickoff();
}

window.onclose = function() {
    monolith.cleanup();
}
