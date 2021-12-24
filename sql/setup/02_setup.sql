BEGIN;

CREATE TYPE role as ENUM('admin','editor','user');

DROP TABLE IF EXISTS "album";
DROP TABLE IF EXISTS "bucket";
DROP TABLE IF EXISTS "album_user_permissions";
DROP TABLE IF EXISTS "album_group_permissions";
DROP TABLE IF EXISTS "tag";
DROP TABLE IF EXISTS "albums_tags";

CREATE TABLE album (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT (now() AT TIME ZONE 'UTC') NOT NULL,
    owner_id TEXT NOT NULL,
    bucket TEXT NOT NULL,
    description TEXT,
    location TEXT
);

CREATE TYPE permission_id as ENUM (
    'album.read',
    'album.write',
    'album.edit',
    'album.delete'
);

CREATE TABLE album_user_permissions (
    user_id TEXT NOT NULL,
    album_id SERIAL REFERENCES album(id),
    permissions permission_id[] NOT NULL,
    CONSTRAINT album_user_permissions_pk PRIMARY KEY (
        user_id,
        album_id
    )
);

CREATE INDEX user_id_idx ON album_user_permissions  (user_id);

CREATE TABLE album_group_permissions (
    group_name TEXT NOT NULL,
    album_id SERIAL REFERENCES album(id),
    permissions permission_id[] NOT NULL,
    CONSTRAINT album_group_permissions_pk PRIMARY KEY (
        group_name,
        album_id
    )
);

CREATE INDEX group_name_idx ON album_group_permissions (group_name);

CREATE TABLE tag (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    user_id TEXT NOT NULL,
    color TEXT 
);

CREATE INDEX tag_user_id_idx ON tag (user_id);

CREATE TABLE albums_tags (
    album_id SERIAL REFERENCES album(id),
    tag_id SERIAL REFERENCES tag(id),
    CONSTRAINT albums_tags_pk PRIMARY KEY (
        album_id,
        tag_id
    )
);

COMMIT;
