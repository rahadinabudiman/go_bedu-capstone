package usecase

import (
	"bufio"
	"errors"
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
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/thanhpk/randstr"
	"golang.org/x/crypto/bcrypt"
)

type AdminUsecase interface {
	LoginAdmin(c echo.Context, req dtos.LoginRequest) (res dtos.LoginResponse, err error)
	LogoutAdmin(c echo.Context) (res dtos.LogoutAdminResponse, err error)
	VerifyEmail(verificationCode any) (res dtos.VerifyEmailResponse, err error)
	UpdateAdminByOTP(otp int, req dtos.ChangePasswordRequest) (res dtos.ForgotPasswordResponse, err error)
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
			Username:  admin.Username,
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
// @Summary      Login Admin with Username and Password
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
// @Router       /admin/login [post]
func (u *adminUsecase) LoginAdmin(c echo.Context, req dtos.LoginRequest) (res dtos.LoginResponse, err error) {
	admin, err := u.adminRepository.GetAdminByUsername(req.Username)
	if err != nil {
		return res, errors.New("Username not registered")
	}

	if !admin.Verified {
		return res, errors.New("Please verify your email first")
	}

	err = helpers.ComparePassword(req.Password, admin.Password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return res, errors.New("Email or Password is wrong")
	}

	token, err := middlewares.CreateToken(int(admin.ID), admin.Username, admin.Email, admin.Role)
	if err != nil {
		return res, errors.New("Failed to generate token")

	}

	admin.Token = token

	middlewares.CreateCookie(c, token)

	res = dtos.LoginResponse{
		Username: admin.Username,
		Token:    admin.Token,
	}

	return res, nil
}

// LogoutAdmin godoc
// @Summary      Logout Administrator
// @Description  Logout Administrator
// @Tags         Admin - Account
// @Accept       json
// @Produce      json
// @Success      200 {object} dtos.LogoutAdminOKResponse
// @Failure      400 {object} dtos.BadRequestResponse
// @Failure      401 {object} dtos.UnauthorizedResponse
// @Failure      403 {object} dtos.ForbiddenResponse
// @Failure      404 {object} dtos.NotFoundResponse
// @Failure      500 {object} dtos.InternalServerErrorResponse
// @Router       /admin/logout [get]
// @Security BearerAuth
func (u *adminUsecase) LogoutAdmin(c echo.Context) (res dtos.LogoutAdminResponse, err error) {
	err = middlewares.DeleteCookie(c)
	if err != nil {
		return res, errors.New("Failed to logout")
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
// @Router       /admin/verifyemail/{verificationCode} [get]
func (u *adminUsecase) VerifyEmail(verificationCode any) (res dtos.VerifyEmailResponse, err error) {
	admin, err := u.adminRepository.GetAdminByVerificationCode(verificationCode)
	if err != nil {
		return res, errors.New("Cannot get admin")
	}

	if admin.Verified {
		return res, errors.New("Email already verified")
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
// @Router       /admin/change-password/{otp} [post]
func (u *adminUsecase) UpdateAdminByOTP(otp int, req dtos.ChangePasswordRequest) (res dtos.ForgotPasswordResponse, err error) {
	admin, err := u.adminRepository.GetAdminOTP(otp)
	if err != nil {
		return res, errors.New("Failed to get admin")
	}

	if req.Password != req.PasswordConfirm {
		return res, errors.New("Password not matches")
	}

	// Reset OTP and OTPReq
	admin.OTP = 0
	admin.OTPReq = false

	// Update Password
	passwordHash, err := helpers.HashPassword(req.Password)
	if err != nil {
		return res, errors.New("Failed to hash password")
	}
	admin.Password = string(passwordHash)

	_, err = u.adminRepository.UpdateAdmin(admin)
	if err != nil {
		return res, errors.New("Failed to update admin")
	}

	res = dtos.ForgotPasswordResponse{
		Email:   admin.Email,
		Message: "Password has been reset successfully",
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
		return res, errors.New("Failed to get admin")
	}

	err = helpers.ComparePassword(req.OldPassword, admin.Password)
	if err != nil {
		return res, errors.New("Wrong Password")
	}

	if req.Password != req.PasswordConfirm {
		return res, errors.New("Password not matches")
	}

	passwordHash, err := helpers.HashPassword(req.Password)
	if err != nil {
		return res, errors.New("Failed to hash password")
	}

	admin.Password = string(passwordHash)

	admin, err = u.adminRepository.UpdateAdmin(admin)
	if err != nil {
		return res, errors.New("Failed to update Admin")
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
// @Router       /admin/register [post]
func (u *adminUsecase) CreateAdmin(req *dtos.RegisterAdminRequest) (dtos.AdminDetailResponse, error) {
	var res dtos.AdminDetailResponse

	err := helpers.ValidateUsername(req.Username)
	if err != nil {
		return res, echo.NewHTTPError(400, err)
	}

	username, _ := u.adminRepository.GetAdminByUsername(req.Username)
	if username.ID > 0 {
		return res, errors.New("Username already in use")
	}

	req.Email = strings.ToLower(req.Email)

	// Check apakah email sudah terdaftar atau belum
	admin, _ := u.adminRepository.GetAdminByEmail(req.Email)
	if admin.ID > 0 {
		return res, errors.New("Email already in use")
	}

	if req.Password != req.PasswordConfirm {
		return res, errors.New("Password does not matches")
	}

	passwordHash, err := helpers.HashPassword(req.Password)
	if err != nil {
		return res, err
	}

	config, err := initializers.LoadConfig(".")
	if err != nil {
		return res, errors.New("Failed to load config")
	}

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

	var firstName = CreateAdmin.Username

	if strings.Contains(firstName, " ") {
		firstName = strings.Split(firstName, " ")[1]
	}

	// 👇 Send Email
	emailData := utils.EmailData{
		URL:       config.ClientOrigin + "/admin/verifyemail/" + url.PathEscape(verification_code),
		FirstName: firstName,
		Subject:   "Your account verification code",
	}

	utils.SendEmail(&CreateAdmin, &emailData)

	resp := dtos.AdminDetailResponse{
		ID:        admins.ID,
		Username:  admins.Username,
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
// @Security BearerAuth
func (u *adminUsecase) GetAdminById(id uint) (res dtos.AdminProfileResponse, err error) {
	admin, err := u.adminRepository.GetAdminById(id)
	if err != nil {
		return res, errors.New("Admin didn't exist")
	}

	res = dtos.AdminProfileResponse{
		ID:       admin.ID,
		Username: admin.Username,
		Nama:     admin.Nama,
		Email:    admin.Email,
		Role:     admin.Role,
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
	admins.Username = req.Username

	// Check Role and save role information from JWT Cookie
	admin, err := u.adminRepository.ReadToken(id)
	if err != nil {
		return res, errors.New("Failed to get admin")
	}

	// Check if admin is not Super Admin and trying to change role
	if admin.Role != req.Role && admin.Role != "Super Admin" {
		return res, errors.New("You are not allowed to change role")
	}

	// Check if Super Admin changes other Role
	if admin.Role == "Super Admin" {
		if req.Role != "Super Admin" && req.Role != "Admin" {
			return res, errors.New("Cannot change role between Super Admin and Admin")
		}
	}

	admins.ID = uint(id)

	passwordHash, err := helpers.HashPassword(admins.Password)
	if err != nil {
		return res, errors.New("Failed to hash password")
	}

	admins.Password = string(passwordHash)

	admins, err = u.adminRepository.UpdateAdmin(admins)
	if err != nil {
		return res, err
	}

	res.Username = admins.Username
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
		return res, errors.New("Failed to get admin")
	}

	err = helpers.ComparePassword(req.Password, admin.Password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		echo.NewHTTPError(400, err.Error())
		return
	}

	err = u.adminRepository.DeleteAdmin(admin)

	if err != nil {
		return res, errors.New("Failed to delete admin")
	}

	return res, nil
}
