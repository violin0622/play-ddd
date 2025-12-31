BEGIN    ;

ALTER     TABLE ONLY public.events
ADD       COLUMN seq BIGSERIAL NOT NULL,
ADD       COLUMN status CHARACTER VARYING(64) NOT NULL DEFAULT 'pending'::CHARACTER VARYING,
ADD       COLUMN reason CHARACTER VARYING(256) NOT NULL DEFAULT ''::CHARACTER VARYING;

CREATE    TABLE IF NOT EXISTS public.event_cursors (
          relay_cursor BIGINT NOT NULL DEFAULT 0,
          target_table CHARACTER VARYING(256) NOT NULL DEFAULT ''::CHARACTER VARYING
          );

INSERT INTO public.event_cursors (relay_cursor, target_table) VALUES (0, "events");

COMMIT   ;
