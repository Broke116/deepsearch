-- create sequence

CREATE SEQUENCE public.textual_id_seq
    INCREMENT 1
    START 1
    MINVALUE 1
    MAXVALUE 2147483647
    CACHE 1;

ALTER SEQUENCE public.textual_id_seq
    OWNER TO admin;

-- table create script
CREATE TABLE public.textual
(
    id integer NOT NULL DEFAULT nextval('textual_id_seq'::regclass),
    title text COLLATE pg_catalog."default",
    content text COLLATE pg_catalog."default",
    tsv tsvector,
    CONSTRAINT textual_pkey PRIMARY KEY (id)
)

TABLESPACE pg_default;

ALTER TABLE public.textual
    OWNER to admin;

-- Index: ix_textual_tsv

-- DROP INDEX public.ix_textual_tsv;

CREATE INDEX ix_textual_tsv
    ON public.textual USING gin
    (tsv)
    TABLESPACE pg_default;

-- Trigger: tsv_vector_insert_trigger

-- DROP TRIGGER tsv_vector_insert_trigger ON public.textual;

-- creating a trigger
CREATE TRIGGER tsv_vector_insert_trigger
    BEFORE INSERT
    ON public.textual
    FOR EACH ROW
    EXECUTE PROCEDURE public.tsv_vector_number_update();

-- creating a postgres function
CREATE OR REPLACE FUNCTION public.tsv_vector_number_update()
  RETURNS trigger AS
$BODY$
BEGIN
   NEW.tsv := setweight(to_tsvector(new.title), 'A') ||
    setweight(to_tsvector(new.content), 'B');
 
   RETURN NEW;
END;
$BODY$
LANGUAGE plpgsql;

-- full text search queries
SELECT
  title, content
FROM
  public.textual, plainto_tsquery('lorem amet') q -- plainto_tsquery is used for creating a tsquery value from simple text by joining lexemes with and operator.
WHERE
  tsv @@ q;

SELECT
  title, ts_headline(content, q) -- returned keywords inside the text will be highlighted
FROM
  public.textual, plainto_tsquery('amet') q
WHERE
  tsv @@ q;

-- below query has two parts. inner query returns the top 10 results

SELECT
  title, ts_headline(content, q) -- highlighting matched results
FROM (
  SELECT
    title, content, ts_rank(tsv, q) as rank, q -- ranking the returned queries
  FROM
    public.textual, plainto_tsquery('lorem amet') q -- keywords to be searched.
  WHERE
    tsv @@ q
  ORDER BY
    rank DESC
  LIMIT
    10
)
ORDER BY
  rank DESC;