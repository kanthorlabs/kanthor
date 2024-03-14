BEGIN;

CREATE TABLE IF NOT EXISTS kanthor_message (
  app_id VARCHAR(64) NOT NULL,
  id VARCHAR(64) NOT NULL,
  PRIMARY KEY (app_id, id),
  created_at BIGINT NOT NULL DEFAULT 0,
  tier VARCHAR(64) NOT NULL,
  tag VARCHAR(256) NOT NULL,
  body TEXT NOT NULL,
  metadata TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS kanthor_request (
  ep_id VARCHAR(64) NOT NULL,
  msg_id VARCHAR(64) NOT NULL,
  id VARCHAR(64) NOT NULL,
  PRIMARY KEY (ep_id, msg_id, id),
  app_id VARCHAR(64) NOT NULL,
  created_at BIGINT NOT NULL DEFAULT 0,
  tier VARCHAR(64) NOT NULL,
  tag VARCHAR(256) NOT NULL,
  metadata TEXT NOT NULL,
  method VARCHAR(64) NOT NULL,
  uri VARCHAR(2048) NOT NULL,
  headers TEXT NOT NULL,
  body TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS kanthor_response (
  ep_id VARCHAR(64) NOT NULL,
  msg_id VARCHAR(64) NOT NULL,
  id VARCHAR(64) NOT NULL,
  PRIMARY KEY (ep_id, msg_id, id),
  app_id VARCHAR(64) NOT NULL,
  created_at BIGINT NOT NULL DEFAULT 0,
  req_id VARCHAR(64) NOT NULL,
  tier VARCHAR(64) NOT NULL,
  tag VARCHAR(256) NOT NULL,
  metadata TEXT NOT NULL,
  status INT NOT NULL,
  uri VARCHAR(2048) NOT NULL,
  headers TEXT NOT NULL,
  body TEXT NOT NULL,
  error TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS kanthor_attempt (
  ep_id VARCHAR(64) NOT NULL,
  msg_id VARCHAR(64) NOT NULL,
  PRIMARY KEY (ep_id, msg_id),
  app_id VARCHAR(64) NOT NULL,
  tier VARCHAR(64) NOT NULL,
  schedule_counter INT NOT NULL,
  schedule_next BIGINT NOT NULL DEFAULT 0,
  scheduled_at BIGINT NOT NULL DEFAULT 0,
  completed_id VARCHAR(64),
  completed_at BIGINT NOT NULL DEFAULT 0
);

COMMIT;