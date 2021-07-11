CREATE TYPE role as ENUM('admin','editor','user');

CREATE TABLE groups (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL
);

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username TEXT NOT NULL,
    role role NOT NULL,
    user_id TEXT NOT NULL,
    can_share BOOLEAN
);

CREATE TABLE users_groups (
    users_id SERIAL REFERENCES users(id),
    groups_id SERIAL REFERENCES groups(id),
    CONSTRAINT users_groups_pk PRIMARY KEY (
        users_id,
        groups_id
    )
);

CREATE TABLE album (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT (now() AT TIME ZONE 'UTC') NOT NULL,
    owner_id SERIAL REFERENCES users(id),
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
    user_id SERIAL REFERENCES users(id),
    album_id SERIAL REFERENCES album(id),
    permissions permission_id[] NOT NULL,
    CONSTRAINT album_user_permissions_pk PRIMARY KEY (
        user_id,
        album_id
    )
);

CREATE TABLE album_group_permissions (
    "group_id" SERIAL REFERENCES groups(id),
    album_id SERIAL REFERENCES album(id),
    permissions permission_id[] NOT NULL,
    CONSTRAINT album_group_permissions_pk PRIMARY KEY (
        "group_id",
        album_id
    )
);

CREATE TABLE tag (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL
);

CREATE TABLE albums_tags (
    album_id SERIAL REFERENCES album(id),
    tag_id SERIAL REFERENCES tag(id),
    CONSTRAINT albums_tags_pk PRIMARY KEY (
        album_id,
        tag_id
    )
);

CREATE TABLE token_blacklist (
    id SERIAL PRIMARY KEY,
    token TEXT NOT NULL
)
