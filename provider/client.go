package main

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/pulumi/pulumi-go-provider/infer"
)

type Config struct {
	ClientEndpoint string `pulumi:"endpoint"`
	ClientUsername string `pulumi:"username"`
	ClientPassword string `pulumi:"password" provider:"secret"`
}

func (c *Config) Annotate(a infer.Annotator) {
	a.Describe(&c.ClientEndpoint, "The cloudbeaver endpoint to connect to.")
	a.Describe(&c.ClientUsername, "The cloudbeaver username to use.")
	a.Describe(&c.ClientPassword, "The cloudbeaver password to use.")
}

type client interface {
	CreateTeam(id string, name string, description string) error
	CreateTeamEntraGroup(id string, groupId string) error
	DeleteTeam(id string) error
	CreateProject(id string, name string) error
	DeleteProject(id string) error
	CreateProjectMember(projectId string, memberId string) error
	DeleteProjectMember(projectId string, memberId string) error
	CreateDatabaseConnection(name string, endpoint string, database string, projectId string) (string, error)
	DeleteDatabaseConnection(id string, projectId string) error
	CreateDatabaseConnectionSecret(databaseConnectionId string, projectId string, teamId string, entraGroupName string) error
	DeleteDatabaseConnectionSecret(databaseConnectionId string, projectId string, teamId string) error
}

type realClient struct {
	endpoint  string
	username  string
	password  string
	sessionId string
}

var cachedClients = map[string]client{}

type clientFactory func(ctx context.Context, config Config) (client, error)

// newRealClient implements clientFactory and creates a real client based on the provider config.
func newRealClient(ctx context.Context, config Config) (client, error) {
	if val, ok := cachedClients[config.ClientUsername]; ok {
		return val, nil
	}
	_client := &realClient{
		endpoint:  config.ClientEndpoint,
		username:  config.ClientUsername,
		password:  config.ClientPassword,
		sessionId: "",
	}
	// Get a session id
	sessionBody := map[string]interface{}{
		"operationName": "serverConfig",
		"query":         "query serverConfig {\n  serverConfig {\n    ...ServerConfig\n  }\n}\n    \n    fragment ServerConfig on ServerConfig {\n  name\n  version\n  workspaceId\n  serverURL\n  rootURI\n  containerId\n  defaultAuthRole\n  defaultUserTeam\n  productConfiguration\n  supportsCustomConnections\n  sessionExpireTime\n  anonymousAccessEnabled\n  adminCredentialsSaveEnabled\n  publicCredentialsSaveEnabled\n  resourceManagerEnabled\n  secretManagerEnabled\n  configurationMode\n  developmentMode\n  redirectOnFederatedAuth\n  distributed\n  enabledFeatures\n  disabledBetaFeatures\n  enabledAuthProviders\n  supportedLanguages {\n    isoCode\n    displayName\n    nativeName\n  }\n  disabledDrivers\n}",
	}

	var responseBody interface{}
	sessionCookieResponse, err := sendPost(fmt.Sprintf("%s/api/gql", config.ClientEndpoint), sessionBody, map[string]string{}, &responseBody)
	if err != nil {
		return nil, err
	}

	for _, cookie := range sessionCookieResponse.Cookies() {
		if cookie.Name == "cb-session-id" {
			_client.sessionId = cookie.Value
			break
		}
	}

	// Authenticate the session
	hash := md5.Sum([]byte(config.ClientPassword))
	hashedPassword := strings.ToUpper(hex.EncodeToString(hash[:]))

	authBody := map[string]interface{}{
		"operationName": "authLogin",
		"query":         "\n query authLogin($provider: ID!, $configuration: ID, $credentials: Object, $linkUser: Boolean, $forceSessionsLogout: Boolean) {\n authInfo: authLogin(\n provider: $provider\n configuration: $configuration\n credentials: $credentials\n linkUser: $linkUser\n forceSessionsLogout: $forceSessionsLogout\n ) {\n redirectLink\n authId\n authStatus\n userTokens {\n ...AuthToken\n }\n }\n}\n \n fragment AuthToken on UserAuthToken {\n authProvider\n authConfiguration\n loginTime\n message\n origin {\n ...ObjectOriginInfo\n }\n}\n \n fragment ObjectOriginInfo on ObjectOrigin {\n type\n subType\n displayName\n icon\n}\n ",
		"variables": map[string]interface{}{
			"credentials": map[string]string{
				"password": hashedPassword,
				"user":     config.ClientUsername,
			},
			"provider":            "local",
			"forceSessionsLogout": false,
			"linkUser":            false,
		},
	}

	_, err = sendPost(fmt.Sprintf("%s/api/gql", config.ClientEndpoint), authBody, map[string]string{"cb-session-id": _client.sessionId}, &responseBody)
	if err != nil {
		return nil, err
	}

	cachedClients[config.ClientUsername] = _client

	return _client, nil
}
