BEGIN;
ALTER TABLE Sheets
    ADD COLUMN Public boolean;
UPDATE
    Sheets
SET
    Public = FALSE;
ALTER TABLE Sheets
    ALTER COLUMN Public SET NOT NULL;
-- Create review table
CREATE TABLE review_sheets (
    IdReview integer NOT NULL,
    IdSheet integer NOT NULL,
    Kind integer CHECK (Kind IN (0, 1, 2, 3)) NOT NULL
);
-- update kind constraint
ALTER TABLE reviews
    DROP CONSTRAINT reviews_kind_check;
ALTER TABLE reviews
    ADD CONSTRAINT reviews_kind_check CHECK (Kind IN (0, 1, 2, 3));
-- add constraints
ALTER TABLE review_sheets
    ADD FOREIGN KEY (IdReview, Kind) REFERENCES reviews (ID, Kind) ON DELETE CASCADE;
ALTER TABLE review_sheets
    ADD CHECK (Kind = 3
    /* ReviewKind.KSheet */);
ALTER TABLE review_sheets
    ADD UNIQUE (IdSheet);
ALTER TABLE review_sheets
    ADD UNIQUE (IdReview);
ALTER TABLE review_sheets
    ADD FOREIGN KEY (IdReview) REFERENCES reviews ON DELETE CASCADE;
ALTER TABLE review_sheets
    ADD FOREIGN KEY (IdSheet) REFERENCES sheets;
COMMIT;

