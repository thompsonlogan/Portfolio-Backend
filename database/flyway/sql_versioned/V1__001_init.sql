CREATE SCHEMA IF NOT EXISTS portfolio;

CREATE TABLE IF NOT EXISTS portfolio.visit(
  id BIGSERIAL PRIMARY KEY,
  source VARCHAR(255) NOT NULL,
  referrer VARCHAR(255) NOT NULL,
  user_agent VARCHAR(255),
  ip VARCHAR(255),
  visit_count INT,
  github_visit_count INT,
  linkedin_visit_count INT,
  resume_download_count INT,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX idx_unique_ip_source ON portfolio.visit (ip, source);