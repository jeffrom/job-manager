CREATE TABLE queues (
    id char(253) PRIMARY KEY,
    v integer not null,
    concurrency smallint not null,
    retries smallint not null,
    unique_args boolean not null,
    duration bigint not null,
    checkin_duration bigint not null,
    claim_duration bigint not null,
    job_schema jsonb,
    created_at timestamp not null default now(),
    updated_at timestamp not null default now(),
    deleted_at timestamp,
    UNIQUE (id, v)
);

CREATE TYPE job_status AS ENUM (
    'queued',
    'running',
    'complete',
    'failed',
    'dead',
    'cancelled',
    'invalid'
);

CREATE TABLE jobs (
    id bigserial PRIMARY KEY,
    root_id bigint not null,
    v integer not null,
    queue char(253) not null REFERENCES queues (id),
    queue_v integer not null,
    attempt smallint,
    status job_status not null,
    args jsonb,
    data jsonb,

    enqueued_at timestamp not null default now(),
    started_at timestamp
    -- completed_at timestamp
);

COMMENT ON COLUMN jobs.root_id IS 'id of v1 job. 0 if row is v1.';

CREATE TABLE job_claims (
    id bigserial PRIMARY KEY,
    job_id bigint not null REFERENCES jobs (id),
    name char(63) not null,
    value char(63) not null,
    UNIQUE (job_id, name, value)
);

CREATE TABLE job_checkins (
    id bigserial PRIMARY KEY,
    job_id bigint not null REFERENCES jobs (id),
    data jsonb,
    created_at timestamp not null default now()
);

CREATE TABLE job_results (
    id bigserial PRIMARY KEY,
    job_id bigint not null REFERENCES jobs (id),
    status job_status not null,
    data jsonb,
    error text,
    started_at timestamp not null,
    completed_at timestamp not null
);