package handlers

import (
	"log"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kodakofidev/kodakofi_server/internal/models"
	"github.com/kodakofidev/kodakofi_server/internal/repositories"
	"github.com/kodakofidev/kodakofi_server/internal/services"
	"github.com/kodakofidev/kodakofi_server/internal/utils"
	"github.com/kodakofidev/kodakofi_server/pkg"
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
	response.Success("Login successful", map[string]string{"token": token})
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
