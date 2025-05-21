package repositories

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kodakofidev/kodakofi_server/internal/models"
	"github.com/kodakofidev/kodakofi_server/pkg"
)

type ProfileRepoInterface interface {
	GetProfile(ctx context.Context, Id string) (*models.Profile, error)
	EditProfile(ctx context.Context, Id string, profile models.ProfileForm, filePath string) (pgconn.CommandTag, error)
}

type RepoProfile struct {
	DB *pgxpool.Pool
}

func NewProfile(db *pgxpool.Pool) *RepoProfile {
	return &RepoProfile{DB: db}
}

func (p *RepoProfile) GetProfile(ctx context.Context, UserId string) (*models.Profile, error) {
	// query := `SELECT fullname, phone, address, image, created_at, updated_at FROM profiles WHERE user_id = $1`
	query := `SELECT p.fullname, u.email, p.phone, p.address, p.image, p.created_at, p.updated_at
	 			FROM profiles p JOIN users u ON p.user_id = u.id
				WHERE p.user_id = $1`

	var profile models.Profile
	err := p.DB.QueryRow(ctx, query, UserId).Scan(
		&profile.Fullname,
		&profile.Email,
		&profile.Phone,
		&profile.Address,
		&profile.ProfileImage,
		&profile.CreatedAt,
		&profile.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &profile, nil

}

func (p *RepoProfile) EditProfile(ctx context.Context, UserId string, profile models.ProfileForm, filePath string) (pgconn.CommandTag, error) {
	// Start transaction
	tx, err := p.DB.Begin(ctx)
	if err != nil {
		return pgconn.CommandTag{}, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	var (
		queryBuilder strings.Builder
		values       []interface{}
		paramCount   = 1
	)

	queryBuilder.WriteString("UPDATE profiles SET ")

	updates := []struct {
		condition bool
		field     string
		value     interface{}
	}{
		{profile.Fullname != "", "fullname", profile.Fullname},
		{profile.Phone != "", "phone", profile.Phone},
		{profile.Address != "", "address", profile.Address},
		{filePath != "", "image", filePath},
	}

	var needsComma bool
	for _, update := range updates {
		if update.condition {
			if needsComma {
				queryBuilder.WriteString(", ")
			}
			queryBuilder.WriteString(fmt.Sprintf("%s = $%d", update.field, paramCount))
			values = append(values, update.value)
			paramCount++
			needsComma = true
		}
	}

	// Only execute profile update if there are fields to update
	var result pgconn.CommandTag
	if len(values) > 0 {
		queryBuilder.WriteString(fmt.Sprintf(" WHERE user_id = $%d", paramCount))
		values = append(values, UserId)

		var err error
		result, err = tx.Exec(ctx, queryBuilder.String(), values...)
		if err != nil {
			return pgconn.CommandTag{}, fmt.Errorf("failed to update profile: %w", err)
		}
	}

	// Handle password update if requested
	if profile.CurrentPassword != "" {
		// Verify current password
		var storedHash string
		if err := tx.QueryRow(ctx, "SELECT password FROM users WHERE id = $1", UserId).Scan(&storedHash); err != nil {
			log.Println("[DEBUG 1]")
			return pgconn.CommandTag{}, fmt.Errorf("failed to get user password: %w", err)
		}

		hash := pkg.InitHashConfig()
		hash.UseDefaultConfig()

		valid, err := hash.CompareHashAndPassword(storedHash, profile.CurrentPassword)
		if err != nil {
			log.Println("[DEBUG 2]")

			return pgconn.CommandTag{}, fmt.Errorf("password comparison failed: %w", err)

		}
		if !valid {
			return pgconn.CommandTag{}, errors.New("invalid current password")
		}

		// Generate and set new password
		hashedPass, err := hash.GenHashedPassword(profile.NewPassword)
		if err != nil {
			log.Println("[DEBUG 3]")

			return pgconn.CommandTag{}, fmt.Errorf("failed to hash new password: %w", err)
		}

		if _, err := tx.Exec(ctx, "UPDATE users SET password = $1 WHERE id = $2", hashedPass, UserId); err != nil {
			log.Println("[DEBUG 4]")

			return pgconn.CommandTag{}, fmt.Errorf("failed to update password: %w", err)
		}
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return pgconn.CommandTag{}, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return result, nil
}

// func (p *RepoProfile) EditProfile(ctx context.Context, UserId string, profile models.ProfileForm, filePath string) (pgconn.CommandTag, error) {
// 	tx, err := p.DB.Begin(ctx)
// 	if err != nil {
// 		return pgconn.CommandTag{}, nil
// 	}
// 	defer tx.Rollback(ctx)

// 	query := `UPDATE profiles SET`
// 	values := []any{}
// 	clauses := []string{}

// 	if profile.Fullname != "" {
// 		clauses = append(clauses, fmt.Sprintf(`fullname = $%d`, len(values)+1))
// 		values = append(values, profile.Fullname)
// 	}

// 	if profile.Phone != "" {
// 		clauses = append(clauses, fmt.Sprintf(`phone = $%d`, len(values)+1))
// 		values = append(values, profile.Phone)
// 	}

// 	if profile.Address != "" {
// 		clauses = append(clauses, fmt.Sprintf(`address = $%d`, len(values)+1))
// 		values = append(values, profile.Address)
// 	}

// 	if filePath != "" {
// 		clauses = append(clauses, fmt.Sprintf(`image = $%d`, len(values)+1))
// 		values = append(values, filePath)
// 	}

// 	query += " " + strings.Join(clauses, ", ")

// 	query += fmt.Sprintf(` WHERE user_id = $%d`, len(values)+1)
// 	values = append(values, UserId)

// 	result, err := tx.Exec(ctx, query, values...)
// 	if err != nil {
// 		return pgconn.CommandTag{}, err
// 	}

// if profile.CurrentPassword != "" {
// 	var storedHash string
// 	queryPass := `SELECT password FROM users WHERE id = $1`

// 	if err := tx.QueryRow(ctx, queryPass, UserId).Scan(&storedHash); err != nil {
// 		if err == pgx.ErrNoRows {
// 			log.Printf(`[DEBUG] 1 %v`, err)
// 			return pgconn.CommandTag{}, err
// 		}

// 	}
// 	hash := pkg.InitHashConfig()
// 	hash.UseDefaultConfig()
// 	valid, err := hash.CompareHashAndPassword(storedHash, profile.CurrentPassword)
// 	if err != nil {
// 		log.Println(err.Error())
// 		return pgconn.CommandTag{}, err
// 	}
// 	if !valid {
// 		log.Printf(`[DEBUG] 2 %v`, err)
// 		return pgconn.CommandTag{}, err
// 	}

// 	hashedPass, err := hash.GenHashedPassword(profile.NewPassword)
// 	if err != nil {
// 		log.Printf(`[DEBUG] 3 %v`, err)
// 		return pgconn.CommandTag{}, nil
// 	}

// 	queryNewPassword := `UPDATE users SET password = $1 WHERE id = $2`
// 	_, err = tx.Exec(ctx, queryNewPassword, hashedPass, UserId)
// 	if err != nil {
// 		log.Printf(`[DEBUG] 4 %v`, err)

// 		return pgconn.CommandTag{}, err
// 	}
// }

// 	if err := tx.Commit(ctx); err != nil {
// 		return pgconn.CommandTag{}, fmt.Errorf("failed to commit transaction: %w", err)
// 	}

// 	return result, nil
// }
