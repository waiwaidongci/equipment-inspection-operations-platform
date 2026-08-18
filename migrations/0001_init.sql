CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    username TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    display_name TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS devices (
    id BIGSERIAL PRIMARY KEY,
    code TEXT NOT NULL UNIQUE,
    name TEXT NOT NULL,
    category TEXT NOT NULL,
    model TEXT NOT NULL,
    serial_number TEXT NOT NULL,
    manufacturer TEXT NOT NULL,
    location TEXT NOT NULL,
    install_date DATE NOT NULL,
    warranty_end DATE NOT NULL,
    owner TEXT NOT NULL,
    status TEXT NOT NULL,
    remark TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS inspection_plans (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    scope_type TEXT NOT NULL,
    device_id BIGINT REFERENCES devices(id) ON DELETE SET NULL,
    category TEXT NOT NULL DEFAULT '',
    period TEXT NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    owner TEXT NOT NULL,
    status TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS plan_check_items (
    id BIGSERIAL PRIMARY KEY,
    plan_id BIGINT NOT NULL REFERENCES inspection_plans(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    item_type TEXT NOT NULL,
    standard_range TEXT NOT NULL,
    required BOOLEAN NOT NULL DEFAULT FALSE,
    description TEXT NOT NULL,
    sort_order INTEGER NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS inspection_tasks (
    id BIGSERIAL PRIMARY KEY,
    plan_id BIGINT NOT NULL REFERENCES inspection_plans(id) ON DELETE CASCADE,
    device_id BIGINT NOT NULL REFERENCES devices(id) ON DELETE CASCADE,
    planned_at TIMESTAMPTZ NOT NULL,
    status TEXT NOT NULL,
    actual_executor_id BIGINT REFERENCES users(id) ON DELETE SET NULL,
    executed_at TIMESTAMPTZ,
    results JSONB NOT NULL DEFAULT '[]'::jsonb,
    remark TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    CONSTRAINT uq_tasks_plan_device_date UNIQUE (plan_id, device_id, planned_at)
);

CREATE TABLE IF NOT EXISTS inspection_records (
    id BIGSERIAL PRIMARY KEY,
    task_id BIGINT NOT NULL REFERENCES inspection_tasks(id) ON DELETE CASCADE,
    device_id BIGINT NOT NULL REFERENCES devices(id) ON DELETE CASCADE,
    plan_id BIGINT NOT NULL REFERENCES inspection_plans(id) ON DELETE CASCADE,
    executor_id BIGINT NOT NULL REFERENCES users(id) ON DELETE SET NULL,
    executed_at TIMESTAMPTZ NOT NULL,
    status TEXT NOT NULL,
    result_summary TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS anomalies (
    id BIGSERIAL PRIMARY KEY,
    device_id BIGINT NOT NULL REFERENCES devices(id) ON DELETE CASCADE,
    task_id BIGINT REFERENCES inspection_tasks(id) ON DELETE SET NULL,
    discoverer_id BIGINT NOT NULL REFERENCES users(id) ON DELETE SET NULL,
    discovered_at TIMESTAMPTZ NOT NULL,
    severity TEXT NOT NULL,
    description TEXT NOT NULL,
    status TEXT NOT NULL,
    assignee_id BIGINT REFERENCES users(id) ON DELETE SET NULL,
    progress TEXT NOT NULL DEFAULT '',
    cause_analysis TEXT NOT NULL DEFAULT '',
    close_note TEXT NOT NULL DEFAULT '',
    verification_result TEXT NOT NULL DEFAULT '',
    closed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS repairs (
    id BIGSERIAL PRIMARY KEY,
    device_id BIGINT NOT NULL REFERENCES devices(id) ON DELETE CASCADE,
    anomaly_id BIGINT REFERENCES anomalies(id) ON DELETE SET NULL,
    repair_type TEXT NOT NULL,
    content TEXT NOT NULL,
    vendor TEXT NOT NULL,
    started_at DATE NOT NULL,
    ended_at DATE NOT NULL,
    cost NUMERIC(12, 2) NOT NULL DEFAULT 0,
    result TEXT NOT NULL,
    remark TEXT NOT NULL,
    created_by BIGINT NOT NULL REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);
