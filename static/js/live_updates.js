(function() {
    // Defines the WebSocket.
    var ws;

    // Handle WebSocket messages.
    function handleWsMessage(event) {
        // The data is organised <packet type>|<additional data...>.
        // This is because not all browsers support JSON, and the polyfill is pretty large!
        var split = event.data.split("|");

        // Get the packet type.
        var packetType = split.shift();

        if (packetType == "a" || packetType == "p") {
            // a/p = append or prepend mode (additional data will be <el name>|<element to append/prepend joined>).
            var elName = split.shift();
            var html = split.join("|");

            // Get the DOM parent.
            var parent = document.getElementById(elName);

            if (packetType == "a") {
                // Append the HTML.
                parent.innerHTML += html;
            } else {
                // Prepend the HTML.
                parent.innerHTML = html + parent.innerHTML;
            }
        } else if (packetType == "P") {
            // P = ping (send P back).
            ws.send("P");
        }
    }

    // Try and establish a WebSocket connection. Loop on error to do this.
    function attemptWsConnection() {
        var args = document.getElementById("ws_args").innerHTML;
        ws = new WebSocket(window.location.origin.replace(/^(h|H)(t|T){2}p|P/, "ws")+window.location.pathname.replace(/\/$/, "")+"/live_updates?"+args);
        function wsConnectionKilled() {
            console.log("Connection killed. Retrying WS connection in 5 seconds.");
            setTimeout(attemptWsConnection, 5000);
        }
        ws.onclose = wsConnectionKilled;
        ws.onmessage = handleWsMessage;
    }
    attemptWsConnection();
})();
