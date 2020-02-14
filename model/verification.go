package model

// Verification stores randomly generated tokens that can be redeemed
type Verification struct {
	Base
	Token  string `json:"token"`
	UserID int    `json:"user_id"`
}
