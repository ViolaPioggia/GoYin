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

运用`gorm`框架实现数据库的增删改查：

![./docs/statistic/gorm.JPG](./docs/static/gorm.JPG)

- 通过事务处理确保数据库中的数据始终保持一致和可靠，避免了数据损坏或不完整的风险。 
- 使用索引帮助数据库快速定位到符合查询条件的数据，从而避免了全表扫描的开销，提高了查询效率。
- 利用gorm本身特性并通过简单逻辑判断一定程度上减少了sql注入风险。

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

### pprof

可通过以下代码查看：

``` go tool pprof -http=:8001 http://127.0.0.1:8080/debug/pprof/profile ```

![pprof.png](./docs/static/pprof.png)

### 日志

#### Kibana实现日志可视化

![kibana1](./docs/static/kibana1.png)

#### 查看指定错误等级日志

![kibana2](./docs/static/kibana2.png)

#### 根据feed流请求ip统计用户来源

![kibana3](./docs/static/kibana3.png)

#### kibana图表面板

- 日志错误监测
- feed流每日不同时间段的请求次数
- feed流请求用户的城市top10

![dashboard1](./docs/static/dashboard1.png)

![dashboard2](./docs/static/dashboard2.png)

#### logstash配置文件展示

- 根据不同类型的日志信息分别映射字段
- 实时收集日志信息，重启后从头读取

![logstash](./docs/static/logstash.png)

### k8s

#### 步骤

##### 创建 master 节点
虚拟机内网搭建 k8s 集群

#### 创建 master 节点

![QQ20230827-213037@2x](./docs/static/QQ20230827-213037@2x.png)

##### 创建 worker 节点并连接到 master 节点

![fa316762dd97c5956cc5e2895992959b](./docs/static/fa316762dd97c5956cc5e2895992959b.png)

##### 创建第一个 deployment 和 service

![51c1eb02abb88b4dd136c6a15f8463cf](./docs/static/51c1eb02abb88b4dd136c6a15f8463cf.png)

![QQ20230829-135029](./docs/static/QQ20230829-135029.png)

##### 创建剩下的几个服务节点

![QQ20230829-183659@2x](./docs/static/QQ20230829-183659@2x.png)

service 转发http请求成功

![QQ20230829-182301](./docs/static/QQ20230829-182301.png)

#### 配置 nginx-ingress

![QQ20230831-202319@2x](./docs/static/QQ20230831-202319@2x.png)

![QQ20230831-202347@2x](./docs/static/QQ20230831-202347@2x.png)

#### service+deployment部署优势
1. 服务发现和负载均衡：Service 可以将多个 Pod 组合成一个逻辑的服务，并为这个服务分配一个唯一的虚拟 IP 地址。这样，无论 Pod 的数量如何变化，服务的访问地址都保持不变，同时还可以实现负载均衡，将请求分发到不同的 Pod 上。
2. 自动伸缩：Deployment 可以根据指定的规则自动调整 Pod 的副本数量，以应对流量的增减。这样可以根据实际需求自动扩展或收缩应用程序的容量。
3. 无缝升级和回滚：Deployment 允许无缝地进行版本升级和回滚。通过逐步替换 Pod 的方式，可以确保服务在升级过程中不中断，并且在遇到问题时可以快速回滚到之前的版本。
4. 故障恢复和自愈：使用Deployment 部署的应用程序可以利用Kubernetes的自愈特性。如果某个 Pod 发生故障，Kubernetes 会自动重新创建一个新的 Pod 来替代。

## 测试

对程序进行了完备的性能测试，并且根据测试结果对代码进行了一定程度的优化

可通过运行以下代码查看：

``` go test -bench='.' -benchmem ```

![bench1.png](./docs/static/bench1.png)

优化后：

![bench2.png](./docs/static/bench2.png)

留给wjj

## 鸣谢
