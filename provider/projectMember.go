package main

import (
	"context"
	"fmt"

	"github.com/pulumi/pulumi-go-provider/infer"
)

func (c *realClient) CreateProjectMember(projectId string, memberId string) error {
	body := map[string]interface{}{
		"operationName": "addProjectsPermissions",
		"query":         "mutation addProjectsPermissions($projectIds: [ID!]!, $subjectIds: [ID!]!, $permissions: [String!]!) {\n rmAddProjectsPermissions(\n projectIds: $projectIds\n subjectIds: $subjectIds\n permissions: $permissions\n )\n}",
		"variables": map[string]interface{}{
			// What are these exactly ?
			"permissions": [2]string{
				"project-datasource-view",
				"project-resource-view",
			},
			"projectIds": [1]string{
				projectId,
			},
			"subjectIds": [1]string{
				memberId,
			},
		},
	}
	var responseBody interface{}
	_, err := sendPost(fmt.Sprintf("%s/api/gql", c.endpoint), body, map[string]string{"cb-session-id": c.sessionId}, &responseBody)
	return err
}

func (c *realClient) DeleteProjectMember(projectId string, memberId string) error {
	body := map[string]interface{}{
		"operationName": "deleteProjectsPermissions",
		"query":         "mutation deleteProjectsPermissions($projectIds: [ID!]!, $subjectIds: [ID!]!, $permissions: [String!]!) {\n rmDeleteProjectsPermissions(\n projectIds: $projectIds\n subjectIds: $subjectIds\n permissions: $permissions\n )\n}",
		"variables": map[string]interface{}{
			"permissions": [2]string{
				"project-datasource-view",
				"project-resource-view",
			},
			"projectIds": [1]string{
				projectId,
			},
			"subjectIds": [1]string{
				memberId,
			},
		},
	}
	var responseBody interface{}
	_, err := sendPost(fmt.Sprintf("%s/api/gql", c.endpoint), body, map[string]string{"cb-session-id": c.sessionId}, &responseBody)
	return err
}

type ProjectMember struct {
	getClient clientFactory
}

type ProjectMemberArgs struct {
	ProjectId string `pulumi:"project_id"`
	MemberId  string `pulumi:"member_id"`
}
type ProjectMemberState struct {
	ProjectMemberArgs
}

func (w *ProjectMember) Create(ctx context.Context, req infer.CreateRequest[ProjectMemberArgs]) (infer.CreateResponse[ProjectMemberState], error) {
	config := infer.GetConfig[Config](ctx)

	state := ProjectMemberState{ProjectMemberArgs: req.Inputs}

	// Use the client factory to create a client based on the current config.
	client, err := w.getClient(ctx, config)
	if err != nil {
		return infer.CreateResponse[ProjectMemberState]{}, err
	}

	// Use the client to create a project.
	err = client.CreateProjectMember(req.Inputs.ProjectId, req.Inputs.MemberId)
	if err != nil {
		return infer.CreateResponse[ProjectMemberState]{}, err
	}

	return infer.CreateResponse[ProjectMemberState]{
		ID:     fmt.Sprintf("%s::%s", req.Inputs.ProjectId, req.Inputs.MemberId),
		Output: state,
	}, nil
}

// func (w *ProjectMember) Read(ctx context.Context, req infer.ReadRequest[ProjectMemberArgs, ProjectMemberState]) (resp infer.ReadResponse[ProjectMemberArgs, ProjectMemberState], err error) {
// 	// Do stuff here
// 	return infer.ReadResponse[ProjectMemberArgs, ProjectMemberState]{
// 		ID:     req.ID,
// 		Inputs: req.Inputs,
// 		State:  req.State,
// 	}, nil
// }

// func (w *ProjectMember) Diff(ctx context.Context, req infer.DiffRequest[ProjectMemberArgs, ProjectMemberState]) (resp infer.DiffResponse, err error) {
// 	// Do stuff here
// 	if false {
// 		return infer.DiffResponse{
// 			HasChanges:   false,
// 			DetailedDiff: map[string]p.PropertyDiff{},
// 		}, nil
// 	}

// 	return infer.DiffResponse{}, nil
// }

func (w *ProjectMember) Delete(ctx context.Context, req infer.DeleteRequest[ProjectMemberState]) (infer.DeleteResponse, error) {
	config := infer.GetConfig[Config](ctx)

	// Use the client factory to create a client based on the current config.
	client, err := w.getClient(ctx, config)
	if err != nil {
		return infer.DeleteResponse{}, err
	}

	// Use the client to delete a project.
	err = client.DeleteProjectMember(req.State.ProjectId, req.State.MemberId)
	if err != nil {
		return infer.DeleteResponse{}, err
	}

	return infer.DeleteResponse{}, nil
}
