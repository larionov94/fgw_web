/**
 * Role Management Module
 * @module RoleManager
 * @description Управление ролями пользователей с поддержкой CRUD операций
 */

// Конфигурация
const CONFIG = {
    API: {
        BASE_URL: '/admin/roles',
        ENDPOINTS: {
            ADD: '/add',
            UPDATE: '/upd',
            DELETE: '/del'
        }
    },
    SELECTORS: {
        EDIT_BTN: '.edit-btn',
        CANCEL_BTN: '.cancel-btn',
        SAVE_BTN: '.save-role-btn',
        ADD_BTN: '.add-role-btn',
        DEL_BTN: '.del-btn',
        ROLE_ROW: 'tr[data-id]',
        ADD_MODAL: '#addRoleModal'
    },
    CLASSES: {
        EDITING: 'editing',
        VIEW_MODE: 'view-mode',
        EDIT_MODE: 'edit-mode',
        EDIT_BUTTONS: '.edit-buttons'
    },
    STORAGE_KEYS: {
        ORIGINAL_DATA: 'roleOriginalData'
    },
    MESSAGES: {
        DELETE_CONFIRM: 'Вы уверены, что хотите удалить эту роль?',
        DELETE_SUCCESS: 'Роль успешно удалена',
        DELETE_ERROR: 'Ошибка при удалении роли'
    }
};

/**
 * Класс для управления состояниями ролей
 */
class RoleStateManager {
    constructor() {
        this.originalData = new Map();
    }

    saveOriginal(roleId, data) {
        this.originalData.set(roleId, {...data});
    }

    getOriginal(roleId) {
        return this.originalData.get(roleId);
    }

    updateOriginal(roleId, data) {
        this.originalData.set(roleId, {...data});
    }

    hasOriginal(roleId) {
        return this.originalData.has(roleId);
    }

    clearTemporary(row) {
        ['originalName', 'originalDesc', 'roleId'].forEach(key => {
            delete row.dataset[key];
        });
    }
}

/**
 * Класс для управления UI состоянием строки роли
 */
class RoleRowManager {
    static enableEditMode(row, originalData) {
        // Сохраняем оригинальные значения
        const roleId = this.getRoleId(row);

        // Устанавливаем значения
        const inputs = this.getInputs(row);
        inputs.name.value = originalData.name;
        inputs.description.value = originalData.description;

        // Сохраняем для возможности отмены
        row.dataset.originalName = originalData.name;
        row.dataset.originalDesc = originalData.description;
        row.dataset.roleId = roleId.toString();

        // Переключаем UI состояния
        this.toggleEditModeUI(row, true);

        // Фокус
        inputs.name.focus();
    }

    static disableEditMode(row) {
        const inputs = this.getInputs(row);
        const originalName = row.dataset.originalName;
        const originalDesc = row.dataset.originalDesc;

        // Восстанавливаем значения
        if (originalName && inputs.name) {
            inputs.name.value = originalName;
        }
        if (originalDesc && inputs.description) {
            inputs.description.value = originalDesc;
        }

        // Переключаем UI состояния
        this.toggleEditModeUI(row, false);
    }

    static toggleEditModeUI(row, isEditMode) {
        // Управляем видимостью элементов
        const viewModeElements = row.querySelectorAll(CONFIG.CLASSES.VIEW_MODE);
        const editModeElements = row.querySelectorAll(CONFIG.CLASSES.EDIT_MODE);
        const editBtn = row.querySelector(CONFIG.SELECTORS.EDIT_BTN);
        const editButtons = row.querySelector(CONFIG.CLASSES.EDIT_BUTTONS);

        // Переключаем классы и стили
        row.classList.toggle(CONFIG.CLASSES.EDITING, isEditMode);
        row.style.backgroundColor = isEditMode ? '#f8f9fa' : '';

        // Переключаем видимость
        viewModeElements.forEach(el => {
            el.style.display = isEditMode ? 'none' : 'table-cell';
        });

        editModeElements.forEach(el => {
            el.style.display = isEditMode ? 'table-cell' : 'none';
        });

        editBtn.style.display = isEditMode ? 'none' : 'block';
        if (editButtons) {
            editButtons.style.display = isEditMode ? 'flex' : 'none';
        }
    }

    static getRoleId(row) {
        return parseInt(row.getAttribute('data-id'), 10);
    }

    static getInputs(row) {
        return {
            name: row.querySelector('input[name="name"]'),
            description: row.querySelector('input[name="description"]')
        };
    }

    static updateRowData(row, data) {
        // Обновляем текстовые представления
        const nameElement = row.querySelector('.forms-name');
        const descElement = row.querySelector('.forms-desc');
        const updateAtElement = row.querySelector('.update-at');
        const updateByElement = row.querySelector('.update-by');

        if (nameElement) nameElement.textContent = data.name;
        if (descElement) descElement.textContent = data.description;
        if (updateAtElement) updateAtElement.textContent = data.updatedAt || '';
        if (updateByElement) updateByElement.textContent = data.updatedBy || '';

        // Обновляем инпуты
        const inputs = this.getInputs(row);
        if (inputs.name) {
            inputs.name.value = data.name;
            inputs.name.setAttribute('data-original', data.name);
        }
        if (inputs.description) {
            inputs.description.value = data.description;
            inputs.description.setAttribute('data-original', data.description);
        }
    }

    static resetSaveButton(button) {
        if (!button) return;
        button.innerHTML = '<span>✓</span>';
        button.disabled = false;
    }

    static setLoadingState(button, isLoading) {
        if (!button) return;

        if (isLoading) {
            button.innerHTML = '<span class="spinner-border spinner-border-sm" role="status"></span>';
            button.disabled = true;
        } else {
            this.resetSaveButton(button);
        }
    }
}

/**
 * Класс для работы с API
 */
class RoleAPI {
    static async addRole(data) {
        return this._makeRequest(CONFIG.API.ENDPOINTS.ADD, {
            RoleId: data.roleId,
            Name: data.name,
            Description: data.description
        });
    }

    static async updateRole(data) {
        return this._makeRequest(CONFIG.API.ENDPOINTS.UPDATE, {
            roleId: data.roleId,
            name: data.name,
            description: data.description
        });
    }

    static async delRole(data) {
        try {
            const response = await fetch(`${CONFIG.API.BASE_URL}${CONFIG.API.ENDPOINTS.DELETE}`, {
                method: 'DELETE',
                headers: {
                    'Content-Type': 'application/json',
                    'Accept': 'application/json'
                },
                body: JSON.stringify({
                    roleId: data.roleId
                })
            });

            if (!response.ok) {
                await this._handleError(response);
            }

            return await response.json();
        } catch (error) {
            console.error('Delete API error:', error);
            throw error;
        }
    }

    static async _makeRequest(endpoint, data) {
        const response = await fetch(`${CONFIG.API.BASE_URL}${endpoint}`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Accept': 'application/json'
            },
            body: JSON.stringify(data)
        });

        if (!response.ok) {
            await this._handleError(response);
        }

        return await response.json();
    }

    static async _makeRequestDel(endpoint, data) {
        const response = await fetch(`${CONFIG.API.BASE_URL}${endpoint}`, {
            method: 'DELETE',
            headers: {
                'Content-Type': 'application/json',
                'Accept': 'application/json'
            },
            body: JSON.stringify(data)
        });

        if (!response.ok) {
            await this._handleError(response);
        }

        return await response.json();
    }

    static async _handleError(response) {
        const contentType = response.headers.get('content-type');

        if (contentType && contentType.includes('application/json')) {
            const result = await response.json();
            throw new Error(result.error || `HTTP ${response.status}`);
        }

        throw new Error(`HTTP ${response.status}`);
    }
}

/**
 * Класс для управления уведомлениями
 */
class NotificationManager {
    static show(message, type = 'info') {
        // Удаляем существующие уведомления
        this.clear();

        // Создаем новое уведомление
        const notification = this._createNotificationElement(message, type);

        // Добавляем на страницу
        document.body.appendChild(notification);

        // Автоматическое скрытие
        this._setupAutoDismiss(notification);
    }

    static _createNotificationElement(message, type) {
        const element = document.createElement('div');
        element.className = `alert alert-${type} alert-dismissible fade show position-fixed`;
        element.style.cssText = `
            top: 20px;
            right: 20px;
            z-index: 9999;
            min-width: 300px;
            max-width: 500px;
            box-shadow: 0 2px 10px rgba(0,0,0,0.1);
        `;

        element.innerHTML = `
            <div class="d-flex align-items-center">
                <div class="flex-grow-1">${this._escapeHtml(message)}</div>
                <button type="button" class="btn-close ms-2" data-bs-dismiss="alert" aria-label="Close"></button>
            </div>
        `;

        return element;
    }

    static _setupAutoDismiss(element) {
        setTimeout(() => {
            if (element.parentNode) {
                element.remove();
            }
        }, 5000);
    }

    static clear() {
        document.querySelectorAll('.alert.position-fixed').forEach(el => el.remove());
    }

    static _escapeHtml(text) {
        const div = document.createElement('div');
        div.textContent = text;
        return div.innerHTML;
    }
}

/**
 * Класс для валидации
 */
class RoleValidator {
    static validateAddForm(roleId, name, description) {
        const errors = [];

        if (!roleId || roleId <= 0) {
            errors.push('Некорректный ID роли');
        }

        if (!name || name.trim() === '') {
            errors.push('Название роли не может быть пустым');
        }

        if (!description || description.trim() === '') {
            errors.push('Описание роли не может быть пустым');
        }

        return {
            isValid: errors.length === 0,
            errors
        };
    }

    static validateEditForm(name, description) {
        const errors = [];

        if (!name || name.trim() === '') {
            errors.push('Название роли не может быть пустым');
        }

        if (!description || description.trim() === '') {
            errors.push('Описание роли не может быть пустым');
        }

        return {
            isValid: errors.length === 0,
            errors
        };
    }
}

/**
 * Главный класс управления ролями
 */
class RoleManager {
    constructor() {
        this.stateManager = new RoleStateManager();
        this.bindEvents();
    }

    bindEvents() {
        // Используем делегирование событий для лучшей производительности
        document.addEventListener('click', this.handleClick.bind(this));

        // Обработка модального окна
        const addModal = document.querySelector(CONFIG.SELECTORS.ADD_MODAL);
        if (addModal) {
            addModal.addEventListener('hidden.bs.modal', this.clearAddForm.bind(this));
        }
    }

    handleClick(event) {
        // Редактирование
        if (event.target.closest(CONFIG.SELECTORS.EDIT_BTN)) {
            const btn = event.target.closest(CONFIG.SELECTORS.EDIT_BTN);
            const row = btn.closest(CONFIG.SELECTORS.ROLE_ROW);
            this.handleEditClick(row);
        }

        // Отмена редактирования
        else if (event.target.closest(CONFIG.SELECTORS.CANCEL_BTN)) {
            const btn = event.target.closest(CONFIG.SELECTORS.CANCEL_BTN);
            const row = btn.closest(CONFIG.SELECTORS.ROLE_ROW);
            this.handleCancelClick(row);
        }

        // Сохранение
        else if (event.target.closest(CONFIG.SELECTORS.SAVE_BTN)) {
            const btn = event.target.closest(CONFIG.SELECTORS.SAVE_BTN);
            const row = btn.closest(CONFIG.SELECTORS.ROLE_ROW);
            this.handleSaveClick(row);
        }

        // Добавление
        else if (event.target.closest(CONFIG.SELECTORS.ADD_BTN)) {
            const btn = event.target.closest(CONFIG.SELECTORS.ADD_BTN);
            this.handleAddClick(btn);
        }

        // Удаление
        else if (event.target.closest(CONFIG.SELECTORS.DEL_BTN)) {
            const btn = event.target.closest(CONFIG.SELECTORS.DEL_BTN);
            const row = btn.closest(CONFIG.SELECTORS.ROLE_ROW);
            this.handleDeleteClick(row);
        }
    }

    handleEditClick(row) {
        const roleId = RoleRowManager.getRoleId(row);

        // Если оригинальные данные не сохранены - сохраняем их
        if (!this.stateManager.hasOriginal(roleId)) {
            const inputs = RoleRowManager.getInputs(row);
            const originalData = {
                name: inputs.name.getAttribute('data-original') || inputs.name.value,
                description: inputs.description.getAttribute('data-original') || inputs.description.value
            };
            this.stateManager.saveOriginal(roleId, originalData);
        }

        const originalData = this.stateManager.getOriginal(roleId);
        RoleRowManager.enableEditMode(row, originalData);
    }

    async handleDeleteClick(row) {
        const roleId = RoleRowManager.getRoleId(row);
        const roleName = row.querySelector('.forms-name')?.textContent || '';

        // Подтверждение удаления
        if (!confirm(`${CONFIG.MESSAGES.DELETE_CONFIRM}\nРоль: ${roleName} (ID: ${roleId})`)) {
            return;
        }

        const deleteBtn = row.querySelector(CONFIG.SELECTORS.DEL_BTN);
        const originalContent = deleteBtn.innerHTML;

        try {
            // Показываем индикатор загрузки
            deleteBtn.innerHTML = '<span class="spinner-border spinner-border-sm" role="status"></span>';
            deleteBtn.disabled = true;

            // Отправляем запрос на удаление
            const result = await RoleAPI.delRole({ roleId });

            if (result.success) {
                // Удаляем строку из таблицы
                row.remove();

                // Очищаем сохраненные данные
                this.stateManager.originalData.delete(roleId);

                // Показываем уведомление
                NotificationManager.show(CONFIG.MESSAGES.DELETE_SUCCESS, 'success');
                window.location.reload();
            } else {
                NotificationManager.show(result.message || CONFIG.MESSAGES.DELETE_ERROR, 'danger');
                this.resetDeleteButton(deleteBtn, originalContent);
            }
        } catch (error) {
            console.error('Delete error:', error);
            NotificationManager.show(`${CONFIG.MESSAGES.DELETE_ERROR}: ${error.message}`, 'danger');
            this.resetDeleteButton(deleteBtn, originalContent);
        }
    }

    resetDeleteButton(button, originalContent) {
        button.innerHTML = originalContent;
        button.disabled = false;
    }

    handleCancelClick(row) {
        RoleRowManager.disableEditMode(row);
        this.stateManager.clearTemporary(row);
    }

    async handleSaveClick(row) {
        const saveBtn = row.querySelector(CONFIG.SELECTORS.SAVE_BTN);

        // Предотвращение двойных кликов
        if (saveBtn.disabled) return;

        RoleRowManager.setLoadingState(saveBtn, true);

        try {
            const roleId = RoleRowManager.getRoleId(row);
            const inputs = RoleRowManager.getInputs(row);
            const name = inputs.name.value.trim();
            const description = inputs.description.value.trim();

            // Валидация
            const validation = RoleValidator.validateEditForm(name, description);
            if (!validation.isValid) {
                validation.errors.forEach(error => {
                    NotificationManager.show(error, 'warning');
                });
                return;
            }

            // Отправка на сервер
            const result = await RoleAPI.updateRole({
                roleId,
                name,
                description
            });

            // Обработка успешного ответа
            this.handleUpdateSuccess(row, result, roleId, name, description);

        } catch (error) {
            console.error('Save error:', error);
            NotificationManager.show(`Ошибка: ${error.message}`, 'danger');
            RoleRowManager.resetSaveButton(saveBtn);
            throw error;
        }
    }

    async handleAddClick(button) {
        if (button.disabled) return;

        button.disabled = true;
        const originalText = button.innerHTML;
        button.innerHTML = '<span class="spinner-border spinner-border-sm" role="status"></span> Добавление...';

        try {
            // Получение данных из формы
            const roleId = parseInt(document.getElementById('newRoleId')?.value || 0, 10);
            const name = document.getElementById('newRoleName')?.value.trim() || '';
            const description = document.getElementById('newRoleDescription')?.value.trim() || '';

            // Валидация
            const validation = RoleValidator.validateAddForm(roleId, name, description);
            if (!validation.isValid) {
                validation.errors.forEach(error => {
                    NotificationManager.show(error, 'warning');
                });
                return;
            }

            // Отправка на сервер
            const result = await RoleAPI.addRole({
                roleId,
                name,
                description
            });

            // Обработка успешного ответа
            this.handleAddSuccess(result, button, originalText);

        } catch (error) {
            console.error('Add error:', error);
            NotificationManager.show(`Ошибка: ${error.message}`, 'danger');
            this.resetAddButton(button, originalText);
            throw error;
        }
    }

    handleUpdateSuccess(row, result, roleId, name, description) {
        // Обновляем оригинальные данные
        this.stateManager.updateOriginal(roleId, {name, description});

        // Обновляем UI
        RoleRowManager.updateRowData(row, {
            name,
            description,
            updatedAt: result.updatedAt,
            updatedBy: result.updatedBy
        });

        // Выходим из режима редактирования
        RoleRowManager.disableEditMode(row);
        this.stateManager.clearTemporary(row);

        // Восстанавливаем кнопку
        const saveBtn = row.querySelector(CONFIG.SELECTORS.SAVE_BTN);
        RoleRowManager.resetSaveButton(saveBtn);

        // Уведомление
        NotificationManager.show(result.message || 'Роль успешно обновлена', 'success');
    }

    handleAddSuccess(result, button, originalText) {
        if (result.success) {
            // Закрываем модальное окно
            const modal = bootstrap.Modal.getInstance(document.querySelector(CONFIG.SELECTORS.ADD_MODAL));
            if (modal) {
                modal.hide();
            }

            // Обновляем таблицу (если необходимо)
            this.refreshTable();

            NotificationManager.show(result.message || 'Роль успешно добавлена', 'success');
        } else {
            NotificationManager.show(result.message || 'Ошибка при добавлении роли', 'danger');
        }

        this.resetAddButton(button, originalText);
    }

    resetAddButton(button, originalText) {
        button.innerHTML = originalText;
        button.disabled = false;
    }

    clearAddForm() {
        const form = document.querySelector('#addRoleForm');
        if (form) {
            form.reset();
        }
    }

    refreshTable() {
        window.location.reload()
        // Здесь может быть логика обновления таблицы
        // Например, через window.location.reload() или обновление через AJAX
        console.log('Table refresh logic here');
    }
}

// Инициализация при загрузке страницы
document.addEventListener('DOMContentLoaded', () => {
    try {
        window.roleManager = new RoleManager();
    } catch (error) {
        console.error('Failed to initialize RoleManager:', error);
        NotificationManager.show('Ошибка инициализации системы ролей', 'danger');
    }
});