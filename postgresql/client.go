package postgresql

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//
// See the License for the specific language governing permissions and
// limitations under the License.

// Note: PostgreSQL service is in preview and APIs have not yet been documented.
// This sample is based on the existing SQL Server client as a starting point and
// should NOT be counted as a stable API.

import (
	"github.com/Azure/go-autorest/autorest"
)

const (
	// DefaultBaseURI is the default URI used for the service Sql
	DefaultBaseURI = "https://management.azure.com"
)

// ManagementClient is the base for PostgreSQL client
type ManagementClient struct {
	autorest.Client
	BaseURI        string
	SubscriptionID string
}

// NewWithBaseURI creates an instance of the ManagementClient client.
func NewWithBaseURI(baseURI string, subscriptionID string) ManagementClient {
	return ManagementClient{
		Client:         autorest.NewClientWithUserAgent(UserAgent()),
		BaseURI:        baseURI,
		SubscriptionID: subscriptionID,
	}
}

// UserAgent returns the UserAgent string to use when sending http.Requests.
// from sql/version.go
func UserAgent() string {
	return "Azure-SDK-For-Go/sample postgresql/"
}

// Version returns the semantic version (see http://semver.org) of the client.
// from sql/version.go
func Version() string {
	return "sample"
}
