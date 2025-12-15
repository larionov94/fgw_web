document.addEventListener('DOMContentLoaded', function () {
    // 1. Хранилище для оригинальных значений
    const originalRoleData = new Map();

    // 2. Обработчик кнопки редактирования
    document.addEventListener('click', function (e) {
        if (e.target.closest('.edit-btn')) {
            const btn = e.target.closest('.edit-btn');
            const row = btn.closest('tr');

            enableRoleEditMode(row);
        }

        if (e.target.closest('.cancel-btn')) {
            const btn = e.target.closest('.cancel-btn');
            const row = btn.closest('tr')

            disableRoleEditMode(row);
        }

        if (e.target.closest('.save-role-btn')) {
            const btn = e.target.closest('.save-role-btn');
            const row = btn.closest('tr');

            // Сохраняем изменения и игнорируем повторные клики
            if (btn.disabled) return;

            saveRoleChanges(row).catch(error => {
                console.error('Save error:', error);
                showRoleNotification('Ошибка при сохранении', 'danger');

                // Восстанавливаем кнопку при ошибке
                const saveBtn = row.querySelector('.save-role-btn');
                saveBtn.innerHTML = '<span>✓</span>';
                saveBtn.disabled = false;
            });
        }

        if (e.target.closest('.add-role-btn')) {
            const btn = e.target.closest('.add-role-btn');
            const row = btn.closest('tr');

            // Сохраняем изменения и игнорируем повторные клики
            if (btn.disabled) return;

            addChanges(row).catch(error => {
                console.error('Save error:', error);
                showRoleNotification('Ошибка при сохранении', 'danger');
            });
        }

    });

    function enableRoleEditMode(row) {
        // 1. Получаем Id из data-id атрибута строки
        const roleIdStr = row.getAttribute('data-id'); // {{ .Obj }}
        const roleId = parseInt(roleIdStr, 10);

        // 2. Получаем элементы input
        const nameInput = row.querySelector('input[name="name"]');
        const descInput = row.querySelector('input[name="description"]');

        // 3. Сохраняем оригинальные значения с сервера
        if (!originalRoleData.has(roleId)) {
            // 3.1. Берем значение из атрибутов select (они содержать оригинальные значения)
            const originalName = nameInput.getAttribute('data-original') || nameInput.value;
            const originalDesc = descInput.getAttribute('data-original') || descInput.value;

            originalRoleData.set(roleId, {
                name: originalName,
                description: originalDesc
            });
        }

        // 4. Получаем сохраненные оригинальные значения
        const original = originalRoleData.get(roleId);

        // 5. Устанавливаем текущее значение в input
        nameInput.value = original.name;
        descInput.value = original.description;

        // 6. Сохраняем текущие значения для возможности отмены
        row.dataset.originalName = original.name;
        row.dataset.originalDesc = original.description;
        row.dataset.roleId = roleId.toString();

        // 7. Показываем поля редактирования
        row.querySelectorAll('.edit-mode').forEach(el => {
            el.style.display = 'table-cell';
        });

        // 8. Скрываем поля просмотра
        row.querySelectorAll('.view-mode').forEach(el => {
            el.style.display = 'none';
        });

        // 9. Показываем кнопки сохранения/отмены
        row.querySelector('.edit-btn').style.display = 'none';
        row.querySelector('.edit-buttons').style.display = 'flex';

        // 10. Добавляем визуальные индикаторы
        row.classList.add('editing');
        row.style.backgroundColor = '#f8f9fa';

        // 11. Фокус на первое поле
        nameInput.focus();
    }

    async function addChanges(row) {
        const roleIdInput = document.getElementById('newRoleId');
        const nameInput = document.getElementById('newRoleName');
        const descriptionInput = document.getElementById('newRoleDescription');

        const roleId = roleIdInput ? parseInt(roleIdInput.value.trim(), 10) : 0;
        const name = nameInput ? nameInput.value.trim() : '';
        const description = descriptionInput ? descriptionInput.value.trim() : '';

        try {
            // 6. Отправляем запрос через Fetch API
            const response = await fetch('/admin/roles/add', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    RoleId: roleId,
                    Name: name,
                    Description: description
                })
            });

            // 7. Проверяем статус ответа
            if (!response.ok) {
                const contentType = response.headers.get('content-type');
                if (contentType && contentType.includes('application/json')) {
                    const result = await response.json();
                    new Error(result.error || `HTTP ${response.status}`);
                } else {
                    new Error(`HTTP ${response.status}`);
                }
            }

            // 9. Парсим JSON ответ
            const result = await response.json();

            // 10. Успешное обновление
            handleSuccessAdd(row, result, roleId, name, description);
        } catch (error) {
            console.error('Save error:', error);

            // Показываем уведомление об ошибке
            showRoleNotification(`Ошибка: ${error.message}`, 'danger');

            throw error; // Пробрасываем ошибку дальше
        }

    }



    function handleSuccessAdd(row, result) {
        if (result.success) {
            // Показываем уведомление об успехе
            showRoleNotification(result.message || 'Роль успешно добавлена', 'success');

            // Закрываем модальное окно
            const modal = bootstrap.Modal.getInstance(document.getElementById('addRoleModal'));
            if (modal) {
                // modal.blur()
                modal.hide();
            }
        } else {
            // Ошибка от сервера
            showRoleNotification(result.message || 'Ошибка при добавлении роли', 'danger');

        }

        // 5. Показываем уведомление
        showRoleNotification(result.message || 'Изменения успешно сохранены', 'success');

    }

    function disableRoleEditMode(row) {
        // 1. Получаем сохраненные значения для восстановления
        const originalName = row.dataset.originalName;
        const originalDesc = row.dataset.originalDesc;

        const nameInput = row.querySelector('input[name="name"]');
        const descInput = row.querySelector('input[name="description"]');

        // 2. Восстанавливаем значения в input
        if (originalName && nameInput) {
            nameInput.value = originalName;
        }

        if (originalDesc && descInput) {
            descInput.value = originalDesc;
        }

        // 3. Переключаем режимы отображения
        row.querySelectorAll('.view-mode').forEach(el => {
            el.style.display = 'table-cell';
        });

        // Скрываем ячейки редактирования
        row.querySelectorAll('.edit-mode').forEach(el => {
            el.style.display = 'none';
        });

        // 4. Переключаем кнопки
        row.querySelector('.edit-btn').style.display = 'block';
        row.querySelector('.edit-buttons').style.display = 'none';

        // 5. Убираем визуальные индикаторы
        row.classList.remove('editing');
        row.style.backgroundColor = '';

        // 6. Очищаем временные data-атрибуты
        delete row.dataset.originalName;
        delete row.dataset.originalDesc;
        delete row.dataset.roleId;
    }

    async function saveRoleChanges(row) {
        // 1. Получаем ID роли
        const roleIdStr = row.getAttribute('data-id');
        const roleId = parseInt(roleIdStr, 10);

        // 2. Получаем элементы input
        const nameInput = row.querySelector('input[name="name"]');
        const descInput = row.querySelector('input[name="description"]');

        // 3. Получаем значения
        const name = nameInput.value.trim();
        const description = descInput.value.trim();

        // 4. Валидация
        if (!name) {
            showRoleNotification('Название роли не может быть пустым', 'warning');
            nameInput.focus();
            return;
        }

        if (!description) {
            showRoleNotification('Описание роли не может быть пустым', 'warning');
            descInput.focus();
            return;
        }

        // 5. Показываем индикатор загрузки
        const saveBtn = row.querySelector('.save-role-btn');
        saveBtn.innerHTML = '<span class="spinner-border spinner-border-sm" role="status"></span>';
        saveBtn.disabled = true;

        try {
            // 6. Отправляем запрос через Fetch API
            const response = await fetch('/admin/roles/upd', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    roleId: roleId,
                    name: name,
                    description: description
                })
            });

            // 7. Проверяем статус ответа
            if (!response.ok) {
                const contentType = response.headers.get('content-type');
                if (contentType && contentType.includes('application/json')) {
                    const result = await response.json();
                    throw new Error(result.error || `HTTP ${response.status}`);
                } else {
                    throw new Error(`HTTP ${response.status}`);
                }
            }

            // 8. Парсим JSON ответ
            const result = await response.json();

            // 9. Успешное обновление
            handleRoleSuccessUpdate(row, result, roleId, name, description);

        } catch (error) {
            console.error('Save error:', error);

            // Показываем уведомление об ошибке
            showRoleNotification(`Ошибка: ${error.message}`, 'danger');

            // Восстанавливаем кнопку
            saveBtn.innerHTML = '<span>✓</span>';
            saveBtn.disabled = false;

            throw error;
        }
    }

    function handleRoleSuccessUpdate(row, result, roleId, name, description) {
        // 1. Обновление оригинальных данных
        originalRoleData.set(roleId, {
            name: name,
            description: description
        });

        row.querySelector('.forms-name').textContent = name;
        row.querySelector('.forms-desc').textContent = description;

        // 3. Обновляем значения в input и data-атрибутах
        const nameInput = row.querySelector('input[name="name"]');
        const descInput = row.querySelector('input[name="description"]');
        const updateAt = row.querySelector('.update-at');
        const updateBy = row.querySelector('.update-by');

        if (nameInput) {
            nameInput.value = name;
            nameInput.setAttribute('data-original', name);
        }

        if (descInput) {
            descInput.value = description;
            descInput.setAttribute('data-original', description);
        }

        // 4. Обновляем дату и пользователя (если пришли с сервера)
        updateAt.textContent = result.updatedAt;
        updateBy.textContent = result.updatedBy;

        // 5. Выходим из режима редактирования
        disableRoleEditMode(row);

        // 6. Показываем уведомление
        showRoleNotification(result.message || 'Роль успешно обновлена', 'success');

        // 7. Восстанавливаем кнопку сохранения
        const saveBtn = row.querySelector('.save-role-btn');
        saveBtn.innerHTML = '<span>✓</span>';
        saveBtn.disabled = false;
    }

    function showRoleNotification(message, type) {
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