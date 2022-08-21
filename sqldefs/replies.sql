CREATE TABLE public.replies (
	thread bigserial NOT NULL,
	poster _int8 NULL,
	"content" _text NULL,
	votes _int4 NULL
);
