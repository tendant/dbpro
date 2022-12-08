package dbpro

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"log"
	"reflect"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
)

var fns = template.FuncMap{
	"plus1": func(x int) int {
		return x + 1
	},
}

const sqlTemplate = `INSERT INTO %s ({{$n := len .}}{{range  $i, $e := .}}{{$e}}{{if lt (plus1 $i) $n}},{{end}}{{end}}) VALUES ({{$n := len .}}{{range  $i, $e := .}}:{{$e}}{{if lt (plus1 $i) $n}},{{end}}{{end}}); select ID = convert(bigint, SCOPE_IDENTITY())`

func GenInsertQuery(driverName string, table string, values interface{}) (string, error) {
	fields := reflect.VisibleFields(reflect.TypeOf(values))
	var names []string
	for _, item := range fields {
		// log.Println(item.Name)
		names = append(names, item.Name)
	}
	formatted := fmt.Sprintf(sqlTemplate, table)
	t := template.Must(template.New("abc").Funcs(fns).Parse(formatted))
	buf := new(bytes.Buffer)
	err := t.Execute(buf, names)
	return buf.String(), err
}

func GenInsertValues(entity interface{}) (map[string]interface{}, error) {
	// https://stackoverflow.com/a/67352492
	// val := reflect.ValueOf(entity).Elem() // could be any underlying type
	val := reflect.ValueOf(entity)

	// if its a pointer, resolve its value
	if val.Kind() == reflect.Ptr {
		// log.Println("it is still a pointer")
		val = reflect.Indirect(val)
	}

	// should double check we now have a struct (could still be anything)
	if val.Kind() != reflect.Struct {
		// log.Fatal("unexpected type")
		return nil, errors.New(fmt.Sprintf("Unexpected type: %s", val.Kind()))
	}

	m := make(map[string]interface{})
	// val = reflect.ValueOf(val).Elem()

	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		typeField := val.Type().Field(i)

		f := valueField.Interface()
		val := reflect.ValueOf(f)
		// t := reflect.TypeOf(f)
		// log.Println("type of f:", t)
		// log.Println("value of f:", val.String())
		// log.Println("testtest:", rand.Intn(2) == 1)
		switch val.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			m[typeField.Name] = strconv.FormatInt(val.Int(), 10)
		case reflect.String:
			m[typeField.Name] = val.String()
		case reflect.Bool:
			m[typeField.Name] = val.Bool()
		case reflect.Float32:
			m[typeField.Name] = val.Float()
		case reflect.Float64:
			m[typeField.Name] = val.Float()
		case reflect.Struct:
			// m[typeField.Name] = val.Interface().(sql.NullString).String
			vi := reflect.ValueOf(val.Interface())
			fmt.Println("type name:", vi.Type().Name())
			pkgPath := vi.Type().PkgPath()
			typeName := vi.Type().Name()
			qualifiedTypeName := fmt.Sprintf("%s.%s", pkgPath, typeName)
			switch qualifiedTypeName {
			case "sql.NullString":
				if vi.FieldByName("Valid").Bool() {
					m[typeField.Name] = vi.FieldByName("String").String()
				}
			case "sql.NullBool":
				if vi.FieldByName("Valid").Bool() {
					m[typeField.Name] = vi.FieldByName("Bool").Bool()
				}
			case "time.Time":
				m[typeField.Name] = val.Interface().(time.Time)
			default:
				log.Println("NOT SUPPORTED TYPE:", typeName)
				return nil, errors.New(fmt.Sprintf("Unsupported type: %s", qualifiedTypeName))
			}
		default:
			log.Println("NOT SUPPORTED KIND:", val.Kind())
			return nil, errors.New(fmt.Sprintf("Unsupported kind: %s", val.Kind()))
		}
	}

	return m, nil
}

func reflectStruct(interface{}) {

}

var supportedDrivers = []string{"postgres", "sqlserver"}

// https://play.golang.org/p/Qg_uv_inCek
// contains checks if a string is present in a slice
func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func InsertRow(db *sqlx.DB, table string, entity interface{}) (int64, error) {

	driverName := db.DriverName()
	if !contains(supportedDrivers, driverName) {
		return -1, errors.New(fmt.Sprintf("Driver(%s) is not supported. Supported drivers: %v", driverName, supportedDrivers))
	}
	stmt, err := GenInsertQuery(db.DriverName(), table, entity)
	if err != nil {
		// log.Fatal("Failed generate insert Query!", err)
		return -1, err
	}

	vals, err := GenInsertValues(entity)
	if err != nil {
		// log.Fatal("Failed generate insert Values!", err)
		return -1, err
	}

	rows, err := db.NamedQuery(stmt, vals)
	if err != nil {
		// log.Fatal(fmt.Sprintf("Failed running query: %s", stmt), err)
		return -1, err
	}

	defer rows.Close()
	var id int64
	for rows.Next() {
		if err := rows.Scan(&id); err != nil {
			return -1, err
		}
		// log.Println("Created record: ", id)
		return id, nil
	}
	return -1, errors.New("Failed getting created record id!")
}
