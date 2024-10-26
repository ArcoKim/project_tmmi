CREATE TABLE album (
	id SERIAL PRIMARY KEY,
	name VARCHAR(100) NOT NULL,
	artist VARCHAR(50) NOT NULL,
	image VARCHAR(200) NOT NULL
);

CREATE TABLE music (
	id SERIAL PRIMARY KEY,
 	name VARCHAR(100) NOT NULL,
	album_id INTEGER NOT NULL,
	lyrics TEXT NOT NULL,
	melon VARCHAR(15) UNIQUE,
	genie VARCHAR(15) UNIQUE,
	flo VARCHAR(15) UNIQUE,
	bugs VARCHAR(15) UNIQUE,
	vibe VARCHAR(15) UNIQUE,
	CONSTRAINT fk_album_id FOREIGN KEY(album_id) REFERENCES album(id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TYPE chart_type AS ENUM ('melon', 'genie', 'flo', 'bugs', 'vibe');
CREATE TABLE ranks (
	id SERIAL PRIMARY KEY,
	ranks smallint NOT NULL,
	crdate DATE NOT NULL,
	types chart_type NOT NULL,
	music_id INTEGER NOT NULL,
	CONSTRAINT fk_music_id FOREIGN 	KEY(music_id) REFERENCES music(id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE EXTENSION pg_trgm;
CREATE EXTENSION aws_s3 CASCADE;

CREATE INDEX music_name_idx ON music USING GIST (name gist_trgm_ops);
CREATE INDEX album_name_idx ON album USING GIST (name gist_trgm_ops);
CREATE INDEX album_artist_idx ON album USING GIST (artist gist_trgm_ops);