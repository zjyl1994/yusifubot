package sqliteadmin

import "strings"

func useQuery(sql string) bool {
	sqlFields := strings.Fields(sql)
	if len(sqlFields) == 0 {
		return false
	}
	firstBlock := strings.ToLower(sqlFields[0])
	return firstBlock == "select" ||
		firstBlock == "with" ||
		firstBlock == "pragma" ||
		firstBlock == "explain"
}
