const userRaw = sessionStorage.getItem("user");
user = JSON.parse(userRaw);
token = sessionStorage.getItem("token")
room = sessionStorage.getItem("room");
if (!token) {
    window.location.replace(`${location.protocol}//${location.hostname}:${location.port}/v1`);
}
const tokenWords = token.split(/(\s+)/).filter( function(e) { return e.trim().length > 0; } );
const wsUrl = "ws://localhost:8080/v1/ws?bearer="+tokenWords[1];
let wsSocket = new WebSocket(wsUrl);
const messagesURL = `${location.protocol}//${location.hostname}:${location.port}/v1/message`

function onLoadMain(){
    const welcomeMsg = document.getElementById("welcome");
    if(room == "free"){
        document.getElementById("chatbox").style.display = "block";
        document.getElementById("chatbox2").style.display = "none";
    }else{
        document.getElementById("chatbox").style.display = "none";
        document.getElementById("chatbox2").style.display = "block";
    }
    
    welcomeMsg.innerText = 'Welcome, '+ user.username;

    const url = new URL(messagesURL);
    const params = [['limit', '50'], ['room', room]];
    url.search = new URLSearchParams(params).toString()
    fetch(url, {method: 'GET', headers: {'Authorization': token}})
        .then(response => response.json())
        .then(messages => {
            messages.forEach(message => {
                appendMessage(new Date(message.created_at), message.user.username, message.message, room);
            });
        }).catch(error =>{
            console.log(error)
    })

}

wsSocket.onmessage = function (event){
    const msg = JSON.parse(event.data);
    appendMessage(new Date(msg.created_at), msg.user.username, msg.message, msg.room);
};

wsSocket.onerror = function (error){
    const alertMsg = document.getElementById("alert-login");
    alertMsg.innerText = "websocket error: " + JSON.stringify(error);
    alertMsg.style.visibility = "visible";
};

document.getElementById("chat-form").addEventListener("submit", function (event){
   event.preventDefault();
   const inText = document.getElementById("usermsg");
   const textMsg = inText.value;
   inText.value = "";
   msg = {
        user: {
           id: user.id,
           username: user.username,
        },
       message: textMsg,
       room: room,
       created_at: new Date().toISOString()
   }
   wsSocket.send(JSON.stringify(msg));
});

document.getElementById("exit").onclick = function (){
    sessionStorage.removeItem("user");
    sessionStorage.removeItem("room");
    window.location.replace(`${location.protocol}//${location.hostname}:${location.port}/v1`);
    return false;
};

function appendMessage(time, username, msg, msg_room){
    if(countChatMessages(msg_room) > 49){
        removeChatMessage(msg_room);
    }
    const msgHTML = `<div class="msgln"><span class="chat-time">${formatDate(time)}</span> <b class="user-name">${username}</b>${msg}<br></div>`;
    
    // command message
    if(msg_room == ""){
        msg_room = sessionStorage.getItem("room")
    }

    // check for rooms
    if(msg_room == "free"){
        document.getElementById("chatbox").insertAdjacentHTML("beforeend", msgHTML);
    }else{
        document.getElementById("chatbox2").insertAdjacentHTML("beforeend", msgHTML);
    }
}

function formatDate(date) {
    const h = "0" + date.getHours();
    const m = "0" + date.getMinutes();

    return `${h.slice(-2)}:${m.slice(-2)}`;
}

function countChatMessages(msg_room){
    if(msg_room == "free"){
        const chat = document.getElementById("chatbox");
        return chat.getElementsByClassName("msgln").length;
    }else{
        const chat = document.getElementById("chatbox2");
        return chat.getElementsByClassName("msgln").length;
    }
}

function removeChatMessage(msg_room){
    if(msg_room == "free"){
        const chat = document.getElementById("chatbox");
        const messages = chat.getElementsByClassName("msgln");
        chat.removeChild(messages[0]);
    }else{
        const chat = document.getElementById("chatbox2");
        const messages = chat.getElementsByClassName("msgln");
        chat.removeChild(messages[0]);
    }
}