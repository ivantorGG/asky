// ================================
// THEME
// ================================

const themeButton = document.getElementById("themeToggle");
const html = document.documentElement;

const moonIcon = `
<svg viewBox="0 0 24 24" width="22" height="22">
<path fill="currentColor"
d="M21 12.79A9 9 0 0 1 11.21 3
a7 7 0 1 0 9.79 9.79z"/>
</svg>
`;

    if (password.length < 6){
        showError('Пароль слишком короткий (минимальная длина - 6 символов)')
    }
    else if (email === '' || !regex.test(email)){
        showError('Почта неверно введена')
    }
    else{
        sendRegs(email, password);
    }



updateThemeButton();

themeButton.addEventListener("click",()=>{

    html.dataset.theme =
        html.dataset.theme === "dark"
        ? "light"
        : "dark";

    localStorage.setItem("theme",html.dataset.theme);

    updateThemeButton();

});

const savedTheme = localStorage.getItem("theme");

if(savedTheme){

    html.dataset.theme = savedTheme;

    updateThemeButton();

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
// VALIDATION
// ================================

function tryToSend(email,password){

    const regex =
    /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/;

    hideError();

    if(password.length<6){

        showError("Password must contain at least 6 characters");

        return;

    }

    if(email==="" || !regex.test(email)){

        showError("Incorrect email");

        return;

    }

    sendRegs(email,password);

}

// ================================
// REQUEST
// ================================

async function sendRegs(email,password){

    const response = await fetch("/api/register",{

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

    if(json.message==="registration_success"){

        location.href="/login";

        return;

    }

    switch(json.error){

        case "server_error":

            showError("Server unavailable");

            break;

        case "bad_request":

            showError("Bad request");

            break;

        case "invalid_input":

            showError("Invalid input");

            break;

        case "email_already_exists":

            showError("Email already exists");

            break;

        default:

            showError("Unknown error");

    }

}

// ================================
// ERROR
// ================================

const errorBox = document.getElementById("loginError");

function showError(message){

    document.getElementById("errorMessage").textContent = message;

    errorBox.classList.add("show");

}

function hideError(){

    errorBox.classList.remove("show");

}

// ================================
// CARD INTRO
// ================================

window.addEventListener("load",()=>{

    document.querySelector(".registerCard")
        .classList.add("loaded");

});

// ================================
// BUTTON RIPPLE
// ================================

const btn = document.querySelector(".registerButton");

btn.addEventListener("mousemove",(e)=>{

    const rect = btn.getBoundingClientRect();

    btn.style.setProperty(
        "--x",
        `${e.clientX-rect.left}px`
    );

    btn.style.setProperty(
        "--y",
        `${e.clientY-rect.top}px`
    );

});