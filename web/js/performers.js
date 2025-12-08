document.addEventListener('DOMContentLoaded', function () {
    // Обработчик кнопки редактирования.
    document.addEventListener('click', function (e) {
        if (e.target.closest('.edit-btn')) {
            const btn = e.target.closest('.edit-btn');
            const row = btn.closest('tr');
            const performerId = row.getAttribute('data-id');

            enableEditMode(row, performerId); // разрешить редактирование
        }

        // Обработчик кнопки сохранения
        if (e.target.closest('.save-btn')) {
            const btn = e.target.closest('.save-btn');
            const row = btn.closest('tr');
            const performerId = row.getAttribute('data-id');

            saveChanges(row, performerId); // сохранить изменения
        }

        if (e.target.closest('.cancel-btn')) {
            const btn = e.target.closest('.cancel-btn');
            const row = btn.closest('tr');

            disableEditMode(row); // отключаем режим редактирования
        }
    });

    function enableEditMode(row, performerId) {
        // 1. Сохраняем исходные значения
        const formsRole = row.querySelector('.forms-role .badge').textContent.trim();
        const fgwRole = row.querySelector('.fgw-role .badge').textContent.trim();

        // 2. Сохраняем эти значения в data-атрибуты строки
        row.setAttribute('data-original-forms', formsRole);
        row.setAttribute('data-original-fgw', fgwRole);

        // 3. Показываем поля редактирования
        row.querySelectorAll('.edit-mode').forEach(el => {
            el.style.display = 'table-cell';
        });

        // 4. Скрываем поля просмотра
        row.querySelectorAll('.view-mode').forEach(el => {
            el.style.display = 'none';
        });

        // 5. Показываем кнопки сохранения и отмены
        row.querySelector('.edit-btn').style.display = 'none';
        row.querySelector('.edit-buttons').style.display = 'flex';

        // 6. Добавляем класс активного редактирования
        row.classList.add('editing');
        row.style.background = '#f8f9fa';
    }

    function disableEditMode(row) {
        // 1. Восстанавливаем исходные значения
        const formSelect = row.querySelector('.role-forms-select');
        const fgwSelect = row.querySelector('.role-fgw-select');

        // 2. Находим option с id соответствующим исходному значению
        const originalForms = row.getAttribute('data-original-forms');
        const originalFgw = row.getAttribute('data-original-fgw');

        // 3. Восстанавливаем селекты к исходным значениям
        if (originalForms) {
            Array.from(formSelect.options).forEach(option => {
                if (option.text === originalForms) {
                    formSelect.value = option.value;
                }
            });
        }

        if (originalFgw) {
            Array.from(fgwSelect.options).forEach(option => {
                if (option.text === originalFgw) {
                    fgwSelect.value = option.value;
                }
            });
        }

        // 4. Скрываем поля редактирования
        row.querySelectorAll('.edit-mode').forEach(el => {
            el.style.display = 'none';
        });

        // 5. Показываем поля просмотра
        row.querySelectorAll('.view-mode').forEach(el => {
            el.style.display = 'table-cell';
        })

        // 6. Показываем кнопку редактирования
        row.querySelector('.edit-btn').style.display = 'block';
        row.querySelector('.edit-buttons').style.display = 'none';

        // 7. Убираем класс активного редактирования
        row.classList.remove('editing');
        row.style.backgroundColor = '';
    }

    // Сохранение изменений через стандартную форму
    function saveChanges(row, performerId) {
        // 1. Выбранные данные
        const formsSelect = row.querySelector('.role-forms-select');
        const fgwSelect = row.querySelector('.role-fgw-select');

        // 2. Заполняем скрытую форму
        document.getElementById('updatePerformerId').value = performerId;
        document.getElementById('updateRoleForms').value = formsSelect.value;
        document.getElementById('updateRoleFGW').value = fgwSelect.value;

        // 3. Создаем временный iframe для отправки формы
        const iframe = document.createElement('iframe');
        iframe.name = 'hiddenFrame';
        iframe.style.display = 'none';
        document.body.appendChild(iframe);

        // 4. Устанавливаем таргет формы на iframe
        const form = document.getElementById('updateForm');
        form.target = 'hiddenFrame';

        // 5. Обработчик загрузки iframe (когда форма отправилась)
        iframe.onload = function () {
            // 5.1. Обновляем отображаемые значения
            const selectedFormsText = formsSelect.options[formsSelect.selectedIndex].text;
            const selectedFgwText = fgwSelect.options[formsSelect.selectedIndex].text;

            row.querySelector('.forms-role .badge').textContent = selectedFormsText;
            row.querySelector('.fgw-role .badge').textContent = selectedFgwText;

            // 5.2. Выходим из режима редактирования.
            disableEditMode(row);

            // 5.3. Удаляем iframe
            setTimeout(() => {
                document.body.removeChild(iframe);
            }, 1000);
        };

        // 6. Отправляем форму
        form.submit();
    }
});

