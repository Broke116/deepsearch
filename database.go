package main

const insertText = `
INSERT INTO public.textual (title, content)
VALUES ($1, $2)
`

const fullTextSearch = `
SELECT
  title, ts_headline(content, q) -- highlighting matched results
FROM (
  SELECT
    title, content, ts_rank(tsv, q) as rank, q -- ranking the returned queries
  FROM
    public.textual, plainto_tsquery($1) q -- keywords to be searched.
  WHERE
    tsv @@ q
  ORDER BY
    rank DESC
  LIMIT
    10
) inner_sql
ORDER BY
  rank DESC;
`