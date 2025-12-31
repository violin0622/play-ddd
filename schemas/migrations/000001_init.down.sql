BEGIN;

DROP TABLE IF EXISTS public.novels;

ALTER TABLE public.events
DROP COLUMN reason,
DROP COLUMN status;

COMMIT;
