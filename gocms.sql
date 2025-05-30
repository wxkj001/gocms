/*
 Navicat Premium Dump SQL

 Source Server         : localhost
 Source Server Type    : MySQL
 Source Server Version : 80031 (8.0.31)
 Source Host           : localhost:3306
 Source Schema         : gocms

 Target Server Type    : MySQL
 Target Server Version : 80031 (8.0.31)
 File Encoding         : 65001

 Date: 30/05/2025 15:05:36
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for medias
-- ----------------------------
DROP TABLE IF EXISTS `medias`;
CREATE TABLE `medias`  (
  `id` int NOT NULL AUTO_INCREMENT,
  `media_url` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `media_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `content_type` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `size` bigint NULL DEFAULT NULL,
  `CreatedAt` datetime NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 3 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of medias
-- ----------------------------
INSERT INTO `medias` VALUES (1, '2025/04/22/ANIME-PICTURES.NET_-_744241-5760x3240-original-kirby+d+a-single-long+hair-looking+at+viewer-fringe.jpg', 'ANIME-PICTURES.NET_-_744241-5760x3240-original-kirby+d+a-single-long+hair-looking+at+viewer-fringe.jpg', 'image/jpeg', 2457786, '2025-04-22 16:32:38');
INSERT INTO `medias` VALUES (2, '2025/04/22/ANIME-PICTURES.NET_-_744241-5760x3240-original-kirby+d+a-single-long+hair-looking+at+viewer-fringe.jpg', 'ANIME-PICTURES.NET_-_744241-5760x3240-original-kirby+d+a-single-long+hair-looking+at+viewer-fringe.jpg', 'image/jpeg', 2457786, '2025-04-22 16:34:25');

-- ----------------------------
-- Table structure for permission
-- ----------------------------
DROP TABLE IF EXISTS `permission`;
CREATE TABLE `permission`  (
  `id` int NOT NULL AUTO_INCREMENT,
  `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `code` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `type` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `parentId` int NULL DEFAULT NULL,
  `path` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `redirect` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `icon` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `component` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `layout` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `keepAlive` tinyint NULL DEFAULT NULL,
  `method` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `description` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `show` tinyint NOT NULL DEFAULT 1 COMMENT '是否展示在页面菜单',
  `enable` tinyint NOT NULL DEFAULT 1,
  `order` int NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `IDX_30e166e8c6359970755c5727a2`(`code` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 28 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of permission
-- ----------------------------
INSERT INTO `permission` VALUES (1, '资源管理', 'Resource_Mgt', 'MENU', 2, '/pms/resource', NULL, 'i-fe:list', '/src/views/pms/resource/index.vue', NULL, NULL, NULL, NULL, 1, 1, 1);
INSERT INTO `permission` VALUES (2, '系统管理', 'SysMgt', 'MENU', NULL, NULL, NULL, 'i-fe:grid', NULL, NULL, NULL, NULL, NULL, 1, 1, 2);
INSERT INTO `permission` VALUES (3, '角色管理', 'RoleMgt', 'MENU', 2, '/pms/role', NULL, 'i-fe:user-check', '/src/views/pms/role/index.vue', NULL, NULL, NULL, NULL, 1, 1, 2);
INSERT INTO `permission` VALUES (4, '用户管理', 'UserMgt', 'MENU', 2, '/pms/user', NULL, 'i-fe:user', '/src/views/pms/user/index.vue', NULL, 1, NULL, NULL, 1, 1, 3);
INSERT INTO `permission` VALUES (5, '分配用户', 'RoleUser', 'MENU', 3, '/pms/role/user/:roleId', NULL, 'i-fe:user-plus', '/src/views/pms/role/role-user.vue', NULL, NULL, NULL, NULL, 0, 1, 1);
INSERT INTO `permission` VALUES (6, '业务示例', 'Demo', 'MENU', NULL, NULL, NULL, 'i-fe:grid', NULL, NULL, NULL, NULL, NULL, 1, 1, 1);
INSERT INTO `permission` VALUES (7, '图片上传', 'ImgUpload', 'MENU', 6, '/demo/upload', NULL, 'i-fe:image', '/src/views/demo/upload/index.vue', NULL, 1, NULL, NULL, 1, 1, 2);
INSERT INTO `permission` VALUES (8, '个人资料', 'UserProfile', 'MENU', NULL, '/profile', NULL, 'i-fe:user', '/src/views/profile/index.vue', NULL, NULL, NULL, NULL, 0, 1, 99);
INSERT INTO `permission` VALUES (9, '基础功能', 'Base', 'MENU', NULL, '/base', NULL, 'i-fe:grid', NULL, NULL, NULL, NULL, NULL, 1, 1, 0);
INSERT INTO `permission` VALUES (10, '基础组件', 'BaseComponents', 'MENU', 9, '/base/components', NULL, 'i-me:awesome', '/src/views/base/index.vue', NULL, NULL, NULL, NULL, 1, 1, 1);
INSERT INTO `permission` VALUES (11, 'Unocss', 'Unocss', 'MENU', 9, '/base/unocss', NULL, 'i-me:awesome', '/src/views/base/unocss.vue', NULL, NULL, NULL, NULL, 1, 1, 2);
INSERT INTO `permission` VALUES (12, 'KeepAlive', 'KeepAlive', 'MENU', 9, '/base/keep-alive', NULL, 'i-me:awesome', '/src/views/base/keep-alive.vue', NULL, 1, NULL, NULL, 1, 1, 3);
INSERT INTO `permission` VALUES (13, '创建新用户', 'AddUser', 'BUTTON', 4, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 1, 1, 1);
INSERT INTO `permission` VALUES (14, '图标 Icon', 'Icon', 'MENU', 9, '/base/icon', NULL, 'i-fe:feather', '/src/views/base/unocss-icon.vue', NULL, NULL, NULL, NULL, 1, 1, 0);
INSERT INTO `permission` VALUES (15, 'MeModal', 'TestModal', 'MENU', 9, '/testModal', NULL, 'i-me:dialog', '/src/views/base/test-modal.vue', NULL, NULL, NULL, NULL, 1, 1, 5);
INSERT INTO `permission` VALUES (22, '系统设置', 'SysConfig', 'MENU', 2, '/pms/config', '', 'i-fe:sliders', '/src/views/pms/config/index.vue', '', 0, '', '', 1, 1, 0);
INSERT INTO `permission` VALUES (23, '配置列表', 'SysConfigList', 'API', 22, '/api/admin/config/list', '', '', '', '', 0, '', '', 1, 1, 0);
INSERT INTO `permission` VALUES (24, '更新配置', 'UPSysConfig', 'API', 22, '/api/admin/config/*', '', '', '', '', 0, '', '', 1, 1, 0);
INSERT INTO `permission` VALUES (25, '查看用户信息', 'GetUserAPI', 'API', 8, '/api/admin/user/detail', '', '', '', '', 0, '', '', 1, 1, 0);
INSERT INTO `permission` VALUES (26, '获取权限列表', 'PermissionTree', 'API', 2, '/api/admin/permission/tree', '', '', '', '', 0, 'GET|POST', '', 1, 1, 0);

-- ----------------------------
-- Table structure for plugin
-- ----------------------------
DROP TABLE IF EXISTS `plugin`;
CREATE TABLE `plugin`  (
  `id` int NOT NULL AUTO_INCREMENT,
  `plugin_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `plugin_router` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `plugin_func` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of plugin
-- ----------------------------

-- ----------------------------
-- Table structure for profile
-- ----------------------------
DROP TABLE IF EXISTS `profile`;
CREATE TABLE `profile`  (
  `id` int NOT NULL AUTO_INCREMENT,
  `gender` int NULL DEFAULT NULL,
  `avatar` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'https://wpimg.wallstcn.com/f778738c-e4f8-4870-b634-56703b4acafe.gif?imageView2/1/w/80/h/80',
  `address` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `email` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `userId` int NOT NULL,
  `nickName` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `IDX_a24972ebd73b106250713dcddd`(`userId` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 5 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of profile
-- ----------------------------
INSERT INTO `profile` VALUES (2, NULL, 'https://wpimg.wallstcn.com/f778738c-e4f8-4870-b634-56703b4acafe.gif?imageView2/1/w/80/h/80', NULL, NULL, 1, 'Admin');
INSERT INTO `profile` VALUES (4, 0, '', '', '', 4, 't1');

-- ----------------------------
-- Table structure for role
-- ----------------------------
DROP TABLE IF EXISTS `role`;
CREATE TABLE `role`  (
  `id` int NOT NULL AUTO_INCREMENT,
  `code` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `enable` tinyint NOT NULL DEFAULT 1,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `IDX_ee999bb389d7ac0fd967172c41`(`code` ASC) USING BTREE,
  UNIQUE INDEX `IDX_ae4578dcaed5adff96595e6166`(`name` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 12 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of role
-- ----------------------------
INSERT INTO `role` VALUES (1, 'ADMIN', '超级管理员', 1);
INSERT INTO `role` VALUES (11, 'TEST', '测试', 1);

-- ----------------------------
-- Table structure for rule
-- ----------------------------
DROP TABLE IF EXISTS `rule`;
CREATE TABLE `rule`  (
  `ptype` varchar(100) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL DEFAULT '',
  `v0` varchar(100) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL DEFAULT '',
  `v1` varchar(100) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL DEFAULT '',
  `v2` varchar(100) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL DEFAULT '',
  `v3` varchar(100) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL DEFAULT '',
  `v4` varchar(100) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL DEFAULT '',
  `v5` varchar(100) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL DEFAULT '',
  INDEX `IDX_rule_v1`(`v1` ASC) USING BTREE,
  INDEX `IDX_rule_v2`(`v2` ASC) USING BTREE,
  INDEX `IDX_rule_v3`(`v3` ASC) USING BTREE,
  INDEX `IDX_rule_v4`(`v4` ASC) USING BTREE,
  INDEX `IDX_rule_v5`(`v5` ASC) USING BTREE,
  INDEX `IDX_rule_ptype`(`ptype` ASC) USING BTREE,
  INDEX `IDX_rule_v0`(`v0` ASC) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb3 COLLATE = utf8mb3_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of rule
-- ----------------------------
INSERT INTO `rule` VALUES ('p', '1', '/api/admin/*', 'GET|POST|PUT|DELETE|PATCH', '', '', '');
INSERT INTO `rule` VALUES ('p', '11', '2', 'role', '', '', '');
INSERT INTO `rule` VALUES ('p', '11', '22', 'role', '', '', '');
INSERT INTO `rule` VALUES ('p', '11', '23', 'role', '', '', '');
INSERT INTO `rule` VALUES ('p', '11', '24', 'role', '', '', '');
INSERT INTO `rule` VALUES ('g', '1', '1', 'user', '', '', '');
INSERT INTO `rule` VALUES ('g', '4', '11', 'user', '', '', '');

-- ----------------------------
-- Table structure for sys_config
-- ----------------------------
DROP TABLE IF EXISTS `sys_config`;
CREATE TABLE `sys_config`  (
  `id` int NOT NULL AUTO_INCREMENT,
  `config_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `config_key` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `config_value` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `config_type` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '设置类型',
  `config_group` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 11 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of sys_config
-- ----------------------------
INSERT INTO `sys_config` VALUES (1, 'AccessKeyID', 'AccessKeyID', 'py8xGC5ri-eumZng8QNd1zONI7q_PfD-IGCbMFtG', 'input', 's3');
INSERT INTO `sys_config` VALUES (2, '类型限制', 'allow_types', 'image/jpeg', 'input', 'upload');
INSERT INTO `sys_config` VALUES (3, '桶名称(BucketName)', 'BucketName', 'ws000', 'input', 's3');
INSERT INTO `sys_config` VALUES (4, 'Endpoint', 'Endpoint', 'ws000.s3.cn-north-1.qiniucs.com', 'input', 's3');
INSERT INTO `sys_config` VALUES (5, '是否启用OSS', 'is_oss', '1', 'switch', 'sys');
INSERT INTO `sys_config` VALUES (6, '地区(Location)', 'Location', 'z1', 'input', 's3');
INSERT INTO `sys_config` VALUES (7, '上传限制(MB)', 'max_size', '20', 'input', 'upload');
INSERT INTO `sys_config` VALUES (8, 'SecretAccessKey', 'SecretAccessKey', 'EwW_-X_4FO-6a0ejlfrrH0NmfwXuEEKz5vWaMopP', 'input', 's3');
INSERT INTO `sys_config` VALUES (9, '上传路径', 'upload_path', './public', 'input', 'upload');
INSERT INTO `sys_config` VALUES (10, '启用SSL', 'UseSSL', '1', 'switch', 's3');

-- ----------------------------
-- Table structure for udo_data
-- ----------------------------
DROP TABLE IF EXISTS `udo_data`;
CREATE TABLE `udo_data`  (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT 'Primary key',
  `tenant_id` bigint NOT NULL COMMENT 'Tenant ID for multi-tenant isolation',
  `object_id` bigint NOT NULL COMMENT 'Reference to udo_object.id',
  `data` json NOT NULL COMMENT 'JSON data containing field values',
  `status` tinyint(1) NOT NULL DEFAULT 1 COMMENT 'Status: 1-Active, 0-Inactive',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Creation time',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Update time',
  `created_by` bigint NULL DEFAULT NULL COMMENT 'Creator user ID',
  `updated_by` bigint NULL DEFAULT NULL COMMENT 'Updater user ID',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_tenant_object`(`tenant_id` ASC, `object_id` ASC) USING BTREE,
  INDEX `idx_object_id`(`object_id` ASC) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = 'UDO Data Storage' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of udo_data
-- ----------------------------

-- ----------------------------
-- Table structure for udo_field
-- ----------------------------
DROP TABLE IF EXISTS `udo_field`;
CREATE TABLE `udo_field`  (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT 'Primary key',
  `tenant_id` bigint NOT NULL COMMENT 'Tenant ID for multi-tenant isolation',
  `object_id` bigint NOT NULL COMMENT 'Reference to udo_object.id',
  `code` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT 'Field code, unique within object',
  `name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT 'Display name of field',
  `description` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL COMMENT 'Field description',
  `field_type` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT 'Field data type: string, number, boolean, date, datetime, enum, text, richtext, file, image',
  `is_required` tinyint(1) NOT NULL DEFAULT 0 COMMENT 'Whether field is required: 1-Required, 0-Optional',
  `min_length` int NULL DEFAULT NULL COMMENT 'Minimum length (for string/text types)',
  `max_length` int NULL DEFAULT NULL COMMENT 'Maximum length (for string/text types)',
  `regex_pattern` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT 'Regular expression for validation',
  `regex_message` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT 'Message to show when regex validation fails',
  `min_value` decimal(15, 6) NULL DEFAULT NULL COMMENT 'Minimum value (for number type)',
  `max_value` decimal(15, 6) NULL DEFAULT NULL COMMENT 'Maximum value (for number type)',
  `enum_options` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL COMMENT 'JSON array of options for enum type: [{\"value\": \"red\", \"label\": \"Red\"}, ...]',
  `default_value` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL COMMENT 'Default value in string format',
  `placeholder` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT 'Input placeholder text',
  `help_text` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT 'Help text for this field',
  `is_unique` tinyint(1) NOT NULL DEFAULT 0 COMMENT 'Whether field values must be unique: 1-Unique, 0-Not unique',
  `is_searchable` tinyint(1) NOT NULL DEFAULT 0 COMMENT 'Whether field should be searchable: 1-Searchable, 0-Not searchable',
  `is_system` tinyint(1) NOT NULL DEFAULT 0 COMMENT 'Whether field is system field: 1-System, 0-Custom',
  `sort_order` int NOT NULL DEFAULT 0 COMMENT 'Display order of the field',
  `status` tinyint(1) NOT NULL DEFAULT 1 COMMENT 'Status: 1-Active, 0-Inactive',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Creation time',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Update time',
  `created_by` bigint NULL DEFAULT NULL COMMENT 'Creator user ID',
  `updated_by` bigint NULL DEFAULT NULL COMMENT 'Updater user ID',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uk_object_code`(`object_id` ASC, `code` ASC) USING BTREE,
  INDEX `idx_tenant_id`(`tenant_id` ASC) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = 'UDO Field Definition' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of udo_field
-- ----------------------------

-- ----------------------------
-- Table structure for udo_log
-- ----------------------------
DROP TABLE IF EXISTS `udo_log`;
CREATE TABLE `udo_log`  (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT 'Primary key',
  `tenant_id` bigint NOT NULL COMMENT 'Tenant ID for multi-tenant isolation',
  `data_id` bigint NOT NULL COMMENT 'Reference to udo_data.id',
  `object_id` bigint NOT NULL COMMENT 'Reference to udo_object.id',
  `action` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT 'Action: create, update, delete',
  `data_before` json NULL COMMENT 'JSON data before change (for update/delete)',
  `data_after` json NULL COMMENT 'JSON data after change (for create/update)',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Log creation time',
  `created_by` bigint NULL DEFAULT NULL COMMENT 'User who performed the action',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_tenant_data`(`tenant_id` ASC, `data_id` ASC) USING BTREE,
  INDEX `idx_tenant_object`(`tenant_id` ASC, `object_id` ASC) USING BTREE,
  INDEX `idx_created_at`(`created_at` ASC) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = 'UDO Change Log' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of udo_log
-- ----------------------------

-- ----------------------------
-- Table structure for udo_object
-- ----------------------------
DROP TABLE IF EXISTS `udo_object`;
CREATE TABLE `udo_object`  (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT 'Primary key',
  `tenant_id` bigint NOT NULL COMMENT 'Tenant ID for multi-tenant isolation',
  `code` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT 'Unique object code within tenant',
  `name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT 'Display name of UDO object',
  `description` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL COMMENT 'Description of the object',
  `status` tinyint(1) NOT NULL DEFAULT 1 COMMENT 'Status: 1-Active, 0-Inactive',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Creation time',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Update time',
  `created_by` bigint NULL DEFAULT NULL COMMENT 'Creator user ID',
  `updated_by` bigint NULL DEFAULT NULL COMMENT 'Updater user ID',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uk_tenant_code`(`tenant_id` ASC, `code` ASC) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = 'UDO Object Definition' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of udo_object
-- ----------------------------

-- ----------------------------
-- Table structure for user
-- ----------------------------
DROP TABLE IF EXISTS `user`;
CREATE TABLE `user`  (
  `id` int NOT NULL AUTO_INCREMENT,
  `username` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `password` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `enable` tinyint(1) NOT NULL DEFAULT 1,
  `createTime` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
  `updateTime` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `IDX_78a916df40e02a9deb1c4b75ed`(`username` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 5 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of user
-- ----------------------------
INSERT INTO `user` VALUES (1, 'admin', 'e10adc3949ba59abbe56e057f20f883e', 1, '2025-04-17 11:29:07.847249', '2025-04-21 10:07:42.535885');
INSERT INTO `user` VALUES (4, 't1', 'e10adc3949ba59abbe56e057f20f883e', 1, '2025-05-27 14:25:03.000000', '2025-05-28 15:30:43.000000');

SET FOREIGN_KEY_CHECKS = 1;
