package model

import "fmt"

type RoleList struct {
	Roles []*Role `json:"roles"`
}

// Role роль.
type Role struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	Desc      string `json:"desc"`
	AuditRec  Audit  `json:"auditRec"`
	IsEditing bool   `json:"isEditing"` // Флаг для редактирования поля.
}

type RoleUpdate struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func ValidateUpdateDataRole(data *Role) error {
	if data == nil {
		return fmt.Errorf("ошибка: не удалось обновить данные, данных нет")
	}

	if data.Name == "" || data.Desc == "" {
		return fmt.Errorf("ошибка: не валидное поле")
	}

	return nil
}
