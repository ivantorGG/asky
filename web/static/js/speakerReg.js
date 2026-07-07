

function tryToSend(email, password){
    const regex = /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/;

    if (password.length < 6){
        showError('Пароль слишком короткий')
    }
    else if (email === '' || !regex.test(email)){
        showError('Почта неверно введена')
    }
    else{
        sendRegs(email, password);
    }
}

async function sendRegs(email, password){
    console.log("Емаил", email);
    console.log("Пароль", password);

    const data = {
        email: email,
        password: password
    };

    const response = await fetch("/api/register", {
        method: 'POST',
        headers: {'Content-Type': 'application/json' },
        body: JSON.stringify(data)
    });

    const response_json = await response.json();
    const msg = response_json.message;
    const err = response_json.error;

    if (msg === 'registration_success'){
        location.href = '/login'
        return;
    }

    switch (err){
        case 'server_error':
            showError('Сервер не отвечает...')
            break;
        case 'bad_request':
            showError('Ошибка на сайте')
            break;
        case 'invalid_input':
            showError('Почта и пароль не могут быть пустыми!')
            break;
        case 'email_already_exists':
            showError('Такой email уже зарегистрирован!')
            break;
    }
}

function showError(message) {
    const alert = document.getElementById('loginError');
    document.getElementById('errorMessage').textContent = message;
    
    // 1. Убираем d-none (элемент становится видимым, но opacity: 0 из-за fade)
    alert.classList.remove('d-none');
    
    // 2. Небольшая задержка, чтобы браузер успел применить display: block
    setTimeout(() => {
        // 3. Добавляем show - запускается transition к opacity: 1
        alert.classList.add('show');
    }, 10);
}

function hideError() {
    const alert = document.getElementById('loginError');
    
    // 1. Убираем show - запускается transition к opacity: 0
    alert.classList.remove('show');
    
    // 2. Ждем окончания анимации (обычно 150ms в Bootstrap)
    setTimeout(() => {
        // 3. Скрываем полностью
        alert.classList.add('d-none');
    }, 150);
}