package postgresql

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/arm/resources/resources"
	"github.com/Azure/go-autorest/autorest/to"
)

// Copyright (c) Microsoft and contributors.  All rights reserved.
//
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

const (
	namespace    = "Microsoft.DBforPostgreSQL"
	resourceType = "servers"
)

// ServersClient provides access to server lifecycle
type ServersClient struct {
	ManagementClient
}

// NewServersClient creates an instance of the ServersClient client.
func NewServersClient(subscriptionID string) ServersClient {
	return NewServersClientWithBaseURI(DefaultBaseURI, subscriptionID)
}

// NewServersClientWithBaseURI creates an instance of the ServersClient client.
func NewServersClientWithBaseURI(baseURI string, subscriptionID string) ServersClient {
	return ServersClient{NewWithBaseURI(baseURI, subscriptionID)}
}

// DeleteServer deletes a PostgreSQL server
// Note this creates a resourceGroup client for generic operations.
func (serversClient ServersClient) DeleteServer(resourceGroupName string, serverName string) (err error) {
	groupClient := resources.NewGroupClient(serversClient.SubscriptionID)
	groupClient.Authorizer = serversClient.Authorizer
	_, errorChannel := groupClient.Delete(resourceGroupName, namespace, "", resourceType, serverName, nil)
	return <-errorChannel
}

// ChangeAdministratorPassword updates the admin password for the server
func (serversClient ServersClient) ChangeAdministratorPassword(resourceGroupName string, serverName string, newPassword string) (err error) {
	// use the GroupClient from the rest_operations.go
	groupClient := NewGroupClient(serversClient.SubscriptionID)
	groupClient.Authorizer = serversClient.Authorizer
	// CreateOrUpdate(resourceGroupName string, resourceProviderNamespace string, parentResourcePath string, resourceType string, resourceName string, parameters resources.GenericResource, cancel <-chan struct{}, useHTTPPatch bool) (<-chan resources.GenericResource, <-chan error)
	parameters := resources.GenericResource{
		Properties: &map[string]interface{}{
			"administratorLoginPassword": newPassword,
		},
	}
	_, errorChannel := groupClient.CreateOrUpdate(resourceGroupName, namespace, "", resourceType, serverName, parameters, nil, true)
	return <-errorChannel
}

// CreateServer creates a new server instance
func (serversClient ServersClient) CreateServer(resourceGroupName string, serverName string, serverProperties ServerProperties) (err error) {
	// use the GroupClient from the rest_operations.go
	groupClient := NewGroupClient(serversClient.SubscriptionID)
	groupClient.Authorizer = serversClient.Authorizer
	// CreateOrUpdate(resourceGroupName string, resourceProviderNamespace string, parentResourcePath string, resourceType string, resourceName string, parameters resources.GenericResource, cancel <-chan struct{}, useHTTPPatch bool) (<-chan resources.GenericResource, <-chan error)

	serverPropertyErrors := validateServerProperties(serverProperties)
	if serverPropertyErrors != nil {
		return serverPropertyErrors
	}
	// crete the sku
	sku := &resources.Sku{
		Name: to.StringPtr("SkuName"),
		Tier: to.StringPtr(string(serverProperties.Tier)),
	}
	// add compute unit/capacity if set
	if serverProperties.ComputeUnits > 0 {
		sku.Capacity = to.Int32Ptr(serverProperties.ComputeUnits)
	}
	// create the properties
	properties := map[string]interface{}{
		"location":                   serverProperties.Location,
		"administratorLogin":         serverProperties.AdministratorLogin,
		"administratorLoginPassword": serverProperties.AdministratorLoginPassword,
		"version":                    serverProperties.Version,
		"sslEnforcement":             "Enabled",
	}
	// add storage if not default
	if serverProperties.StorageMB > 0 {
		//GB to MB
		properties["storageMB"] = to.Int32Ptr(serverProperties.StorageMB * 1024)
	}
	// create the resource
	serverResource := resources.GenericResource{
		Location:   to.StringPtr(serverProperties.Location),
		Properties: &properties,
		Sku:        sku,
	}
	_, errorChannel := groupClient.CreateOrUpdate(resourceGroupName, namespace, "", resourceType, serverName, serverResource, nil, false)
	return <-errorChannel
}

// RestoreServer creates a new server as point in time copy of source server
// {
//   "createMode": "PointInTimeRestore",
//   "sourceServerId": "/subscriptions/xxx/rg/providers/Microsoft.DBforMySQL/servers/sourceDb",
//   "restorePointInTime": "2017-05-22T05:01:02.344444Z"},
//   "location": "westeurope"
// }
func (serversClient ServersClient) RestoreServer(
	srcResourceGroup string,
	srcServerName string,
	targetResourceGroup string,
	targetServerName string,
	restorePointInTime time.Time,
) (err error) {
	// get the source server
	//resourceClient := resources.NewGroupClient(serversClient.SubscriptionID)
	//resourceClient.Authorizer = serversClient.Authorizer
	// use the GroupClient from the rest_operations.go
	groupClient := NewGroupClient(serversClient.SubscriptionID)
	groupClient.Authorizer = serversClient.Authorizer

	// resourceGroupName string, resourceProviderNamespace string, parentResourcePath string, resourceType string, resourceName string
	srcServer, resourceError := groupClient.Get(srcResourceGroup, namespace, "", resourceType, srcServerName)

	if resourceError != nil {
		fmt.Printf("Error getting source server %s", resourceError)
		return resourceError
	}
	// CreateOrUpdate(resourceGroupName string, resourceProviderNamespace string, parentResourcePath string, resourceType string, resourceName string, parameters resources.GenericResource, cancel <-chan struct{}, useHTTPPatch bool) (<-chan resources.GenericResource, <-chan error)

	srcServerResourceID := srcServer.ID
	fmt.Printf("srcServer ResourceId %s", *srcServerResourceID)

	// create the properties
	properties := map[string]interface{}{
		"sourceServerId":     srcServer.ID,
		"location":           srcServer.Location,
		"createMode":         CreateModePointInTimeRestore,
		"restorePointInTime": restorePointInTime,
	}
	// create the resource
	serverResource := resources.GenericResource{
		Location:   srcServer.Location,
		Properties: &properties,
	}
	_, errorChannel := groupClient.CreateOrUpdate(targetResourceGroup, namespace, "", resourceType, targetServerName, serverResource, nil, false)
	return <-errorChannel
}

// CreateFirewallRule creates a firewall rule
func (serversClient ServersClient) CreateFirewallRule(resourceGroupName string, serverName string, ruleName string, startIP string, endIP string) (err error) {
	// use the GroupClient from the rest_operations.go
	groupClient := NewGroupClient(serversClient.SubscriptionID)
	groupClient.Authorizer = serversClient.Authorizer
	// CreateOrUpdate(resourceGroupName string, resourceProviderNamespace string, parentResourcePath string, resourceType string, resourceName string, parameters resources.GenericResource, cancel <-chan struct{}, useHTTPPatch bool) (<-chan resources.GenericResource, <-chan error)
	// /subscriptions/31f97be2-2566-44f2-bb14-14d6924c8caa/resourceGroups/postgresql_from_go/providers/Microsoft.DBforPostgreSQL/servers/dar-db3/firewallRules/any-ip
	firewallRuleResource := resources.GenericResource{
		Properties: &map[string]interface{}{
			"startIpAddress": startIP,
			"endIpAddress":   endIP,
		},
	}
	_, errorChannel := groupClient.CreateOrUpdate(resourceGroupName, namespace, "servers/"+serverName, "firewallRules", ruleName, firewallRuleResource, nil, false)
	return <-errorChannel
}

//TODO: validate properties
func validateServerProperties(serverProperties ServerProperties) error {
	return nil
}
