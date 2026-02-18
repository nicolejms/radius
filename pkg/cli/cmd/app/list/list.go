/*
Copyright 2023 The Radius Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package list

import (
	"context"

	"github.com/radius-project/radius/pkg/cli"
	"github.com/radius-project/radius/pkg/cli/cmd/commonflags"
	"github.com/radius-project/radius/pkg/cli/connections"
	"github.com/radius-project/radius/pkg/cli/framework"
	"github.com/radius-project/radius/pkg/cli/objectformats"
	"github.com/radius-project/radius/pkg/cli/output"
	"github.com/radius-project/radius/pkg/cli/workspaces"
	corerp "github.com/radius-project/radius/pkg/corerp/api/v20231001preview"
	"github.com/spf13/cobra"
)

// NewCommand creates an instance of the `rad app list` command and runner.
//

// NewCommand creates a new Cobra command and a new Runner, and configures the command with flags and usage information
// to list Radius Applications in a resource group associated with the default environment.
func NewCommand(factory framework.Factory) (*cobra.Command, framework.Runner) {
	runner := NewRunner(factory)

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List Radius Applications",
		Long:  `Lists Radius Applications deployed in the resource group associated with the default environment`,
		Args:  cobra.NoArgs,
		Example: `
# List applications
rad app list

# List applications in a specific resource group
rad app list --group my-group
`,
		RunE: framework.RunCommand(runner),
	}

	commonflags.AddWorkspaceFlag(cmd)
	commonflags.AddResourceGroupFlag(cmd)
	commonflags.AddOutputFlag(cmd)

	return cmd, runner
}

// Runner is the Runner implementation for the `rad app list` command.
type Runner struct {
	ConfigHolder      *framework.ConfigHolder
	ConnectionFactory connections.Factory
	Workspace         *workspaces.Workspace
	Output            output.Interface

	Format string
}

// NewRunner creates an instance of the runner for the `rad app list` command.
func NewRunner(factory framework.Factory) *Runner {
	return &Runner{
		ConnectionFactory: factory.GetConnectionFactory(),
		ConfigHolder:      factory.GetConfigHolder(),
		Output:            factory.GetOutput(),
	}
}

// Validate runs validation for the `rad app list` command.
//

// Validate checks the workspace, scope, and output format of the command and sets them in the Runner struct,
// returning an error if any of these checks fail.
func (r *Runner) Validate(cmd *cobra.Command, args []string) error {
	workspace, err := cli.RequireWorkspace(cmd, r.ConfigHolder.Config, r.ConfigHolder.DirectoryConfig)
	if err != nil {
		return err
	}
	r.Workspace = workspace

	// Allow '--group' to override scope
	scope, err := cli.RequireScope(cmd, *r.Workspace)
	if err != nil {
		return err
	}
	r.Workspace.Scope = scope

	format, err := cli.RequireOutput(cmd)
	if err != nil {
		return err
	}

	r.Format = format

	return nil
}

// Run runs the `rad app list` command.
//

// Run() creates an ApplicationsManagementClient using the provided ConnectionFactory, then lists the applications,
// filters them by the default environment if one is set, and writes the output in the specified format,
// returning an error if any of these steps fail.
func (r *Runner) Run(ctx context.Context) error {
	client, err := r.ConnectionFactory.CreateApplicationsManagementClient(ctx, *r.Workspace)
	if err != nil {
		return err
	}

	apps, err := client.ListApplications(ctx)
	if err != nil {
		return err
	}

	// Filter applications by the default environment if one is configured in the workspace
	if r.Workspace.Environment != "" {
		apps = filterApplicationsByEnvironment(apps, r.Workspace.Environment)
	}

	return r.Output.WriteFormatted(r.Format, apps, objectformats.GetResourceTableFormat())
}

// filterApplicationsByEnvironment filters the list of applications to only include those
// associated with the specified environment. The environment can be either a resource ID
// or just an environment name.
func filterApplicationsByEnvironment(apps []corerp.ApplicationResource, environment string) []corerp.ApplicationResource {
	var filtered []corerp.ApplicationResource
	
	for _, app := range apps {
		if app.Properties != nil && app.Properties.Environment != nil {
			// Match either by full resource ID or by environment name at the end of the ID
			if matchesEnvironment(*app.Properties.Environment, environment) {
				filtered = append(filtered, app)
			}
		}
	}
	
	return filtered
}

// matchesEnvironment checks if the application's environment matches the filter environment.
// It handles both full resource IDs and simple environment names.
func matchesEnvironment(appEnvironment, filterEnvironment string) bool {
	// Direct match (full resource ID)
	if appEnvironment == filterEnvironment {
		return true
	}
	
	// Check if the filter environment appears at the end of the app's environment resource ID
	// This handles cases where the workspace has just the name but the app has the full ID
	if len(appEnvironment) > len(filterEnvironment) {
		// Check if it ends with /environments/{name}
		suffix := "/" + filterEnvironment
		if len(appEnvironment) >= len(suffix) && appEnvironment[len(appEnvironment)-len(suffix):] == suffix {
			return true
		}
	}
	
	// Check if the app environment ends with the filter (for cases where filter is a full ID but app is shorter)
	if len(filterEnvironment) > len(appEnvironment) {
		suffix := "/" + appEnvironment
		if len(filterEnvironment) >= len(suffix) && filterEnvironment[len(filterEnvironment)-len(suffix):] == suffix {
			return true
		}
	}
	
	return false
}
