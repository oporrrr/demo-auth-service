package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"demo-auth-center/config"
)

type RoleClient struct {
	baseURL    string
	apiKey     string
	systemCode string
	http       *http.Client
}

func NewRoleClient() *RoleClient {
	return &RoleClient{
		baseURL:    config.Cfg.RoleServiceURL,
		apiKey:     config.Cfg.RoleServiceAPIKey,
		systemCode: config.Cfg.SystemCode,
		http:       &http.Client{Timeout: 5 * time.Second},
	}
}

func (c *RoleClient) do(method, path string, body any) ([]byte, error) {
	var r io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		r = bytes.NewBuffer(b)
	}
	req, err := http.NewRequest(method, c.baseURL+path, r)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", c.apiKey)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

// ── Permission Check ──────────────────────────────────────

type checkResponse struct {
	Code    string `json:"code"`
	Allowed bool   `json:"allowed"`
}

func (c *RoleClient) Check(accountID, resource, action string) bool {
	body, err := c.do("POST", "/internal/check", map[string]any{
		"accountId":  accountID,
		"systemCode": c.systemCode,
		"resource":   resource,
		"action":     action,
	})
	if err != nil {
		log.Printf("warn: role check error: %v", err)
		return false
	}
	var resp checkResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return false
	}
	return resp.Allowed
}

// ── Get User Permissions ──────────────────────────────────

func (c *RoleClient) GetPermissions(accountID string) ([]string, error) {
	body, err := c.do("GET", fmt.Sprintf("/internal/permissions?accountId=%s&system=%s", accountID, c.systemCode), nil)
	if err != nil {
		return nil, err
	}
	var resp struct {
		Code string   `json:"code"`
		Data []string `json:"data"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// ── Get User Menus ────────────────────────────────────────

func (c *RoleClient) GetMenus(accountID string) ([]any, error) {
	body, err := c.do("GET", fmt.Sprintf("/internal/menus?accountId=%s&system=%s", accountID, c.systemCode), nil)
	if err != nil {
		return nil, err
	}
	var resp struct {
		Code string `json:"code"`
		Data []any  `json:"data"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}
	return resp.Data, nil
}
