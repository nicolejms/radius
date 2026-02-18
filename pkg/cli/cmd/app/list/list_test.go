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
	"testing"

	"go.uber.org/mock/gomock"
	"github.com/radius-project/radius/pkg/cli/clients"
	"github.com/radius-project/radius/pkg/cli/connections"
	"github.com/radius-project/radius/pkg/cli/framework"
	"github.com/radius-project/radius/pkg/cli/objectformats"
	"github.com/radius-project/radius/pkg/cli/output"
	"github.com/radius-project/radius/pkg/cli/workspaces"
	"github.com/radius-project/radius/pkg/corerp/api/v20231001preview"
	"github.com/radius-project/radius/pkg/to"
	"github.com/radius-project/radius/test/radcli"
	"github.com/stretchr/testify/require"
)

func Test_CommandValidation(t *testing.T) {
	radcli.SharedCommandValidation(t, NewCommand)
}

func Test_Validate(t *testing.T) {
	configWithWorkspace := radcli.LoadConfigWithWorkspace(t)
	testcases := []radcli.ValidateInput{
		{
			Name:          "List Command with incorrect args",
			Input:         []string{"group"},
			ExpectedValid: false,
			ConfigHolder: framework.ConfigHolder{
				ConfigFilePath: "",
				Config:         configWithWorkspace,
			},
		},
		{
			Name:          "List Command with correct options but bad workspace",
			Input:         []string{"-w", "doesnotexist"},
			ExpectedValid: false,
			ConfigHolder: framework.ConfigHolder{
				ConfigFilePath: "",
				Config:         configWithWorkspace,
			},
		},
		{
			Name:          "List Command with valid workspace specified",
			Input:         []string{"-w", radcli.TestWorkspaceName},
			ExpectedValid: true,
			ConfigHolder: framework.ConfigHolder{
				ConfigFilePath: "",
				Config:         configWithWorkspace,
			},
		},
		{
			Name:          "List Command with fallback workspace",
			Input:         []string{"--group", "test-group"},
			ExpectedValid: true,
			ConfigHolder: framework.ConfigHolder{
				ConfigFilePath: "",
				Config:         radcli.LoadEmptyConfig(t),
			},
		},
	}
	radcli.SharedValidateValidation(t, NewCommand, testcases)
}

func Test_Run(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	applications := []v20231001preview.ApplicationResource{
		{
			Name: to.Ptr("A"),
		},
		{
			Name: to.Ptr("B"),
		},
	}

	appManagementClient := clients.NewMockApplicationsManagementClient(ctrl)
	appManagementClient.EXPECT().
		ListApplications(gomock.Any()).
		Return(applications, nil).
		Times(1)

	workspace := &workspaces.Workspace{
		Connection: map[string]any{
			"kind":    "kubernetes",
			"context": "kind-kind",
		},
		Name:  "kind-kind",
		Scope: "/planes/radius/local/resourceGroups/test-group",
	}
	outputSink := &output.MockOutput{}
	runner := &Runner{
		ConnectionFactory: &connections.MockFactory{ApplicationsManagementClient: appManagementClient},
		Workspace:         workspace,
		Format:            "table",
		Output:            outputSink,
	}

	err := runner.Run(context.Background())
	require.NoError(t, err)

	expected := []any{
		output.FormattedOutput{
			Format:  "table",
			Obj:     applications,
			Options: objectformats.GetResourceTableFormat(),
		},
	}

	require.Equal(t, expected, outputSink.Writes)
}

func Test_filterApplicationsByEnvironment(t *testing.T) {
testCases := []struct {
name        string
apps        []v20231001preview.ApplicationResource
environment string
expected    []v20231001preview.ApplicationResource
}{
{
name: "filter by environment name - single match",
apps: []v20231001preview.ApplicationResource{
{
Name: to.Ptr("app1"),
Properties: &v20231001preview.ApplicationProperties{
Environment: to.Ptr("/planes/radius/local/resourceGroups/test-group/providers/Applications.Core/environments/dev"),
},
},
{
Name: to.Ptr("app2"),
Properties: &v20231001preview.ApplicationProperties{
Environment: to.Ptr("/planes/radius/local/resourceGroups/test-group/providers/Applications.Core/environments/prod"),
},
},
},
environment: "dev",
expected: []v20231001preview.ApplicationResource{
{
Name: to.Ptr("app1"),
Properties: &v20231001preview.ApplicationProperties{
Environment: to.Ptr("/planes/radius/local/resourceGroups/test-group/providers/Applications.Core/environments/dev"),
},
},
},
},
{
name: "filter by full environment ID",
apps: []v20231001preview.ApplicationResource{
{
Name: to.Ptr("app1"),
Properties: &v20231001preview.ApplicationProperties{
Environment: to.Ptr("/planes/radius/local/resourceGroups/test-group/providers/Applications.Core/environments/dev"),
},
},
{
Name: to.Ptr("app2"),
Properties: &v20231001preview.ApplicationProperties{
Environment: to.Ptr("/planes/radius/local/resourceGroups/test-group/providers/Applications.Core/environments/prod"),
},
},
},
environment: "/planes/radius/local/resourceGroups/test-group/providers/Applications.Core/environments/dev",
expected: []v20231001preview.ApplicationResource{
{
Name: to.Ptr("app1"),
Properties: &v20231001preview.ApplicationProperties{
Environment: to.Ptr("/planes/radius/local/resourceGroups/test-group/providers/Applications.Core/environments/dev"),
},
},
},
},
{
name: "no matches",
apps: []v20231001preview.ApplicationResource{
{
Name: to.Ptr("app1"),
Properties: &v20231001preview.ApplicationProperties{
Environment: to.Ptr("/planes/radius/local/resourceGroups/test-group/providers/Applications.Core/environments/dev"),
},
},
},
environment: "staging",
expected:    []v20231001preview.ApplicationResource{},
},
{
name: "app with nil environment property",
apps: []v20231001preview.ApplicationResource{
{
Name:       to.Ptr("app1"),
Properties: &v20231001preview.ApplicationProperties{},
},
{
Name: to.Ptr("app2"),
Properties: &v20231001preview.ApplicationProperties{
Environment: to.Ptr("/planes/radius/local/resourceGroups/test-group/providers/Applications.Core/environments/dev"),
},
},
},
environment: "dev",
expected: []v20231001preview.ApplicationResource{
{
Name: to.Ptr("app2"),
Properties: &v20231001preview.ApplicationProperties{
Environment: to.Ptr("/planes/radius/local/resourceGroups/test-group/providers/Applications.Core/environments/dev"),
},
},
},
},
}

for _, tc := range testCases {
t.Run(tc.name, func(t *testing.T) {
result := filterApplicationsByEnvironment(tc.apps, tc.environment)
require.Equal(t, len(tc.expected), len(result), "number of filtered applications should match")
for i := range tc.expected {
require.Equal(t, *tc.expected[i].Name, *result[i].Name)
if tc.expected[i].Properties != nil && tc.expected[i].Properties.Environment != nil {
require.Equal(t, *tc.expected[i].Properties.Environment, *result[i].Properties.Environment)
}
}
})
}
}

func Test_matchesEnvironment(t *testing.T) {
testCases := []struct {
name             string
appEnvironment   string
filterEnvironment string
expected         bool
}{
{
name:             "exact match",
appEnvironment:   "/planes/radius/local/resourceGroups/test-group/providers/Applications.Core/environments/dev",
filterEnvironment: "/planes/radius/local/resourceGroups/test-group/providers/Applications.Core/environments/dev",
expected:         true,
},
{
name:             "match by environment name",
appEnvironment:   "/planes/radius/local/resourceGroups/test-group/providers/Applications.Core/environments/dev",
filterEnvironment: "dev",
expected:         true,
},
{
name:             "no match",
appEnvironment:   "/planes/radius/local/resourceGroups/test-group/providers/Applications.Core/environments/dev",
filterEnvironment: "prod",
expected:         false,
},
{
name:             "partial match but not at end",
appEnvironment:   "/planes/radius/local/resourceGroups/dev/providers/Applications.Core/environments/prod",
filterEnvironment: "dev",
expected:         false,
},
}

for _, tc := range testCases {
t.Run(tc.name, func(t *testing.T) {
result := matchesEnvironment(tc.appEnvironment, tc.filterEnvironment)
require.Equal(t, tc.expected, result)
})
}
}

func Test_Run_WithEnvironmentFiltering(t *testing.T) {
ctrl := gomock.NewController(t)
defer ctrl.Finish()

// All applications from different environments
allApplications := []v20231001preview.ApplicationResource{
{
Name: to.Ptr("app-dev"),
Properties: &v20231001preview.ApplicationProperties{
Environment: to.Ptr("/planes/radius/local/resourceGroups/test-group/providers/Applications.Core/environments/dev"),
},
},
{
Name: to.Ptr("app-prod"),
Properties: &v20231001preview.ApplicationProperties{
Environment: to.Ptr("/planes/radius/local/resourceGroups/test-group/providers/Applications.Core/environments/prod"),
},
},
}

// Expected filtered applications (only dev)
filteredApplications := []v20231001preview.ApplicationResource{
{
Name: to.Ptr("app-dev"),
Properties: &v20231001preview.ApplicationProperties{
Environment: to.Ptr("/planes/radius/local/resourceGroups/test-group/providers/Applications.Core/environments/dev"),
},
},
}

appManagementClient := clients.NewMockApplicationsManagementClient(ctrl)
appManagementClient.EXPECT().
ListApplications(gomock.Any()).
Return(allApplications, nil).
Times(1)

workspace := &workspaces.Workspace{
Connection: map[string]any{
"kind":    "kubernetes",
"context": "kind-kind",
},
Name:        "kind-kind",
Scope:       "/planes/radius/local/resourceGroups/test-group",
Environment: "dev", // Default environment is set
}
outputSink := &output.MockOutput{}
runner := &Runner{
ConnectionFactory: &connections.MockFactory{ApplicationsManagementClient: appManagementClient},
Workspace:         workspace,
Format:            "table",
Output:            outputSink,
}

err := runner.Run(context.Background())
require.NoError(t, err)

// Verify only filtered applications are output
expected := []any{
output.FormattedOutput{
Format:  "table",
Obj:     filteredApplications,
Options: objectformats.GetResourceTableFormat(),
},
}

require.Equal(t, expected, outputSink.Writes)
}
