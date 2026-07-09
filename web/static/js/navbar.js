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
        const isMobile = window.innerWidth <= 768;
        if (email.length > 8 && isMobile)  {
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
    const logoutButton = document.getElementById('logoutButton');

    if (!authLink || !logoutButton) {
        return;
    }

    if (email) {

        authLink.href = '/events';
        authLink.textContent = 'Открыть мероприятия';
        authLink.classList.remove('loginButton');
        authLink.classList.add('eventsLink');

        logoutButton.style.display = 'flex';

    } else {

        authLink.href = '/login';
        authLink.textContent = 'Войти';
        authLink.classList.remove('eventsLink');
        authLink.classList.add('loginButton');

        logoutButton.style.display = 'none';

    }

}