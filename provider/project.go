package main

import (
	"context"
	"fmt"

	"github.com/pulumi/pulumi-go-provider/infer"
)

type Project struct {
}

type ProjectArgs struct {
	Name string `pulumi:"name"`
}
type ProjectState struct {
	ProjectArgs
}

func (w *Project) Create(ctx context.Context, req infer.CreateRequest[ProjectArgs]) (infer.CreateResponse[ProjectState], error) {
	id := fmt.Sprintf("s_%s", req.Inputs.Name)
	state := ProjectState{ProjectArgs: req.Inputs}

	config := infer.GetConfig[CloudbeaverProviderConfig](ctx)

	body := map[string]interface{}{
		"operationName": "createProject",
		"query":         "mutation rmCreateProject($projectId: ID, $projectName: String!, $description: String) { rmCreateProject(projectId: $projectId, projectName: $projectName, description: $description) { id } }",
		"variables": map[string]interface{}{
			"projectId":   id,
			"projectName": req.Inputs.Name,
			"description": "",
		},
	}
	var responseBody interface{}
	_, err := sendPost(fmt.Sprintf("%s/api/gql", config.Endpoint), body, map[string]string{"cb-session-id": config.SessionId}, &responseBody)
	if err != nil {
		return infer.CreateResponse[ProjectState]{}, err
	}

	return infer.CreateResponse[ProjectState]{
		ID:     id,
		Output: state,
	}, nil
}

// func (w *Project) Read(ctx context.Context, req infer.ReadRequest[ProjectArgs, ProjectState]) (resp infer.ReadResponse[ProjectArgs, ProjectState], err error) {
// 	// Do stuff here
// 	return infer.ReadResponse[ProjectArgs, ProjectState]{
// 		ID:     req.ID,
// 		Inputs: req.Inputs,
// 		State:  req.State,
// 	}, nil
// }

// func (w *Project) Diff(ctx context.Context, req infer.DiffRequest[ProjectArgs, ProjectState]) (resp infer.DiffResponse, err error) {
// 	// Do stuff here
// 	if false {
// 		return infer.DiffResponse{
// 			HasChanges:   false,
// 			DetailedDiff: map[string]p.PropertyDiff{},
// 		}, nil
// 	}

// 	return infer.DiffResponse{}, nil
// }

func (w *Project) Delete(ctx context.Context, req infer.DeleteRequest[ProjectState]) (infer.DeleteResponse, error) {
	config := infer.GetConfig[CloudbeaverProviderConfig](ctx)

	body := map[string]interface{}{
		"operationName": "deleteProject",
		"query":         "mutation rmDeleteProject($projectId: ID!) { rmDeleteProject(projectId: $projectId) }",
		"variables": map[string]interface{}{
			"projectId": req.ID,
		},
	}
	var responseBody interface{}
	_, err := sendPost(fmt.Sprintf("%s/api/gql", config.Endpoint), body, map[string]string{"cb-session-id": config.SessionId}, &responseBody)
	if err != nil {
		return infer.DeleteResponse{}, err
	}

	return infer.DeleteResponse{}, nil
}
