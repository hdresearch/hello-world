CREATE TABLE IF NOT EXISTS visits (
    singleton BOOLEAN PRIMARY KEY DEFAULT TRUE CHECK (singleton),
    count BIGINT NOT NULL DEFAULT 0 CHECK (count >= 0)
);

INSERT INTO visits (singleton, count)
VALUES (TRUE, 0)
ON CONFLICT (singleton) DO NOTHING;
