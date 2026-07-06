function tryToSend(email, password){
    const regex = /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/;

    if (password.length < 6){
        showError('Пароль слишком короткий');
    }
    else if (email === '' || !regex.test(email)){
        showError('Почта неверно введена');
    }
    else {
        sendLogs(email, password);
    }

}

async function sendLogs(email, password){
    console.log("Емаил", email);
    console.log("Пароль", password);
    const data = {
        email: email,
        password: password
    };

    const response = await fetch("/api/login", {
        method: 'POST',
        headers: {'Content-Type': 'application/json' },
        body: JSON.stringify(data)
    });

    const response_json = await response.json();
    const msg = response_json.message
    if (msg === 'login_success'){
        location.href = '/events'
    }
    else if (err === 'bad_creditans'){
        showError('Неправильная почта или пароль!')
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