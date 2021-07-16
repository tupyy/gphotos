do $$
BEGIN

TRUNCATE TABLE users, groups, users_groups, album CASCADE;


INSERT INTO users VALUES 
(1,'batman','admin','user1',TRUE),
(2,'cosmin','user','user2',TRUE),
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

-- insert some albums
INSERT INTO album (id, name, created_at, owner_id, description, location) VALUES
(1, 'album1', '2021-01-01 12:00:00', 1, 'album1', NULL),
(2, 'album2', '2021-02-01 12:00:00', 1, 'album2', 'craiova'),
(3, 'album3', '2021-02-11 12:00:00', 1, 'album3', 'oltenia'),
(4, 'album4', '2021-04-21 12:00:00', 1, 'album4', NULL),
(5, 'album5', '2021-01-11 12:00:00', 2, 'album5', NULL),
(6, 'album6', '2021-07-01 12:00:00', 2, 'album6', 'craiova'),
(7, 'album7', '2021-08-11 12:00:00', 2, 'album7', 'oltenia'),
(8, 'album8', '2021-09-21 12:00:00', 2, 'album8', NULL),
(9, 'jane1', '2021-08-11 12:00:00',4, 'jane', 'oltenia'),
(10, 'jane2', '2021-09-21 12:00:00', 4, 'jane', NULL);

INSERT INTO album_user_permissions (user_id, album_id, permissions) VALUES
(3, 1, '{album.read, album.write}'),
(3, 2, '{album.read}'),
(3, 3, '{album.read, album.write, album.edit}'),
(4, 3, '{album.read}'),
(4, 5, '{album.read, album.write}'),
(4, 6, '{album.read}'),
(4, 7, '{album.read, album.write, album.edit}'),
(3, 7, '{album.read, album.write}'),
(2, 9, '{album.read, album.write}'),
(2, 10, '{album.read, album.write}');

INSERT INTO album_group_permissions (group_id, album_id, permissions) VALUES
(1, 1, '{album.read, album.write}'),
(1, 2, '{album.read, album.write}'),
(1, 3, '{album.read, album.write}'),
(1, 4, '{album.read, album.write}'),
(2, 1, '{album.read}'),
(3, 1, '{album.read, album.write}');


END;
$$;
