CREATE TABLE public.users (
	id bigserial NOT NULL,
	username varchar NOT NULL,
	"password" varchar NOT NULL,
	"role" varchar NOT NULL DEFAULT 0
);
