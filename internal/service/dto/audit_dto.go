package dto

// AuditDTO аудит для отслеживания изменений данных.
type AuditDTO struct {
	CreatedAt string // CreatedAt - дата создания записи.
	CreatedBy int    // CreatedBy - табельный номер сотрудника.
	UpdatedAt string // UpdatedAt - дата изменения записи.
	UpdatedBy int    // UpdatedBy - табельный номер сотрудника изменивший запись.
}
