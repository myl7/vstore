create table sources (
  sid serial primary key,
  name varchar(120) not null
);
create table videos (
  vid serial primary key,
  sid int not null references sources (sid),
  title varchar(120) not null,
  description varchar(510) not null,
  file oid not null
);
