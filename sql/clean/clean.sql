-- Close connections
SELECT *, pg_terminate_backend(pid) FROM pg_stat_activity WHERE pid <> pg_backend_pid() AND datname = 'paperless';

-- Clean databases
DROP DATABASE IF EXISTS paperless;

-- Clean table roles 
DROP ROLE IF EXISTS core_readwrite;
DROP ROLE IF EXISTS core_readonly;

-- Clean user roles 
DROP USER IF EXISTS paperless_service;
