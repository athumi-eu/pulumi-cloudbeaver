package main

import (
	"context"
	"fmt"

	"github.com/pulumi/pulumi-go-provider/infer"
)

type ProjectMember struct {
}

type ProjectMemberArgs struct {
	ProjectId string `pulumi:"project_id"`
	MemberId  string `pulumi:"member_id"`
}
type ProjectMemberState struct {
	ProjectMemberArgs
}

func (w *ProjectMember) Create(ctx context.Context, req infer.CreateRequest[ProjectMemberArgs]) (infer.CreateResponse[ProjectMemberState], error) {
	state := ProjectMemberState{ProjectMemberArgs: req.Inputs}

	config := infer.GetConfig[CloudbeaverProviderConfig](ctx)

	body := map[string]interface{}{
		"operationName": "rmAddProjectsPermissions",
		"query":         "mutation rmAddProjectsPermissions($projectIds: [ID!]!, $subjectIds: [ID!]!, $permissions: [String!]!) { rmAddProjectsPermissions(projectIds: $projectIds, subjectIds: $subjectIds, permissions: $permissions) }",
		"variables": map[string]interface{}{
			// What are these exactly ?
			"permissions": [2]string{
				"project-datasource-view",
				"project-resource-view",
			},
			"projectIds": [1]string{
				req.Inputs.ProjectId,
			},
			"subjectIds": [1]string{
				req.Inputs.MemberId,
			},
		},
	}
	var responseBody interface{}
	_, err := sendPost(fmt.Sprintf("%s/api/gql", config.Endpoint), body, map[string]string{"cb-session-id": config.SessionId}, &responseBody)
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
	config := infer.GetConfig[CloudbeaverProviderConfig](ctx)

	body := map[string]interface{}{
		"operationName": "rmDeleteProjectsPermissions",
		"query":         "mutation rmDeleteProjectsPermissions($projectIds: [ID!]!, $subjectIds: [ID!]!, $permissions: [String!]!) { rmDeleteProjectsPermissions(projectIds: $projectIds, subjectIds: $subjectIds, permissions: $permissions) }",
		"variables": map[string]interface{}{
			"permissions": [2]string{
				"project-datasource-view",
				"project-resource-view",
			},
			"projectIds": [1]string{
				req.State.ProjectId,
			},
			"subjectIds": [1]string{
				req.State.MemberId,
			},
		},
	}
	var responseBody interface{}
	_, err := sendPost(fmt.Sprintf("%s/api/gql", config.Endpoint), body, map[string]string{"cb-session-id": config.SessionId}, &responseBody)
	if err != nil {
		return infer.DeleteResponse{}, err
	}

	return infer.DeleteResponse{}, nil
}
