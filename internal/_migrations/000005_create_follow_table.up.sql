BEGIN;

CREATE TABLE IF NOT EXISTS followings (
    followee_id int not null,
    follower_id int not null,
    followed_on timestamptz not null default now(),
    primary key (followee_id, follower_id),
    constraint fk_following foreign key(followee_id) references users(id),
    constraint fk_follower foreign key(follower_id) references users(id)
);

COMMIT;
