package common

import (
	"OrderManager/pb"
	"database/sql"
)

func Process(rows *sql.Rows) ([]*pb.Task, error) {
	result := make([]*pb.Task, 0)
	for rows.Next() {
		task := &pb.Task{}
		err := rows.Scan(&task.Comment, &task.TaskId, &task.EmergencyLevel, &task.Deadline,
			&task.Principal, &task.ReqNo, &task.EstimatedWorkHours, &task.State, &task.TypeId)
		if err != nil {
			return result, err
		}
		result = append(result, task)
	}
	return result, nil
}
