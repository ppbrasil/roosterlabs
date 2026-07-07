CREATE TABLE IF NOT EXISTS leads (
    id BIGSERIAL PRIMARY KEY,
    token TEXT NOT NULL UNIQUE,
    language TEXT NOT NULL,
    profile TEXT NOT NULL,
    goal TEXT NOT NULL,
    maturity TEXT NOT NULL,
    challenge TEXT NOT NULL,
    email TEXT NOT NULL,
    linkedin_url TEXT NOT NULL,
    utm_source TEXT,
    utm_medium TEXT,
    utm_campaign TEXT,
    utm_term TEXT,
    utm_content TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS funnel_events (
    id BIGSERIAL PRIMARY KEY,
    token TEXT NOT NULL,
    event_type TEXT NOT NULL,
    step SMALLINT,
    language TEXT NOT NULL,
    page_path TEXT NOT NULL,
    payload JSONB,
    utm_source TEXT,
    utm_medium TEXT,
    utm_campaign TEXT,
    utm_term TEXT,
    utm_content TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_funnel_events_token_created_at
    ON funnel_events (token, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_funnel_events_event_type_step
    ON funnel_events (event_type, step);

CREATE INDEX IF NOT EXISTS idx_funnel_events_language
    ON funnel_events (language);
