package pgxquery

// TableName - used for setting model's related table name, like this:
//
// struct sample {
//    pgxquery.TableName `db:"samples"`
//    ID int `db:"id,primaryKey"`
//
//    ...
// }
type TableName struct{}
