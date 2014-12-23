package api

import (
	"encoding/base64"
	"encoding/json"
)

type AuthConfig struct {
	Username      string `json:"username"`
	Password      string `json:"password"`
	Email         string `json:"email"`
	ServerAddress string `json:"serveraddress"`
	Auth          string `json:"auth"`
}

func (authConfig *AuthConfig) Encode() string {
	buf, err := json.Marshal(authConfig)
	if err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(buf)
}
