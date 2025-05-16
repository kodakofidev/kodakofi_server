package pkg

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims represents JWT claims data - renamed and exported
type Claims struct {
	Uuid  string `json:"uuid"`
	Email string `json:"email"`
	Role  string `json:"role"`
	jwt.RegisteredClaims
}

type JWTErr struct {
	Type string
	Err  error
}

func NewJWT(uuid, email, role string) *Claims {
	return &Claims{
		Uuid:  uuid,
		Email: email,
		Role:  role,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "koda kofi",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 60)),
		},
	}
}

func (c *Claims) GenerateToken() (string, error) {
	secret := os.Getenv("JWT_SECRET")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	return token.SignedString([]byte(secret))
}

func (c *Claims) VerifyToken(token string) JWTErr {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return JWTErr{
			Type: "System",
			Err:  errors.New("Secret not provided"),
		}
	}
	parsedToken, err := jwt.ParseWithClaims(token, c, func(t *jwt.Token) (any, error) {
		return []byte(jwtSecret), nil
	})
	if err != nil {
		return JWTErr{
			Type: "Token",
			Err:  err,
		}
	}
	if !parsedToken.Valid {
		return JWTErr{
			Type: "Token",
			Err:  errors.New("Expired Token"),
		}
	}

	return JWTErr{}
}
