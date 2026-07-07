// Логика переключения темы оформления
document.addEventListener('DOMContentLoaded', () => {
    const themeToggleBtn = document.getElementById('themeToggle');
    const htmlElement = document.documentElement;

    if (!themeToggleBtn) return;

    const savedTheme = localStorage.getItem('theme') || 'dark';
    applyTheme(savedTheme, themeToggleBtn, htmlElement);

    themeToggleBtn.addEventListener('click', () => {
        const currentTheme = htmlElement.getAttribute('data-theme') || 'dark';
        const newTheme = currentTheme === 'dark' ? 'light' : 'dark';

        applyTheme(newTheme, themeToggleBtn, htmlElement);
    });
});

function applyTheme(theme, button, htmlElement) {
    const normalizedTheme = theme === 'light' ? 'light' : 'dark';

    htmlElement.setAttribute('data-theme', normalizedTheme);
    htmlElement.setAttribute('data-bs-theme', normalizedTheme);
    localStorage.setItem('theme', normalizedTheme);
    document.dispatchEvent(new CustomEvent('themechange', { detail: { theme: normalizedTheme } }));

    if (!button) return;

    button.innerHTML = normalizedTheme === 'dark'
        ? '<svg viewBox="0 0 24 24" aria-hidden="true"><path d="M12 3v2m0 14v2m9-9h-2M5 12H3m9-6.5 1.4 1.4M9.6 14.4 8.2 15.8m8.6 0-1.4 1.4M9.6 9.6 8.2 8.2m6.8 0 1.4-1.4m-1.4 8.6 1.4 1.4" /></svg>'
        : '<svg viewBox="0 0 24 24" aria-hidden="true"><path d="M21 12.8A9 9 0 1 1 11.2 3a7 7 0 0 0 9.8 9.8Z" /></svg>';

    button.setAttribute(
        'aria-label',
        normalizedTheme === 'dark' ? 'Переключить на светлую тему' : 'Переключить на тёмную тему'
    );
}
