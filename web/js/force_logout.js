(function () {
    document.getElementById('status').textContent = 'Очистка истории браузера...';

    // 1. Заменяем текущее состояние
    if (window.history.replaceState) {
        try {
            history.replaceState(null, '', window.location.href);
        } catch (e) {
            console.error('Ошибка замены состояния:', e);
        }
    }

    // 2. Добавляем временное состояние для логина
    if (window.history.pushState) {
        try {
            history.pushState({authed: false}, '', '/');
        } catch (e) {
            console.error('Ошибка добавления состояния:', e);
        }
    }

    // 3. Возвращаемся назад, удаляя текущую страницу из истории
    try {
        history.back();
    } catch (e) {
        console.error('Ошибка навигации назад:', e);
        // Если не удалось - сразу редирект
        setTimeout(function () {
            window.location.replace('/login');
        }, 100);
        return;
    }

    // 4. Редирект на логин без добавления в историю
    setTimeout(function () {
        document.getElementById('status').textContent = 'Перенаправление на страницу входа...';
        window.location.replace('/login');
    }, 100);

    // 5. Фолбэк на случай если что-то пошло не так
    setTimeout(function () {
        window.location.replace('/login');
    }, 2000);

})();