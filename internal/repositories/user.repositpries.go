package repositories

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kodakofidev/kodakofi_server/internal/models"
	"github.com/kodakofidev/kodakofi_server/internal/utils"
)

type UserRepoInterface interface {
	GetAllUsers(ctx context.Context, search string) (models.Users, error)
	UpdateUserByAdmin(ctx context.Context, req models.UpdateUserByAdminReq) (*models.UserUpdateRes, error)
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
			u.id, u.email, u.role, u.is_verified, u.created_at
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
			image     string
		)
		if err := rows.Scan(
			&user.Fullname,
			&user.Phone,
			&user.Address,
			&image,
			&user.ID,
			&user.Email,
			&user.Role,
			&user.IsVerified,
			&createdAt,
		); err != nil {
			return nil, err
		}

		// Format image URL
		if image != "" {
			user.Image = fmt.Sprintf("%s%s", utils.BaseImgProfileURL, image)
		} else {
			user.Image = image
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

func (r *RepoUser) UpdateUserByAdmin(ctx context.Context, req models.UpdateUserByAdminReq) (*models.UserUpdateRes, error) {
	tx, err := r.DB.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	// Ambil data lama untuk gambar
	var oldImage string
	queryGet := `
	SELECT p.image 
	FROM profiles p 
	WHERE p.user_id = $1
	`
	err = tx.QueryRow(ctx, queryGet, req.ID).Scan(&oldImage)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Update table users
	queryUser := `
	UPDATE users 
	SET role = COALESCE(NULLIF($1, ''), role),
	    updated_at = NOW()
	WHERE id = $2
	`
	_, err = tx.Exec(ctx, queryUser, req.Role, req.ID)
	if err != nil {
		return nil, err
	}

	// Update table profiles
	queryProfile := `
	UPDATE profiles 
	SET fullname = COALESCE(NULLIF($1, ''), fullname),
	    phone = COALESCE(NULLIF($2, ''), phone),
	    address = COALESCE(NULLIF($3, ''), address),
	    image = COALESCE(NULLIF($4, ''), image),
	    updated_at = NOW()
	WHERE user_id = $5
	`
	newImage := ""
	if req.Image != nil {
		newImage = req.Image.Filename
	}
	_, err = tx.Exec(ctx, queryProfile, req.Fullname, req.Phone, req.Address, newImage, req.ID)
	if err != nil {
		return nil, err
	}

	// Hapus gambar lama jika berbeda dan bukan default
	if newImage != "" && oldImage != "" && oldImage != "avatar_default.svg" {
		path := filepath.Join("public/profile-images", oldImage)
		if err := os.Remove(path); err != nil && !errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("failed to delete old image: %w", err)
		}
	}

	// Ambil data hasil update
	query := `
	SELECT u.id, p.fullname, u.email, p.phone, u.role, p.address, p.image, u.is_verified, 
	       u.created_at, p.updated_at
	FROM users u
	JOIN profiles p ON p.user_id = u.id
	WHERE u.id = $1
	`

	var (
		res       models.UserUpdateRes
		createdAt time.Time
		updatedAt time.Time
	)
	err = tx.QueryRow(ctx, query, req.ID).Scan(
		&res.ID,
		&res.Fullname,
		&res.Email,
		&res.Phone,
		&res.Role,
		&res.Address,
		&res.Image,
		&res.IsVerified,
		&createdAt,
		&updatedAt,
	)

	if err != nil {
		return nil, err
	}

	res.CreatedAt = createdAt.Format("2006-01-02 15:04")
	res.UpdatedAt = updatedAt.Format("2006-01-02 15:04")

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return &res, nil
}
