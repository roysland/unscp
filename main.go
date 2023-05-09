package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

const exampleCode = 42271903

type unspscObject struct {
	CommodityCode  int64  `db:"commodity_code" json:"commodity_code"`
	CommodityTitle string `db:"commodity_title" json:"commodity_title"`
	ClassCode      int64  `db:"class_code" json:"class_code"`
	ClassTitle     string `db:"class_title" json:"class_title"`
	FamilyCode     int64  `db:"family_code" json:"family_code"`
	FamilyTitle    string `db:"family_title" json:"family_title"`
	SegmentCode    int64  `db:"segment_code" json:"segment_code"`
	SegmentTitle   string `db:"segment_title" json:"segment_title"`
}

func main() {
	// connect
	db, err := sql.Open("sqlite3", "unscp.sqlite")
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/code/", func(w http.ResponseWriter, r *http.Request) {
		getCode(w, r, db)
	})

	log.Fatal(http.ListenAndServe(":8090", nil))
	defer db.Close()
}

func getCode(w http.ResponseWriter, req *http.Request, db *sql.DB) {
	code := req.URL.Query().Get("number")
	code = code + "%"
	rows, err := runSQL(db, "SELECT * FROM unspsc_reference WHERE commodity_code LIKE ?", code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var objects []unspscObject
	for rows.Next() {
		var obj unspscObject
		err := rows.Scan(&obj.CommodityCode, &obj.CommodityTitle, &obj.ClassCode, &obj.ClassTitle, &obj.FamilyCode, &obj.FamilyTitle, &obj.SegmentCode, &obj.SegmentTitle)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		objects = append(objects, obj)
	}

	if len(objects) == 0 {
		http.Error(w, fmt.Sprintf("No rows found for commodity code %d\n", exampleCode), http.StatusNotFound)
		return
	}

	jsonBytes, err := json.Marshal(objects)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(jsonBytes)
	if err != nil {
		log.Printf("Error writing response: %v\n", err)
	}
}

// Run the sql statement and return the rows
func runSQL(db *sql.DB, sql string, args ...interface{}) (*sql.Rows, error) {
	rows, err := db.Query(sql, args...)
	if err != nil {
		return nil, err
	}
	return rows, nil
}
