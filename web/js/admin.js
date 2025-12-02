// Защита от навигации по истории для защищенных страниц
(function () {
    let sessionCheckInterval;
    let isCheckingSession = false;
    let lastActivity = Date.now();

    // Функция проверки сессии
    function checkSession() {
        if (isCheckingSession) return;

        isCheckingSession = true;

        fetch('/api/session-check', {
            method: 'HEAD',
            credentials: 'include'
        })
            .then(function(response) {
                if (!response.ok) {
                    // Сессия невалидна - logout
                    console.log('Сессия истекла, выполняется выход...');
                    window.location.href = '/logout';
                    return;
                }

                // Проверяем заголовки если нужно
                const sessionStatus = response.headers.get('Session-Status');
                if (sessionStatus && sessionStatus !== 'active') {
                    console.log('Статус сессии:', sessionStatus);
                    window.location.href = '/logout';
                }

                // Обновляем время последней активности
                lastActivity = Date.now();
            })
            .catch(function(error) {
                console.error('Ошибка проверки сессии:', error);
                // При ошибке сети не делаем logout сразу
                // можно попробовать снова через некоторое время
            })
            .finally(function() {
                isCheckingSession = false;
            });
    }

    // Функция проверки не активности пользователя
    function checkInactivity() {
        const now = Date.now();
        const inactiveTime = now - lastActivity;

        // Если неактивны более 10 минут - показываем предупреждение
        if (inactiveTime > 10 * 60 * 1000) {
            const userConfirmed = confirm('Вы неактивны более 10 минут. Сессия будет завершена через 5 минут.');
            if (userConfirmed) {
                // Сброс времени активности
                lastActivity = Date.now();
                // Принудительная проверка сессии
                checkSession();
            }
        }

        // Если неактивны более 15 минут - принудительный выход
        if (inactiveTime > 15 * 60 * 1000) {
            alert('Сессия завершена из-за не активности.');
            window.location.href = '/logout';
        }
    }

    // Обновление времени активности при действиях пользователя
    function updateActivity() {
        lastActivity = Date.now();
    }

    // События активности пользователя
    ['click', 'keypress', 'mousemove', 'scroll', 'touchstart'].forEach(function(eventName) {
        document.addEventListener(eventName, updateActivity, { passive: true });
    });

    // Инициализация
    function initSessionMonitoring() {
        // Первая проверка через 1 минуту после загрузки
        setTimeout(checkSession, 60000);

        // Затем каждые 5 минут
        sessionCheckInterval = setInterval(checkSession, 300000);

        // Проверка не активности каждую минуту
        setInterval(checkInactivity, 60000);
    }

    // Заменяем состояние в истории
    if (window.history.replaceState) {
        history.replaceState({
            authed: true,
            timestamp: new Date().toISOString()
        }, '', window.location.pathname);
    }

    // Обработчик навигации
    window.addEventListener('popstate', function (event) {
        if (!event.state || event.state.authed !== true) {
            // Пытаемся уйти с защищенной страницы
            window.location.href = '/logout';
        }
    });

    // Защита от кэширования
    window.addEventListener('pageshow', function (event) {
        if (event.persisted) {
            window.location.reload();
        }
    });

    // Запускаем мониторинг сессии
    initSessionMonitoring();

    // Очистка при закрытии страницы
    window.addEventListener('beforeunload', function() {
        if (sessionCheckInterval) {
            clearInterval(sessionCheckInterval);
        }
    });

})();