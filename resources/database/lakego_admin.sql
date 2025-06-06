# DROP TABLE IF EXISTS `pre__admin`;
# CREATE TABLE `pre__admin` (
#                               `id` char(36) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
#                               `name` varchar(30) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
#                               `password` char(32) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
#                               `password_salt` char(6) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
#                               `nickname` varchar(150) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
#                               `email` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
#                               `avatar` char(36) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
#                               `introduce` mediumtext COLLATE utf8mb4_unicode_ci,
#                               `is_root` tinyint(1) DEFAULT NULL,
#                               `status` tinyint(1) NOT NULL,
#                               `refresh_time` int(10) NOT NULL DEFAULT '0' COMMENT '刷新时间',
#                               `refresh_ip` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '刷新ip',
#                               `last_active` int(10) DEFAULT NULL,
#                               `last_ip` varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
#                               `update_time` int(10) DEFAULT NULL,
#                               `update_ip` varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
#                               `add_time` int(10) DEFAULT NULL,
#                               `add_ip` varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
#                               PRIMARY KEY (`id`)
# ) ENGINE=MyISAM DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
#
#
#
# -- 场地模板表
# DROP TABLE IF EXISTS `pre__venue_templates`;
# CREATE TABLE `pre__venue_templates` (
#                                    `id` INT AUTO_INCREMENT PRIMARY KEY,         -- 模板ID
#                                    `template_name` VARCHAR(100) COLLATE utf8mb4_unicode_ci NOT NULL, -- 模板名称
#                                    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- 创建时间
#                                    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP -- 更新时间
# ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
#
# -- 模板字段表
# DROP TABLE IF EXISTS `pre__template_fields`;
# CREATE TABLE `pre__template_fields` (
#                                    `id` INT AUTO_INCREMENT PRIMARY KEY,         -- 字段ID
#                                    `template_id` INT NOT NULL,                  -- 关联的模板ID
#                                    `field_name` VARCHAR(100) COLLATE utf8mb4_unicode_ci NOT NULL, -- 字段名称
#                                    `field_type` VARCHAR(50) COLLATE utf8mb4_unicode_ci NOT NULL,  -- 字段类型
#                                    `field_order` INT NOT NULL,                  -- 字段顺序
#                                    FOREIGN KEY (`template_id`) REFERENCES `pre__venue_templates`(`id`) -- 外键关联
# ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
#
# -- 场地记录表
# DROP TABLE IF EXISTS `pre__venue_records`;
# CREATE TABLE `pre__venue_records` (
#                                  `id` INT AUTO_INCREMENT PRIMARY KEY,         -- 场地ID
#                                  `template_id` INT NOT NULL,                  -- 关联的模板ID
#                                  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- 创建时间
#                                  `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP, -- 更新时间
#                                  FOREIGN KEY (`template_id`) REFERENCES `pre__venue_templates`(`id`) -- 外键关联
# ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
#
#
# -- 场地记录属性表
# DROP TABLE IF EXISTS `pre__venue_record_attributes`;
# CREATE TABLE `pre__venue_record_attributes` (
#                                            `id` INT AUTO_INCREMENT PRIMARY KEY,         -- 属性ID
#                                            `record_id` INT NOT NULL,                    -- 关联的场地记录ID
#                                            `field_id` INT NOT NULL,                     -- 添加 field_id 字段
#                                            `field_name` VARCHAR(100) COLLATE utf8mb4_unicode_ci NOT NULL, -- 字段名称
#                                            `field_value` TEXT COLLATE utf8mb4_unicode_ci, -- 字段值
#                                            FOREIGN KEY (`record_id`) REFERENCES `pre__venue_records`(`id`), -- 外键关联
#                                            FOREIGN KEY (`field_id`) REFERENCES `pre__template_fields`(`id`) -- 新增外键约束，关联到 pre__template_fields 表的 id
# ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
#
# -- 首先，插入一个标准模板
# INSERT INTO `pre__venue_templates` (template_name)
# VALUES ('标准模板');
#
# -- 然后，获取新插入模板的 ID
# SET @template_id = LAST_INSERT_ID();
#
# -- 插入模板字段到模板字段表
# INSERT INTO `pre__template_fields` (template_id, field_name, field_type, field_order)
# VALUES
#     (@template_id, '场地编码', 'VARCHAR', 1),
#     (@template_id, '场地名称', 'VARCHAR', 2),
#     (@template_id, '所在国家', 'VARCHAR', 3),
#     (@template_id, '场地地址', 'VARCHAR', 4);
#
#
#
# -- 观察者链接表
# DROP TABLE IF EXISTS `pre__link_info`;
# CREATE TABLE `pre__link_info` (
#                                       `id` INT AUTO_INCREMENT PRIMARY KEY,            -- 场地记录ID
#                                       `site_name` VARCHAR(50) COLLATE utf8mb4_unicode_ci NOT NULL, -- 场地名称
#                                       `sub_account` VARCHAR(50) COLLATE utf8mb4_unicode_ci NOT NULL, -- 子账号
#                                       `antpool_link` VARCHAR(100) COLLATE utf8mb4_unicode_ci NOT NULL, -- antpool链接
#                                       `f2pool_link` VARCHAR(100) COLLATE utf8mb4_unicode_ci NOT NULL,  -- f2pool链接
#                                       `sort_order` INT DEFAULT 0 NOT NULL                        -- 排序字段，默认值为0
# ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
#
#
# DROP TABLE IF EXISTS `pre__venue_report`;
# CREATE TABLE `pre__venue_report` (
#                                              `id` INT AUTO_INCREMENT PRIMARY KEY,                        -- 记录ID
#                                              `site_id` INT NOT NULL,                                     -- 关联的场地ID
#                                              `sub_account` VARCHAR(20) COLLATE utf8mb4_unicode_ci NOT NULL, -- 子账号
#                                              `record_date` DATE NOT NULL,                                -- 记录日期
#                                              `record_year` VARCHAR(10) NOT NULL,                                -- 记录年份
#                                              `record_month` TINYINT NOT NULL,                           -- 记录月份
#                                              `antpool_hash_rate` VARCHAR(20) NOT NULL,                -- antpool矿池算力 (单位: T)
#                                              `f2pool_hash_rate` VARCHAR(20) NOT NULL,                -- f2pool矿池算力 (单位: T)
#                                              `antpool_daily_income` VARCHAR(20) NOT NULL,           -- antpool日收益
#                                              `f2pool_daily_income` VARCHAR(20) NOT NULL,            -- f2pool日收益
#                                              `fb_income` VARCHAR(20) NOT NULL,                       -- FB收益
#                                              FOREIGN KEY (`site_id`) REFERENCES `pre__link_info`(`id`), -- 外键关联场地记录
#                                              UNIQUE (`site_id`, `sub_account`, `record_date`),          -- 唯一约束
#                                              INDEX (`record_year`),                                      -- 为年份创建索引
#                                              INDEX (`record_month`),                                     -- 为月份创建索引
#                                              INDEX (`site_id`, `record_year`, `record_month`)           -- 为场地和年月组合创建索引
# ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
#
# # -- 托管信息表
# DROP TABLE IF EXISTS `pre__custody_info`;
# CREATE TABLE `pre__custody_info` (
#                                      `id` INT AUTO_INCREMENT PRIMARY KEY,                                -- 托管信息记录ID
#                                      `venue_name` VARCHAR(50) COLLATE utf8mb4_unicode_ci NOT NULL,     -- 场地名
#                                      `sub_account_name` VARCHAR(50) COLLATE utf8mb4_unicode_ci NOT NULL, -- 子账户名
#                                      `observer_link` VARCHAR(100) COLLATE utf8mb4_unicode_ci,            -- 观察者链接（可为空）
#                                      `energy_ratio` VARCHAR(20) COLLATE utf8mb4_unicode_ci DEFAULT NULL,  -- 能耗比（百分比字符串，例如"15.75"）
#                                      `basic_hosting_fee` VARCHAR(20) COLLATE utf8mb4_unicode_ci DEFAULT NULL, -- 基础托管费（字符串格式的金额）
#                                      `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP -- 创建时间
# ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
#
# # -- 托管统计表
# DROP TABLE IF EXISTS `pre__custody_statistics`;
# CREATE TABLE `pre__custody_statistics` (
#                                            `id` INT AUTO_INCREMENT PRIMARY KEY,                                 -- 统计记录ID
#                                            `custody_id` INT NOT NULL,                                           -- 关联的托管信息ID
#                                            `energy_ratio` VARCHAR(20) COLLATE utf8mb4_unicode_ci DEFAULT NULL,  -- 能耗比（百分比字符串，例如"15.75"）
#                                            `basic_hosting_fee` VARCHAR(20) COLLATE utf8mb4_unicode_ci DEFAULT NULL, -- 基础托管费
#                                            `hourly_computing_power` VARCHAR(20) COLLATE utf8mb4_unicode_ci DEFAULT NULL, -- 24小时算力
#                                            `total_income_btc` VARCHAR(20) COLLATE utf8mb4_unicode_ci DEFAULT NULL, -- 24小时算力
#                                            `total_hosting_fee` VARCHAR(20) COLLATE utf8mb4_unicode_ci DEFAULT NULL, -- 总托管费
#                                            `total_income_usd` VARCHAR(20) COLLATE utf8mb4_unicode_ci DEFAULT NULL, -- 总收益
#                                            `net_income` VARCHAR(20) COLLATE utf8mb4_unicode_ci DEFAULT NULL, -- 净收益
#                                            `hosting_fee_ratio` VARCHAR(20) COLLATE utf8mb4_unicode_ci DEFAULT NULL, -- 托管费占比（字符串）
#                                            `report_date` VARCHAR(20) COLLATE utf8mb4_unicode_ci DEFAULT NULL, -- 统计日期
#                                            `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- 创建时间
#                                            FOREIGN KEY (`custody_id`) REFERENCES `pre__custody_info`(`id`),
#                                            INDEX `idx_report_date` (`report_date`), -- 现有的索引
#     -- 联合唯一索引，确保 `custody_id` 和 `report_date` 组合唯一
#                                            UNIQUE KEY `uniq_custody_report` (`custody_id`, `report_date`)
# ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
#
#
# -- 删除已存在的表（如果有）
# DROP TABLE IF EXISTS `product_historic_rates`;
#
# -- 创建新表存储历史行情数据
# CREATE TABLE `pre__btc_usd_candle` (
#                                           `timestamp` TIMESTAMP NOT NULL UNIQUE,                                        -- 时间戳（或标识时间点）
#                                           `price_low` DECIMAL(15,2) DEFAULT NULL,                                       -- 该时间段内最低价格
#                                           `price_high` DECIMAL(15,2) DEFAULT NULL,                                      -- 最高价格
#                                           `price_open` DECIMAL(15,2) DEFAULT NULL,                                      -- 开盘价格
#                                           `price_close` DECIMAL(15,2) DEFAULT NULL,                                     -- 收盘价格
#                                           `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,                              -- 创建时间
#     -- 索引优化
#                                           INDEX `idx_timestamp` (`timestamp`)
# ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
#
#
# CREATE TABLE pre__daily_average_price (
#                                      date VARCHAR(20) NOT NULL Unique,
#                                      cst_avg_price DECIMAL(15, 2) NOT NULL,
#                                      utc_avg_price DECIMAL(15, 2) NOT NULL,
#                                      created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
#                                          INDEX idx_date (date)
# );
#
#
# INSERT INTO `pre__admin` VALUES ('01cabd82-060d-405f-ba47-4d79fc47efcf','lakego','8966aff5289184448a004af81373c8f9','gazqzd','lakego','lakego@admin.com','5acfcd19-3a4c-4a28-8386-ae877952fd11','lakego-admin 是基于 gin、jwt 和 rbac 的 go 后台管理系统',0,1,0,'',1652759635,'127.0.0.1',1652587697,'127.0.0.1',1652545221,'127.0.0.1'),('642eb7b3-91ea-4808-bba6-f5f10938929a','admin','2a9b6b430ebe2f4257639e62ff9321bb','chNI7n','管理员','lakego-admin@admin.com','1f3cd4fb-f7e4-4b41-8663-167ca23ea5ab','lakego-admin 是基于 gin、jwt 和 rbac 的 go 后台管理系统',1,1,0,'',1675937003,'127.0.0.1',1652587697,'127.0.0.1',1652545221,'127.0.0.1');



# CREATE TABLE `pre__settlement_data` (
#                                    `settlement_point_name` VARCHAR(20) NOT NULL,                        -- 结算点名称 (最大20个字符)
#                                    `settlement_point_type` VARCHAR(20) NOT NULL,                        -- 结算点类型 (最大20个字符)
#                                    `delivery_date` DATE NOT NULL,                                         -- 交付日期
#                                    `delivery_hour` TINYINT NOT NULL,                                    -- 交付小时 (0-99)
#                                    `delivery_interval` TINYINT NOT NULL,                                -- 交付间隔 (0-99)
#                                    `settlement_point_price` FLOAT NOT NULL,                             -- 结算点价格
#                                    INDEX idx_settlement (`settlement_point_name`, `settlement_point_type`, `delivery_date`, `settlement_point_price`),  -- 联合索引
#                                    INDEX idx_price (`settlement_point_price`)                           -- 单列索引
# ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
#
#
#
# CREATE TABLE `pre__settlement_points` (
#                                      `settlement_point_id` INT AUTO_INCREMENT PRIMARY KEY,               -- 结算点唯一标识
#                                      `settlement_point_name` VARCHAR(20) NOT NULL,                      -- 结算点名称 (最大20个字符)
#                                      `settlement_point_type` VARCHAR(20) NOT NULL,                      -- 结算点类型 (最大20个字符)
#                                      UNIQUE KEY `unique_settlement_point` (`settlement_point_name`, `settlement_point_type`)  -- 复合唯一索引
# ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


# CREATE TABLE `pre__settlement_data_t` (
#                                         `settlement_point_name` VARCHAR(20) NOT NULL,                        -- 结算点名称 (最大20个字符)
#                                         `delivery_date` DATE NOT NULL,                                         -- 交付日期
#                                         `delivery_hour` TINYINT NOT NULL,                                    -- 交付小时 (0-99)
#                                         `settlement_point_price` FLOAT NOT NULL,                             -- 结算点价格
#                                         INDEX idx_settlement (`settlement_point_name`, `delivery_date`, `settlement_point_price`),  -- 联合索引
#                                         INDEX idx_price (`settlement_point_price`)                           -- 单列索引
# ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
#
#
#
# CREATE TABLE `pre__settlement_points_t` (
#                                           `settlement_point_id` INT AUTO_INCREMENT PRIMARY KEY,               -- 结算点唯一标识
#                                           `settlement_point_name` VARCHAR(20) NOT NULL,                      -- 结算点名称 (最大20个字符)
#                                           UNIQUE KEY `unique_settlement_point` (`settlement_point_name`)  -- 复合唯一索引
# ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


CREATE TABLE pre__mining_pools (
                              id INT AUTO_INCREMENT PRIMARY KEY,  -- 主键，自增ID
                              pool_type ENUM('NS', 'CANG') NOT NULL,  -- 矿池类型
                              pool_category ENUM('主矿池', '备用矿池') NOT NULL, -- 矿池类别
                              pool_name VARCHAR(255) NOT NULL,  -- 账户名
                              theoretical_hashrate DECIMAL(10, 2),  -- 理论算力
                              link VARCHAR(255),                     -- 链接
                              country VARCHAR(100),                     -- 所属地区
                              `sort_order` INT DEFAULT 0,                                  -- 排序字段，默认为0
                              `is_enabled` TINYINT(1) DEFAULT 1,                           -- 启用状态，默认为1（启用）
                              updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP, -- 更新时间
                              UNIQUE KEY unique_pool (pool_type, pool_category, pool_name) -- 唯一索引，确保组合独一无二
);

CREATE TABLE pre__mining_settlement_records (
                                         id INT AUTO_INCREMENT PRIMARY KEY,                -- 主键，自增ID
                                         pool_id INT NOT NULL,                              -- 外键，关联到 pre__mining_pools 表的 id
                                         settlement_date VARCHAR(30) NOT NULL,                     -- 结算日期
                                         settlement_hashrate DECIMAL(15, 2) NOT NULL,      -- 结算算力
                                         settlement_theoretical_hashrate DECIMAL(10, 2),  -- 结算时理论算力
                                         settlement_profit_btc DECIMAL(15, 8) NOT NULL,    -- 结算收益 BTC
                                         settlement_profit_fb DECIMAL(15, 8) NOT NULL,   -- 结算算力 FB
                                         created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,    -- 创建时间
                                         FOREIGN KEY (pool_id) REFERENCES pre__mining_pools(id), -- 关联外键，确保 pool_id 存在于 pre__mining_pools 表
                                         UNIQUE KEY unique_pool_date (pool_id, settlement_date)     -- 唯一索引，确保 pool_id 和 created_at 的组合独一无二
);

CREATE TABLE pre__mining_pool_status (
                                         id INT AUTO_INCREMENT PRIMARY KEY,                    -- 主键，自增ID
                                         pool_id INT NOT NULL,                                  -- 外键，关联到 pre__mining_pools 表的 id
                                         current_hashrate DECIMAL(30, 2) NOT NULL,            -- 实时算力
                                         online_machines INT NOT NULL,                          -- 实时在线机器数
                                         offline_machines INT NOT NULL,                         -- 离线机器数
                                         last_24h_hashrate DECIMAL(30, 2) NOT NULL,           -- 24小时算力
                                         last_update TIMESTAMP DEFAULT CURRENT_TIMESTAMP,      -- 最后更新时间
                                         FOREIGN KEY (pool_id) REFERENCES pre__mining_pools(id) -- 关联外键，确保 pool_id 存在于 pre__mining_pools 表
);