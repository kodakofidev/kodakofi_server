package repositories

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kodakofidev/kodakofi_server/internal/models"
)

type ProfileRepoInterface interface {
	GetProfile(ctx *context.Context, Id string) (*models.Profile, error)
	EditProfile()
}

type RepoProfile struct {
	DB *pgxpool.Pool
}

func NewProfile(db *pgxpool.Pool) *RepoProfile {
	return &RepoProfile{DB: db}
}

func (p *RepoProfile) GetProfile(ctx context.Context, Id string) (*models.Profile, error) {
	query := `SELECT fullname, phone, address, image FROM profiles WHERE user_id = $1`

	var profile models.Profile
	err := p.DB.QueryRow(ctx, query, Id).Scan(
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

func (p *RepoProfile) EditProfile() {

}
