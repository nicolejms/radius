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

package clients

import (
	"errors"
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	v1 "github.com/radius-project/radius/pkg/armrpc/api/v1"
	"github.com/radius-project/radius/pkg/to"
	"github.com/stretchr/testify/require"
)

func TestIs404Error(t *testing.T) {
	var err error

	// Test with a ResponseError with an ErrorCode of "NotFound"
	err = &azcore.ResponseError{ErrorCode: v1.CodeNotFound}
	if !Is404Error(err) {
		t.Errorf("Expected Is404Error to return true for ResponseError with ErrorCode of 'NotFound', but it returned false")
	}

	// Test with a ResponseError with a StatusCode of 404
	err = &azcore.ResponseError{StatusCode: http.StatusNotFound}
	if !Is404Error(err) {
		t.Errorf("Expected Is404Error to return true for ResponseError with StatusCode of 404, but it returned false")
	}

	// Test with an ErrorResponse with an Error Code of "NotFound"
	err = errors.New(`{"error": {"code": "NotFound"}}`)
	if !Is404Error(err) {
		t.Errorf("Expected Is404Error to return true for ErrorResponse with Error Code of 'NotFound', but it returned false")
	}

	// Test with an ErrorResponse with a different Error Code
	err = errors.New(`{"error": {"code": "SomeOtherCode"}}`)
	if Is404Error(err) {
		t.Errorf("Expected Is404Error to return false for ErrorResponse with Error Code of 'SomeOtherCode', but it returned true")
	}

	// Test with a different error type
	err = errors.New("Some other error")
	if Is404Error(err) {
		t.Errorf("Expected Is404Error to return false for error of type %T, but it returned true", err)
	}

	// Test with a nil error
	if Is404Error(nil) {
		t.Errorf("Expected Is404Error to return false for nil error, but it returned true")
	}

	// Test with a fake server not found response
	err = errors.New(fakeServerNotFoundResponse)
	if !Is404Error(err) {
		t.Errorf("Expected Is404Error to return true for fake server not found response, but it returned false")
	}
}

func TestConvertAzureErrorResponse(t *testing.T) {
	t.Run("nil error response", func(t *testing.T) {
		result := ConvertAzureErrorResponse(nil)
		require.Nil(t, result)
	})

	t.Run("simple error without nested details", func(t *testing.T) {
		azErr := &armresources.ErrorResponse{
			Code:    to.Ptr("DeploymentFailed"),
			Message: to.Ptr("The deployment failed"),
			Target:  to.Ptr("resource"),
		}

		result := ConvertAzureErrorResponse(azErr)
		require.NotNil(t, result)
		require.Equal(t, "DeploymentFailed", result.Code)
		require.Equal(t, "The deployment failed", result.Message)
		require.Equal(t, "resource", result.Target)
		require.Empty(t, result.Details)
		require.Empty(t, result.AdditionalInfo)
	})

	t.Run("error with nested details", func(t *testing.T) {
		azErr := &armresources.ErrorResponse{
			Code:    to.Ptr("DeploymentFailed"),
			Message: to.Ptr("The deployment failed"),
			Details: []*armresources.ErrorResponse{
				{
					Code:    to.Ptr("InvalidTemplate"),
					Message: to.Ptr("Template validation failed"),
				},
				{
					Code:    to.Ptr("ResourceNotFound"),
					Message: to.Ptr("Resource does not exist"),
				},
			},
		}

		result := ConvertAzureErrorResponse(azErr)
		require.NotNil(t, result)
		require.Equal(t, "DeploymentFailed", result.Code)
		require.Equal(t, "The deployment failed", result.Message)
		require.Len(t, result.Details, 2)
		require.Equal(t, "InvalidTemplate", result.Details[0].Code)
		require.Equal(t, "Template validation failed", result.Details[0].Message)
		require.Equal(t, "ResourceNotFound", result.Details[1].Code)
		require.Equal(t, "Resource does not exist", result.Details[1].Message)
	})

	t.Run("error with additional info", func(t *testing.T) {
		azErr := &armresources.ErrorResponse{
			Code:    to.Ptr("DeploymentFailed"),
			Message: to.Ptr("The deployment failed"),
			AdditionalInfo: []*armresources.ErrorAdditionalInfo{
				{
					Type: to.Ptr("PolicyViolation"),
					Info: map[string]any{
						"policyName": "RequiredTags",
						"severity":   "High",
					},
				},
			},
		}

		result := ConvertAzureErrorResponse(azErr)
		require.NotNil(t, result)
		require.Equal(t, "DeploymentFailed", result.Code)
		require.Len(t, result.AdditionalInfo, 1)
		require.Equal(t, "PolicyViolation", result.AdditionalInfo[0].Type)
		require.Equal(t, "RequiredTags", result.AdditionalInfo[0].Info["policyName"])
		require.Equal(t, "High", result.AdditionalInfo[0].Info["severity"])
	})

	t.Run("deeply nested error details", func(t *testing.T) {
		azErr := &armresources.ErrorResponse{
			Code:    to.Ptr("DeploymentFailed"),
			Message: to.Ptr("The deployment failed"),
			Details: []*armresources.ErrorResponse{
				{
					Code:    to.Ptr("InvalidTemplate"),
					Message: to.Ptr("Template validation failed"),
					Details: []*armresources.ErrorResponse{
						{
							Code:    to.Ptr("MissingParameter"),
							Message: to.Ptr("Required parameter not provided"),
						},
					},
				},
			},
		}

		result := ConvertAzureErrorResponse(azErr)
		require.NotNil(t, result)
		require.Equal(t, "DeploymentFailed", result.Code)
		require.Len(t, result.Details, 1)
		require.Equal(t, "InvalidTemplate", result.Details[0].Code)
		require.Len(t, result.Details[0].Details, 1)
		require.Equal(t, "MissingParameter", result.Details[0].Details[0].Code)
		require.Equal(t, "Required parameter not provided", result.Details[0].Details[0].Message)
	})

	t.Run("error with all fields populated", func(t *testing.T) {
		azErr := &armresources.ErrorResponse{
			Code:    to.Ptr("DeploymentFailed"),
			Message: to.Ptr("The deployment failed"),
			Target:  to.Ptr("Microsoft.Resources/deployments"),
			AdditionalInfo: []*armresources.ErrorAdditionalInfo{
				{
					Type: to.Ptr("PolicyViolation"),
					Info: map[string]any{
						"policyName": "RequiredTags",
					},
				},
			},
			Details: []*armresources.ErrorResponse{
				{
					Code:    to.Ptr("InvalidTemplate"),
					Message: to.Ptr("Template validation failed"),
					Target:  to.Ptr("template"),
				},
			},
		}

		result := ConvertAzureErrorResponse(azErr)
		require.NotNil(t, result)
		require.Equal(t, "DeploymentFailed", result.Code)
		require.Equal(t, "The deployment failed", result.Message)
		require.Equal(t, "Microsoft.Resources/deployments", result.Target)
		require.Len(t, result.AdditionalInfo, 1)
		require.Len(t, result.Details, 1)
		require.Equal(t, "InvalidTemplate", result.Details[0].Code)
		require.Equal(t, "template", result.Details[0].Target)
	})
}
