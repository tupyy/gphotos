INSERT INTO users VALUES (
    DEFAULT,
    'batman',
    'admin',
    'Wayne',
    'Wayne',
    'batman@gotham.com',
    TRUE
);


INSERT INTO users VALUES (
    DEFAULT,
    'superman',
    'user',
    'superman',
    'superman',
    'superman@gotham.com',
    TRUE
);

INSERT INTO groups VALUES (
    1,
    'admins'
);

INSERT INTO groups VALUES (
    2,
    'users'
);

INSERT INTO groups VALUES (
    3,
    'editors'
);

INSERT INTO users_groups VALUES (
    1,
    1
);

INSERT INTO users_groups VALUES (
    1,
    2
);

INSERT INTO users_groups VALUES (
    2,
    2
);
