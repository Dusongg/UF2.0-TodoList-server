package models

import (
	"time"
)

//const (
//	Online  = true
//	Offline = false
//)

type TaskInfo struct {
	TaskID             string    `gorm:"column:task_id;type:varchar(25);primaryKey;not null;comment:任务单号"`
	Comment            string    `gorm:"column:comment;type:varchar(100);comment:任务描述"`
	EmergencyLevel     int       `gorm:"column:emergency_level;default:0;comment:紧急程度"`
	Deadline           time.Time `gorm:"column:deadline;type:date;default:(date_format((now() + interval 3 day),_utf8mb4'%Y-%m-%d'));comment:截止日期"`
	Principal          string    `gorm:"column:principal;type:varchar(20);not null;comment:负责人"`
	ReqNo              string    `gorm:"column:req_no;type:varchar(20);not null;comment:需求号;index:tasklist_req_no_index"`
	EstimatedWorkHours int8      `gorm:"column:estimated_work_hours;default:16;comment:预计工时"`
	State              string    `gorm:"column:state;type:varchar(20);default:带启动;comment:任务状态"`
	Type               int       `gorm:"column:type;default:0;comment:任务类型"`
}

func (TaskInfo) TableName() string {
	return "tasklist_table"
}

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
