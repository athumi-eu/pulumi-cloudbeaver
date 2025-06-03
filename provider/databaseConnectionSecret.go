package main

import (
	"context"
	"fmt"

	"github.com/pulumi/pulumi-go-provider/infer"
)

func (c *realClient) CreateDatabaseConnectionSecret(databaseConnectionId string, projectId string, teamId string, entraGroupName string) error {
	body := map[string]interface{}{
		"operationName": "setConnectionSecret",
		"query":         "mutation setConnectionSecret($teamId: ID!, $projectId: ID!, $dataSourceId: ID!, $credentials: Object) {\n secret: secretSetToTeam(\n teamId: $teamId\n projectId: $projectId\n dataSourceId: $dataSourceId\n credentials: $credentials\n ) {\n subjectId\n authProperties {\n ...UserConnectionAuthProperties\n }\n }\n}\n \n fragment UserConnectionAuthProperties on ObjectPropertyInfo {\n id\n displayName\n description\n category\n dataType\n value\n validValues\n defaultValue\n length\n features\n required\n order\n conditions {\n ...Condition\n }\n}\n \n fragment Condition on Condition {\n expression\n conditionType\n}",
		"variables": map[string]interface{}{
			"dataSourceId": databaseConnectionId,
			"projectId":    projectId,
			"teamId":       teamId,
			"credentials": map[string]string{
				"azureGroupName": entraGroupName,
			},
		},
	}
	var responseBody interface{}
	_, err := sendPost(fmt.Sprintf("%s/api/gql", c.endpoint), body, map[string]string{"cb-session-id": c.sessionId}, &responseBody)
	return err
}

func (c *realClient) DeleteDatabaseConnectionSecret(databaseConnectionId string, projectId string, teamId string) error {
	body := map[string]interface{}{
		"operationName": "deleteConnectionSecret",
		"query":         "mutation deleteConnectionSecret($teamId: ID!, $projectId: ID!, $dataSourceId: ID!) {\n secretDeleteFromTeam(\n teamId: $teamId\n projectId: $projectId\n dataSourceId: $dataSourceId\n )\n}",
		"variables": map[string]interface{}{
			"dataSourceId": databaseConnectionId,
			"projectId":    projectId,
			"teamId":       teamId,
		},
	}
	var responseBody interface{}
	_, err := sendPost(fmt.Sprintf("%s/api/gql", c.endpoint), body, map[string]string{"cb-session-id": c.sessionId}, &responseBody)
	return err
}

type DatabaseConnectionSecret struct {
	getClient clientFactory
}

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
	config := infer.GetConfig[Config](ctx)

	state := DatabaseConnectionSecretState{DatabaseConnectionSecretArgs: req.Inputs}

	// Use the client factory to create a client based on the current config.
	client, err := w.getClient(ctx, config)
	if err != nil {
		return infer.CreateResponse[DatabaseConnectionSecretState]{}, err
	}

	// Use the client to create a project.
	err = client.CreateDatabaseConnectionSecret(req.Inputs.DatabaseConnectionId, req.Inputs.ProjectId, req.Inputs.TeamId, req.Inputs.EntraGroupName)
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
	config := infer.GetConfig[Config](ctx)

	// Use the client factory to create a client based on the current config.
	client, err := w.getClient(ctx, config)
	if err != nil {
		return infer.DeleteResponse{}, err
	}

	// Use the client to delete a project.
	err = client.DeleteDatabaseConnectionSecret(req.State.DatabaseConnectionId, req.State.ProjectId, req.State.TeamId)
	if err != nil {
		return infer.DeleteResponse{}, err
	}

	return infer.DeleteResponse{}, nil
}
