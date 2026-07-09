document.addEventListener('DOMContentLoaded', () => {
    showEmail();
});

async function showEmail() {
    const response = await fetch('/api/email', {
        method: 'GET',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include'
    });

    if (!response.ok) {
        updateAuthNav(null);
        return;
    }

    const data = await response.json();
    const email = data?.email;

    const emailEl = document.getElementById('email');
    if (emailEl && email) {
        if (email.length > 8) {
            let tEmail = email.slice(0, 8) + "...";
            emailEl.textContent = tEmail;
        } else {

            emailEl.textContent = email;
        }
    }

    updateAuthNav(email || null);
}

function updateAuthNav(email) {
    const authLink = document.getElementById('authNavLink');

    if (!authLink) {
        return;
    }

    if (email) {
        authLink.href = '/events';
        authLink.textContent = 'Открыть мероприятия';
        authLink.classList.remove('loginButton');
        authLink.classList.add('eventsLink');
    } else {
        authLink.href = '/login';
        authLink.textContent = 'Войти';
        authLink.classList.remove('eventsLink');
        authLink.classList.add('loginButton');
    }
}