BEGIN;

DROP TABLE IF EXISTS "album";
DROP TABLE IF EXISTS "bucket";
DROP TABLE IF EXISTS "album_user_permissions";
DROP TABLE IF EXISTS "album_group_permissions";
DROP TABLE IF EXISTS "tag";
DROP TABLE IF EXISTS "albums_tags";

CREATE TYPE role as ENUM('admin','editor','user');

CREATE TABLE album (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT (now() AT TIME ZONE 'UTC') NOT NULL,
    owner_id TEXT NOT NULL,
    bucket TEXT NOT NULL,
    description TEXT,
    location TEXT,
    thumbnail VARCHAR(200)
);

CREATE TYPE permission_id as ENUM (
    'album.read',
    'album.write',
    'album.edit',
    'album.delete'
);

CREATE TYPE owner_kind as ENUM (
    'user',
    'group'
);

CREATE TABLE album_permissions (
    owner_id TEXT NOT NULL,
    owner_kind owner_kind NOT NULL,
    album_id TEXT REFERENCES album(id) ON DELETE CASCADE,
    permissions permission_id[] NOT NULL,
    CONSTRAINT album_user_permissions_pk PRIMARY KEY (
        owner_id,
        album_id
    )
);

CREATE TABLE tag (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    color TEXT 
    owner_id TEXT NOT NULL,
);

CREATE INDEX tag_user_id_idx ON tag (user_id);

CREATE TABLE albums_tags (
    album_id TEXT REFERENCES album(id) ON DELETE CASCADE,
    tag_id TEXT REFERENCES tag(id) ON DELETE CASCADE,
    CONSTRAINT albums_tags_pk PRIMARY KEY (
        album_id,
        tag_id
    ) 
);

COMMIT;
