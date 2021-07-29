do $$
BEGIN

TRUNCATE TABLE album, album_user_permissions, album_group_permissions CASCADE;


-- insert some albums
INSERT INTO album (id, name, created_at, owner_id, description, location) VALUES
(1, 'album1', '2021-01-01 12:00:00', 'user1', 'album1', NULL),
(2, 'album2', '2021-02-01 12:00:00', 'user1', 'album2', 'craiova'),
(3, 'album3', '2021-02-11 12:00:00', 'user1', 'album3', 'oltenia'),
(4, 'album4', '2021-04-21 12:00:00', 'user1', 'album4', NULL),
(5, 'album5', '2021-01-11 12:00:00', 'user2', 'album5', NULL),
(6, 'album6', '2021-07-01 12:00:00', 'user2', 'album6', 'craiova'),
(7, 'album7', '2021-08-11 12:00:00', 'user2', 'album7', 'oltenia'),
(8, 'album8', '2021-09-21 12:00:00', 'user2', 'album8', NULL),
(9, 'jane1', '2021-08-11 12:00:00','user4', 'jane', 'oltenia'),
(10, 'jane2', '2021-09-21 12:00:00', 'user4', 'jane', NULL);

INSERT INTO album_user_permissions (user_id, album_id, permissions) VALUES
('user3', 1, '{album.read, album.write}'),
('user3', 2, '{album.read}'),
('user3', 3, '{album.read, album.write, album.edit}'),
('user4', 3, '{album.read}'),
('user4', 5, '{album.read, album.write}'),
('user4', 6, '{album.read}'),
('user3', 7, '{album.read, album.write, album.edit}'),
('user2', 7, '{album.read, album.write}'),
('user2', 9, '{album.read, album.write}'),
('user2', 10, '{album.read, album.write}');

INSERT INTO album_group_permissions (group_name, album_id, permissions) VALUES
('group1', 1, '{album.read, album.write}'),
('group1', 2, '{album.read, album.write}'),
('group1', 3, '{album.read, album.write}'),
('group1', 4, '{album.read, album.write}'),
('group3', 1, '{album.read}'),
('group2', 1, '{album.read, album.write}');


END;
$$;
