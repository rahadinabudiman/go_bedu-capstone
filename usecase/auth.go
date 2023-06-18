package usecase

import (
	"errors"
	"go_bedu/dtos"
	"go_bedu/helpers"
	"go_bedu/initializers"
	"go_bedu/repositories"
	"go_bedu/utils"
	"net/url"
	"strconv"
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

// ForgotPassword godoc
// @Summary      Forgot Password Request OTP
// @Description  Forgot Password an Account
// @Tags         Utils - Authentikasi
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
func (u *authUsecase) ForgotPassword(req dtos.ForgotPasswordRequest) (res dtos.ForgotPasswordResponse, err error) {
	admin, err := u.adminRepository.GetAdminByEmail(req.Email)
	if err != nil {
		user, err := u.userRepository.GetUserByEmail(req.Email)
		if err != nil {
			return res, errors.New("Email not registered")
		}

		// Mengenerate OTP
		otp := helpers.GenerateRandomOTP(6)
		NewOTP, err := strconv.Atoi(otp)
		if err != nil {
			return res, errors.New("Failed to Generate OTP")
		}
		user.OTP = NewOTP
		user.OTPReq = true

		u.userRepository.UpdateUser(user)

		// 👇 Kirim Email
		config, _ := initializers.LoadConfig(".")
		emailData := utils.EmailData{
			URL:       config.ClientOrigin + "/change-password/" + url.PathEscape(otp),
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
		return res, errors.New("Failed to generate OTP")
	}
	admin.OTP = NewOTP
	admin.OTPReq = true

	u.adminRepository.UpdateAdmin(admin)

	// 👇 Kirim Email
	config, _ := initializers.LoadConfig(".")
	emailData := utils.EmailData{
		URL:       config.ClientOrigin + "/admin/change-password/" + url.PathEscape(otp),
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
