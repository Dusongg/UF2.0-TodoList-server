# 安装部署

## 1.1 Docker

1. `git pull git@github.com:Dusongg/UF2.0-TodoList-server.git`或`git pull https://github.com/Dusongg/UF2.0-TodoList-server.git`

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

3. `git pull git@github.com:Dusongg/UF2.0-TodoList-server.git`或`git pull https://github.com/Dusongg/UF2.0-TodoList-server.git`拉取远程仓库到本地

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



5. 运行`OrderManager.exe`文件



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
