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

type UserUsecase interface {
	LoginUser(c echo.Context, req dtos.LoginRequest) (res dtos.LoginResponse, err error)
	LogoutUser(c echo.Context) (res dtos.LogoutUserResponse, err error)
	VerifyEmail(verificationCode any) (res dtos.VerifyEmailResponse, err error)
	UpdateUserByOTP(otp int, req dtos.ChangePasswordRequest) (res dtos.ForgotPasswordResponse, err error)
	ChangePassword(id uint, req dtos.ChangePasswordUserRequest) (res helpers.ResponseMessage, err error)
	MustDispEmailDom() (dispEmailDomains []string, err error)
	GetUsers() ([]dtos.UserDetailResponse, error)
	GetUserById(id uint) (res dtos.UserProfileResponse, err error)
	CreateUser(req *dtos.RegisterUserRequest) (dtos.UserDetailResponse, error)
	UpdateUser(id uint, req dtos.UpdateUserRequest) (res dtos.UpdateUserResponse, err error)
	DeleteUser(id uint, req dtos.DeleteUserRequest) (res helpers.ResponseMessage, err error)
}

type userUsecase struct {
	userRepository repositories.UserRepository
}

func NewUserUsecase(userRepository repositories.UserRepository) *userUsecase {
	return &userUsecase{userRepository}
}

func (u *userUsecase) MustDispEmailDom() (dispEmailDomains []string, err error) {
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

// UserLogin godoc
// @Summary      Login User with Username and Password
// @Description  Login an account
// @Tags         User - Auth
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
func (u *userUsecase) LoginUser(c echo.Context, req dtos.LoginRequest) (res dtos.LoginResponse, err error) {
	user, err := u.userRepository.GetUserByUsername(req.Username)
	if err != nil {
		return res, errors.New("Username not found")
	}

	if !user.Verified {
		return res, errors.New("Please verify your email first")
	}

	err = helpers.ComparePassword(req.Password, user.Password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return res, errors.New("Email or password is incorrect")

	}

	token, err := middlewares.CreateToken(int(user.ID), user.Username, user.Email, user.Role)
	if err != nil {
		return res, errors.New("Failed to generate token")

	}

	user.Token = token

	middlewares.CreateCookie(c, token)

	res = dtos.LoginResponse{
		Username: user.Username,
		Token:    user.Token,
	}

	return res, nil
}

// LogoutUser godoc
// @Summary      Logout User
// @Description  Logout User
// @Tags         User - Account
// @Accept       json
// @Produce      json
// @Success      200 {object} dtos.LogoutUserResponse
// @Failure      400 {object} dtos.BadRequestResponse
// @Failure      401 {object} dtos.UnauthorizedResponse
// @Failure      403 {object} dtos.ForbiddenResponse
// @Failure      404 {object} dtos.NotFoundResponse
// @Failure      500 {object} dtos.InternalServerErrorResponse
// @Router       /user/logout [get]
// @Security BearerAuth
func (u *userUsecase) LogoutUser(c echo.Context) (res dtos.LogoutUserResponse, err error) {
	err = middlewares.DeleteCookie(c)
	if err != nil {
		return res, errors.New("Failed to logout")
	}

	return res, err
}

// UserVerify godoc
// @Summary      Verify Email by Verification Code
// @Description  Verif an account
// @Tags         User - Auth
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
func (u *userUsecase) VerifyEmail(verificationCode any) (res dtos.VerifyEmailResponse, err error) {
	user, err := u.userRepository.GetUserByVerificationCode(verificationCode)
	if err != nil {
		return res, errors.New("Failed to get user")
	}

	if user.Verified {
		return res, errors.New("Email already verified")
	}

	user.VerificationCode = ""
	user.Verified = true

	u.userRepository.UpdateUser(user)

	res = dtos.VerifyEmailResponse{
		Username: user.Username,
		Message:  "Email has been verified",
	}

	return res, nil
}

// UpdateUserOTP godoc
// @Summary      Change Password by OTP
// @Description  Change Password an Account
// @Tags         User - Auth
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
func (u *userUsecase) UpdateUserByOTP(otp int, req dtos.ChangePasswordRequest) (res dtos.ForgotPasswordResponse, err error) {
	user, err := u.userRepository.GetUserOTP(otp)
	if err != nil {
		return res, errors.New("Failed to get user")
	}

	if req.Password != req.PasswordConfirm {
		return res, errors.New("Password does not matches")
	}

	// Reset OTP and OTP Request
	user.OTP = 0
	user.OTPReq = false

	// Update Password
	passwordHash, err := helpers.HashPassword(req.Password)
	if err != nil {
		return res, errors.New("Failed to hash password")
	}
	user.Password = string(passwordHash)

	_, err = u.userRepository.UpdateUser(user)
	if err != nil {
		return res, errors.New("Failed to update user")
	}

	res = dtos.ForgotPasswordResponse{
		Email:   user.Email,
		Message: "Password has been reset successfully",
	}

	return res, nil
}

// ChangePassword godoc
// @Summary      Change Password User
// @Description  Change Password User
// @Tags         User - Account
// @Accept       json
// @Produce      json
// @Param        request body dtos.ChangePasswordUserRequest true "Payload Body [RAW]"
// @Success      200 {object} dtos.ChangePasswordUserOKResponse
// @Failure      400 {object} dtos.BadRequestResponse
// @Failure      401 {object} dtos.UnauthorizedResponse
// @Failure      403 {object} dtos.ForbiddenResponse
// @Failure      404 {object} dtos.NotFoundResponse
// @Failure      500 {object} dtos.InternalServerErrorResponse
// @Router       /user/change-password [post]
// @Security BearerAuth
func (u *userUsecase) ChangePassword(id uint, req dtos.ChangePasswordUserRequest) (res helpers.ResponseMessage, err error) {
	user, err := u.userRepository.GetUserById(id)
	if err != nil {
		return res, errors.New("Failed to get user")
	}

	err = helpers.ComparePassword(req.OldPassword, user.Password)
	if err != nil {
		return res, errors.New("Wrong Password")
	}

	if req.Password != req.PasswordConfirm {
		return res, errors.New("Password does not match")
	}

	passwordHash, err := helpers.HashPassword(req.Password)
	if err != nil {
		return res, errors.New("Failed to hash password")
	}

	user.Password = string(passwordHash)

	user, err = u.userRepository.UpdateUser(user)
	if err != nil {
		return res, errors.New("Failed to update user")
	}

	res = helpers.NewResponseMessage(
		http.StatusOK,
		"Password has been changed successfully",
	)

	return res, nil
}

// GetAllUsers godoc
// @Summary      Get all users
// @Description  Get all users
// @Tags         User - Account
// @Accept       json
// @Produce      json
// @Success      200 {object} dtos.GetAllUserResponse
// @Failure      400 {object} dtos.BadRequestResponse
// @Failure      401 {object} dtos.UnauthorizedResponse
// @Failure      403 {object} dtos.ForbiddenResponse
// @Failure      404 {object} dtos.NotFoundResponse
// @Failure      500 {object} dtos.InternalServerErrorResponse
// @Router       /user [get]
// @Security BearerAuth
func (u *userUsecase) GetUsers() ([]dtos.UserDetailResponse, error) {
	users, err := u.userRepository.GetUsers()
	if err != nil {
		return nil, err
	}

	var userResponse []dtos.UserDetailResponse
	for _, user := range users {
		userResponse = append(userResponse, dtos.UserDetailResponse{
			ID:        user.ID,
			Username:  user.Username,
			Nama:      user.FullName,
			Email:     user.Email,
			Role:      user.Role,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		})
	}

	return userResponse, nil
}

// GetUserByID godoc
// @Summary      Get user by ID
// @Description  Get user by ID
// @Tags         User - Account
// @Accept       json
// @Produce      json
// @Param id path integer true "ID user"
// @Success      200 {object} dtos.UserStatusOKResponse
// @Failure      400 {object} dtos.BadRequestResponse
// @Failure      401 {object} dtos.UnauthorizedResponse
// @Failure      403 {object} dtos.ForbiddenResponse
// @Failure      404 {object} dtos.NotFoundResponse
// @Failure      500 {object} dtos.InternalServerErrorResponse
// @Router       /user/{id} [get]
// @Security BearerAuth
func (u *userUsecase) GetUserById(id uint) (res dtos.UserProfileResponse, err error) {
	user, err := u.userRepository.GetUserById(id)
	if err != nil {
		return res, errors.New("User not found")
	}

	res = dtos.UserProfileResponse{
		ID:       user.ID,
		Username: user.Username,
		Nama:     user.FullName,
		Email:    user.Email,
		Role:     user.Role,
	}

	return res, nil
}

// UserRegister godoc
// @Summary      Register User
// @Description  Register an account
// @Tags         User - Auth
// @Accept       json
// @Produce      json
// @Param        request body dtos.RegisterUserRequest true "Payload Body [RAW]"
// @Success      201 {object} dtos.UserCreeatedResponse
// @Failure      400 {object} dtos.BadRequestResponse
// @Failure      401 {object} dtos.UnauthorizedResponse
// @Failure      403 {object} dtos.ForbiddenResponse
// @Failure      404 {object} dtos.NotFoundResponse
// @Failure      500 {object} dtos.InternalServerErrorResponse
// @Router       /register [post]
func (u *userUsecase) CreateUser(req *dtos.RegisterUserRequest) (dtos.UserDetailResponse, error) {
	var res dtos.UserDetailResponse

	err := helpers.ValidateUsername(req.Username)
	if err != nil {
		return res, echo.NewHTTPError(400, err)
	}

	username, _ := u.userRepository.GetUserByUsername(req.Username)
	if username.ID > 0 {
		return res, errors.New("Username already in use")
	}

	req.Email = strings.ToLower(req.Email)

	// Check apakah email sudah terdaftar atau belum
	user, _ := u.userRepository.GetUserByEmail(req.Email)
	if user.ID > 0 {
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

	CreateUser := models.User{
		FullName:         req.Nama,
		Username:         req.Username,
		Email:            req.Email,
		Password:         passwordHash,
		VerificationCode: verification_code,
	}
	users, err := u.userRepository.CreateUser(CreateUser)

	if err != nil {
		return res, err
	}

	var firstName = CreateUser.Username

	if strings.Contains(firstName, " ") {
		firstName = strings.Split(firstName, " ")[1]
	}

	// 👇 Send Email
	emailData := utils.EmailData{
		URL:       config.ClientOrigin + "/verifyemail/" + url.PathEscape(verification_code),
		FirstName: firstName,
		Subject:   "Your account verification code",
	}

	utils.SendEmailUser(&CreateUser, &emailData)

	resp := dtos.UserDetailResponse{
		ID:        users.ID,
		Username:  users.Username,
		Nama:      users.FullName,
		Email:     users.Email,
		Role:      users.Role,
		CreatedAt: users.CreatedAt,
		UpdatedAt: users.UpdatedAt,
	}

	return resp, nil

}

// UserUpdate godoc
// @Summary      Update Information
// @Description  User update an information
// @Tags         User - Account
// @Accept       json
// @Produce      json
// @Param        request body dtos.UpdateUserRequest true "Payload Body [RAW]"
// @Success      200 {object} dtos.UserStatusOKResponse
// @Failure      400 {object} dtos.BadRequestResponse
// @Failure      401 {object} dtos.UnauthorizedResponse
// @Failure      403 {object} dtos.ForbiddenResponse
// @Failure      404 {object} dtos.NotFoundResponse
// @Failure      500 {object} dtos.InternalServerErrorResponse
// @Router       /user/{id} [put]
// @Security BearerAuth
func (u *userUsecase) UpdateUser(id uint, req dtos.UpdateUserRequest) (res dtos.UpdateUserResponse, err error) {
	var (
		users models.User
	)

	users, err = u.userRepository.GetUserById(id)
	if err != nil {
		return res, err
	}

	users.FullName = req.Nama
	users.Email = req.Email
	users.Password = req.Password
	users.Role = req.Role
	users.Username = req.Username

	// Check Role and save role information from JWT Cookie
	user, err := u.userRepository.ReadToken(id)
	if err != nil {
		return res, errors.New("Failed to get user")
	}

	users.Role = user.Role
	users.ID = uint(id)

	passwordHash, err := helpers.HashPassword(users.Password)
	if err != nil {
		return res, errors.New("Failed to hash password")
	}

	users.Password = string(passwordHash)

	users, err = u.userRepository.UpdateUser(users)
	if err != nil {
		return res, err
	}

	res.Username = users.Username
	res.Nama = users.FullName
	res.Email = users.Email
	res.Password = users.Password
	res.Role = users.Role

	return res, nil

}

// DeleteUser godoc
// @Summary      Delete an User
// @Description  Delete an User
// @Tags         User - Account
// @Accept       json
// @Produce      json
// @Param        request body dtos.DeleteUserRequest true "Payload Body [RAW]"
// @Success      200 {object} dtos.StatusOKDeletedResponse
// @Failure      400 {object} dtos.BadRequestResponse
// @Failure      401 {object} dtos.UnauthorizedResponse
// @Failure      403 {object} dtos.ForbiddenResponse
// @Failure      404 {object} dtos.NotFoundResponse
// @Failure      500 {object} dtos.InternalServerErrorResponse
// @Router       /user/{id} [delete]
// @Security BearerAuth
func (u *userUsecase) DeleteUser(id uint, req dtos.DeleteUserRequest) (res helpers.ResponseMessage, err error) {
	user, err := u.userRepository.ReadToken(id)

	if err != nil {
		return res, errors.New("Failed to get user")

	}

	err = helpers.ComparePassword(req.Password, user.Password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return res, errors.New("Password is incorrect")
	}

	err = u.userRepository.DeleteUser(user)

	if err != nil {
		return res, errors.New("Failed to delete user")
	}

	return res, nil
}
