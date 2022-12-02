# dbpro
dbpro is a library which provides a set of utility function to simplify the task of generating data in database tables.

## install

    go get github.com/tendant/dbpro
    
## usage

```go
type Person struct {
    FirstName string
    LastName  string
    Email     string
}

person := Person{
    FirstName: "test first name"
    LastName: "test last name"
    Email: "test@example.com"
}

// query, err := dbpro.GenInsertQuery("postgres", "person"", person)

// vals, err := dbpro.GenInsertValues(person)

// rows, err := db.NamedQuery(stmt, vals)

db, err := sqlx.Connect("postgres", "user=foo dbname=bar sslmode=disable")
if err != nil {
     log.Fatalln(err)
}

id, err := dbpro.InsertRow(db, "person", person)
```
