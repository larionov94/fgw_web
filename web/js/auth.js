// Глобальные обработчики для защиты авторизации
document.addEventListener('DOMContentLoaded', function() {
    // Для формы логина
    const loginForm = document.getElementById('loginForm');
    if (loginForm) {
        loginForm.addEventListener('submit', function(e) {
            // Можно добавить дополнительную валидацию
            const performerId = document.getElementById('performerId');
            const password = document.getElementById('performerPassword');

            if (!performerId.value || !password.value) {
                e.preventDefault();
                alert('Заполните все поля');
                return false;
            }

            // Блокируем повторную отправку
            const submitBtn = this.querySelector('input[type="submit"]');
            if (submitBtn) {
                submitBtn.disabled = true;
                submitBtn.value = 'Вход...';
            }

            return true;
        });
    }

    // Защита от вставки пароля
    const passwordFields = document.querySelectorAll('input[type="password"]');
    passwordFields.forEach(function(field) {
        field.addEventListener('copy', function(e) {
            e.preventDefault();
            return false;
        });

        field.addEventListener('paste', function(e) {
            e.preventDefault();
            return false;
        });
    });
});

// Блокировка кэширования и управление историей
window.addEventListener('pageshow', function(event) {
    // Если страница загружена из кэша - перезагружаем
    if (event.persisted) {
        window.location.reload();
    }
});

// Очищаем историю при загрузке
if (window.history.replaceState) {
    history.replaceState(null, '', '/login');
    // Блокируем кнопку "назад" на странице логина
    history.pushState(null, '', window.location.href);
}

// Предотвращаем навигацию назад/вперед
window.addEventListener('popstate', function() {
    if (window.location.pathname === '/login') {
        history.pushState(null, '', window.location.href);
    }
});