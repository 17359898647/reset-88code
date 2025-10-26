package reset

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

// GetBeijingTime returns current time in Beijing timezone
func GetBeijingTime() time.Time {
	// Beijing is UTC+8
	location, _ := time.LoadLocation("Asia/Shanghai")
	return time.Now().In(location)
}

const (
	// ResetCreditsAPI endpoint base URL
	ResetCreditsAPI = "https://www.88code.org/admin-api/cc-admin/system/subscription/my/reset-credits"
	// RequestDelay between reset requests (seconds)
	RequestDelay = 1
)

// ResetCredits resets credits for a single subscription
func (c *HTTPClient) ResetCredits(subscriptionID int) error {
	url := fmt.Sprintf("%s/%d", ResetCreditsAPI, subscriptionID)
	log.Printf("发送重置请求到: %s", url)

	// Send POST request with "null" body
	resp, err := c.DoRequest("POST", url, strings.NewReader("null"))
	if err != nil {
		return fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	log.Printf("响应状态码: %d", resp.StatusCode)

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("读取响应失败: %w", err)
	}

	// Log response for debugging
	log.Printf("响应内容: %s", string(body))

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("重置失败，状态码: %d, 响应: %s", resp.StatusCode, string(body))
	}

	return nil
}

// ResetAllCredits resets credits for all active subscriptions
func (c *HTTPClient) ResetAllCredits(subscriptions []ActiveSubscription) Stats {
	stats := Stats{
		Total: len(subscriptions),
	}

	// Get current Beijing time
	beijingTime := GetBeijingTime()
	log.Printf("当前北京时间: %s", beijingTime.Format("2006-01-02 15:04:05"))
	log.Printf("开始批量重置额度，共 %d 个订阅", stats.Total)

	for i, sub := range subscriptions {
		log.Printf("正在重置第 %d/%d 个订阅 (ID: %d, 当前额度: %.2f)...",
			i+1, stats.Total, sub.ID, sub.CurrentCredits)

		if err := c.ResetCredits(sub.ID); err != nil {
			log.Printf("❌ 订阅 ID %d 重置失败: %v", sub.ID, err)
			stats.Failed++
		} else {
			log.Printf("✅ 订阅 ID %d 重置成功", sub.ID)
			stats.Success++
		}

		// Add delay between requests (except for the last one)
		if i < len(subscriptions)-1 {
			time.Sleep(RequestDelay * time.Second)
		}
	}

	return stats
}

// ResetAt18 handles reset logic at 18:30 Beijing time
// Rule: Only reset if resetTimes >= 1
func (c *HTTPClient) ResetAt18(subscriptions []ActiveSubscription) Stats {
	stats := Stats{
		Total: len(subscriptions),
	}

	log.Printf("执行 18:30 重置策略，共 %d 个订阅", stats.Total)

	// Step 1: Classify subscriptions by resetTimes
	var resetTimes0 []ActiveSubscription
	var resetTimes1 []ActiveSubscription
	var resetTimes2 []ActiveSubscription

	for _, sub := range subscriptions {
		switch sub.ResetTimes {
		case 0:
			resetTimes0 = append(resetTimes0, sub)
		case 1:
			resetTimes1 = append(resetTimes1, sub)
		case 2:
			resetTimes2 = append(resetTimes2, sub)
		default:
			log.Printf("⚠️  订阅 ID %d 重置次数异常: %d", sub.ID, sub.ResetTimes)
		}
	}

	log.Printf("分类完成 - 重置次数 0: %d 个, 1: %d 个, 2: %d 个",
		len(resetTimes0), len(resetTimes1), len(resetTimes2))

	// Step 2: Filter resetTimes2 subscriptions that need reset (currentCredits <= 20% of creditLimit)
	var needReset []ActiveSubscription
	for _, sub := range resetTimes2 {
		usageRate := sub.CurrentCredits / sub.CreditLimit
		if usageRate <= 0.2 {
			needReset = append(needReset, sub)
			log.Printf("  订阅 ID %d 需要重置 (当前额度: %.2f, 总额度: %.0f, 使用率: %.1f%%)",
				sub.ID, sub.CurrentCredits, sub.CreditLimit, usageRate*100)
		}
	}

	if len(needReset) == 0 {
		log.Println("没有订阅需要重置")
		return stats
	}

	log.Printf("开始并发重置 %d 个订阅", len(needReset))

	// Step 3: Concurrent reset for filtered subscriptions
	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, sub := range needReset {
		wg.Add(1)
		go func(s ActiveSubscription) {
			defer wg.Done()

			log.Printf("正在重置订阅 ID %d (当前额度: %.2f, 总额度: %.0f)...",
				s.ID, s.CurrentCredits, s.CreditLimit)

			if err := c.ResetCredits(s.ID); err != nil {
				log.Printf("❌ 订阅 ID %d 重置失败: %v", s.ID, err)
				mu.Lock()
				stats.Failed++
				mu.Unlock()
			} else {
				log.Printf("✅ 订阅 ID %d 重置成功", s.ID)
				mu.Lock()
				stats.Success++
				mu.Unlock()
			}
		}(sub)
	}

	wg.Wait()

	return stats
}

// ResetAt23 handles reset logic at 23:45 Beijing time
// Rule: Reset all active subscriptions unconditionally
func (c *HTTPClient) ResetAt23(subscriptions []ActiveSubscription) Stats {
	stats := Stats{
		Total: len(subscriptions),
	}

	log.Printf("执行 23:45 重置策略，共 %d 个订阅", stats.Total)
	log.Printf("无条件重置所有活跃订阅")

	// Reset all subscriptions sequentially
	for i, sub := range subscriptions {
		log.Printf("正在重置第 %d/%d 个订阅 (ID: %d, 重置次数: %d, 当前额度: %.2f, 总额度: %.0f)...",
			i+1, stats.Total, sub.ID, sub.ResetTimes, sub.CurrentCredits, sub.CreditLimit)

		if err := c.ResetCredits(sub.ID); err != nil {
			log.Printf("❌ 订阅 ID %d 重置失败: %v", sub.ID, err)
			stats.Failed++
		} else {
			log.Printf("✅ 订阅 ID %d 重置成功", sub.ID)
			stats.Success++
		}

		// Add delay between requests (except for the last one)
		if i < len(subscriptions)-1 {
			time.Sleep(RequestDelay * time.Second)
		}
	}

	return stats
}

// ExecuteReset determines and executes the appropriate reset strategy based on current time
func (c *HTTPClient) ExecuteReset(subscriptions []ActiveSubscription) Stats {
	beijingTime := GetBeijingTime()
	hour := beijingTime.Hour()

	log.Printf("当前北京时间: %s", beijingTime.Format("2006-01-02 15:04:05"))

	if hour >= 18 && hour < 19 {
		// 18:00-18:59: Execute 18:30 strategy
		log.Println("🕐 当前时间段: 18:00-18:59，执行 18:30 重置策略")
		return c.ResetAt18(subscriptions)
	} else if hour >= 23 || hour < 1 {
		// 23:00-00:59: Execute 23:45 strategy
		log.Println("🕐 当前时间段: 23:00-00:59，执行 23:45 重置策略")
		return c.ResetAt23(subscriptions)
	} else {
		// Other times: Manual execution, reset all unconditionally
		log.Printf("⚠️  当前时间 %s 不在定时重置时间段内", beijingTime.Format("15:04:05"))
		log.Println("判定为手动执行，执行全量重置策略")
		return c.ResetAt23(subscriptions)
	}
}
