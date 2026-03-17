package service

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"demo-auth-center/cache"
	"demo-auth-center/entity"
	"demo-auth-center/model"
	"demo-auth-center/repository"

	"gorm.io/gorm"
)

type UserService struct {
	authSvc   *AuthCenterService
	userRepo  *repository.UserRepository
	userCache *cache.UserCache
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{
		authSvc:   GetAuthCenterService(),
		userRepo:  repository.NewUserRepository(db),
		userCache: cache.NewUserCache(15 * time.Minute),
	}
}

// GetUser returns user from cache → DB (in that order)
func (s *UserService) GetUser(accountID string) (*entity.User, error) {
	if user, ok := s.userCache.Get(accountID); ok {
		return user, nil
	}
	user, err := s.userRepo.FindByAccountID(accountID)
	if err != nil {
		return nil, err
	}
	s.userCache.Set(user)
	return user, nil
}

// ── Auth ──────────────────────────────────────────────────

func (s *UserService) Register(req model.RegisterRequest) ([]byte, int, error) {
	req.IsIncludeAKID = false

	body, statusCode, err := s.authSvc.Register(req)
	if err != nil {
		return nil, 0, err
	}

	log.Printf("auth center register response [%d]: %s", statusCode, string(body))

	var resp model.RegisterResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return body, statusCode, nil
	}

	if resp.Code != "SUCCESS" {
		return body, statusCode, nil
	}

	user := &entity.User{}

	// call account/information to get full profile
	if resp.Data.AccessToken != "" {
		infoBody, _, err := s.authSvc.GetAccountInformation(resp.Data.AccessToken)
		if err == nil {
			log.Printf("account information response: %s", string(infoBody))
			var info model.AccountInfoResponse
			if err := json.Unmarshal(infoBody, &info); err == nil && info.Code == "SUCCESS" {
				user.AccountID = info.Data.ID
				user.FirstName = info.Data.FirstName
				user.LastName = info.Data.LastName
				user.Email = info.Data.Email
				user.PhoneNumber = info.Data.PhoneNumber
				user.CountryCode = info.Data.CountryCode
				user.PrefixName = info.Data.PrefixName
				user.Gender = info.Data.Gender
				user.DateOfBirth = info.Data.DateOfBirth
				user.AccountStatus = info.Data.AccountStatus
				user.CisNumber = info.Data.CisNumber
			}
		} else {
			log.Printf("warn: failed to get account information: %v", err)
		}
	}

	if err := s.userRepo.Upsert(user); err != nil {
		log.Printf("error: failed to save user %s to DB: %v", user.AccountID, err)
		return nil, 500, fmt.Errorf("failed to save user: %w", err)
	}
	s.userCache.Set(user)
	log.Printf("user %s saved to DB and cached", user.AccountID)

	out, _ := json.Marshal(map[string]any{"statusCode": 201, "code": "SUCCESS"})
	return out, 201, nil
}

func (s *UserService) Login(req model.LoginRequest) ([]byte, int, error) {
	body, statusCode, err := s.authSvc.Login(req)
	if err != nil {
		return nil, 0, err
	}

	if statusCode == 200 {
		var resp model.LoginResponse
		if err := json.Unmarshal(body, &resp); err != nil || resp.Code != "SUCCESS" || resp.Data.AccessToken == "" {
			return body, statusCode, nil
		}

		ttl := 15 * time.Minute
		if resp.Data.ExpiresIn > 0 {
			ttl = time.Duration(resp.Data.ExpiresIn) * time.Second
		}

		// get accountId from /information
		infoBody, _, err := s.authSvc.GetAccountInformation(resp.Data.AccessToken)
		if err != nil {
			log.Printf("warn: login could not get account information: %v", err)
			return body, statusCode, nil
		}
		var info model.AccountInfoResponse
		if err := json.Unmarshal(infoBody, &info); err != nil || info.Code != "SUCCESS" || info.Data.ID == "" {
			log.Printf("warn: login account information parse failed")
			return body, statusCode, nil
		}

		// cache with token TTL — no DB write needed
		// profile data lives in auth center, our DB only stores role/custom fields
		user := &entity.User{AccountID: info.Data.ID}
		s.userCache.SetWithTTL(user, ttl)
		log.Printf("login: user %s cached with TTL %s", info.Data.ID, ttl)
	}

	return body, statusCode, nil
}

// ── Passthrough (delegate to authSvc directly) ────────────

func (s *UserService) LoginWithOTP(req model.LoginWithOTPRequest) ([]byte, int, error) {
	return s.authSvc.LoginWithOTP(req)
}

func (s *UserService) RefreshToken(req model.RefreshTokenRequest) ([]byte, int, error) {
	return s.authSvc.RefreshToken(req)
}

func (s *UserService) Logout(req model.LogoutRequest, userToken string) ([]byte, int, error) {
	return s.authSvc.Logout(req, userToken)
}

func (s *UserService) UpdatePassword(req model.UpdatePasswordRequest, userToken string) ([]byte, int, error) {
	return s.authSvc.UpdatePassword(req, userToken)
}
