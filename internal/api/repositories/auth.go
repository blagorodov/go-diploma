package auth

import (
	"diploma/internal/drivers"
	"diploma/internal/models"
	"errors"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

type Repository struct {
	db drivers.Database
}

func NewRepository(db drivers.Database) Repository {
	return Repository{
		db: db,
	}
}

func (r *Repository) Get(login string) (*models.User, error) {
	var user models.User
	result := r.db.DB.Where("Login = ?", login).First(&user)
	if result.Error != nil {
		// check for unique violation error
		var pgErr *pgconn.PgError
		if errors.As(result.Error, &pgErr) && pgErr.Code == pgerrcode.NoDataFound {
			return nil, models.ErrLoginOrPasswordNotFound
		}
	}
	return &user, result.Error
}

func (r *Repository) Register(user *models.User) error {
	result := r.db.DB.Create(user)
	if result.Error != nil {
		// check for unique violation error
		var pgErr *pgconn.PgError
		if errors.As(result.Error, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return models.ErrLoginUniqueViolation
		}
	}
	// other errors or no error
	return result.Error
}
