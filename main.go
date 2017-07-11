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

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/adal"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/dave-read/azure-postgresql-go-sample/postgresql"
)

var (
	// Hard coded values for instance creation/update
	resourceGroupName          = "postgresql_from_go"
	location                   = "westus"
	administratorLogin         = "azadmin"
	administratorLoginPassword = "Welcome1234"

	// postgresql client
	postgresqlClient postgresql.ServersClient
)

func main() {
	//var input string
	/*
		createServer(resourceGroupName, "dar-95-50-0", location, administratorLogin, administratorLoginPassword, "9.5", postgresql.SkuTierBasic, 50, 0)
		createServer(resourceGroupName, "dar-95-100-0", location, administratorLogin, administratorLoginPassword, "9.5", postgresql.SkuTierBasic, 100, 0)
		createServer(resourceGroupName, "dar-96-50-0", location, administratorLogin, administratorLoginPassword, "9.6", postgresql.SkuTierBasic, 50, 0)
		createServer(resourceGroupName, "dar-96-100-0", location, administratorLogin, administratorLoginPassword, "9.6", postgresql.SkuTierBasic, 100, 0)
	*/

	/*
		deleteServer(resourceGroupName, "dar-95-50-0")
		deleteServer(resourceGroupName, "dar-95-100-0")
		deleteServer(resourceGroupName, "dar-96-50-0")
		deleteServer(resourceGroupName, "dar-96-100-0")
	*/

	/*
		createFirewallRule(resourceGroupName, "dar-db2", "all-ips", "0.0.0.0", "255.255.255.255")
	*/

	// Predix storage sizes 50/100/200 GB

	createServer(resourceGroupName, "dar-95-50-50", location, administratorLogin, administratorLoginPassword, "9.5", postgresql.SkuTierBasic, 50, 50)
	createServer(resourceGroupName, "dar-95-50-100", location, administratorLogin, administratorLoginPassword, "9.5", postgresql.SkuTierBasic, 100, 100)
	createServer(resourceGroupName, "dar-95-50-200", location, administratorLogin, administratorLoginPassword, "9.5", postgresql.SkuTierBasic, 100, 200)

	/*
		deleteServer(resourceGroupName, "dar-95-50-50")
		deleteServer(resourceGroupName, "dar-95-50-100")
		deleteServer(resourceGroupName, "dar-95-50-200")
	*/

	/*
		updateAdministratorPassword(resourceGroupName, "serverName", "Welcome12345")
	*/
}

// createServerGroup creates a resource group
/*
func createResourceGroup() resources.Group {
	fmt.Println("Create resource group:" + resourceGroupName)
	rgParms := resources.Group{
		Location: to.StringPtr(location),
	}
	rg, err := groupsClient.CreateOrUpdate(resourceGroupName, rgParms)
	onErrorFail(err, "CreateOrUpdate resource group failed")
	return rg
}
*/

// createServer creates a generic resource
func createServer(
	resourceGroup string,
	serverName string,
	location string,
	administratorLogin string,
	administratorLoginPassword string,
	serverVersion postgresql.ServerVersion,
	serverTier postgresql.SkuTier,
	computeUnits int32, //optional
	storageMB int32, // optional
) {

	serverProperties := postgresql.ServerProperties{
		Location:                   location,
		AdministratorLogin:         administratorLogin,
		AdministratorLoginPassword: administratorLoginPassword,
		Version:                    serverVersion,
		Tier:                       serverTier,
	}

	if computeUnits > 0 {
		serverProperties.ComputeUnits = computeUnits
	}
	if storageMB > 0 {
		serverProperties.StorageMB = storageMB
	}

	err := postgresqlClient.CreateServer(resourceGroupName, serverName, serverProperties)
	if err != nil {
		onErrorFail(err, "Create failed")
	}
	fmt.Println("Create server done")
}

// create firewall rule
func createFirewallRule(
	resourceGroup string,
	serverName string,
	firewallRuleName string,
	startIPAddress string,
	endIPAddress string,
) {

	err := postgresqlClient.CreateFirewallRule(resourceGroup, serverName, firewallRuleName, startIPAddress, endIPAddress)
	//CreateFirewallRule(resourceGroupName string, serverName string, ruleName string, startIP string, endIP string) (err error)
	if err != nil {
		onErrorFail(err, "firewall create failed")
	}
	fmt.Println("Create firewall done")
}

// deleteServer deletes a server
func deleteServer(resourceGroupName string, serverName string) {
	fmt.Println("Delete server:" + resourceGroupName + "/" + serverName)
	err := postgresqlClient.DeleteServer(resourceGroupName, serverName)
	if err != nil {
		onErrorFail(err, "Delete failed")
	}
	fmt.Println("Delete server done")
}

func updateAdministratorPassword(resourceGroupName string, serverName string, newPassword string) {
	fmt.Println("changing password:" + resourceGroupName + "/" + serverName)
	err := postgresqlClient.ChangeAdministratorPassword(resourceGroupName, serverName, newPassword)
	if err != nil {
		onErrorFail(err, "changing password failed")
	}
	fmt.Println("password change done")
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
		os.Exit(1)
	}
}

func createClients(subscriptionID string, authorizer *autorest.BearerAuthorizer) {
	postgresqlClient = postgresql.NewServersClient(subscriptionID)
	postgresqlClient.Authorizer = authorizer
}

func toJSON(v interface{}) (string, error) {
	j, err := json.MarshalIndent(v, "", "  ")
	return string(j), err
}

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
