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
// REGISTER FORM
// ================================

const form = document.getElementById("registerForm");

form.addEventListener("submit",(e)=>{

    e.preventDefault();

    const email =
        document.getElementById("email").value.trim();

    const password =
        document.getElementById("password").value;
    
    const name =
        document.getElementById("name").value;


    tryToSend(email,password,name);

});


// ================================
// VALIDATION
// ================================

function tryToSend(email,password,name){

    const regex =
    /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/;

    hideError();

    if(name.length === 0){
        showError("Введите имя");
        return;
    }

    if(password.length < 6){

        showError(
            "Пароль слишком короткий (минимум 6 символов)"
        );

        return;

    }

    if(email === "" || !regex.test(email)){

        showError(
            "Некорректная почта"
        );

        return;

    }

    sendRegs(email,password,name);

}


// ================================
// REQUEST
// ================================

async function sendRegs(email,password,name){

    try{

        const response = await fetch("/api/register",{

            method:"POST",

            headers:{
                "Content-Type":"application/json"
            },

            body:JSON.stringify({
                email: email,
                password: password,
                name: name
            })

        });


        const json = await response.json();


        if(json.message === "registration_success"){

            location.href="/login";

            return;

        }


        switch(json.error){

            case "server_error":

                showError(
                    "Сервер недоступен"
                );

                break;


            case "bad_request":

                showError(
                    "Некорректный запрос"
                );

                break;


            case "invalid_input":

                showError(
                    "Некорректные данные"
                );

                break;


            case "email_already_exists":

                showError(
                    "Такая почта уже зарегистрирована"
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
        document.querySelector(".registerCard");

    if(card){

        card.classList.add("loaded");

    }

});


// ================================
// BUTTON RIPPLE
// ================================

const btn =
    document.querySelector(".registerButton");


if(btn){

    btn.addEventListener("mousemove",(e)=>{

        const rect =
            btn.getBoundingClientRect();

        btn.style.setProperty(
            "--x",
            `${e.clientX - rect.left}px`
        );

        btn.style.setProperty(
            "--y",
            `${e.clientY - rect.top}px`
        );

    });

}