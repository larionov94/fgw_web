document.addEventListener('DOMContentLoaded', function () {
    // 1. Хранилище для оригинальных значений
    const originalData = new Map();

    // 2. Обработчик кнопки редактирования
    document.addEventListener('click', function (e) {
        if (e.target.closest('.edit-btn')) {
            const btn = e.target.closest('.edit-btn');
            const row = btn.closest('tr');

            enablePerformersEditMode(row);
        }

        if (e.target.closest('.cancel-btn')) {
            const btn = e.target.closest('.cancel-btn');
            const row = btn.closest('tr')

            disablePerformersEditMode(row);
        }

        if (e.target.closest('.save-btn')) {
            const btn = e.target.closest('.save-btn');
            const row = btn.closest('tr');

            // Сохраняем изменения и игнорируем повторные клики
            if (btn.disabled) return;

            saveChanges(row).catch(error => {
                console.error('Save error:', error);
                showPerformersNotification('Ошибка при сохранении', 'danger');

                // Восстанавливаем кнопку при ошибке
                const saveBtn = row.querySelector('.save-btn');
                saveBtn.innerHTML = '<span>✓</span>';
                saveBtn.disabled = false;
            });
        }



    });

    function enablePerformersEditMode(row) {
        // 1. Получаем Id из data-id атрибута строки
        const performerIdStr = row.getAttribute('data-id'); // {{ .Obj }}
        const performerId = parseInt(performerIdStr, 10);

        // 2. Получаем элементы select
        const formSelect = row.querySelector('.role-forms-select');
        const fgwSelect = row.querySelector('.role-fgw-select');

        // 3. Сохраняем оригинальные значения с сервера
        if (!originalData.has(performerId)) {
            // 3.1. Берем значение из атрибутов select (они содержать оригинальные значения)
            const originalFormValue = formSelect.getAttribute('data-original') || formSelect.value;
            const originalFgwValue = fgwSelect.getAttribute('data-original') || fgwSelect.value;

            let originalFormText = '';
            let originalFgwText = ''

            for (let option of formSelect.options) {
                if (option.value === originalFormValue) {
                    originalFormText = option.text.trim()
                    break;
                }
            }

            for (let option of fgwSelect.options) {
                if (option.value === originalFgwValue) {
                    originalFgwText = option.text.trim()
                    break;
                }
            }

            originalData.set(performerId, {
                formsValue: originalFormValue,
                fgwValue: originalFgwValue,
                formText: originalFormText,
                fgwText: originalFgwText
            });
        }

        // 4. Получаем сохраненные оригинальные значения
        const original = originalData.get(performerId);

        // 5. Устанавливаем текущее значение в select
        formSelect.value = original.formsValue;
        fgwSelect.value = original.fgwValue;

        // 6. Сохраняем текущие значения для возможности отмены
        row.dataset.originalFormsValue = original.formsValue;
        row.dataset.originalFgwValue = original.fgwValue;
        row.dataset.performerId = performerId.toString();

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
    }

    function disablePerformersEditMode(row) {
        // 1. Получаем сохраненные значения для восстановления
        const originalFormsValue = row.dataset.originalFormsValue;
        const originalFgwValue = row.dataset.originalFgwValue;

        const formSelect = row.querySelector('.role-forms-select');
        const fgwSelect = row.querySelector('.role-fgw-select');

        // 2. Восстанавливаем значение в select
        if (originalFormsValue && formSelect) {
            formSelect.value = originalFormsValue;
        }

        if (originalFgwValue && fgwSelect) {
            fgwSelect.value = originalFgwValue;
        }

        // 3. Скрываем поля редактирования
        row.querySelectorAll('.edit-mode').forEach(el => {
            el.style.display = 'none';
        });

        // 4. Показываем поля просмотра
        row.querySelectorAll('.view-mode').forEach(el => {
            el.style.display = 'table-cell';
        });

        // 5. Показываем кнопку редактирования
        row.querySelector('.edit-btn').style.display = 'block';
        row.querySelector('.edit-buttons').style.display = 'none';

        // 6. Убираем визуальные индикаторы
        row.classList.remove('editing');
        row.style.backgroundColor = '';

        // Очищаем временные data-атрибуты
        delete row.dataset.originalFormsValue;
        delete row.dataset.originalFgwValue;
        delete row.dataset.performerId;
    }

    async function saveChanges(row) {
        // 1. Получаем Id из data-id атрибута строки
        const performerIdStr = row.getAttribute('data-id'); // {{ .Obj }}
        const performerId = parseInt(performerIdStr, 10);

        // 2. Получаем элементы select
        const formSelect = row.querySelector('.role-forms-select');
        const fgwSelect = row.querySelector('.role-fgw-select');

        // 3. Получаем текстовые значения выбранных опций
        const selectedFormsText = formSelect.options[formSelect.selectedIndex].text;
        const selectedFgwText = fgwSelect.options[fgwSelect.selectedIndex].text;

        // 4. Преобразуем значения в числа
        const idRoleAForms = parseInt(formSelect.value, 10);
        const idRoleAFGW = parseInt(fgwSelect.value, 10);

        // 5. Показываем индикатор загрузки
        const saveBtn = row.querySelector('.save-btn');
        saveBtn.innerHTML = '<span class="spinner-border spinner-border-sm" role="status"></span>';
        saveBtn.disabled = true;

        try {
            // 6. Отправляем запрос через Fetch API
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
            handleSuccessUpdate(row, result, selectedFormsText, selectedFgwText, performerId, idRoleAForms, idRoleAFGW);
        } catch (error) {
            console.error('Save error:', error);

            // Показываем уведомление об ошибке
            showPerformersNotification(`Ошибка: ${error.message}`, 'danger');

            // Восстанавливаем кнопку
            saveBtn.innerHTML = '<span>✓</span>';
            saveBtn.disabled = false;

            throw error; // Пробрасываем ошибку дальше
        }

    }





    function handleSuccessUpdate(row, result, roleId, name, description) {
        // 1. Обновление оригинальных данных
        originalData.set(roleId, {

            formsValue: idRoleAForms.toString(),
            fgwValue: idRoleAFGW.toString(),
            formsText: selectedFormsText,
            fgwText: selectedFgwText
        });

        // 2. Обновляем отображение ролей
        row.querySelector('.forms-role .badge').textContent = selectedFormsText;
        row.querySelector('.fgw-role .badge').textContent = selectedFgwText;

        // 3. Обновляем значения в select'ах
        const formSelect = row.querySelector('.role-forms-select');
        const fgwSelect = row.querySelector('.role-fgw-select');
        const updateAt = row.querySelector('.update-at');
        const updateBy = row.querySelector('.update-by');

        if (formSelect) {
            formSelect.value = idRoleAForms.toString();
            // 3.1. Обновляем атрибут data-original
            formSelect.setAttribute('data-original', idRoleAForms.toString());
        }

        if (fgwSelect) {
            fgwSelect.value = idRoleAFGW.toString();
            fgwSelect.setAttribute('data-original', idRoleAFGW.toString());
        }

        updateAt.textContent = result.updatedAt;
        updateBy.textContent = result.updatedBy;

        // 4. Выходим из режима редактирования
        disablePerformersEditMode(row)

        // 5. Показываем уведомление
        showPerformersNotification(result.message || 'Изменения успешно сохранены', 'success');

        // 6. Восстанавливаем кнопку сохранения
        const saveBtn = row.querySelector('.save-btn');
        saveBtn.innerHTML = '<span>✓</span>';
        saveBtn.disabled = false;
    }

    function showPerformersNotification(message, type) {
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

document.addEventListener('DOMContentLoaded', function() {
    const searchInput = document.getElementById('searchInput');
    const searchForm = document.getElementById('searchForm');
    let debounceTimer;

    // 1. Авто-поиск при вводе
    searchInput.addEventListener('input', function(e) {
        clearTimeout(debounceTimer);

        // 2. Если поле очищено - сразу отправляем
        if (e.target.value === '') {
            searchForm.submit();
            return;
        }

        // 3. Ждем 800ms после последнего ввода
        debounceTimer = setTimeout(() => {
            searchForm.submit();
        }, 800);
    });

    // 4. Фокус на поле поиска если есть поисковый запрос
    if (searchInput.value) {
        searchInput.focus();
        // 5. Помещаем курсор в конец
        searchInput.setSelectionRange(searchInput.value.length, searchInput.value.length);
    }
});