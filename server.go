package main

import (
	"OrderManager/common"
	"OrderManager/pb"
	"context"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"strings"
	"time"
)

type server struct {
	pb.UnimplementedServiceServer
}

var Server = &server{}

func (s *server) SayHello(ctx context.Context, in *pb.TestRequest) (*pb.TestReply, error) {
	return &pb.TestReply{Answer: "Hello " + in.Name}, nil
}

func (s *server) Login(ctx context.Context, in *pb.LoginRequest) (*pb.LoginReply, error) {
	var user UserInfo
	res := db.Where("name = ?", in.Name).First(&user)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("this user does not exist")
		}
		return nil, res.Error
	}
	//if user.State == models.Online {
	//	return nil, errors.New("this user already online")
	//}
	if NotificationServer.checkIfLoggedIn(user.Name) {
		return nil, errors.New("this user already online")
	}
	//算法标识符 + 成本因子 + 盐值
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(in.Password))
	if err != nil {
		return nil, errors.New("wrong password")
	} else {
		return &pb.LoginReply{}, nil
	}
}

func (s *server) Register(ctx context.Context, in *pb.RegisterRequest) (*pb.RegisterReply, error) {
	var existingUser UserInfo
	res := db.Where("name = ?", in.User.Name).First(&existingUser)
	if res.RowsAffected > 0 {
		return &pb.RegisterReply{}, errors.New("user already exists")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(in.User.Password), bcrypt.DefaultCost)
	if err != nil {
		return &pb.RegisterReply{}, err
	}
	userinfo := UserInfo{
		Name:     in.User.Name,
		JobNo:    int32(in.User.JobNum),
		Password: string(hashedPassword),
		Email:    in.User.Email,
	}
	if err := db.Create(&userinfo).Error; err != nil {
		return nil, err
	}
	return &pb.RegisterReply{}, nil
}
func (s *server) GetUserInfo(ctx context.Context, in *pb.GetUserInfoRequest) (*pb.GetUserInfoReply, error) {
	var info UserInfo
	if err := db.Where("name = ?", in.UserName).First(&info).Error; err != nil {
		return nil, err
	}
	return &pb.GetUserInfoReply{
		JobNO:  info.JobNo,
		Email:  info.Email,
		Group:  int32(info.Group),
		RoleNo: int32(info.RoleNo),
	}, nil
}

func (s *server) ModUserInfo(ctx context.Context, in *pb.ModUserInfoRequest) (*pb.ModUserInfoReply, error) {
	if in.ModPass {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(in.Pass), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		if err := db.Where("name = ?", in.Name).Updates(UserInfo{
			Password: string(hashedPassword),
			Email:    in.Email,
			Group:    int8(in.Group),
			RoleNo:   int8(in.RoleNo),
		}).Error; err != nil {
			return nil, err
		}
	} else {
		if err := db.Where("name = ?", in.Name).Updates(UserInfo{
			Email:  in.Email,
			Group:  int8(in.Group),
			RoleNo: int8(in.RoleNo),
		}).Error; err != nil {
			return nil, err
		}
	}
	return &pb.ModUserInfoReply{}, nil
}

func (s *server) GetTaskListAll(ctx context.Context, in *pb.GetTaskListAllRequest) (*pb.GetTaskListAllReply, error) {
	var tasks []TaskInfo
	res := db.Find(&tasks)
	if res.Error != nil {
		return nil, res.Error
	}
	reply := &pb.GetTaskListAllReply{}
	reply.Tasks = common.AllTaskInfoToPbTask(tasks)
	return reply, nil
}
func (s *server) GetTaskListOne(ctx context.Context, in *pb.GetTaskListOneRequest) (*pb.GetTaskListOneReply, error) {
	var tasks []TaskInfo
	if err := db.Where("principal = ?", in.Name).Find(&tasks).Error; err != nil {
		return nil, err
	}
	reply := &pb.GetTaskListOneReply{}
	reply.Tasks = common.AllTaskInfoToPbTask(tasks)
	return reply, nil
}

// TODO: 考虑导入任务时，是否查询联系到的补丁的时间，并修改
func (s *server) ImportXLSToTaskTable(ctx context.Context, in *pb.ImportToTaskListRequest) (*pb.ImportToTaskListReply, error) {
	taskInfos := common.AllPbTaskToTaskInfo(in.Tasks)
	tx := db.Begin()
	for _, taskInfo := range taskInfos {
		res := db.Model(&PatchsInfo{}).Where("req_no = ?", taskInfo.ReqNo).Select("deadline")
		if res.Error == nil {
			res.Scan(&taskInfo.Deadline)
		}
		if err := tx.Create(&taskInfo).Error; err != nil {
			logrus.Warning(err)
			continue
		}

		msg := fmt.Sprintf("<%s> -> import tasks counts: %d -> <%s>", in.User, len(in.Tasks), taskInfo.Principal)
		NotificationServer.updateDatabaseAndNotify(msg)
	}
	tx.Commit()
	return &pb.ImportToTaskListReply{}, nil
}
func (s *server) DelTask(ctx context.Context, in *pb.DelTaskRequest) (*pb.DelTaskReply, error) {
	if !isAdmin(in.User) && in.User != in.Principal {
		return nil, errors.New("you don't have permission")
	}
	task := TaskInfo{}
	if err := db.Where("task_id = ?", in.TaskNo).Delete(&task).Error; err != nil {
		return nil, err
	} else {
		msg := fmt.Sprintf("<%s> -> delete task: %s -> <%s>", in.User, in.TaskNo, in.Principal)
		NotificationServer.updateDatabaseAndNotify(msg)
		return &pb.DelTaskReply{}, nil
	}
}
func (s *server) ModTask(ctx context.Context, in *pb.ModTaskRequest) (*pb.ModTaskReply, error) {
	if !isAdmin(in.User) && in.User != in.T.Principal {
		return nil, errors.New("you don't have permission")
	}
	if err := db.Save(common.OnePbTaskToTaskInfo(in.T)).Error; err != nil {
		return nil, err
	} else {
		msg := fmt.Sprintf("<%s> -> modifiy task: %s -> <%s>", in.User, in.T.TaskId, in.T.Principal)

		NotificationServer.updateDatabaseAndNotify(msg)
		return &pb.ModTaskReply{}, nil

	}
}
func (s *server) AddTask(ctx context.Context, in *pb.AddTaskRequest) (*pb.AddTaskReply, error) {
	var task TaskInfo
	if err := db.Where("task_id = ?", in.T.TaskId).First(&task).Error; err == nil { //有重复task_id
		return nil, errors.New("task already exists")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) { //查找出错，且不是没有找到的错误
		return nil, err
	}
	if err := db.Create(common.OnePbTaskToTaskInfo(in.T)).Error; err != nil {
		return nil, err
	} else {
		msg := fmt.Sprintf("<%s> -> add task: %s -> <%s>", in.User, in.T.TaskId, in.T.Principal)
		NotificationServer.updateDatabaseAndNotify(msg)
		return &pb.AddTaskReply{}, nil
	}
}
func (s *server) QueryTaskWithSQL(ctx context.Context, in *pb.QueryTaskWithSQLRequest) (*pb.QueryTaskWithSQLReply, error) {
	var tasks []TaskInfo
	db.Raw(in.Sql).Scan(tasks)
	reply := &pb.QueryTaskWithSQLReply{}
	reply.Tasks = common.AllTaskInfoToPbTask(tasks)
	return reply, nil
}
func (s *server) QueryTaskWithField(ctx context.Context, in *pb.QueryTaskWithFieldRequest) (*pb.QueryTaskWithFieldReply, error) {
	var tasks []TaskInfo
	whereCond := fmt.Sprintf(" %s = ?", in.Field)
	if err := db.Where(whereCond, in.FieldValue).Find(&tasks).Error; err != nil {
		return nil, err
	}
	return &pb.QueryTaskWithFieldReply{Tasks: common.AllTaskInfoToPbTask(tasks)}, nil
}

func (s *server) ImportXLSToPatchTable(ctx context.Context, in *pb.ImportXLSToPatchRequest) (*pb.ImportXLSToPatchReply, error) {
	tx := db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	for _, patch := range in.Patchs {
		go func() {
			newDeadline := patch.Deadline
			err := s.ModDeadLineInPatchsAndTasks(context.Background(), &modDeadlineInfo{patchNo: patch.PatchNo, newDeadline: newDeadline, user: in.User}, false)
			if err != nil {
				logrus.Warning(err)
			}
		}()
	}

	if err := tx.Save(common.AllPbPatchsToPatchsInfo(in.Patchs)).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	msg := fmt.Sprintf("<%s> -> import patchs counts: %d -> <ALL>", in.User, len(in.Patchs))
	NotificationServer.updateDatabaseAndNotify(msg)
	return &pb.ImportXLSToPatchReply{}, nil
}

// 修改补丁时间，其下的修改单日期也一同修改

func (s *server) ModDeadLineInPatchsAndTasks(ctx context.Context, in *modDeadlineInfo, modPatchs bool) error {
	tx := db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	newDeadline, patchNo, reqNo := in.newDeadline, in.patchNo, in.reqNo

	// 更新 patch_table 表中的 deadline
	if modPatchs {
		if err := tx.Model(&PatchsInfo{}).
			Where("patch_no = ?", patchNo).
			Update("deadline", newDeadline).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	// 更新 tasklist_table 表中的 deadline
	reqNos := strings.Split(reqNo, ",")
	if err := db.Model(&TaskInfo{}).
		Where("req_no IN ?", reqNos).
		Update("deadline", gorm.Expr("LEAST(deadline, ?)", newDeadline)).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	//TODO: 导入补丁或修改补丁（时间被修改）时调用该函数，这里已经不需要在发布消息了（外层已经发布了）
	//msg := fmt.Sprintf("<%s> -> modifiy patchs's:%s deadline to %s -> <ALL>", in.user, in.patchNo, in.newDeadline)
	//NotificationServer.updateDatabaseAndNotify(msg)

	return nil
}

func (s *server) DelPatch(ctx context.Context, in *pb.DelPatchRequest) (*pb.DelPatchReply, error) {
	if !isAdmin(in.User) {
		return nil, errors.New("you don't have permission")
	}
	patchNo := in.PatchNo

	tx := db.Begin()
	if tx.Error != nil {
		logrus.Errorf("failed to begin transaction: %v", tx.Error)
		return nil, tx.Error
	}

	// 删除 tasklist_table 中相关的记录
	if err := tx.Exec("DELETE FROM tasklist_table WHERE req_no IN (SELECT req_no FROM patch_table WHERE patch_no = ?)", patchNo).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// 删除 patch_table 中的记录
	if err := tx.Where("patch_no = ?", patchNo).Delete(&PatchsInfo{}).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		logrus.Errorf("failed to commit transaction: %v", err)
		return nil, err
	}

	msg := fmt.Sprintf("<%s> -> delete patchs:%s -> <ALL>", in.User, in.PatchNo)
	NotificationServer.updateDatabaseAndNotify(msg)
	return &pb.DelPatchReply{}, nil
}

func (s *server) GetPatchsAll(ctx context.Context, in *pb.GetPatchsAllRequest) (*pb.GetPatchsAllReply, error) {
	var patchs []PatchsInfo
	if er := db.Find(&patchs).Error; er != nil {
		return nil, er
	}
	return &pb.GetPatchsAllReply{Patchs: common.AllPatchsInfoToPbPatchs(patchs)}, nil
}

// TODO:确认表的patch_no是否唯一
func (s *server) GetOnePatchs(ctx context.Context, in *pb.GetOnePatchsRequest) (*pb.GetOnePatchsReply, error) {
	var patchs PatchsInfo
	if err := db.Where("patch_no = ?", in.PatchNo).Find(&patchs).Error; err != nil {
		return nil, err
	}
	return &pb.GetOnePatchsReply{P: common.OnePatchsInfoToPbPatchs(patchs)}, nil
}

func (s *server) QueryPatchsWithField(ctx context.Context, in *pb.QueryPatchsWithFieldRequest) (*pb.QueryPatchsWithFieldReply, error) {
	var patchs []PatchsInfo
	qurey := fmt.Sprintf("%s = ?", in.FieldName)
	if err := db.Where(qurey, in.FieldValue).Find(&patchs).Error; err != nil {
		return nil, err
	}
	return &pb.QueryPatchsWithFieldReply{Ps: common.AllPatchsInfoToPbPatchs(patchs)}, nil
}

func (s *server) ModPatch(ctx context.Context, in *pb.ModPatchRequest) (*pb.ModPatchReply, error) {
	if !isAdmin(in.User) {
		return nil, errors.New("you don't have permission")
	}
	patchs := &PatchsInfo{PatchNo: in.P.PatchNo, ReqNo: in.P.ReqNo}
	if err := db.Model(&patchs).Updates(&PatchsInfo{
		ClientName: in.P.ClientName,
		Reason:     in.P.Reason,
		Describe:   in.P.Describe,
		Sponsor:    in.P.Sponsor,
		State:      in.P.State}).Error; err != nil {
		return nil, err
	}

	var deadline time.Time
	db.Table("patch_table").Where("patch_no = ?", patchs.PatchNo).Scan(&deadline)

	if deadline.Format("2006-01-02") != in.P.Deadline {
		err := s.ModDeadLineInPatchsAndTasks(context.Background(), &modDeadlineInfo{
			patchNo:     in.P.PatchNo,
			newDeadline: in.P.Deadline,
			user:        in.User,
			reqNo:       in.P.ReqNo,
		}, true)
		return nil, err
	} else {
		msg := fmt.Sprintf("<%s> -> modifiy patchs:%s -> <ALL>", in.User, in.P.PatchNo)
		NotificationServer.updateDatabaseAndNotify(msg)
		return &pb.ModPatchReply{}, nil
	}
}
