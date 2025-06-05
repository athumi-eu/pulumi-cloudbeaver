package main

import (
	"context"
	"fmt"

	"github.com/pulumi/pulumi-go-provider/infer"
)

type Team struct {
}

type TeamArgs struct {
	Name         string  `pulumi:"name"`
	Description  string  `pulumi:"description,optional"`
	EntraGroupId *string `pulumi:"entra_group_id,optional"`
}
type TeamState struct {
	TeamArgs
}

func (w *Team) Create(ctx context.Context, req infer.CreateRequest[TeamArgs]) (infer.CreateResponse[TeamState], error) {
	id := fmt.Sprintf("t_%s", req.Inputs.Name)
	state := TeamState{TeamArgs: req.Inputs}

	config := infer.GetConfig[CloudbeaverProviderConfig](ctx)

	if req.DryRun {
		return infer.CreateResponse[TeamState]{
			ID:     id,
			Output: state,
		}, nil
	}

	body := map[string]interface{}{
		"operationName": "createTeam",
		"query":         "query createTeam($teamId: ID!, $teamName: String, $description: String) { team: createTeam(teamId: $teamId, teamName: $teamName, description: $description) { teamId }}",
		"variables": map[string]interface{}{
			"teamId":      id,
			"teamName":    req.Inputs.Name,
			"description": req.Inputs.Description,
		},
	}
	var responseBody interface{}
	_, err := sendPost(fmt.Sprintf("%s/api/gql", config.Endpoint), body, map[string]string{"cb-session-id": config.SessionId}, &responseBody)
	if err != nil {
		return infer.CreateResponse[TeamState]{
			ID:     id,
			Output: state,
		}, err
	}

	if req.Inputs.EntraGroupId != nil {
		body := map[string]interface{}{
			"operationName": "setTeamMetaParameterValues",
			"query":         "query setTeamMetaParameterValues($teamId: ID!, $parameters: Object!) { setTeamMetaParameterValues(teamId: $teamId, parameters: $parameters) }",
			"variables": map[string]interface{}{
				"teamId": id,
				"parameters": map[string]string{
					"aad.group-id": *req.Inputs.EntraGroupId,
				},
			},
		}
		var responseBody interface{}
		_, err = sendPost(fmt.Sprintf("%s/api/gql", config.Endpoint), body, map[string]string{"cb-session-id": config.SessionId}, &responseBody)
		if err != nil {
			return infer.CreateResponse[TeamState]{
				ID:     id,
				Output: state,
			}, err
		}
	}

	return infer.CreateResponse[TeamState]{
		ID:     id,
		Output: state,
	}, nil
}

// func (w *Team) Read(ctx context.Context, req infer.ReadRequest[TeamArgs, TeamState]) (resp infer.ReadResponse[TeamArgs, TeamState], err error) {
// 	// Do stuff here
// 	return infer.ReadResponse[TeamArgs, TeamState]{
// 		ID:     req.ID,
// 		Inputs: req.Inputs,
// 		State:  req.State,
// 	}, nil
// }

// func (w *Team) Diff(ctx context.Context, req infer.DiffRequest[TeamArgs, TeamState]) (resp infer.DiffResponse, err error) {
// 	// Do stuff here
// 	if false {
// 		return infer.DiffResponse{
// 			HasChanges:   false,
// 			DetailedDiff: map[string]p.PropertyDiff{},
// 		}, nil
// 	}

// 	return infer.DiffResponse{}, nil
// }

func (w *Team) Delete(ctx context.Context, req infer.DeleteRequest[TeamState]) (infer.DeleteResponse, error) {
	config := infer.GetConfig[CloudbeaverProviderConfig](ctx)

	body := map[string]interface{}{
		"operationName": "deleteTeam",
		"query":         "query deleteTeam($teamId: ID!, $force: Boolean) { deleteTeam(teamId: $teamId, force: $force) }",
		"variables": map[string]interface{}{
			"teamId": req.ID,
			"force":  true,
		},
	}
	var responseBody interface{}
	_, err := sendPost(fmt.Sprintf("%s/api/gql", config.Endpoint), body, map[string]string{"cb-session-id": config.SessionId}, &responseBody)
	if err != nil {
		return infer.DeleteResponse{}, err
	}

	return infer.DeleteResponse{}, nil
}
