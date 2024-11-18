package models

// struct of the incoming request for general bot.
type GeneralRequest struct {
	Message   string `json:"message"`
	SessionID string `json:"sessionID"`
}
