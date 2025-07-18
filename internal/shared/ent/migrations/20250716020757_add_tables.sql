-- Create "products" table
CREATE TABLE `products` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT,
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `title` varchar(255) NOT NULL,
    `description` varchar(255) NULL,
    PRIMARY KEY (`id`)
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "seckill_activities" table
CREATE TABLE `seckill_activities` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT,
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `name` varchar(255) NOT NULL,
    `sku_id` bigint unsigned NOT NULL,
    `price` double NOT NULL,
    `seckill_price` double NOT NULL,
    `stock` bigint NOT NULL,
    `start_time` timestamp NOT NULL,
    `end_time` timestamp NOT NULL,
    PRIMARY KEY (`id`),
    INDEX `seckillactivity_end_time` (`end_time`),
    INDEX `seckillactivity_sku_id` (`sku_id`),
    INDEX `seckillactivity_start_time` (`start_time`)
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "seckill_orders" table
CREATE TABLE `seckill_orders` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT,
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `user_id` bigint unsigned NOT NULL,
    `activity_id` bigint unsigned NOT NULL,
    `sku_id` bigint unsigned NOT NULL,
    `price` double NOT NULL,
    `status` bigint NOT NULL DEFAULT 0,
    PRIMARY KEY (`id`),
    INDEX `seckillorder_activity_id` (`activity_id`),
    INDEX `seckillorder_sku_id` (`sku_id`),
    UNIQUE INDEX `seckillorder_user_id_activity_id_sku_id` (`user_id`, `activity_id`, `sku_id`)
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "skus" table
CREATE TABLE `skus` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT,
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `product_id` bigint unsigned NOT NULL,
    `specs` json NOT NULL,
    `price` double NOT NULL,
    `stock` bigint NOT NULL,
    PRIMARY KEY (`id`)
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "users" table
CREATE TABLE `users` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT,
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `phone` varchar(255) NOT NULL,
    `password_hash` varchar(255) NOT NULL,
    `status` enum('active','locked','deleted') NOT NULL DEFAULT "active",
    PRIMARY KEY (`id`),
    UNIQUE INDEX `phone` (`phone`)
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
