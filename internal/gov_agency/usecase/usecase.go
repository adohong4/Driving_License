package usecase

import (
	"context"

	"github.com/adohong4/driving-license/config"
	govagency "github.com/adohong4/driving-license/internal/gov_agency"
	"github.com/adohong4/driving-license/internal/models"
	"github.com/adohong4/driving-license/pkg/httpErrors"
	"github.com/adohong4/driving-license/pkg/logger"
	"github.com/adohong4/driving-license/pkg/utils"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type GovAgencyUC struct {
	cfg           *config.Config
	GovAgencyRepo govagency.Repository
	logger        logger.Logger
}

func NewGovAgencyUseCase(cfg *config.Config, GovAgencyRepo govagency.Repository, logger logger.Logger) govagency.UseCase {
	return &GovAgencyUC{cfg: cfg, GovAgencyRepo: GovAgencyRepo, logger: logger}
}

func (u *GovAgencyUC) CreateGovAgency(ctx context.Context, gov *models.GovAgency) (*models.GovAgency, error) {
	if err := gov.PrepareCreate(); err != nil {
		return nil, httpErrors.NewBadRequestError(errors.Wrap(err, "GovAgencyUC.CreateGovAgency.PrepareCreate"))
	}

	if err := utils.ValidateStruct(ctx, gov); err != nil {
		return nil, httpErrors.NewBadRequestError(errors.Wrap(err, "GovAgencyUC.CreateGovAgency.ValidateStruct"))
	}

	n, err := u.GovAgencyRepo.CreateGovAgency(ctx, gov)
	if err != nil {
		return nil, err
	}
	return n, nil
}

func (u *GovAgencyUC) UpdateGovAgency(ctx context.Context, gov *models.GovAgency) (*models.GovAgency, error) {
	if err := gov.PrepareCreate(); err != nil {
		return nil, httpErrors.NewBadRequestError(errors.Wrap(err, "GovAgencyUC.CreateGovAgency.PrepareCreate"))
	}

	if err := utils.ValidateStruct(ctx, gov); err != nil {
		return nil, httpErrors.NewBadRequestError(errors.WithMessage(err, "GovAgencyUC.UpdateGovAgency.ValidateStruct"))
	}

	updatedGovAgency, err := u.GovAgencyRepo.UpdateGovAgency(ctx, gov)
	if err != nil {
		return nil, err
	}

	return updatedGovAgency, nil
}

func (u *GovAgencyUC) DeleteGovAgency(ctx context.Context, gov *models.GovAgency) (*models.GovAgency, error) {
	if err := gov.PrepareCreate(); err != nil {
		return nil, httpErrors.NewBadRequestError(errors.Wrap(err, "GovAgencyUC.CreateGovAgency.PrepareCreate"))
	}

	if err := utils.ValidateStruct(ctx, gov); err != nil {
		return nil, httpErrors.NewBadRequestError(errors.WithMessage(err, "GovAgencyUC.DeleteGovAgency.ValidateStruct"))
	}

	deleteGovAgency, err := u.GovAgencyRepo.DeleteGovAgency(ctx, gov)
	if err != nil {
		return nil, err
	}

	return deleteGovAgency, nil
}

func (u *GovAgencyUC) GetGovAgency(ctx context.Context, pq *utils.PaginationQuery) (*models.GovAgencyList, error) {
	return u.GovAgencyRepo.GetGovAgency(ctx, pq)
}

func (u *GovAgencyUC) GetGovAgencyByID(ctx context.Context, Id uuid.UUID) (*models.GovAgency, error) {
	n, err := u.GovAgencyRepo.GetGovAgencyByID(ctx, Id)
	if err != nil {
		return nil, err
	}
	return n, nil
}

func (u *GovAgencyUC) SearchByName(ctx context.Context, name string, query *utils.PaginationQuery) (*models.GovAgencyList, error) {
	return u.GovAgencyRepo.SearchByName(ctx, name, query)
}
