package main

import (
	"context"
	"fmt"
	"os"

	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/infer"
	"github.com/pulumi/pulumi/sdk/v3/go/common/tokens"
)

func main() {
	// Create a provider that uses the real client factory.
	provider, err := provider(newRealClient)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s", err.Error())
		os.Exit(1)
	}

	err = provider.Run(context.Background(), "cloudbeaver", "0.1.0")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s", err.Error())
		os.Exit(1)
	}
}

func provider(factory clientFactory) (p.Provider, error) {
	return infer.NewProviderBuilder().
		WithDisplayName("Cloudbeaver").
		WithPublisher("athumi-eu").
		WithPluginDownloadURL("github://api.github.com/athumi-eu").
		WithNamespace("cloudbeaver").
		WithConfig(infer.Config(&Config{})).
		WithResources(
			infer.Resource(&Team{getClient: factory}),
			infer.Resource(&Project{getClient: factory}),
			infer.Resource(&ProjectMember{getClient: factory}),
			infer.Resource(&DatabaseConnection{getClient: factory}),
			infer.Resource(&DatabaseConnectionSecret{getClient: factory}),
		).
		WithModuleMap(map[tokens.ModuleName]tokens.ModuleName{
			"cloudbeaver": "index",
		}).
		Build()
}
