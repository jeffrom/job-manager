CREATE TABLE queues (
    id bigserial PRIMARY KEY,
    name varchar(253) not null,
    v integer not null,
    retries smallint not null,
    unique_args boolean not null,
    duration bigint not null,
    checkin_duration bigint not null,
    claim_duration bigint not null,
    job_schema jsonb,
    backoff_initial_duration bigint not null default 0,
    backoff_max_duration bigint not null default 0,
    backoff_factor float not null default 0,
    created_at timestamp not null default now(),
    updated_at timestamp not null default now(),
    deleted_at timestamp,
    UNIQUE (name, v)
);

CREATE INDEX queues_name ON queues (name);

CREATE TABLE queue_labels (
    queue varchar(253) not null,
    name varchar(63) not null,
    value varchar(63) not null,
    PRIMARY KEY (queue, name, value)
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
    v integer not null,
    queue_id bigint not null REFERENCES queues (id),
    attempt smallint,
    status job_status not null,
    args jsonb,
    data jsonb,

    enqueued_at timestamp not null default now(),
    started_at timestamp,
    completed_at timestamp
);

CREATE INDEX jobs_status ON jobs (status);

CREATE TABLE job_claims (
    id bigserial PRIMARY KEY,
    job_id bigint not null REFERENCES jobs (id),
    name varchar(63) not null,
    value varchar(63) not null,
    UNIQUE (job_id, name, value)
);

CREATE INDEX job_claims_job_id ON job_claims (job_id);

CREATE TABLE job_checkins (
    id bigserial PRIMARY KEY,
    job_id bigint not null REFERENCES jobs (id),
    data jsonb,
    created_at timestamp not null default now()
);

CREATE INDEX job_checkins_job_id ON job_checkins (job_id);

CREATE TABLE job_results (
    id bigserial PRIMARY KEY,
    job_id bigint not null REFERENCES jobs (id),
    status job_status not null,
    data jsonb,
    error text,
    started_at timestamp not null,
    completed_at timestamp not null
);

CREATE INDEX job_results_job_id ON job_results (job_id);

CREATE TABLE job_uniqueness (
    job_id bigint not null REFERENCES jobs (id),
    key bytea PRIMARY KEY,
    created_at timestamp not null default now()
);

CREATE INDEX job_uniqueness_key ON job_uniqueness (key);
