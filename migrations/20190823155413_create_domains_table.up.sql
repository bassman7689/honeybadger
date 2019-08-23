CREATE TABLE domains (
	id uuid PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
	domain text UNIQUE NOT NULL
);
