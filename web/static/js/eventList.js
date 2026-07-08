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

    loadEvents();

});

// =========================================

async function loadEvents(){

    const container = document.getElementById("events");

    const response = await fetch("/api/events/teacher");

    const data = await response.json();

    if(data.error === "unauthorized"){

        location.href="/login";

        return;

    }

    const events = data;

    container.innerHTML = "";

    // ----------------------------

    if(events.length===0){

        container.innerHTML=`

        <div class="col-lg-9">

            <div class="emptyState reveal active">

                <i class="bi bi-calendar2-plus"></i>

                <h2>

                    Пока нет мероприятий

                </h2>

                <p>

                    Создайте своё первое мероприятие,
                    чтобы студенты могли подключиться.

                </p>

            </div>

        </div>

        `;

        return;

    }

    // ----------------------------

    const cardsHtml = await Promise.all(

        events.map(async event => {

            const questionCount = await getQuestionsCount(event.code);

            return `

<div class="col-lg-10 reveal active">

<div class="eventCard">

<a

href="/events/${event.code}/teacher"

class="stretched-link">

</a>

<div class="eventBody">

<div class="eventInfo">

<h2 class="eventTitle">

${event.title}

</h2>

<div class="eventSubtitle">

Вопросов: ${questionCount}

</div>

</div>

<button

class="deleteButton"

onclick="event.preventDefault();event.stopPropagation();deleteEvent('${event.code}')">

<i class="bi bi-trash"></i>

Удалить

</button>

</div>

</div>

</div>

`;

        })

    );

    container.innerHTML = cardsHtml.join("");

    enableGlow();

}

// =========================================

async function getQuestionsCount(code){

    try{

        const response = await fetch(`/api/events/${code}/questions/count`);

        if(!response.ok){

            return 0;

        }

        const data = await response.json();

        return data.question_count || 0;

    }

    catch(error){

        console.error("Failed to load questions count", error);

        return 0;

    }

}

// =========================================

function enableGlow(){

    document.querySelectorAll(".eventCard").forEach(card=>{

        card.addEventListener("mousemove",e=>{

            const rect=card.getBoundingClientRect();

            card.style.setProperty("--x",(e.clientX-rect.left)+"px");

            card.style.setProperty("--y",(e.clientY-rect.top)+"px");

        });

    });

}

// =========================================

async function deleteEvent(id){

    const card=document.activeElement;

    try{

        await fetch(`/api/events/${id}`,{

            method:"DELETE"

        });

    }

    finally{

        loadEvents();

    }

}

// =========================================

async function createEvent(name){

    name=name.trim();

    if(name===""){

        document.getElementById("eventName").focus();

        return;

    }

    const response=await fetch("/api/events",{

        method:"POST",

        headers:{

            "Content-Type":"application/json"

        },

        body:JSON.stringify({

            title:name

        })

    });

    const data=await response.json();

    if(data.message==="event_created"){

        document.querySelector(

            '#newEventModal [data-bs-dismiss="modal"]'

        ).click();

        document.getElementById("eventName").value="";

        loadEvents();

    }

}