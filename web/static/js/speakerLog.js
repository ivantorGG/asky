async function sendLogs(email, password){
    console.log("Емаил", email);
    console.log("Пароль", password);

    const data = {
        email: email,
        password: password
    };

    const data = await fetch("http://127.0.0.1:8080/login", {
        method: 'POST',
        headers: {'Content-Type': 'application/json' },
        body: JSON.stringify(data)
    }).json();
    const err = data.error
    if (err === ''){
        location.href = '/eventList'
    }
    else{
        
    }
}