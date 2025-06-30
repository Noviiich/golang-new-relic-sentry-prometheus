package postgres

import (
	"database/sql"
	"fmt"

	"github.com/Noviiich/golang-new-relic-sentry-prometheus/internal/domain"
	"github.com/Noviiich/golang-new-relic-sentry-prometheus/internal/repository"
	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db *sqlx.DB
}

// NewUserRepository создает новый репозиторий пользователей
func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) CreateUser(user domain.User) (domain.User, error) {
	const op = "repository.postgres.CreateUser"

	query := `
		INSERT INTO users (name, age, created_date)
		VALUES ($1, $2, $3)
		RETURNING id
	`

	var id int
	err := r.db.QueryRow(query, user.Name, user.Age, user.CreatedDate).Scan(&id)
	if err != nil {
		return domain.User{}, fmt.Errorf("%s: %w", op, repository.ErrUserCreationFailed)
	}

	user.ID = id

	return user, nil
}

func (r *UserRepository) GetUserByID(id int) (domain.User, error) {
	const op = "repository.postgres.GetUserByID"

	query := `
		SELECT id, name, age, created_date
		FROM users
		WHERE id = $1
	`

	var user domain.User
	if err := r.db.Get(&user, query, id); err != nil {
		if err == sql.ErrNoRows {
			return domain.User{}, fmt.Errorf("%s: %w", op, repository.ErrUserNotFound)
		}
		return domain.User{}, fmt.Errorf("%s: %w", op, repository.ErrUserRetrievalFailed)
	}

	return user, nil
}

func (r *UserRepository) UpdateUser(user domain.User) (domain.User, error) {
	const op = "repository.postgres.UpdateUser"

	query := `
		UPDATE users
		SET name = $1, age = $2, created_date = NOW()
		WHERE id = $3
	`

	_, err := r.db.Exec(query, user.Name, user.Age, user.ID)
	if err != nil {
		return domain.User{}, fmt.Errorf("%s: %w", op, repository.ErrUserUpdateFailed)
	}
	return user, nil
}

func (r *UserRepository) DeleteUser(id int) error {
	const op = "repository.postgres.DeleteUser"

	query := `
		DELETE FROM users
		WHERE id = $1
	`

	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, repository.ErrUserDeletionFailed)
	}
	return nil
}
