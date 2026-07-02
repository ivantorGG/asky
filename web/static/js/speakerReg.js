async function sendRegs(email, password){
    console.log("Емаил", email);
    console.log("Пароль", password);

    const data = {
        email: email,
        password: password
    };

    const response = await fetch("http://127.0.0.1:8080/register", {
        method: 'POST',
        headers: {'Content-Type': 'application/json' },
        body: JSON.stringify(data)
    });
}