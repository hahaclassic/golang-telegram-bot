CREATE TABLE users (
    id BIGINT NOT NULL,
    username VARCHAR(255) NOT NULL,
    premium BOOLEAN NOT NULL,
    lang VARCHAR(3) NOT NULL,
    PRIMARY KEY ( id )
);

CREATE TABLE folders (
    id VARCHAR(12),
    name TEXT NOT NULL,
    PRIMARY KEY ( id )
);

CREATE TABLE pages (
    url TEXT,
    tag VARCHAR(60),
    folder_id VARCHAR(12) REFERENCES folders (id) ON DELETE CASCADE
);

CREATE TABLE passwords (
    password VARCHAR(8),
    access_lvl SMALLINT,
    folder_id VARCHAR(12) REFERENCES folders (id) ON DELETE CASCADE,
    PRIMARY KEY(folder_id, access_lvl)
);

CREATE TABLE users_folders (
    user_id BIGINT REFERENCES users (id) ON DELETE CASCADE,
    folder_id VARCHAR(12) REFERENCES folders (id) ON DELETE CASCADE,
    access_lvl SMALLINT,

    FOREIGN KEY (folder_id, access_lvl)
        REFERENCES passwords (folder_id, access_lvl)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION
);
