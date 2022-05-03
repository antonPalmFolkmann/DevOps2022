drop table if exists users cascade;
create table users (
  id serial primary key,
  username varchar(100) not null,
  email varchar(100) not null,
  pw_hash varchar(5000) not null
);

drop table if exists follows cascade;
create table follows (
  user_id integer,
  whom_id integer,

  primary key (user_id, whom_id)
);

drop table if exists messages cascade;
create table messages (
  id serial primary key,
  user_id integer not null,
  text varchar(5000) not null,
  pub_date integer,
  flagged integer
);

COPY follows
FROM '/restore/dumps/followers.csv'
DELIMITER ','
CSV HEADER;

COPY messages 
FROM '/restore/dumps/messages.csv'
DELIMITER ','
CSV HEADER;

COPY users
FROM '/restore/dumps/users.csv'
DELIMITER ','
CSV HEADER;

BEGIN;
LOCK TABLE users IN EXCLUSIVE MODE;
SELECT setval('users_id_seq', COALESCE((SELECT MAX(id)+1 FROM users), 1), false);
COMMIT;

BEGIN;
LOCK TABLE messages IN EXCLUSIVE MODE;
SELECT setval('messages_id_seq', COALESCE((SELECT MAX(id)+1 FROM messages), 1), false);
COMMIT;

