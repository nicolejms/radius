package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// propertiesCmd represents the properties command
var propertiesCmd = &cobra.Command{
	Use:   "properties [resource ID]",
	Short: "List the properties of a resource",
	Args:  cobra.ExactArgs(1),
	RunE:  runProperties,
}

func runProperties(cmd *cobra.Command, args []string) error {
	resourceID := args[0]
	props, err := getResourceProperties(resourceID)
	if err != nil {
		return fmt.Errorf("error retrieving properties for resource %q: %w", resourceID, err)
	}

	for key, value := range props {
		fmt.Printf("%s: %s\n", key, value)
	}
	return nil
}

func getResourceProperties(id string) (map[string]string, error) {
	// Implement the logic to fetch resource properties
	// Return clierrors.Message or clierrors.MessageWithCause for expected errors
	// Return the error as-is for unexpected errors
	fmt.Print("hello!")
	return nil, nil
}