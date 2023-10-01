# About

[Alin.io](http://Alin.io) pkgstore is a simple NPM and Pypi registry server, which also acts as a proxy to the generic public registries. It is built for easy maintainability and performance.

pkgstore is built with an extendable structure that allows adding more storage backends or databases to keep the package metadata information. Currently, by default, the storage backend is an AWS S3 bucket or Minio Bucket if you have a self-hosted environment.

The database is a simple SQLite file, which is configurable from the environment variable of `DATABASE_URL`, and it acts as a database type selector based on the given database URL prefix, like if you have a `postgresql://...` then the database instance will act with a PostgreSQL driver. Otherwise, it will fall back to SQLite.

## Running Locally

To run locally, you will need the Minio service instance, which is configured in Docker-Compose

```bash
~# container-compose up -d minio
~# go run ./cmd/server

[GIN-debug] [WARNING] Creating an Engine instance with the Logger and Recovery middleware already attached.

[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:   export GIN_MODE=release
 - using code:  gin.SetMode(gin.ReleaseMode)

[GIN-debug] GET    /                         --> github.com/alin-io/pkgstore/services.HealthCheckHandler (6 handlers)
[GIN-debug] GET    /npm/*path                --> github.com/alin-io/pkgstore/router.PackageRouter.HandleFetch.func1 (6 handlers)
[GIN-debug] GET    /pypi/*path               --> github.com/alin-io/pkgstore/router.PackageRouter.HandleFetch.func2 (6 handlers)
[GIN-debug] PUT    /npm/*path                --> github.com/alin-io/pkgstore/services/npm.(*Service).UploadHandler-fm (6 handlers)
[GIN-debug] POST   /pypi/*path               --> github.com/alin-io/pkgstore/services/pypi.(*Service).UploadHandler-fm (6 handlers)
[GIN-debug] [WARNING] You trusted all proxies, this is NOT safe. We recommend you to set a value.
Please check https://pkg.go.dev/github.com/gin-gonic/gin#readme-don-t-trust-all-proxies for details.
[GIN-debug] Listening and serving HTTP on :8080
```