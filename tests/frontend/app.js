const API_URL = ''; // Оставляем пустую строку для работы на одном порту

async function createEvent() {
    const titleInput = document.getElementById('eventTitle');
    const postOutput = document.getElementById('postOutput');
    const postStatus = document.getElementById('postStatus');
    
    if (!titleInput.value) {
        alert("Введи название ивента!");
        return;
    }

    try {
        postStatus.textContent = "Отправка...";
        const response = await fetch(`${API_URL}/events/new`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ title: titleInput.value })
        });
        
        postStatus.textContent = `Статус ответа сервера: ${response.status} ${response.statusText}`;
        const data = await response.json();
        postOutput.textContent = JSON.stringify(data, null, 2);
        
        // Автоматически обновляем список после создания
        getEvents();
    } catch (err) {
        postStatus.textContent = "Ошибка соединения";
        postOutput.textContent = JSON.stringify({ error: err.message }, null, 2);
    }
}

async function getEvents() {
    const getOutput = document.getElementById('getOutput');
    const getStatus = document.getElementById('getStatus');
    
    try {
        getStatus.textContent = "Получение данных...";
        const response = await fetch(`${API_URL}/events`);
        
        getStatus.textContent = `Статус ответа сервера: ${response.status} ${response.statusText}`;
        const data = await response.json();
        getOutput.textContent = JSON.stringify(data, null, 2);
    } catch (err) {
        getStatus.textContent = "Ошибка соединения";
        getOutput.textContent = JSON.stringify({ error: err.message }, null, 2);
    }
}