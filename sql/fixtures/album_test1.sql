do $$
BEGIN

TRUNCATE TABLE users, groups, users_groups, album CASCADE;

(1,'admins'),
(2,'users'),
(3,'editors');

INSERT INTO users_groups VALUES 
(1,1),
(1,2),
(2,2),
(3,2),
(4,2);

END;
$$;
