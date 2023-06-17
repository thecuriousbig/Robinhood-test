package userhdl

import (
	"net/http"
	"robinhood/internal/core/domains"
	"robinhood/internal/core/ports"
	"robinhood/internal/dto"
	"robinhood/pkg/auth"

	"github.com/asaskevich/govalidator"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	s ports.UserService
}

func New(s ports.UserService) *Handler {
	return &Handler{s: s}
}

// @Summary      Register
// @Tags         User
// @Accept       json
// @Produce      json
// @Router       /user/register [post]
// @Param username body string true "username"
// @Param password body string true "password"
// @Param email body string true "email"
// @Response 200 {object} dto.BaseResponse
// @Response 400 {object} dto.BaseErrorResponse
// @Response 500 {object} dto.BaseErrorResponse
func (h *Handler) Register(c echo.Context) error {
	ctx := c.Request().Context()
	var req dto.RegisterRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	// validate request
	if _, err := govalidator.ValidateStruct(req); err != nil {
		return err
	}

	// create user
	if err := h.s.Register(ctx, &domains.RegisterRequest{
		Username: req.Username,
		Password: req.Password,
		Email:    req.Email,
	}); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, dto.BaseResponse{
		Code: 0,
	})
}

// @Summary      Login
// @Tags         User
// @Accept       json
// @Produce      json
// @Router       /user/login [post]
// @Param username body string true "username"
// @Param password body string true "password"
// @Response 200 {object} dto.BaseResponseWithData[dto.LoginResponse]
// @Response 400 {object} dto.BaseErrorResponse
// @Response 500 {object} dto.BaseErrorResponse
func (h *Handler) Login(c echo.Context) error {
	ctx := c.Request().Context()
	var req dto.LoginRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	// validate request
	if _, err := govalidator.ValidateStruct(req); err != nil {
		return err
	}

	// login
	res, err := h.s.Login(ctx, &domains.LoginRequest{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, dto.BaseResponseWithData[dto.LoginResponse]{
		BaseResponse: dto.BaseResponse{
			Code: 0,
		},
		Data: dto.LoginResponse{
			Token: res.Token,
		},
	})
}

// @Summary      Update user
// @Tags         User
// @Accept       json
// @Produce      json
// @Router       /user [put]
// @Security     ApiKeyAuth
// @Param profileImage body string true "url of profile image"
// @Response 200 {object} dto.BaseResponseWithData[dto.User]
// @Response 400 {object} dto.BaseErrorResponse
// @Response 500 {object} dto.BaseErrorResponse
func (h *Handler) UpdateUser(c echo.Context) error {
	ctx := c.Request().Context()
	user := c.Get("user").(*jwt.Token)
	claims, ok := user.Claims.(*auth.JWTCustomClaims)
	if !ok {
		return echo.ErrUnauthorized
	}
	userId := claims.UserId

	var req dto.UpdateUserRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	// validate request
	if _, err := govalidator.ValidateStruct(req); err != nil {
		return err
	}

	// update user
	updatedUser, err := h.s.Update(ctx, &domains.UpdateUserRequest{
		UserId:       userId,
		ProfileImage: req.ProfileImage,
	})
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, dto.BaseResponseWithData[dto.User]{
		BaseResponse: dto.BaseResponse{
			Code: 0,
		},
		Data: dto.User{
			ID:           updatedUser.ID.Hex(),
			Username:     updatedUser.Username,
			Email:        updatedUser.Email,
			ProfileImage: updatedUser.ProfileImage,
		},
	})

}
