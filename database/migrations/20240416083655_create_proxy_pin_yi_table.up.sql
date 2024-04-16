CREATE TABLE proxy_pin_yi
(
    id         bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    created_at datetime(3)         NOT NULL,
    updated_at datetime(3)         NOT NULL,
    deleted_at datetime(3)         NULL,
    user_id    varchar(255)        NOT NULL default '' comment '用户id',
    appkey     varchar(255)        NOT NULL default '' comment '用户密钥',
    PRIMARY KEY (id),
    KEY idx_proxy_pin_yi_created_at (created_at),
    KEY idx_proxy_pin_yi_updated_at (updated_at),
    KEY idx_proxy_pin_yi_deleted_at (deleted_at)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4;
