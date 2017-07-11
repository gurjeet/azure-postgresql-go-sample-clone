# azure-postgresql-go-sample
Initial sample of API access to Azure PostgreSQL service.  This is not production quality code and is intended to provide an interim ability to support (via API) the core Predix data servcies operations:

- Create instance
- Create firewall rule
- Destroy instance
- Set/Reset master user credentials
- Point-in-time-recovery

The swagger for postgresql is available at: https://github.com/Azure/azure-rest-api-specs/tree/current/specification/postgresql

There is a service facade patterned after the SQL Server API in postgresql/servers.go.  The main.go file includes scenarios for invoking the API operations through the facade.

Note: the code generated from the swagger API (when available) will be different than this example.  Progress on adding support for MySQL and PostgreSQL to the go SDK can be tracked here: [plans to support MySQL database service #654](https://github.com/Azure/azure-sdk-for-go/issues/654)

