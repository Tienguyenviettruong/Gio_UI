package Table

import (
	page "Gio_UI/UI/app"
	"database/sql"
	"gioui.org/text"
	"gioui.org/widget"
	"gioui.org/x/component"
	"log"
	"time"
)

type Page struct {
	widget.List
	*page.Router
	ReadBtn                widget.Clickable
	SearchBtn              widget.Clickable
	Data                   []DataRow
	SearchEditor           widget.Editor
	CurrentPage            int
	RowsPerPage            int
	PreviousBtn            widget.Clickable
	NextBtn                widget.Clickable
	SearchID               string
	ShowDeleteConfirmation bool
	//ConfirmDeleteBtn widget.Clickable
	EditBtns                                                    []widget.Clickable
	AddBtns                                                     []widget.Clickable
	DelBtns                                                     []widget.Clickable
	DeleteID                                                    int // ID của bản ghi cần xóa
	ShowConfirm                                                 bool
	ShowEditConfirmation                                        bool
	SelectedRow                                                 *DataRow
	ConfirmBtn                                                  widget.Clickable
	CancelBtn                                                   widget.Clickable
	FilenameInput, FolderInput, priceInput, tweetInput, IDInput component.TextField
	inputAlignment                                              text.Alignment
	inputAlignmentEnum                                          widget.Enum
}

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

type DataRow struct {
	ID       int
	Filename string
	Modified time.Time
	Folder   string
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
func (p *Page) deleteDataByID(dbPath string, id int) error {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec("DELETE FROM tbl_fileinfo WHERE id = ?", id)
	return err
}
