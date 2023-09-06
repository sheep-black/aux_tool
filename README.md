# auxTool

## Introduction

auxTool是一款前后端分离的数据驱动建模辅助工具

**Tech Stack:**

Vue3+vite+JavaScript+View UI Plus+Gin+MySql

## Install

node: [20.4.6](https://nodejs.cn/download/)

go: [1.20](https://go.dev/dl/).*

mySQL: [5.7.43](https://dev.mysql.com/downloads/installer/)

```
git clone https://github.com/yafeiya/aux_tool.git
```

## Configurarion

```
# 1.前端
# 路径 ./auxTool-frontEnd-main/url_config.ys  
    "frontEndUrl":'192.168.0.103',  # 可选项，部署时用

    "backEndUrl": 'http://192.168.0.103:8080'  # 必填，向后端发送请求的地址
# 2. 后端
# 路径 ./congfig/application.yml
  # 后端服务
   server:        # 默认自动获取本机地址
  
  # SQL数据库
   datasource:    # 在3016局域网内默认配置，否则依据本机mysql配置

```

## Build & Run

```
# 开启后端服务器
go run main.go

# 进入前端目录
cd ./auxTool-frontEnd-main

# 安装前端依赖
npm install

# 开启json虚拟后端
npm run json  # TODO后续删除

# 开启前端服务 (./前端代码目录)
npm run dev

```

## SQL

```
# 初次部署时，mySQL配置
# 开启mySQL
mysql -u root -p

# 创建数据库
CREATE DATABASE auxtool
    DEFAULT CHARACTER SET = 'utf8mb4';
## 1. sql文件创建
# 导入数据库
use auxtool;
source /xx/xxxx.sql;

# 导出数据库
mysqldump -u root -p auxtool > /xx/xxxx.sql;

## 2. sql语句创建
# 创建用户表
use auxtool;
create table user
(
Id INT(11) not null auto_increment PRIMARY KEY comment '主键' ,
UserName VARCHAR(15),
PassWord VARCHAR(15)
);

# 注册用户

```

## Note

```
# 允许其他ip访问mysql
use mysql;
select Host, User,Password from user;
update user set Host='%' where User='root';
flush privileges;
```

## Reference

[gorm](https://gorm.io/docs/)

[View UI Plus](https://www.iviewui.com/view-ui-plus/guide/introduce)
