package main

import (
	"context"
	"fmt"
	"os"

	"github.com/pulumi/pulumi-go-provider/infer"
	"github.com/pulumi/pulumi/sdk/v3/go/common/tokens"
)

var version string = "0.3.1"

func main() {
	p, err := infer.NewProviderBuilder().
		WithDisplayName("Cloudbeaver").
		WithPublisher("athumi-eu").
		WithPluginDownloadURL("github://api.github.com/athumi-eu").
		WithNamespace("pulumi").
		WithConfig(infer.Config(&CloudbeaverProviderConfig{})).
		WithFunctions(
			infer.Function(&EnableUser{}),
		).
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

	err = p.Run(context.Background(), "cloudbeaver", version)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s", err.Error())
		os.Exit(1)
	}
}
