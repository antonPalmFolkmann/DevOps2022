drop table if exists "user";
create table "user" (
  user_id serial primary key,
  username varchar(100) not null,
  email varchar(100) not null,
  pw_hash varchar(100) not null
);

drop table if exists follower;
create table follower (
  who_id integer,
  whom_id integer
);

drop table if exists message;
create table message (
  message_id serial primary key,
  author_id integer not null,
  text varchar(5000) not null,
  pub_date integer,
  flagged integer
);

COPY follower
FROM '/restore/dumps/follower.csv'
DELIMITER ','
CSV HEADER;

COPY "message"
FROM '/restore/dumps/message.csv'
DELIMITER ','
CSV HEADER;

COPY "user"
FROM '/restore/dumps/user.csv'
DELIMITER ','
CSV HEADER;