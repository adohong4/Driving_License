package vehicleRegistration

import (
	"context"

	"github.com/adohong4/driving-license/internal/models"
	"github.com/adohong4/driving-license/pkg/utils"
	"github.com/google/uuid"
)

type UseCase interface {
	CreateVehicleDoc(ctx context.Context, veDoc *models.VehicleRegistration) (*models.VehicleRegistration, error)
	UpdateVehicleDoc(ctx context.Context, veDoc *models.VehicleRegistration) (*models.VehicleRegistration, error)
	DeleteVehicleDoc(ctx context.Context, veDoc *models.VehicleRegistration) (*models.VehicleRegistration, error)
	GetVehicleDocs(ctx context.Context, pq *utils.PaginationQuery) (*models.VehicleRegistrationList, error)
	GetVehicleByID(ctx context.Context, vehicleID uuid.UUID) (*models.VehicleRegistration, error)
	FindByVehiclePlateNO(ctx context.Context, vePlaNO string, query *utils.PaginationQuery) (*models.VehicleRegistrationList, error)
}
