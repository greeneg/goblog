--
-- File generated with SQLiteStudio v3.4.4 on Sun Jan 28 15:45:10 2024
--
-- Text encoding used: UTF-8
--
PRAGMA foreign_keys = off;
BEGIN TRANSACTION;

-- Table: blog
DROP TABLE IF EXISTS blog;

CREATE TABLE IF NOT EXISTS blog (
    id      INTEGER  PRIMARY KEY AUTOINCREMENT,
    author  TEXT,
    title   TEXT,
    image   BLOB,
    content TEXT,
    ctime   DATETIME
);


COMMIT TRANSACTION;
PRAGMA foreign_keys = on;
