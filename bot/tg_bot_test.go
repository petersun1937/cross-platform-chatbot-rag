package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/stretchr/testify/mock"
)

// Mock for tgbotapi.BotAPI to simulate interactions with Telegram API
type MockTelegramAPI struct {
	mock.Mock
}

func (m *MockTelegramAPI) Send(c tgbotapi.Chattable) (tgbotapi.Message, error) {
	args := m.Called(c)
	return args.Get(0).(tgbotapi.Message), args.Error(1)
}

func (m *MockTelegramAPI) SetWebhook(config tgbotapi.WebhookConfig) (tgbotapi.APIResponse, error) {
	args := m.Called(config)
	return args.Get(0).(tgbotapi.APIResponse), args.Error(1)
}

func (m *MockTelegramAPI) GetUpdatesChan(config tgbotapi.UpdateConfig) (tgbotapi.UpdatesChannel, error) {
	args := m.Called(config)
	return args.Get(0).(tgbotapi.UpdatesChannel), args.Error(1)
}

// TestNewTGBot tests the creation of a new Telegram bot instance
/*func TestNewTGBot(t *testing.T) {
	// Mock the config and service
	mockConfig := &config.Config{
		TelegramBotToken: "mock_token",
	}
	mockService := new(service.Service)

	// Call the NewTGBot function
	tgBot, err := NewTGBot(mockConfig, mockService)

	// Ensure the tgBot instance was created successfully
	assert.NotNil(t, tgBot)
	assert.NotNil(t, tgBot.botApi)
}*/

// TestTGBot_SetWebhook tests the SetWebhook method of the tgBot struct
/*func TestTGBot_SetWebhook(t *testing.T) {
	// Mock the dependencies
	mockService := new(service.Service)
	mockTelegramAPI := new(MockTelegramAPI)

	// Initialize the tgBot instance
	tgBot := &tgBot{
		service:    mockService,
		lineClient: mockTelegramAPI,
	}

	// Mock the SetWebhook method
	mockTelegramAPI.On("SetWebhook", mock.Anything).Return(tgbotapi.APIResponse{Ok: true}, nil).Once()

	// Call the SetWebhook method
	err := tgBot.SetWebhook("https://example.com/webhook")

	// Assertions
	assert.NoError(t, err)
	mockTelegramAPI.AssertExpectations(t)
}*/

// TestTGBot_HandleTelegramUpdate tests the HandleTelegramUpdate method
/*func TestTGBot_HandleTelegramUpdate(t *testing.T) {
	// Mock the dependencies
	mockService := new(service.Service)
	mockTelegramAPI := new(MockTelegramAPI)

	// Initialize the tgBot instance
	tgBot := &TgBot{
		service:    mockService,
		lineClient: mockTelegramAPI,
	}

	// Create a mock update from Telegram
	mockUpdate := tgbotapi.Update{
		Message: &tgbotapi.Message{
			Chat: &tgbotapi.Chat{
				ID: 12345,
			},
			Text: "Hello",
		},
	}

	// Mock the Send method
	mockTelegramAPI.On("Send", mock.Anything).Return(tgbotapi.Message{}, nil).Once()

	// Call the HandleTelegramUpdate method
	tgBot.HandleTelegramUpdate(mockUpdate)

	// Assertions
	mockTelegramAPI.AssertExpectations(t)
}*/
