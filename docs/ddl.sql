-- Create users table
CREATE TABLE IF NOT EXISTS `users` (
    `id` INT AUTO_INCREMENT PRIMARY KEY,
    `username` VARCHAR(20) NOT NULL UNIQUE COMMENT 'Username',
    `password` VARCHAR(255) NOT NULL COMMENT 'Hashed password',
    `display_name` VARCHAR(20) DEFAULT NULL COMMENT 'Display name',
    `role` INT DEFAULT 1 COMMENT 'Role: 1-common user, 2-admin',
    `status` INT DEFAULT 1 COMMENT 'Status: 1-enabled, 2-disabled',
    `email` VARCHAR(50) DEFAULT NULL COMMENT 'Email address',
    `github_id` VARCHAR(255) DEFAULT NULL COMMENT 'GitHub ID',
    `discord_id` VARCHAR(255) DEFAULT NULL COMMENT 'Discord ID',
    `oidc_id` VARCHAR(255) DEFAULT NULL COMMENT 'OIDC ID',
    `wechat_id` VARCHAR(255) DEFAULT NULL COMMENT 'WeChat ID',
    `telegram_id` VARCHAR(255) DEFAULT NULL COMMENT 'Telegram ID',
    `access_token` CHAR(32) DEFAULT NULL UNIQUE COMMENT 'Access token',
    `quota` INT DEFAULT 0 COMMENT 'Total quota',
    `used_quota` INT DEFAULT 0 COMMENT 'Used quota',
    `request_count` INT DEFAULT 0 COMMENT 'Request count',
    `group` VARCHAR(64) DEFAULT 'default' COMMENT 'User group',
    `aff_code` VARCHAR(32) DEFAULT NULL UNIQUE COMMENT 'Affiliate code',
    `aff_count` INT DEFAULT 0 COMMENT 'Affiliate count',
    `aff_quota` INT DEFAULT 0 COMMENT 'Remaining affiliate quota',
    `aff_history` INT DEFAULT 0 COMMENT 'History affiliate quota',
    `inviter_id` INT DEFAULT NULL COMMENT 'Inviter ID',
    `linux_do_id` VARCHAR(255) DEFAULT NULL COMMENT 'LinuxDO ID',
    `setting` TEXT COMMENT 'User settings',
    `remark` VARCHAR(255) DEFAULT NULL COMMENT 'Remark',
    `stripe_customer` VARCHAR(64) DEFAULT NULL COMMENT 'Stripe customer ID',
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT 'Creation time',
    `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Update time',
    `deleted_at` DATETIME DEFAULT NULL COMMENT 'Deletion time (soft delete)',
    INDEX `idx_username` (`username`),
    INDEX `idx_display_name` (`display_name`),
    INDEX `idx_email` (`email`),
    INDEX `idx_github_id` (`github_id`),
    INDEX `idx_discord_id` (`discord_id`),
    INDEX `idx_oidc_id` (`oidc_id`),
    INDEX `idx_wechat_id` (`wechat_id`),
    INDEX `idx_telegram_id` (`telegram_id`),
    INDEX `idx_access_token` (`access_token`),
    INDEX `idx_inviter_id` (`inviter_id`),
    INDEX `idx_linux_do_id` (`linux_do_id`),
    INDEX `idx_stripe_customer` (`stripe_customer`),
    INDEX `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='User table';

-- Insert test data
INSERT INTO `users` (`username`, `password`, `display_name`, `role`, `status`, `email`, `quota`, `used_quota`, `request_count`) VALUES
('admin', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'Administrator', 2, 1, 'admin@example.com', 1000, 0, 0),
('user1', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'User One', 1, 1, 'user1@example.com', 100, 10, 5),
('user2', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'User Two', 1, 1, 'user2@example.com', 100, 0, 0);

create table `uni-search-hub`.options
(
    `key` varchar(191) not null
        primary key,
    value longtext     null
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Options table';
