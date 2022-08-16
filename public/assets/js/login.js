document.getElementById("login-form").addEventListener("submit", function (event){
    event.preventDefault();
    const url = "/v1/login"
    let data = {
        username: document.getElementById("username").value,
        password: document.getElementById("password").value,
        room: document.getElementById("room").value,
    }
    fetch(url, {
        method: "POST",
        headers: {"Content-Type": "application/json;charset=utf-8"},
        body: JSON.stringify(data)
    }).then(result => {
        result.json().then(body => {
            const response = body;
            if (result.ok){
                sessionStorage.setItem("token", result.headers.get("Authorization"));
                sessionStorage.setItem("user", JSON.stringify(response));
                sessionStorage.setItem("room", result.headers.get("room"));
                window.location.href = "chat.html";
            }else{
                const alertMsg = document.getElementById("alert-login");
                alertMsg.innerText = "Error user authentication: " + JSON.stringify(response);
                alertMsg.style.visibility = "visible";
            }
        });
    });
    return false
});