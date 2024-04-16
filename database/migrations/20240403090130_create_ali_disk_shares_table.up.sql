CREATE TABLE ali_disk_shares
(
    id         bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    created_at datetime(3)         NOT NULL,
    updated_at datetime(3)         NOT NULL,
    deleted_at datetime(3)         NULL,
    share_id   varchar(255)        NOT NULL default '',
    password   varchar(255)        NOT NULL default '',
    flag       varchar(255)        NOT NULL default '',
    PRIMARY KEY (id),
    KEY idx_ali_disk_shares_created_at (created_at),
    KEY idx_ali_disk_shares_updated_at (updated_at),
    KEY idx_ali_disk_shares_deleted_at (deleted_at)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4;
