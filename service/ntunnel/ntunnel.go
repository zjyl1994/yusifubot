package ntunnel

// Navicat Http tunnel for SQLite
// Tested on Navicat Premium Lite 17.1

import (
	"bytes"
	"database/sql"
	"encoding/base64"
	"encoding/binary"
	"io"
	"net/http"
	"strings"
)

type nTunnel struct {
	opendb  func(string) (*sql.DB, error)
	errfunc func(error) (int, string)
	version string
}

func NewNTunnel(openFunc func(string) (*sql.DB, error), getErrCode func(error) (int, string), version string) *nTunnel {
	return &nTunnel{
		opendb:  openFunc,
		errfunc: getErrCode,
		version: version,
	}
}

func (nt *nTunnel) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(2 * 1024 * 1024)

	var action string
	var queries []string
	var encodeBase64 bool
	var dbFile string

	for key, value := range r.Form {
		switch key {
		case "actn":
			action = value[0]
		case "encodeBase64":
			encodeBase64 = value[0] == "1"
		case "q[]":
			queries = append(queries, value...)
		case "dbfile":
			dbFile = value[0]
		}
	}

	if encodeBase64 {
		var decoded []string
		for _, v := range queries {
			ret, err := base64.StdEncoding.DecodeString(v)
			if err == nil {
				decoded = append(decoded, string(ret))
			}
		}
		queries = decoded
	}

	w.Header().Set("Content-Type", "text/plain; charset=x-user-defined")

	db, err := nt.opendb(dbFile)
	if err != nil {
		nt.writeHeader(w, 202)
		nt.writeBlock(w, []byte(err.Error()))
		return
	}

	nt.writeHeader(w, 0)

	switch action {
	case "C":
		for range 3 {
			nt.writeBlock(w, []byte(nt.version))
		}
	case "Q":
		for i, v := range queries {
			if v == "" {
				continue
			}

			nt.runSQL(w, db, v)

			if i < len(queries)-1 {
				w.Write([]byte{0x01})
			} else {
				w.Write([]byte{0x00})
			}
		}
	}
}

func (nt *nTunnel) runSQL(w io.Writer, db *sql.DB, sql string) {
	if nt.isExec(sql) {
		ret, err := db.Exec(sql)
		if err != nil {
			errno, errmsg := nt.errfunc(err)
			nt.writeResultSetHeader(w, errno, 0, 0, 0, 0)
			nt.writeBlock(w, []byte(errmsg))
			return
		}
		affectRows, _ := ret.RowsAffected()
		insertId, _ := ret.LastInsertId()
		nt.writeResultSetHeader(w, 0, int(affectRows), int(insertId), 0, 0)
		w.Write([]byte{0x00})
		return
	}

	ret, err := db.Query(sql)
	if err != nil {
		errno, errmsg := nt.errfunc(err)
		nt.writeResultSetHeader(w, errno, 0, 0, 0, 0)
		nt.writeBlock(w, []byte(errmsg))
		return
	}
	defer ret.Close()

	cols, err := ret.Columns()
	if err != nil {
		nt.writeResultSetHeader(w, 1, 0, 0, 0, 0)
		nt.writeBlock(w, []byte(err.Error()))
		return
	}

	rowCounter := 0
	var resultBuffer bytes.Buffer
	for ret.Next() {
		values := make([]interface{}, len(cols))
		rawValues := make([][]byte, len(cols))
		for i := range values {
			values[i] = &rawValues[i]
		}

		err = ret.Scan(values...)
		if err != nil {
			nt.writeResultSetHeader(w, 1, 0, 0, 0, 0)
			nt.writeBlock(w, []byte(err.Error()))
			return
		}
		for _, v := range rawValues {
			if v == nil {
				resultBuffer.Write([]byte{0xFF})
			} else {
				nt.writeBlock(&resultBuffer, v)
			}
			// database/sql can't get raw sqlite type const value, use 0 instead.
			binary.Write(&resultBuffer, binary.BigEndian, uint32(0))
		}
		rowCounter++
	}

	nt.writeResultSetHeader(w, 0, 0, 0, len(cols), rowCounter)
	nt.writeFieldsHeader(w, cols)
	resultBuffer.WriteTo(w)
}

func (nt *nTunnel) isExec(sql string) bool {
	segments := make([]string, 0)
	for _, seg := range strings.Fields(strings.ToLower(sql)) {
		val := strings.TrimSpace(seg)
		if len(val) > 0 {
			segments = append(segments, val)
		}
	}

	if len(segments) == 0 {
		return true
	}
	if segments[0] == "select" || segments[0] == "with" || segments[0] == "explain" {
		return false
	} else if segments[0] == "pragma" {
		return strings.ContainsRune(sql, '=')
	} else {
		return true
	}
}

func (nt *nTunnel) writeHeader(w io.Writer, errno int) {
	const versionNum = 203 //CC_HTTP_TUNNEL_SCRIPT_LATEST_VERSION_SQLITE

	binary.Write(w, binary.BigEndian, uint32(1111))
	binary.Write(w, binary.BigEndian, uint16(versionNum))
	binary.Write(w, binary.BigEndian, uint32(errno))
	w.Write(make([]byte, 6))
}

func (nt *nTunnel) writeBlock(w io.Writer, val []byte) {
	l := len(val)
	if l < 254 {
		w.Write([]byte{byte(l)})
		w.Write(val)
	} else {
		w.Write([]byte{0xFE})
		binary.Write(w, binary.BigEndian, uint32(l))
		w.Write(val)
	}
}

func (nt *nTunnel) writeResultSetHeader(w io.Writer, errno, affectRows, insertId, numFields, numRows int) {
	binary.Write(w, binary.BigEndian, uint32(errno))
	binary.Write(w, binary.BigEndian, uint32(affectRows))
	binary.Write(w, binary.BigEndian, uint32(insertId))
	binary.Write(w, binary.BigEndian, uint32(numFields))
	binary.Write(w, binary.BigEndian, uint32(numRows))
	w.Write(make([]byte, 12))
}

func (nt *nTunnel) writeFieldsHeader(w io.Writer, fieldNames []string) {
	for _, v := range fieldNames {
		nt.writeBlock(w, []byte(v))
		w.Write([]byte{0x00})
		binary.Write(w, binary.BigEndian, uint32(0))
		binary.Write(w, binary.BigEndian, uint32(0))
		binary.Write(w, binary.BigEndian, uint32(0))
	}
}
