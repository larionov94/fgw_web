// 1. Защита от навигации по истории для защищенных страниц
(function () {
    let sessionCheckInterval;
    let isCheckingSession = false;
    let lastActivity = Date.now();

    // 2. Функция проверки сессии
    function checkSession() {
        if (isCheckingSession) return;

        isCheckingSession = true;

        fetch('/api/session-check', {
            method: 'HEAD',
            credentials: 'include'
        })
            .then(function (response) {
                if (!response.ok) {
                    // 2.1. Сессия невалидна - logout
                    console.log('Сессия истекла, выполняется выход...');
                    window.location.href = '/logout';
                    return;
                }

                // 2.2. Проверяем заголовки если нужно
                const sessionStatus = response.headers.get('Session-Status');
                if (sessionStatus && sessionStatus !== 'active') {
                    console.log('Статус сессии:', sessionStatus);
                    window.location.href = '/logout';
                }

                // 2.3. Обновляем время последней активности
                lastActivity = Date.now();
            })
            .catch(function (error) {
                console.error('Ошибка проверки сессии:', error);
                // 2.4 При ошибке сети не делаем logout сразу можно попробовать снова через некоторое время
            })
            .finally(function () {
                isCheckingSession = false;
            });
    }

    // 1. Функция проверки не активности пользователя
    function checkInactivity() {
        const now = Date.now();
        const inactiveTime = now - lastActivity;

        // 2. Если неактивны более 10 минут - показываем предупреждение
        if (inactiveTime > 10 * 60 * 1000) {
            const userConfirmed = confirm('Вы неактивны более 10 минут. Сессия будет завершена через 5 минут.');
            if (userConfirmed) {
                // 2.1. Сброс времени активности
                lastActivity = Date.now();
                // 2.2. Принудительная проверка сессии
                checkSession();
            }
        }

        // 3. Если неактивны более 15 минут - принудительный выход
        if (inactiveTime > 15 * 60 * 1000) {
            alert('Сессия завершена из-за не активности.');
            window.location.href = '/logout';
        }
    }

    // 1. Обновление времени активности при действиях пользователя
    function updateActivity() {
        lastActivity = Date.now();
    }

    // 2. События активности пользователя
    ['click', 'keypress', 'mousemove', 'scroll', 'touchstart'].forEach(function (eventName) {
        document.addEventListener(eventName, updateActivity, {passive: true});
    });

    // 3. Инициализация
    function initSessionMonitoring() {
        // 3.1. Первая проверка через 1 минуту после загрузки
        setTimeout(checkSession, 60000);

        // 3.2. Затем каждые 5 минут
        sessionCheckInterval = setInterval(checkSession, 300000);

        // 3.3. Проверка не активности каждую минуту
        setInterval(checkInactivity, 60000);
    }

    // 4. Заменяем состояние в истории
    if (window.history.replaceState) {
        history.replaceState({
            authed: true,
            timestamp: new Date().toISOString()
        }, '', window.location.pathname);
    }

    // 5. Обработчик навигации
    window.addEventListener('popstate', function (event) {
        if (!event.state || event.state.authed !== true) {
            // Пытаемся уйти с защищенной страницы
            window.location.href = '/logout';
        }
    });

    // 6. Защита от кэширования
    window.addEventListener('pageshow', function (event) {
        if (event.persisted) {
            window.location.reload();
        }
    });

    // 7. Запускаем мониторинг сессии
    initSessionMonitoring();

    // 8. Очистка при закрытии страницы
    window.addEventListener('beforeunload', function () {
        if (sessionCheckInterval) {
            clearInterval(sessionCheckInterval);
        }
    });
})();
