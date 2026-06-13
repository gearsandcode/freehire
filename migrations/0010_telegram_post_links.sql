-- Preserve a post's outbound hyperlinks (anchor text + href) alongside its plain text.
-- The crawl parser drops hrefs from `text` (it stays clean for the LLM), so the links go
-- here as a JSON array of {text, url}. The extraction worker follows them to fetch full
-- vacancies from destination sites (e.g. career.habr.com) instead of the thin teaser.
--
-- Applied by Postgres on first volume init (same as 0001-0009) and read by sqlc. On an
-- existing volume this must be applied manually (no versioned migration runner yet).
ALTER TABLE telegram_posts
    ADD COLUMN IF NOT EXISTS links JSONB NOT NULL DEFAULT '[]';
