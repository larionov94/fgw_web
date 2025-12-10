package model

import "fmt"

type PerformerList struct {
	Performers []*Performer
	Roles      []*Role
}

type Performer struct {
	Id           int    `json:"id"`           // Id - табельный номер.
	FIO          string `json:"fio"`          // FIO - ФИО сотрудника.
	BC           string `json:"bc"`           // BC - код доступа сотрудника.
	Pass         string `json:"password"`     // Pass - пароль сотрудника.
	Archive      bool   `json:"archive"`      // Archive - флаг архивного сотрудника.
	IdRoleAForms int    `json:"idRoleAForms"` // IdRoleAForms - id роли.
	IdRoleAFGW   int    `json:"idRoleAFGW"`   // IdRoleAFGW - id роли.
	AuditRec     Audit  `json:"auditRec"`     // AuditRec - аудит для отслеживания изменений данных.
}

type AuthPerformer struct {
	Success   bool      `json:"success"`
	Performer Performer `json:"performer"`
	Message   string    `json:"message"`
}

type PerformerUpdate struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func ValidateUpdateDataPerformer(data *Performer) error {
	if data == nil {
		return fmt.Errorf("ошибка: не удалось обновить данные, данных нет")
	}

	if data.AuditRec.UpdatedBy == 0 || data.AuditRec.UpdatedBy < 0 {
		return fmt.Errorf("ошибка: невалидное поле")
	}

	if data.IdRoleAForms < 0 {
		return fmt.Errorf("ошибка: невалидное поле")
	}

	if data.IdRoleAFGW < 0 {
		return fmt.Errorf("ошибка: невалидное поле")
	}

	return nil
}
