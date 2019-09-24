CREATE TABLE ts_connections (
	id uuid PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
	time TIMESTAMPTZ NOT NULL,
    value REAL NOT NULL
);
