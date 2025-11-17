package dto

type PerformerDTOList struct {
	Performers []PerformerDTO
}

type PerformerDTO struct {
	Id           int      `json:"id"`
	FIO          string   `json:"fio"`
	BC           string   `json:"bc"`
	Archive      bool     `json:"archive"`
	IdRoleAForms int      `json:"idRoleAForms"`
	IdRoleAFGW   int      `json:"idRoleAFGW"`
	Audit        AuditDTO `json:"audit"`
}

type AuthPerformerDTO struct {
	Success   bool         `json:"success"`
	Performer PerformerDTO `json:"performer"`
	Message   string       `json:"message"`
}
