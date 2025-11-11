package repository

import (
	"FGW_WEB/internal/config/db"
	"FGW_WEB/internal/model"
	"FGW_WEB/pkg/common"
	"FGW_WEB/pkg/common/msg"
	"FGW_WEB/pkg/convert"
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
	All(ctx context.Context) ([]model.Performer, error)
	AuthByIdAndPass(ctx context.Context, id int, password string) (bool, error)
	FindById(ctx context.Context, id int) (*model.Performer, error)
	UpdById(ctx context.Context, id int, performer *model.Performer) error
	ExistById(ctx context.Context, id int) (bool, error)
}

// All получить всех сотрудников из БД.
func (p *PerformerRepo) All(ctx context.Context) ([]model.Performer, error) {
	rows, err := p.mssql.QueryContext(ctx, FGWsvPerformerAllQuery)
	if err != nil {
		p.logg.LogE(msg.E3202, err)

		return nil, fmt.Errorf("%s: %v", msg.E3202, err)
	}
	defer db.RowsClose(rows)

	var performers []model.Performer
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
			&performer.AuditRec.CreatedBy,
			&performer.AuditRec.CreatedAt,
			&performer.AuditRec.UpdatedBy,
			&performer.AuditRec.UpdatedAt,
		); err != nil {
			p.logg.LogE(msg.E3204, err)

			return nil, fmt.Errorf("%s: %v", msg.E3204, err)
		}
		performer.FIO, _ = convert.Win1251ToUTF8(performer.FIO)

		performers = append(performers, performer)
	}

	if err = rows.Err(); err != nil {
		p.logg.LogE(msg.E3205, err)

		return nil, fmt.Errorf("%s: %v", msg.E3205, err)
	}

	return performers, nil
}

// AuthByIdAndPass проверка существования в БД сотрудника.
func (p *PerformerRepo) AuthByIdAndPass(ctx context.Context, id int, password string) (bool, error) {
	var authSuccess bool

	err := p.mssql.QueryRowContext(ctx, FGWsvPerformerAuthQuery,
		sql.Named("id", id),
		sql.Named("pass", password)).Scan(&authSuccess)
	if err != nil {
		p.logg.LogE(msg.E3202, err)

		return false, fmt.Errorf("%s: %v", msg.E3202, err)
	}

	return authSuccess, nil
}

// FindById ищет сотрудника по ИД.
func (p *PerformerRepo) FindById(ctx context.Context, id int) (*model.Performer, error) {
	var performer model.Performer

	if err := p.mssql.QueryRowContext(ctx, FGWsvPerformerFindByIdQuery, sql.Named("id", id)).Scan(
		&performer.Id,
		&performer.FIO,
		&performer.BC,
		&performer.Pass,
		&performer.Archive,
		&performer.IdRoleAForms,
		&performer.IdRoleAFGW,
		&performer.AuditRec.CreatedBy,
		&performer.AuditRec.CreatedAt,
		&performer.AuditRec.UpdatedBy,
		&performer.AuditRec.UpdatedAt,
	); err != nil {
		p.logg.LogE(msg.E3202, err)

		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: %v", msg.E3206, err)
		}
		return nil, fmt.Errorf("%s: %v", msg.E3202, err)
	}

	performer.FIO, _ = convert.Win1251ToUTF8(performer.FIO)

	return &performer, nil
}

// UpdById обновить данные сотрудника по табельному номеру в БД.
func (p *PerformerRepo) UpdById(ctx context.Context, id int, performer *model.Performer) error {
	result, err := p.mssql.ExecContext(ctx, FGWsvPerformerUpdByIdQuery,
		sql.Named("id", id),
		sql.Named("id_role_a_forms", performer.IdRoleAForms),
		sql.Named("id_role_a_fgw", performer.IdRoleAFGW),
		sql.Named("updated_by", performer.AuditRec.UpdatedBy), // TODO: в поле updated_by - нужно подставить табельный номер авторизованного сотрудника.
	)
	if err != nil {
		p.logg.LogE(msg.E3202, err)

		return fmt.Errorf("%s: %v", msg.E3202, err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		p.logg.LogE(msg.E3207, err)

		return fmt.Errorf("%s: %v", msg.E3207, err)
	}

	if rows == 0 {
		p.logg.LogE(msg.E3208, err)

		return fmt.Errorf("%s: %v", msg.E3208, err)
	}

	return nil
}

// ExistById проверяет существование сотрудника.
func (p *PerformerRepo) ExistById(ctx context.Context, id int) (bool, error) {
	var exists bool

	err := p.mssql.QueryRowContext(ctx, FGWsvPerformerExistsByIdQuery, sql.Named("id", id)).Scan(&exists)
	if err != nil {
		p.logg.LogE(msg.E3206, err)

		return false, fmt.Errorf("%s: %v", msg.E3206, err)
	}

	return exists, nil
}
