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
	"fmt"
	"sort"
	"strings"

	"github.com/radius-project/radius/pkg/corerp/api/v20231001preview"
	"github.com/radius-project/radius/pkg/ucp/resources"
)

// displayMermaid builds a Mermaid diagram representation of the application graph.
// The output can be used directly in GitHub Markdown files for visualization.
func displayMermaid(applicationResources []*v20231001preview.ApplicationGraphResource, applicationName string) string {
	output := &strings.Builder{}
	
	// Write Mermaid header and application title
	output.WriteString("```mermaid\n")
	output.WriteString("graph TB\n")
	output.WriteString(fmt.Sprintf("    %s[\"Application: %s\"]\n", sanitizeNodeID("app"), applicationName))
	output.WriteString("\n")
	
	if len(applicationResources) == 0 {
		output.WriteString("    empty[\"(empty)\"]\n")
		output.WriteString(fmt.Sprintf("    %s --> empty\n", sanitizeNodeID("app")))
		output.WriteString("```\n")
		return output.String()
	}
	
	// Sort resources for consistent output
	containerType := "Applications.Core/containers"
	sort.Slice(applicationResources, func(i, j int) bool {
		if strings.EqualFold(*applicationResources[i].Type, containerType) !=
			strings.EqualFold(*applicationResources[j].Type, containerType) {
			return strings.EqualFold(*applicationResources[i].Type, containerType)
		}
		if *applicationResources[i].Type != *applicationResources[j].Type {
			return *applicationResources[i].Type < *applicationResources[j].Type
		}
		if *applicationResources[i].Name != *applicationResources[j].Name {
			return *applicationResources[i].Name < *applicationResources[j].Name
		}
		return *applicationResources[i].ID < *applicationResources[j].ID
	})
	
	// Track which nodes we've already defined
	definedNodes := make(map[string]bool)
	
	// Define all resource nodes with their types
	for _, resource := range applicationResources {
		nodeID := sanitizeNodeID(*resource.Name)
		if !definedNodes[nodeID] {
			// Use different shapes based on resource type
			shape := getNodeShape(*resource.Type)
			label := fmt.Sprintf("%s<br/>%s", *resource.Name, shortenType(*resource.Type))
			output.WriteString(fmt.Sprintf("    %s%s%s%s\n", nodeID, shape.open, label, shape.close))
			definedNodes[nodeID] = true
		}
	}
	
	output.WriteString("\n")
	
	// Track connections to avoid duplicates
	connections := make(map[string]bool)
	
	// Add connections between resources
	for _, resource := range applicationResources {
		sourceNodeID := sanitizeNodeID(*resource.Name)
		
		for _, connection := range resource.Connections {
			connectionID, err := resources.Parse(*connection.ID)
			if err != nil {
				continue
			}
			
			targetNodeID := sanitizeNodeID(connectionID.Name())
			
			// Create connection string based on direction
			var connStr string
			if *connection.Direction == v20231001preview.DirectionOutbound {
				connStr = fmt.Sprintf("    %s --> %s\n", sourceNodeID, targetNodeID)
			} else {
				// For inbound, we skip because outbound will handle it
				continue
			}
			
			// Add connection if not already added
			if !connections[connStr] {
				output.WriteString(connStr)
				connections[connStr] = true
			}
		}
	}
	
	output.WriteString("\n")
	
	// Add styling for different resource types
	output.WriteString("    classDef container fill:#e1f5ff,stroke:#0078d4,stroke-width:2px\n")
	output.WriteString("    classDef datastore fill:#fff4ce,stroke:#d83b01,stroke-width:2px\n")
	output.WriteString("    classDef gateway fill:#f3e5f5,stroke:#6a1b9a,stroke-width:2px\n")
	output.WriteString("    classDef default fill:#e8eaf6,stroke:#3f51b5,stroke-width:2px\n")
	
	output.WriteString("\n")
	
	// Apply classes to nodes based on their types
	for _, resource := range applicationResources {
		nodeID := sanitizeNodeID(*resource.Name)
		class := getNodeClass(*resource.Type)
		if class != "" {
			output.WriteString(fmt.Sprintf("    class %s %s\n", nodeID, class))
		}
	}
	
	output.WriteString("```\n")
	return output.String()
}

// nodeShape represents the opening and closing characters for a Mermaid node shape
type nodeShape struct {
	open  string
	close string
}

// getNodeShape returns the appropriate Mermaid shape based on resource type
func getNodeShape(resourceType string) nodeShape {
	switch {
	case strings.Contains(resourceType, "containers"):
		return nodeShape{"[\"", "\"]"} // Rectangle
	case strings.Contains(resourceType, "Datastores"):
		return nodeShape{"[(\"", "\")]"} // Cylinder (database)
	case strings.Contains(resourceType, "gateway"):
		return nodeShape{"{\"", "\"}"} // Diamond
	default:
		return nodeShape{"[\"", "\"]"} // Rectangle
	}
}

// getNodeClass returns the CSS class to apply based on resource type
func getNodeClass(resourceType string) string {
	switch {
	case strings.Contains(resourceType, "containers"):
		return "container"
	case strings.Contains(resourceType, "Datastores"):
		return "datastore"
	case strings.Contains(resourceType, "gateway"):
		return "gateway"
	default:
		return "default"
	}
}

// sanitizeNodeID converts a resource name to a valid Mermaid node ID
func sanitizeNodeID(name string) string {
	// Replace characters that might cause issues in Mermaid with underscores
	replacer := strings.NewReplacer(
		"-", "_",
		".", "_",
		" ", "_",
		"/", "_",
		":", "_",
	)
	return replacer.Replace(name)
}

// shortenType shortens a resource type for display in the diagram
func shortenType(fullType string) string {
	parts := strings.Split(fullType, "/")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return fullType
}
