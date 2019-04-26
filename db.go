package main

import (
	"database/sql"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/spf13/viper"
)

func NewDBConn() *sql.DB {
	db, err := sql.Open("mssql", viper.GetString("db.url"))
	must(err)
	return db
}

var SQLSysTable = `SELECT
		s.Name AS SchemaName,
		t.NAME AS TableName,
		p.rows AS RowCounts,
		SUM(a.total_pages) * 8 AS TotalSpaceBytes,
		SUM(a.used_pages) * 8 AS UsedSpaceBytes,
		(SUM(a.total_pages) - SUM(a.used_pages)) * 8 AS UnusedSpaceBytes
	FROM sys.tables t
		INNER JOIN sys.indexes i ON t.OBJECT_ID = i.object_id
		INNER JOIN sys.partitions p ON i.object_id = p.OBJECT_ID AND i.index_id = p.index_id
		INNER JOIN sys.allocation_units a ON p.partition_id = a.container_id
		LEFT OUTER JOIN sys.schemas s ON t.schema_id = s.schema_id
	WHERE 1=1
		AND t.NAME NOT LIKE 'dt%'
		AND t.is_ms_shipped = 0
		AND i.OBJECT_ID > 255
	GROUP BY t.Name, s.Name, p.Rows
	ORDER BY t.Name`

var SQLSysTableRowCount = `SELECT
		s.Name as [Schema],
		t.NAME AS TableName,
		p.rows AS RowCounts
	FROM sys.tables t
		INNER JOIN sys.indexes i ON t.OBJECT_ID = i.object_id
		INNER JOIN sys.partitions p ON i.object_id = p.OBJECT_ID AND i.index_id = p.index_id
		INNER JOIN sys.allocation_units a ON p.partition_id = a.container_id
		LEFT OUTER JOIN sys.schemas s ON t.schema_id = s.schema_id
	WHERE 1=1
		AND t.NAME NOT LIKE 'dt%'
		AND t.is_ms_shipped = 0
		AND i.OBJECT_ID > 255
	GROUP BY t.Name, s.Name, p.Rows
	ORDER BY t.Name`

var SQLSysTableTotalSpaceBytes = `SELECT
		s.Name as [Schema],
		t.NAME AS TableName,
		SUM(a.total_pages) * 8 AS TotalSpaceBytes
	FROM sys.tables t
		INNER JOIN sys.indexes i ON t.OBJECT_ID = i.object_id
		INNER JOIN sys.partitions p ON i.object_id = p.OBJECT_ID AND i.index_id = p.index_id
		INNER JOIN sys.allocation_units a ON p.partition_id = a.container_id
		LEFT OUTER JOIN sys.schemas s ON t.schema_id = s.schema_id
	WHERE 1=1
		AND t.NAME NOT LIKE 'dt%'
		AND t.is_ms_shipped = 0
		AND i.OBJECT_ID > 255
	GROUP BY t.Name, s.Name, p.Rows
	ORDER BY t.Name`

var SQLSysTableUsedSpaceBytes = `SELECT
		s.Name as [Schema],
		t.NAME AS TableName,
		SUM(a.used_pages) * 8 AS UsedSpaceBytes
	FROM sys.tables t
		INNER JOIN sys.indexes i ON t.OBJECT_ID = i.object_id
		INNER JOIN sys.partitions p ON i.object_id = p.OBJECT_ID AND i.index_id = p.index_id
		INNER JOIN sys.allocation_units a ON p.partition_id = a.container_id
		LEFT OUTER JOIN sys.schemas s ON t.schema_id = s.schema_id
	WHERE 1=1
		AND t.NAME NOT LIKE 'dt%'
		AND t.is_ms_shipped = 0
		AND i.OBJECT_ID > 255
	GROUP BY t.Name, s.Name, p.Rows
	ORDER BY t.Name`
