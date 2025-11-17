package model

// Audit аудит для отслеживания изменений данных.
type Audit struct {
	CreatedAt string `json:"createdAt"` // CreatedAt - дата создания записи.
	CreatedBy int    `json:"createdBy"` // CreatedBy - табельный номер сотрудника.
	UpdatedAt string `json:"updatedAt"` // UpdatedAt - дата изменения записи.
	UpdatedBy int    `json:"updatedBy"` // UpdatedBy - табельный номер сотрудника изменивший запись.
}
