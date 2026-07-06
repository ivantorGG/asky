let questions = [];

async function loadInitialQuestions() {
    try {
        const response = await fetch('/events/{code}/questions');
        if (!response.ok) throw new Error('Ошибка загрузки');
        
        questions = await response.json();
        renderQuestions(questions);
        
        // После успешной загрузки подключаемся к реальному времени
        connectToRealTimeUpdates(); 
    } catch (error) {
        console.error('Ошибка:', error);
    }
}

function renderQuestions(list) {
    const container = document.getElementById('questions-list');
    container.innerHTML = ''; // Очищаем список
    
    list.forEach(q => {
        const div = document.createElement('div');
        div.className = 'question-item';
        div.id = `q-${q.id}`; // Важно для точечного обновления
        div.textContent = q.text;
        container.appendChild(div);
    });
}