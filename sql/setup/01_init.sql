CREATE USER resources_admin WITH CREATEDB CREATEROLE PASSWORD :resources_admin_pwd;
GRANT resources_admin TO postgres;

-- Create databases and remove default permissions on public schema to ensure readonly permissions
-- are well applied 
CREATE DATABASE gophoto OWNER resources_admin;
REVOKE ALL ON DATABASE gophoto FROM PUBLIC;

CREATE DATABASE keycloak OWNER resources_admin;
REVOKE ALL ON DATABASE keycloak FROM PUBLIC;

REVOKE CREATE ON SCHEMA public FROM PUBLIC;

\connect keycloak

-- Create core device management role

DROP ROLE IF EXISTS keycloak_readonly;
CREATE ROLE keycloak_readwrite;
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO keycloak_readwrite;

-- Create users
DROP USER IF EXISTS keycloak;
CREATE USER keycloak WITH PASSWORD :keycloak_pwd;
GRANT CONNECT ON DATABASE keycloak TO keycloak;
GRANT keycloak_readwrite TO keycloak;

-- Setup roles

\connect gophoto

-- Create core device management role
DROP ROLE IF EXISTS core_readonly;

CREATE ROLE core_readonly;
GRANT SELECT ON ALL TABLES IN SCHEMA public TO core_readonly;

DROP ROLE IF EXISTS core_readwrite;

CREATE ROLE core_readwrite;
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO core_readwrite;

-- Create users
DROP USER IF EXISTS gophoto;

CREATE USER gophoto WITH PASSWORD :gophoto_pwd;
GRANT CONNECT ON DATABASE gophoto TO gophoto;
GRANT core_readwrite TO gophoto;
