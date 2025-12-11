package repository

import (
	"FGW_WEB/internal/config/db"
	"FGW_WEB/internal/model"
	"FGW_WEB/pkg/common"
	"FGW_WEB/pkg/common/msg"
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type PerformerRepo struct {
	mssql *sql.DB
	logg  *common.Logger
}

func NewPerformerRepo(mssql *sql.DB, logger *common.Logger) *PerformerRepo {
	return &PerformerRepo{mssql: mssql, logg: logger}
}

type PerformerRepository interface {
	All(ctx context.Context) ([]*model.Performer, error)
	AuthByIdAndPass(ctx context.Context, id int, password string) (bool, error)
	FindById(ctx context.Context, id int) (*model.Performer, error)
	UpdById(ctx context.Context, id int, performer *model.Performer) error
	ExistById(ctx context.Context, id int) (bool, error)
	GetPerformersCount(ctx context.Context) (int, error)
	GetPerformersWithPagination(ctx context.Context, offset, limit int) ([]*model.Performer, error)
	FilterById(ctx context.Context, pattern string) ([]*model.Performer, error)
}

// All получить всех сотрудников из БД.
func (p *PerformerRepo) All(ctx context.Context) ([]*model.Performer, error) {
	rows, err := p.mssql.QueryContext(ctx, FGWsvPerformerAllQuery)
	if err != nil {
		p.logg.LogE(msg.E3202, err)

		return nil, err
	}
	defer db.RowsClose(rows)

	var performers []*model.Performer
	for rows.Next() {
		var performer model.Performer
		if err = rows.Scan(
			&performer.Id,
			&performer.FIO,
			&performer.BC,
			&performer.Pass,
			&performer.Archive,
			&performer.IdRoleAForms,
			&performer.IdRoleAFGW,
			&performer.AuditRec.CreatedAt,
			&performer.AuditRec.CreatedBy,
			&performer.AuditRec.UpdatedAt,
			&performer.AuditRec.UpdatedBy,
		); err != nil {
			p.logg.LogE(msg.E3204, err)

			return nil, err
		}

		performers = append(performers, &performer)
	}

	if err = rows.Err(); err != nil {
		p.logg.LogE(msg.E3205, err)

		return nil, err
	}

	return performers, nil
}

// AuthByIdAndPass проверка существования в БД сотрудника.
func (p *PerformerRepo) AuthByIdAndPass(ctx context.Context, id int, password string) (bool, error) {
	var authSuccess bool

	err := p.mssql.QueryRowContext(ctx, FGWsvPerformerAuthQuery, id, password).Scan(&authSuccess)
	if err != nil {
		p.logg.LogE(msg.E3202, err)

		return false, err
	}

	return authSuccess, nil
}

// FindById ищет сотрудника по ИД.
func (p *PerformerRepo) FindById(ctx context.Context, id int) (*model.Performer, error) {
	var performer model.Performer

	if err := p.mssql.QueryRowContext(ctx, FGWsvPerformerFindByIdQuery, id).Scan(
		&performer.Id,
		&performer.FIO,
		&performer.BC,
		&performer.Pass,
		&performer.Archive,
		&performer.IdRoleAForms,
		&performer.IdRoleAFGW,
		&performer.AuditRec.CreatedAt,
		&performer.AuditRec.CreatedBy,
		&performer.AuditRec.UpdatedAt,
		&performer.AuditRec.UpdatedBy,
	); err != nil {
		p.logg.LogE(msg.E3204, err)

		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: %v", msg.E3206, err)
		}
		return nil, err
	}

	return &performer, nil
}

// UpdById обновить данные сотрудника по табельному номеру в БД.
func (p *PerformerRepo) UpdById(ctx context.Context, id int, performer *model.Performer) error {
	_, err := p.mssql.ExecContext(ctx, FGWsvPerformerUpdByIdQuery, id, performer.IdRoleAForms,
		performer.IdRoleAFGW, performer.AuditRec.UpdatedBy)
	if err != nil {
		p.logg.LogE(msg.E3202, err)

		return err
	}

	return nil
}

// ExistById проверяет существование сотрудника.
func (p *PerformerRepo) ExistById(ctx context.Context, id int) (bool, error) {
	var exists bool

	err := p.mssql.QueryRowContext(ctx, FGWsvPerformerExistsByIdQuery, id).Scan(&exists)
	if err != nil {
		p.logg.LogE(msg.E3206, err)

		return false, err
	}

	return exists, nil
}

// GetPerformersCount кол-во сотрудников.
func (p *PerformerRepo) GetPerformersCount(ctx context.Context) (int, error) {
	var count int
	if err := p.mssql.QueryRowContext(ctx, FGWsvPerformersCountQuery).Scan(&count); err != nil {
		p.logg.LogE(msg.E3217, err)

		return 0, err
	}

	return count, nil
}

// GetPerformersWithPagination получает сотрудников с нумерации страниц.
func (p *PerformerRepo) GetPerformersWithPagination(ctx context.Context, offset, limit int) ([]*model.Performer, error) {
	startRow := offset
	endRow := offset + limit

	rows, err := p.mssql.QueryContext(ctx, FGWsvPerformersPaginationQuery, startRow, endRow)
	if err != nil {
		return nil, err
	}
	defer db.RowsClose(rows)

	var performers []*model.Performer
	for rows.Next() {
		var performer model.Performer

		if err = rows.Scan(
			&performer.Id,
			&performer.FIO,
			&performer.BC,
			&performer.Pass,
			&performer.Archive,
			&performer.IdRoleAForms,
			&performer.IdRoleAFGW,
			&performer.AuditRec.CreatedAt,
			&performer.AuditRec.CreatedBy,
			&performer.AuditRec.UpdatedAt,
			&performer.AuditRec.UpdatedBy,
		); err != nil {
			p.logg.LogE(msg.E3204, err)

			return nil, err
		}

		performers = append(performers, &performer)
	}

	if err = rows.Err(); err != nil {
		p.logg.LogE(msg.E3205, err)

		return nil, err
	}

	return performers, nil
}

func (p *PerformerRepo) FilterById(ctx context.Context, pattern string) ([]*model.Performer, error) {
	rows, err := p.mssql.QueryContext(ctx, FGWsvPerformerFilterByIdQuery, pattern)
	if err != nil {
		p.logg.LogE(msg.E3206, err)

		return nil, err
	}
	defer db.RowsClose(rows)

	var performers []*model.Performer

	for rows.Next() {
		var performer model.Performer

		if err = rows.Scan(
			&performer.Id,
			&performer.FIO,
			&performer.BC,
			&performer.Pass,
			&performer.Archive,
			&performer.IdRoleAForms,
			&performer.IdRoleAFGW,
			&performer.AuditRec.CreatedAt,
			&performer.AuditRec.CreatedBy,
			&performer.AuditRec.UpdatedAt,
			&performer.AuditRec.UpdatedBy,
		); err != nil {
			p.logg.LogE(msg.E3204, err)

			return nil, err
		}

		performers = append(performers, &performer)
	}
	if err = rows.Err(); err != nil {
		p.logg.LogE(msg.E3205, err)

		return nil, err
	}

	return performers, nil
}
