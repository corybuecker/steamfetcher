package test

type FakeDatabase struct {
	Entry             map[string]interface{}
	DatabaseCreated   string
	ExistingDatabases []string
	TableCreated      string
	ExistingTables    []string
	Rows              map[string]interface{}
}

func (rethinkDB *FakeDatabase) Upsert(databaseName string, tableName string, record map[string]interface{}) error {
	rethinkDB.Entry = record
	return nil
}

func (rethinkDB *FakeDatabase) CreateTable(databaseName string, tableName string) error {
	rethinkDB.TableCreated = tableName
	return nil
}
func (rethinkDB *FakeDatabase) CreateDatabase(databaseName string) error {
	rethinkDB.DatabaseCreated = databaseName
	return nil
}
func (rethinkDB *FakeDatabase) ListDatabases() ([]string, error) {
	return rethinkDB.ExistingDatabases, nil
}
func (rethinkDB *FakeDatabase) ListTables(databaseName string) ([]string, error) {
	return rethinkDB.ExistingTables, nil
}

func (rethinkDB *FakeDatabase) RowsWithoutFields(databaseName string, tableName string, fieldsToExclude []string) ([]map[string]interface{}, error) {
	return []map[string]interface{}{
		{
			"name": "mario",
		},
	}, nil
}

func (rethinkDB *FakeDatabase) RowsWithField(databaseName string, tableName string, fieldsToInclude string) ([]map[string]interface{}, error) {
	return []map[string]interface{}{
		{
			"giantbomb_id": float64(10),
			"id":           float64(34),
		},
	}, nil
}

func (rethinkDB *FakeDatabase) GetRow(databaseName string, tableName string, field string) (map[string]interface{}, error) {
	return rethinkDB.Rows, nil
}
