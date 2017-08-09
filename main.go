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

//
// Notes:
// - in preview most properties can not be changed and only Basic SKU can be used
// - need to have provider registered:
//   az provider register --namespace Microsoft.DBforPostgreSQL
// - service instance parameters are hard coded as vars
// - credentials are read from environment
//

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/Azure/azure-sdk-for-go/arm/postgresql"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/adal"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/date"
	"github.com/Azure/go-autorest/autorest/to"
)

var (
	// Hard coded values for instance creation/update
	resourceGroupName          = "postgresql_from_go"
	location                   = "westus"
	administratorLogin         = "azadmin"
	administratorLoginPassword = "Welcome1234"

	// resource clients
	serversClient       postgresql.ServersClient
	firewallRulesClient postgresql.FirewallRulesClient
)

func main() {

	// storage sizes
	// default 0 -> 50 GB
	// 179200 MB -> 175 GB
	// 307200 MB ->300 GB

	createServer(resourceGroupName, "dar-95-50-175", location, administratorLogin, administratorLoginPassword, postgresql.NineFullStopFive, postgresql.Basic, 50, 179200)
	createFirewallRule(resourceGroupName, "dar-95-50-175", "all", "0.0.0.0", "255.255.255.255")

	createServer(resourceGroupName, "dar-96-100-300", location, administratorLogin, administratorLoginPassword, postgresql.NineFullStopSix, postgresql.Basic, 100, 307200)
	createFirewallRule(resourceGroupName, "dar-96-100-300", "myip", "0.0.0.0", "255.255.255.255")

	wait("check logins ... then press Enter")
	updateAdministratorPassword(resourceGroupName, "dar-95-50-175", "Welcome0000")
	updateAdministratorPassword(resourceGroupName, "dar-96-100-300", "Welcome0000")
	wait("check logins with new passwords ... then press Enter")

	wait("creating backup ... press Enter to start")

	utcNow := time.Now().UTC()
	restorePoint := utcNow.Add(time.Minute * 5 * -1)
	fmt.Printf("UTC now %v  restorePoint %v\n", utcNow, restorePoint)

	restoreServer(resourceGroupName, "dar-95-50-175", resourceGroupName, "dar-95-50-175-restored", restorePoint)
	wait("check backup server ... then Enter to delete servers")

	deleteServer(resourceGroupName, "dar-95-50-175")
	deleteServer(resourceGroupName, "dar-95-50-175-restored")
	deleteServer(resourceGroupName, "dar-96-100-300")

	fmt.Println("Done")
}

// createServer creates a server
func createServer(
	resourceGroup string,
	serverName string,
	location string,
	administratorLogin string,
	administratorLoginPassword string,
	serverVersion postgresql.ServerVersion,
	serverTier postgresql.SkuTier,
	computeUnits int32, //optional
	storageMB int64, // optional
) {

	fmt.Println("Creating server:" + resourceGroupName + "/" + serverName)
	spfdc := postgresql.ServerPropertiesForDefaultCreate{
		AdministratorLogin:         to.StringPtr(administratorLogin),
		AdministratorLoginPassword: to.StringPtr(administratorLoginPassword),
		SslEnforcement:             postgresql.SslEnforcementEnumEnabled,
		CreateMode:                 postgresql.CreateModeDefault,
		Version:                    serverVersion,
		StorageMB:                  to.Int64Ptr(storageMB),
	}

	properties, _ := spfdc.AsServerPropertiesForDefaultCreate()
	serverForCreate := postgresql.ServerForCreate{

		Location:   to.StringPtr(location),
		Properties: properties,
		Sku: &postgresql.Sku{
			Name:     to.StringPtr("SkuName"),
			Tier:     postgresql.Basic,
			Capacity: to.Int32Ptr(computeUnits),
		},
		Tags: &map[string]*string{
			"Tag1": to.StringPtr("1"),
		},
	}

	serverChannel, errChannel := serversClient.CreateOrUpdate(resourceGroupName, serverName, serverForCreate, nil)
	err := <-errChannel
	if err != nil {
		onErrorFail(err, "Create failed")
	}
	server := <-serverChannel

	fmt.Printf("Create server done. Response type: %s \n", toJSON(server))
}

// restore creates server from point-in-time state of source server
/*
 {
  "id": "/subscriptions/31f97be2-2566-44f2-bb14-14d6924c8caa/resourceGroups/postgresql_from_go/providers/Microsoft.DBforPostgreSQL/servers/dr-pwd-change",
  "name": "dr-pwd-change",
  "type": "Microsoft.DBforPostgreSQL/servers",
  "location": "westus",
  "sku": {
    "name": "PGSQLB100",
    "tier": "Basic",
    "capacity": 100
  },
  "properties": {
    "administratorLogin": "azadmin",
    "storageMB": 51200,
    "version": "9.6",
    "sslEnforcement": "Enabled",
    "userVisibleState": "Ready",
    "fullyQualifiedDomainName": "dr-pwd-change.postgres.database.azure.com"
  }
*/
func restoreServer(
	srcResourceGroup string,
	srcServerName string,
	targetResourceGroup string,
	targetServerName string,
	restorePoint time.Time,
) {
	fmt.Printf("Restore server source %s/%s target %s/%s point-in-time %s\n", srcResourceGroup, srcServerName, targetResourceGroup, targetServerName, restorePoint.String())
	srcServer, getServerErr := serversClient.Get(srcResourceGroup, srcServerName)
	if getServerErr != nil {
		onErrorFail(getServerErr, "Get source server details failed")
	}

	srcServerResourceID := srcServer.ID
	fmt.Printf("srcServer %s\n", toJSON(srcServer))
	fmt.Printf("srcServer ResourceId %s\n", *srcServerResourceID)

	spfr := postgresql.ServerPropertiesForRestore{
		CreateMode:         postgresql.CreateModePointInTimeRestore,
		SourceServerID:     srcServerResourceID,
		RestorePointInTime: &date.Time{Time: restorePoint},
	}

	properties, _ := spfr.AsServerPropertiesForRestore()
	serverForCreate := postgresql.ServerForCreate{
		Location:   srcServer.Location,
		Properties: properties,
	}

	fmt.Printf("Calling CreateOrUpdate %s\n", toJSON(serverForCreate))
	serverChannel, errChannel := serversClient.CreateOrUpdate(targetResourceGroup, targetServerName, serverForCreate, nil)
	err := <-errChannel
	if err != nil {
		onErrorFail(err, "Restore failed")
	}
	server := <-serverChannel
	fmt.Printf("Restore server done. Response type: %s \n", toJSON(server))
}

// create firewall rule
func createFirewallRule(
	resourceGroup string,
	serverName string,
	firewallRuleName string,
	startIPAddress string,
	endIPAddress string,
) {

	firewallRuleProperties := postgresql.FirewallRuleProperties{
		StartIPAddress: to.StringPtr(startIPAddress),
		EndIPAddress:   to.StringPtr(endIPAddress),
	}

	firewallRule := postgresql.FirewallRule{
		Name: to.StringPtr(firewallRuleName),
	}

	firewallRule.FirewallRuleProperties = &firewallRuleProperties
	fmt.Printf("Creating firewall %s/%s %s [%s][%s]\n", resourceGroup, serverName, firewallRuleName, startIPAddress, endIPAddress)
	_, errChannel := firewallRulesClient.CreateOrUpdate(resourceGroup, serverName, firewallRuleName, firewallRule, nil)
	err := <-errChannel
	if err != nil {
		onErrorFail(err, "firewall create failed")
	}
	fmt.Println("Creating firewall rule done")
}

// deleteServer deletes a server
func deleteServer(resourceGroupName string, serverName string) {
	fmt.Println("Delete server:" + resourceGroupName + "/" + serverName)
	responseChannel, errChannel := serversClient.Delete(resourceGroupName, serverName, nil)
	err := <-errChannel
	if err != nil {
		onErrorFail(err, "Delete failed")
	}
	response := <-responseChannel
	fmt.Printf("Delete server done.  Response: %s \n", toJSON(response))

}

func updateAdministratorPassword(resourceGroupName string, serverName string, newPassword string) {
	fmt.Println("changing password:" + resourceGroupName + "/" + serverName)

	serverUpdateParametersProperties := postgresql.ServerUpdateParametersProperties{
		AdministratorLoginPassword: to.StringPtr(newPassword),
	}

	serverUpdateParameters := postgresql.ServerUpdateParameters{}
	serverUpdateParameters.ServerUpdateParametersProperties = &serverUpdateParametersProperties

	serverChannel, errChannel := serversClient.Update(resourceGroupName, serverName, serverUpdateParameters, nil)
	err := <-errChannel
	if err != nil {
		onErrorFail(err, "Create failed")
	}
	server := <-serverChannel
	fmt.Printf("Parameter update done. Response: %s \n", toJSON(server))

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
	serversClient = postgresql.NewServersClient(subscriptionID)
	serversClient.Authorizer = authorizer
	firewallRulesClient = postgresql.FirewallRulesClient(serversClient)
}

func toJSON(v interface{}) string {
	j, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return "ERROR:" + err.Error()
	}
	return string(j)
}

func wait(prompt string) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prompt)
	reader.ReadString('\n')
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
