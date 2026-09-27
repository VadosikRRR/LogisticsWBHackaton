ALTER TABLE raw_records
ADD COLUMN IF NOT EXISTS target_2h DOUBLE PRECISION;

UPDATE raw_records
SET target_2h = 0
WHERE target_2h IS NULL;

ALTER TABLE raw_records
ALTER COLUMN target_2h SET NOT NULL;
