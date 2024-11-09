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

package resourceproviders

import (
	"net/http"
	"testing"

	"github.com/radius-project/radius/pkg/ucp/frontend/api"
	"github.com/radius-project/radius/pkg/ucp/integrationtests/testserver"
	"github.com/stretchr/testify/require"
)

const (
	locationEmptyListResponseFixture = "testdata/location_v20231001preview_emptylist_responsebody.json"
	locationListResponseFixture      = "testdata/location_v20231001preview_list_responsebody.json"
)

func Test_Location_Lifecycle(t *testing.T) {
	server := testserver.StartWithETCD(t, api.DefaultModules)
	defer server.Close()

	createRadiusPlane(server)
	createResourceProvider(server)

	// We don't use t.Run() here because we want the test to fail if *any* of these steps fail.

	// List should start empty
	response := server.MakeRequest(http.MethodGet, locationCollectionURL, nil)
	response.EqualsFixture(200, locationEmptyListResponseFixture)

	// Getting a specific resource type should return 404 with the correct resource ID.
	response = server.MakeRequest(http.MethodGet, locationURL, nil)
	response.EqualsErrorCode(404, "NotFound")
	require.Equal(t, locationID, response.Error.Error.Target)

	// Create a resource provider
	createLocation(server)

	// List should now contain the resource provider
	response = server.MakeRequest(http.MethodGet, locationCollectionURL, nil)
	response.EqualsFixture(200, locationListResponseFixture)

	response = server.MakeRequest(http.MethodGet, locationURL, nil)
	response.EqualsFixture(200, locationResponseFixture)

	deleteLocation(server)
	deleteResourceProvider(server)
}