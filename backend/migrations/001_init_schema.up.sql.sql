-- backend/migrations/001_init_schema.up.sql
CREATE TABLE users (
    id UUID PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE agents (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    hostname VARCHAR(255) NOT NULL,
    ip_address VARCHAR(45) NOT NULL,
    os VARCHAR(255) NOT NULL,
    architecture VARCHAR(255) NOT NULL,
    version VARCHAR(255) NOT NULL,
    last_heartbeat TIMESTAMP,
    status VARCHAR(50) NOT NULL,
    tags TEXT[],
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE subscriptions (
    id UUID PRIMARY KEY,
    user_id UUID REFERENCES users(id),
    plan_id VARCHAR(50) NOT NULL,
    status VARCHAR(50) NOT NULL,
    start_date TIMESTAMP NOT NULL,
    end_date TIMESTAMP NOT NULL,
    renewal_date TIMESTAMP NOT NULL,
    payment_method_id VARCHAR(255),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE usage_records (
    id UUID PRIMARY KEY,
    subscription_id UUID REFERENCES subscriptions(id),
    server_count INTEGER NOT NULL,
    cpu_hours FLOAT NOT NULL,
    memory_gb_hours FLOAT NOT NULL,
    network_gb FLOAT NOT NULL,
    storage_gb FLOAT NOT NULL,
    recorded_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE jobs (
    id UUID PRIMARY KEY,
    agent_id UUID REFERENCES agents(id),
    command TEXT NOT NULL,
    args TEXT[] NOT NULL,
    timeout INTERVAL NOT NULL,
    status VARCHAR(50) NOT NULL,
    result TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    started_at TIMESTAMP,
    completed_at TIMESTAMP
);

CREATE TABLE patches (
    id UUID PRIMARY KEY,
    agent_id UUID REFERENCES agents(id),
    updates TEXT[] NOT NULL,
    status VARCHAR(50) NOT NULL,
    logs TEXT,
    started_at TIMESTAMP NOT NULL,
    completed_at TIMESTAMP
);

CREATE TABLE vulnerabilities (
    id UUID PRIMARY KEY,
    agent_id UUID REFERENCES agents(id),
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    severity VARCHAR(50) NOT NULL,
    cve TEXT,
    fix TEXT,
    scanned_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE ssh_keys (
    id UUID PRIMARY KEY,
    agent_id UUID REFERENCES agents(id),
    public_key TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMP NOT NULL
);

CREATE INDEX idx_agents_status ON agents(status);
CREATE INDEX idx_subscriptions_user_id ON subscriptions(user_id);
CREATE INDEX idx_usage_records_subscription_id ON usage_records(subscription_id);
CREATE INDEX idx_jobs_agent_id ON jobs(agent_id);
CREATE INDEX idx_patches_agent_id ON patches(agent_id);
CREATE INDEX idx_vulnerabilities_agent_id ON vulnerabilities(agent_id);
CREATE INDEX idx_ssh_keys_agent_id ON ssh_keys(agent_id);