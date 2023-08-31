# GoYin
第六届字节跳动青训营后端极简版抖音

实现了用户功能，互动功能，社交功能，视频功能，聊天功能

## 项目特点

### 易维护

· 日志全覆盖，额外日志组件实现日志过滤，日志可视化，日志分类提取有效信息

· 可观测组件实现链路追踪，数据监控并实现可视化

· 采用微服务框架，并实现消息队列，实现结构解耦化

· 使用配置中心对配置文件进行统一管理，可进行版本回溯和监听查询

### 高可用

· k8s集群部署，实现高可拓展性

· 消息队列实现流量削峰，减小流量压力

· 实现了限流，熔断

### 高性能

· 采用kitex和hertz，高性能的微服务框架

· 适当使用并发处理减少请求时间

· 消息队列实现聊天消息处理

· 对象存储数据库存储视频和视频封面，节省了服务器资源

### 可靠性

· 敏感信息加密加盐处理

· 对敏感信息进行屏蔽过滤

· 具备CI功能，对整体项目代码进行了高覆盖率的单元测试和性能测试

· k8s集群搭建项目，具备容错性，支持横向拓展使项目管理更加便捷，负载均衡

## 功能介绍

· 基础接口：

​		· 视频流接口

​		· 用户注册接口

​		· 用户登录接口

​		· 用户信息

​		· 视频投稿

​		· 发布列表

· 互动接口：

​		· 赞操作

​		· 喜欢列表

​		· 评论操作

​		· 视频评论操作

· 社交接口：

​		· 关系操作

​		· 用户专注操作

​		· 用户粉丝操作

​		· 用户好友列表

​		· 消息：

​			· 聊天记录

​			· 用户好友列表

### 目录结构
```
├── deployment
│   ├── GoYin-k8s
│   │   ├── api
│   │   ├── chat
│   │   ├── ingress
│   │   │   └── common
│   │   │       └── rbac
│   │   ├── intetaction
│   │   ├── sociality
│   │   ├── user
│   │   └── video
│   ├── ip_info
│   └── opentelemetry
├── docs
│   └── static
├── server
│   ├── common
│   │   ├── consts
│   │   ├── middleware
│   │   ├── test
│   │   └── tools
│   ├── idl
│   ├── kitex_gen
│   │   ├── api
│   │   │   └── apiservice
│   │   ├── base
│   │   ├── chat
│   │   │   └── chatservice
│   │   ├── errno
│   │   ├── interaction
│   │   │   └── interactionserver
│   │   ├── sociality
│   │   │   └── socialityservice
│   │   ├── user
│   │   │   └── userservice
│   │   └── video
│   │       └── videoservice
│   └── service
│       ├── api
│       │   ├── biz
│       │   │   ├── handler
│       │   │   │   └── api
│       │   │   ├── model
│       │   │   │   ├── api
│       │   │   │   └── base
│       │   │   └── router
│       │   │       └── api
│       │   ├── config
│       │   ├── initialize
│       │   │   └── rpc
│       │   ├── models
│       │   ├── pkg
│       │   └── script
│       ├── chat
│       │   ├── config
│       │   ├── dao
│       │   ├── initialize
│       │   ├── model
│       │   ├── pkg
│       │   └── script
│       ├── interaction
│       │   ├── config
│       │   ├── dao
│       │   ├── initialize
│       │   ├── model
│       │   ├── pkg
│       │   └── script
│       ├── sociality
│       │   ├── config
│       │   ├── dao
│       │   ├── initialize
│       │   ├── model
│       │   ├── pkg
│       │   └── script
│       ├── user
│       │   ├── config
│       │   ├── dao
│       │   ├── initialize
│       │   ├── model
│       │   ├── pkg
│       │   └── script
│       └── video
│           ├── config
│           ├── dao
│           ├── initialize
│           ├── model
│           ├── pkg
│           └── script
```
## 代码架构图

​		![](./docs/static/架构图1-2023-08-30-1216.png)

![QQ20230830-230759](./docs/static/QQ20230830-230759.png)

## 技术栈

微服务框架：hertz，kitex

配置中心：nacos

服务注册中心：nacos

消息队列：nsq

数据库：mysql，redis，minio

可观测性组件：jaeger，victoria-metrics，grafana

日志组件：logstash，elasticsearch，kibana

云原生技术栈：k8s，docker，nginx-ingress

## 具体功能

### 数据库

#### mysql

![./docs/statistic/QQ20230830-220032.png](./docs/static/QQ20230830-220032.png)

![](./docs/static/QQ20230830-220107.png)



![QQ20230830-220132](./docs/static/QQ20230830-220132.png)

![QQ20230830-220141](./docs/static/QQ20230830-220141.png)

![QQ20230830-220157](./docs/static/QQ20230830-220157.png)

#### redis

![QQ20230830-223757](./docs/static/QQ20230830-223757.png)

![QQ20230830-223943](./docs/static/QQ20230830-223943.png)

#### minio

![QQ20230830-223911](./docs/static/QQ20230830-223911.png)

### nacos

![QQ20230830-223232](./docs/static/QQ20230830-223232.png)

### 可观测性

![QQ20230830-224259](./docs/static/QQ20230830-224259.png)

![QQ20230830-224308](./docs/static/QQ20230830-224308.png)

![QQ20230831-104652](./docs/static/QQ20230831-104652.png)

### 消息队列

![QQ20230830-230924](./docs/static/QQ20230830-230924.png)

![QQ20230830-230933](./docs/static/QQ20230830-230933.png)

![QQ20230831-103730](./docs/static/QQ20230831-103730.png)

### 日志

留给yjx

### k8s

#### 创建 master 节点

![QQ20230827-213037@2x](./docs/static/QQ20230827-213037@2x.png)

#### 创建 worker 节点并连接到 master 节点

![fa316762dd97c5956cc5e2895992959b](./docs/static/fa316762dd97c5956cc5e2895992959b.png)

#### 创建第一个 deployment 和 service

![51c1eb02abb88b4dd136c6a15f8463cf](./docs/static/51c1eb02abb88b4dd136c6a15f8463cf.png)

![QQ20230829-135029](./docs/static/QQ20230829-135029.png)

#### 创建剩下的几个服务节点

![QQ20230829-183659@2x](./docs/static/QQ20230829-183659@2x.png)

service 转发http请求成功

![QQ20230829-182301](./docs/static/QQ20230829-182301.png)



## 测试

留给wjj和spm

## 鸣谢
