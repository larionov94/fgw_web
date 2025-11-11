package model

// Audit аудит для отслеживания изменений данных.
type Audit struct {
	CreatedBy int    // CreatedBy - табельный номер сотрудника.
	CreatedAt string // CreatedAt - дата создания записи.
	UpdatedBy int    // UpdatedBy - табельный номер сотрудника изменивший запись.
	UpdatedAt string // UpdatedAt - дата изменения записи.
}
