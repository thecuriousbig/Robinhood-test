package usersvc

import (
	"context"
	"log"
	"robinhood/internal/core/domains"
	"robinhood/internal/core/ports"
	"robinhood/internal/errmsg"
	"robinhood/pkg/auth"
	"robinhood/pkg/utils"
)

type userService struct {
	ur ports.UserRepository
}

func New(ur ports.UserRepository) ports.UserService {
	return &userService{ur: ur}
}

func (s *userService) Register(ctx context.Context, req *domains.RegisterRequest) error {
	user, err := s.ur.GetByUsername(ctx, req.Username)
	if err != nil {
		log.Printf("[userService::Register::GetByUsername] error => %+v", err)
		return errmsg.UserRegisterFailed
	}

	if user != nil {
		return errmsg.UserExisted
	}

	hashedPassword, err := utils.HashPassword(req.Password, utils.DefaultCost)
	if err != nil {
		log.Printf("[userService::Register::HashPassword] error => %+v", err)
		return errmsg.UserRegisterFailed
	}

	if _, err := s.ur.Create(ctx, &domains.CreateUserRequest{
		Username: req.Username,
		Password: hashedPassword,
		Email:    req.Email,
	}); err != nil {
		log.Printf("[userService::Register::Create] error => %+v", err)
		return errmsg.UserRegisterFailed
	}

	return nil
}

func (s *userService) Login(ctx context.Context, req *domains.LoginRequest) (*domains.LoginResponse, error) {
	user, err := s.ur.GetByUsername(ctx, req.Username)
	if err != nil {
		log.Printf("[userService::Login::GetByUsername] error => %+v", err)
		return nil, errmsg.UserLoginFailed
	}

	if user == nil {
		return nil, errmsg.UsernameOrPasswordIncorrect
	}

	// check password
	if !utils.CheckPasswordHash(req.Password, user.Password) {
		return nil, errmsg.UsernameOrPasswordIncorrect
	}

	// generate custom claims JWT token
	token, err := auth.GenerateToken(user.ID.Hex())
	if err != nil {
		log.Printf("[userService::Login::GenerateToken] error => %+v", err)
		return nil, errmsg.UserLoginFailed
	}

	return &domains.LoginResponse{
		Token: token,
	}, nil
}

func (s *userService) Update(ctx context.Context, req *domains.UpdateUserRequest) (*domains.User, error) {
	return s.ur.Update(ctx, req)
}
