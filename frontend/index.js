"use strict";

const button = document.getElementById("button");
button.addEventListener("click", function() {
    const socket = new WebSocket("ws://localhost:8000/ws/update-order"); 

    // open
    socket.addEventListener("open", function(event) {
        console.log("WebSocket is connected:", event);
    });

    // receives a message
    socket.addEventListener("message", function(event) {
        console.log(event.data);
    });

    // close
    socket.addEventListener("close", function(event) {
        console.log(event);
    });

    // error
    socket.addEventListener("error", function(event) {
        console.error(event);
    });
});
