CREATE TABLE IF NOT EXISTS `payment_type` (
 `type_id`           int(10) NOT NULL AUTO_INCREMENT ,
 `type_name`         varchar(100) NOT NULL ,
 `payment_member_id` varchar(100) NULL ,
 `type_token`        varchar(50) NULL COMMENT '金鑰' ,
 `type_token_secret` varchar(100) NULL COMMENT '加密金鑰' ,
 `status`            varchar(1) NULL DEFAULT '1' COMMENT '狀態，0為關閉，1為啟用' ,
 `create_date`       timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ,
 `update_date`       timestamp NOT NULL ,
  PRIMARY KEY (`type_id`)
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS `payment_type_url` (
 `url_id`       int(10) NOT NULL AUTO_INCREMENT ,
 `type_id`      int(10) NOT NULL ,
 `url_name`     varchar(100) NULL COMMENT '網址功能名稱' ,
 `url_describe` varchar(10) NULL COMMENT '自定義編碼，可按照流程編號以利自動接續處理' ,
 `stage_type`   int(3) NULL COMMENT '環境類型，1為正式環境，0為測試環境' ,
 `request_url`  varchar(1000) NULL COMMENT '請求網址' ,
 `status`       varchar(1) NULL DEFAULT '1' COMMENT '狀態，0為關閉，1為啟用' ,
 `create_date`  timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ,
 `update_date`  timestamp NOT NULL ,
 `url_code`     varchar(10) NOT NULL ,
  PRIMARY KEY (`url_id`)
) ENGINE=InnoDB;

CREATE TABLE   IF NOT EXISTS `payment_platform_group_auth`(
 `auth_id`     int(10) NOT NULL AUTO_INCREMENT ,
 `group_id`    int(10) NOT NULL ,
 `type_id`     int(10) NOT NULL ,
 `create_date` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ,
 `update_date` timestamp NOT NULL ,
PRIMARY KEY (`auth_id`),
KEY `fkIdx_105` (`group_id`),
CONSTRAINT `FK_105` FOREIGN KEY `fkIdx_105` (`group_id`) REFERENCES `payment_platform_group` (`group_id`),
KEY `fkIdx_108` (`type_id`),
CONSTRAINT `FK_108` FOREIGN KEY `fkIdx_108` (`type_id`) REFERENCES `payment_type` (`type_id`)
) ENGINE=INNODB;


CREATE TABLE IF NOT EXISTS  `payment_platform` (
  `platform_id` int NOT NULL AUTO_INCREMENT,
  `platform_account` varchar(100) NOT NULL COMMENT '帳號',
  `platform_password` varchar(100) NOT NULL COMMENT '密碼',
  `platform_name` varchar(100) NOT NULL,
  `platform_group_id` int NOT NULL DEFAULT '1' COMMENT '1 為空群組',
  `platform_email` varchar(100) NOT NULL,
  `platform_token` varchar(50) NOT NULL COMMENT '金鑰',
  `platform_token_secret` varchar(100) NOT NULL COMMENT '加密金鑰',
  `status` varchar(1) DEFAULT '0' COMMENT '狀態，0為禁用，1為啟用',
  `create_date` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `update_date` timestamp NOT NULL,
  PRIMARY KEY (`platform_id`),
    UNIQUE KEY `platform_account_UNIQUE` (`platform_account`),
  KEY `fkIdx_101` (`platform_group_id`),
  CONSTRAINT `FK_101` FOREIGN KEY (`platform_group_id`) REFERENCES `payment_platform_group` (`group_id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS  `payment_platform_group`(
 `group_id`    int(10) NOT NULL AUTO_INCREMENT COMMENT '編號', 
 `group_name`  varchar(20) NOT NULL COMMENT '群組名稱',
 `group_describe`  varchar(45)  COMMENT '詳細',
 `create_date` timestamp NOT NULL  DEFAULT CURRENT_TIMESTAMP,
 `update_date` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,

PRIMARY KEY (`group_id`)
);

INSERT INTO `payment_platform_group` ( `group_name`, `group_describe`) VALUES ('nil', '預設群組');

CREATE TABLE IF NOT EXISTS `payment_platform_url` (
  `url_id` int(11) NOT NULL AUTO_INCREMENT,
  `payment_platform_id` int(11) NOT NULL COMMENT 'payment_type ID',
  `url_name` varchar(100) COMMENT '網址功能名稱',
  `url_code` varchar(10) COMMENT '自定義編碼，可按照流程編號以利自動接續處理',
  `stage_type` int(3) COMMENT '環境類型，1為正式環境，0為測試環境',
  `request_url` varchar(1000) COMMENT '請求網址',
  `status` varchar(1) DEFAULT '1' COMMENT '狀態，0為關閉，1為啟用',
  `create_date` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `update_date` timestamp NOT NULL ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`url_id`)
) ENGINE=InnoDB;


CREATE TABLE IF NOT EXISTS `order_payment` (
  `order_id` int(11) NOT NULL AUTO_INCREMENT,
  `payment_type_id` int(11) NOT NULL COMMENT 'payment_type ID',
  `platform_id` int(11) NOT NULL COMMENT 'platform ID',
  `order_client_id` varchar(50) NOT NULL COMMENT '客戶端訂單編號',
  `order_date` timestamp NOT NULL COMMENT '發起訂單時間',
  `order_original_data` varchar(1000) COMMENT '發起訂單原始資料',
  `order_price` decimal(10, 2) NOT NULL COMMENT '發起訂單原始資料',
  `redirect_url` varchar(1000) COMMENT '導向網址',
  `callback_original_data` varchar(1000),
  `received_callback_date` timestamp,
  `callback_url` varchar(1000) COMMENT '回調網址',
  `payment_id` varchar(50) NOT NULL COMMENT '支付端訂單編號',
  `stage_type` int(3) COMMENT '環境類型，1為正式環境，0為測試環境',
  `status` varchar(1) DEFAULT '1' COMMENT '狀態，0為禁用，1為啟用',
  `create_date` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `update_date` timestamp NOT NULL ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`order_id`)
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS `log_request` (
  `log_id` int(11) NOT NULL AUTO_INCREMENT,
  `platform_id` int(11) NOT NULL COMMENT 'payment_platform ID',
  `payment_platform_url_id` int(11) NOT NULL COMMENT 'payment platform url ID',
  `request_date` timestamp NOT NULL COMMENT '發起訂單時間',
  `request_original_data` varchar(1000) COMMENT '發起訂單原始資料',
  `status` varchar(1) DEFAULT '1' COMMENT '狀態，0為禁用，1為啟用',
  `create_date` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `update_date` timestamp NOT NULL ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`log_id`)
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS `log_connect` (
`log_id` int(11) NOT NULL AUTO_INCREMENT,
`statusCode` int(5) NOT NULL COMMENT '狀態碼',
`latencyTime` int(11) NOT NULL COMMENT '執行時間',
`clientIP` varchar(15) NOT NULL COMMENT '發起訂單時間',
`reqMethod` varchar(100) NOT NULL COMMENT '請求方式',
`reqURL` varchar(100) NOT NULL COMMENT '請求路由',
`create_date` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
`update_date` timestamp NOT NULL ON UPDATE CURRENT_TIMESTAMP,
PRIMARY KEY (`log_id`)
) ENGINE=InnoDB;
