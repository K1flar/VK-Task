CREATE TABLE actors(
	id SERIAL PRIMARY KEY,
	full_name VARCHAR NOT NULL,
	gender VARCHAR(10) CHECK(gender IN ('male', 'female')) NOT NULL,
	birthday DATE NOT NULL
);

CREATE TABLE films(
	id SERIAL PRIMARY KEY,
	name VARCHAR(150) CHECK(length(name)>0) UNIQUE NOT NULL,
	description VARCHAR(1000),
	release_date DATE NOT NULL,
	rating SMALLINT CHECK(rating BETWEEN 0 AND 10) NOT NULL
);

CREATE TABLE film_actor(
	film_id INTEGER REFERENCES films(id) ON DELETE CASCADE NOT NULL,
	actor_id INTEGER REFERENCES actors(id) ON DELETE CASCADE NOT NULL ,
	PRIMARY KEY(film_id, actor_id)
);

CREATE TABLE users(
	id SERIAL PRIMARY KEY,
	login VARCHAR UNIQUE NOT NULL,
	password VARCHAR NOT NULL,
	role VARCHAR(20) CHECK(role IN ('viewer', 'admin')) NOT NULL
);