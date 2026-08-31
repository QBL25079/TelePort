package marzban

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/QBL25079/TelePort/internal/domain"
)

type Client struct {
	baseURL    string
	token      string
	httpClient *http.Client
}

func NewClient(baseURL, token string, httpClient *http.Client) *Client {
	return &Client{
		baseURL:    baseURL,
		token:      token,
		httpClient: httpClient,
	}
}

func (c *Client) CreateOrUpdate(ctx context.Context, u domain.GateWay) (string, error) {
	userName := u.UserName
	if userName == "" {
		userName = fmt.Sprintf("tg_%d", u.TelegramID)
	}

	body := map[string]any{
		"username": userName,
		"expires":  u.ExpiresAt.Unix(),
		"status":   "active",
	}

	raw, _ := json.Marshal(body)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/api/user", bytes.NewReader(raw))
	if err != nil {
		return "", fmt.Errorf("cannot create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")

	res, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("cannot get response: %w", err)
	}
	defer res.Body.Close()

	b, _ := io.ReadAll(res.Body)
	if res.StatusCode >= 300 {
		return "", fmt.Errorf("marzban create user: %s: %s", res.Status, string(b))
	}

	var parsed struct {
		SubscriptionURL string `json:"subscription_url"`
	}
	_ = json.Unmarshal(b, &parsed)
	if parsed.SubscriptionURL != "" {
		return parsed.SubscriptionURL, nil
	}

	return fmt.Sprintf("%s/sub/%s", c.baseURL, userName), nil
}

func (c *Client) Disable(ctx context.Context, userName string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPut,
		fmt.Sprintf("%s/api/user/%s", c.baseURL, userName),
		bytes.NewReader([]byte(`{"status":"disabled"}`)))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")

	res, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode >= 300 {
		b, _ := io.ReadAll(res.Body)
		return fmt.Errorf("marzban disable: %s: %s", res.Status, string(b))
	}

	return nil
}
