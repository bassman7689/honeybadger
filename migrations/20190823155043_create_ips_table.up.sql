CREATE TABLE ips (
	id uuid PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
	ip text UNIQUE NOT NULL
);
