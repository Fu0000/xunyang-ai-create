CREATE TABLE IF NOT EXISTS system_settings (
    `key` VARCHAR(50) NOT NULL PRIMARY KEY COMMENT 'Setting Key',
    `value` TEXT NOT NULL COMMENT 'Setting Value',
    `description` VARCHAR(255) DEFAULT '' COMMENT 'Description',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Updated at'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='System Settings';

-- Initialize default value for auto-approve (optional, will be handled by code fallback mostly)
-- Wait, the environment variable logic will fallback. 
-- We can pre-insert the default:
INSERT IGNORE INTO system_settings (`key`, `value`, `description`) 
VALUES ('inspiration_auto_approve', 'true', '灵感分享是否自动过审展览');
