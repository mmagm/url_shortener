CREATE SEQUENCE links_id_seq
  INCREMENT 1
  START 1
  MINVALUE 1
  MAXVALUE 9223372036854775807
  CACHE 1;

CREATE TABLE links
(
  id integer NOT NULL DEFAULT nextval('links_id_seq'::regclass),
  url character varying COLLATE pg_catalog."default",
  token character varying COLLATE pg_catalog."default",
  created_at timestamp without time zone NOT NULL,
  updated_at timestamp without time zone NOT NULL,
  CONSTRAINT links_pkey PRIMARY KEY (id)
);

CREATE INDEX links_token_idx ON links(token);
