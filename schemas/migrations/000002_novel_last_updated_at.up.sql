BEGIN;
ALTER TABLE public.novels ADD COLUMN
last_updated_at bigint NOT NULL DEFAULT 0;
END;
