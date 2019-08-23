CREATE TABLE addresses (
	id uuid PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
	address text UNIQUE NOT NULL
);
