package model

type Performer struct {
	Id           int    // Id - табельный номер.
	FIO          string // FIO - ФИО сотрудника.
	BC           string // BC - код доступа сотрудника.
	Pass         string // Pass - пароль сотрудника.
	Archive      bool   // Archive - флаг архивного сотрудника.
	IdRoleAForms int    // IdRoleAForms - id роли.
	IdRoleAFGW   int    // IdRoleAFGW - id роли.
	AuditRec     Audit  // AuditRec - аудит для отслеживания изменений данных.
}
