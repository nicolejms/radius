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
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	v1 "github.com/radius-project/radius/pkg/armrpc/api/v1"
	"github.com/radius-project/radius/pkg/corerp/api/v20231001preview"
)

const (
	fakeServerNotFoundResponse = "unexpected status code 404. acceptable values are http.StatusOK"
)

// Is404Error returns true if the error is a 404 payload from an autorest operation.
//

// "Is404Error" checks if the given error is a 404 error by checking if it is one of:
// a ResponseError with an ErrorCode of "NotFound", or
// a ResponseError with a StatusCode of 404, or
// an ErrorResponse with an Error Code of "NotFound".
func Is404Error(err error) bool {
	if err == nil {
		return false
	}

	// NotFound Response from Fake Server - used for testing
	if strings.Contains(err.Error(), fakeServerNotFoundResponse) {
		return true
	}

	// The error might already be an ResponseError
	responseError := &azcore.ResponseError{}
	if errors.As(err, &responseError) && responseError.ErrorCode == v1.CodeNotFound || responseError.StatusCode == http.StatusNotFound {
		return true
	} else if errors.As(err, &responseError) {
		return false
	}

	// OK so it's not an ResponseError, can we turn it into an ErrorResponse?
	errorResponse := v20231001preview.ErrorResponse{}
	marshallErr := json.Unmarshal([]byte(err.Error()), &errorResponse)
	if marshallErr != nil {
		return false
	}

	if errorResponse.Error != nil && *errorResponse.Error.Code == v1.CodeNotFound {
		return true
	}

	return false
}

// ConvertAzureErrorResponse converts Azure SDK's ErrorResponse to Radius ErrorDetails.
// This function handles nested error details recursively.
func ConvertAzureErrorResponse(azErr *armresources.ErrorResponse) *v1.ErrorDetails {
	if azErr == nil {
		return nil
	}

	errDetails := &v1.ErrorDetails{}

	// Convert basic fields
	if azErr.Code != nil {
		errDetails.Code = *azErr.Code
	}
	if azErr.Message != nil {
		errDetails.Message = *azErr.Message
	}
	if azErr.Target != nil {
		errDetails.Target = *azErr.Target
	}

	// Convert additional info if present
	if len(azErr.AdditionalInfo) > 0 {
		errDetails.AdditionalInfo = make([]*v1.ErrorAdditionalInfo, len(azErr.AdditionalInfo))
		for i, info := range azErr.AdditionalInfo {
			additionalInfo := &v1.ErrorAdditionalInfo{}
			if info.Type != nil {
				additionalInfo.Type = *info.Type
			}
			if info.Info != nil {
				// Info is an 'any' type from Azure SDK, we need to handle it carefully
				if infoMap, ok := info.Info.(map[string]any); ok {
					additionalInfo.Info = infoMap
				}
			}
			errDetails.AdditionalInfo[i] = additionalInfo
		}
	}

	// Recursively convert nested error details
	if len(azErr.Details) > 0 {
		errDetails.Details = make([]*v1.ErrorDetails, len(azErr.Details))
		for i, detail := range azErr.Details {
			errDetails.Details[i] = ConvertAzureErrorResponse(detail)
		}
	}

	return errDetails
}
