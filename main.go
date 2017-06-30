package main

// Copyright (c) Microsoft.  All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//--------------------------------------------------------------------------

// Example of using resources.GenericResource from the Azure SDK for Go
// https://github.com/Azure/azure-sdk-for-go
//
// based on samples here: https://github.com/Azure-Samples/resource-manager-go-resources-and-groups
//
// Notes:
// - in preview most properties can not be changed and only Basic SKU can be used
// - need to have provider registered:
//   az provider register --namespace Microsoft.DBforPostgreSQL
// - service instance parameters are hard coded as vars
// - credentials are read from environment
//
// Open Questions for PG:
// - Is is better to use template? https://gallery.azure.com/artifact/20161101/Microsoft.PostgreSQLServer.1.0.18/DeploymentTemplates/NewPostgreSqlServer.json
// - Udpate resource looks like all initial parms are needed.  Is there ability to just sent the changed values
// - If there's a failure during deploy, it looks like the entire resource group is removed
//
// Open Questions for Predix:
//   Does current implementation use polling?
//   golang version?
//   how does naming work today (PostgreSQL instance is a global name)
//

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/Azure/azure-sdk-for-go/arm/resources/resources"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/adal"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/to"
)

var (
	// Hard coded values for instance creation/update
	resourceGroupName                = "postgresql_from_goy"
	location                         = "westus"
	namespace                        = "Microsoft.DBforPostgreSQL"
	resourceType                     = "servers"
	administratorLogin               = "azadmin"
	administratorLoginPassword       = "Welcome1234"
	version                          = "9.5"
	storageMB                        = 307200
	sslEnforcement                   = "Disabled"
	tier                             = "Basic"
	capacity                   int32 = 100

	// PostgreSQL instance name.  Must be globally unique
	serverName = "azcat-db7"

	// client to create resource groups
	groupsClient resources.GroupsClient
	// client to create resources
	resourcesClient resources.GroupClient
)

// Create the clients
func init() {
	// credentials read from environment
	subscriptionID := getEnvVarOrExit("AZURE_SUBSCRIPTION_ID")
	tenantID := getEnvVarOrExit("AZURE_TENANT_ID")

	oauthConfig, err := adal.NewOAuthConfig(azure.PublicCloud.ActiveDirectoryEndpoint, tenantID)
	onErrorFail(err, "Error getting OAuth configuration")

	clientID := getEnvVarOrExit("AZURE_CLIENT_ID")
	clientSecret := getEnvVarOrExit("AZURE_CLIENT_SECRET")
	spToken, err := adal.NewServicePrincipalToken(*oauthConfig, clientID, clientSecret, azure.PublicCloud.ResourceManagerEndpoint)
	authorizer := autorest.NewBearerAuthorizer(spToken)
	onErrorFail(err, "NewServicePrincipalToken failed")

	createClients(subscriptionID, authorizer)
}

func main() {
	var input string
	// ensure the resource group is created
	fmt.Println("Creating resource group:" + resourceGroupName)
	rg := createResourceGroup()
	rgAsJSON, merr := toJSON(rg)
	if merr == nil {
		fmt.Println("Resource Group result:" + string(rgAsJSON))
	}
	/*
		server := createServer()
		serverAsJSON, serr := toJSON(server)
		if serr == nil {
			fmt.Println("Create Server result:" + string(serverAsJSON))
		}
		fmt.Print("Server created. Press enter to update Server")
		var input string
		fmt.Scanln(&input)
	*/
	sslEnforcement = "Enabled"
	updateServer()
	fmt.Print("Server updated. Press enter to delete Server")
	fmt.Scanln(&input)

	deleteServer()
}

// createServerGroup creates a resource group
func createResourceGroup() resources.Group {
	fmt.Println("Create resource group:" + resourceGroupName)
	rgParms := resources.Group{
		Location: to.StringPtr(location),
	}
	rg, err := groupsClient.CreateOrUpdate(resourceGroupName, rgParms)
	onErrorFail(err, "CreateOrUpdate resource group failed")
	return rg
}

// createServer creates a generic resource
func createServer() resources.GenericResource {
	fmt.Println("PostgreSQL instance via Generic Resource Put")
	sku := &resources.Sku{
		Name:     to.StringPtr("SkuName"),
		Tier:     to.StringPtr(tier),
		Capacity: to.Int32Ptr(capacity),
	}

	genericResource := resources.GenericResource{
		Location: to.StringPtr(location),
		Properties: &map[string]interface{}{
			"location":                   location,
			"administratorLogin":         administratorLogin,
			"administratorLoginPassword": administratorLoginPassword,
			"version":                    version,
			"storageMB":                  storageMB,
			"sslEnforcement":             sslEnforcement,
		},
		Sku: sku,
	}
	respChannel, errChannel := resourcesClient.CreateOrUpdate(resourceGroupName, namespace, "", resourceType, serverName, genericResource, nil)
	select {
	case resp := <-respChannel:
		genericResource = resp
	case err := <-errChannel:
		onErrorFail(err, "Create failed")
	}
	return genericResource
}

// updateResource updates a generic resource
// TODO:
//    - need to find out if it's possible to just send the changed attributes
//    - if not, check on how to read back the current set of full configuration properites

func updateServer() resources.GenericResource {

	sku := &resources.Sku{
		Name:     to.StringPtr("SkuName"),
		Tier:     to.StringPtr(tier),
		Capacity: to.Int32Ptr(capacity),
	}

	genericResource := resources.GenericResource{
		Location: to.StringPtr(location),
		Properties: &map[string]interface{}{
			"location":                   location,
			"administratorLogin":         administratorLogin,
			"administratorLoginPassword": administratorLoginPassword,
			"version":                    version,
			"storageMB":                  storageMB,
			"sslEnforcement":             sslEnforcement,
		},
		Sku: sku,
	}
	resultChannel, errorChannel := resourcesClient.CreateOrUpdate(resourceGroupName, namespace, "", resourceType, serverName, genericResource, nil)
	err := <-errorChannel
	if err != nil {
		onErrorFail(err, "Update failed")
	}
	return <-resultChannel
}

// deleteServer deletes a generic resource

func deleteServer() {
	fmt.Println("Delete a resource")
	resultChannel, errorChannel := resourcesClient.Delete(resourceGroupName, namespace, "", resourceType, serverName, nil)
	err := <-errorChannel
	if err != nil {
		onErrorFail(err, "Delete failed")
	}
	deleteResult, merr := toJSON(<-resultChannel)
	if merr == nil {
		fmt.Println("Delete result:" + deleteResult)
	}
}

// getEnvVarOrExit returns the value of specified environment variable or terminates if it's not defined.
func getEnvVarOrExit(varName string) string {
	value := os.Getenv(varName)
	if value == "" {
		fmt.Printf("Missing environment variable %s\n", varName)
		os.Exit(1)
	}

	return value
}

// onErrorFail prints a failure message and exits the program if err is not nil.
// it also deletes the resource group created in the sample
func onErrorFail(err error, message string) {
	if err != nil {
		fmt.Printf("%s: %s\n", message, err)
		groupsClient.Delete(resourceGroupName, nil)
		os.Exit(1)
	}
}

func createClients(subscriptionID string, authorizer *autorest.BearerAuthorizer) {
	groupsClient = resources.NewGroupsClient(subscriptionID)
	groupsClient.Authorizer = authorizer

	resourcesClient = resources.NewGroupClient(subscriptionID)
	resourcesClient.Authorizer = authorizer
	//resourcesClient.APIVersion = "2017-04-30-preview"
}

func toJSON(v interface{}) (string, error) {
	j, err := json.MarshalIndent(v, "", "  ")
	return string(j), err
}
