create table sources (
  sid serial primary key,
  name varchar(120) not null
);
create table users (
  uid serial primary key,
  name varchar(60) not null,
  token varchar(120) not null
);
create table videos (
  vid serial primary key,
  sid int not null references sources (sid),
  uid int not null references users (uid),
  title varchar(120) not null,
  description varchar(510) not null,
  file oid not null
);
