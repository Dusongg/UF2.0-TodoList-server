package main

import (
	"OrderManager/pb"
	"context"
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
		return &pb.LoginReply{Ok: false}, nil
	}
	row2, _ := db.Query("select * from user_table where name = ? and password = ?", in.Name, in.Password)
	defer row2.Close()
	if !row2.Next() {
		return &pb.LoginReply{Ok: false}, nil
	} else {
		return &pb.LoginReply{Ok: true}, nil
	}
}

func (s *server) GetTaskListAll(ctx context.Context, in *pb.GetTaskListAllRequest) (*pb.GetTaskListAllReply, error) {
	tasks := &pb.GetTaskListAllReply{}
	rows, err := db.Query("select * from tasklist_table")
	if err != nil {
		log.Fatal(err)
		return tasks, err
	}
	defer rows.Close()
	for rows.Next() {
		task := &pb.Task{}
		err := rows.Scan(&task.Comment, &task.TaskId, &task.EmergencyLevel, &task.Deadline,
			&task.Principal, &task.ReqNo, &task.EstimatedWorkHours, &task.State, &task.TypeId)
		if err != nil {
			return tasks, err
		}
		tasks.Tasks = append(tasks.Tasks, task)
	}
	return tasks, nil
}

func (s *server) GetTaskListOne(ctx context.Context, in *pb.GetTaskListOneRequest) (*pb.GetTaskListOneReply, error) {
	tasks := &pb.GetTaskListOneReply{}
	rows, err := db.Query("select * from tasklist_table where principal = ?", in.Name)
	if err != nil {
		return tasks, err
	}
	defer rows.Close()
	for rows.Next() {
		task := &pb.Task{}
		err := rows.Scan(&task.Comment, &task.TaskId, &task.EmergencyLevel, &task.Deadline,
			&task.Principal, &task.ReqNo, &task.EstimatedWorkHours, &task.State, &task.TypeId)
		if err != nil {
			return tasks, err
		}
		tasks.Tasks = append(tasks.Tasks, task)
	}
	return tasks, nil
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
	}

	// 提交事务
	err = tx.Commit()
	if err != nil {
		log.Fatalf("failed to commit transaction: %v", err)
	}
	return &pb.ImportToTaskListReply{}, nil
}

func (s *server) ImportXLSToPatchTable(ctx context.Context, in *pb.ImportXLSToPatchRequest) (*pb.ImportXLSToPatchReply, error) {
	//for _, t := range patchs {
	//	fmt.Printf("reqNo:%s\n patchNo: %s\n describe: %s\n clientName: %s\n reason: %s\n deadline: %s\n sponser: %s\n",
	//		t.reqNo, t.patchNo, t.describe, t.clientName, t.reason, t.deadline, t.sponsor)
	//
	//}

	tx, err := db.Begin()
	if err != nil {
		log.Fatalf("failed to begin transaction: %v", err)
	}

	//批量插入数据
	//遇到重复任务直接跳过
	//若需要覆盖，则需要将insert加上ON DUPLICATE KEY UPDATE
	stmt, err := tx.Prepare("INSERT IGNORE INTO patch_table (patch_no, req_no, `describe`,client_name, deadline, reason, sponsor) VALUES (?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	defer stmt.Close()

	for _, item := range in.Patchs {
		_, err := stmt.Exec(item.PatchNo, item.ReqNo, item.Describe, item.ClientName,
			item.Deadline, item.Reason, item.Sponsor)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	// 提交事务
	err = tx.Commit()
	if err != nil {
		log.Fatalf("failed to commit transaction: %v", err)
	}
	return &pb.ImportXLSToPatchReply{}, nil
}
