/*
 * Copyright © 2022 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

CREATE TABLE users
(
    id SERIAL PRIMARY KEY,
    username character varying(50) NOT NULL,
    password character varying(100) NOT NULL,
    email character varying(355) NOT NULL,
    created_on timestamp without time zone,
    last_login timestamp without time zone,
    image character varying DEFAULT 'https://static.productionready.io/images/smiley-cyrus.jpg'::character varying NOT NULL,
    bio character varying(280)
);

CREATE TABLE articles
(
    id SERIAL PRIMARY KEY,
    author_id integer NOT NULL,
    title character varying(255) NOT NULL,
    description character varying(255) NOT NULL,
    body character varying(255) NOT NULL,
    slug character varying(255) NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    tags character varying(255)
);

CREATE TABLE article_comments
(
    id SERIAL PRIMARY KEY,
    article_id integer NOT NULL,
    author_image character varying(255) NOT NULL,
    author_id integer NOT NULL,
    author_username character varying(255) NOT NULL,
    body character varying(255) NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);

CREATE TABLE articles_favorites
(
    id SERIAL PRIMARY KEY,
    article_id integer NOT NULL,
    user_id integer NOT NULL,
    value boolean NOT NULL
);

CREATE TABLE sessions
(
    id SERIAL PRIMARY KEY,
    user_id integer NOT NULL,
    session_one character varying(255) NOT NULL,
    session_two character varying(255) NOT NULL,
    CONSTRAINT fk_user
        FOREIGN KEY(user_id)
            REFERENCES users(id)
);
