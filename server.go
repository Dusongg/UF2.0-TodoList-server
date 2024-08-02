package main

import (
	"OrderManager/common"
	"OrderManager/pb"
	"context"
	"errors"
	"fmt"
	"log"
)

type server struct {
	pb.UnimplementedServiceServer
}

var Server = &server{}

func (s *server) Login(ctx context.Context, in *pb.LoginRequest) (*pb.LoginReply, error) {
	row, _ := db.Query("select * from user_table where name = ?", in.Name)
	defer row.Close()
	if !row.Next() {
		return &pb.LoginReply{}, errors.New("this user does not exist")
	}
	row2, _ := db.Query("select * from user_table where name = ? and password = ?", in.Name, in.Password)
	defer row2.Close()
	if !row2.Next() {
		return &pb.LoginReply{}, errors.New("wrong password")
	} else {
		return &pb.LoginReply{}, nil
	}
}

func (s *server) Register(ctx context.Context, in *pb.RegisterRequest) (*pb.RegisterReply, error) {
	row, _ := db.Query("select * from user_table where name = ?", in.User.Name)
	defer row.Close()
	if row.Next() {
		return &pb.RegisterReply{}, errors.New("Registration Failed Username already exists")
	}
	_, err := db.Exec("insert into user_table(name, job_no, password, email) values (?, ?, ?, ?)", in.User.Name, in.User.JobNum, in.User.Password, in.User.Email)
	if err != nil {
		return &pb.RegisterReply{}, err
	}
	return &pb.RegisterReply{}, nil
}

func (s *server) GetTaskListAll(ctx context.Context, in *pb.GetTaskListAllRequest) (*pb.GetTaskListAllReply, error) {
	rows, err := db.Query("select * from tasklist_table")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	reply := &pb.GetTaskListAllReply{}
	reply.Tasks, err = common.Process(rows)
	return reply, err
}
func (s *server) GetTaskListOne(ctx context.Context, in *pb.GetTaskListOneRequest) (*pb.GetTaskListOneReply, error) {
	rows, err := db.Query("select * from tasklist_table where principal = ?", in.Name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	reply := &pb.GetTaskListOneReply{}
	reply.Tasks, err = common.Process(rows)
	return reply, err
}
func (s *server) ImportToTaskListTable(ctx context.Context, in *pb.ImportToTaskListRequest) (*pb.ImportToTaskListReply, error) {
	tx, err := db.Begin()
	if err != nil {
		log.Fatalf("failed to begin transaction: %v", err)
	}

	// 批量插入数据
	//task_id,  principal,s tate,   升级说明： task_id, req_no,comment
	//遇到重复任务直接跳过
	//若需要覆盖，则需要将insert加上ON DUPLICATE KEY UPDATE
	var insertCnt int32 = 0
	stmt, err := tx.Prepare("INSERT IGNORE INTO tasklist_table (task_id, req_no, comment, principal, state) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	defer stmt.Close()

	for _, item := range in.Tasks {
		_, err := stmt.Exec(item.TaskId, item.ReqNo, item.Comment, item.Principal, item.State)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		insertCnt++
	}

	// 提交事务
	err = tx.Commit()
	if err != nil {
		log.Fatalf("failed to commit transaction: %v", err)
	}
	return &pb.ImportToTaskListReply{InsertCnt: insertCnt}, nil
}
func (s *server) DelTask(ctx context.Context, in *pb.DelTaskRequest) (*pb.DelTaskReply, error) {
	_, err := db.Exec("delete from tasklist_table where task_id = ?", in.TaskNo)
	if err != nil {
		return nil, err
	}
	return &pb.DelTaskReply{}, nil
}
func (s *server) ModTask(ctx context.Context, in *pb.ModTaskRequest) (*pb.ModTaskReply, error) {
	_, err := db.Exec("update tasklist_table set comment = ?, emergency_level = ?, deadline = ?, principal = ?, estimated_work_hours = ?, state = ?, type = ? where task_id = ?",
		in.T.Comment, in.T.EmergencyLevel, in.T.Deadline, in.T.Principal, in.T.EstimatedWorkHours, in.T.State, in.T.TypeId, in.T.TaskId)
	if err != nil {
		return nil, err
	}
	return &pb.ModTaskReply{}, nil
}
func (s *server) AddTask(ctx context.Context, in *pb.AddTaskRequest) (*pb.AddTaskReply, error) {
	res, err := db.Query("select * from tasklist_table where task_id = ?", in.T.TaskId)
	if err != nil {
		return nil, err
	}
	defer res.Close()
	if res.Next() {
		return nil, errors.New("task already exists")
	}
	_, err = db.Exec("insert into tasklist_table values (?, ?, ?, ?, ?, ?, ?, ?, ?)",
		in.T.Comment, in.T.TaskId, in.T.EmergencyLevel, in.T.Deadline, in.T.Principal, in.T.ReqNo, in.T.EstimatedWorkHours, in.T.State, in.T.TypeId)
	if err != nil {
		return nil, err
	}
	return &pb.AddTaskReply{}, nil
}
func (s *server) QueryTaskWithSQL(ctx context.Context, in *pb.QueryTaskWithSQLRequest) (*pb.QueryTaskWithSQLReply, error) {
	rows, err := db.Query(in.Sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	reply := &pb.QueryTaskWithSQLReply{}
	reply.Tasks, err = common.Process(rows)
	return reply, err
}
func (s *server) QueryTaskWithField(ctx context.Context, in *pb.QueryTaskWithFieldRequest) (*pb.QueryTaskWithFieldReply, error) {
	sql := fmt.Sprintf("select * from tasklist_table where %s = ?", in.Field)
	rows, err := db.Query(sql, in.FieldValue)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	reply := &pb.QueryTaskWithFieldReply{}
	reply.Tasks, err = common.Process(rows)
	return reply, err
}

func (s *server) ImportXLSToPatchTable(ctx context.Context, in *pb.ImportXLSToPatchRequest) (*pb.ImportXLSToPatchReply, error) {
	tx, err := db.Begin()
	if err != nil {
		log.Fatalf("failed to begin transaction: %v", err)
	}

	//批量插入数据
	//遇到重复任务直接跳过
	//若需要覆盖，则需要将insert加上ON DUPLICATE KEY UPDATE
	var insertCnt int32 = 0
	stmt, err := tx.Prepare("INSERT IGNORE INTO patch_table (patch_no, req_no, `describe`,client_name, deadline, reason, sponsor) VALUES (?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	defer stmt.Close()

	for i, item := range in.Patchs {
		_, err := stmt.Exec(item.PatchNo, item.ReqNo, item.Describe, item.ClientName,
			item.Deadline, item.Reason, item.Sponsor)
		if err != nil {
			fmt.Printf("Error inserting row %d: %v\n", i, err) // 打印错误行号和错误信息
			tx.Rollback()
			return nil, err
		}
		insertCnt++
	}

	// 提交事务
	err = tx.Commit()
	if err != nil {
		log.Fatalf("failed to commit transaction: %v", err)
	}
	return &pb.ImportXLSToPatchReply{InsertCnt: insertCnt}, nil
}

// 修改补丁时间，其下的修改单日期也一同修改
func (s *server) ModDeadLineInPatchs(ctx context.Context, in *pb.MDLIPRequest) (*pb.MDLIPReply, error) {
	tx, err := db.Begin()
	if err != nil {
		log.Fatalf("failed to begin transaction: %v", err)
	}

	newDeadline, patchNo := in.NewDeadline, in.PatchNo
	_, err = tx.Exec("update patch_table set deadline = ? where patch_no = ?", newDeadline, patchNo)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	_, err = tx.Exec("update tasklist_table tt set deadline = LEAST(tt.deadline, ?) where tt.req_no in (select req_no from patch_table where patch_no = ?)", newDeadline, patchNo)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		log.Fatalf("failed to commit transaction: %v", err)
	}
	return &pb.MDLIPReply{}, nil
}
func (s *server) DelPatch(ctx context.Context, in *pb.DelPatchRequest) (*pb.DelPatchReply, error) {
	patchNo := in.PatchNo

	tx, err := db.Begin()
	if err != nil {
		log.Fatalf("failed to begin transaction: %v", err)
	}
	_, err = tx.Exec("delete from tasklist_table where tasklist_table.req_no in (select req_no from patch_table where patch_no = ?)", patchNo)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	_, err = tx.Exec("delete from patch_table where patch_no = ?", patchNo)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		log.Fatalf("failed to commit transaction: %v", err)
	}
	return &pb.DelPatchReply{}, nil
}
func (s *server) GetPatchsAll(ctx context.Context, in *pb.GetPatchsAllRequest) (*pb.GetPatchsAllReply, error) {
	rows, err := db.Query("select * from patch_table")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	reply := &pb.GetPatchsAllReply{}
	for rows.Next() {
		task := &pb.Patch{}
		err := rows.Scan(&task.PatchNo, &task.ReqNo, &task.Describe, &task.ClientName,
			&task.Deadline, &task.Reason, &task.Sponsor)
		if err != nil {
			return reply, err
		}
		reply.Patchs = append(reply.Patchs, task)
	}
	return reply, err
}
func (s *server) GetOnePatchs(ctx context.Context, in *pb.GetOnePatchsRequest) (*pb.GetOnePatchsReply, error) {
	res, err := db.Query("select * from patch_table where patch_no = ?", in.PatchNo)
	if err != nil {
		return nil, err
	}
	defer res.Close()
	patch := &pb.Patch{}
	for res.Next() {
		err := res.Scan(&patch.PatchNo, &patch.ReqNo, &patch.Describe, &patch.ClientName,
			&patch.Deadline, &patch.Reason, &patch.Sponsor)
		if err != nil {
			return nil, err
		}
	}
	return &pb.GetOnePatchsReply{P: patch}, nil
}
func (s *server) ModPatch(ctx context.Context, in *pb.ModPatchRequest) (*pb.ModPatchReply, error) {
	_, err := db.Exec("update patch_table set `describe` = ? , client_name = ?, reason = ?, sponsor = ?", in.P.Describe, in.P.ClientName, in.P.Reason, in.P.Sponsor)
	if err != nil {
		return nil, err
	}
	rows, err := db.Query("select deadline from patch_table where patch_no = ?", in.P.PatchNo)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var preDeadline string
		rows.Scan(&preDeadline)
		if preDeadline != in.P.Deadline {
			_, err := s.ModDeadLineInPatchs(context.Background(), &pb.MDLIPRequest{PatchNo: in.P.PatchNo, NewDeadline: in.P.Deadline})
			if err != nil {
				return nil, err
			}
		}
	}
	return &pb.ModPatchReply{}, nil
}
