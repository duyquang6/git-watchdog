## Dependencies
- Docker
- Go 1.17
- MySQL 8.0.25

## Bootstrap
- Run `chmod +x start.sh` if start.sh script does not have privileged to run
- Run `./start.sh --bootstrap` quick bootstrap app (include build, start docker, migrate schema and start app), it will ready to accept connection to :8080 local
- Run `make docker.local.stop` to cleanup

## For developing
- Get tools for developing: `make install-go-tools`
- Build app docker image: run `make build.docker.image`
- Startup local docker compose: `make docker.local.start`
- Stop local docker compose: `make docker.local.stop`
- Migrate schema database: run `./start.sh --migrate`

## Automate CI CD local
Or you can use skaffold to automate that pipeline
- Install skaffold, helm latest, minikube latest version
- Run `skaffold dev --port-forward`
- Every change to source code, will trigger build, unit-test and deploy locally

## Testing & Coverage
- Run integration test(docker, go1.17 required): `./start.sh --integration`
- Run unittest: `make test.unit`
- Check coverage: `make coverage`
- Clean up report files: `make clean`
