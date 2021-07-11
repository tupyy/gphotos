#####################
#  					# 
#   Infra targets   #
#					#
#####################

INFRA_NETWORK=gophoto_network

.PHONY: run.infra run.infra.stop run.docker.network run.docker.network.stop

#####################
# ENV				#
#####################
ifeq ($(ENV),dev)
	ENV_SET=1
else ifeq ($(ENV),prod)
	ENV_SET=1
endif

#help deploy.check: check required ENV variable
run.check:
	@if [ "$(ENV_SET)" = "" ]; then echo "$(COLOR_RED)ERROR: ENV is mandatory, could be 'dev' or 'prod'$(RESET_COLOR)"; exit 1; fi

#help run.infra: start postgres and keycloak
run.infra: 
	@if [ "$(ENV)" == "dev" ]; then\
		echo "$(COLOR_YELLOW)Run for dev$(RESET_COLOR)"; \
		make run.docker.network run.docker.postgres; \
		echo "$(COLOR_YELLOW)Wait for postgres to start$(RESET_COLOR)";\
		sleep 20;\
		make postgres.setup run.docker.keycloak; \
	else \
		echo "run infra for prod...unavailable";\
	fi

#help run.infra.stop: stop postgre and keycloak
run.infra.stop: run.docker.keycloak.stop run.docker.postgres.stop run.docker.network.stop

#help run.docker.network: create the infra network
run.docker.network:
	$(DOCKER_CMD) network create -d bridge $(INFRA_NETWORK) || true

#help run.docker.network.stop: remove the infra network
run.docker.network.stop:
	$(DOCKER_CMD) network rm $(INFRA_NETWORK) || true

##################################
#
# 			Postgres  
#
##################################

DB_HOST=localhost
DB_PORT=5432
POSTGRES_CONTAINER=postgresql
IMAGE_NAME=postgres
IMAGE_TAG=13
PG_DATA=/home/cosmin/tmp/pgdata
ROOT_USER=postgres
ROOT_PWD=$(shell cat $(PGPASSFILE) | head -n 1 | cut -d":" -f5)
USER_ID=$(shell id `whoami` -u)
GROUP_ID=$(shell id `whoami` -g)
RESOURCE_ADMIN_USER=resources_admin
RESOURCE_ADMIN_PWD=$(shell cat $(PGPASSFILE) | grep $(RESOURCE_ADMIN_USER) | cut -d":" -f5)
PGPASSFILE=$(CURDIR)/sql/.pgpass
PSQL_COMMAND=PGPASSFILE=$(PGPASSFILE) psql --quiet --host=$(DB_HOST) --port=$(DB_PORT) -v ON_ERROR_STOP=on
GOPHOTO_USER=gophoto
GOPHOTO_PWD=$(shell cat $(PGPASSFILE) | grep $(GOPHOTO_USER) | cut -d":" -f5)
KEYCLOAK_DB_USER=keycloak
KEYCLOAK_DB_PWD=$(shell cat $(PGPASSFILE) | grep $(KEYCLOAK_DB_USER) | cut -d":" -f5)

.PHONY: run.docker.postgres run.docker.postgres.stop run.docker.postgres.stop run.docker.postgres.restart

#help run.docker.postgres: run postgres using docker
run.docker.postgres:
	@if [ "$(ENV)" == "dev" ]; then \
		$(DOCKER_CMD) run --rm -d -p $(DB_PORT):5432 \
		--network=$(INFRA_NETWORK) \
		-e POSTGRES_USER=$(ROOT_USER) \
		-e POSTGRES_PASSWORD=$(ROOT_PWD) \
		-e VERBOSE=1 \
		--name $(POSTGRES_CONTAINER) $(IMAGE_NAME):$(IMAGE_TAG); \
	else \
		$(DOCKER_CMD) run --rm -d -p $(DB_PORT):5432 \
		--network=$(INFRA_NETWORK) \
		-e POSTGRES_USER=$(ROOT_USER) \
		-e POSTGRES_PASSWORD=$(ROOT_PWD) \
		-e VERBOSE=1 \
		-v $(PG_DATA):/var/lib/postgresql/data \
		--user $(USER_ID):$(GROUP_ID) \
		--name $(POSTGRES_CONTAINER) $(IMAGE_NAME):$(IMAGE_TAG); \
	fi

#help run.docker.postgres.logs: show logs from postgres
run.docker.postgres.logs:
	$(DOCKER_CMD) logs -f $(POSTGRES_CONTAINER)

#help run.docker.postgres.stop: stop postgres docker
run.docker.postgres.stop:
	$(DOCKER_CMD) stop $(POSTGRES_CONTAINER)

#help run.docker.postgres.restart: run.postgres.docker.restart
run.docker.postgres.restart: 
	$(DOCKER_CMD) restart $(POSTGRES_CONTAINER)


#################
# Setup targets #
#################

.PHONY: postgres.setup.clean postgres.setup.init postgres.setup.tables

# help postgres.setup: Setup postgres from scratch
postgres.setup: postgres.setup.init postgres.setup.tables

#help postgres.setup.clean: cleans postgres from all created resources
postgres.setup.clean:
	$(PSQL_COMMAND) --user=$(ROOT_USER) -f sql/clean/clean.sql

#help postgres.setup.init: init the database
postgres.setup.init:
	$(PSQL_COMMAND) --dbname=postgres --user=$(ROOT_USER) \
		-v resources_admin_pwd="'$(RESOURCE_ADMIN_PWD)'" \
		-v gophoto_pwd="'$(GOPHOTO_PWD)'" \
		-v keycloak_pwd="'$(KEYCLOAK_DB_PWD)'" \
		-f sql/setup/01_init.sql

#help postgres.setup.users: init postgres users
postgres.setup.tables:
	$(PSQL_COMMAND) --dbname=gophoto --user=$(RESOURCE_ADMIN_USER) \
		-f sql/setup/02_setup.sql


############################
# Model generation targets #
############################

BASE_CONNSTR="postgresql://$(RESOURCE_ADMIN_USER):$(RESOURCE_ADMIN_PWD)@$(DB_HOST):$(DB_PORT)"
GEN_CMD=$(TOOLS_DIR)/gen --sqltype=postgres \
	--module=github.com/tupyy/gophoto \
	--gorm --no-json --no-xml --overwrite --mapping tools/mappings.json

.PHONY: generate.models

#help generate.models: generate models for the gophoto database
generate.models:
	sh -c '$(GEN_CMD) --connstr "$(BASE_CONNSTR)/gophoto?sslmode=disable"  --model=models --database gophoto' 						# Generate models for the DB tables

##################################
#
# 			Keycloak  
#
##################################

KEYCLOAK_IMAGE=jboss/keycloak
KEYCLOAK_TAG=12.0.4
KEYCLOAK_CONTAINER=keycloak
KEYCLOAK_PORT=9000
KEYCLOAK_USER=$(shell cat $(CURDIR)/resources/keycloak/.pass | cut -d":" -f1)
KEYCLOAK_PWD=$(shell cat $(CURDIR)/resources/keycloak/.pass | cut -d":" -f2)
KEYCLOAK_REALM_FILE=$(CURDIR)/resources/keycloak/gophoto-realm-$(ENV).json

.PHONY: run.docker.keycloak run.docker.keycloak.stop run.docker.keycloak.restart

#help run.docker.keycloak.setup: run keycloak without the realm file. used for setup realm
run.docker.keycloak.setup:
	$(DOCKER_CMD) run --rm -d -p $(KEYCLOAK_PORT):8080 \
		--network=$(INFRA_NETWORK) \
		-e KEYCLOAK_USER=$(KEYCLOAK_USER) \
		-e KEYCLOAK_PASSWORD=$(KEYCLOAK_PWD) \
		--name $(KEYCLOAK_CONTAINER) $(KEYCLOAK_IMAGE):$(KEYCLOAK_TAG)

#help run.docker.keycloak: run keycloak
run.docker.keycloak:
		$(DOCKER_CMD) run --rm -d -p $(KEYCLOAK_PORT):8080 \
		--network=$(INFRA_NETWORK) \
		-e KEYCLOAK_USER=$(KEYCLOAK_USER) \
		-e KEYCLOAK_PASSWORD=$(KEYCLOAK_PWD) \
		-e KEYCLOAK_IMPORT='/tmp/gophoto-realm.json' \
		-e DB_VENDOR='postgres' \
		-e DB_ADDR="$(POSTGRES_CONTAINER)" \
		-e DB_PORT="$(DB_PORT)" \
		-e DB_DATABASE='keycloak' \
		-e DB_USER="$(KEYCLOAK_DB_USER)" \
		-e DB_PASSWORD="$(KEYCLOAK_DB_PWD)" \
		-v $(KEYCLOAK_REALM_FILE):'/tmp/gophoto-realm.json' \
		--name $(KEYCLOAK_CONTAINER) $(KEYCLOAK_IMAGE):$(KEYCLOAK_TAG)
	
#help run.docker.keycloak.stop: stop keycloak
run.docker.keycloak.stop:
	$(DOCKER_CMD) stop $(KEYCLOAK_CONTAINER)

#help run.docker.keycloak.restart: restart keycloak
run.docker.keycloak.restart:
	$(DOCKER_CMP) restart $(KEYCLOAK_CONTAINER)
