BEGIN TRANSACTION;

ALTER TABLE short_links
    ADD COLUMN deleted BOOLEAN NOT NULL DEFAULT FALSE;

COMMIT;