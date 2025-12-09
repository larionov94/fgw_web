document.addEventListener('DOMContentLoaded', function () {
    // Хранилище для оригинальных значений из сервера
    const originalData = new Map();

    // Обработчик кнопки редактирования
    document.addEventListener('click', function (e) {
        if (e.target.closest('.edit-btn')) {
            const btn = e.target.closest('.edit-btn');
            const row = btn.closest('tr');
            enableEditMode(row);
        }

        // Обработчик кнопки сохранения
        if (e.target.closest('.save-btn')) {
            const btn = e.target.closest('.save-btn');
            const row = btn.closest('tr');

            // Сохраняем изменения и игнорируем повторные клики
            if (btn.disabled) return;

            saveChanges(row).catch(error => {
                console.error('Save error:', error);
                showNotification('Ошибка при сохранении', 'danger');

                // Восстанавливаем кнопку при ошибке
                const saveBtn = row.querySelector('.save-btn');
                saveBtn.innerHTML = '<span>✓</span>';
                saveBtn.disabled = false;
            });
        }

        if (e.target.closest('.cancel-btn')) {
            const btn = e.target.closest('.cancel-btn');
            const row = btn.closest('tr');
            disableEditMode(row);
        }
    });

    function enableEditMode(row) {
        // Получаем ID из data-id атрибута строки
        const performerIdStr = row.getAttribute('data-id');
        const performerId = parseInt(performerIdStr, 10);

        if (isNaN(performerId)) {
            console.error('Invalid performer ID:', performerIdStr);
            return;
        }

        // Получаем элементы select
        const formSelect = row.querySelector('.role-forms-select');
        const fgwSelect = row.querySelector('.role-fgw-select');

        // Если это первое редактирование, сохраняем оригинальные значения из сервера
        if (!originalData.has(performerId)) {
            // Берем значения из атрибутов select'ов (они содержат оригинальные значения из сервера)
            const originalFormsValue = formSelect.getAttribute('data-original') || formSelect.value;
            const originalFgwValue = fgwSelect.getAttribute('data-original') || fgwSelect.value;

            // Находим текст для этих значений
            let originalFormsText = '';
            let originalFgwText = '';

            for (let option of formSelect.options) {
                if (option.value === originalFormsValue) {
                    originalFormsText = option.text.trim();
                    break;
                }
            }

            for (let option of fgwSelect.options) {
                if (option.value === originalFgwValue) {
                    originalFgwText = option.text.trim();
                    break;
                }
            }

            originalData.set(performerId, {
                formsValue: originalFormsValue,
                fgwValue: originalFgwValue,
                formsText: originalFormsText,
                fgwText: originalFgwText
            });
        }

        // Получаем сохраненные оригинальные значения
        const original = originalData.get(performerId);

        // Устанавливаем текущие значения в select'ах
        formSelect.value = original.formsValue;
        fgwSelect.value = original.fgwValue;

        // Сохраняем текущие значения для возможной отмены
        row.dataset.originalFormsValue = original.formsValue;
        row.dataset.originalFgwValue = original.fgwValue;
        row.dataset.performerId = performerId.toString();

        // Показываем поля редактирования
        row.querySelectorAll('.edit-mode').forEach(el => {
            el.style.display = 'table-cell';
        });

        // Скрываем поля просмотра
        row.querySelectorAll('.view-mode').forEach(el => {
            el.style.display = 'none';
        });

        // Показываем кнопки сохранения/отмены
        row.querySelector('.edit-btn').style.display = 'none';
        row.querySelector('.edit-buttons').style.display = 'flex';

        // Добавляем визуальные индикаторы
        row.classList.add('editing');
        row.style.backgroundColor = '#f8f9fa';
    }

    function disableEditMode(row) {
        // Получаем сохраненные значения для восстановления
        const originalFormsValue = row.dataset.originalFormsValue;
        const originalFgwValue = row.dataset.originalFgwValue;

        const formSelect = row.querySelector('.role-forms-select');
        const fgwSelect = row.querySelector('.role-fgw-select');

        // Восстанавливаем значения в select'ах
        if (originalFormsValue && formSelect) {
            formSelect.value = originalFormsValue;
        }

        if (originalFgwValue && fgwSelect) {
            fgwSelect.value = originalFgwValue;
        }

        // Скрываем поля редактирования
        row.querySelectorAll('.edit-mode').forEach(el => {
            el.style.display = 'none';
        });

        // Показываем поля просмотра
        row.querySelectorAll('.view-mode').forEach(el => {
            el.style.display = 'table-cell';
        });

        // Показываем кнопку редактирования
        row.querySelector('.edit-btn').style.display = 'block';
        row.querySelector('.edit-buttons').style.display = 'none';

        // Убираем визуальные индикаторы
        row.classList.remove('editing');
        row.style.backgroundColor = '';

        // Очищаем временные data-атрибуты
        delete row.dataset.originalFormsValue;
        delete row.dataset.originalFgwValue;
        delete row.dataset.performerId;
    }

    async function saveChanges(row) {
        // Получаем performerId из строки
        const performerIdStr = row.getAttribute('data-id');
        const performerId = parseInt(performerIdStr, 10);

        if (isNaN(performerId)) {
            showNotification('Ошибка: неверный ID сотрудника', 'danger');
            return;
        }

        // Получаем элементы select
        const formsSelect = row.querySelector('.role-forms-select');
        const fgwSelect = row.querySelector('.role-fgw-select');

        // Получаем текстовые значения выбранных опций
        const selectedFormsText = formsSelect.options[formsSelect.selectedIndex].text;
        const selectedFgwText = fgwSelect.options[fgwSelect.selectedIndex].text;

        // Преобразуем значения в числа
        const idRoleAForms = parseInt(formsSelect.value, 10);
        const idRoleAFGW = parseInt(fgwSelect.value, 10);

        // Валидация
        if (isNaN(idRoleAForms) || isNaN(idRoleAFGW)) {
            showNotification('Ошибка: неверные значения ролей', 'danger');
            throw new Error('Invalid role values');
        }

        // Показываем индикатор загрузки
        const saveBtn = row.querySelector('.save-btn');
        saveBtn.innerHTML = '<span class="spinner-border spinner-border-sm" role="status"></span>';
        saveBtn.disabled = true;

        try {
            // Отправляем запрос через Fetch API
            const response = await fetch('/admin/performers/upd', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    performerId: performerId,
                    idRoleAForms: idRoleAForms,
                    idRoleAFGW: idRoleAFGW
                })
            });

            // Проверяем статус ответа
            if (!response.ok) {
                const contentType = response.headers.get('content-type');
                if (contentType && contentType.includes('application/json')) {
                    const result = await response.json();
                    throw new Error(result.error || `HTTP ${response.status}`);
                } else {
                    throw new Error(`HTTP ${response.status}`);
                }
            }

            // Парсим JSON ответ
            const result = await response.json();

            // Успешное обновление
            handleSuccessUpdate(row, result, selectedFormsText, selectedFgwText, performerId, idRoleAForms, idRoleAFGW);

        } catch (error) {
            console.error('Save error:', error);

            // Показываем уведомление об ошибке
            showNotification(`Ошибка: ${error.message}`, 'danger');

            // Восстанавливаем кнопку
            saveBtn.innerHTML = '<span>✓</span>';
            saveBtn.disabled = false;

            throw error; // Пробрасываем ошибку дальше
        }
    }

    function handleSuccessUpdate(row, result, selectedFormsText, selectedFgwText, performerId, idRoleAForms, idRoleAFGW) {
        // ОБНОВЛЯЕМ хранилище оригинальных данных
        originalData.set(performerId, {
            formsValue: idRoleAForms.toString(),
            fgwValue: idRoleAFGW.toString(),
            formsText: selectedFormsText,
            fgwText: selectedFgwText
        });

        // Обновляем отображение ролей
        row.querySelector('.forms-role .badge').textContent = selectedFormsText;
        row.querySelector('.fgw-role .badge').textContent = selectedFgwText;

        // Обновляем значения в select'ах
        const formSelect = row.querySelector('.role-forms-select');
        const fgwSelect = row.querySelector('.role-fgw-select');

        if (formSelect) {
            formSelect.value = idRoleAForms.toString();
            // Обновляем атрибут data-original
            formSelect.setAttribute('data-original', idRoleAForms.toString());
        }

        if (fgwSelect) {
            fgwSelect.value = idRoleAFGW.toString();
            fgwSelect.setAttribute('data-original', idRoleAFGW.toString());
        }
        

        // Выходим из режима редактирования
        disableEditMode(row);

        // Показываем уведомление
        showNotification(result.message || 'Изменения успешно сохранены', 'success');

        // Восстанавливаем кнопку сохранения
        const saveBtn = row.querySelector('.save-btn');
        saveBtn.innerHTML = '<span>✓</span>';
        saveBtn.disabled = false;
    }

    // Остальные функции остаются без изменений...
    function handleErrorResponse(status, errorMessage) {
        console.error('Server error:', status, errorMessage);

        switch (status) {
            case 400:
                showNotification('Неверные данные: ' + errorMessage, 'danger');
                break;
            case 401:
                showNotification('Сессия истекла. Перенаправление...', 'warning');
                setTimeout(() => {
                    window.location.href = '/login';
                }, 2000);
                break;
            case 403:
                showNotification('У вас нет прав для этого действия', 'danger');
                break;
            case 404:
                showNotification('Сотрудник не найден', 'danger');
                break;
            case 500:
                showNotification('Ошибка сервера: ' + errorMessage, 'danger');
                break;
            default:
                showNotification('Ошибка: ' + errorMessage, 'danger');
        }
    }

    function formatDate(dateString) {
        try {
            const date = new Date(dateString);
            return date.toLocaleDateString('ru-RU') + ' ' +
                date.toLocaleTimeString('ru-RU', {hour: '2-digit', minute: '2-digit'});
        } catch (e) {
            return dateString; // Возвращаем как есть если не удалось распарсить
        }
    }

    function showNotification(message, type) {
        // Удаляем существующие уведомления
        document.querySelectorAll('.alert.position-fixed').forEach(el => el.remove());

        // Создаем элемент уведомления
        const notification = document.createElement('div');
        notification.className = `alert alert-${type} alert-dismissible fade show position-fixed`;
        notification.style.cssText = `
            top: 20px;
            right: 20px;
            z-index: 9999;
            min-width: 300px;
            max-width: 500px;
            box-shadow: 0 2px 10px rgba(0,0,0,0.1);
        `;
        notification.innerHTML = `
            <div class="d-flex align-items-center">
                <div class="flex-grow-1">${message}</div>
                <button type="button" class="btn-close ms-2" data-bs-dismiss="alert"></button>
            </div>
        `;

        // Добавляем на страницу
        document.body.appendChild(notification);

        // Автоматически удаляем через 5 секунд
        setTimeout(() => {
            if (notification.parentNode) {
                notification.remove();
            }
        }, 5000);
    }
});