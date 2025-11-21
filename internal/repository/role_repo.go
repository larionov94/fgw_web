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

type RoleRepo struct {
	mssql *sql.DB
	logg  *common.Logger
}

func NewRoleRepo(mssql *sql.DB, logger *common.Logger) *RoleRepo {
	return &RoleRepo{mssql: mssql, logg: logger}
}

type RoleRepository interface {
	All(ctx context.Context) ([]*model.Role, error)
	Add(ctx context.Context, role *model.Role) error
	UpdById(ctx context.Context, id int, role *model.Role) error
	FindById(ctx context.Context, id int) (*model.Role, error)
	ExistById(ctx context.Context, id int) (bool, error)
}

func (r *RoleRepo) All(ctx context.Context) ([]*model.Role, error) {
	rows, err := r.mssql.QueryContext(ctx, FGWsvRoleAllQuery)
	if err != nil {
		r.logg.LogE(msg.E3202, err)

		return nil, err
	}
	defer db.RowsClose(rows)

	var roles []*model.Role
	for rows.Next() {
		var role model.Role
		if err = rows.Scan(
			&role.Id,
			&role.Name,
			&role.Desc,
			&role.AuditRec.CreatedAt,
			&role.AuditRec.CreatedBy,
			&role.AuditRec.UpdatedAt,
			&role.AuditRec.UpdatedBy,
		); err != nil {
			r.logg.LogE(msg.E3204, err)

			return nil, err
		}

		roles = append(roles, &role)
	}

	if err = rows.Err(); err != nil {
		r.logg.LogE(msg.E3205, err)

		return nil, err
	}

	return roles, nil
}

func (r *RoleRepo) Add(ctx context.Context, role *model.Role) error {
	if _, err := r.mssql.ExecContext(ctx, FGWsvRoleAddQuery,
		role.Id,
		role.Name,
		role.Desc,
		role.AuditRec.CreatedBy,
	); err != nil {
		r.logg.LogE(msg.E3215, err)

		return err
	}

	return nil
}

func (r *RoleRepo) UpdById(ctx context.Context, id int, role *model.Role) error {
	_, err := r.mssql.ExecContext(ctx, FGWsvRoleUpdByIdQuery, id, role.Name, role.Desc, role.AuditRec.UpdatedBy)
	if err != nil {
		r.logg.LogE(msg.E3216, err)

		return err
	}

	return nil
}

func (r *RoleRepo) FindById(ctx context.Context, id int) (*model.Role, error) {
	var role model.Role

	if err := r.mssql.QueryRowContext(ctx, FGWsvRoleFindByIdQuery, id).Scan(
		&role.Id,
		&role.Name,
		&role.Desc,
		&role.AuditRec.CreatedAt,
		&role.AuditRec.CreatedBy,
		&role.AuditRec.UpdatedAt,
		&role.AuditRec.UpdatedBy,
	); err != nil {
		r.logg.LogE(msg.E3204, err)

		if errors.Is(err, sql.ErrNoRows) {
			r.logg.LogE(msg.E3206, err)

			return nil, err
		}
		return nil, fmt.Errorf("%s: %v", msg.E3202, err)
	}

	return &role, nil
}

func (r *RoleRepo) ExistById(ctx context.Context, id int) (bool, error) {
	var exists bool

	err := r.mssql.QueryRowContext(ctx, FGWsvRoleExistsByIdQuery, id).Scan(&exists)
	if err != nil {
		r.logg.LogE(msg.E3206, err)

		return false, err
	}

	return exists, nil
}
