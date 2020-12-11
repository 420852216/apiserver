-- ----------------------------
--  Table structure for `user`
-- ----------------------------
CREATE TABLE user (
  id int(11) NOT NULL AUTO_INCREMENT,
  password varchar(128) NOT NULL COMMENT '密码',
  name varchar(128) NOT NULL COMMENT '昵称',
  email varchar(150) NOT NULL COMMENT '邮箱',
  phone varchar(11) DEFAULT '' COMMENT '手机号',
  created_at datetime NOT NULL COMMENT '创建时间',
  updated_at datetime NOT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `email` (`email`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
-- ----------------------------
--  Table structure for `dept`
-- ----------------------------
CREATE TABLE dept (
  id int(11) NOT NULL AUTO_INCREMENT,
  name varchar(128) NOT NULL COMMENT '部门名称',
  remark varchar(128) DEFAULT '' COMMENT '部门描述',
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
-- ----------------------------
--  Table structure for `Role`
-- ----------------------------
CREATE TABLE role (
  id int(11) NOT NULL AUTO_INCREMENT,
  name varchar(128) NOT NULL COMMENT '角色名称',
  remark varchar(128) DEFAULT '' COMMENT '角色描述',
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
-- ----------------------------
--  Table structure for `user_dept`
-- ----------------------------
CREATE TABLE user_dept (
  id int(11) NOT NULL AUTO_INCREMENT,
  user_id int(11) NOT NULL COMMENT '用户ID',
  dept_id int(11) NOT NULL COMMENT '部门ID',
  PRIMARY KEY (`id`),
  UNIQUE KEY `user_id_dept_id` (`user_id`,`dept_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
-- ----------------------------
--  Table structure for `dept_role`
-- ----------------------------
CREATE TABLE dept_role (
  id int(11) NOT NULL AUTO_INCREMENT,
  dept_id int(11) NOT NULL COMMENT '部门ID',
  role_id int(11) NOT NULL COMMENT '角色ID',
  PRIMARY KEY (`id`),
  UNIQUE KEY `dept_id_role_id` (`dept_id`,`role_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
-- ----------------------------
--  Table structure for `user_role`
-- ----------------------------
CREATE TABLE user_role (
  id int(11) NOT NULL AUTO_INCREMENT,
  user_id int(11) NOT NULL COMMENT '用户ID',
  role_id int(11) NOT NULL COMMENT '角色ID',
  PRIMARY KEY (`id`),
  UNIQUE KEY `user_id_role_id` (`user_id`,`role_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
-- ----------------------------
--  Table structure for `user_role`
-- ----------------------------
CREATE TABLE permission (
  id int(11) NOT NULL AUTO_INCREMENT,
  name varchar(11) NOT NULL COMMENT '权限名',
  remark varchar(128) DEFAULT '' COMMENT '权限描述',
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
-- ----------------------------
--  Table structure for `user_role`
-- ----------------------------
CREATE TABLE role_permission (
  id int(11) NOT NULL AUTO_INCREMENT,
  role_id varchar(128) NOT NULL COMMENT '角色ID',
  permission_id int(11) NOT NULL COMMENT '权限ID',
  PRIMARY KEY (`id`),
  UNIQUE KEY `role_id_permission_id` (`role_id`,`permission_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;



