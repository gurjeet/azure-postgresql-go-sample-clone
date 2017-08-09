package postgresql

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
//
// Code generated by Microsoft (R) AutoRest Code Generator 1.1.0.0
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

import (
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"net/http"
)

// DatabasesClient is the the Microsoft Azure management API provides create, read, update, and delete functionality
// for Azure PostgreSQL resources including servers, databases, firewall rules, log files and configurations.
type DatabasesClient struct {
	ManagementClient
}

// NewDatabasesClient creates an instance of the DatabasesClient client.
func NewDatabasesClient(subscriptionID string) DatabasesClient {
	return NewDatabasesClientWithBaseURI(DefaultBaseURI, subscriptionID)
}

// NewDatabasesClientWithBaseURI creates an instance of the DatabasesClient client.
func NewDatabasesClientWithBaseURI(baseURI string, subscriptionID string) DatabasesClient {
	return DatabasesClient{NewWithBaseURI(baseURI, subscriptionID)}
}

// CreateOrUpdate creates a new database or updates an existing database. This method may poll for completion. Polling
// can be canceled by passing the cancel channel argument. The channel will be used to cancel polling and any
// outstanding HTTP requests.
//
// resourceGroupName is the name of the resource group that contains the resource. You can obtain this value from the
// Azure Resource Manager API or the portal. serverName is the name of the server. databaseName is the name of the
// database. parameters is the required parameters for creating or updating a database.
func (client DatabasesClient) CreateOrUpdate(resourceGroupName string, serverName string, databaseName string, parameters Database, cancel <-chan struct{}) (<-chan Database, <-chan error) {
	resultChan := make(chan Database, 1)
	errChan := make(chan error, 1)
	go func() {
		var err error
		var result Database
		defer func() {
			if err != nil {
				errChan <- err
			}
			resultChan <- result
			close(resultChan)
			close(errChan)
		}()
		req, err := client.CreateOrUpdatePreparer(resourceGroupName, serverName, databaseName, parameters, cancel)
		if err != nil {
			err = autorest.NewErrorWithError(err, "postgresql.DatabasesClient", "CreateOrUpdate", nil, "Failure preparing request")
			return
		}

		resp, err := client.CreateOrUpdateSender(req)
		if err != nil {
			result.Response = autorest.Response{Response: resp}
			err = autorest.NewErrorWithError(err, "postgresql.DatabasesClient", "CreateOrUpdate", resp, "Failure sending request")
			return
		}

		result, err = client.CreateOrUpdateResponder(resp)
		if err != nil {
			err = autorest.NewErrorWithError(err, "postgresql.DatabasesClient", "CreateOrUpdate", resp, "Failure responding to request")
		}
	}()
	return resultChan, errChan
}

// CreateOrUpdatePreparer prepares the CreateOrUpdate request.
func (client DatabasesClient) CreateOrUpdatePreparer(resourceGroupName string, serverName string, databaseName string, parameters Database, cancel <-chan struct{}) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"databaseName":      autorest.Encode("path", databaseName),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"serverName":        autorest.Encode("path", serverName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2017-04-30-preview"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsJSON(),
		autorest.AsPut(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DBforPostgreSQL/servers/{serverName}/databases/{databaseName}", pathParameters),
		autorest.WithJSON(parameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare(&http.Request{Cancel: cancel})
}

// CreateOrUpdateSender sends the CreateOrUpdate request. The method will close the
// http.Response Body if it receives an error.
func (client DatabasesClient) CreateOrUpdateSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client,
		req,
		azure.DoPollForAsynchronous(client.PollingDelay))
}

// CreateOrUpdateResponder handles the response to the CreateOrUpdate request. The method always
// closes the http.Response Body.
func (client DatabasesClient) CreateOrUpdateResponder(resp *http.Response) (result Database, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusCreated, http.StatusAccepted),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// Delete deletes a database. This method may poll for completion. Polling can be canceled by passing the cancel
// channel argument. The channel will be used to cancel polling and any outstanding HTTP requests.
//
// resourceGroupName is the name of the resource group that contains the resource. You can obtain this value from the
// Azure Resource Manager API or the portal. serverName is the name of the server. databaseName is the name of the
// database.
func (client DatabasesClient) Delete(resourceGroupName string, serverName string, databaseName string, cancel <-chan struct{}) (<-chan autorest.Response, <-chan error) {
	resultChan := make(chan autorest.Response, 1)
	errChan := make(chan error, 1)
	go func() {
		var err error
		var result autorest.Response
		defer func() {
			if err != nil {
				errChan <- err
			}
			resultChan <- result
			close(resultChan)
			close(errChan)
		}()
		req, err := client.DeletePreparer(resourceGroupName, serverName, databaseName, cancel)
		if err != nil {
			err = autorest.NewErrorWithError(err, "postgresql.DatabasesClient", "Delete", nil, "Failure preparing request")
			return
		}

		resp, err := client.DeleteSender(req)
		if err != nil {
			result.Response = resp
			err = autorest.NewErrorWithError(err, "postgresql.DatabasesClient", "Delete", resp, "Failure sending request")
			return
		}

		result, err = client.DeleteResponder(resp)
		if err != nil {
			err = autorest.NewErrorWithError(err, "postgresql.DatabasesClient", "Delete", resp, "Failure responding to request")
		}
	}()
	return resultChan, errChan
}

// DeletePreparer prepares the Delete request.
func (client DatabasesClient) DeletePreparer(resourceGroupName string, serverName string, databaseName string, cancel <-chan struct{}) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"databaseName":      autorest.Encode("path", databaseName),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"serverName":        autorest.Encode("path", serverName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2017-04-30-preview"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsDelete(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DBforPostgreSQL/servers/{serverName}/databases/{databaseName}", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare(&http.Request{Cancel: cancel})
}

// DeleteSender sends the Delete request. The method will close the
// http.Response Body if it receives an error.
func (client DatabasesClient) DeleteSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client,
		req,
		azure.DoPollForAsynchronous(client.PollingDelay))
}

// DeleteResponder handles the response to the Delete request. The method always
// closes the http.Response Body.
func (client DatabasesClient) DeleteResponder(resp *http.Response) (result autorest.Response, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusAccepted, http.StatusNoContent),
		autorest.ByClosing())
	result.Response = resp
	return
}

// Get gets information about a database.
//
// resourceGroupName is the name of the resource group that contains the resource. You can obtain this value from the
// Azure Resource Manager API or the portal. serverName is the name of the server. databaseName is the name of the
// database.
func (client DatabasesClient) Get(resourceGroupName string, serverName string, databaseName string) (result Database, err error) {
	req, err := client.GetPreparer(resourceGroupName, serverName, databaseName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "postgresql.DatabasesClient", "Get", nil, "Failure preparing request")
		return
	}

	resp, err := client.GetSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "postgresql.DatabasesClient", "Get", resp, "Failure sending request")
		return
	}

	result, err = client.GetResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "postgresql.DatabasesClient", "Get", resp, "Failure responding to request")
	}

	return
}

// GetPreparer prepares the Get request.
func (client DatabasesClient) GetPreparer(resourceGroupName string, serverName string, databaseName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"databaseName":      autorest.Encode("path", databaseName),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"serverName":        autorest.Encode("path", serverName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2017-04-30-preview"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DBforPostgreSQL/servers/{serverName}/databases/{databaseName}", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare(&http.Request{})
}

// GetSender sends the Get request. The method will close the
// http.Response Body if it receives an error.
func (client DatabasesClient) GetSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req)
}

// GetResponder handles the response to the Get request. The method always
// closes the http.Response Body.
func (client DatabasesClient) GetResponder(resp *http.Response) (result Database, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// ListByServer list all the databases in a given server.
//
// resourceGroupName is the name of the resource group that contains the resource. You can obtain this value from the
// Azure Resource Manager API or the portal. serverName is the name of the server.
func (client DatabasesClient) ListByServer(resourceGroupName string, serverName string) (result DatabaseListResult, err error) {
	req, err := client.ListByServerPreparer(resourceGroupName, serverName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "postgresql.DatabasesClient", "ListByServer", nil, "Failure preparing request")
		return
	}

	resp, err := client.ListByServerSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "postgresql.DatabasesClient", "ListByServer", resp, "Failure sending request")
		return
	}

	result, err = client.ListByServerResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "postgresql.DatabasesClient", "ListByServer", resp, "Failure responding to request")
	}

	return
}

// ListByServerPreparer prepares the ListByServer request.
func (client DatabasesClient) ListByServerPreparer(resourceGroupName string, serverName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"serverName":        autorest.Encode("path", serverName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2017-04-30-preview"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DBforPostgreSQL/servers/{serverName}/databases", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare(&http.Request{})
}

// ListByServerSender sends the ListByServer request. The method will close the
// http.Response Body if it receives an error.
func (client DatabasesClient) ListByServerSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req)
}

// ListByServerResponder handles the response to the ListByServer request. The method always
// closes the http.Response Body.
func (client DatabasesClient) ListByServerResponder(resp *http.Response) (result DatabaseListResult, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}
