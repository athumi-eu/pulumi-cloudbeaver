package main

import (
	"context"
	"fmt"

	"github.com/pulumi/pulumi-go-provider/infer"
)

type EnableUser struct{}

type EnableUserArgs struct {
	UserName string `pulumi:"user_name"`
	Enabled  bool   `pulumi:"enabled"`
}

type EnableUserResponse struct {
	EnableUserArgs
}

func (EnableUser) Invoke(ctx context.Context, req infer.FunctionRequest[EnableUserArgs]) (infer.FunctionResponse[EnableUserResponse], error) {
	config := infer.GetConfig[CloudbeaverProviderConfig](ctx)

	body := map[string]interface{}{
		"operationName": "enableUser",
		"query":         "query enableUser($userId: ID!, $enabled: Boolean!) { enableUser(userId: $userId, enabled: $enabled) }",
		"variables": map[string]interface{}{
			"userId":  req.Input.UserName,
			"enabled": req.Input.Enabled,
		},
	}

	var responseBody interface{}
	_, err := sendPost(fmt.Sprintf("%s/api/gql", config.Endpoint), body, map[string]string{"cb-session-id": config.SessionId}, &responseBody)
	if err != nil {
		return infer.FunctionResponse[EnableUserResponse]{
			Output: EnableUserResponse{req.Input},
		}, err
	}
	return infer.FunctionResponse[EnableUserResponse]{
		Output: EnableUserResponse{req.Input},
	}, nil
}
