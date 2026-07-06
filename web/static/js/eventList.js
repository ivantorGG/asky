document.addEventListener('DOMContentLoaded', () => {
    loadEvents();
});

async function loadEvents(){
    const container = document.getElementById('events');
    const response = await fetch("/api/events/teacher");
    const events = await response.json();

    if (events.length === 0){
        return;
    }

    console.log("Тип данных:", typeof events);
    console.log("Содержимое:", events);

    events.forEach(event => {
        htmlContent += `
            <div class="card shadow-sm mt-5 card-hover-bg" style="cursor: pointer;">
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

    container.innerHTML = htmlContent;
}

async function createEvent(name){
    const response = await fetch("/api/events",
        {
            method: "POST",
            headers: {'Content-Type': 'application/json' },
            body: JSON.stringify({title: name})
        }
    )

    const data = await response.json()
    msg = data.messsage
    if (msg === 'event_created'){
        loadEvents();
    }
}