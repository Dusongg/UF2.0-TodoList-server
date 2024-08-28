# Why use ？

- **Efficiency**：采用gRPC与客户端进行通信，Protobuf压缩传输数据，并依赖Golang的高并发
- **Function**：redis + gRPC流 实现消息广播
- **Readability & concise**：采用GORM管理数据库
- **Stability**：采用Nginx作反向代理，负载均衡，切换挂掉的服务节点
- **Containerization**：`docker compose up`一键部署运行



# 安装部署

## 1.1 Docker

1. `git clone git@github.com:Dusongg/UF2.0-TodoList-server.git`

   或`git clone https://github.com/Dusongg/UF2.0-TodoList-server.git`

2. 在docker-compose.yml文件目录下，运行`docker compose up`

3. 查看是否运行成功

   1. `grpcurl -plaintext -d '{"name": "dusong"}' localhost:8001 notification.Service/SayHello`

      ![image-20240815111932773](https://typora-dusong.oss-cn-chengdu.aliyuncs.com/image-20240815111932773.png)

   2. `docker ps`![image-20240815131511204](https://typora-dusong.oss-cn-chengdu.aliyuncs.com/image-20240815131511204.png)

4. **环境变量说明**

   在`./.env`文件中（`docker compose up`时会读取该文件到环境变量中）

   1. Redis![image-20240816151657078](https://typora-dusong.oss-cn-chengdu.aliyuncs.com/image-20240816151657078.png)
   2. MySQL![image-20240816151812767](https://typora-dusong.oss-cn-chengdu.aliyuncs.com/image-20240816151812767.png)
   3. SMTP![image-20240816152117264](https://typora-dusong.oss-cn-chengdu.aliyuncs.com/image-20240816152117264.png)
   4. Admin![image-20240816152206912](https://typora-dusong.oss-cn-chengdu.aliyuncs.com/image-20240816152206912.png)



## 1.2 手动部署

1. 安装redis和mysql

2. 配置好mysql用户和密码，进入mysql`mysql -uroot -p`， 创建数据库`create database OrderManager`（数据库名可以更改，同时修改config.json文件内容）

3. `git clone git@github.com:Dusongg/UF2.0-TodoList-server.git`或`git clone https://github.com/Dusongg/UF2.0-TodoList-server.git` 拉取远程仓库到本地

4. 在`./config/config.json`文件下修改相应配置，大部分与**1.1**中的环境变量配置相同

   ```json
   {
     "redis": {
       "host": "localhost",
       "port": "6379"
     },
   
     "mysql": {
       "gorm_dns": "root:123123@tcp(127.0.0.1:3306)/OrderManager?charset=utf8mb4&parseTime=True&loc=Local" 
     },
   
     "smtp": {
       "host": "smtp.gmail.com",
       "port": "587",
       "sender" : "dusong700@gmail.com",
       "password": "xxxx xxxx xxxx xxxx",  
       "send_time_point_1": "9",
       "send_time_point_2": "13"
   
     },
   
     "log_path": "./logs/app.log",
   
     "admins" : [
       "dusong",
       "游洋"
     ]
   
   }
   ```



5. 以Linux为例
   1. `cd UF2.0-TodoList-server `
   2. `go mod tidy`
   3. `go run .`




# 数据库字段说明

## 2.1 任务表(tasklist_table)

### 字段说明

- 主键：task_id

- 唯一键：req_no

![image-20240816161009673](https://typora-dusong.oss-cn-chengdu.aliyuncs.com/image-20240816161009673.png)

### DDL

```mysql
-- auto-generated definition
create table tasklist_table
(
    comment              varchar(100)                                                                    null comment '任务描述',
    task_id              varchar(25)                                                                     not null comment '任务单号'
        primary key,
    emergency_level      int         default 0                                                           null comment '紧急程度',
    deadline             date        default (date_format((now() + interval 3 day), _utf8mb4'%Y-%m-%d')) null comment '截止日期',
    principal            varchar(20)                                                                     not null comment '负责人',
    req_no               varchar(20)                                                                     not null comment '需求号',
    estimated_work_hours double      default 16                                                          null comment '预计工时',
    state                varchar(20) default '带启动'                                                    null comment '任务状态',
    type                 int         default 0                                                           null comment '任务类型'
);

create index tasklist_req_no_index
    on tasklist_table (req_no);


```

### GORM模型

```go
type TaskInfo struct {
	TaskID             string    `gorm:"column:task_id;type:varchar(25);primaryKey;not null;comment:任务单号"`
	Comment            string    `gorm:"column:comment;type:varchar(100);comment:任务描述"`
	EmergencyLevel     int       `gorm:"column:emergency_level;default:0;comment:紧急程度"`
	Deadline           time.Time `gorm:"column:deadline;type:date;default:(date_format((now() + interval 3 day),_utf8mb4'%Y-%m-%d'));comment:截止日期"`
	Principal          string    `gorm:"column:principal;type:varchar(20);not null;comment:负责人"`
	ReqNo              string    `gorm:"column:req_no;type:varchar(20);not null;comment:需求号;index:tasklist_req_no_index"`
	EstimatedWorkHours float64   `gorm:"column:estimated_work_hours;default:16;comment:预计工时"`
	State              string    `gorm:"column:state;type:varchar(20);default:带启动;comment:任务状态"`
	Type               int       `gorm:"column:type;default:0;comment:任务类型"`
}

func (TaskInfo) TableName() string {
	return "tasklist_table"
}
```



## 2.2 补丁表(patch_table)

### 字段说明

- 主键：patch_no

- 唯一键：req_no

![image-20240816161402947](https://typora-dusong.oss-cn-chengdu.aliyuncs.com/image-20240816161402947.png)

### DDL

```mysql
-- auto-generated definition
create table patch_table
(
    patch_no    varchar(20)  not null comment '补丁号'
        primary key,
    req_no      varchar(40)  not null comment '需求号',
    `describe`  text         null comment '问题描述',
    client_name varchar(20)  not null comment '客户名称',
    deadline    date         not null comment '预计发布时间',
    reason      varchar(100) null comment '补丁原因',
    sponsor     varchar(20)  not null comment '发起人',
    state       varchar(10)  null comment '发布状态'
);

create index patch_table_req_no_index
    on patch_table (req_no);

```

### GORM模型

```go
type PatchsInfo struct {
	PatchNo    string    `gorm:"column:patch_no;type:varchar(20);primaryKey;not null;comment:补丁号"`
	ReqNo      string    `gorm:"column:req_no;type:varchar(40);not null;comment:需求号;index:patch_table_req_no_index"`
	Describe   string    `gorm:"column:describe;type:text;comment:问题描述"`
	ClientName string    `gorm:"column:client_name;type:varchar(20);not null;comment:客户名称"`
	Deadline   time.Time `gorm:"column:deadline;type:date;not null;comment:预计发布时间"`
	Reason     string    `gorm:"column:reason;type:varchar(100);comment:补丁原因"`
	Sponsor    string    `gorm:"column:sponsor;type:varchar(20);not null;comment:发起人"`
	State      string    `gorm:"column:state;type:varchar(10);comment:发布状态"`
}

func (PatchsInfo) TableName() string {
	return "patch_table"
}
```



## 2.3 用户表(user_table)

### 字段说明

- 普通索引：name

![image-20240816161457805](https://typora-dusong.oss-cn-chengdu.aliyuncs.com/image-20240816161457805.png)

### DDL

```mysql
-- auto-generated definition
create table user_table
(
    name     varchar(20)      not null comment '姓名',
    job_no   bigint           not null comment '工号',
    password varchar(60)      not null comment '密码',
    email    varchar(50)      not null comment '邮箱',
    `group`  bigint default 0 null comment '分组编号',
    role_no  bigint default 0 null comment '角色编号'
);

create index idx_user_table_name
    on user_table (name);


```

### GORM模型

```go
type UserInfo struct {
	Name     string `gorm:"column:name;type:varchar(20);not null;index;comment:姓名"`
	JobNo    int32  `gorm:"column:job_no;not null;comment:工号"`
	Password string `gorm:"column:password;type:varchar(60);not null;comment:密码"`
	Email    string `gorm:"column:email;type:varchar(50);not null;comment:邮箱"`
	Group    int8   `gorm:"column:group;default:0;comment:分组编号"`
	RoleNo   int8   `gorm:"column:role_no;default:0;comment:角色编号"`
	//State    bool   `gorm:"column:state;default:false;comment:登录状态"`
}

func (UserInfo) TableName() string {
	return "user_table"
}
```



# 部分API说明

![image-20240816153407027](https://typora-dusong.oss-cn-chengdu.aliyuncs.com/image-20240816153407027.png)

![image-20240816155011703](https://typora-dusong.oss-cn-chengdu.aliyuncs.com/image-20240816155011703.png)

## 2.1 任务相关

## 2.1.1 AddTask & ImportXLSToTaskTable

- 导入或新增任务如果有主键或者唯一键冲突，**那么会插入失败，并将报错写入日志**

### 2.1.2 DelTask & ModTask

- 当当前用户不是管理员时，无法删除或修改其他用户任务
- 当修改或删除的是他人的任务时，会将消息发送给他人



## 2.2 补丁相关

### 2.2.1 ImportXLSToPatchTable

- 如果任务表中有与导入补丁相关联时，任务的截止日期会与补丁的截止日期两者取较小值，最终将补丁导入的消息通知给所有人
- 导入补丁与导入任务不同， 导入补丁如果遇见主键或唯一键冲突会**覆盖改行内容，而不会报错**



### 2.2.2 DelPatch

- 调用该函数会将补丁以及与补丁相关联的所有任务都删除

### 2.2.3 ModPatch

- 如果修改了补丁的时间，此时会将与补丁相关联的任务时间取两者较小值修改





## 2.3 发布订阅模式

### 2.3.1 Subscribe

- 订阅redis的updates频道，并接受消息通过stream发送给客户端





# Why gRPC ?

### 1. **协议层次和数据格式**：

- **HTTP**：HTTP是一种应用层协议，通常用于Web服务，通过请求和响应传输数据。常见的数据格式是JSON或XML，易于人类读取和调试。
- **gRPC**：gRPC基于HTTP/2协议，使用Protobuf（Protocol Buffers）作为数据序列化格式。Protobuf是一种二进制格式，比JSON更紧凑、更高效，尤其在传输大规模数据时优势明显。

### 2. **性能和效率**：

- **HTTP**：由于使用了文本格式（如JSON），HTTP请求和响应的数据量相对较大，解析速度相对较慢。
- **gRPC**：gRPC利用了HTTP/2的多路复用、流控制、头部压缩等特性，结合Protobuf的高效数据编码，在网络和处理性能上都有明显优势。

### 3. **流式通信**：

- **HTTP**：标准的HTTP是基于请求-响应的模型，不支持双向流通信。WebSocket可以实现实时双向通信，但不是HTTP本身的特性。
- **gRPC**：gRPC原生支持双向流式通信（双向流、服务器流、客户端流），适合需要实时数据交换的场景。

### 4. **类型安全**：

- **HTTP**：使用JSON时，由于没有严格的类型定义，容易在服务端和客户端之间产生类型不匹配的错误。
- **gRPC**：gRPC使用Protobuf定义服务和消息的结构，提供了强类型的接口，编译时即能发现数据结构不匹配的问题。

### 5. **生态系统和兼容性**：

- **HTTP**：作为Web应用的基础协议，HTTP有着广泛的生态系统支持和兼容性。几乎所有编程语言和平台都支持HTTP。
- **gRPC**：gRPC提供了跨语言支持，但在一些特定场景下（如浏览器端）不如HTTP普遍。此外，gRPC对客户端和服务端都需要依赖Protobuf定义文件，这可能增加复杂性。





# 进入容器内部，使用mysql客户端与redis-cli

```bash
docker exec -it ordermanager1 bash

mysql -uroot -p -h db   
 
redis-cli -h redis --raw
```



# 问题

## 1. mysql中文乱码

1. ![image-20240828155001025](https://typora-dusong.oss-cn-chengdu.aliyuncs.com/image-20240828155001025.png)
2. `SET NAMES 'utf8mb4';`
3. ![image-20240828155046648](https://typora-dusong.oss-cn-chengdu.aliyuncs.com/image-20240828155046648.png)

## 2. redis-cli中文显示十六进制编码

- 添加`--raw`选项

```
redis-cli -h redis --raw
```



