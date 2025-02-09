package handlers

import (
	config "crossplatform_chatbot/configs"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func (h *Handler) HandlerGetAIConfig(c *gin.Context) {
	const maxRetries = 10                        // Maximum number of retries
	const retryInterval = 500 * time.Millisecond // Wait 500ms between retries
	timeout := time.After(5 * time.Second)       // Timeout after 5 seconds

	var aiConfig gin.H
	updated := false

	for i := 0; i < maxRetries; i++ {
		select {
		case <-timeout:
			// If timeout occurs, return error response
			c.JSON(http.StatusRequestTimeout, gin.H{"error": "Config update timeout"})
			return
		default:
			cfg := config.GetConfig() // Fetch the latest config

			aiConfig = gin.H{
				"UseOpenAI":     cfg.BotConfig.UseOpenAI,
				"UseMistral":    cfg.BotConfig.UseMistral,
				"UseMETA":       cfg.BotConfig.UseMETA,
				"UseDialogflow": cfg.BotConfig.UseDialogflow,
			}

			// Check if at least one value is updated
			if aiConfig["UseOpenAI"].(bool) || aiConfig["UseMistral"].(bool) || aiConfig["UseMETA"].(bool) || aiConfig["UseDialogflow"].(bool) {
				updated = true
				break
			}

			time.Sleep(retryInterval) // Wait before retrying
		}
	}

	if updated {
		c.JSON(http.StatusOK, aiConfig)
	} else {
		c.JSON(http.StatusRequestTimeout, gin.H{"error": "Config update timeout"})
	}
}
