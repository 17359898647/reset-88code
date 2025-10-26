package main

import (
	"log"
	"os"

	"github.com/ggg/reset-88code/pkg/reset"
)

func main() {
	log.SetFlags(0) // Remove timestamp from log output

	// Load configuration from environment
	config, err := reset.LoadConfig()
	if err != nil {
		log.Printf("❌ %v", err)
		os.Exit(1)
	}

	log.Println("开始重置额度...")
	log.Printf("Token: %s", reset.MaskToken(config.Token))

	// Create HTTP client
	client := reset.NewHTTPClient(config.Token)

	// Fetch subscription data
	result, err := client.FetchSubscriptions()
	if err != nil {
		log.Printf("❌ 获取订阅信息失败: %v", err)
		os.Exit(1)
	}

	log.Printf("✅ 获取订阅信息成功，共 %d 个订阅", len(result.Data))

	// Extract and filter active subscriptions
	activeSubscriptions := result.ExtractActiveSubscriptions()
	log.Printf("活跃订阅数: %d / %d", len(activeSubscriptions), len(result.Data))

	if len(activeSubscriptions) == 0 {
		log.Println("⚠️  没有活跃订阅")
		return
	}

	// Print active subscriptions with essential fields
	for _, sub := range activeSubscriptions {
		log.Printf("  [活跃] ID: %d, 重置次数: %d, 当前额度: %.2f, 总额度: %.0f",
			sub.ID, sub.ResetTimes, sub.CurrentCredits, sub.CreditLimit)
	}

	// Execute reset based on current time
	stats := client.ExecuteReset(activeSubscriptions)

	// Print final statistics
	log.Println("\n📊 重置完成统计：")
	log.Printf("   总订阅数: %d", stats.Total)
	log.Printf("   成功数量: %d", stats.Success)
	log.Printf("   失败数量: %d", stats.Failed)

	if stats.Failed == 0 && stats.Success > 0 {
		log.Println("🎉 所有订阅额度重置成功！")
	} else if stats.Success == 0 && stats.Failed == 0 {
		log.Println("ℹ️  无订阅需要重置")
	} else {
		log.Printf("⚠️  有 %d 个订阅重置失败，请检查日志", stats.Failed)
		os.Exit(1)
	}
}
