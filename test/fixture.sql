insert into users (uid, name, token)
values (1, 'User name', 'Token');
insert into sources (sid, name)
values (1, 'Source');
insert into videos (vid, sid, uid, title, description, file)
values (1, 1, 1, 'Title 1', 'Description 1', 1),
       (2, 1, 1, 'Title 2', 'Description 2', 2);
