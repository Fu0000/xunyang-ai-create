-- Add attempts counter to email_verifications for brute-force protection
ALTER TABLE email_verifications ADD COLUMN IF NOT EXISTS attempts INT NOT NULL DEFAULT 0 COMMENT 'Number of failed verification attempts';
