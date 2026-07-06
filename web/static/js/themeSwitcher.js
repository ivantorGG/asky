// Логика переключения темы оформления
document.addEventListener('DOMContentLoaded', () => {
    const themeToggleBtn = document.getElementById('themeToggle');
    const htmlElement = document.documentElement;

    // 1. Проверяем, есть ли сохраненная тема в памяти браузера
    const savedTheme = localStorage.getItem('theme') || 'light';
    
    // 2. Применяем сохраненную тему и меняем текст кнопки
    htmlElement.setAttribute('data-bs-theme', savedTheme);
    updateButtonText(themeToggleBtn, savedTheme);

    // 3. Отслеживаем клик по кнопке
    themeToggleBtn.addEventListener('click', () => {
        const currentTheme = htmlElement.getAttribute('data-bs-theme');
        // Меняем тему на противоположную
        const newTheme = currentTheme === 'dark' ? 'light' : 'dark';
        
        htmlElement.setAttribute('data-bs-theme', newTheme);
        localStorage.setItem('theme', newTheme); // Сохраняем выбор
        updateButtonText(themeToggleBtn, newTheme);
    });
});

// Функция для изменения текста и иконки на кнопке
function updateButtonText(button, theme) {
    if (!button) return;
    if (theme === 'dark') {
        button.innerHTML = '☀️';
    } else {
        button.innerHTML = '🌙';
    }
}
