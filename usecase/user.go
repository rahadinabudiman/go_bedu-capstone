package usecase

import (
	"go_bedu/dtos"
	"go_bedu/helpers"
	"go_bedu/initializers"
	"go_bedu/middlewares"
	"go_bedu/models"
	"go_bedu/repositories"
	"go_bedu/utils"
	"net/url"
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
	ForgotPassword(req dtos.ForgotPasswordRequest) (res dtos.ForgotPasswordResponse, err error)
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

func (u *userUsecase) LoginUser(c echo.Context, req dtos.LoginRequest) (res dtos.LoginResponse, err error) {
	user, err := u.userRepository.GetUserByUsername(req.Username)
	if err != nil {
		echo.NewHTTPError(400, "Username not registered")
		return
	}

	if !user.Verified {
		return res, echo.NewHTTPError(400, "Please verify your email first")
	}

	err = helpers.ComparePassword(req.Password, user.Password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		echo.NewHTTPError(400, err.Error())
		return
	}

	token, err := middlewares.CreateToken(int(user.ID), user.Username, user.Email, user.Role)
	if err != nil {
		echo.NewHTTPError(400, "Failed to generate token")
		return
	}

	user.Token = token

	middlewares.CreateCookie(c, token)

	res = dtos.LoginResponse{
		Username: user.Username,
		Token:    user.Token,
	}

	return res, nil
}

func (u *userUsecase) LogoutUser(c echo.Context) (res dtos.LogoutUserResponse, err error) {
	err = middlewares.DeleteCookie(c)
	if err != nil {
		return res, echo.NewHTTPError(400, "Failed to logout")
	}

	return res, err
}

func (u *userUsecase) VerifyEmail(verificationCode any) (res dtos.VerifyEmailResponse, err error) {
	user, err := u.userRepository.GetUserByVerificationCode(verificationCode)
	if err != nil {
		return res, echo.NewHTTPError(400, "Failed to get admin")
	}

	if user.Verified {
		return res, echo.NewHTTPError(400, "Email already verified")
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

func (u *userUsecase) GetUserById(id uint) (res dtos.UserProfileResponse, err error) {
	user, err := u.userRepository.GetUserById(id)
	if err != nil {
		echo.NewHTTPError(401, "User didn't exist")
		return
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

func (u *userUsecase) CreateUser(req *dtos.RegisterUserRequest) (dtos.UserDetailResponse, error) {
	var res dtos.UserDetailResponse

	err := helpers.ValidateUsername(req.Username)
	if err != nil {
		return res, echo.NewHTTPError(400, err)
	}

	username, _ := u.userRepository.GetUserByUsername(req.Username)
	if username.ID > 0 {
		return res, echo.NewHTTPError(400, "Username already in use")
	}

	req.Email = strings.ToLower(req.Email)

	// Check apakah email sudah terdaftar atau belum
	user, _ := u.userRepository.GetUserByEmail(req.Email)
	if user.ID > 0 {
		return res, echo.NewHTTPError(400, "Email already in use")
	}

	if req.Password != req.PasswordConfirm {
		return res, echo.NewHTTPError(400, "Password does not match")
	}

	passwordHash, err := helpers.HashPassword(req.Password)
	if err != nil {
		return res, err
	}

	config, err := initializers.LoadConfig(".")
	if err != nil {
		return res, echo.NewHTTPError(400, "Failed to load config")
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
		URL:       "http://" + config.ClientOrigin + "/verifyemail/" + url.PathEscape(verification_code),
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
		return res, echo.NewHTTPError(400, "Failed to get User")
	}

	users.Role = user.Role
	users.ID = uint(id)

	passwordHash, err := helpers.HashPassword(users.Password)
	if err != nil {
		return res, echo.NewHTTPError(400, "Failed to hash password")
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

func (u *userUsecase) DeleteUser(id uint, req dtos.DeleteUserRequest) (res helpers.ResponseMessage, err error) {
	user, err := u.userRepository.ReadToken(id)

	if err != nil {
		echo.NewHTTPError(400, "Failed to get User")
		return
	}

	err = helpers.ComparePassword(req.Password, user.Password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		echo.NewHTTPError(400, err.Error())
		return
	}

	err = u.userRepository.DeleteUser(user)

	if err != nil {
		return res, echo.NewHTTPError(500, "Failed to delete user")
	}

	return res, nil
}
