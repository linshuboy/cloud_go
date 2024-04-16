CREATE TABLE ali_disk_share_files
(
    id             bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    created_at     datetime(3)         NOT NULL,
    updated_at     datetime(3)         NOT NULL,
    deleted_at     datetime(3)         NULL,
    drive_id       varchar(255)        NOT NULL default '' comment '驱动器id',
    domain_id      varchar(255)        NOT NULL default '' comment '域id',
    file_id        varchar(255)        NOT NULL default '' comment '文件id',
    share_id       varchar(255)        NOT NULL default '' comment '共享id',
    name           varchar(255)        NOT NULL default '' comment '名称',
    type           varchar(255)        NOT NULL default '' comment '类型 file 文件 folder 文件夹',
    parent_file_id varchar(255)        NOT NULL default '' comment '父文件id',
    next_marker    varchar(255)        NOT NULL default '' comment '下一页标记',
    PRIMARY KEY (id),
    KEY idx_ali_disk_share_files_created_at (created_at),
    KEY idx_ali_disk_share_files_updated_at (updated_at),
    KEY idx_ali_disk_share_files_deleted_at (deleted_at)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4;
