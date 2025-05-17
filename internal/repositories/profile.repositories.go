package repositories

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kodakofidev/kodakofi_server/internal/models"
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
	query := `SELECT fullname, phone, address, image FROM profiles WHERE user_id = $1`

	var profile models.Profile
	err := p.DB.QueryRow(ctx, query, UserId).Scan(
		&profile.Fullname,
		&profile.Phone,
		&profile.Address,
		&profile.ProfileImage,
	)
	if err != nil {
		return nil, err
	}

	return &profile, nil

}

func (p *RepoProfile) EditProfile(ctx context.Context, UserId string, profile models.ProfileForm, filePath string) (pgconn.CommandTag, error) {
	query := `UPDATE profiles SET`
	values := []any{}
	clauses := []string{}

	if profile.Fullname != "" {
		clauses = append(clauses, fmt.Sprintf(`fullname = $%d`, len(values)+1))
		values = append(values, profile.Fullname)
	}

	if profile.Phone != "" {
		clauses = append(clauses, fmt.Sprintf(`phone = $%d`, len(values)+1))
		values = append(values, profile.Phone)
	}

	if profile.Address != "" {
		clauses = append(clauses, fmt.Sprintf(`address = $%d`, len(values)+1))
		values = append(values, profile.Address)
	}

	if filePath != "" {
		clauses = append(clauses, fmt.Sprintf(`image = $%d`, len(values)+1))
		values = append(values, filePath)
	}

	query += " " + strings.Join(clauses, ", ")

	query += fmt.Sprintf(`WHERE user_id = $%d`, len(values)+1)
	values = append(values, UserId)

	result, err := p.DB.Exec(ctx, query, values...)
	if err != nil {
		return pgconn.CommandTag{}, err
	}
	return result, nil
}
