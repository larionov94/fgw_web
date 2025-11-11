package model

// Role роль.
type Role struct {
	Id       int
	Name     string
	Desc     string
	AuditRec Audit
}
