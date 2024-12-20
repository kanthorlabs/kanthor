BEGIN;

CREATE TABLE IF NOT EXISTS kanthor_workspace (
  id VARCHAR(64) NOT NULL PRIMARY KEY,
  created_at BIGINT NOT NULL DEFAULT 0,
  updated_at BIGINT NOT NULL DEFAULT 0,
  owner_id VARCHAR(64) NOT NULL,
  name VARCHAR(256) NOT NULL,
  tier VARCHAR(256) NOT NULL
);

CREATE INDEX IF NOT EXISTS kanthor_ws_owner ON kanthor_workspace(owner_id ASC);

CREATE TABLE IF NOT EXISTS kanthor_application (
  id VARCHAR(64) NOT NULL PRIMARY KEY,
  created_at BIGINT NOT NULL DEFAULT 0,
  updated_at BIGINT NOT NULL DEFAULT 0,
  ws_id VARCHAR(64) NOT NULL,
  name VARCHAR(256) NOT NULL,
  FOREIGN KEY (ws_id) REFERENCES kanthor_workspace (id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS kanthor_app_ws_ref ON kanthor_application(ws_id ASC);

CREATE TABLE IF NOT EXISTS kanthor_endpoint (
  id VARCHAR(64) NOT NULL PRIMARY KEY,
  created_at BIGINT NOT NULL DEFAULT 0,
  updated_at BIGINT NOT NULL DEFAULT 0,
  app_id VARCHAR(64) NOT NULL,
  secret_key TEXT NOT NULL,
  name VARCHAR(256) NOT NULL,
  uri TEXT NOT NULL,
  method VARCHAR(64) NOT NULL,
  FOREIGN KEY (app_id) REFERENCES kanthor_application (id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS kanthor_ep_app_ref ON kanthor_endpoint(app_id ASC);

CREATE TABLE IF NOT EXISTS kanthor_route (
  id VARCHAR(64) NOT NULL PRIMARY KEY,
  created_at BIGINT NOT NULL DEFAULT 0,
  updated_at BIGINT NOT NULL DEFAULT 0,
  ep_id VARCHAR(64) NOT NULL,
  name VARCHAR(256) NOT NULL,
  condition_source VARCHAR(256) NOT NULL,
  condition_expression TEXT NOT NULL,
  exclusionary BOOLEAN NOT NULL DEFAULT FALSE,
  priority SMALLINT NOT NULL DEFAULT 0,
  FOREIGN KEY (ep_id) REFERENCES kanthor_endpoint (id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS kanthor_rt_ep_ref ON kanthor_route(ep_id ASC);

COMMIT;