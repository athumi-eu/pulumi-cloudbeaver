package main

import (
	"context"
	"fmt"
	"os"

	"github.com/pulumi/pulumi-go-provider/infer"
	"github.com/pulumi/pulumi/sdk/v3/go/common/tokens"
)

func main() {
	p, err := infer.NewProviderBuilder().
		WithDisplayName("Cloudbeaver").
		WithPublisher("athumi-eu").
		WithPluginDownloadURL("github://api.github.com/athumi-eu").
		WithNamespace("cloudbeaver").
		WithConfig(infer.Config(&CloudbeaverProviderConfig{})).
		WithResources(
			infer.Resource(&Team{}),
			infer.Resource(&Project{}),
			infer.Resource(&ProjectMember{}),
			infer.Resource(&DatabaseConnection{}),
			infer.Resource(&DatabaseConnectionSecret{}),
		).
		WithModuleMap(map[tokens.ModuleName]tokens.ModuleName{
			"cloudbeaver": "index",
		}).
		Build()

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s", err.Error())
		os.Exit(1)
	}

	err = p.Run(context.Background(), "cloudbeaver", "0.1.0")

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s", err.Error())
		os.Exit(1)
	}
}
