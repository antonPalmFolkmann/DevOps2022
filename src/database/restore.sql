drop table if exists users cascade;
create table "users" (
  user_id serial primary key,
  username varchar(100) not null,
  email varchar(100) not null,
  pw_hash varchar(100) not null
);

drop table if exists followers cascade;
create table follower (
  user_id integer,
  whom_id integer,

  primary key (user_id, whom_id)
);

drop table if exists messages cascade;
create table messages (
  message_id serial primary key,
  user_id integer not null,
  text varchar(5000) not null,
  pub_date integer,
  flagged integer
);

COPY followers
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