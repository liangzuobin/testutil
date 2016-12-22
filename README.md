```go
var _ = Describe("test...", func() {
	db := &testutil.MySQL{
		DataSource:               "root:@tcp(127.0.0.1:3306)/",
		Database:                 fmt.Sprintf("test%d", time.Now().Unix()),
		User:                     "liangzuobin",
		Pwd:                      "123",
		ScriptFile:               "create_table.sql",
		RestartMySQLFirst:        false,
		DropExistsDatabaseFirst:  true,
		DropDatabaseAfterTesting: true,
	}

	BeforeSuite(func() {
		db.Prepare()
	})

	AfterSuite(func() {
		db.Close()
	})

	It("unit test...", func() {})
	...
})

```
