package dbpro

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type SampleTable struct {
	ColumnA string
	ColumnB string
	ColumnC int
}

func TestGenInsertQuery(t *testing.T) {
	query, err := GenInsertQuery("mssql", "SampleTable", SampleTable{
		ColumnA: "column1",
		ColumnB: "column2",
		ColumnC: 3,
	})

	if err != nil {
	}

	expected := `INSERT INTO SampleTable (ColumnA,ColumnB,ColumnC) VALUES (:ColumnA,:ColumnB,:ColumnC); select ID = convert(bigint, SCOPE_IDENTITY())`
	assert.Equal(t, query, expected)

}
