#####################
#  					# 
#   Infra targets   #
#					#
#####################

.PHONY: run.infra run.infra.stop 

#####################
# ENV				#
#####################
ifeq ($(ENV),dev)
	ENV_SET=1
else ifeq ($(ENV),prod)
	ENV_SET=1
endif

DOCKER_COMPOSE_COMMAND=docker-compose
#help deploy.check: check required ENV variable
run.check:
	@if [ "$(ENV_SET)" = "" ]; then echo "$(COLOR_RED)ERROR: ENV is mandatory, could be 'dev' or 'prod'$(RESET_COLOR)"; exit 1; fi

#help run.infra: start postgres and keycloak
run.infra: 
	@if [ "$(ENV)" == "dev" ]; then\
		echo "$(COLOR_YELLOW)Run for dev$(RESET_COLOR)"; \
		$(DOCKER_COMPOSE_COMMAND) -f $(CURDIR)/resources/docker-compose.yaml --env-file $(CURDIR)/resources/dev.env up -d; \
	elif [ "$(ENV)" == "prod" ]; then\
		echo "$(COLOR_YELLOW)Run for prod$(RESET_COLOR)"; \
		IMAGE_NAME=$(IMAGE_NAME) IMAGE_TAG=$(IMAGE_TAG) RESOURCES_FOLDER=$(CURDIR)/resources STATICS_FOLDER=$(CURDIR)/assets $(DOCKER_COMPOSE_COMMAND) -f $(CURDIR)/resources/docker-compose-prod.yaml --env-file $(CURDIR)/resources/prod.env up -d; \
	fi

#help run.infra.stop: shutdown infra 
run.infra.stop: 
	$(DOCKER_COMPOSE_COMMAND) -f $(CURDIR)/resources/docker-compose.yaml --env-file $(CURDIR)/resources/dev.env down --remove-orphans

##################################
#
# 			Postgres  
#
##################################

DB_HOST=localhost
DB_PORT=5432
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

#################
# Setup targets #
#################

.PHONY: postgres.setup.clean postgres.setup.init postgres.setup.tables postgres.setup.migrations

#help postgres.setup: Setup postgres from scratch
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

postgres.setup.migrations:
	$(PSQL_COMMAND) --dbname=gophoto --user=$(RESOURCE_ADMIN_USER) \
		-f sql/setup/04_migrations.sql


