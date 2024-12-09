import Game from "./game";
import EventMsg from "./event";


let game = new Game();

/**
 * login will send a login request to the server and then connect websocket
 */
export function login() {
    connectWebsocket();

    return false;
}

document.getElementById("login-btn").addEventListener('click', login);

/**
 * connectWebsocket will connect to websocket and add listeners
 */
function connectWebsocket() {
    if (window["WebSocket"]) {
        console.log("supports websockets");
        game.conn = new WebSocket("ws://localhost:8080/ws");

        // Onopen
        game.conn.onopen = function (evt) {
            let username = document.getElementById("username").value;
            sendEvent("login", username);
        };

        game.conn.onmessage = function (event) {
            // parse websocket message as JSON
            const eventData = JSON.parse(event.data);
            console.log(eventData);
            routeEvent(eventData, game);

        };

        game.conn.onclose = function (event) {
            cancelAnimationFrame(game.animationFrame);
        }
    } else {
        alert("Not supporting websockets");
    }
}

/**
     * routeEvent is a proxy function that routes
     * events into their correct Handler
     * based on the type field
     * */
function routeEvent(event) {
    if (event.type === undefined) {
        alert("no 'type' field in event");
    }
    console.log("EVENT FROM SERVER", event);
    switch (event.type) {
        case "state":
            game.setCurrentState(event.payload);
            break;
        case "start":
            game.start(event.payload)
            game.animationFrame = requestAnimationFrame(() => game.update());
            
            addListenters();
            break
        case "end":
            alert(event.payload.data);
            break
        default:
            alert("unsupported message type");
            break;
    }
}

/**
 * sendEvent
 * eventName - the event name to send on
 * payload - the data payload
 * */
function sendEvent (eventName, payload) {
    // Create a event Object with a event named send_message
    const event = new EventMsg(eventName, payload);
    console.log("SENDEVENT", event);
    game.conn.send(JSON.stringify(event));
}

function addListenters() {
    document.addEventListener("keydown", function (event) {
        if (event.code == "KeyA") {
            sendEvent("keydown", "left");
        }
        if (event.code == "KeyD") {
            sendEvent("keydown", "right");
        }
        if (event.code == "KeyW") {
            sendEvent("keydown", "forward");
        }
        if (event.code == "KeyS") {
            sendEvent("keydown", "back");
        }
        if (event.code == "Space") {
            sendEvent("keydown", "space");
        }
    });
    document.addEventListener("keyup", function (event) {
        if (event.code == "KeyA") {
            sendEvent("keyup", "left");
        }
        if (event.code == "KeyD") {
            sendEvent("keyup", "right");
        }
        if (event.code == "KeyW") {
            sendEvent("keyup", "forward");
        }
        if (event.code == "KeyS") {
            sendEvent("keyup", "back");
        }
        if (event.code == "Space") {
            sendEvent("keyup", "space");
        }
    });
}