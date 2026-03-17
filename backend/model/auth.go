package model

// ── Client Token ──────────────────────────────────────────

type ClientTokenResponse struct {
	StatusCode int    `json:"statusCode"`
	Code       string `json:"code"`
	Data       struct {
		AccessToken string `json:"accessToken"`
	} `json:"data"`
}

// ── Auth ──────────────────────────────────────────────────

type RegisterRequest struct {
	FirstName     string `json:"firstName,omitempty"`
	LastName      string `json:"lastName,omitempty"`
	CountryCode   string `json:"countryCode,omitempty"`
	PhoneNumber   string `json:"phoneNumber,omitempty"`
	Email         string `json:"email,omitempty"`
	PrefixName    string `json:"prefixName,omitempty"`
	Gender        string `json:"gender,omitempty"`
	DateOfBirth   string `json:"dateOfBirth,omitempty"`
	Password      string `json:"password"`
	IsIncludeAKID bool   `json:"isIncludeAKID,omitempty"`
}

type LoginRequest struct {
	LoginType   string `json:"loginType"`
	CountryCode string `json:"countryCode,omitempty"`
	PhoneNumber string `json:"phoneNumber,omitempty"`
	Email       string `json:"email,omitempty"`
	Password    string `json:"password"`
}

type LoginWithOTPRequest struct {
	OTP         string `json:"otp"`
	PhoneNumber string `json:"phoneNumber"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken"`
}

type LogoutRequest struct {
	RefreshToken string `json:"refreshToken"`
}

type UpdatePasswordRequest struct {
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}

// ── Account ───────────────────────────────────────────────

type CheckExistValueRequest struct {
	Field string `json:"field"`
	Value string `json:"value"`
}

type LinkCISRequest struct {
	CisNumber string `json:"cisNumber"`
}

type UpdateUsernameRequest struct {
	Username string `json:"username"`
}

type UpdateProfileRequest struct {
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	PhoneNumber string `json:"phoneNumber"`
}

type GetCISNumberRequest struct {
	AccountID string `json:"accountId"`
}

// ── Auth Center Typed Responses ───────────────────────────

type RegisterResponse struct {
	StatusCode int    `json:"statusCode"`
	Code       string `json:"code"`
	Data       struct {
		AccountID   string `json:"accountId"`
		AccessToken string `json:"accessToken"`
	} `json:"data"`
	Message string `json:"message,omitempty"`
}

type AccountInfoResponse struct {
	StatusCode int    `json:"statusCode"`
	Code       string `json:"code"`
	Data       struct {
		ID                string  `json:"id"`
		FirstName         string  `json:"firstName"`
		LastName          string  `json:"lastName"`
		Email             string  `json:"email"`
		PhoneNumber       string  `json:"phoneNumber"`
		CountryCode       string  `json:"countryCode"`
		PrefixName        string  `json:"prefixName"`
		Gender            string  `json:"gender"`
		DateOfBirth       string  `json:"dateOfBirth"`
		AccountStatus     string  `json:"accountStatus"`
		Username          string  `json:"username"`
		CisNumber         *string `json:"cisNumber"`
		ProfilePictureURL *string `json:"profilePictureUrl"`
	} `json:"data"`
}

type LoginResponse struct {
	StatusCode int    `json:"statusCode"`
	Code       string `json:"code"`
	Data       struct {
		AccessToken      string `json:"accessToken"`
		RefreshToken     string `json:"refreshToken"`
		ExpiresIn        int    `json:"expiresIn"`        // seconds
		RefreshExpiresIn int    `json:"refreshExpiresIn"` // seconds
	} `json:"data"`
	Message string `json:"message,omitempty"`
}

// ── Generic Response ──────────────────────────────────────

type AuthResponse struct {
	StatusCode int         `json:"statusCode"`
	Code       string      `json:"code"`
	Data       interface{} `json:"data,omitempty"`
	Message    string      `json:"message,omitempty"`
}
