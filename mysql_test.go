package testutil

import (
	"fmt"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("test testutil", func() {
	Context("private funcs", func() {
		var m *MySQL
		BeforeEach(func() {
			m = &MySQL{
				DataSource:               "root:@tcp(127.0.0.1:3306)/",
				Database:                 fmt.Sprintf("test%d", time.Now().Unix()),
				User:                     "liangzuobin",
				UserHost:                 "%",
				Pwd:                      "123",
				ScriptFile:               "create_table.sql",
				RestartMySQLFirst:        false,
				DropUserAfterTesting:     true,
				DropExistsDatabaseFirst:  true,
				DropDatabaseAfterTesting: true,
			}
		})

		It("m.dbScript() with m.DropExistsDatabseFirst is true", func() {
			Expect(m.dbScripts()).To(HaveLen(5))
		})

		It("m.dbScript() wiht m.DropExistsDatabseFirst is false", func() {
			m.DropExistsDatabaseFirst = false
			Expect(m.dbScripts()).To(HaveLen(4))
		})

		It("m.tableScript() with m.ScriptFile is not empty", func() {
			scripts, err := m.tableScripts()
			Expect(err).NotTo(HaveOccurred())
			Expect(scripts).To(HaveLen(4)) // 3 sentence and a empty one.
		})

		It("m.tableScript() with m.ScriptFile is invalid", func() {
			By("m.ScriptFile == ''")
			m.ScriptFile = "../create_table.sql"
			scripts, err := m.tableScripts()
			Expect(err).To(HaveOccurred())
			Expect(scripts).To(BeEmpty())
		})

		It("m.tableScript() with m.ScriptFile is abs path", func() {
			By("m.ScriptFile == ''")
			m.ScriptFile = "/Users/liangzuobin/Code/Go/src/testutil/create_table.sql"
			scripts, err := m.tableScripts()
			Expect(err).NotTo(HaveOccurred())
			Expect(scripts).To(HaveLen(4))
		})

		It("m.tableScript() with m.ScriptFile is empty", func() {
			By("m.ScriptFile == ''")
			m.ScriptFile = ""
			scripts, err := m.tableScripts()
			Expect(err).NotTo(HaveOccurred())
			Expect(scripts).To(BeEmpty())
		})

		Describe("m.execScript need m.oepnDB() first", func() {
			BeforeEach(func() {
				Expect(m.openDB()).NotTo(HaveOccurred())
			})

			It("m.execScript() with valid script", func() {
				script, err := m.execScripts([]string{"use test"})
				Expect(err).To(BeNil())
				Expect(script).To(BeZero())
			})

			It("m.execScript() with invaild script", func() {
				script, err := m.execScripts([]string{"create table"})
				Expect(err).To(HaveOccurred())
				Expect(script).To(Equal("create table"))
			})

			It("m.execScript() with empty script", func() {
				script, err := m.execScripts([]string{})
				Expect(err).To(HaveOccurred())
				Expect(script).To(BeZero())
			})
		})
	})

	Context("exports func", func() {
		m := NewMySQL()
		BeforeEach(func() {
			m.DataSource = "root:@tcp(127.0.0.1:3306)/"
			m.Database = fmt.Sprintf("test%d", time.Now().Unix())
			m.User = "liangzuobin"
			m.Pwd = "123"
			m.ScriptFile = "create_table.sql"
			m.RestartMySQLFirst = false
			m.DropUserAfterTesting = false
			m.DropExistsDatabaseFirst = true
			m.DropDatabaseAfterTesting = true
		})

		It("restart mysql first", func() {
			m.RestartMySQLFirst = true
			Expect(m.Prepare()).NotTo(HaveOccurred())
			Expect(m.Close()).NotTo(HaveOccurred())
		})

		It("no script file", func() {
			m.ScriptFile = ""
			Expect(m.Prepare()).NotTo(HaveOccurred())
			Expect(m.Close()).NotTo(HaveOccurred())
		})

		It("not drop exists databse", func() {
			m.DropExistsDatabaseFirst = false
			Expect(m.Prepare()).NotTo(HaveOccurred())
			Expect(m.Close()).NotTo(HaveOccurred())
		})

		It("not drop data base after testing", func() {
			m.DropDatabaseAfterTesting = false
			Expect(m.Prepare()).NotTo(HaveOccurred())
			Expect(m.Close()).NotTo(HaveOccurred())
		})

		It("drop user after testing", func() {
			m.DropUserAfterTesting = true
			Expect(m.Prepare()).NotTo(HaveOccurred())
			Expect(m.Close()).NotTo(HaveOccurred())
		})
	})
})
