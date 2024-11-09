package cmd

import (
	"context"
	"fmt"

	corerp "github.com/radius-project/radius/pkg/corerp/api/v20231001preview"
	"github.com/radius-project/radius/pkg/to"
	"github.com/spf13/cobra"

	v1 "github.com/radius-project/radius/pkg/armrpc/api/v1"
	"github.com/radius-project/radius/pkg/cli"
	"github.com/radius-project/radius/pkg/cli/clients"
	"github.com/radius-project/radius/pkg/cli/clierrors"
	"github.com/radius-project/radius/pkg/cli/cmd/commonflags"
	"github.com/radius-project/radius/pkg/cli/cmd/env/namespace"
	"github.com/radius-project/radius/pkg/cli/connections"
	"github.com/radius-project/radius/pkg/cli/framework"
	"github.com/radius-project/radius/pkg/cli/kubernetes"
	"github.com/radius-project/radius/pkg/cli/output"
	"github.com/radius-project/radius/pkg/cli/workspaces"
	"github.com/radius-project/radius/pkg/ucp/resources"
	resources_radius "github.com/radius-project/radius/pkg/ucp/resources/radius
)

// propertiesCmd represents the properties command
var propCmd = &cobra.Command{
	Use:   "props [resource ID]",
	Short: "List the properties of a resource",
	Args:  cobra.ExactArgs(1),
	RunE:  runProperties,
}

func runProperties(cmd *cobra.Command, args []string) error {
	resourceID := args[0]
	props, err := getResourceProperties(resourceID)
	if err != nil {
		if clierrors.IsExpectedError(err) {
			return clierrors.Message("The resource %q could not be found.", resourceID)
		}
		return fmt.Errorf("error retrieving properties for resource %q: %w", resourceID, err)
	}

	for key, value := range props {
		fmt.Printf("%s: %s\n", key, value)
	}
	return nil
}

var getResourceProperties = func(id string) (map[string]string, error) {
	// Implement the logic to fetch resource properties
	// Return clierrors.Message or clierrors.MessageWithCause for expected errors
	// Return the error as-is for unexpected errors
	fmt.Printf("hello from nj")
	return nil, nil
}
