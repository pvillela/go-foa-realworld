BEGIN;

CREATE EXTENSION IF NOT EXISTS citext;

CREATE TABLE IF NOT EXISTS users
(
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL UNIQUE,
    email CITEXT NOT NULL UNIQUE,
--     email VARCHAR(255) NOT NULL UNIQUE,
    bio TEXT,
    image TEXT,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMIT;
BEGIN;

CREATE TABLE IF NOT EXISTS articles
(
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    body TEXT NOT NULL,
    description TEXT,
    slug VARCHAR(255) NOT NULL UNIQUE,
    author_id  INT NOT NULL,
    favorites_count numeric NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_author FOREIGN KEY(author_id) REFERENCES users(id)
);

COMMIT;
BEGIN;

CREATE TABLE IF NOT EXISTS tags
(
    id serial primary KEY,
    name citext unique NOT NULL
--     name VARCHAR(255) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS article_tags
(
    article_id INT NOT NULL,
    tag_id int not null,
    primary key(article_id, tag_id),
    constraint fk_article foreign key(article_id) references articles(id) on delete cascade ,
    constraint fk_tag foreign key(tag_id) references tags(id) on delete cascade
);

COMMIT;
BEGIN;

CREATE TABLE IF NOT EXISTS followings
(
    followee_id int not null,
    follower_id int not null,
    followed_on timestamptz not null default now(),
    primary key (followee_id, follower_id),
    constraint fk_following foreign key(followee_id) references users(id),
    constraint fk_follower foreign key(follower_id) references users(id)
);

COMMIT;
BEGIN;

CREATE TABLE IF NOT EXISTS comments
(
    id serial primary key,
    article_id int not null,
    author_id int not null,
    body text not null,
    created_at timestamptz not null default now(),
    constraint fk_article foreign key(article_id) references articles(id) on delete cascade,
    constraint fk_author foreign key(author_id) references users(id) on delete cascade
);

COMMIT;
BEGIN;

CREATE TABLE IF NOT EXISTS favorites
(
    article_id int not null,
    user_id int not null,
    created_at timestamptz not null default now(),
    primary key(article_id, user_id),
    constraint fk_article foreign key(article_id) references articles(id) on delete cascade,
    constraint fk_user foreign key(user_id) references users(id) on delete cascade
);

COMMIT;
