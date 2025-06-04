package main

import (
	"context"

	"github.com/pulumi/pulumi-go-provider/infer"
)

type CloudbeaverProviderConfig struct {
	Endpoint  string `pulumi:"endpoint"`
	ApiKey    string `pulumi:"apiKey" provider:"secret"`
	SessionId string
}

func (c *CloudbeaverProviderConfig) Annotate(a infer.Annotator) {
	a.Describe(&c.Endpoint, "The cloudbeaver endpoint to connect to.")
	a.Describe(&c.ApiKey, "The cloudbeaver API key to use.")
}

func (p *CloudbeaverProviderConfig) Configure(ctx context.Context) error {
	sessionId, err := NewCloudbeaverSession(p.ApiKey, p.Endpoint)
	if err != nil {
		return err
	}
	p.SessionId = sessionId

	return nil
}
