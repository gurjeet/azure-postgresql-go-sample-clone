# azure-postgresql-go-sample
Initial sample of API access to Azure PostgreSQL service.  This is not production quality code and is intended to provide an interim ability to support (via API) the core Predix data servcies operations:

- Create instance
    - __Includes preliminary example of async server create rather than polling with the following changes:__
    - _Updated [server.go](vendor/github.com/Azure/azure-sdk-for-go/arm/postgresql/servers.go#L53) to return HTTP header rather than polling for Server struct_
    - _Updated [main.go](main.go) ..._
      - _added function to get the polling URL from the header_
      - _added function to get and parse the results of calling the polling URL_
      - _added loop at top of main to poll for status other than Provisioning or timeout_
- Create firewall rule
- Destroy instance
- Set/Reset master user credentials
- Point-in-time-recovery

# Other notes
- main.go provides the example
- Credentials are provided through environment variables.
- The swagger for postgresql is available at: https://github.com/Azure/azure-rest-api-specs/tree/current/specification/postgresql
- This update includes the go SDK with initial (beta) support for postgresql service: https://github.com/Azure/azure-sdk-for-go/releases/tag/v10.2.1-beta

