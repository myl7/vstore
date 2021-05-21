create table sources (
  sid serial primary key,
  name varchar(120) not null
);
create table users (
  uid serial primary key,
  name varchar(60) not null,
  token varchar(120) not null
);
create index on users (name);
create table videos (
  vid serial primary key,
  sid int not null references sources (sid),
  uid int not null references users (uid),
  title varchar(120) not null,
  description varchar(510) not null,
  file oid not null
);
create table comments (
  mid serial primary key,
  vid int not null references videos (vid),
  uid int not null references users (uid),
  text varchar(510) not null,
  time timestamp not null
);
