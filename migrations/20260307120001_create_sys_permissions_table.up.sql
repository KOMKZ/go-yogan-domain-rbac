CREATE TABLE IF NOT EXISTS `sys_permissions` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `permission_code` VARCHAR(100) NOT NULL,
    `permission_name` VARCHAR(100) NOT NULL,
    `permission_type` VARCHAR(10) NOT NULL COMMENT 'READ, WRITE',
    `resource_code` VARCHAR(50) NOT NULL,
    `group_code` VARCHAR(50) NOT NULL,
    `description` VARCHAR(200) NOT NULL DEFAULT '',
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_permission_code` (`permission_code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
