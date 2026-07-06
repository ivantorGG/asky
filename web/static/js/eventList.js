document.addEventListener('DOMContentLoaded', () => {
    loadEvents();
});

async function loadEvents(){
    console.log('!!!!')
    const container = document.getElementById('events');
    console.log(container);
    const response = await fetch("/api/events/teacher");
    const events = await response.json();

    if (!container) {
        console.error("Контейнер #events не найден на странице!");
        return;
    }

    if (events.length === 0){
        console.log('events пустой!')
        return;
    }

    console.log("Тип данных:", typeof events);
    console.log("Содержимое:", events);

    let htmlContent = '';

    events.forEach(event => {
        htmlContent += `
            <div class="card shadow-sm mt-5 card-hover-bg position-relative" style="cursor: pointer;">
                <a href="/events/${event.code}/teacher" class="stretched-link text-decoration-none"></a>
                <div class="card-body p-4 d-flex justify-content-between align-items-center">
                    <h2 class="card-title mb-0">${event.title}</h2>
                    <button type="button" class="btn btn-outline-danger btn-lg ms-5 me-1 position-relative z-1" >
                        <div class="d-flex justify-content-between">
                            Удалить<i class="bi bi-trash"></i>
                        </div>
                    </button>
                </div>
            </div>
        `;
    });
    console.log(htmlContent)

    container.innerHTML = '';
    container.innerHTML = htmlContent;
}

async function createEvent(name){
    console.log('try to load events...', name);
    const response = await fetch("/api/events",
        {
            method: "POST",
            headers: {'Content-Type': 'application/json' },
            body: JSON.stringify({title: name})
        }
    )

    const data = await response.json()
    msg = data.message
    err = data.error
    console.log(msg)
    if (err !== ''){
        console.log(err)
    }
    if (msg === 'event_created'){
        console.log('try to load events...');
        loadEvents();
    }
}