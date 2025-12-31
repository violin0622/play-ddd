BEGIN;

CREATE TABLE IF NOT EXISTS public.events (
    id character(26) NOT NULL,
    seq BIGSERIAL NOT NULL,
    created_ts bigint NOT NULL DEFAULT 0,
    updated_ts bigint NOT NULL DEFAULT 0,
    deleted_ts bigint NOT NULL DEFAULT 0,
    aggregate_id character(26) NOT NULL DEFAULT ''::character,
    aggregate_kind character varying(128) NOT NULL DEFAULT ''::character varying,
    version integer NOT NULL DEFAULT 0,
    kind character varying(128) NOT NULL DEFAULT ''::character varying,
    payload jsonb NOT NULL DEFAULT '{}'::jsonb
);

ALTER TABLE ONLY public.events
    ADD CONSTRAINT events_pkey PRIMARY KEY (id);
CREATE INDEX idx_events_aggregate_id ON public.events USING btree (aggregate_id);
CREATE UNIQUE INDEX uni_events_aggregate_id_version ON public.events USING btree (aggregate_id, version)
CREATE INDEX idx_events_aggregate_kind ON public.events USING btree (aggregate_kind);
CREATE INDEX idx_events_deleted_ts ON public.events USING btree (deleted_ts);
CREATE INDEX idx_events_kind ON public.events USING btree (kind);
COMMIT;
