-- Close connections
SELECT *, pg_terminate_backend(pid) FROM pg_stat_activity WHERE pid <> pg_backend_pid() AND datname = 'gophoto';

-- Clean databases
DROP DATABASE IF EXISTS gophoto;

-- Clean table roles 
DROP ROLE IF EXISTS core_readwrite;
DROP ROLE IF EXISTS core_readonly;
DROP USER IF EXISTS resources_admin;

