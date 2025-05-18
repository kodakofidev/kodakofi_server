package utils

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// GetUserInfo fetches user information from Google using an access token
// Kept for backward compatibility
func GetUserInfo(accessToken string) (map[string]any, error) {
	userInfoEndpoint := "https://www.googleapis.com/oauth2/v2/userinfo"
	req, err := http.NewRequest("GET", userInfoEndpoint, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var userInfo map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, err
	}
	return userInfo, nil
}

// SignJWT creates a JWT token from user information
func SignJWT(userInfo map[string]any) (string, error) {
	// Extract user ID or use email as ID if not available
	userID, ok := userInfo["id"]
	if !ok {
		userID = userInfo["email"]
	}

	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		secretKey = "your-secret-key" // Fallback for development
	}

	claims := jwt.MapClaims{
		"sub":   userID,
		"name":  userInfo["name"],
		"email": userInfo["email"],
		"iss":   "kodakofi-server",
		"exp":   time.Now().Add(time.Hour * 72).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

// SignJWTFromOAuthInfo creates a JWT token from structured OAuth user info
func SignJWTFromOAuthInfo(userInfo *OAuthUserInfo) (string, error) {
	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		secretKey = "your-secret-key" // Fallback for development
	}

	claims := jwt.MapClaims{
		"sub":     userInfo.ID,
		"email":   userInfo.Email,
		"name":    userInfo.Name,
		"picture": userInfo.Picture,
		"iss":     "kodakofi-server",
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}
