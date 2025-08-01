-- users definition

CREATE TABLE users (
	username TEXT NOT NULL,
	password TEXT,
	max_storage INTEGER DEFAULT (-1) NOT NULL,
	CONSTRAINT users_pk PRIMARY KEY (username)
);


-- files definition

CREATE TABLE files (
	file_id TEXT NOT NULL,
	filename TEXT NOT NULL,
	created_at TEXT DEFAULT (date()) NOT NULL,
	updated_at TEXT DEFAULT (date()) NOT NULL,
	username TEXT, directory TEXT DEFAULT ('/') NOT NULL,
	CONSTRAINT files_pk PRIMARY KEY (file_id),
	CONSTRAINT files_users_FK FOREIGN KEY (username) REFERENCES users(username) ON DELETE SET NULL ON UPDATE CASCADE
);

CREATE UNIQUE INDEX files_filename_IDX ON files (filename,directory);