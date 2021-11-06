BEGIN;

ALTER TABLE album ADD COLUMN thumbnail varchar(200);

COMMIT;
