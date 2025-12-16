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
        ADD_MODAL: '#addRoleModal',
        ROLES_TABLE: '#rolesTable',
        ROLES_TABLE_BODY: '#rolesTable tbody',
        ROLES_COUNT: '.roles-count'
    },
    CLASSES: {
        EDITING: 'editing',
        VIEW_MODE: 'view-mode',
        EDIT_MODE: 'edit-mode',
        EDIT_BUTTONS: '.edit-buttons'
    },
    MESSAGES: {
        DELETE_CONFIRM: 'Вы уверены, что хотите удалить эту роль?',
        DELETE_SUCCESS: 'Роль успешно удалена',
        DELETE_ERROR: 'Ошибка при удалении роли'
    }
};

/**
 * Класс для управления таблицей ролей
 */
class RolesTableManager {
    /**
     * Удаляет строку из таблицы с анимацией
     * @param {HTMLElement} row - Строка для удаления
     * @returns {Promise<void>}
     */
    static async removeRowWithAnimation(row) {
        return new Promise(resolve => {
            // Анимация удаления
            row.style.transition = 'all 0.3s ease';
            row.style.transform = 'translateX(-100%)';
            row.style.opacity = '0';

            // Ждем завершения анимации
            setTimeout(() => {
                row.remove();
                this.updateRolesCount();
                resolve();
            }, 300);
        });
    }

    /**
     * Добавляет новую строку в таблицу
     * @param {Object} roleData - Данные роли
     */
    static addNewRow(roleData) {
        const tbody = document.querySelector(CONFIG.SELECTORS.ROLES_TABLE_BODY);
        if (!tbody) return;

        // Создаем HTML для новой строки
        const newRowHtml = this.createRowHtml(roleData);

        // Создаем временный элемент для вставки
        const tempDiv = document.createElement('div');
        tempDiv.innerHTML = newRowHtml;
        const newRow = tempDiv.firstElementChild;

        // Анимация появления
        newRow.style.opacity = '0';
        newRow.style.transform = 'translateY(-20px)';
        tbody.prepend(newRow);

        // Анимация
        requestAnimationFrame(() => {
            newRow.style.transition = 'all 0.3s ease';
            newRow.style.opacity = '1';
            newRow.style.transform = 'translateY(0)';
        });

        this.updateRolesCount();
    }

    /**
     * Создает HTML для строки таблицы
     * @param {Object} roleData - Данные роли
     * @returns {string} HTML строка
     */
    static createRowHtml(roleData) {
        return `
            <tr id="role-${roleData.id}" data-id="${roleData.id}">
                <!-- ИД (всегда в режиме просмотра) -->
                <td class="fw-semibold">${roleData.id}</td>

                <!-- Наименование - режим просмотра -->
                <td class="view-mode forms-name">${this._escapeHtml(roleData.name)}</td>

                <!-- Описание - режим просмотра -->
                <td class="view-mode forms-desc">${this._escapeHtml(roleData.description)}</td>

                <!-- Наименование - режим редактирования (скрыт) -->
                <td class="edit-mode" style="display: none;">
                    <label style="width: 75%">
                        <input type="text"
                               name="name"
                               value="${this._escapeHtml(roleData.name)}"
                               class="form-control form-control-sm"
                               data-original="${this._escapeHtml(roleData.name)}"
                               required>
                    </label>
                </td>

                <!-- Описание - режим редактирования (скрыт) -->
                <td class="edit-mode" style="display: none;">
                    <label style="width: 95%">
                        <input type="text"
                               name="description"
                               value="${this._escapeHtml(roleData.description)}"
                               class="form-control form-control-sm"
                               data-original="${this._escapeHtml(roleData.description)}"
                               required>
                    </label>
                </td>

                <!-- Дата создания -->
                <td>${roleData.createdAt || ''}</td>

                <!-- ТН создателя -->
                <td>${roleData.createdBy || ''}</td>

                <!-- Дата изменения -->
                <td class="update-at">${roleData.updatedAt || ''}</td>

                <!-- ТН редактора -->
                <td class="update-by">${roleData.updatedBy || ''}</td>

                <!-- Кнопки операций -->
                <td>
                    <div class="d-flex justify-content-center gap-2">
                        <!-- Кнопка редактирования (отображается в view-mode) -->
                        <button class="btn btn-sm btn-outline-primary edit-btn"
                                title="Редактировать">
                            <span>✏️</span>
                        </button>

                        <button class="btn btn-sm btn-outline-primary del-btn"
                                title="Удалить">
                            <span>🗑️</span>
                        </button>

                        <!-- Кнопки сохранения/отмены (скрыты в view-mode) -->
                        <div class="edit-buttons" style="display: none;">
                            <button class="btn btn-sm btn-success save-role-btn" title="Сохранить">
                                <span>✓</span>
                            </button>
                            <button class="btn btn-sm btn-secondary cancel-btn" title="Отмена">
                                <span>✗</span>
                            </button>
                        </div>
                    </div>
                </td>
            </tr>
        `;
    }

    /**
     * Обновляет счетчик ролей
     */
    static updateRolesCount() {
        const rows = document.querySelectorAll(CONFIG.SELECTORS.ROLE_ROW);
        const countElement = document.querySelector(CONFIG.SELECTORS.ROLES_COUNT);

        if (countElement) {
            countElement.textContent = `Всего ролей: ${rows.length}`;
        } else {
            // Ищем элемент по тексту, если нет специального класса
            const elements = document.querySelectorAll('p');
            elements.forEach(el => {
                if (el.textContent.includes('Всего ролей')) {
                    el.textContent = `Всего ролей: ${rows.length}`;
                }
            });
        }
    }

    /**
     * Обновляет строку с данными
     * @param {HTMLElement} row - Строка таблицы
     * @param {Object} data - Новые данные
     */
    static updateRow(row, data) {
        RoleRowManager.updateRowData(row, data);
        this.updateRolesCount();
    }

    /**
     * Экранирует HTML-сущности
     * @param {string} text - Текст для экранирования
     * @returns {string} Экранированный текст
     */
    static _escapeHtml(text) {
        const div = document.createElement('div');
        div.textContent = text;
        return div.innerHTML;
    }
}

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

    removeOriginal(roleId) {
        this.originalData.delete(roleId);
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
        const roleId = this.getRoleId(row);
        const inputs = this.getInputs(row);
        inputs.name.value = originalData.name;
        inputs.description.value = originalData.description;

        row.dataset.originalName = originalData.name;
        row.dataset.originalDesc = originalData.description;
        row.dataset.roleId = roleId.toString();

        this.toggleEditModeUI(row, true);
        inputs.name.focus();
    }

    static disableEditMode(row) {
        const inputs = this.getInputs(row);
        const originalName = row.dataset.originalName;
        const originalDesc = row.dataset.originalDesc;

        if (originalName && inputs.name) {
            inputs.name.value = originalName;
        }
        if (originalDesc && inputs.description) {
            inputs.description.value = originalDesc;
        }

        this.toggleEditModeUI(row, false);
    }

    static toggleEditModeUI(row, isEditMode) {
        const viewModeElements = row.querySelectorAll(CONFIG.CLASSES.VIEW_MODE);
        const editModeElements = row.querySelectorAll(CONFIG.CLASSES.EDIT_MODE);
        const editBtn = row.querySelector(CONFIG.SELECTORS.EDIT_BTN);
        const editButtons = row.querySelector(CONFIG.CLASSES.EDIT_BUTTONS);

        row.classList.toggle(CONFIG.CLASSES.EDITING, isEditMode);
        row.style.backgroundColor = isEditMode ? '#f8f9fa' : '';

        viewModeElements.forEach(el => {
            el.style.display = isEditMode ? 'none' : 'table-cell';
        });

        editModeElements.forEach(el => {
            el.style.display = isEditMode ? 'table-cell' : 'none';
        });

        if (editBtn) editBtn.style.display = isEditMode ? 'none' : 'block';
        if (editButtons) editButtons.style.display = isEditMode ? 'flex' : 'none';
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
        const nameElement = row.querySelector('.forms-name');
        const descElement = row.querySelector('.forms-desc');
        const updateAtElement = row.querySelector('.update-at');
        const updateByElement = row.querySelector('.update-by');

        if (nameElement) nameElement.textContent = data.name;
        if (descElement) descElement.textContent = data.description;
        if (updateAtElement) updateAtElement.textContent = data.updatedAt || '';
        if (updateByElement) updateByElement.textContent = data.updatedBy || '';

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
        return this._makeRequest(CONFIG.API.ENDPOINTS.DELETE, {
            roleId: data.roleId
        }, 'DELETE');
    }

    static async _makeRequest(endpoint, data, method = 'POST') {
        const response = await fetch(`${CONFIG.API.BASE_URL}${endpoint}`, {
            method: method,
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
        this.clear();

        const notification = this._createNotificationElement(message, type);
        document.body.appendChild(notification);

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
        document.addEventListener('click', this.handleClick.bind(this));

        const addModal = document.querySelector(CONFIG.SELECTORS.ADD_MODAL);
        if (addModal) {
            addModal.addEventListener('hidden.bs.modal', this.clearAddForm.bind(this));
        }
    }

    handleClick(event) {
        if (event.target.closest(CONFIG.SELECTORS.EDIT_BTN)) {
            const btn = event.target.closest(CONFIG.SELECTORS.EDIT_BTN);
            const row = btn.closest(CONFIG.SELECTORS.ROLE_ROW);
            this.handleEditClick(row);
        }
        else if (event.target.closest(CONFIG.SELECTORS.CANCEL_BTN)) {
            const btn = event.target.closest(CONFIG.SELECTORS.CANCEL_BTN);
            const row = btn.closest(CONFIG.SELECTORS.ROLE_ROW);
            this.handleCancelClick(row);
        }
        else if (event.target.closest(CONFIG.SELECTORS.SAVE_BTN)) {
            const btn = event.target.closest(CONFIG.SELECTORS.SAVE_BTN);
            const row = btn.closest(CONFIG.SELECTORS.ROLE_ROW);
            this.handleSaveClick(row);
        }
        else if (event.target.closest(CONFIG.SELECTORS.ADD_BTN)) {
            const btn = event.target.closest(CONFIG.SELECTORS.ADD_BTN);
            this.handleAddClick(btn);
        }
        else if (event.target.closest(CONFIG.SELECTORS.DEL_BTN)) {
            const btn = event.target.closest(CONFIG.SELECTORS.DEL_BTN);
            const row = btn.closest(CONFIG.SELECTORS.ROLE_ROW);
            this.handleDeleteClick(row);
        }
    }

    handleEditClick(row) {
        const roleId = RoleRowManager.getRoleId(row);

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

        if (!confirm(`${CONFIG.MESSAGES.DELETE_CONFIRM}\nРоль: ${roleName} (ID: ${roleId})`)) {
            return;
        }

        const deleteBtn = row.querySelector(CONFIG.SELECTORS.DEL_BTN);
        const originalContent = deleteBtn.innerHTML;

        try {
            deleteBtn.innerHTML = '<span class="spinner-border spinner-border-sm" role="status"></span>';
            deleteBtn.disabled = true;

            const result = await RoleAPI.delRole({ roleId });

            if (result.success) {
                // Удаляем строку с анимацией вместо перезагрузки
                await RolesTableManager.removeRowWithAnimation(row);

                // Очищаем сохраненные данные
                this.stateManager.removeOriginal(roleId);

                NotificationManager.show(CONFIG.MESSAGES.DELETE_SUCCESS, 'success');
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

        if (saveBtn.disabled) return;

        RoleRowManager.setLoadingState(saveBtn, true);

        try {
            const roleId = RoleRowManager.getRoleId(row);
            const inputs = RoleRowManager.getInputs(row);
            const name = inputs.name.value.trim();
            const description = inputs.description.value.trim();

            const validation = RoleValidator.validateEditForm(name, description);
            if (!validation.isValid) {
                validation.errors.forEach(error => {
                    NotificationManager.show(error, 'warning');
                });
                return;
            }

            const result = await RoleAPI.updateRole({
                roleId,
                name,
                description
            });

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
            const roleId = parseInt(document.getElementById('newRoleId')?.value || 0, 10);
            const name = document.getElementById('newRoleName')?.value.trim() || '';
            const description = document.getElementById('newRoleDescription')?.value.trim() || '';

            const validation = RoleValidator.validateAddForm(roleId, name, description);
            if (!validation.isValid) {
                validation.errors.forEach(error => {
                    NotificationManager.show(error, 'warning');
                });
                return;
            }

            const result = await RoleAPI.addRole({
                roleId,
                name,
                description
            });

            this.handleAddSuccess(result, button, originalText);

        } catch (error) {
            console.error('Add error:', error);
            NotificationManager.show(`Ошибка: ${error.message}`, 'danger');
            this.resetAddButton(button, originalText);
            throw error;
        }
    }

    handleUpdateSuccess(row, result, roleId, name, description) {
        this.stateManager.updateOriginal(roleId, {name, description});

        // Обновляем данные в таблице
        RolesTableManager.updateRow(row, {
            name,
            description,
            updatedAt: result.updatedAt,
            updatedBy: result.updatedBy
        });

        RoleRowManager.disableEditMode(row);
        this.stateManager.clearTemporary(row);

        const saveBtn = row.querySelector(CONFIG.SELECTORS.SAVE_BTN);
        RoleRowManager.resetSaveButton(saveBtn);

        NotificationManager.show(result.message || 'Роль успешно обновлена', 'success');
    }

    handleAddSuccess(result, button, originalText) {
        if (result.success) {
            const modal = bootstrap.Modal.getInstance(document.querySelector(CONFIG.SELECTORS.ADD_MODAL));
            if (modal) {
                modal.hide();
            }

            // Добавляем новую строку в таблицу динамически
            if (result.role) {
                RolesTableManager.addNewRow(result.role);
            } else {
                // Если сервер не вернул данные роли, перезагружаем таблицу через AJAX
                this.refreshTable();
            }

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

    async refreshTable() {
        try {
            // AJAX запрос для обновления таблицы
            const response = await fetch('/admin/roles', {
                headers: {
                    'X-Requested-With': 'XMLHttpRequest'
                }
            });

            if (response.ok) {
                const html = await response.text();
                // Парсим HTML и обновляем только таблицу
                const parser = new DOMParser();
                const doc = parser.parseFromString(html, 'text/html');
                const newTable = doc.querySelector(CONFIG.SELECTORS.ROLES_TABLE);

                if (newTable) {
                    const currentTable = document.querySelector(CONFIG.SELECTORS.ROLES_TABLE);
                    currentTable.parentNode.replaceChild(newTable, currentTable);

                    // Обновляем счетчик
                    RolesTableManager.updateRolesCount();
                }
            }
        } catch (error) {
            console.error('Refresh table error:', error);
            // В крайнем случае - обычная перезагрузка
            window.location.reload();
        }
    }
}

// Инициализация
document.addEventListener('DOMContentLoaded', () => {
    try {
        window.roleManager = new RoleManager();
    } catch (error) {
        console.error('Failed to initialize RoleManager:', error);
        NotificationManager.show('Ошибка инициализации системы ролей', 'danger');
    }
});