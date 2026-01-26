CREATE TABLE `user` (
    id bigint AUTO_INCREMENT,
    name varchar(255) NULL COMMENT '用户名',
    password varchar(255) NOT NULL DEFAULT '' COMMENT '密码',
    mobile varchar(255) NOT NULL DEFAULT '' COMMENT '手机号码',
    gender char(10) NOT NULL DEFAULT 'male' COMMENT '性别,male|female|unknown',
    nickname varchar(255) NULL DEFAULT '' COMMENT '昵称',
    type tinyint(1) NULL DEFAULT 0 COMMENT '用户类型, 0:普通用户,1:VIP用户,用于测试golang关键字',
    create_at timestamp NULL,
    update_at timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE mobile_index (mobile),
    UNIQUE name_index (name),
    PRIMARY KEY (id)
) ENGINE = InnoDB COLLATE utf8mb4_general_ci COMMENT '用户表';