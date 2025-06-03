package main

import (
	"context"
	"fmt"

	"github.com/pulumi/pulumi-go-provider/infer"
)

func (c *realClient) CreateTeam(id string, name string, description string) error {
	body := map[string]interface{}{
		"operationName": "createTeam",
		"query":         "query createTeam($teamId: ID!, $teamName: String, $description: String) {\n  team: createTeam(\n    teamId: $teamId\n    teamName: $teamName\n    description: $description\n  ) {\n    ...AdminTeamInfo\n  }\n}\n    \n    fragment AdminTeamInfo on AdminTeamInfo {\n  teamId\n  teamName\n  description\n  teamPermissions\n}",
		"variables": map[string]interface{}{
			"teamId":            id,
			"teamName":          name,
			"description":       description,
			"customIncludeBase": false,
		},
	}
	var responseBody interface{}
	_, err := sendPost(fmt.Sprintf("%s/api/gql", c.endpoint), body, map[string]string{"cb-session-id": c.sessionId}, &responseBody)
	return err
}

func (c *realClient) CreateTeamEntraGroup(id string, groupId string) error {
	body := map[string]interface{}{
		"operationName": "saveTeamMetaParameters",
		"query":         "query saveTeamMetaParameters($teamId: ID!, $parameters: Object!) {\n setTeamMetaParameterValues(teamId: $teamId, parameters: $parameters)\n}",
		"variables": map[string]interface{}{
			"teamId": id,
			"parameters": map[string]string{
				"aad.group-id": groupId,
			},
		},
	}
	var responseBody interface{}
	_, err := sendPost(fmt.Sprintf("%s/api/gql", c.endpoint), body, map[string]string{"cb-session-id": c.sessionId}, &responseBody)
	return err
}

func (c *realClient) DeleteTeam(id string) error {
	body := map[string]interface{}{
		"operationName": "deleteTeam",
		"query":         "query deleteTeam($teamId: ID!, $force: Boolean) {\n  deleteTeam(teamId: $teamId, force: $force)\n}",
		"variables": map[string]interface{}{
			"teamId": id,
			"force":  true,
		},
	}
	var responseBody interface{}
	_, err := sendPost(fmt.Sprintf("%s/api/gql", c.endpoint), body, map[string]string{"cb-session-id": c.sessionId}, &responseBody)
	return err
}

type Team struct {
	getClient clientFactory
}

type TeamArgs struct {
	Name         string  `pulumi:"name"`
	Description  string  `pulumi:"description"`
	EntraGroupId *string `pulumi:"entra_group_id,optional"`
}
type TeamState struct {
	TeamArgs
}

func (w *Team) Create(ctx context.Context, req infer.CreateRequest[TeamArgs]) (infer.CreateResponse[TeamState], error) {
	config := infer.GetConfig[Config](ctx)

	id := fmt.Sprintf("t_%s", req.Inputs.Name)
	state := TeamState{TeamArgs: req.Inputs}

	// Use the client factory to create a client based on the current config.
	client, err := w.getClient(ctx, config)
	if err != nil {
		return infer.CreateResponse[TeamState]{}, err
	}

	// Use the client to create a team.
	err = client.CreateTeam(id, req.Inputs.Name, req.Inputs.Description)
	if err != nil {
		return infer.CreateResponse[TeamState]{}, err
	}

	if req.Inputs.EntraGroupId != nil {
		err = client.CreateTeamEntraGroup(id, *req.Inputs.EntraGroupId)
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
	config := infer.GetConfig[Config](ctx)

	// Use the client factory to create a client based on the current config.
	client, err := w.getClient(ctx, config)
	if err != nil {
		return infer.DeleteResponse{}, err
	}

	// Use the client to delete a team.
	err = client.DeleteTeam(req.ID)
	if err != nil {
		return infer.DeleteResponse{}, err
	}

	return infer.DeleteResponse{}, nil
}
