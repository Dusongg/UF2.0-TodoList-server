package main

import (
	"OrderManager/pb"
	"database/sql"
	"github.com/extrame/xls"
	_ "github.com/go-sql-driver/mysql"
	"google.golang.org/grpc"
	"log"
	"net"
	"time"
)

const (
	TASK_STATE_WAIT   = 0
	TASK_STATE_ING    = 1
	TASK_STATE_FINISH = 2

	EMERGENCY_LEVEL_0 = 0
	EMERGENCY_LEVEL_1 = 1
	EMERGENCY_LEVEL_2 = 2
)

var (
	dsn = "root:123123@tcp(127.0.0.1:3306)/OrderManager"
	db  *sql.DB
)

func init() {
	db, _ = sql.Open("mysql", dsn)

	db.SetMaxOpenConns(25)                 // 最大打开连接数
	db.SetMaxIdleConns(25)                 // 最大闲置连接数
	db.SetConnMaxLifetime(5 * time.Minute) // 连接的最大生命周期
	if err := db.Ping(); err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}
}

func main() {

	grpcServer := grpc.NewServer()
	pb.RegisterServiceServer(grpcServer, Server)
	listener, err := net.Listen("tcp", ":8001")
	if err != nil {
		log.Fatal("服务监听失败", err)
	}
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatal(err)
	}
	defer db.Close()
}

type task struct {
	comment            string
	taskId             string
	emergencyLevel     int
	deadline           string
	principal          string
	reqNo              string
	estimatedWorkHours float64
	state              string
	typeId             int
}

type patch struct {
	patchNo    string
	reqNo      string
	describe   string
	clientName string
	deadline   string
	reason     string
	sponsor    string
}

// 规范化导出文件的导入
// 修改单信息： task_id,  principal,s tate,   升级说明： task_id, req_no,comment
// TODO:front
func importXLSForTaskList(file string) {
	// 打开.xls文件
	workbook, err := xls.Open("./规范化导出文件.xls", "utf-8")
	if err != nil {
		log.Fatalf("无法打开文件: %v", err)
	}

	allInsert := make(map[string]*task)

	// 读取“修改单信息”工作表中的数据
	sheet := workbook.GetSheet(2)
	if sheet == nil {
		log.Fatalf("没有找到工作表：修改单信息")
	}

	// 读取B,C,D列的数据
	for i := 2; i <= int(sheet.MaxRow); i++ {
		row := sheet.Row(i)
		colTaskID := row.Col(1)
		colState := row.Col(2)
		colPrincipal := row.Col(3)

		allInsert[colTaskID] = &task{taskId: colTaskID, state: colState, principal: colPrincipal}
	}

	sheet = workbook.GetSheet(3)
	if sheet == nil {
		log.Fatalf("没有找到工作表：升级说明")
	}

	// 读取C,D,I列的数据
	for i := 2; i <= int(sheet.MaxRow); i++ {
		row := sheet.Row(i)
		colTaskID2 := row.Col(2)
		colComment := row.Col(3)
		colReqNo := row.Col(8)

		if task, ok := allInsert[colTaskID2]; ok {
			task.comment = colComment
			task.reqNo = colReqNo
		}
	}

	tasks := make([]task, 0)
	for _, task := range allInsert {
		tasks = append(tasks, *task)
	}

	//TODO:调用rpc

}

// TODO:front
func importXLSForPatchTable(file string) {
	workbook, err := xls.Open("./补丁导出_2024-07-25 13-56-06.xls", "utf-8")
	if err != nil {
		log.Fatalf("无法打开文件: %v", err)
	}

	sheet := workbook.GetSheet(0)
	if sheet == nil {
		log.Fatalf("没有找到工作表：修改单信息")
	}
	patchs := make([]patch, 0)
	for i := 2; i <= int(sheet.MaxRow); i++ {
		row := sheet.Row(i)
		//t, err := time.Parse("20060102", row.Col(14))
		//if err != nil {
		//	log.Println("err to parse time", err)
		//}

		patchs = append(patchs, patch{reqNo: row.Col(0), patchNo: row.Col(1),
			describe: row.Col(2), clientName: row.Col(3), reason: row.Col(12),
			deadline: row.Col(14), sponsor: row.Col(19)})
	}

	//TODO:调用rpc

}
