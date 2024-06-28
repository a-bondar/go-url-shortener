BEGIN TRANSACTION;

CREATE TABLE short_links
(
    id           INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    short_url    VARCHAR(255) NOT NULL UNIQUE,
    original_url VARCHAR(255) NOT NULL UNIQUE
);

COMMIT;