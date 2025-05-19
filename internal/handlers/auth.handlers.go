package handlers

import (
	"fmt"
	"log"
	"net/http"
	"net/url" // Changed from "url" to "net/url"
	"os"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"github.com/kodakofidev/kodakofi_server/internal/models"
	"github.com/kodakofidev/kodakofi_server/internal/repositories"
	"github.com/kodakofidev/kodakofi_server/internal/services"
	"github.com/kodakofidev/kodakofi_server/internal/utils"
	"github.com/kodakofidev/kodakofi_server/pkg"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
)

type AuthHandlers struct {
	repo   repositories.AuthRepoInterface
	mailer services.MailerService
}

func NewAuthHandlers(repo repositories.AuthRepoInterface) *AuthHandlers {
	return &AuthHandlers{
		repo:   repo,
		mailer: services.NewGmailMailer(),
	}
}

func (a *AuthHandlers) Login(ctx *gin.Context) {
	var userReq models.UserReq
	response := models.NewResponse(ctx)
	if err := ctx.ShouldBindJSON(&userReq); err != nil {
		if strings.Contains(err.Error(), "Field validation for 'Password'") {
			response.BadRequest("Password length must be at least 8 characters", err.Error())
			return
		}
		log.Println(err.Error())
		response.BadRequest("Invalid input", err.Error())
		return
	}
	if !isValidEmail(userReq.Email) {
		response.BadRequest("Invalid email format", nil)
		return
	}
	result, err := a.repo.Login(ctx.Request.Context(), userReq.Email)
	if err != nil {
		response.InternalServerError("Failed to login", err.Error())
		return
	}
	if result.Email == "" {
		response.BadRequest("Email or password is incorrect", nil)
		return
	}
	hash := pkg.InitHashConfig()
	hash.UseDefaultConfig()
	valid, err := hash.CompareHashAndPassword(result.Pass, userReq.Password)
	if err != nil {
		log.Println(err.Error())
		response.InternalServerError("Internal Server Error", err.Error())
		return
	}
	if !valid {
		response.BadRequest("Email or password is incorrect", nil)
		return
	}

	// Check if user is verified
	if !result.IsVerified {
		response.BadRequest("Email not verified. Please verify your email first.", nil)
		return
	}

	payload := pkg.NewJWT(result.AuthID, result.Email, result.Role)
	token, err := payload.GenerateToken()
	if err != nil {
		response.InternalServerError("Failed to generate token", err.Error())
		return
	}
	// Response now includes id as authID
	response.Success("Login successful", map[string]string{
		"token": token,
		"role":  result.Role,
		"email": result.Email,
		"id":    result.AuthID,
	})
}

func (h *AuthHandlers) Register(ctx *gin.Context) {
	var userReq models.UserReq
	response := models.NewResponse(ctx)
	if err := ctx.ShouldBindJSON(&userReq); err != nil {
		if strings.Contains(err.Error(), "Field validation for 'Password'") {
			response.BadRequest("Password length must be at least 8 characters", err.Error())
			return
		}
		log.Println(err.Error())
		response.BadRequest("Invalid input", err.Error())
		return
	}
	if !isValidEmail(userReq.Email) {
		response.BadRequest("Invalid email format", nil)
		return
	}
	hash := pkg.InitHashConfig()
	hash.UseDefaultConfig()
	hashedPass, err := hash.GenHashedPassword(userReq.Password)
	if err != nil {
		response.InternalServerError("Failed to hash password", err.Error())
		return
	}
	result, findEmail, err := h.repo.Register(ctx.Request.Context(), userReq, hashedPass)
	if err != nil {
		response.InternalServerError("Failed to register user", err.Error())
		return
	}
	if findEmail.Email != "" {
		response.BadRequest("Email already exists", nil)
		return
	}

	// Generate OTP for email verification
	otp := utils.GenerateOTP()
	expiry := utils.GenerateOTPExpiry()

	// Store OTP in database (type 1 for email verification)
	if err := h.repo.StoreOTP(ctx.Request.Context(), result.AuthID, otp, 1, expiry); err != nil {
		response.InternalServerError("Failed to generate verification code", err.Error())
		return
	}

	// Send verification email
	if err := h.mailer.SendVerificationEmail(result.Email, otp); err != nil {
		log.Printf("Failed to send verification email: %v", err)
		response.Success("User registered successfully, but failed to send verification email", map[string]string{"email": result.Email})
		return
	}

	response.Success("User registered successfully. Please check your email for verification code.", map[string]string{"email": result.Email})
}

func (h *AuthHandlers) VerifyEmail(ctx *gin.Context) {
	var req models.OTPVerificationReq
	response := models.NewResponse(ctx)

	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest("Invalid input", err.Error())
		return
	}

	// Log verification attempt for debugging
	log.Printf("Attempting to verify OTP: %s for email: %s", req.OTP, req.Email)

	// Verify OTP - pass the type_id parameter (1 for email verification)
	valid, err := h.repo.VerifyOTP(ctx.Request.Context(), req.Email, req.OTP, req.TypeID)
	if err != nil {
		log.Printf("OTP verification failed: %v", err)
		response.BadRequest("Verification failed", err.Error())
		return
	}

	if !valid {
		response.BadRequest("Invalid verification code", nil)
		return
	}

	// Get user ID
	userResult, err := h.repo.Login(ctx.Request.Context(), req.Email)
	if err != nil || userResult.AuthID == "" {
		response.BadRequest("User not found", nil)
		return
	}

	// Update user verification status
	if err := h.repo.UpdateUserVerificationStatus(ctx.Request.Context(), userResult.AuthID); err != nil {
		response.InternalServerError("Failed to verify user", err.Error())
		return
	}

	response.Success("Email verified successfully", nil)
}

func (h *AuthHandlers) Logout(ctx *gin.Context) {
	// Implementation for logout
}

func (h *AuthHandlers) SendOTP(ctx *gin.Context) {
	var req models.OTPSend
	response := models.NewResponse(ctx)

	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest("Invalid input", err.Error())
		return
	}

	if !isValidEmail(req.Email) {
		response.BadRequest("Invalid email format", nil)
		return
	}

	// Check if user exists
	userResult, err := h.repo.Login(ctx.Request.Context(), req.Email)
	if err != nil {
		response.InternalServerError("Failed to find user", err.Error())
		return
	}

	if userResult.Email == "" {
		response.BadRequest("Email not registered", nil)
		return
	}

	// Generate new OTP
	otp := utils.GenerateOTP()
	expiry := utils.GenerateOTPExpiry()

	// Store OTP in database
	if err := h.repo.StoreOTP(ctx.Request.Context(), userResult.AuthID, otp, req.TypeID, expiry); err != nil {
		response.InternalServerError("Failed to generate verification code", err.Error())
		return
	}

	// Send verification email
	if err := h.mailer.SendOTPEmail(req.Email, otp, req.TypeID); err != nil {
		log.Printf("Failed to send email with OTP: %v", err)
		response.InternalServerError("Failed to send email", err.Error())
		return
	}

	response.Success("Verification code sent successfully", map[string]string{"email": req.Email})
}

func isValidEmail(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}

const (
	key    = "kodakofi_auth_key"
	maxAge = 60 * 60 * 24
	isProd = false
)

func InitAuth() {
	googleClientID := os.Getenv("GOOGLE_CLIENT_ID")
	googleClientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	callbackURL := os.Getenv("GOOGLE_CALLBACK_URL")

	// Fallback to GOOGLE_REDIRECT_URL if GOOGLE_CALLBACK_URL is not set
	if callbackURL == "" {
		callbackURL = os.Getenv("GOOGLE_REDIRECT_URL")
	}

	log.Printf("Initializing OAuth providers with:")
	log.Printf("  GOOGLE_CLIENT_ID: %s", googleClientID)
	log.Printf("  GOOGLE_CLIENT_SECRET: %s", strings.Repeat("*", len(googleClientSecret)))
	log.Printf("  CALLBACK_URL: %s", callbackURL)

	if googleClientID == "" || googleClientSecret == "" || callbackURL == "" {
		log.Printf("WARNING: Missing OAuth environment variables. OAuth login will not work!")
		return
	}

	// Clear any existing providers to avoid conflicts
	goth.ClearProviders()

	// Setup cookie store for sessions
	store := sessions.NewCookieStore([]byte(key))
	store.MaxAge(maxAge)
	store.Options = &sessions.Options{
		Path:     "/",
		HttpOnly: true,
		Secure:   isProd,
	}

	// Set the session store for gothic
	gothic.Store = store

	// Register the Google provider
	provider := google.New(
		googleClientID,
		googleClientSecret,
		callbackURL,
		"email", "profile",
	)
	goth.UseProviders(provider)

	log.Printf("Google OAuth provider registered successfully")
}

// GoogleLogin initiates OAuth flow with Google
func (a *AuthHandlers) GoogleLogin(ctx *gin.Context) {
	// Clear session to prevent stale data
	session, _ := gothic.Store.Get(ctx.Request, gothic.SessionName)
	session.Values = map[any]any{}
	session.Save(ctx.Request, ctx.Writer)

	// Set the provider explicitly
	ctx.Request.URL.RawQuery = "provider=google"

	log.Printf("Starting Google OAuth flow with URL: %s", ctx.Request.URL.String())

	// Begin the OAuth flow
	gothic.BeginAuthHandler(ctx.Writer, ctx.Request)
}

// GoogleCallback handles the callback from Google OAuth
func (a *AuthHandlers) GoogleCallback(ctx *gin.Context) {
	response := models.NewResponse(ctx)

	// Ensure provider is set in query parameters
	q := ctx.Request.URL.Query()
	q.Set("provider", "google")
	ctx.Request.URL.RawQuery = q.Encode()

	log.Printf("Processing OAuth callback for provider: google")

	// Complete the OAuth flow
	user, err := gothic.CompleteUserAuth(ctx.Writer, ctx.Request)
	if err != nil {
		log.Printf("Error during OAuth callback: %v", err)

		// Return the error as JSON response instead of redirecting
		response.BadRequest("Authentication failed", err.Error())
		return
	}

	log.Printf("Google auth successful for email: %s", user.Email)

	// Check if user exists in our database
	result, err := a.repo.Login(ctx.Request.Context(), user.Email)
	if err != nil {
		// If there's an error but it's not because the user doesn't exist
		if err.Error() != "no rows in result set" {
			response.InternalServerError("Failed to check user", err.Error())
			return
		}
	}

	// User does not exist, register them
	if result.Email == "" {
		log.Printf("Registering new user from Google OAuth: %s", user.Email)

		// Generate a random password for OAuth users
		hash := pkg.InitHashConfig()
		hash.UseDefaultConfig()
		randomPass := utils.GenerateRandomString(12)

		hashedPass, err := hash.GenHashedPassword(randomPass)
		if err != nil {
			response.InternalServerError("Failed to hash password", err.Error())
			return
		}

		// Create new user request with name from profile
		userReq := models.UserReq{
			Email:    user.Email,
			Password: randomPass,
		}

		// Register the user
		newUser, _, err := a.repo.Register(ctx.Request.Context(), userReq, hashedPass)
		if err != nil {
			response.InternalServerError("Failed to register user", err.Error())
			return
		}

		// Since this is OAuth, we can mark them as verified directly
		if err := a.repo.UpdateUserVerificationStatus(ctx.Request.Context(), newUser.AuthID); err != nil {
			response.InternalServerError("Failed to verify user", err.Error())
			return
		}

		log.Printf("Successfully registered user from Google OAuth: %s", user.Email)

		// Update result to use the newly created user
		result = newUser
	} else {
		log.Printf("Existing user logged in with Google OAuth: %s", user.Email)
	}

	// Generate JWT token for the user
	payload := pkg.NewJWT(result.AuthID, result.Email, result.Role)
	token, err := payload.GenerateToken()
	if err != nil {
		response.InternalServerError("Failed to generate token", err.Error())
		return
	}

	// Check if this is API call or a redirect
	if ctx.Request.Header.Get("X-Requested-With") == "XMLHttpRequest" {
		// API call, return JSON with id field
		response.Success("Google authentication successful", map[string]string{
			"token": token,
			"email": user.Email,
			"name":  user.Name,
			"id":    result.AuthID, // Include the real user ID
			"role":  result.Role,
		})
	} else {
		// Browser flow, redirect to frontend with token, user ID, and role
		frontendURL := os.Getenv("FRONTEND_URL")
		if frontendURL == "" {
			frontendURL = "http://localhost:5173"
		}

		redirectURL := fmt.Sprintf("%s/auth/callback?token=%s&email=%s&name=%s&id=%s&role=%s",
			frontendURL,
			url.QueryEscape(token),
			url.QueryEscape(user.Email),
			url.QueryEscape(user.Name),
			url.QueryEscape(result.AuthID),
			url.QueryEscape(result.Role))

		log.Printf("Redirecting to: %s", redirectURL)
		ctx.Redirect(http.StatusFound, redirectURL)
	}
}

// TokenLogin handles login with a Google ID token
func (a *AuthHandlers) TokenLogin(ctx *gin.Context) {
	var req struct {
		IDToken string `json:"id_token" binding:"required"`
	}
	response := models.NewResponse(ctx)

	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest("Invalid input", err.Error())
		return
	}

	// Validate the token
	userInfo, err := utils.ValidateGoogleIDToken(req.IDToken)
	if err != nil {
		response.BadRequest("Invalid token", err.Error())
		return
	}

	// Check if user exists in our database
	result, err := a.repo.Login(ctx.Request.Context(), userInfo.Email)
	if err != nil {
		// If there's an error but it's not because the user doesn't exist
		if err.Error() != "no rows in result set" {
			response.InternalServerError("Failed to check user", err.Error())
			return
		}
	}

	// User does not exist, register them
	if result.Email == "" {
		// Generate a random password for OAuth users
		hash := pkg.InitHashConfig()
		hash.UseDefaultConfig()
		randomPass := utils.GenerateRandomString(12)

		hashedPass, err := hash.GenHashedPassword(randomPass)
		if err != nil {
			response.InternalServerError("Failed to hash password", err.Error())
			return
		}

		// Create new user request
		userReq := models.UserReq{
			Email:    userInfo.Email,
			Password: randomPass,
		}

		// Register the user
		newUser, _, err := a.repo.Register(ctx.Request.Context(), userReq, hashedPass)
		if err != nil {
			response.InternalServerError("Failed to register user", err.Error())
			return
		}

		// Since this is OAuth, we can mark them as verified directly
		if err := a.repo.UpdateUserVerificationStatus(ctx.Request.Context(), newUser.AuthID); err != nil {
			response.InternalServerError("Failed to verify user", err.Error())
			return
		}

		// Update result to use the newly created user
		result = newUser
	}

	// Generate JWT token for the user
	payload := pkg.NewJWT(result.AuthID, result.Email, result.Role)
	token, err := payload.GenerateToken()
	if err != nil {
		response.InternalServerError("Failed to generate token", err.Error())
		return
	}

	response.Success("Google authentication successful", map[string]any{
		"token": token,
		"user": map[string]any{
			"email":   userInfo.Email,
			"name":    userInfo.Name,
			"picture": userInfo.Picture,
		},
	})
}
