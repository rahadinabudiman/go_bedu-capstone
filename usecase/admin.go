package usecase

import (
	"bufio"
	"fmt"
	"go_bedu/dtos"
	"go_bedu/helpers"
	"go_bedu/initializers"
	"go_bedu/middlewares"
	"go_bedu/models"
	"go_bedu/repositories"
	"go_bedu/utils"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	emailverifier "github.com/AfterShip/email-verifier"
	"github.com/labstack/echo/v4"
	"github.com/thanhpk/randstr"
	"golang.org/x/crypto/bcrypt"
)

type AdminUsecase interface {
	LoginAdmin(c echo.Context, req dtos.LoginRequest) (res dtos.LoginResponse, err error)
	LogoutAdmin(c echo.Context) (res dtos.LogoutAdminResponse, err error)
	VerifyEmail(verificationCode any) (res dtos.VerifyEmailResponse, err error)
	UpdateAdminByOTP(otp int, req dtos.ChangePasswordRequest) (res dtos.ForgotPasswordResponse, err error)
	ForgotPassword(req dtos.ForgotPasswordRequest) (res dtos.ForgotPasswordResponse, err error)
	ChangePassword(id uint, req dtos.ChangePasswordAdminRequest) (res helpers.ResponseMessage, err error)
	MustDispEmailDom() (dispEmailDomains []string, err error)
	GetAdmin() ([]dtos.AdminDetailResponse, error)
	GetAdminById(id uint) (res dtos.AdminProfileResponse, err error)
	UpdateAdmin(id uint, req dtos.UpdateAdminRequest) (res dtos.UpdateAdminResponse, err error)
	CreateAdmin(req *dtos.RegisterAdminRequest) (dtos.AdminDetailResponse, error)
	DeleteAdmin(id uint, req dtos.DeleteAdminRequest) (res helpers.ResponseMessage, err error)
}

type adminUsecase struct {
	adminRepository repositories.AdminRepository
}

func NewAdminUsecase(adminRepository repositories.AdminRepository) *adminUsecase {
	return &adminUsecase{adminRepository}
}

func (u *adminUsecase) MustDispEmailDom() (dispEmailDomains []string, err error) {
	file, err := os.Open("utils/disposable_email_blocklist.txt")
	if err != nil {
		return nil, err
	}

	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		dispEmailDomains = append(dispEmailDomains, scanner.Text())
	}
	return dispEmailDomains, nil
}

// GetAllAdmins godoc
// @Summary      Get all admins
// @Description  Get all admins
// @Tags         Admin - Account
// @Accept       json
// @Produce      json
// @Success      200 {object} dtos.GetAllAdminsResponse
// @Failure      400 {object} dtos.BadRequestResponse
// @Failure      401 {object} dtos.UnauthorizedResponse
// @Failure      403 {object} dtos.ForbiddenResponse
// @Failure      404 {object} dtos.NotFoundResponse
// @Failure      500 {object} dtos.InternalServerErrorResponse
// @Router       /admin [get]
// @Security BearerAuth
func (u *adminUsecase) GetAdmin() ([]dtos.AdminDetailResponse, error) {
	admins, err := u.adminRepository.GetAdmins()
	if err != nil {
		return nil, err
	}

	var adminResponse []dtos.AdminDetailResponse
	for _, admin := range admins {
		adminResponse = append(adminResponse, dtos.AdminDetailResponse{
			ID:        admin.ID,
			Nama:      admin.Nama,
			Email:     admin.Email,
			Role:      admin.Role,
			CreatedAt: admin.CreatedAt,
			UpdatedAt: admin.UpdatedAt,
			// Article:   admin.Articles,
		})
	}

	return adminResponse, nil
}

// AdminLogin godoc
// @Summary      Login Admin with Email and Password
// @Description  Login an account
// @Tags         Admin - Auth
// @Accept       json
// @Produce      json
// @Param        request body dtos.LoginRequest true "Payload Body [RAW]"
// @Success      200 {object} dtos.LoginStatusOKResponse
// @Failure      400 {object} dtos.BadRequestResponse
// @Failure      401 {object} dtos.UnauthorizedResponse
// @Failure      403 {object} dtos.ForbiddenResponse
// @Failure      404 {object} dtos.NotFoundResponse
// @Failure      500 {object} dtos.InternalServerErrorResponse
// @Router       /login [post]
func (u *adminUsecase) LoginAdmin(c echo.Context, req dtos.LoginRequest) (res dtos.LoginResponse, err error) {
	admin, err := u.adminRepository.GetAdminByUsername(req.Username)
	if err != nil {
		echo.NewHTTPError(400, "Username not registered")
		return
	}

	if !admin.Verified {
		return res, echo.NewHTTPError(400, "Please verify your email first")
	}

	err = helpers.ComparePassword(req.Password, admin.Password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		echo.NewHTTPError(400, err.Error())
		return
	}

	token, err := middlewares.CreateToken(int(admin.ID), admin.Username, admin.Email, admin.Role)
	if err != nil {
		echo.NewHTTPError(400, "Failed to generate token")
		return
	}

	admin.Token = token

	middlewares.CreateCookie(c, token)

	res = dtos.LoginResponse{
		Username: admin.Username,
		Token:    admin.Token,
	}

	return
}

func (u *adminUsecase) LogoutAdmin(c echo.Context) (res dtos.LogoutAdminResponse, err error) {
	err = middlewares.DeleteCookie(c)
	if err != nil {
		return res, echo.NewHTTPError(400, "Failed to logout")
	}

	return res, err
}

// AdminVerif godoc
// @Summary      Verify Email by Verification Code
// @Description  Verif an account
// @Tags         Admin - Auth
// @Accept       json
// @Produce      json
// @Param verification_code path string true "Verification Code"
// @Success      200 {object} dtos.VerifyEmailOKResponse
// @Failure      400 {object} dtos.BadRequestResponse
// @Failure      401 {object} dtos.UnauthorizedResponse
// @Failure      403 {object} dtos.ForbiddenResponse
// @Failure      404 {object} dtos.NotFoundResponse
// @Failure      500 {object} dtos.InternalServerErrorResponse
// @Router       /verifyemail/{verificationCode} [get]
func (u *adminUsecase) VerifyEmail(verificationCode any) (res dtos.VerifyEmailResponse, err error) {
	admin, err := u.adminRepository.GetAdminByVerificationCode(verificationCode)
	if err != nil {
		return res, echo.NewHTTPError(400, "Failed to get admin")
	}

	if admin.Verified {
		return res, echo.NewHTTPError(400, "Email already verified")
	}

	admin.VerificationCode = ""
	admin.Verified = true

	u.adminRepository.UpdateAdmin(admin)

	res = dtos.VerifyEmailResponse{
		Username: admin.Username,
		Message:  "Email has been verified",
	}

	return res, nil
}

// UpdateAdminByOTP godoc
// @Summary      Change Password by OTP
// @Description  Change Password an Account
// @Tags         Admin - Auth
// @Accept       json
// @Produce      json
// @Param        request body dtos.ChangePasswordRequest true "Payload Body [RAW]"
// @Success      200 {object} dtos.ChangePasswordOKResponse
// @Failure      400 {object} dtos.BadRequestResponse
// @Failure      401 {object} dtos.UnauthorizedResponse
// @Failure      403 {object} dtos.ForbiddenResponse
// @Failure      404 {object} dtos.NotFoundResponse
// @Failure      500 {object} dtos.InternalServerErrorResponse
// @Router       /change-password/{otp} [post]
func (u *adminUsecase) UpdateAdminByOTP(otp int, req dtos.ChangePasswordRequest) (res dtos.ForgotPasswordResponse, err error) {
	admin, err := u.adminRepository.GetAdminOTP(otp)
	if err != nil {
		return res, echo.NewHTTPError(400, "Failed to get admin")
	}

	if req.Password != req.PasswordConfirm {
		return res, echo.NewHTTPError(400, "Password not matches")
	}

	// Reset OTP and OTPReq
	admin.OTP = 0
	admin.OTPReq = false

	// Update Password
	passwordHash, err := helpers.HashPassword(req.Password)
	if err != nil {
		return res, echo.NewHTTPError(400, "Failed to hash password")
	}
	admin.Password = string(passwordHash)

	_, err = u.adminRepository.UpdateAdmin(admin)
	if err != nil {
		return res, echo.NewHTTPError(400, "Failed to update admin")
	}

	res = dtos.ForgotPasswordResponse{
		Email:   admin.Email,
		Message: "Password has been reset successfully",
	}

	return res, nil
}

// ForgotPassword godoc
// @Summary      Forgot Password Request OTP
// @Description  Forgot Password an Account
// @Tags         Admin - Auth
// @Accept       json
// @Produce      json
// @Param        request body dtos.ForgotPasswordRequest true "Payload Body [RAW]"
// @Success      200 {object} dtos.ForgotPasswordOKResponse
// @Failure      400 {object} dtos.BadRequestResponse
// @Failure      401 {object} dtos.UnauthorizedResponse
// @Failure      403 {object} dtos.ForbiddenResponse
// @Failure      404 {object} dtos.NotFoundResponse
// @Failure      500 {object} dtos.InternalServerErrorResponse
// @Router       /forgot-password [post]
func (u *adminUsecase) ForgotPassword(req dtos.ForgotPasswordRequest) (res dtos.ForgotPasswordResponse, err error) {
	admin, err := u.adminRepository.GetAdminByEmail(req.Email)
	if err != nil {
		return res, echo.NewHTTPError(400, "Email not registered")
	}

	config, _ := initializers.LoadConfig(".")

	// Generate OTP
	otp := helpers.GenerateRandomOTP(6)
	NewOTP, err := strconv.Atoi(otp)
	if err != nil {
		return res, echo.NewHTTPError(400, "Failed to generate OTP")
	}
	admin.OTP = NewOTP
	admin.OTPReq = true

	u.adminRepository.UpdateAdmin(admin)

	// 👇 Send Email
	emailData := utils.EmailData{
		URL:       "http://" + config.ClientOrigin + "/change-password/" + url.PathEscape(otp),
		FirstName: admin.Nama,
		Subject:   "Your OTP to reset password",
	}

	utils.SendEmail(&admin, &emailData)

	res = dtos.ForgotPasswordResponse{
		Email:   admin.Email,
		Message: "OTP has been sent to your email",
	}

	return res, nil
}

// ChangePassword godoc
// @Summary      Change Password Admin
// @Description  Change Password Admin
// @Tags         Admin - Account
// @Accept       json
// @Produce      json
// @Param        request body dtos.ChangePasswordAdminRequest true "Payload Body [RAW]"
// @Success      200 {object} dtos.ChangePasswordAdminOKResponse
// @Failure      400 {object} dtos.BadRequestResponse
// @Failure      401 {object} dtos.UnauthorizedResponse
// @Failure      403 {object} dtos.ForbiddenResponse
// @Failure      404 {object} dtos.NotFoundResponse
// @Failure      500 {object} dtos.InternalServerErrorResponse
// @Router       /admin/change-password [post]
// @Security BearerAuth
func (u *adminUsecase) ChangePassword(id uint, req dtos.ChangePasswordAdminRequest) (res helpers.ResponseMessage, err error) {
	admin, err := u.adminRepository.GetAdminById(id)
	if err != nil {
		return res, echo.NewHTTPError(400, "Failed to get admin")
	}

	err = helpers.ComparePassword(req.OldPassword, admin.Password)
	if err != nil {
		return res, echo.NewHTTPError(400, "Wrong password")
	}

	if req.Password != req.PasswordConfirm {
		return res, echo.NewHTTPError(400, "Password not matches")
	}

	passwordHash, err := helpers.HashPassword(req.Password)
	if err != nil {
		return res, echo.NewHTTPError(400, "Failed to hash password")
	}

	admin.Password = string(passwordHash)

	admin, err = u.adminRepository.UpdateAdmin(admin)
	if err != nil {
		return res, echo.NewHTTPError(400, "Failed to update admin")
	}

	res = helpers.NewResponseMessage(
		http.StatusOK,
		"Password has been changed successfully",
	)

	return res, nil
}

// AdminRegister godoc
// @Summary      Register Admin
// @Description  Register an account
// @Tags         Admin - Auth
// @Accept       json
// @Produce      json
// @Param        request body dtos.RegisterAdminRequest true "Payload Body [RAW]"
// @Success      201 {object} dtos.AdminCreeatedResponse
// @Failure      400 {object} dtos.BadRequestResponse
// @Failure      401 {object} dtos.UnauthorizedResponse
// @Failure      403 {object} dtos.ForbiddenResponse
// @Failure      404 {object} dtos.NotFoundResponse
// @Failure      500 {object} dtos.InternalServerErrorResponse
// @Router       /register [post]
func (u *adminUsecase) CreateAdmin(req *dtos.RegisterAdminRequest) (dtos.AdminDetailResponse, error) {
	var (
		verifier = emailverifier.
				NewVerifier().
				EnableAutoUpdateDisposable().
				EnableSMTPCheck().
				DisableCatchAllCheck()

		res dtos.AdminDetailResponse
	)

	err := helpers.ValidateUsername(req.Username)
	if err != nil {
		return res, echo.NewHTTPError(400, err)
	}

	username, _ := u.adminRepository.GetAdminByUsername(req.Username)
	if username.ID > 0 {
		return res, echo.NewHTTPError(400, "Username already in use")
	}

	req.Email = strings.ToLower(req.Email)
	emailDomain := utils.GetEmailDomain(req.Email)
	usernameDomain := utils.GetEmailUsername(req.Email)

	// Check apakah domain email typo atau tidak
	suggestion := verifier.SuggestDomain(emailDomain)
	if suggestion != "" {
		return res, echo.NewHTTPError(400, "Did you mean "+suggestion+"?")
	}

	// Check SMTP apakah domain email valid atau tidak
	ret, err := verifier.CheckSMTP(emailDomain, usernameDomain)
	if err != nil {
		return res, echo.NewHTTPError(400, err)
	}
	fmt.Println("smtp validation result: ", ret)

	// Check apakah email disposable atau tidak
	if verifier.IsDisposable(emailDomain) {
		return res, echo.NewHTTPError(400, "Sorry, we do not accept disposable email addresses")
	}

	// Check apakah email sudah terdaftar atau belum
	admin, _ := u.adminRepository.GetAdminByEmail(req.Email)
	if admin.ID > 0 {
		return res, echo.NewHTTPError(400, "Email already in use")
	}

	if req.Password != req.PasswordConfirm {
		return res, echo.NewHTTPError(400, "Password does not match")
	}

	passwordHash, err := helpers.HashPassword(req.Password)
	if err != nil {
		return res, err
	}

	config, _ := initializers.LoadConfig(".")

	// Generate Verification Code
	code := randstr.String(20)
	verification_code := utils.Encode(code)

	CreateAdmin := models.Administrator{
		Nama:             req.Nama,
		Username:         req.Username,
		Email:            req.Email,
		Password:         passwordHash,
		VerificationCode: verification_code,
	}
	admins, err := u.adminRepository.CreateAdmin(CreateAdmin)

	if err != nil {
		return res, err
	}

	var firstName = CreateAdmin.Nama

	if strings.Contains(firstName, " ") {
		firstName = strings.Split(firstName, " ")[1]
	}

	// 👇 Send Email
	emailData := utils.EmailData{
		URL:       "http://" + config.ClientOrigin + "/verifyemail/" + url.PathEscape(verification_code),
		FirstName: firstName,
		Subject:   "Your account verification code",
	}

	utils.SendEmail(&CreateAdmin, &emailData)

	resp := dtos.AdminDetailResponse{
		ID:        admins.ID,
		Nama:      admins.Nama,
		Email:     admins.Email,
		Role:      admins.Role,
		CreatedAt: admins.CreatedAt,
		UpdatedAt: admins.UpdatedAt,
	}

	return resp, nil
}

// GetAdminByID godoc
// @Summary      Get admin by ID
// @Description  Get admin by ID
// @Tags         Admin - Account
// @Accept       json
// @Produce      json
// @Param id path integer true "ID admin"
// @Success      200 {object} dtos.AdminStatusOKResponse
// @Failure      400 {object} dtos.BadRequestResponse
// @Failure      401 {object} dtos.UnauthorizedResponse
// @Failure      403 {object} dtos.ForbiddenResponse
// @Failure      404 {object} dtos.NotFoundResponse
// @Failure      500 {object} dtos.InternalServerErrorResponse
// @Router       /admin/{id} [get]
// @Security BearerAut
func (u *adminUsecase) GetAdminById(id uint) (res dtos.AdminProfileResponse, err error) {
	admin, err := u.adminRepository.GetAdminById(id)

	if err != nil {
		echo.NewHTTPError(401, "This routes for admin only")
		return
	}

	res = dtos.AdminProfileResponse{
		ID:    admin.ID,
		Nama:  admin.Nama,
		Email: admin.Email,
		Role:  admin.Role,
	}

	return res, nil
}

// AdminUpdate godoc
// @Summary      Update Information
// @Description  Admin update an information
// @Tags         Admin - Account
// @Accept       json
// @Produce      json
// @Param        request body dtos.UpdateAdminRequest true "Payload Body [RAW]"
// @Success      200 {object} dtos.AdminStatusOKResponse
// @Failure      400 {object} dtos.BadRequestResponse
// @Failure      401 {object} dtos.UnauthorizedResponse
// @Failure      403 {object} dtos.ForbiddenResponse
// @Failure      404 {object} dtos.NotFoundResponse
// @Failure      500 {object} dtos.InternalServerErrorResponse
// @Router       /admin/{id} [put]
// @Security BearerAuth
func (u *adminUsecase) UpdateAdmin(id uint, req dtos.UpdateAdminRequest) (dtos.UpdateAdminResponse, error) {
	var (
		admins models.Administrator
		res    dtos.UpdateAdminResponse
	)

	admins, err := u.adminRepository.GetAdminById(id)
	if err != nil {
		return res, err
	}

	admins.Nama = req.Nama
	admins.Email = req.Email
	admins.Password = req.Password
	admins.Role = req.Role

	// Check Role and save role information from JWT Cookie
	admin, err := u.adminRepository.ReadToken(id)
	if err != nil {
		return res, echo.NewHTTPError(400, "Failed to get Admin")
	}

	// Check if admin is not Super Admin and trying to change role
	if admin.Role != req.Role && admin.Role != "Super Admin" {
		return res, echo.NewHTTPError(401, "You are not allowed to change role")
	}

	// Check if Super Admin changes other Role
	if admin.Role == "Super Admin" {
		if req.Role != "Super Admin" && req.Role != "Admin" {
			return res, echo.NewHTTPError(401, "Tidak bisa mengubah role menjadi selain Super Admin atau Admin")
		}
	}

	admins.ID = uint(id)

	passwordHash, err := helpers.HashPassword(admins.Password)
	if err != nil {
		return res, echo.NewHTTPError(400, "Failed to hash password")
	}

	admins.Password = string(passwordHash)

	admins, err = u.adminRepository.UpdateAdmin(admins)
	if err != nil {
		return res, err
	}

	res.Nama = admins.Nama
	res.Email = admins.Email
	res.Password = admins.Password
	res.Role = admins.Role

	return res, nil
}

// DeleteAdmin godoc
// @Summary      Delete an Admin
// @Description  Delete an Admin
// @Tags         Admin - Account
// @Accept       json
// @Produce      json
// @Param        request body dtos.DeleteAdminRequest true "Payload Body [RAW]"
// @Success      200 {object} dtos.StatusOKDeletedResponse
// @Failure      400 {object} dtos.BadRequestResponse
// @Failure      401 {object} dtos.UnauthorizedResponse
// @Failure      403 {object} dtos.ForbiddenResponse
// @Failure      404 {object} dtos.NotFoundResponse
// @Failure      500 {object} dtos.InternalServerErrorResponse
// @Router       /admin/{id} [delete]
// @Security BearerAuth
func (u *adminUsecase) DeleteAdmin(id uint, req dtos.DeleteAdminRequest) (res helpers.ResponseMessage, err error) {
	admin, err := u.adminRepository.ReadToken(id)

	if err != nil {
		echo.NewHTTPError(400, "Failed to get Admin")
		return
	}

	err = helpers.ComparePassword(req.Password, admin.Password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		echo.NewHTTPError(400, err.Error())
		return
	}

	err = u.adminRepository.DeleteAdmin(admin)

	if err != nil {
		return res, echo.NewHTTPError(500, "Failed to delete admin")
	}

	return res, nil
}
