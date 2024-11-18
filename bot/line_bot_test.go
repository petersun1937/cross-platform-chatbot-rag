package bot

import (
	"net/http"

	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/stretchr/testify/mock"
)

type MockService struct {
	mock.Mock
}

type MockLineClient struct {
	linebot.Client // Embed the real linebot.Client
	mock.Mock      // Testify mock
}

// Define a MockConfig that mocks the config.Config struct
type MockConfig struct {
	mock.Mock
}

type MockDB struct {
	mock.Mock
}

// Mock the GetLineSecret method
func (m *MockConfig) GetLineSecret() string {
	args := m.Called()
	return args.String(0)
}

// Mock the GetLineToken method
func (m *MockConfig) GetLineToken() string {
	args := m.Called()
	return args.String(0)
}

// func (m *MockService) GetDB() service.Database {
// 	args := m.Called()
// 	return args.Get(0).(service.Database)
// }

func (m *MockLineClient) GetProfile(userID string) (*linebot.UserProfileResponse, error) {
	args := m.Called(userID)
	return args.Get(0).(*linebot.UserProfileResponse), args.Error(1)
}

func (m *MockLineClient) ReplyMessage(replyToken string, messages ...linebot.SendingMessage) *linebot.ReplyMessageCall {
	args := m.Called(replyToken, messages)
	return args.Get(0).(*linebot.ReplyMessageCall)
}

type MockReplyMessageCall struct {
	mock.Mock
}

func (m *MockReplyMessageCall) Do() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockLineClient) ParseRequest(req *http.Request) ([]*linebot.Event, error) {
	args := m.Called(req)
	return args.Get(0).([]*linebot.Event), args.Error(1)
}

// Mock for GORM methods like Where, First, and Save
type MockGormDB struct {
	mock.Mock
}

/*func TestNewLineBot(t *testing.T) {
	// Reset the config singleton before the test
	config.ResetConfig()

	// Set environment variables for the test
	os.Setenv("DATABASE_URL", "mock_db_string")
	os.Setenv("TELEGRAM_BOT_TOKEN", "mock_telegram_token")
	os.Setenv("LINE_CHANNEL_SECRET", "mock_secret")
	os.Setenv("LINE_CHANNEL_TOKEN", "mock_token")
	os.Setenv("TELEGRAM_API_URL", "https://api.telegram.org/bot")
	os.Setenv("SERVER_HOST", "localhost")
	os.Setenv("APP_PORT", "8080")
	os.Setenv("SERVER_TIMEOUT", "30s")
	os.Setenv("SERVER_MAX_CONN", "100")

	// Initialize the config
	mockConfig := config.GetConfig()

	// Initialize a mock service
	mockService := new(service.Service)

	// Call the NewLineBot constructor
	lineBot, err := NewLineBot(mockConfig, mockService)

	assert.NoError(t, err)
	assert.NotNil(t, lineBot)

	// Check that the secret and token are set correctly from the environment variables
	assert.Equal(t, "mock_secret", lineBot.secret)
	assert.Equal(t, "mock_token", lineBot.token)

	// test other aspects?
	assert.NotNil(t, lineBot.lineClient)
}*/

/*
	func TestLineBot_Run(t *testing.T) {
		lineBot := &lineBot{
			secret: "mock_secret",
			token:  "mock_token",
		}

		err := lineBot.Run()
		assert.NoError(t, err)
		assert.NotNil(t, lineBot.lineClient)
	}
*/
/*
func TestLineBot_GetUserProfile(t *testing.T) {
	// Set environment variables for the test
	os.Setenv("LINE_CHANNEL_SECRET", "mock_secret")
	os.Setenv("LINE_CHANNEL_TOKEN", "mock_token")

	// Initialize the real linebot.Client
	realLineClient, err := linebot.New("mock_secret", "mock_token")
	assert.NoError(t, err)

	// Initialize service
	Service := new(service.Service)

	// Create the lineBot instance
	lineBot := &lineBot{
		BaseBot:    &BaseBot{},     // Initialize BaseBot to avoid nil issues
		lineClient: realLineClient, // Use the real linebot.Client
		service:    Service,
	}

	// gock to mock HTTP requests
	defer gock.Off() // Ensure gock is disabled after the test
	gock.New("https://api.line.me").
		Get("/v2/bot/profile/mock_user_id").
		Reply(200).
		JSON(map[string]string{
			"userId":      "mock_user_id",
			"displayName": "Mock User",
		})

	// Call the method to be tested
	userProfile, err := lineBot.getUserProfile("mock_user_id")

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, "mock_user_id", userProfile.UserID)
	assert.Equal(t, "Mock User", userProfile.DisplayName)

	// ensure no unexpected requests are made
	assert.True(t, gock.IsDone())
}*/

/*func TestLineBot_validateAndGenerateToken_UserNotFound(t *testing.T) {
	// Create mock dependencies
	mockDB := new(service.MockDB)
	mockGormDB := new(MockGormDB) // Mock for GORM-like behavior
	// Initialize the real linebot.Client
	realLineClient, err := linebot.New("mock_secret", "mock_token")
	assert.NoError(t, err)

	// Mock the GetDB method to return the mock GORM DB
	mockDB.On("GetDB").Return(mockGormDB).Once()

	// Mock GORM's Where, First, and Save methods
	mockGormDB.On("Where", mock.Anything, mock.Anything).Return(mockGormDB).Once()
	mockGormDB.On("First", mock.Anything).Return(errors.New("record not found")).Once()
	mockGormDB.On("Save", mock.Anything).Return(nil).Once()

	// Mock the Line Client's GetProfile method
	// mockLineClient.On("GetProfile", "mock_user_id").Return(&linebot.UserProfileResponse{
	// 	UserID:      "mock_user_id",
	// 	DisplayName: "Mock User",
	// }, nil).Once()

	// Create the service and lineBot instances
	serviceInstance := service.NewService(mockDB) // Use the mocked database

	lineBot := &lineBot{
		BaseBot:    &BaseBot{},
		service:    serviceInstance,
		lineClient: realLineClient, // Use the mocked Line client
	}

	// Call the method under test
	userExists, err := lineBot.validateAndGenerateToken(&linebot.UserProfileResponse{
		UserID:      "mock_user_id",
		DisplayName: "Mock User",
	}, &linebot.Event{}, "mock_user_id")

	userProfile, err := lineBot.getUserProfile("mock_user_id")

	// Assertions
	assert.NoError(t, err)
	assert.False(t, userExists)

	// Verify the expectations
	assert.NoError(t, err)
	assert.Equal(t, "mock_user_id", userProfile.UserID)
	assert.Equal(t, "Mock User", userProfile.DisplayName)
}*/

/*
func TestLineBot_HandleLineMessage(t *testing.T) {
	// Create mock dependencies
	mockService := new(MockService)
	mockLineClient := new(MockLineClient)
	mockReplyMessageCall := new(MockReplyMessageCall)

	lineBot := &lineBot{
		service:    mockService,
		lineClient: mockLineClient, // Pass the mock Line client
	}

	// Mock the GetProfile method to return a mock user profile
	mockUserProfile := &linebot.UserProfileResponse{
		UserID:      "mock_user_id",
		DisplayName: "Mock User",
	}
	mockLineClient.On("GetProfile", "mock_user_id").Return(mockUserProfile, nil).Once()

	// Mock the ReplyMessage method and chaining with Do()
	mockLineClient.On("ReplyMessage", "mock_reply_token", mock.Anything).Return(mockReplyMessageCall).Once()
	mockReplyMessageCall.On("Do").Return(nil).Once()

	// Mock user profile and event
	mockEvent := &linebot.Event{
		Source: &linebot.EventSource{
			UserID: "mock_user_id",
		},
		ReplyToken: "mock_reply_token",
	}
	mockTextMessage := &linebot.TextMessage{Text: "Hello"}

	// Call the HandleLineMessage method
	lineBot.HandleLineMessage(mockEvent, mockTextMessage)

	// Verify the expectations
	mockLineClient.AssertExpectations(t)
	mockReplyMessageCall.AssertExpectations(t)
}
*/
/*
func TestLineBot_ParseRequest(t *testing.T) {
	mockLineClient := new(MockLineClient)
	lineBot := &lineBot{lineClient: mockLineClient}

	mockRequest := &http.Request{}
	mockEvent := []*linebot.Event{}

	mockLineClient.On("ParseRequest", mockRequest).Return(mockEvent, nil)

	events, err := lineBot.ParseRequest(mockRequest)
	assert.NoError(t, err)
	assert.NotNil(t, events)
}*/
