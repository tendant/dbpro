package dbpro

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
)

type SampleTable struct {
	ColumnA string
	ColumnB string
	ColumnC int
	ColumnD sql.NullString
	ColumnE sql.NullBool
	ColumnF sql.NullFloat64
}

func TestGenInsertQuery(t *testing.T) {
	query, err := GenInsertQuery("mssql", "SampleTable", SampleTable{
		ColumnA: "column1",
		ColumnB: "column2",
		ColumnC: 0,
		// ColumnD: sql.NullString{
		// 	String: "columnd",
		// 	Valid:  true,
		// },
		// ColumnE: sql.NullBool{
		// 	Bool:  false,
		// 	Valid: false,
		// },
	})

	if err != nil {
	}

	expected := `INSERT INTO SampleTable (ColumnA,ColumnB,ColumnC) VALUES (:ColumnA,:ColumnB,:ColumnC); select ID = convert(bigint, SCOPE_IDENTITY())`
	assert.Equal(t, query, expected)

}

func TestGenInsertValues(t *testing.T) {
	actual, err := GenInsertValues(SampleTable{
		ColumnA: "column1",
		ColumnB: "column2",
		ColumnC: 3,
	})

	expected := map[string]interface{}{"ColumnA": "column1", "ColumnB": "column2", "ColumnC": "3"} // "ColumnD": sql.NullString{
	// 	String: "",
	// 	Valid:  false,
	// },
	// "ColumnE": sql.NullBool{
	// 	Bool:  false,
	// 	Valid: false,
	// },

	assert.Equal(t, expected, actual)
	assert.Nil(t, err)
}

func TestGenInsertValuesWithNotNullString(t *testing.T) {
	actual, err := GenInsertValues(SampleTable{
		ColumnA: "column1",
		ColumnB: "column2",
		ColumnC: 3,
		ColumnD: sql.NullString{
			String: "columnd",
			Valid:  true,
		},
	})

	expected := map[string]interface{}{"ColumnA": "column1", "ColumnB": "column2", "ColumnC": "3", "ColumnD": "columnd"} // "ColumnE": sql.NullBool{
	// 	Bool:  false,
	// 	Valid: false,
	// },

	assert.Equal(t, expected, actual)
	assert.Nil(t, err)
}

func TestGenInsertValuesWithNotNullBool(t *testing.T) {
	actual, err := GenInsertValues(SampleTable{
		ColumnA: "column1",
		ColumnB: "column2",
		ColumnC: 3,
		ColumnD: sql.NullString{
			String: "columnd",
			Valid:  true,
		},
		ColumnE: sql.NullBool{
			Bool:  true,
			Valid: true,
		},
	})

	expected := map[string]interface{}{"ColumnA": "column1", "ColumnB": "column2", "ColumnC": "3", "ColumnD": "columnd", "ColumnE": true}

	assert.Equal(t, expected, actual)
	assert.Nil(t, err)
}
