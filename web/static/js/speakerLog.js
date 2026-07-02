async function sendLogs(email, password){
    const inputData = { email, password };

    const response = await fetch("http://127.0.0.1:8080/login", {
        method: 'POST',
        headers: {'Content-Type': 'application/json'},
        body: JSON.stringify(inputData)
    });

    const data = await response.json();

    if (response.ok){
        location.href = '/eventList';
    } else {
        document.getElementById("errorMessage").innerText = data.error;
        document.getElementById("loginError").classList.remove("d-none");
    }
}