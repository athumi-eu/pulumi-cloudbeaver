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

func (c *realClient) CreateDatabaseConnection(name string, endpoint string, database string, projectId string) (string, error) {
	body := map[string]interface{}{
		"operationName": "createConnection",
		"query":         "mutation createConnection($projectId: ID!, $config: ConnectionConfig!, $includeAuthProperties: Boolean!, $includeNetworkHandlersConfig: Boolean!, $includeCredentialsSaved: Boolean!, $includeAuthNeeded: Boolean!, $includeProperties: Boolean!, $includeProviderProperties: Boolean!, $customIncludeOptions: Boolean!) {\n  connection: createConnection(projectId: $projectId, config: $config) {\n    ...DatabaseConnection\n  }\n}\n    \n    fragment DatabaseConnection on ConnectionInfo {\n  id\n  projectId\n  name\n  description\n  driverId\n  connected\n  readOnly\n  saveCredentials\n  credentialsSaved @include(if: $includeCredentialsSaved)\n  sharedCredentials\n  folder\n  nodePath\n  mainPropertyValues\n  configurationType @include(if: $customIncludeOptions)\n  useUrl @include(if: $customIncludeOptions)\n  host @include(if: $customIncludeOptions)\n  port @include(if: $customIncludeOptions)\n  serverName @include(if: $customIncludeOptions)\n  databaseName @include(if: $customIncludeOptions)\n  url @include(if: $customIncludeOptions)\n  properties @include(if: $includeProperties)\n  providerProperties @include(if: $includeProviderProperties)\n  requiredAuth\n  features\n  supportedDataFormats\n  authNeeded @include(if: $includeAuthNeeded)\n  authModel\n  authProperties @include(if: $includeAuthProperties) {\n    ...UserConnectionAuthProperties\n  }\n  networkHandlersConfig @skip(if: $includeNetworkHandlersConfig) {\n    ...NetworkHandlerBasics\n  }\n  networkHandlersConfig @include(if: $includeNetworkHandlersConfig) {\n    ...NetworkHandlerBasics\n    authType\n    userName\n    password\n    key\n    savePassword\n    properties\n    secureProperties\n  }\n  navigatorSettings {\n    ...AllNavigatorSettings\n  }\n  canViewSettings\n  canEdit\n  canDelete\n}\n    \n    fragment UserConnectionAuthProperties on ObjectPropertyInfo {\n  id\n  displayName\n  description\n  category\n  dataType\n  value\n  validValues\n  defaultValue\n  length\n  features\n  required\n  order\n  conditions {\n    ...Condition\n  }\n}\n    \n    fragment Condition on Condition {\n  expression\n  conditionType\n}\n    \n\n    fragment NetworkHandlerBasics on NetworkHandlerConfig {\n  id\n  enabled\n}\n    \n\n    fragment AllNavigatorSettings on NavigatorSettings {\n  showSystemObjects\n  showUtilityObjects\n  showOnlyEntities\n  mergeEntities\n  hideFolders\n  hideSchemas\n  hideVirtualModel\n}",
		"variables": map[string]interface{}{
			"config": map[string]interface{}{
				"authModelId":       "azure_ad_postgresql",
				"configurationType": "MANUAL",
				"credentials": map[string]bool{
					"useLegacyToken": false,
				},
				"databaseName": database,
				"driverId":     "postgresql:postgres-jdbc",
				"host":         endpoint,
				"mainPropertyValues": map[string]string{
					"database": database,
					"host":     endpoint,
					"port":     "5432",
				},
				"name":                  name,
				"networkHandlersConfig": []string{},
				"port":                  "5432",
				"properties":            map[string]string{},
				"providerProperties":    map[string]string{},
				"saveCredentials":       true,
				"sharedCredentials":     true,
			},
			"includeNetworkHandlersConfig": false,
			"includeAuthProperties":        false,
			"includeAuthNeeded":            false,
			"includeCredentialsSaved":      false,
			"includeProperties":            false,
			"includeProviderProperties":    false,
			"customIncludeOptions":         false,
			"customIncludeBase":            true,
			"projectId":                    projectId,
		},
	}
	var responseBody CreateDatabaseConnectionResponse
	_, err := sendPost(fmt.Sprintf("%s/api/gql", c.endpoint), body, map[string]string{"cb-session-id": c.sessionId}, &responseBody)

	return responseBody.Data.Connection.ID, err
}

func (c *realClient) DeleteDatabaseConnection(id string, projectId string) error {
	body := map[string]interface{}{
		"operationName": "deleteConnection",
		"query":         "mutation deleteConnection($projectId: ID!, $connectionId: ID!) {\n deleteConnection(projectId: $projectId, id: $connectionId)\n}",
		"variables": map[string]interface{}{
			"connectionId": id,
			"projectId":    projectId,
		},
	}
	var responseBody interface{}
	_, err := sendPost(fmt.Sprintf("%s/api/gql", c.endpoint), body, map[string]string{"cb-session-id": c.sessionId}, responseBody)
	return err
}

type DatabaseConnection struct {
	getClient clientFactory
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
	config := infer.GetConfig[Config](ctx)

	state := DatabaseConnectionState{DatabaseConnectionArgs: req.Inputs}

	// Use the client factory to create a client based on the current config.
	client, err := w.getClient(ctx, config)
	if err != nil {
		return infer.CreateResponse[DatabaseConnectionState]{}, err
	}

	// Use the client to create a project.
	connectionId, err := client.CreateDatabaseConnection(req.Inputs.Name, req.Inputs.Endpoint, req.Inputs.Database, req.Inputs.ProjectId)
	if err != nil {
		return infer.CreateResponse[DatabaseConnectionState]{}, err
	}

	return infer.CreateResponse[DatabaseConnectionState]{
		ID:     connectionId,
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
	config := infer.GetConfig[Config](ctx)

	// Use the client factory to create a client based on the current config.
	client, err := w.getClient(ctx, config)
	if err != nil {
		return infer.DeleteResponse{}, err
	}

	// Use the client to delete a project.
	err = client.DeleteDatabaseConnection(req.ID, req.State.ProjectId)
	if err != nil {
		return infer.DeleteResponse{}, err
	}

	return infer.DeleteResponse{}, nil
}
