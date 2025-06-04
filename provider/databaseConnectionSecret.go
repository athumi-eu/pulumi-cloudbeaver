package main

import (
	"context"
	"fmt"

	"github.com/pulumi/pulumi-go-provider/infer"
)

type DatabaseConnectionSecret struct{}

type DatabaseConnectionSecretArgs struct {
	DatabaseConnectionId string `pulumi:"database_connection_id"`
	ProjectId            string `pulumi:"project_id"`
	TeamId               string `pulumi:"team_id"`
	EntraGroupName       string `pulumi:"entra_group_name"`
}
type DatabaseConnectionSecretState struct {
	DatabaseConnectionSecretArgs
}

func (w *DatabaseConnectionSecret) Create(ctx context.Context, req infer.CreateRequest[DatabaseConnectionSecretArgs]) (infer.CreateResponse[DatabaseConnectionSecretState], error) {
	state := DatabaseConnectionSecretState{DatabaseConnectionSecretArgs: req.Inputs}

	config := infer.GetConfig[CloudbeaverProviderConfig](ctx)

	body := map[string]interface{}{
		"operationName": "secretSetToTeam",
		"query":         "mutation secretSetToTeam($teamId: ID!, $projectId: ID!, $dataSourceId: ID!, $credentials: Object) { secretSetToTeam(teamId: $teamId, projectId: $projectId, dataSourceId: $dataSourceId, credentials: $credentials) { subjectId } }",
		"variables": map[string]interface{}{
			"dataSourceId": req.Inputs.DatabaseConnectionId,
			"projectId":    req.Inputs.ProjectId,
			"teamId":       req.Inputs.TeamId,
			"credentials": map[string]string{
				"azureGroupName": req.Inputs.EntraGroupName,
			},
		},
	}
	var responseBody interface{}
	_, err := sendPost(fmt.Sprintf("%s/api/gql", config.Endpoint), body, map[string]string{"cb-session-id": config.SessionId}, &responseBody)
	if err != nil {
		return infer.CreateResponse[DatabaseConnectionSecretState]{}, err
	}

	return infer.CreateResponse[DatabaseConnectionSecretState]{
		ID:     fmt.Sprintf("%s::%s::%s", req.Inputs.DatabaseConnectionId, req.Inputs.ProjectId, req.Inputs.TeamId),
		Output: state,
	}, nil
}

// func (w *DatabaseConnectionSecret) Read(ctx context.Context, req infer.ReadRequest[DatabaseConnectionSecretArgs, DatabaseConnectionSecretState]) (resp infer.ReadResponse[DatabaseConnectionSecretArgs, DatabaseConnectionSecretState], err error) {
// 	// Do stuff here
// 	return infer.ReadResponse[DatabaseConnectionSecretArgs, DatabaseConnectionSecretState]{
// 		ID:     req.ID,
// 		Inputs: req.Inputs,
// 		State:  req.State,
// 	}, nil
// }

// func (w *DatabaseConnectionSecret) Diff(ctx context.Context, req infer.DiffRequest[DatabaseConnectionSecretArgs, DatabaseConnectionSecretState]) (resp infer.DiffResponse, err error) {
// 	// Do stuff here
// 	if false {
// 		return infer.DiffResponse{
// 			HasChanges:   false,
// 			DetailedDiff: map[string]p.PropertyDiff{},
// 		}, nil
// 	}

// 	return infer.DiffResponse{}, nil
// }

func (w *DatabaseConnectionSecret) Delete(ctx context.Context, req infer.DeleteRequest[DatabaseConnectionSecretState]) (infer.DeleteResponse, error) {
	config := infer.GetConfig[CloudbeaverProviderConfig](ctx)

	body := map[string]interface{}{
		"operationName": "secretDeleteFromTeam",
		"query":         "mutation secretDeleteFromTeam($teamId: ID!, $projectId: ID!, $dataSourceId: ID!) { secretDeleteFromTeam(teamId: $teamId, projectId: $projectId, dataSourceId: $dataSourceId) }",
		"variables": map[string]interface{}{
			"dataSourceId": req.State.DatabaseConnectionId,
			"projectId":    req.State.ProjectId,
			"teamId":       req.State.TeamId,
		},
	}
	var responseBody interface{}
	_, err := sendPost(fmt.Sprintf("%s/api/gql", config.Endpoint), body, map[string]string{"cb-session-id": config.SessionId}, &responseBody)
	if err != nil {
		return infer.DeleteResponse{}, err
	}

	return infer.DeleteResponse{}, nil
}
