BEGIN TRANSACTION;

ALTER TABLE short_links
    DROP CONSTRAINT short_links_short_url_key;

ALTER TABLE short_links
    ADD COLUMN deleted BOOLEAN NOT NULL DEFAULT FALSE;

CREATE UNIQUE INDEX unique_short_url_when_not_deleted
    ON short_links (short_url)
    WHERE deleted = FALSE;

COMMIT;