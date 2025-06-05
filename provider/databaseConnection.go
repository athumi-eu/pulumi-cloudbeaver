package main

import (
	"context"
	"fmt"

	"github.com/pulumi/pulumi-go-provider/infer"
)

type CreateDatabaseConnectionResponseDataConnection struct {
	// This is the only property we are currently intrested in
	ID string `json:"id"`
}

type CreateDatabaseConnectionResponseData struct {
	Connection CreateDatabaseConnectionResponseDataConnection `json:"connection"`
}

type CreateDatabaseConnectionResponse struct {
	Data CreateDatabaseConnectionResponseData `json:"data"`
}

type DatabaseConnection struct {
}

type DatabaseConnectionArgs struct {
	Name     string `pulumi:"name"`
	Endpoint string `pulumi:"endpoint"`
	Database string `pulumi:"database"`

	ProjectId string `pulumi:"project_id"`
}
type DatabaseConnectionState struct {
	DatabaseConnectionArgs
}

func (w *DatabaseConnection) Create(ctx context.Context, req infer.CreateRequest[DatabaseConnectionArgs]) (infer.CreateResponse[DatabaseConnectionState], error) {
	state := DatabaseConnectionState{DatabaseConnectionArgs: req.Inputs}

	config := infer.GetConfig[CloudbeaverProviderConfig](ctx)

	if req.DryRun {
		return infer.CreateResponse[DatabaseConnectionState]{
			ID:     "",
			Output: state,
		}, nil
	}

	body := map[string]interface{}{
		"operationName": "createConnection",
		"query":         "mutation createConnection($projectId: ID!, $config: ConnectionConfig!) { createConnection(projectId: $projectId, config: $config) { id } }",
		"variables": map[string]interface{}{
			"projectId": req.Inputs.ProjectId,
			"config": map[string]interface{}{
				"driverId":          "postgresql:postgres-jdbc",
				"authModelId":       "azure_ad_postgresql",
				"name":              req.Inputs.Name,
				"configurationType": "MANUAL",
				"host":              req.Inputs.Endpoint,
				"databaseName":      req.Inputs.Database,
				"port":              "5432",
				"sharedCredentials": true,
			},
		},
	}
	var responseBody CreateDatabaseConnectionResponse
	_, err := sendPost(fmt.Sprintf("%s/api/gql", config.Endpoint), body, map[string]string{"cb-session-id": config.SessionId}, &responseBody)
	if err != nil {
		return infer.CreateResponse[DatabaseConnectionState]{}, err
	}

	return infer.CreateResponse[DatabaseConnectionState]{
		ID:     responseBody.Data.Connection.ID,
		Output: state,
	}, nil
}

// func (w *DatabaseConnection) Read(ctx context.Context, req infer.ReadRequest[DatabaseConnectionArgs, DatabaseConnectionState]) (resp infer.ReadResponse[DatabaseConnectionArgs, DatabaseConnectionState], err error) {
// 	// Do stuff here
// 	return infer.ReadResponse[DatabaseConnectionArgs, DatabaseConnectionState]{
// 		ID:     req.ID,
// 		Inputs: req.Inputs,
// 		State:  req.State,
// 	}, nil
// }

// func (w *DatabaseConnection) Diff(ctx context.Context, req infer.DiffRequest[DatabaseConnectionArgs, DatabaseConnectionState]) (resp infer.DiffResponse, err error) {
// 	// Do stuff here
// 	if false {
// 		return infer.DiffResponse{
// 			HasChanges:   false,
// 			DetailedDiff: map[string]p.PropertyDiff{},
// 		}, nil
// 	}

// 	return infer.DiffResponse{}, nil
// }

func (w *DatabaseConnection) Delete(ctx context.Context, req infer.DeleteRequest[DatabaseConnectionState]) (infer.DeleteResponse, error) {
	config := infer.GetConfig[CloudbeaverProviderConfig](ctx)

	body := map[string]interface{}{
		"operationName": "deleteConnection",
		"query":         "mutation deleteConnection($projectId: ID!, $connectionId: ID!) { deleteConnection(projectId: $projectId, id: $connectionId) }",
		"variables": map[string]interface{}{
			"connectionId": req.ID,
			"projectId":    req.State.ProjectId,
		},
	}
	var responseBody interface{}
	_, err := sendPost(fmt.Sprintf("%s/api/gql", config.Endpoint), body, map[string]string{"cb-session-id": config.SessionId}, &responseBody)
	if err != nil {
		return infer.DeleteResponse{}, err
	}

	return infer.DeleteResponse{}, nil
}
