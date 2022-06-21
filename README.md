## pqxquery library

## Add to your project

```bash
$ go get github.com/daniilty/pgxquery
```
## Usage

```golang
type user struct {
	pgxquery.TableName `db:"users"` // set here the name of related table

	id int `db:"id,primarykey"` // primarykey tag used for updating by in update query
	name string `db:"name,omitempty"` // omitempty is used for update queries(field will be omitted if zero value)
	email string `db:"email"`
}

func insertUsers(...) {
	q, _ := pgxquery.GenerateNamedUpdate(&user{
		id: 1,
		email: "pablo@gmail.com",
	})

	//will result in:
	//update users set email=:email where id=:id;
}
```
