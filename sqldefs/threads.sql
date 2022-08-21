CREATE TABLE public.threads (
	id bigserial NOT NULL,
	op int4 NOT NULL,
	title varchar NOT NULL,
	"content" text NOT NULL,
	pubd time NOT NULL,
	votes int4 NOT NULL DEFAULT 0
);
