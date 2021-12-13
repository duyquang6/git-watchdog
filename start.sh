#!/bin/sh
set -e

usage()
{
    echo "usage: start.sh [[[-i/--integration] (to run integration tests too) | [-h]]"
}

# Run local migration for developing
migrate() {
  make build.migrate
  export $(cat .env_template | xargs)
  ./bin/migrate
}

# Run local app for developing
start() {
  make build
  export $(cat .env_template | xargs)
  ./bin/app
}

# Run local app for developing
bootstrap() {
  echo "----------- SETUP DOCKER -------------"
  make docker.local.stop
  make docker.local.start

  echo "Wait 30 secs for db start complete & run migration"
  declare -ir MAX_SECONDS=30
  declare -ir TIMEOUT=$SECONDS+$MAX_SECONDS
  while (( $SECONDS < $TIMEOUT )); do
      if [ $(healthcheck) == "200" ]; then
          break
      fi
      echo "starting ...."
      sleep 3
  done
  echo "----------- DONE -------------"
  echo "Tryout api spec at http://0.0.0.0:8080/swagger/index.html !"
  echo "Run make docker.local.stop to stop docker instances"
}

healthcheck(){
  curl -s -o /dev/null --head -w "%{http_code}" -X GET "http://0.0.0.0:8080/ping"
}

# Boot up docker & run integration test
integration() {
  echo "----------- SETUP DOCKER -------------"
  make docker.integration.stop
  make docker.integration.start

  echo "Wait 30 secs for db start complete"
  sleep 30
  echo "----------- INTEGRATION TEST -------------"
  RULE_FILE_PATH=$(pwd)/assets/rules.json make test.integration

  echo "----------- CLEAN UP -------------"
  make docker.integration.stop
}

while [ "$1" != "" ]; do
    case $1 in
        -m | --migrate )        migrate
                                ;;
        -s | --start )          start
                                ;;
        -i | --integration )    integration
                                ;;
        -b | --bootstrap )      bootstrap
                                ;;
        -h | --help )           usage
                                exit
                                ;;
        * )                     usage
                                exit 1
    esac
    shift
done
