package service

import (
	"FGW_WEB/internal/model"
	"FGW_WEB/internal/repository"
	"FGW_WEB/internal/service/dto"
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
	GetAllPerformers(ctx context.Context) ([]dto.PerformerDTO, error)
	AuthPerformer(ctx context.Context, id int, password string) (*model.AuthPerformer, error)
}

func (p *PerformerService) GetAllPerformers(ctx context.Context) ([]dto.PerformerDTO, error) {
	performers, err := p.performerRepo.All(ctx)
	if err != nil {
		p.logg.LogE(msg.E3209, err)

		return nil, fmt.Errorf("%s: %v", msg.E3209, err)
	}

	var performersDTO []dto.PerformerDTO
	for _, performer := range performers {
		performersDTO = append(performersDTO, p.toPerformerDTO(performer))
	}

	return performersDTO, nil
}

func (p *PerformerService) AuthPerformer(ctx context.Context, id int, password string) (*model.AuthPerformer, error) {
	if id <= 0 || password == "" {
		p.logg.LogE(msg.E3211, nil)

		return &model.AuthPerformer{Success: false, Message: msg.E3211}, errors.New("ТН или пароль не должны быть пустыми")
	}

	authOK, err := p.performerRepo.AuthByIdAndPass(ctx, id, password)
	if err != nil || !authOK {
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

func (p *PerformerService) toPerformerDTO(performer model.Performer) dto.PerformerDTO {
	return dto.PerformerDTO{
		Id:           performer.Id,
		FIO:          performer.FIO,
		BC:           performer.BC,
		Archive:      performer.Archive,
		IdRoleAForms: performer.IdRoleAForms,
		IdRoleAFGW:   performer.IdRoleAFGW,
		Audit: dto.AuditDTO(model.Audit{
			CreatedAt: performer.AuditRec.CreatedAt,
			CreatedBy: 6680,
			UpdatedAt: performer.AuditRec.UpdatedAt,
			UpdatedBy: 6680,
		}),
	}
}
