package main

import (
	"database/sql"
	"github.com/extrame/xls"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"time"
)

const (
	TASK_STATE_WAIT   = 0
	TASK_STATE_ING    = 1
	TASK_STATE_FINISH = 2
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

	//rows, err := db.Query("select * from Login")
	//if err != nil {
	//	log.Fatal("query err", err)
	//}
	//defer rows.Close()
	//for rows.Next() {
	//	var id int
	//	var username, password, email string
	//	err = rows.Scan(&id, &username, &password, &email)
	//	if err != nil {
	//		fmt.Println("Error scanning row:", err)
	//		return
	//	}
	//	fmt.Printf("ID: %d, Username: %s, Password: %s, Email: %s\n", id, username, password, email)
	//}
	//improtXLSToTaskList("xxx")
	importXLSToPatch("xxx")

	defer db.Close()
}

func Login(name string, password string) (bool, string) {
	row, _ := db.Query("select * from user where name = ?", name)
	defer row.Close()
	if !row.Next() {
		return false, "用户不存在"
	}
	row2, _ := db.Query("select * from user where name = ? and password = ?", name, password)
	defer row2.Close()
	if !row2.Next() {
		return false, "密码错误"
	} else {
		return true, "登录成功"
	}
}

type task struct {
	comment            string
	taskId             string
	emergencyLevel     int
	deadline           time.Time
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
	deadline   time.Time
	reason     string
	sponsor    string
}

func getTaskListAll() []task {
	rows, err := db.Query("select * from TaskList")
	if err != nil {
		log.Fatal(err)
		return nil
	}
	defer rows.Close()
	var tasks []task
	for rows.Next() {
		var t task
		rows.Scan(&t.comment, &t.taskId, &t.emergencyLevel, &t.deadline, &t.principal, &t.reqNo, &t.estimatedWorkHours, &t.state, &t.typeId)
		tasks = append(tasks, t)
	}
	return tasks

}
func getTaskListOne(name string) []task {
	rows, err := db.Query("select * from TaskList where principal = ?", name)
	if err != nil {
		log.Fatal("query err", err)
	}
	result := make([]task, 0)
	defer rows.Close()
	for rows.Next() {
		t := task{}
		rows.Scan(&t.comment, &t.taskId, &t.emergencyLevel, &t.deadline, &t.principal, &t.reqNo, &t.estimatedWorkHours, &t.state, &t.typeId)
		result = append(result, t)
	}
	return result
}

// 规范化导出文件的导入
// 修改单信息： task_id,  principal,s tate,   升级说明： task_id, req_no,comment
func improtXLSToTaskList(file string) {
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

	// 读取“升级说明”工作表中的数据
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

	// 输出结果以验证
	//for _, task := range allInsert {
	//	fmt.Printf("TaskID: %s, State: %s, Principal: %s, Comment: %s, ReqNo: %s\n",
	//		task.taskId, task.state, task.principal, task.comment, task.reqNo)
	//}

	tx, err := db.Begin()
	if err != nil {
		log.Fatalf("failed to begin transaction: %v", err)
	}

	// 批量插入数据
	//task_id,  principal,s tate,   升级说明： task_id, req_no,comment
	//遇到重复任务直接跳过
	//若需要覆盖，则需要将insert加上ON DUPLICATE KEY UPDATE
	stmt, err := tx.Prepare("INSERT IGNORE INTO tasklist (task_id, req_no, comment, principal, state) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		tx.Rollback()
		log.Fatalf("failed to prepare statement: %v", err)
	}
	defer stmt.Close()

	for _, item := range allInsert {
		_, err := stmt.Exec(item.taskId, item.reqNo, item.comment, item.principal, item.state)
		if err != nil {
			tx.Rollback()
			log.Fatalf("failed to execute statement: %v", err)
		}
	}

	// 提交事务
	err = tx.Commit()
	if err != nil {
		log.Fatalf("failed to commit transaction: %v", err)
	}
}

// req_no:A   patch_no:B  describe：C client_name:D  reason: M deadline:O  sponser:T
func importXLSToPatch(path string) {
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
		t, err := time.Parse("20060102", row.Col(14))
		if err != nil {
			log.Println("err to parse time", err)
		}

		patchs = append(patchs, patch{reqNo: row.Col(0), patchNo: row.Col(1),
			describe: row.Col(2), clientName: row.Col(3), reason: row.Col(12),
			deadline: t, sponsor: row.Col(19)})
	}

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
	stmt, err := tx.Prepare("INSERT IGNORE INTO patch (patch_no, req_no, `describe`,client_name, deadline, reason, sponsor) VALUES (?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		tx.Rollback()
		log.Fatalf("failed to prepare statement: %v", err)
	}
	defer stmt.Close()

	for _, item := range patchs {
		_, err := stmt.Exec(item.patchNo, item.reqNo, item.describe, item.clientName,
			item.deadline.Format("2006-01-02"), item.reason, item.sponsor)
		if err != nil {
			tx.Rollback()
			log.Fatalf("failed to execute statement: %v", err)
		}
	}

	// 提交事务
	err = tx.Commit()
	if err != nil {
		log.Fatalf("failed to commit transaction: %v", err)
	}
}

// 规范化导出文件的导入
// 修改单信息： task_id,  principal,s tate,   升级说明： task_id, req_no,comment
//TODO:xlsx格式导入？
//func improtXLSXToTaskList(exclFile string) {
//	//f, err := excelize.OpenFile(exclFile)
//	//debug
//	f, err := excelize.OpenFile("./规范化导出文件.xls")
//
//	if err != nil {
//		fmt.Println("open file err: ", err)
//		return
//	}
//	defer f.Close()
//
//	allInsert := make(map[string]*task)
//	cols, err := f.GetCols("修改单信息")
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//
//	//B,C,D
//	colTaskID, colState, colPrincipal := cols[1], cols[2], cols[3]
//	for i := 2; i < len(colTaskID); i++ {
//		allInsert[colTaskID[i]] = &task{taskId: colTaskID[i], state: colState[i], principal: colPrincipal[i]}
//	}
//
//	cols, err = f.GetCols("升级说明")
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//
//	//C,D,I
//	colTaskID2, colComment, colReqNo := cols[2], cols[3], cols[8]
//	for i := 2; i < len(colTaskID2); i++ {
//		if _, ok := allInsert[colTaskID2[i]]; !ok {
//			allInsert[colTaskID2[i]].comment = colComment[i]
//			allInsert[colTaskID2[i]].reqNo = colReqNo[i]
//		}
//	}
//
//	//debug
//	for _, task := range allInsert {
//		fmt.Println(task)
//	}
//
//}
