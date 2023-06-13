package usecase

import (
	"go_bedu/dtos"
	"go_bedu/helpers"
	"go_bedu/initializers"
	"go_bedu/repositories"
	"go_bedu/utils"
	"net/url"
	"strconv"

	"github.com/labstack/echo/v4"
)

type AuthUsecase interface {
	ForgotPassword(req dtos.ForgotPasswordRequest) (res dtos.ForgotPasswordResponse, err error)
}

type authUsecase struct {
	adminRepository repositories.AdminRepository
	userRepository  repositories.UserRepository
}

func NewAuthUsecase(adminRepository repositories.AdminRepository, userRepository repositories.UserRepository) AuthUsecase {
	return &authUsecase{
		adminRepository: adminRepository,
		userRepository:  userRepository,
	}
}

func (u *authUsecase) ForgotPassword(req dtos.ForgotPasswordRequest) (res dtos.ForgotPasswordResponse, err error) {
	admin, err := u.adminRepository.GetAdminByEmail(req.Email)
	if err != nil {
		user, err := u.userRepository.GetUserByEmail(req.Email)
		if err != nil {
			return res, echo.NewHTTPError(400, "Email not registered")
		}

		// Mengenerate OTP
		otp := helpers.GenerateRandomOTP(6)
		NewOTP, err := strconv.Atoi(otp)
		if err != nil {
			return res, echo.NewHTTPError(400, "Failed to generate OTP")
		}
		user.OTP = NewOTP
		user.OTPReq = true

		u.userRepository.UpdateUser(user)

		// 👇 Kirim Email
		config, _ := initializers.LoadConfig(".")
		emailData := utils.EmailData{
			URL:       "http://" + config.ClientOrigin + "/change-password/" + url.PathEscape(otp),
			FirstName: user.Username,
			Subject:   "Your OTP to reset password",
		}

		utils.SendEmailUser(&user, &emailData)

		res = dtos.ForgotPasswordResponse{
			Email:   user.Email,
			Message: "OTP has been sent to your email",
		}

		return res, nil
	}

	// Mengenerate OTP
	otp := helpers.GenerateRandomOTP(6)
	NewOTP, err := strconv.Atoi(otp)
	if err != nil {
		return res, echo.NewHTTPError(400, "Failed to generate OTP")
	}
	admin.OTP = NewOTP
	admin.OTPReq = true

	u.adminRepository.UpdateAdmin(admin)

	// 👇 Kirim Email
	config, _ := initializers.LoadConfig(".")
	emailData := utils.EmailData{
		URL:       "http://" + config.ClientOrigin + "/admin/change-password/" + url.PathEscape(otp),
		FirstName: admin.Username,
		Subject:   "Your OTP to reset password",
	}

	utils.SendEmail(&admin, &emailData)

	res = dtos.ForgotPasswordResponse{
		Email:   admin.Email,
		Message: "OTP has been sent to your email",
	}

	return res, nil
}
