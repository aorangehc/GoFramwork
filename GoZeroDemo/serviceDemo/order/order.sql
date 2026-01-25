CREATE TABLE `order` (
    `order_sn` VARCHAR(64) NOT NULL COMMENT '订单号',
    `user_id` BIGINT NOT NULL DEFAULT 0 COMMENT '用户ID',
    `goods_id` BIGINT NOT NULL DEFAULT 0 COMMENT '商品ID',
    `num` BIGINT NOT NULL DEFAULT 0 COMMENT '数量',
    `amount` BIGINT NOT NULL DEFAULT 0 COMMENT '金额',
    `status` INT NOT NULL DEFAULT 0 COMMENT '状态',
    `create_at` DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `delete_at` DATETIME DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`order_sn`)
) ENGINE=InnoDB COLLATE utf8mb4_general_ci COMMENT='订单表';
