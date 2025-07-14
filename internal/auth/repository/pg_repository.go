package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/adohong4/driving-license/internal/auth"
	"github.com/adohong4/driving-license/internal/models"
	"github.com/adohong4/driving-license/pkg/utils"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// Auth Repository
type authRepo struct {
	db *sqlx.DB
}

// Auth new constructor
func NewAuthRepository(db *sqlx.DB) auth.Repository {
	return &authRepo{db: db}
}

func (r *authRepo) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	u := &models.User{}
	if err := r.db.QueryRowxContext(ctx, createUserQuery,
		user.Id, user.IdentityNo, user.Password, user.HashPassword, user.Active, user.Role,
		user.Version, user.CreatorId, user.ModifierId, user.CreatedAt, user.UpdatedAt,
	).StructScan(u); err != nil {
		return nil, errors.Wrap(err, "authRepo.CreateUser.StructScan")
	}
	return u, nil
}

func (r *authRepo) Update(ctx context.Context, user *models.User) (*models.User, error) {
	u := &models.User{}
	if err := r.db.QueryRowxContext(ctx, updateUserQuery,
		user.IdentityNo, user.Password, user.HashPassword, user.Active, user.Role,
		user.CreatorId, user.ModifierId, user.Id, user.Version,
	).StructScan(u); err != nil {
		return nil, errors.Wrap(err, "authRepo.Update.QueryRowxContext")
	}
	return u, nil
}

func (r *authRepo) Delete(ctx context.Context, id uuid.UUID, modifierId uuid.UUID, version int) error {
	result, err := r.db.ExecContext(ctx, deleteUserQuery, modifierId, time.Now(), id, version)
	if err != nil {
		return errors.Wrap(err, "authRepo.Delete.ExecContext")
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "authRepo.Delete.RowsAffected")
	}
	if rowsAffected == 0 {
		return errors.Wrap(sql.ErrNoRows, "authRepo.Delete.rowsAffected")
	}
	return nil
}

func (r *authRepo) GetUserById(ctx context.Context, id uuid.UUID) (*models.User, error) {
	user := &models.User{}
	if err := r.db.QueryRowxContext(ctx, getUserQuery, id).StructScan(user); err != nil {
		return nil, errors.Wrap(err, "authRepo.GetByID.QueryRowxContext")
	}
	return user, nil
}

func (r *authRepo) FindByIdentityNO(ctx context.Context, identity string, query *utils.PaginationQuery) (*models.UsersList, error) {
	var totalCount int
	if err := r.db.GetContext(ctx, &totalCount, getTotalCount, identity); err != nil {
		return nil, errors.Wrap(err, "authRepo.FindByIdentityNO.GetContext.totalCount")
	}

	if totalCount == 0 {
		return &models.UsersList{
			TotalCount: totalCount,
			TotalPages: utils.GetTotalPage(totalCount, query.GetSize()),
			Page:       query.GetPage(),
			Size:       query.GetSize(),
			HasMore:    utils.GetHasMore(query.GetPage(), totalCount, query.GetSize()),
			Users:      make([]*models.User, 0),
		}, nil
	}

	var users []*models.User
	if err := r.db.SelectContext(ctx, &users, findUsers, identity, query.GetOffset(), query.GetLimit()); err != nil {
		return nil, errors.Wrap(err, "authRepo.FindByIdentityNO.SelectContext")
	}

	return &models.UsersList{
		TotalCount: totalCount,
		TotalPages: utils.GetTotalPage(totalCount, query.GetSize()),
		Page:       query.GetPage(),
		Size:       query.GetSize(),
		HasMore:    utils.GetHasMore(query.GetPage(), totalCount, query.GetSize()),
		Users:      users,
	}, nil
}

func (r *authRepo) GetUsers(ctx context.Context, pq *utils.PaginationQuery) (*models.UsersList, error) {
	var totalCount int
	if err := r.db.GetContext(ctx, &totalCount, getTotal); err != nil {
		return nil, errors.Wrap(err, "authRepo.GetUsers.GetContext.totalCount")
	}

	if totalCount == 0 {
		return &models.UsersList{
			TotalCount: totalCount,
			TotalPages: utils.GetTotalPage(totalCount, pq.GetSize()),
			Page:       pq.GetPage(),
			Size:       pq.GetSize(),
			HasMore:    utils.GetHasMore(pq.GetPage(), totalCount, pq.GetSize()),
			Users:      make([]*models.User, 0),
		}, nil
	}

	var users []*models.User
	if err := r.db.SelectContext(ctx, &users, getUsers, pq.GetOrderBy(), pq.GetOffset(), pq.GetLimit()); err != nil {
		return nil, errors.Wrap(err, "authRepo.GetUsers.SelectContext")
	}

	return &models.UsersList{
		TotalCount: totalCount,
		TotalPages: utils.GetTotalPage(totalCount, pq.GetSize()),
		Page:       pq.GetPage(),
		Size:       pq.GetSize(),
		HasMore:    utils.GetHasMore(pq.GetPage(), totalCount, pq.GetSize()),
		Users:      users,
	}, nil
}

func (r *authRepo) FindByIdentity(ctx context.Context, user *models.User) (*models.User, error) {
	foundUser := &models.User{}
	err := r.db.QueryRowxContext(ctx, findUserByIdentity, user.IdentityNo).StructScan(foundUser)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, errors.Wrap(err, "authRepo.FindByIdentity.QueryRowxContext")
	}
	return foundUser, nil
}
