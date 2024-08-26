package common

import (
	"OrderManager/models"
	"OrderManager/pb"
	"time"
)

func AllTaskInfoToPbTask(tasks []models.TaskInfo) []*pb.Task {
	result := make([]*pb.Task, len(tasks))
	for i, task := range tasks {
		result[i] = &pb.Task{
			Comment:            task.Comment,
			TaskId:             task.TaskID,
			EmergencyLevel:     int32(task.EmergencyLevel),
			Deadline:           task.Deadline.Format("2006-01-02"), // 格式化日期
			Principal:          task.Principal,
			ReqNo:              task.ReqNo,
			EstimatedWorkHours: int32(task.EstimatedWorkHours),
			State:              task.State,
			TypeId:             int32(task.Type),
		}
	}
	return result
}

func OneTaskInfoToPbTask(task models.TaskInfo) *pb.Task {
	return &pb.Task{
		Comment:            task.Comment,
		TaskId:             task.TaskID,
		EmergencyLevel:     int32(task.EmergencyLevel),
		Deadline:           task.Deadline.Format("2006-01-02"), // 格式化日期
		Principal:          task.Principal,
		ReqNo:              task.ReqNo,
		EstimatedWorkHours: int32(task.EstimatedWorkHours),
		State:              task.State,
		TypeId:             int32(task.Type),
	}

}

func AllPbTaskToTaskInfo(tasks []*pb.Task) []models.TaskInfo {
	res := make([]models.TaskInfo, len(tasks))
	for i, task := range tasks {
		t, _ := time.Parse("2006-01-02", task.Deadline)
		res[i] = models.TaskInfo{
			Comment:            task.Comment,
			TaskID:             task.TaskId,
			EmergencyLevel:     int(task.EmergencyLevel),
			Deadline:           t, // 格式化日期
			Principal:          task.Principal,
			ReqNo:              task.ReqNo,
			EstimatedWorkHours: int8(float64(task.EstimatedWorkHours)),
			State:              task.State,
			Type:               int(task.TypeId),
		}
	}
	return res
}

func OnePbTaskToTaskInfo(task *pb.Task) *models.TaskInfo {
	t, _ := time.Parse("2006-01-02", task.Deadline)
	res := &models.TaskInfo{
		Comment:            task.Comment,
		TaskID:             task.TaskId,
		EmergencyLevel:     int(task.EmergencyLevel),
		Deadline:           t, // 格式化日期
		Principal:          task.Principal,
		ReqNo:              task.ReqNo,
		EstimatedWorkHours: int8(task.EstimatedWorkHours),
		State:              task.State,
		Type:               int(task.TypeId),
	}
	return res
}

func AllPbPatchsToPatchsInfo(patchs []*pb.Patch) []models.PatchsInfo {
	result := make([]models.PatchsInfo, len(patchs))
	for i, p := range patchs {
		t, _ := time.Parse("2006-01-02", p.Deadline)
		result[i] = models.PatchsInfo{
			PatchNo:    p.PatchNo,
			ReqNo:      p.ReqNo,
			Describe:   p.Describe,
			ClientName: p.ClientName,
			Deadline:   t,
			Reason:     p.Reason,
			Sponsor:    p.Sponsor,
			State:      p.State,
		}
	}
	return result
}

func AllPatchsInfoToPbPatchs(patchs []models.PatchsInfo) []*pb.Patch {
	result := make([]*pb.Patch, len(patchs))
	for i, p := range patchs {
		result[i] = &pb.Patch{
			PatchNo:    p.PatchNo,
			ClientName: p.ClientName,
			Reason:     p.Reason,
			ReqNo:      p.ReqNo,
			Describe:   p.Describe,
			Sponsor:    p.Sponsor,
			Deadline:   p.Deadline.Format("2006-01-02"),
			State:      p.State,
		}
	}
	return result
}

func OnePatchsInfoToPbPatchs(patchs models.PatchsInfo) *pb.Patch {
	return &pb.Patch{
		PatchNo:    patchs.PatchNo,
		ClientName: patchs.ClientName,
		Reason:     patchs.Reason,
		ReqNo:      patchs.ReqNo,
		Describe:   patchs.Describe,
		Sponsor:    patchs.Sponsor,
		Deadline:   patchs.Deadline.Format("2006-01-02"),
		State:      patchs.State,
	}

}
