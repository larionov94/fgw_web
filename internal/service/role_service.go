package service

import (
	"FGW_WEB/internal/model"
	"FGW_WEB/internal/repository"
	"FGW_WEB/pkg/common"
	"FGW_WEB/pkg/common/msg"
	"context"
)

type RoleService struct {
	roleRepo repository.RoleRepository
	logg     *common.Logger
}

func NewRoleService(roleRepo repository.RoleRepository, logger *common.Logger) *RoleService {
	return &RoleService{roleRepo: roleRepo, logg: logger}
}

type RoleUseCase interface {
	GetAllRole(ctx context.Context) ([]*model.Role, error)
	UpdRole(ctx context.Context, id int, role *model.Role) error
	AddRole(ctx context.Context, role *model.Role) error
	ExistRole(ctx context.Context, id int) (bool, error)
	FindRoleById(ctx context.Context, id int) (*model.Role, error)
	DelRoleById(ctx context.Context, id int) error
}

func (r *RoleService) GetAllRole(ctx context.Context) ([]*model.Role, error) {
	roles, err := r.roleRepo.All(ctx)
	if err != nil {
		r.logg.LogE(msg.E3209, err)

		return nil, err
	}

	return roles, nil
}

func (r *RoleService) UpdRole(ctx context.Context, id int, role *model.Role) error {
	if err := model.ValidateUpdateDataRole(role); err != nil {
		r.logg.LogE(msg.E3213, err)

		return err
	}

	if err := r.roleRepo.UpdById(ctx, id, role); err != nil {
		r.logg.LogE(msg.E3216, err)

		return err
	}

	return nil
}

func (r *RoleService) AddRole(ctx context.Context, role *model.Role) error {
	if err := model.ValidateUpdateDataRole(role); err != nil {
		r.logg.LogE(msg.E3213, err)

		return err
	}

	if err := r.roleRepo.Add(ctx, role); err != nil {
		r.logg.LogE(msg.E3215, err)

		return err
	}

	return nil
}

func (r *RoleService) FindRoleById(ctx context.Context, id int) (*model.Role, error) {
	role, err := r.roleRepo.FindById(ctx, id)
	if err != nil {
		r.logg.LogE(msg.E3212, err)

		return nil, err
	}

	return role, nil
}

func (r *RoleService) ExistRole(ctx context.Context, id int) (bool, error) {
	return r.roleRepo.ExistById(ctx, id)
}

func (r *RoleService) DelRoleById(ctx context.Context, id int) error {
	if err := r.roleRepo.DelById(ctx, id); err != nil {
		r.logg.LogE(msg.E3216, err)

		return err
	}

	return nil
}
