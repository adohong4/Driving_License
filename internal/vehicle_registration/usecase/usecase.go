package usecase

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"github.com/adohong4/driving-license/config"
	"github.com/adohong4/driving-license/internal/models"
	vehicleRegistration "github.com/adohong4/driving-license/internal/vehicle_registration"
	"github.com/adohong4/driving-license/pkg/httpErrors"
	"github.com/adohong4/driving-license/pkg/logger"
	"github.com/adohong4/driving-license/pkg/utils"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type vehicleRegUC struct {
	cfg            *config.Config
	vehicleRegRepo vehicleRegistration.Repository
	logger         logger.Logger
}

// Vehicle Registration Usecase Constructor
func NewVehicleRegUseCase(cfg *config.Config, vehicleRegRepo vehicleRegistration.Repository, log logger.Logger) vehicleRegistration.UseCase {
	return &vehicleRegUC{cfg: cfg, vehicleRegRepo: vehicleRegRepo, logger: log}
}

func (v *vehicleRegUC) CreateVehicleDoc(ctx context.Context, veDoc *models.VehicleRegistration) (*models.VehicleRegistration, error) {
	existsVehiclePlateNO, err := v.vehicleRegRepo.FindVehiclePlateNO(ctx, veDoc)
	if existsVehiclePlateNO != nil {
		return nil, httpErrors.NewRestErrorWithMessage(http.StatusBadRequest, httpErrors.ErrVehicleAlreadyExists, nil)
	}
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, errors.Wrap(err, "vehicleRegUC.CreateVehicleDoc.FindVehiclePlateNO")
	}

	if err = veDoc.PrepareCreate(); err != nil {
		return nil, httpErrors.NewBadRequestError(errors.Wrap(err, "vehicleRegUC.CreateVehicleDoc.PrepareCreate"))
	}

	user, err := utils.GetUserFromCtx(ctx)
	if err != nil {
		return nil, httpErrors.NewUnauthorizedError(errors.WithMessage(err, "vehicleRegUC.Create.GetUserFromCtx"))
	}

	veDoc.CreatorId = user.Id

	if err = utils.ValidateStruct(ctx, veDoc); err != nil {
		return nil, httpErrors.NewBadRequestError(errors.WithMessage(err, "vehicleRegUC.Create.ValidateStruct"))
	}

	n, err := v.vehicleRegRepo.CreateVehicleDoc(ctx, veDoc)
	if err != nil {
		return nil, err
	}

	return n, nil
}

func (v *vehicleRegUC) UpdateVehicleDoc(ctx context.Context, veDoc *models.VehicleRegistration) (*models.VehicleRegistration, error) {
	user, err := utils.GetUserFromCtx(ctx)
	if err != nil {
		return nil, httpErrors.NewUnauthorizedError(errors.WithMessage(err, "vehicleRegUC.UpdateVehicleDoc.GetUserFromCtx"))
	}

	veDoc.ModifierId = &user.Id
	veDoc.UpdatedAt = time.Now()

	if err := utils.ValidateStruct(ctx, veDoc); err != nil {
		return nil, httpErrors.NewBadRequestError(errors.WithMessage(err, "vehicleRegUC.UpdateVehicleDoc.ValidateStruct"))
	}

	updatedVeReg, err := v.vehicleRegRepo.UpdateVehicleDoc(ctx, veDoc)
	if err != nil {
		return nil, err
	}

	return updatedVeReg, nil
}

func (v *vehicleRegUC) DeleteVehicleDoc(ctx context.Context, veDoc *models.VehicleRegistration) (*models.VehicleRegistration, error) {
	user, err := utils.GetUserFromCtx(ctx)
	if err != nil {
		return nil, httpErrors.NewUnauthorizedError(errors.WithMessage(err, "vehicleRegUC.DeleteVehicleDoc.GetUserFromCtx"))
	}

	veDoc.ModifierId = &user.Id
	veDoc.Active = false

	if err = utils.ValidateStruct(ctx, veDoc); err != nil {
		return nil, httpErrors.NewBadRequestError(errors.WithMessage(err, "vehicleRegUC.DeleteVehicleDoc.ValidateStruct"))
	}

	deletedVeReg, err := v.vehicleRegRepo.DeleteVehicleDoc(ctx, veDoc)
	if err != nil {
		return nil, err
	}

	return deletedVeReg, nil
}

func (v *vehicleRegUC) GetVehicleDocs(ctx context.Context, pq *utils.PaginationQuery) (*models.VehicleRegistrationList, error) {
	return v.vehicleRegRepo.GetVehicleDocs(ctx, pq)
}

func (v *vehicleRegUC) GetVehicleByID(ctx context.Context, vehicleID uuid.UUID) (*models.VehicleRegistration, error) {
	n, err := v.vehicleRegRepo.GetVehicleByID(ctx, vehicleID)
	if err != nil {
		return nil, err
	}
	return n, nil
}

func (v *vehicleRegUC) FindByVehiclePlateNO(ctx context.Context, vePlaNO string, query *utils.PaginationQuery) (*models.VehicleRegistrationList, error) {
	return v.vehicleRegRepo.SearchByVehiclePlateNO(ctx, vePlaNO, query)
}
