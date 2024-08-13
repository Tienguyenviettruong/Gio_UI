package Table

import (
	page "Gio_UI/UI/app"
	"database/sql"
	"gioui.org/widget"
	"log"
	"time"
)

func New(router *page.Router) *Page {
	numRows := 10
	return &Page{
		Router:      router,
		CurrentPage: 0,
		RowsPerPage: 10,
		EditBtns:    make([]widget.Clickable, numRows),
		AddBtns:     make([]widget.Clickable, numRows),
		DelBtns:     make([]widget.Clickable, numRows),
	}
}

//type Page struct {
//	widget.List
//	*page.Router
//	ReadBtn                widget.Clickable
//	SearchBtn              widget.Clickable
//	Data                   []DataRow
//	SearchEditor           widget.Editor
//	CurrentPage            int
//	RowsPerPage            int
//	PreviousBtn            widget.Clickable
//	NextBtn                widget.Clickable
//	SearchID               string
//	ShowDeleteConfirmation bool
//	SelectedRowToDelete    *DataRow
//	ConfirmDeleteBtn       widget.Clickable
//	CancelDeleteBtn        widget.Clickable
//
//}

type DataRow struct {
	ID       int
	Filename string
	Modified time.Time
	Folder   string
	//EditBtn  widget.Clickable
	//AddBtn   widget.Clickable
	//DelBtn   widget.Clickable
}

func (p *Page) readDataFromDB(dbFile string) []DataRow {
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, filename, modified, absolute_path FROM tbl_fileinfo")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var data []DataRow
	for rows.Next() {
		var row DataRow
		var modifiedStr string
		if err := rows.Scan(&row.ID, &row.Filename, &modifiedStr, &row.Folder); err != nil {
			log.Fatal(err)
		}
		row.Modified, err = time.Parse("2006-01-02T15:04:05Z", modifiedStr)
		if err != nil {
			log.Fatal(err)
		}
		data = append(data, row)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	return data
}
func (p *Page) searchDataByID(dbFile, searchID string) []DataRow {
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, filename, modified, absolute_path FROM tbl_fileinfo WHERE id = ?", searchID)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var data []DataRow
	for rows.Next() {
		var row DataRow
		var modifiedStr string
		if err := rows.Scan(&row.ID, &row.Filename, &modifiedStr, &row.Folder); err != nil {
			log.Fatal(err)
		}
		row.Modified, err = time.Parse("2006-01-02T15:04:05Z", modifiedStr)
		if err != nil {
			log.Fatal(err)
		}
		data = append(data, row)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	return data
}
func (p *Page) deleteData(id int) {
	// Kết nối đến cơ sở dữ liệu
	db, err := sql.Open("sqlite3", "UI/access.sqlite")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
		return
	}
	defer db.Close()

	// Câu lệnh SQL để xóa dữ liệu theo ID
	stmt, err := db.Prepare("DELETE FROM tbl_fileinfo WHERE id = ?")
	if err != nil {
		log.Fatalf("Failed to prepare statement: %v", err)
		return
	}
	defer stmt.Close()

	// Thực thi câu lệnh SQL
	_, err = stmt.Exec(id)
	if err != nil {
		log.Fatalf("Failed to execute statement: %v", err)
		return
	}

	// Đăng nhập thông báo thành công (có thể thay thế bằng thông báo giao diện người dùng)
	log.Printf("Deleted record with ID %d", id)
}
func (p *Page) DeleteData(id int) {
	// Xử lý xóa dữ liệu từ cơ sở dữ liệu
	// Ví dụ:
	db, err := sql.Open("sqlite3", "UI/access.sqlite")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec("DELETE FROM tbl_fileinfo WHERE id = ?", id)
	if err != nil {
		log.Fatal(err)
	}

	// Cập nhật dữ liệu sau khi xóa
	p.Data = p.readDataFromDB("UI/access.sqlite")
}
