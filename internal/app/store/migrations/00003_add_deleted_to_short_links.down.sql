BEGIN TRANSACTION;

DROP INDEX IF EXISTS unique_short_url_when_not_deleted;

ALTER TABLE short_links
    DROP COLUMN deleted;

ALTER TABLE short_links
    ADD CONSTRAINT short_links_short_url_key UNIQUE (short_url);

COMMIT;