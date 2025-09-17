// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package dm

import (
	"context"
	"fmt"
	"strings"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/util/gutil"
)

const (
	// tableFieldsSqlTmp = `SELECT * FROM ALL_TAB_COLUMNS WHERE Table_Name= '%s' AND OWNER = '%s'`
	tableFieldsSqlTmp = `SELECT ATC.COLUMN_NAME, ATC.DATA_TYPE, ATC.NULLABLE, ATC.DATA_DEFAULT, DJC.DATA_TYPE AS DJC_DATA_TYPE, DCC.COMMENTS, CASE WHEN PK.COLUMN_NAME IS NOT NULL THEN 'P' ELSE NULL END AS CONSTRAINT_TYPE FROM ALL_TAB_COLUMNS ATC LEFT JOIN DBA_JSON_COLUMNS DJC ON ATC.OWNER = DJC.OWNER AND ATC.TABLE_NAME = DJC.TABLE_NAME AND ATC.COLUMN_NAME = DJC.COLUMN_NAME LEFT JOIN DBA_COL_COMMENTS DCC ON ATC.OWNER = DCC.OWNER AND ATC.TABLE_NAME = DCC.TABLE_NAME AND ATC.COLUMN_NAME = DCC.COLUMN_NAME LEFT JOIN (SELECT acc.OWNER, acc.TABLE_NAME, acc.COLUMN_NAME FROM ALL_CONSTRAINTS AC JOIN ALL_CONS_COLUMNS ACC ON AC.OWNER = ACC.OWNER AND AC.TABLE_NAME = ACC.TABLE_NAME AND AC.CONSTRAINT_NAME = ACC.CONSTRAINT_NAME WHERE AC.CONSTRAINT_TYPE = 'P') PK ON ATC.OWNER = PK.OWNER AND ATC.TABLE_NAME = PK.TABLE_NAME AND ATC.COLUMN_NAME = PK.COLUMN_NAME WHERE ATC.TABLE_NAME = '%s' AND ATC.OWNER = '%s'`
)

// TableFields retrieves and returns the fields' information of specified table of current schema.
func (d *Driver) TableFields(
	ctx context.Context, table string, schema ...string,
) (fields map[string]*gdb.TableField, err error) {
	var (
		result gdb.Result
		link   gdb.Link
		// When no schema is specified, the configuration item is returned by default
		usedSchema = gutil.GetOrDefaultStr(d.GetSchema(), schema...)
	)
	// When usedSchema is empty, return the default link
	if link, err = d.SlaveLink(usedSchema); err != nil {
		return nil, err
	}
	// The link has been distinguished and no longer needs to judge the owner
	result, err = d.DoSelect(
		ctx, link,
		fmt.Sprintf(
			tableFieldsSqlTmp,
			strings.ToUpper(table),
			strings.ToUpper(d.GetSchema()),
		),
	)
	if err != nil {
		return nil, err
	}
	fields = make(map[string]*gdb.TableField)
	for i, m := range result {
		// m[NULLABLE] returns "N" "Y"
		// "N" means not null
		// "Y" means could be null
		var nullable bool
		if m["NULLABLE"].String() != "N" {
			nullable = true
		}

		dataType := m["DATA_TYPE"].String()
		if !m["DJC_DATA_TYPE"].IsNil() && (m["DJC_DATA_TYPE"].String() != "JSON" || m["DJC_DATA_TYPE"].String() != "JSONB") {
			dataType = m["DJC_DATA_TYPE"].String()
		}

		key := ""
		if !m["CONSTRAINT_TYPE"].IsNil() && m["CONSTRAINT_TYPE"].String() == "P" {
			key = "PRI"
		}

		fields[m["COLUMN_NAME"].String()] = &gdb.TableField{
			Index: i,
			Name:  m["COLUMN_NAME"].String(),
			// Type:    m["DATA_TYPE"].String(),
			Type:    dataType,
			Null:    nullable,
			Default: m["DATA_DEFAULT"].Val(),
			// Key:     m["Key"].String(),
			Key: key,
			// Extra:   m["Extra"].String(),
			// Comment: m["Comment"].String(),
			Comment: m["COMMENTS"].String(),
		}
	}
	return fields, nil
}
