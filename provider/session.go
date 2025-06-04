package main

import (
	"fmt"
)

func NewCloudbeaverSession(apiKey string, endpoint string) (string, error) {
	var sessionId string
	var responseBody interface{}

	body := map[string]interface{}{
		"operationName": "authLogin",
		"query":         "query authLogin($provider: ID!, $credentials: Object) { authLogin(provider: $provider, credentials: $credentials) { userTokens { userId } } }",
		"variables": map[string]interface{}{
			"credentials": map[string]string{
				"token": apiKey,
			},
			"provider": "token",
		},
	}

	response, err := sendPost(fmt.Sprintf("%s/api/gql", endpoint), body, map[string]string{}, &responseBody)
	if err != nil {
		return sessionId, err
	}

	for _, cookie := range response.Cookies() {
		if cookie.Name == "cb-session-id" {
			sessionId = cookie.Value
			break
		}
	}

	return sessionId, nil
}

// TODO: Call this somewhere, not really necessary when using API keys but needed when authenticating with username/password
func CloseCloudbeaverSession(sessionId string, endpoint string) error {
	var responseBody interface{}

	body := map[string]interface{}{
		"operationName": "authLogout",
		"query":         "query authLogout($provider: ID!) { authLogout(provider: $provider) }",
		"variables": map[string]interface{}{
			"provider": "token",
		},
	}

	_, err := sendPost(fmt.Sprintf("%s/api/gql", endpoint), body, map[string]string{"cb-session-id": sessionId}, &responseBody)
	if err != nil {
		return err
	}

	return nil
}
