document.addEventListener("DOMContentLoaded", () => {

    const saveButton = document.getElementById("saveEvent");
    const input = document.getElementById("eventName");

    saveButton.addEventListener("click", () => {
        createEvent(input.value);
    });

    input.addEventListener("keydown", e => {
        if (e.key === "Enter") {
            e.preventDefault();
            createEvent(input.value);
        }
    });

    // Переключение вкладок
    document.querySelectorAll(".eventsTab").forEach(tab => {

        tab.addEventListener("click", () => {

            document
                .querySelectorAll(".eventsTab")
                .forEach(t => t.classList.remove("active"));

            tab.classList.add("active");

            loadEvents(tab.dataset.type);

        });

    });

    loadEvents("teacher");

});

// =========================================

async function loadEvents(type = "teacher") {

    const container = document.getElementById("events");

    const response = await fetch(`/api/events/${type}`);

    const data = await response.json();

    if (data.error === "unauthorized") {
        location.href = "/login";
        return;
    }

    container.innerHTML = "";

    if (data.length === 0) {

        const empty = type === "teacher"
            ? {
                title: "Пока нет мероприятий",
                text: "Создайте своё первое мероприятие, чтобы студенты могли подключиться."
            }
            : {
                title: "Вы ещё нигде не участвовали",
                text: "После подключения к мероприятию оно появится здесь."
            };

        container.innerHTML = `

        <div class="col-lg-9">

            <div class="emptyState reveal active">

                <i class="bi bi-calendar2-plus"></i>

                <h2>${empty.title}</h2>

                <p>${empty.text}</p>

            </div>

        </div>

        `;

        return;
    }

    const cardsHtml = await Promise.all(

        data.map(async event => {

            const questionCount = await getQuestionsCount(event.code);

            const deleteButton = type === "teacher"
                ? `
                <button
                    class="deleteButton"
                    onclick="event.preventDefault();event.stopPropagation();deleteEvent('${event.code}')">
                    <i class="bi bi-trash"></i>
                    Удалить
                </button>
                `
                : "";

            return `

<div class="col-lg-10 reveal active">

    <div class="eventCard">

        <a
            href="/events/${event.code}/${type}"
            class="stretched-link">

        <div class="eventBody">

            <div class="eventInfo">

                <h2 class="eventTitle">

                    ${event.title}

                </h2>

                <div class="eventSubtitle">

                    Вопросов: ${questionCount}

                </div>

            </div>

            ${deleteButton}

        </div>
        </a>

    </div>

</div>

`;

        })

    );

    container.innerHTML = cardsHtml.join("");

    enableGlow();

}

// =========================================

async function getQuestionsCount(code) {
        try {

        const response = await fetch(`/api/events/${code}/questions/count`);

        if (!response.ok) {
            return 0;
        }

        const data = await response.json();

        return data.question_count || 0;

    } catch (error) {

        console.error("Failed to load questions count", error);

        return 0;

    }

}

// =========================================

function enableGlow() {

    document.querySelectorAll(".eventCard").forEach(card => {

        card.addEventListener("mousemove", e => {

            const rect = card.getBoundingClientRect();

            card.style.setProperty("--x", (e.clientX - rect.left) + "px");
            card.style.setProperty("--y", (e.clientY - rect.top) + "px");

        });

    });

}

// =========================================

async function deleteEvent(id) {

    try {

        await fetch(`/api/events/${id}`, {
            method: "DELETE"
        });

    } finally {

        loadEvents("teacher");

    }

}

// =========================================

async function createEvent(name) {

    name = name.trim();

    if (name === "") {

        document.getElementById("eventName").focus();

        return;

    }

    const response = await fetch("/api/events", {

        method: "POST",

        headers: {
            "Content-Type": "application/json"
        },

        body: JSON.stringify({
            title: name
        })

    });

    const data = await response.json();

    if (data.message === "event_created") {

        document
            .querySelector('#newEventModal [data-bs-dismiss="modal"]')
            .click();

        document.getElementById("eventName").value = "";

        // Обновляем текущую открытую вкладку
        const activeTab = document.querySelector(".eventsTab.active");

        loadEvents(activeTab ? activeTab.dataset.type : "teacher");

    }

}