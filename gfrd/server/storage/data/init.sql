-- GFRD Framework - Test Database Script
-- 创建测试数据库和示例表

-- 创建数据库
CREATE DATABASE IF NOT EXISTS `gfrd` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE `gfrd`;

-- 1. 部门表
DROP TABLE IF EXISTS `sys_dept`;
CREATE TABLE `sys_dept` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '部门 ID',
  `parent_id` bigint DEFAULT '0' COMMENT '父部门 ID',
  `name` varchar(50) NOT NULL COMMENT '部门名称',
  `code` varchar(50) DEFAULT '' COMMENT '部门编码',
  `leader` varchar(20) DEFAULT '' COMMENT '负责人',
  `phone` varchar(11) DEFAULT '' COMMENT '联系电话',
  `email` varchar(100) DEFAULT '' COMMENT '邮箱',
  `status` tinyint DEFAULT '1' COMMENT '状态 (1 正常 0 禁用)',
  `sort` int DEFAULT '0' COMMENT '排序',
  `remark` varchar(255) DEFAULT '' COMMENT '备注',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_parent_id` (`parent_id`),
  KEY `idx_code` (`code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='系统部门表';

-- 插入测试数据
INSERT INTO `sys_dept` (`id`, `parent_id`, `name`, `code`, `leader`, `phone`, `email`, `status`, `sort`) VALUES
(1, 0, '总公司', '001', '张三', '13800138000', 'admin@example.com', 1, 1),
(2, 0, '上海分公司', '002', '李四', '13800138001', 'shanghai@example.com', 1, 2),
(3, 1, '研发部', '001001', '王五', '13800138002', 'dev@example.com', 1, 1),
(4, 1, '市场部', '001002', '赵六', '13800138003', 'market@example.com', 1, 2),
(5, 2, '技术部', '002001', '孙七', '13800138004', 'tech.sh@example.com', 1, 1);

-- 2. 用户表
DROP TABLE IF EXISTS `sys_user`;
CREATE TABLE `sys_user` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '用户 ID',
  `dept_id` bigint DEFAULT '0' COMMENT '部门 ID',
  `username` varchar(50) NOT NULL COMMENT '用户名',
  `password` varchar(100) NOT NULL COMMENT '密码',
  `nickname` varchar(50) DEFAULT '' COMMENT '昵称',
  `avatar` varchar(255) DEFAULT '' COMMENT '头像',
  `email` varchar(100) DEFAULT '' COMMENT '邮箱',
  `phone` varchar(20) DEFAULT '' COMMENT '手机号',
  `gender` tinyint DEFAULT '0' COMMENT '性别 (0 未知 1 男 2 女)',
  `status` tinyint DEFAULT '1' COMMENT '状态 (1 正常 0 禁用)',
  `last_login_at` datetime DEFAULT NULL COMMENT '最后登录时间',
  `last_login_ip` varchar(50) DEFAULT '' COMMENT '最后登录 IP',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_username` (`username`),
  KEY `idx_dept_id` (`dept_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='系统用户表';

-- 插入测试数据 (密码：admin123)
INSERT INTO `sys_user` (`id`, `dept_id`, `username`, `password`, `nickname`, `email`, `phone`, `gender`, `status`) VALUES
(1, 3, 'admin', '$2a$10$rO16LpwMlWjQJIjmQlV2p.4s3kVvNqXOvFkZkPqZ1K5sQ7X8X8X8X', '管理员', 'admin@example.com', '13800138000', 1, 1),
(2, 3, 'zhangsan', '$2a$10$rO16LpwMlWjQJIjmQlV2p.4s3kVvNqXOvFkZkPqZ1K5sQ7X8X8X8X', '张三', 'zhangsan@example.com', '13800138001', 1, 1),
(3, 4, 'lisi', '$2a$10$rO16LpwMlWjQJIjmQlV2p.4s3kVvNqXOvFkZkPqZ1K5sQ7X8X8X8X', '李四', 'lisi@example.com', '13800138002', 2, 1);

-- 3. 角色表
DROP TABLE IF EXISTS `sys_role`;
CREATE TABLE `sys_role` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '角色 ID',
  `name` varchar(50) NOT NULL COMMENT '角色名称',
  `code` varchar(50) NOT NULL COMMENT '角色编码',
  `description` varchar(255) DEFAULT '' COMMENT '角色描述',
  `data_scope` tinyint DEFAULT '1' COMMENT '数据范围 (1 全部 2 本部门 3 本人)',
  `status` tinyint DEFAULT '1' COMMENT '状态 (1 正常 0 禁用)',
  `sort` int DEFAULT '0' COMMENT '排序',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_code` (`code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='系统角色表';

-- 插入测试数据
INSERT INTO `sys_role` (`id`, `name`, `code`, `description`, `data_scope`, `status`, `sort`) VALUES
(1, '超级管理员', 'super_admin', '系统超级管理员', 1, 1, 1),
(2, '管理员', 'admin', '普通管理员', 2, 1, 2),
(3, '普通用户', 'user', '普通用户', 3, 1, 3);

-- 4. 菜单表
DROP TABLE IF EXISTS `sys_menu`;
CREATE TABLE `sys_menu` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '菜单 ID',
  `parent_id` bigint DEFAULT '0' COMMENT '父菜单 ID',
  `name` varchar(50) NOT NULL COMMENT '菜单名称',
  `path` varchar(200) DEFAULT '' COMMENT '路由地址',
  `component` varchar(255) DEFAULT '' COMMENT '组件路径',
  `icon` varchar(100) DEFAULT '' COMMENT '菜单图标',
  `type` tinyint DEFAULT '1' COMMENT '菜单类型 (1 目录 2 菜单 3 按钮)',
  `permission` varchar(50) DEFAULT '' COMMENT '权限标识',
  `sort` int DEFAULT '0' COMMENT '排序',
  `visible` tinyint DEFAULT '1' COMMENT '是否可见 (1 是 0 否)',
  `status` tinyint DEFAULT '1' COMMENT '状态 (1 正常 0 禁用)',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_parent_id` (`parent_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='系统菜单表';

-- 插入测试数据
INSERT INTO `sys_menu` (`id`, `parent_id`, `name`, `path`, `component`, `icon`, `type`, `permission`, `sort`, `visible`, `status`) VALUES
(1, 0, '系统管理', '/system', 'Layout', 'setting', 1, '', 1, 1, 1),
(2, 1, '部门管理', '/system/dept', 'system/dept/index', 'folder', 2, 'system:dept:list', 1, 1, 1),
(3, 1, '用户管理', '/system/user', 'system/user/index', 'user', 2, 'system:user:list', 2, 1, 1),
(4, 1, '角色管理', '/system/role', 'system/role/index', 'peoples', 2, 'system:role:list', 3, 1, 1);

-- 5. 字典表
DROP TABLE IF EXISTS `sys_dict_type`;
CREATE TABLE `sys_dict_type` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '字典类型 ID',
  `name` varchar(50) NOT NULL COMMENT '字典名称',
  `code` varchar(50) NOT NULL COMMENT '字典编码',
  `remark` varchar(255) DEFAULT '' COMMENT '备注',
  `status` tinyint DEFAULT '1' COMMENT '状态 (1 正常 0 禁用)',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_code` (`code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='字典类型表';

-- 插入测试数据
INSERT INTO `sys_dict_type` (`name`, `code`, `remark`, `status`) VALUES
('性别', 'gender', '', 1),
('状态', 'status', '', 1),
('是/否', 'yes_no', '', 1);

-- 字典数据表
DROP TABLE IF EXISTS `sys_dict_data`;
CREATE TABLE `sys_dict_data` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '字典数据 ID',
  `type_id` bigint NOT NULL COMMENT '字典类型 ID',
  `label` varchar(50) NOT NULL COMMENT '字典标签',
  `value` varchar(50) NOT NULL COMMENT '字典值',
  `sort` int DEFAULT '0' COMMENT '排序',
  `color` varchar(20) DEFAULT '' COMMENT '标签颜色',
  `status` tinyint DEFAULT '1' COMMENT '状态 (1 正常 0 禁用)',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_type_id` (`type_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='字典数据表';

-- 插入测试数据
INSERT INTO `sys_dict_data` (`type_id`, `label`, `value`, `sort`, `color`, `status`) VALUES
(1, '未知', '0', 1, '', 1),
(1, '男', '1', 2, 'blue', 1),
(1, '女', '2', 3, 'pink', 1),
(2, '正常', '1', 1, 'green', 1),
(2, '禁用', '0', 2, 'red', 1),
(3, '是', '1', 1, 'green', 1),
(3, '否', '0', 2, 'red', 1);
