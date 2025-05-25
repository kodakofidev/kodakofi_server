package repositories

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kodakofidev/kodakofi_server/internal/models"
)

type UserRepoInterface interface {
	GetAllUsers(ctx context.Context, search string) (models.Users, error)
}

type RepoUser struct {
	DB *pgxpool.Pool
}

func NewUser(db *pgxpool.Pool) *RepoUser {
	return &RepoUser{DB: db}
}

func (r *RepoUser) GetAllUsers(ctx context.Context, search string) (models.Users, error) {
	query := `
		SELECT 
			p.fullname, p.phone, p.address, p.image,
			u.email, u.role, u.is_verified, u.created_at
		FROM users u
		JOIN profiles p ON p.user_id = u.id
	`

	// Search filter
	var args []any
	if search != "" {
		query += `
			WHERE 
				u.email ILIKE $1 OR 
				u.role ILIKE $1 OR 
				p.fullname ILIKE $1 OR 
				p.phone ILIKE $1
		`
		args = append(args, "%"+search+"%")
	}

	query += " ORDER BY u.created_at DESC"

	rows, err := r.DB.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users models.Users
	for rows.Next() {
		var (
			user      models.User
			createdAt time.Time
		)
		if err := rows.Scan(
			&user.Fullname,
			&user.Phone,
			&user.Address,
			&user.Image,
			&user.Email,
			&user.Role,
			&user.IsVerified,
			&createdAt,
		); err != nil {
			return nil, err
		}

		// Format waktu ke string (contoh: 2006-01-02 15:04)
		user.CreatedAt = createdAt.Format("2006-01-02 15:04")

		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
