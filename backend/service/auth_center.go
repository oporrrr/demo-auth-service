package service

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"demo-auth-center/config"
	"demo-auth-center/model"
)

type AuthCenterService struct {
	baseURL     string
	clientToken string
	tokenExpiry time.Time
	mu          sync.RWMutex
	httpClient  *http.Client
}

var authCenterSvc *AuthCenterService
var once sync.Once

func GetAuthCenterService() *AuthCenterService {
	once.Do(func() {
		authCenterSvc = &AuthCenterService{
			baseURL:    config.Cfg.AuthCenterBaseURL,
			httpClient: &http.Client{Timeout: 30 * time.Second},
		}
	})
	return authCenterSvc
}

// ── Step 1: Get Client Token ──────────────────────────────
// Base64(clientId:clientSecret) → GET /api/v1/client/token

func (s *AuthCenterService) getClientToken() (string, error) {
	s.mu.RLock()
	if s.clientToken != "" && time.Now().Before(s.tokenExpiry) {
		token := s.clientToken
		s.mu.RUnlock()
		return token, nil
	}
	s.mu.RUnlock()

	// สร้าง Basic Auth header
	cred := fmt.Sprintf("%s:%s", config.Cfg.AuthClientID, config.Cfg.AuthClientSecret)
	encoded := base64.StdEncoding.EncodeToString([]byte(cred))

	req, err := http.NewRequest("GET", s.baseURL+"/api/v1/client/token", nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("accept", "application/json")
	req.Header.Set("Authorization", "Basic "+encoded)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to get client token: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var result model.ClientTokenResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("failed to parse client token response: %w", err)
	}

	if result.Code != "SUCCESS" || result.Data.AccessToken == "" {
		return "", fmt.Errorf("get client token failed: %s", string(body))
	}

	// Cache token 50 นาที (ปรับตาม expiry จริงของ server)
	s.mu.Lock()
	s.clientToken = result.Data.AccessToken
	s.tokenExpiry = time.Now().Add(50 * time.Minute)
	s.mu.Unlock()

	return result.Data.AccessToken, nil
}

// ── Helper: request พร้อม Bearer token ────────────────────

func (s *AuthCenterService) doRequest(method, path string, body interface{}, userToken string) ([]byte, int, error) {
	var reqBody io.Reader
	if body != nil {
		jsonBytes, err := json.Marshal(body)
		if err != nil {
			return nil, 0, err
		}
		reqBody = bytes.NewBuffer(jsonBytes)
	}

	req, err := http.NewRequest(method, s.baseURL+path, reqBody)
	if err != nil {
		return nil, 0, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("accept", "application/json")

	// ถ้ามี userToken ใช้อันนั้น ไม่งั้นใช้ clientToken
	if userToken != "" {
		req.Header.Set("Authorization", "Bearer "+userToken)
	} else {
		clientToken, err := s.getClientToken()
		if err != nil {
			return nil, 0, err
		}
		req.Header.Set("Authorization", "Bearer "+clientToken)
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	return respBody, resp.StatusCode, nil
}

// ── Auth APIs ─────────────────────────────────────────────

func (s *AuthCenterService) Register(req model.RegisterRequest) ([]byte, int, error) {
	return s.doRequest("POST", "/api/v1/auth/register", req, "")
}

func (s *AuthCenterService) Login(req model.LoginRequest) ([]byte, int, error) {
	return s.doRequest("POST", "/api/v1/auth/login", req, "")
}

func (s *AuthCenterService) LoginWithOTP(req model.LoginWithOTPRequest) ([]byte, int, error) {
	return s.doRequest("POST", "/api/v1/auth/login-with-otp", req, "")
}

func (s *AuthCenterService) RefreshToken(req model.RefreshTokenRequest) ([]byte, int, error) {
	return s.doRequest("POST", "/api/v1/auth/refresh-token", req, "")
}

func (s *AuthCenterService) Logout(req model.LogoutRequest, userToken string) ([]byte, int, error) {
	return s.doRequest("POST", "/api/v1/auth/logout", req, userToken)
}

func (s *AuthCenterService) UpdatePassword(req model.UpdatePasswordRequest, userToken string) ([]byte, int, error) {
	return s.doRequest("PUT", "/api/v1/auth/update-password", req, userToken)
}

// ── Account APIs ──────────────────────────────────────────

func (s *AuthCenterService) CheckExistValue(req model.CheckExistValueRequest) ([]byte, int, error) {
	return s.doRequest("POST", "/api/v1/account/check-exist-value", req, "")
}

func (s *AuthCenterService) GetAccountInformation(userToken string) ([]byte, int, error) {
	return s.doRequest("GET", "/api/v1/account/information", nil, userToken)
}

func (s *AuthCenterService) LinkCIS(req model.LinkCISRequest, userToken string) ([]byte, int, error) {
	return s.doRequest("POST", "/api/v1/account/link-cis", req, userToken)
}

func (s *AuthCenterService) UpdateUsername(req model.UpdateUsernameRequest, userToken string) ([]byte, int, error) {
	return s.doRequest("PUT", "/api/v1/account/update-username", req, userToken)
}

func (s *AuthCenterService) UpdateProfile(req model.UpdateProfileRequest, userToken string) ([]byte, int, error) {
	return s.doRequest("PUT", "/api/v1/account/update-profile", req, userToken)
}

func (s *AuthCenterService) GetCISNumber(req model.GetCISNumberRequest, userToken string) ([]byte, int, error) {
	return s.doRequest("POST", "/api/v1/account/get-cis-number", req, userToken)
}
