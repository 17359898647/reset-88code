package reset

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

const (
	// SubscriptionAPI endpoint
	SubscriptionAPI = "https://www.88code.org/admin-api/cc-admin/system/subscription/my"
)

// FetchSubscriptions fetches user subscription information
func (c *HTTPClient) FetchSubscriptions() (*SubscriptionResponse, error) {
	log.Printf("发送请求到: %s", SubscriptionAPI)

	resp, err := c.DoRequest("GET", SubscriptionAPI, nil)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	log.Printf("响应状态码: %d", resp.StatusCode)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("请求失败，状态码: %d", resp.StatusCode)
	}

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	// Parse JSON response
	var result SubscriptionResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析JSON失败: %w", err)
	}

	if result.Code != 0 || !result.OK {
		return nil, fmt.Errorf("API返回错误: %s", result.Msg)
	}

	return &result, nil
}

// ExtractActiveSubscriptions extracts active subscriptions with essential fields
func (r *SubscriptionResponse) ExtractActiveSubscriptions() []ActiveSubscription {
	var active []ActiveSubscription

	for _, sub := range r.Data {
		// Filter out inactive subscriptions
		if !sub.IsActive {
			continue
		}

		active = append(active, ActiveSubscription{
			ID:             sub.ID,
			ResetTimes:     sub.ResetTimes,
			CurrentCredits: sub.CurrentCredits,
			CreditLimit:    sub.SubscriptionPlan.CreditLimit,
		})
	}

	return active
}
