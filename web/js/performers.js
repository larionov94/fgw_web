/**
 * Performers Management Module
 * @module PerformersManager
 * @description Управление исполнителями с поддержкой CRUD операций и поиска
 */

// Конфигурация модуля
const PERFORMERS_CONFIG = {
    API: {
        BASE_URL: '/admin/performers',
        ENDPOINTS: {
            UPDATE: '/upd'
        }
    },
    SELECTORS: {
        EDIT_BTN: '.edit-btn',
        CANCEL_BTN: '.cancel-btn',
        SAVE_BTN: '.save-btn',
        SEARCH_INPUT: '#searchInput',
        SEARCH_FORM: '#searchForm',
        PERFORMER_ROW: 'tr[data-id]',
        FORM_SELECT: 'select.role-forms-select, .role-forms-select',
        FGW_SELECT: 'select.role-fgw-select, .role-fgw-select',
        EDIT_BUTTONS: '.edit-buttons'
    },
    CLASSES: {
        EDITING: 'editing',
        VIEW_MODE: '.view-mode',
        EDIT_MODE: '.edit-mode'
    },
    UI: {
        DEBOUNCE_DELAY: 800,
        NOTIFICATION_TIMEOUT: 5000,
        ROW_BG_EDITING: '#f8f9fa'
    },
    MESSAGES: {
        SAVE_SUCCESS: 'Изменения успешно сохранены',
        SAVE_ERROR: 'Ошибка при сохранении',
        SEARCH_ERROR: 'Ошибка при поиске'
    }
};

/**
 * Класс для управления состояниями исполнителей
 */
class PerformersStateManager {
    constructor() {
        this.originalData = new Map();
    }

    saveOriginal(performerId, data) {
        this.originalData.set(performerId, {...data});
    }

    getOriginal(performerId) {
        return this.originalData.get(performerId);
    }

    updateOriginal(performerId, data) {
        this.originalData.set(performerId, {...data});
    }

    hasOriginal(performerId) {
        return this.originalData.has(performerId);
    }

    clearTemporaryData(row) {
        ['originalFormsValue', 'originalFgwValue', 'performerId'].forEach(key => {
            delete row.dataset[key];
        });
    }
}

/**
 * Класс для управления UI состоянием строки исполнителя
 */
class PerformerRowManager {
    static enableEditMode(row, originalData) {
        const performerId = this.getPerformerId(row);
        const selects = this.getSelects(row);

        if (!selects || !selects.form || !selects.fgw) {
            console.error('Select elements not found in row:', row);
            return false;
        }

        // Устанавливаем значения
        selects.form.value = originalData.formsValue;
        selects.fgw.value = originalData.fgwValue;

        // Сохраняем для возможности отмены
        row.dataset.originalFormsValue = originalData.formsValue;
        row.dataset.originalFgwValue = originalData.fgwValue;
        row.dataset.performerId = performerId.toString();

        // Переключаем UI
        this.toggleEditModeUI(row, true);

        // Фокус на первом select
        selects.form.focus();
        return true;
    }

    static disableEditMode(row) {
        const selects = this.getSelects(row);
        const originalFormsValue = row.dataset.originalFormsValue;
        const originalFgwValue = row.dataset.originalFgwValue;

        // Восстанавливаем значения
        if (originalFormsValue && selects && selects.form) {
            selects.form.value = originalFormsValue;
        }
        if (originalFgwValue && selects && selects.fgw) {
            selects.fgw.value = originalFgwValue;
        }

        // Переключаем UI
        this.toggleEditModeUI(row, false);
    }

    static toggleEditModeUI(row, isEditMode) {
        const viewModeElements = row.querySelectorAll(PERFORMERS_CONFIG.CLASSES.VIEW_MODE);
        const editModeElements = row.querySelectorAll(PERFORMERS_CONFIG.CLASSES.EDIT_MODE);
        const editBtn = row.querySelector(PERFORMERS_CONFIG.SELECTORS.EDIT_BTN);
        const editButtons = row.querySelector(PERFORMERS_CONFIG.SELECTORS.EDIT_BUTTONS);

        row.classList.toggle('editing', isEditMode);
        row.style.backgroundColor = isEditMode ? PERFORMERS_CONFIG.UI.ROW_BG_EDITING : '';

        // Переключаем видимость элементов
        viewModeElements.forEach(el => {
            el.style.display = isEditMode ? 'none' : '';
        });

        editModeElements.forEach(el => {
            el.style.display = isEditMode ? '' : 'none';
        });

        if (editBtn) editBtn.style.display = isEditMode ? 'none' : '';
        if (editButtons) editButtons.style.display = isEditMode ? 'flex' : 'none';
    }

    static getPerformerId(row) {
        const id = row.getAttribute('data-id');
        if (!id) {
            console.error('data-id attribute not found in row:', row);
            return null;
        }
        return parseInt(id, 10);
    }

    static getSelects(row) {
        if (!row) {
            console.error('Row is null in getSelects');
            return null;
        }

        const formSelect = row.querySelector(PERFORMERS_CONFIG.SELECTORS.FORM_SELECT);
        const fgwSelect = row.querySelector(PERFORMERS_CONFIG.SELECTORS.FGW_SELECT);

        // Проверяем, что элементы существуют
        if (!formSelect || !fgwSelect) {
            console.warn('Select elements not found in row. Form:', formSelect, 'FGW:', fgwSelect);
            console.warn('Row HTML:', row.outerHTML);
            return null;
        }

        return {
            form: formSelect,
            fgw: fgwSelect
        };
    }

    static getSelectedOptionText(select) {
        if (!select) return '';
        return select.options[select.selectedIndex]?.text?.trim() || '';
    }

    static updateRowData(row, data) {
        // Обновляем отображение ролей
        const formsBadge = row.querySelector('.forms-role .badge');
        const fgwBadge = row.querySelector('.fgw-role .badge');
        const updateAt = row.querySelector('.update-at');
        const updateBy = row.querySelector('.update-by');

        if (formsBadge) formsBadge.textContent = data.formsText;
        if (fgwBadge) fgwBadge.textContent = data.fgwText;
        if (updateAt) updateAt.textContent = data.updatedAt || '';
        if (updateBy) updateBy.textContent = data.updatedBy || '';

        // Обновляем select элементы
        const selects = this.getSelects(row);
        if (selects && selects.form) {
            selects.form.value = data.formsValue;
            selects.form.setAttribute('data-original', data.formsValue);
        }
        if (selects && selects.fgw) {
            selects.fgw.value = data.fgwValue;
            selects.fgw.setAttribute('data-original', data.fgwValue);
        }
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

    static resetSaveButton(button) {
        if (!button) return;
        button.innerHTML = '<span>✓</span>';
        button.disabled = false;
    }
}

/**
 * Класс для работы с API исполнителей
 */
class PerformersAPI {
    static async updatePerformer(data) {
        return this._makeRequest(PERFORMERS_CONFIG.API.ENDPOINTS.UPDATE, {
            performerId: data.performerId,
            idRoleAForms: data.idRoleAForms,
            idRoleAFGW: data.idRoleAFGW
        });
    }

    static async _makeRequest(endpoint, data) {
        const response = await fetch(`${PERFORMERS_CONFIG.API.BASE_URL}${endpoint}`, {
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
class PerformersNotificationManager {
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
        }, PERFORMERS_CONFIG.UI.NOTIFICATION_TIMEOUT);
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
 * Класс для управления поиском с debounce
 */
class SearchManager {
    constructor() {
        this.debounceTimer = null;
        this.searchInput = document.querySelector(PERFORMERS_CONFIG.SELECTORS.SEARCH_INPUT);
        this.searchForm = document.querySelector(PERFORMERS_CONFIG.SELECTORS.SEARCH_FORM);

        this.init();
    }

    init() {
        if (!this.searchInput || !this.searchForm) {
            console.warn('Search elements not found');
            return;
        }

        this.searchInput.addEventListener('input', this.handleInput.bind(this));
        this.focusSearchInput();
    }

    handleInput(event) {
        clearTimeout(this.debounceTimer);

        if (event.target.value === '') {
            this.searchForm.submit();
            return;
        }

        this.debounceTimer = setTimeout(() => {
            this.searchForm.submit();
        }, PERFORMERS_CONFIG.UI.DEBOUNCE_DELAY);
    }

    focusSearchInput() {
        if (this.searchInput && this.searchInput.value) {
            this.searchInput.focus();
            const length = this.searchInput.value.length;
            this.searchInput.setSelectionRange(length, length);
        }
    }

    cleanup() {
        clearTimeout(this.debounceTimer);
    }
}

/**
 * Класс для валидации данных исполнителей
 */
class PerformersValidator {
    static validateEditForm(idRoleAForms, idRoleAFGW) {
        const errors = [];

        if (!idRoleAForms || idRoleAForms <= 0) {
            errors.push('Необходимо выбрать роль для форм');
        }

        if (!idRoleAFGW || idRoleAFGW <= 0) {
            errors.push('Необходимо выбрать роль для ФГВ');
        }

        return {
            isValid: errors.length === 0,
            errors
        };
    }
}

/**
 * Главный класс управления исполнителями
 */
class PerformersManager {
    constructor() {
        this.stateManager = new PerformersStateManager();
        this.searchManager = null;
        this.bindEvents();
    }

    bindEvents() {
        document.addEventListener('click', this.handleClick.bind(this));
    }

    handleClick(event) {
        // Редактирование
        if (event.target.closest(PERFORMERS_CONFIG.SELECTORS.EDIT_BTN)) {
            const btn = event.target.closest(PERFORMERS_CONFIG.SELECTORS.EDIT_BTN);
            const row = btn.closest(PERFORMERS_CONFIG.SELECTORS.PERFORMER_ROW);
            if (row) {
                this.handleEditClick(row);
            } else {
                console.warn('Row not found for edit button');
            }
        }

        // Отмена редактирования
        else if (event.target.closest(PERFORMERS_CONFIG.SELECTORS.CANCEL_BTN)) {
            const btn = event.target.closest(PERFORMERS_CONFIG.SELECTORS.CANCEL_BTN);
            const row = btn.closest(PERFORMERS_CONFIG.SELECTORS.PERFORMER_ROW);
            if (row) {
                this.handleCancelClick(row);
            }
        }

        // Сохранение
        else if (event.target.closest(PERFORMERS_CONFIG.SELECTORS.SAVE_BTN)) {
            const btn = event.target.closest(PERFORMERS_CONFIG.SELECTORS.SAVE_BTN);
            const row = btn.closest(PERFORMERS_CONFIG.SELECTORS.PERFORMER_ROW);
            if (row) {
                this.handleSaveClick(row);
            }
        }
    }

    handleEditClick(row) {
        if (!row) {
            console.error('Row is null in handleEditClick');
            return;
        }

        const performerId = PerformerRowManager.getPerformerId(row);
        if (!performerId) {
            console.error('Could not get performer ID from row:', row);
            return;
        }

        // Если оригинальные данные не сохранены - сохраняем их
        if (!this.stateManager.hasOriginal(performerId)) {
            const selects = PerformerRowManager.getSelects(row);
            if (!selects) {
                console.error('Could not get selects for performer:', performerId);
                PerformersNotificationManager.show('Ошибка: элементы формы не найдены', 'danger');
                return;
            }

            const originalData = {
                formsValue: selects.form.getAttribute('data-original') || selects.form.value,
                fgwValue: selects.fgw.getAttribute('data-original') || selects.fgw.value,
                formsText: PerformerRowManager.getSelectedOptionText(selects.form),
                fgwText: PerformerRowManager.getSelectedOptionText(selects.fgw)
            };
            this.stateManager.saveOriginal(performerId, originalData);
        }

        const originalData = this.stateManager.getOriginal(performerId);
        if (originalData) {
            PerformerRowManager.enableEditMode(row, originalData);
        }
    }

    handleCancelClick(row) {
        if (!row) return;
        PerformerRowManager.disableEditMode(row);
        this.stateManager.clearTemporaryData(row);
    }

    async handleSaveClick(row) {
        if (!row) return;

        const saveBtn = row.querySelector(PERFORMERS_CONFIG.SELECTORS.SAVE_BTN);
        if (!saveBtn) return;

        // Предотвращение двойных кликов
        if (saveBtn.disabled) return;

        PerformerRowManager.setLoadingState(saveBtn, true);

        try {
            const performerId = PerformerRowManager.getPerformerId(row);
            if (!performerId) {
                throw new Error('Не удалось получить ID исполнителя');
            }

            const selects = PerformerRowManager.getSelects(row);
            if (!selects) {
                throw new Error('Элементы формы не найдены');
            }

            // Получаем значения
            const idRoleAForms = parseInt(selects.form.value, 10);
            const idRoleAFGW = parseInt(selects.fgw.value, 10);
            const formsText = PerformerRowManager.getSelectedOptionText(selects.form);
            const fgwText = PerformerRowManager.getSelectedOptionText(selects.fgw);

            // Валидация
            const validation = PerformersValidator.validateEditForm(idRoleAForms, idRoleAFGW);
            if (!validation.isValid) {
                validation.errors.forEach(error => {
                    PerformersNotificationManager.show(error, 'warning');
                });
                PerformerRowManager.resetSaveButton(saveBtn);
                return;
            }

            // Отправка на сервер
            const result = await PerformersAPI.updatePerformer({
                performerId,
                idRoleAForms,
                idRoleAFGW
            });

            // Обработка успешного ответа
            this.handleUpdateSuccess(row, result, performerId, {
                formsValue: idRoleAForms.toString(),
                fgwValue: idRoleAFGW.toString(),
                formsText,
                fgwText,
                updatedAt: result.updatedAt,
                updatedBy: result.updatedBy
            });

        } catch (error) {
            console.error('Save error:', error);
            PerformersNotificationManager.show(`${PERFORMERS_CONFIG.MESSAGES.SAVE_ERROR}: ${error.message}`, 'danger');
            PerformerRowManager.resetSaveButton(saveBtn);
        }
    }

    handleUpdateSuccess(row, result, performerId, data) {
        // Обновляем оригинальные данные
        this.stateManager.updateOriginal(performerId, {
            formsValue: data.formsValue,
            fgwValue: data.fgwValue,
            formsText: data.formsText,
            fgwText: data.fgwText
        });

        // Обновляем UI
        PerformerRowManager.updateRowData(row, data);

        // Выходим из режима редактирования
        PerformerRowManager.disableEditMode(row);
        this.stateManager.clearTemporaryData(row);

        // Восстанавливаем кнопку
        const saveBtn = row.querySelector(PERFORMERS_CONFIG.SELECTORS.SAVE_BTN);
        PerformerRowManager.resetSaveButton(saveBtn);

        // Уведомление
        PerformersNotificationManager.show(
            result.message || PERFORMERS_CONFIG.MESSAGES.SAVE_SUCCESS,
            'success'
        );
    }

    initSearch() {
        this.searchManager = new SearchManager();
    }

    cleanup() {
        if (this.searchManager) {
            this.searchManager.cleanup();
        }
    }
}

// Инициализация при загрузке страницы
document.addEventListener('DOMContentLoaded', () => {
    try {
        window.performersManager = new PerformersManager();
        window.performersManager.initSearch();

        window.addEventListener('beforeunload', () => {
            if (window.performersManager) {
                window.performersManager.cleanup();
            }
        });
    } catch (error) {
        console.error('Failed to initialize PerformersManager:', error);
        PerformersNotificationManager.show('Ошибка инициализации системы исполнителей', 'danger');
    }
});