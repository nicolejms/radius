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

package graph

import (
	"strings"
	"testing"

	corerpv20231001preview "github.com/radius-project/radius/pkg/corerp/api/v20231001preview"
	"github.com/stretchr/testify/require"
)

func Test_displayMermaid(t *testing.T) {
	t.Run("empty graph", func(t *testing.T) {
		graph := []*corerpv20231001preview.ApplicationGraphResource{}
		actual := displayMermaid(graph, "cool-app")
		
		require.Contains(t, actual, "```mermaid")
		require.Contains(t, actual, "graph TB")
		require.Contains(t, actual, "Application: cool-app")
		require.Contains(t, actual, "(empty)")
		require.Contains(t, actual, "```")
	})

	t.Run("simple application with two resources", func(t *testing.T) {
		backendID := "/planes/radius/local/resourcegroups/default/providers/Applications.Core/containers/backend"
		backendType := "Applications.Core/containers"
		backendName := "backend"

		redisID := "/planes/radius/local/resourcegroups/default/providers/Applications.Datastores/redisCaches/redis"
		redisName := "redis"
		redisType := "Applications.Datastores/redisCaches"

		provisioningStateSuccess := "Succeeded"
		dirOutbound := corerpv20231001preview.DirectionOutbound
		dirInbound := corerpv20231001preview.DirectionInbound

		graph := []*corerpv20231001preview.ApplicationGraphResource{
			{
				ID:                &backendID,
				Name:              &backendName,
				Type:              &backendType,
				ProvisioningState: &provisioningStateSuccess,
				OutputResources:   []*corerpv20231001preview.ApplicationGraphOutputResource{},
				Connections: []*corerpv20231001preview.ApplicationGraphConnection{
					{
						Direction: &dirOutbound,
						ID:        &redisID,
					},
				},
			},
			{
				ID:                &redisID,
				Name:              &redisName,
				Type:              &redisType,
				ProvisioningState: &provisioningStateSuccess,
				OutputResources:   []*corerpv20231001preview.ApplicationGraphOutputResource{},
				Connections: []*corerpv20231001preview.ApplicationGraphConnection{
					{
						Direction: &dirInbound,
						ID:        &backendID,
					},
				},
			},
		}

		actual := displayMermaid(graph, "test-app")

		// Check for Mermaid syntax
		require.Contains(t, actual, "```mermaid")
		require.Contains(t, actual, "graph TB")
		require.Contains(t, actual, "```")

		// Check for application name
		require.Contains(t, actual, "Application: test-app")

		// Check for resource nodes
		require.Contains(t, actual, "backend")
		require.Contains(t, actual, "redis")

		// Check for resource types in labels
		require.Contains(t, actual, "containers")
		require.Contains(t, actual, "redisCaches")

		// Check for connection
		require.Contains(t, actual, "backend --> redis")

		// Check for styling
		require.Contains(t, actual, "classDef container")
		require.Contains(t, actual, "classDef datastore")
	})

	t.Run("complex application with multiple resources and connections", func(t *testing.T) {
		backendID := "/planes/radius/local/resourcegroups/default/providers/Applications.Core/containers/backend"
		backendType := "Applications.Core/containers"
		backendName := "backend"

		frontendID := "/planes/radius/local/resourcegroups/default/providers/Applications.Core/containers/frontend"
		frontendName := "frontend"
		frontendType := "Applications.Core/containers"

		sqlDbID := "/planes/radius/local/resourcegroups/default/providers/Applications.Datastores/sqlDatabases/sql-db"
		sqlDbName := "sql-db"
		sqlDbType := "Applications.Datastores/sqlDatabases"

		provisioningStateSuccess := "Succeeded"
		dirOutbound := corerpv20231001preview.DirectionOutbound

		graph := []*corerpv20231001preview.ApplicationGraphResource{
			{
				ID:                &frontendID,
				Name:              &frontendName,
				Type:              &frontendType,
				ProvisioningState: &provisioningStateSuccess,
				OutputResources:   []*corerpv20231001preview.ApplicationGraphOutputResource{},
				Connections: []*corerpv20231001preview.ApplicationGraphConnection{
					{
						Direction: &dirOutbound,
						ID:        &backendID,
					},
				},
			},
			{
				ID:                &backendID,
				Name:              &backendName,
				Type:              &backendType,
				ProvisioningState: &provisioningStateSuccess,
				OutputResources:   []*corerpv20231001preview.ApplicationGraphOutputResource{},
				Connections: []*corerpv20231001preview.ApplicationGraphConnection{
					{
						Direction: &dirOutbound,
						ID:        &sqlDbID,
					},
				},
			},
			{
				ID:                &sqlDbID,
				Name:              &sqlDbName,
				Type:              &sqlDbType,
				ProvisioningState: &provisioningStateSuccess,
				OutputResources:   []*corerpv20231001preview.ApplicationGraphOutputResource{},
			},
		}

		actual := displayMermaid(graph, "complex-app")

		// Check for all resources
		require.Contains(t, actual, "frontend")
		require.Contains(t, actual, "backend")
		require.Contains(t, actual, "sql_db")

		// Check for connections
		require.Contains(t, actual, "frontend --> backend")
		require.Contains(t, actual, "backend --> sql_db")

		// Verify proper Mermaid structure
		require.True(t, strings.HasPrefix(actual, "```mermaid\n"))
		require.True(t, strings.HasSuffix(actual, "```\n"))
	})
}

func Test_sanitizeNodeID(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple name",
			input:    "backend",
			expected: "backend",
		},
		{
			name:     "name with hyphens",
			input:    "sql-db",
			expected: "sql_db",
		},
		{
			name:     "name with dots",
			input:    "my.service",
			expected: "my_service",
		},
		{
			name:     "name with spaces",
			input:    "my service",
			expected: "my_service",
		},
		{
			name:     "name with slashes",
			input:    "app/service",
			expected: "app_service",
		},
		{
			name:     "name with colons",
			input:    "app:service",
			expected: "app_service",
		},
		{
			name:     "complex name",
			input:    "my-app.service/v1:latest",
			expected: "my_app_service_v1_latest",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := sanitizeNodeID(tt.input)
			require.Equal(t, tt.expected, actual)
		})
	}
}

func Test_shortenType(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple type",
			input:    "containers",
			expected: "containers",
		},
		{
			name:     "namespaced type",
			input:    "Applications.Core/containers",
			expected: "containers",
		},
		{
			name:     "deep namespaced type",
			input:    "Applications.Datastores/sqlDatabases",
			expected: "sqlDatabases",
		},
		{
			name:     "multiple slashes",
			input:    "a/b/c/d",
			expected: "d",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := shortenType(tt.input)
			require.Equal(t, tt.expected, actual)
		})
	}
}

func Test_getNodeShape(t *testing.T) {
	tests := []struct {
		name         string
		resourceType string
		expectedOpen string
	}{
		{
			name:         "container type",
			resourceType: "Applications.Core/containers",
			expectedOpen: "[\"",
		},
		{
			name:         "datastore type",
			resourceType: "Applications.Datastores/sqlDatabases",
			expectedOpen: "[(\"",
		},
		{
			name:         "gateway type",
			resourceType: "Applications.Core/gateway",
			expectedOpen: "{\"",
		},
		{
			name:         "default type",
			resourceType: "Applications.Core/unknown",
			expectedOpen: "[\"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			shape := getNodeShape(tt.resourceType)
			require.Equal(t, tt.expectedOpen, shape.open)
		})
	}
}

func Test_getNodeClass(t *testing.T) {
	tests := []struct {
		name         string
		resourceType string
		expected     string
	}{
		{
			name:         "container type",
			resourceType: "Applications.Core/containers",
			expected:     "container",
		},
		{
			name:         "datastore type",
			resourceType: "Applications.Datastores/sqlDatabases",
			expected:     "datastore",
		},
		{
			name:         "gateway type",
			resourceType: "Applications.Core/gateway",
			expected:     "gateway",
		},
		{
			name:         "default type",
			resourceType: "Applications.Core/unknown",
			expected:     "default",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			class := getNodeClass(tt.resourceType)
			require.Equal(t, tt.expected, class)
		})
	}
}
