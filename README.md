# azure-postgresql-go-sample
Initial sample of API access to Azure PostgreSQL service.  This is not production quality code and is intended to provide an interim ability to support (via API) the core Predix data servcies operations:

- Create instance
- Create firewall rule
- Destroy instance
- Set/Reset master user credentials
- Point-in-time-recovery

# Other notes
- Credentials are provided through environment variables.
- The swagger for postgresql is available at: https://github.com/Azure/azure-rest-api-specs/tree/current/specification/postgresql
- This update includes the go SDK with initial (beta) support for postgresql service: https://github.com/Azure/azure-sdk-for-go/releases/tag/v10.2.1-beta

