// ================================
// THEME
// ================================

const themeButton = document.getElementById("themeToggle");
const html = document.documentElement;

function updateThemeButton(){

    if(html.dataset.theme === "dark"){
        themeButton.textContent = "☀️";
    }
    else{
        themeButton.textContent = "🌙";
    }

}

const savedTheme = localStorage.getItem("theme");

if(savedTheme){
    html.dataset.theme = savedTheme;
}

updateThemeButton();

if(themeButton){

    themeButton.addEventListener("click",()=>{

        html.dataset.theme =
            html.dataset.theme === "dark"
            ? "light"
            : "dark";

        localStorage.setItem(
            "theme",
            html.dataset.theme
        );

        updateThemeButton();

    });

}

// ================================
// INPUT EFFECT
// ================================

document.querySelectorAll("input").forEach(input=>{

    input.addEventListener("focus",()=>{

        input.parentElement.classList.add("focused");

    });

    input.addEventListener("blur",()=>{

        if(input.value===""){

            input.parentElement.classList.remove("focused");

        }

    });

});

// ================================
// LOGIN FORM
// ================================

const form = document.getElementById("loginForm");

form.addEventListener("submit",(e)=>{

    e.preventDefault();

    hideError();

    const email =
        document.getElementById("email").value.trim();

    const password =
        document.getElementById("password").value;

    tryToSend(email,password);

});

// ================================
// VALIDATION
// ================================

function tryToSend(email,password){

    const regex =
    /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/;

    if(password.length < 6){

        showError("Пароль слишком короткий");

        return;

    }

    if(email === "" || !regex.test(email)){

        showError("Почта неверно введена");

        return;

    }

    sendLogs(email,password);

}

// ================================
// REQUEST
// ================================

async function sendLogs(email,password){

    try{

        const response = await fetch("/api/login",{

            method:"POST",

            headers:{
                "Content-Type":"application/json"
            },

            body:JSON.stringify({

                email,
                password

            })

        });

        const json = await response.json();

        if(json.message === "login_success"){

            location.href="/events";

            return;

        }

        switch(json.error){

            case "bad_credentials":

                showError(
                    "Неправильная почта или пароль!"
                );

                break;

            case "server_error":

                showError(
                    "Сервер не отвечает..."
                );

                break;

            default:

                showError(
                    "Неизвестная ошибка"
                );

        }

    }
    catch(error){

        showError(
            "Ошибка соединения с сервером"
        );

    }

}

// ================================
// ERROR
// ================================

const errorBox =
    document.getElementById("loginError");

const errorMessage =
    document.getElementById("errorMessage");

function showError(message){

    errorMessage.textContent = message;

    errorBox.classList.add("show");

}

function hideError(){

    errorBox.classList.remove("show");

}

// ================================
// CARD INTRO
// ================================

window.addEventListener("load",()=>{

    const card =
        document.querySelector(".loginCard");

    if(card){

        card.classList.add("loaded");

    }

});

// ================================
// BUTTON RIPPLE
// ================================

const btn =
    document.querySelector(".loginButtonMain");

if(btn){

    btn.addEventListener("mousemove",(e)=>{

        const rect =
            btn.getBoundingClientRect();

        btn.style.setProperty(
            "--x",
            `${e.clientX-rect.left}px`
        );

        btn.style.setProperty(
            "--y",
            `${e.clientY-rect.top}px`
        );

    });

}