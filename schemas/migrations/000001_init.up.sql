BEGIN;

-- SET client_encoding = 'UTF8';
CREATE TABLE IF NOT EXISTS public.novels (
    id character(26) NOT NULL,
    created_at bigint NOT NULL,
    updated_at bigint NOT NULL,
    deleted_at bigint NOT NULL,
    title character varying(256) NOT NULL,
    category character varying(256) NOT NULL,
    description character varying(1024) DEFAULT ''::character varying,
    tags jsonb DEFAULT '[]'::jsonb NOT NULL,
    author_id character(26) NOT NULL,
    toc jsonb DEFAULT '[]'::jsonb NOT NULL,
    status bigint DEFAULT 1 NOT NULL,
    word_count bigint
);

ALTER TABLE public.novels OWNER TO postgres;

COMMENT ON COLUMN public.novels.created_at IS 'autofilled create milliseconds unix timestamp.';
COMMENT ON COLUMN public.novels.updated_at IS 'autofilled update milliseconds unix timestamp.';
COMMENT ON COLUMN public.novels.deleted_at IS 'autofilled delete milliseconds unix timestamp. Used as soft deletion.';

ALTER TABLE ONLY public.novels
    ADD CONSTRAINT novels_pkey PRIMARY KEY (id);


CREATE INDEX idx_novels_deleted_at ON public.novels USING btree (deleted_at, deleted_at);

COMMIT;
