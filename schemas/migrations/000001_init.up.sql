BEGIN;

-- SET client_encoding = 'UTF8';
CREATE TABLE IF NOT EXISTS public.novels (
    id character(26) NOT NULL,
    created_ts bigint NOT NULL DEFAULT 0,
    updated_ts bigint NOT NULL DEFAULT 0,
    deleted_ts bigint NOT NULL DEFAULT 0,
    author_id character(26) NOT NULL DEFAULT ''::character varying,
    title character varying(256) NOT NULL DEFAULT ''::character varying,
    category character varying(256) NOT NULL DEFAULT ''::character varying,
    description character varying(1024) NOT NULL DEFAULT ''::character varying,
    tags jsonb NOT NULL DEFAULT '[]'::jsonb,
    toc jsonb NOT NULL DEFAULT '[]'::jsonb,
    status bigint NOT NULL DEFAULT 1,
    word_count bigint NOT NULL DEFAULT 0
);

ALTER TABLE public.novels OWNER TO postgres;

COMMENT ON COLUMN public.novels.created_ts IS 'create milliseconds unix timestamp.';
COMMENT ON COLUMN public.novels.updated_ts IS 'update milliseconds unix timestamp.';
COMMENT ON COLUMN public.novels.deleted_ts IS 'delete milliseconds unix timestamp. Used as soft deletion.';

ALTER TABLE ONLY public.novels
    ADD CONSTRAINT novels_pkey PRIMARY KEY (id);


CREATE INDEX idx_novels_deleted_ts ON public.novels USING btree (deleted_ts, deleted_ts);

COMMIT;
