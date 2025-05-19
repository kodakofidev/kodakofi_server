package repositories

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kodakofidev/kodakofi_server/internal/models"
)

type AuthRepoInterface interface {
	Register(c context.Context, userReq models.UserReq, hashedPass string) (models.UserRes, models.UserRes, error)
	Login(c context.Context, email string) (models.UserRes, error)
	Logout()
	StoreOTP(c context.Context, userID string, otp string, typeID int, expiry time.Time) error
	VerifyOTP(c context.Context, email string, otp string, typeID int) (bool, error) // Added typeID parameter
	UpdateUserVerificationStatus(c context.Context, userID string) error
}

type RepoAuth struct {
	DB *pgxpool.Pool
}

func NewAuth(db *pgxpool.Pool) *RepoAuth {
	return &RepoAuth{DB: db}
}

func (r *RepoAuth) Register(c context.Context, userReq models.UserReq, hashedPass string) (models.UserRes, models.UserRes, error) {
	// Begin transaction
	tx, err := r.DB.Begin(c)
	if err != nil {
		return models.UserRes{}, models.UserRes{}, err
	}
	// Defer rollback in case of error
	defer tx.Rollback(c)

	queryCheckMail := `SELECT email FROM users WHERE email = $1`
	valuesMail := []any{userReq.Email}
	var findUser models.UserRes
	if err := tx.QueryRow(c, queryCheckMail, valuesMail...).Scan(&findUser.Email); err != nil && err != pgx.ErrNoRows {
		return models.UserRes{}, models.UserRes{}, err
	}
	if findUser.Email == userReq.Email {
		return models.UserRes{}, findUser, nil
	}

	query := `INSERT INTO users (email, "password") VALUES ($1, $2) RETURNING email, "role", id;`
	values := []any{userReq.Email, hashedPass}
	var result models.UserRes
	if err := tx.QueryRow(c, query, values...).Scan(&result.Email, &result.Role, &result.AuthID); err != nil {
		return models.UserRes{}, models.UserRes{}, err
	}

	queryProfile := `INSERT INTO profiles (user_id, fullname, phone, address, image) VALUES ($1, $2, $3, $4, $5);`
	valuesProfile := []any{result.AuthID, userReq.Fullname, "", "", ""}
	if _, err := tx.Exec(c, queryProfile, valuesProfile...); err != nil {
		return models.UserRes{}, models.UserRes{}, err
	}

	// Commit transaction
	if err := tx.Commit(c); err != nil {
		return models.UserRes{}, models.UserRes{}, err
	}

	return result, models.UserRes{}, nil
}

func (r *RepoAuth) Login(c context.Context, email string) (models.UserRes, error) {
	query := `
		SELECT u.id, u.email, u."role", u.password, u.is_verified, COALESCE(p.fullname, '') as fullname 
		FROM users u 
		LEFT JOIN profiles p ON u.id = p.user_id 
		WHERE u.email = $1
	`
	values := []any{email}
	var result models.UserRes
	if err := r.DB.QueryRow(c, query, values...).Scan(&result.AuthID, &result.Email, &result.Role, &result.Pass, &result.IsVerified, &result.Fullname); err != nil && err != pgx.ErrNoRows {
		return models.UserRes{}, err
	}
	return result, nil
}

func (r *RepoAuth) StoreOTP(c context.Context, userID string, otp string, typeID int, expiry time.Time) error {
	query := `INSERT INTO code_otp (user_id, code, type_id, expired_at) VALUES ($1, $2, $3, $4)`
	_, err := r.DB.Exec(c, query, userID, otp, typeID, expiry)
	return err
}

func (r *RepoAuth) VerifyOTP(c context.Context, email string, otp string, typeID int) (bool, error) {
	// Get user ID from email
	queryUser := `SELECT id FROM users WHERE email = $1`
	var userID string
	if err := r.DB.QueryRow(c, queryUser, email).Scan(&userID); err != nil {
		if err == pgx.ErrNoRows {
			return false, errors.New("user not found")
		}
		return false, err
	}

	// First check if the OTP exists at all (regardless of expiration)
	checkOTPQuery := `
		SELECT id, expired_at FROM code_otp 
		WHERE user_id = $1 AND code = $2 AND type_id = $3
		ORDER BY created_at DESC LIMIT 1
	` // Modified query to use parameterized type_id
	var otpID int
	var expiry time.Time
	err := r.DB.QueryRow(c, checkOTPQuery, userID, otp, typeID).Scan(&otpID, &expiry)
	if err != nil {
		if err == pgx.ErrNoRows {
			return false, errors.New("invalid OTP code")
		}
		return false, err
	}

	// Check if the OTP is expired
	if time.Now().After(expiry) {
		return false, errors.New("OTP has expired, please request a new one")
	}

	return true, nil
}

func (r *RepoAuth) UpdateUserVerificationStatus(c context.Context, userID string) error {
	query := `UPDATE users SET is_verified = true WHERE id = $1`
	_, err := r.DB.Exec(c, query, userID)
	return err
}

func (r *RepoAuth) Logout() {
	// Implementation for logout
}
