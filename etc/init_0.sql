create schema if not exists af;

create table af.users
(
	userid bigserial not null constraint users_pk primary key,
	name varchar(255) not null,
	surname varchar(255) not null,
	email varchar(255) not null,
	created date not null DEFAULT '2020-12-10',
	modified timestamp not null DEFAULT now()
);

create unique index email_unique on af.users (email);

create table af.albums
(
	album_id bigserial not null constraint albums_pk primary key,
	album_name varchar(255),
	album_owner bigint constraint albums_users__fk references af.users,
	created date not null DEFAULT '2020-12-10',
	modified timestamp not null DEFAULT now()
);

create table af.images
(
	image_id bigserial constraint images_pk primary key,
	image_path varchar(255) not null,
	image_name varchar(255) not null,
	created date not null DEFAULT '2020-12-10',
	modified timestamp not null DEFAULT now(),
    user_id bigint constraint users_fk references af.users
);

create table af.images_to_albums
(
	id bigserial constraint images_albums_pk primary key,
	album_id bigint constraint albums__fk references af.albums,
	image_id bigint constraint images__fk references af.images
);

create table af.users_audit (
   id bigserial constraint users_audit_pk primary key,
   user_id bigint NOT NULL,
   name varchar(250) NOT NULL,
   new_name varchar(250) NOT NULL,
   surname varchar(250) NOT NULL,
   new_surname varchar(250) NOT NULL,
   email varchar(250) NOT NULL,
   new_email varchar(250) NOT NULL,
   changed_on timestamp(6) NOT NULL
);

alter table af.users_audit
	add constraint users_audit_users_userid_fk
		foreign key (user_id) references af.users;

create table af.hashtags
(
	tag_id bigserial constraint hashtags_pk primary key,
	tag_name VARCHAR(250) NOT NULL,
	created timestamp not null
);

create table af.hashtags_to_images
(
    id bigserial constraint hashtags_to_images_pk primary key,
	tag_id bigint constraint hashtags__fk references af.hashtags,
	image_id bigint constraint images__fk references af.images
);

create table af.users_to_hashtags
(
	id bigserial constraint users_to_hashtags_pk primary key,
	user_id bigint not null,
	tag_id bigint not null
);

alter table af.users_to_hashtags
	add constraint users_to_hashtags_hashtags_tag_id_fk
		foreign key (tag_id) references af.hashtags;

alter table af.users_to_hashtags
	add constraint users_to_hashtags_users_userid_fk
		foreign key (user_id) references af.users;

alter table af.users owner to marius;
alter table af.albums owner to marius;
alter table af.images owner to marius;
alter table af.images_to_albums owner to marius;
alter table af.users_audit owner to marius;
alter table af.hashtags owner to marius;
alter table af.hashtags_to_images owner to marius;
alter table af.users_to_hashtags owner to marius;

-- Trigger pentru scrierea in tabela de users_audit a modificarilor inregistrate pe tabela de users, campurile email, name si surname
CREATE OR REPLACE FUNCTION af.log_user_changes()
  RETURNS TRIGGER
  LANGUAGE PLPGSQL
  AS
$$
BEGIN
    IF NEW.email <> OLD.email THEN
        INSERT INTO af.users_audit(user_id,name,new_name,surname, new_surname, email, new_email, changed_on) VALUES (OLD.userid, OLD.name, NEW.name, OLD.surname, NEW.surname, OLD.email, NEW.email, now());
    END IF;

    IF NEW.name <> OLD.name THEN
        INSERT INTO af.users_audit(user_id,name,new_name,surname, new_surname, email, new_email, changed_on) VALUES (OLD.userid, OLD.name, NEW.name, OLD.surname, NEW.surname, OLD.email, NEW.email, now());
    END IF;

    IF NEW.surname <> OLD.surname THEN
        INSERT INTO af.users_audit(user_id,name,new_name,surname, new_surname, email, new_email, changed_on) VALUES(OLD.userid, OLD.name, NEW.name, OLD.surname, NEW.surname, OLD.email, NEW.email, now());
    END IF;
    RETURN NEW;
END;
$$;

CREATE TRIGGER assign_role
  BEFORE UPDATE
  ON af.users
  FOR EACH ROW
  EXECUTE PROCEDURE af.log_user_changes();


-- Inseram 3 useri
INSERT INTO af.users (userid, name, surname, email, created, modified) VALUES (1, 'Marius', 'Costache', 'marius.costache.b@gmail.com', '2020-12-10', '2020-12-09 23:02:57.310822');
INSERT INTO af.users (userid, name, surname, email, created, modified) VALUES (2, 'Ion', 'Tiriac', 'ion@tiriac.ro', '2020-12-10', '2020-12-09 23:03:27.916636');
INSERT INTO af.users (userid, name, surname, email, created, modified) VALUES (3, 'John', 'Cena', 'hisnameis@johncena.com', '2020-12-10', '2020-12-10 23:03:27.916636');

-- Inseram 10 imagini
INSERT INTO af.images (image_id, image_path, image_name, created, modified, user_id) VALUES (1, '/images/cat.jpg', 'Pretty cat', '2020-12-10', '2020-12-09 23:18:55.281761', 1);
INSERT INTO af.images (image_id, image_path, image_name, created, modified, user_id) VALUES (2, '/images/dog.jpg', 'Cute dog', '2020-12-10', '2020-12-09 23:20:30.672537', 1);
INSERT INTO af.images (image_id, image_path, image_name, created, modified, user_id) VALUES (3, '/images/squirrel.jpg', 'Jack the squirrel', '2020-12-10', '2020-12-09 23:21:26.715517', 1);
INSERT INTO af.images (image_id, image_path, image_name, created, modified, user_id) VALUES (4, '/images/marius.jpg', 'MArius portret', '2020-12-10', '2020-12-09 23:28:24.208036', 1);
INSERT INTO af.images (image_id, image_path, image_name, created, modified, user_id) VALUES (5, '/images/iontiriac.jpg', 'Ion tiriac portret', '2020-12-10', '2020-12-09 23:28:43.503305', 2);
INSERT INTO af.images (image_id, image_path, image_name, created, modified, user_id) VALUES (6, '/images/motorcycle.jpg', 'Honda CBR motorcycle', '2020-12-10', '2020-12-09 23:29:19.047579', 2);
INSERT INTO af.images (image_id, image_path, image_name, created, modified, user_id) VALUES (7, '/images/motorcycle2.jpg', 'Ducatti motorcycle', '2020-12-10', '2020-12-09 23:29:50.686807', 2);
INSERT INTO af.images (image_id, image_path, image_name, created, modified, user_id) VALUES (8, '/images/motorcycle2.jpg', 'BMW motorcycle', '2020-12-10', '2020-12-09 23:29:50.686807', 2);
INSERT INTO af.images (image_id, image_path, image_name, created, modified, user_id) VALUES (9, '/images/motorcycle3.jpg', 'Aprilia', '2020-12-10', '2020-12-09 23:29:50.686807', 2);
INSERT INTO af.images (image_id, image_path, image_name, created, modified, user_id) VALUES (10, '/images/motorcycle4.jpg', 'KTM 690 SMC R', '2020-12-10', '2020-12-10 23:29:50.686807', 2);

-- Inseram 4 albume
INSERT INTO af.albums (album_id, album_name, album_owner, created, modified) VALUES (1, 'motorcycles', 2, '2020-12-10', '2020-12-09 23:34:56.089792');
INSERT INTO af.albums (album_id, album_name, album_owner, created, modified) VALUES (2, 'Animals', 1, '2020-12-10', '2020-12-09 23:35:05.527477');
INSERT INTO af.albums (album_id, album_name, album_owner, created, modified) VALUES (3, 'Random images', 1, '2020-12-10', '2020-12-09 23:35:24.953676');
INSERT INTO af.albums (album_id, album_name, album_owner, created, modified) VALUES (4, 'Portraits', 1, '2020-12-10', '2020-12-09 23:35:24.953676');

-- Legam albumele la imagini
INSERT INTO af.images_to_albums (id, album_id, image_id)VALUES (1, 2, 1);
INSERT INTO af.images_to_albums (id, album_id, image_id)VALUES (2, 2, 2);
INSERT INTO af.images_to_albums (id, album_id, image_id)VALUES (3, 2, 3);
INSERT INTO af.images_to_albums (id, album_id, image_id)VALUES (4, 3, 4);
INSERT INTO af.images_to_albums (id, album_id, image_id)VALUES (5, 4, 4);

-- Inseram 3 hashtags
INSERT INTO af.hashtags (tag_name, created) VALUES ('motorcycles', now());
INSERT INTO af.hashtags (tag_name, created) VALUES ('animals', now());
INSERT INTO af.hashtags (tag_name, created) VALUES ('portraits', now());
INSERT INTO af.hashtags (tag_name, created) VALUES ('school', now());

-- legam imagini la hashtags
INSERT INTO af.hashtags_to_images (tag_id, image_id) VALUES (1, 6);
INSERT INTO af.hashtags_to_images (tag_id, image_id) VALUES (1, 7);
INSERT INTO af.hashtags_to_images (tag_id, image_id) VALUES (1, 8);
INSERT INTO af.hashtags_to_images (tag_id, image_id) VALUES (1, 9);
INSERT INTO af.hashtags_to_images (tag_id, image_id) VALUES (2, 1);
INSERT INTO af.hashtags_to_images (tag_id, image_id) VALUES (2, 2);
INSERT INTO af.hashtags_to_images (tag_id, image_id) VALUES (2, 3);
INSERT INTO af.hashtags_to_images (tag_id, image_id) VALUES (3, 4);
INSERT INTO af.hashtags_to_images (tag_id, image_id) VALUES (3, 5);

-- legam useri la hashtags pentru a putea urmari hashtags publice
INSERT INTO af.users_to_hashtags (user_id, tag_id) VALUES (1, 3);
INSERT INTO af.users_to_hashtags (user_id, tag_id) VALUES (1, 4);
INSERT INTO af.users_to_hashtags (user_id, tag_id) VALUES (2, 3);

-- SELECTEAZA TOATE IMAGINILE CARE APARTIN USERULUI 1 DAR CARE APARTIN UNUI ALBUM
SELECT a.image_id, a.image_name, c.album_name FROM af.images a
INNER JOIN af.images_to_albums b on a.image_id = b.image_id
INNER JOIN af.albums c on c.album_id = b.album_id
WHERE a.user_id = 1;

-- CREARE VIEW cu userii care nu au imagini adaugate
CREATE VIEW af.materialized_users_no_images AS
    SELECT u.userid, u.email, u.created FROM af.users u
    LEFT JOIN af.images i on i.user_id = u.userid
    WHERE i.user_id IS NULL;

-- CREARE VIEW imaginile care apartin tagurilor urmarite de userul 1.
CREATE VIEW af.materialized_images_from_hashtags_followed_by_userid1 AS
    SELECT i.image_id, i.image_path, i.image_name, h.tag_name, u.name FROM af.hashtags h
    LEFT JOIN af.hashtags_to_images hti on h.tag_id = hti.tag_id
    LEFT JOIN af.images i on hti.image_id = i.image_id
    LEFT JOIN af.users_to_hashtags uth on uth.tag_id = h.tag_id
    LEFT JOIN af.users u on uth.user_id = u.userid
    WHERE uth.user_id = 1;

--UPDATE pentru o inregistrare
UPDATE af.albums SET album_name = 'superbikes' WHERE album_id = 1;

--DELETE pentru o inregistrare
DELETE FROM af.images WHERE image_id = 10;

--Actualizam un user pentru a proba triggerul
UPDATE af.users SET email = 'marius.costache@icloud.com' WHERE af.users.userid = 1;

--Afisam toate rezultatele din tabela users_audit
SELECT * FROM af.users_audit;