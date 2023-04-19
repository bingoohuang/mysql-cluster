-- 创建数据库
create database `docker` default character set utf8mb4 collate utf8mb4_general_ci;

use docker;

-- 建表
DROP TABLE IF EXISTS `user`;

CREATE TABLE `user` (
    `id` bigint(20) NOT NULL,
    `created_at` bigint(40) DEFAULT NULL,
    `last_modified` bigint(40) DEFAULT NULL,
    `email` varchar(255) DEFAULT NULL,
    `first_name` varchar(255) DEFAULT NULL,
    `last_name` varchar(255) DEFAULT NULL,
    `username` varchar(255) DEFAULT NULL,
    PRIMARY KEY (`id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4;

-- 插入数据
INSERT INTO
    `user` (
        `id`,
        `created_at`,
        `last_modified`,
        `email`,
        `first_name`,
        `last_name`,
        `username`
    )
VALUES
    (
        0,
        1490257904,
        1490257904,
        'elliot@example.com',
        'Elliot',
        'Chen',
        'user'
    );