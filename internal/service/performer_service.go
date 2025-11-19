package service

import (
	"FGW_WEB/internal/model"
	"FGW_WEB/internal/repository"
	"FGW_WEB/pkg/common"
	"FGW_WEB/pkg/common/msg"
	"context"
	"errors"
	"fmt"
)

type PerformerService struct {
	performerRepo repository.PerformerRepository
	logg          *common.Logger
}

func NewPerformerService(performerRepo repository.PerformerRepository, logger *common.Logger) *PerformerService {
	return &PerformerService{performerRepo: performerRepo, logg: logger}
}

type PerformerUseCase interface {
	GetAllPerformers(ctx context.Context) ([]*model.Performer, error)
	AuthPerformer(ctx context.Context, id int, password string) (*model.AuthPerformer, error)
	UpdPerformer(ctx context.Context, id int, performer *model.Performer) error
	ExistPerformer(ctx context.Context, id int) (bool, error)
}

func (p *PerformerService) GetAllPerformers(ctx context.Context) ([]*model.Performer, error) {
	performers, err := p.performerRepo.All(ctx)
	if err != nil {
		p.logg.LogE(msg.E3209, err)

		return nil, fmt.Errorf("%s: %v", msg.E3209, err)
	}

	return performers, nil
}

func (p *PerformerService) AuthPerformer(ctx context.Context, id int, password string) (*model.AuthPerformer, error) {
	if id <= 0 || password == "" {
		p.logg.LogE(msg.E3211, nil)

		return &model.AuthPerformer{Success: false, Message: msg.E3211}, errors.New("ТН или пароль не должны быть пустыми")
	}

	authOK, err := p.performerRepo.AuthByIdAndPass(ctx, id, password)
	if err != nil {
		p.logg.LogE(msg.E3210, err)

		return &model.AuthPerformer{Success: false, Message: msg.E3210 + " AuthPerformer.AuthByIdAndPass()"}, err
	}

	if !authOK {
		p.logg.LogE(msg.E3210, err)

		return &model.AuthPerformer{Success: false, Message: msg.E3210 + " AuthPerformer.AuthByIdAndPass()"}, err
	}

	performer, err := p.performerRepo.FindById(ctx, id)
	if err != nil {
		return &model.AuthPerformer{Success: false, Message: msg.E3212 + "AuthPerformer.FindById"}, err
	}

	return &model.AuthPerformer{
		Success:   true,
		Performer: *performer,
		Message:   "Успешный вход",
	}, nil
}

func (p *PerformerService) UpdPerformer(ctx context.Context, id int, performer *model.Performer) error {
	if err := model.ValidateUpdateDataPerformer(performer); err != nil {
		p.logg.LogE(msg.E3213, err)

		return fmt.Errorf("%s: %v", msg.E3213, err)
	}

	if err := p.performerRepo.UpdById(ctx, id, performer); err != nil {
		p.logg.LogE(msg.E3216, err)

		return fmt.Errorf("%s: %v", msg.E3216, err)
	}

	return nil
}

func (p *PerformerService) ExistPerformer(ctx context.Context, id int) (bool, error) {
	return p.performerRepo.ExistById(ctx, id)
}
