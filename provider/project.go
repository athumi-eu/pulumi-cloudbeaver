package main

import (
	"context"
	"fmt"

	"github.com/pulumi/pulumi-go-provider/infer"
)

func (c *realClient) CreateProject(id string, name string) error {
	body := map[string]interface{}{
		"operationName": "createProject",
		"query":         "mutation createProject($projectId: ID, $projectName: String!, $description: String) {\n project: rmCreateProject(\n projectId: $projectId\n projectName: $projectName\n description: $description\n ) {\n ...SharedProject\n }\n}\n \n fragment SharedProject on RMProject {\n id\n name\n shared\n global\n description\n projectPermissions\n resourceTypes {\n ...ResourceType\n }\n}\n \n fragment ResourceType on RMResourceType {\n id\n displayName\n icon\n fileExtensions\n rootFolder\n}",
		"variables": map[string]interface{}{
			"projectId":   id,
			"projectName": name,
		},
	}
	var responseBody interface{}
	_, err := sendPost(fmt.Sprintf("%s/api/gql", c.endpoint), body, map[string]string{"cb-session-id": c.sessionId}, &responseBody)
	return err
}

func (c *realClient) DeleteProject(id string) error {
	body := map[string]interface{}{
		"operationName": "deleteProject",
		"query":         "mutation deleteProject($projectId: ID!) {\n rmDeleteProject(projectId: $projectId)\n}",
		"variables": map[string]interface{}{
			"projectId": id,
		},
	}
	var responseBody interface{}
	_, err := sendPost(fmt.Sprintf("%s/api/gql", c.endpoint), body, map[string]string{"cb-session-id": c.sessionId}, &responseBody)
	return err
}

type Project struct {
	getClient clientFactory
}

type ProjectArgs struct {
	Name string `pulumi:"name"`
}
type ProjectState struct {
	ProjectArgs
}

func (w *Project) Create(ctx context.Context, req infer.CreateRequest[ProjectArgs]) (infer.CreateResponse[ProjectState], error) {
	config := infer.GetConfig[Config](ctx)

	id := fmt.Sprintf("s_%s", req.Inputs.Name)
	state := ProjectState{ProjectArgs: req.Inputs}

	// Use the client factory to create a client based on the current config.
	client, err := w.getClient(ctx, config)
	if err != nil {
		return infer.CreateResponse[ProjectState]{}, err
	}

	// Use the client to create a project.
	err = client.CreateProject(id, req.Inputs.Name)
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
	config := infer.GetConfig[Config](ctx)

	// Use the client factory to create a client based on the current config.
	client, err := w.getClient(ctx, config)
	if err != nil {
		return infer.DeleteResponse{}, err
	}

	// Use the client to delete a project.
	err = client.DeleteProject(req.ID)
	if err != nil {
		return infer.DeleteResponse{}, err
	}

	return infer.DeleteResponse{}, nil
}
