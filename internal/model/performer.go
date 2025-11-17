package model

type PerformerList struct {
	Performers []Performer
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
