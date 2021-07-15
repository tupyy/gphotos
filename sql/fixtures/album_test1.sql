do $$
BEGIN

TRUNCATE TABLE users, groups, users_groups, album CASCADE;


INSERT INTO users VALUES 
(1,'batman','admin','user1',TRUE),
(2,'superman','user','user2',TRUE),
(3,'joedoe','user','joe',TRUE),
(4,'jane','user','jane',TRUE);

INSERT INTO groups VALUES 
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
