package dao

import (
	"database/sql"
	"os"
	"testing"

	"github.com/bojand/ghz-web/model"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/stretchr/testify/assert"
)

const dbName = "../test/project_test.db"

func createTestData() error {
	os.Remove(dbName)

	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		return err
	}
	defer db.Close()

	sqlStmt := `CREATE TABLE "projects" ("id" integer primary key autoincrement,"created_at" datetime,"updated_at" datetime,"deleted_at" datetime,"name" varchar(255),"description" varchar(255) );`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		return err
	}

	sqlStmt = `CREATE INDEX idx_projects_deleted_at ON "projects"(deleted_at);`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		return err
	}

	sqlStmt = `CREATE UNIQUE INDEX uix_projects_email ON "projects"("name");`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		return err
	}

	sqlStmt = `INSERT INTO "projects" ("created_at","updated_at","deleted_at","name","description") VALUES ('2018-05-06 20:42:37','2018-05-06 20:42:37',NULL,'testproject123','test project description goes here');`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		return err
	}

	return nil
}

func TestProjectService_FindByID(t *testing.T) {
	defer os.Remove(dbName)

	err := createTestData()
	if err != nil {
		assert.FailNow(t, err.Error())
	}

	db, err := gorm.Open("sqlite3", dbName)
	if err != nil {
		assert.FailNow(t, err.Error())
	}
	defer db.Close()

	dao := ProjectService{DB: db}

	t.Run("test existing", func(t *testing.T) {
		p := model.Project{}
		err := dao.FindByID(1, &p)

		assert.NoError(t, err)
		assert.Equal(t, uint(1), p.ID)
		assert.Equal(t, "testproject123", p.Name)
		assert.Equal(t, "test project description goes here", p.Description)
		assert.NotNil(t, p.CreatedAt)
		assert.NotNil(t, p.UpdatedAt)
		assert.Nil(t, p.DeletedAt)
	})

	t.Run("test not found", func(t *testing.T) {
		p := model.Project{}
		err := dao.FindByID(2, &p)

		assert.Error(t, err)
		assert.Equal(t, uint(0), p.ID)
		assert.Equal(t, "", p.Name)
		assert.Equal(t, "", p.Description)
	})
}

func TestProjectService_FindByName(t *testing.T) {
	defer os.Remove(dbName)

	err := createTestData()
	if err != nil {
		assert.FailNow(t, err.Error())
	}

	db, err := gorm.Open("sqlite3", dbName)
	if err != nil {
		assert.FailNow(t, err.Error())
	}
	defer db.Close()

	dao := ProjectService{DB: db}

	t.Run("test existing", func(t *testing.T) {
		p := model.Project{}
		err := dao.FindByName("testproject123", &p)

		assert.NoError(t, err)
		assert.Equal(t, uint(1), p.ID)
		assert.Equal(t, "testproject123", p.Name)
		assert.Equal(t, "test project description goes here", p.Description)
		assert.NotNil(t, p.CreatedAt)
		assert.NotNil(t, p.UpdatedAt)
		assert.Nil(t, p.DeletedAt)
	})

	t.Run("test not found", func(t *testing.T) {
		p := model.Project{}
		err := dao.FindByName("testproject999", &p)

		assert.Error(t, err)
		assert.Equal(t, uint(0), p.ID)
		assert.Equal(t, "", p.Name)
		assert.Equal(t, "", p.Description)
	})
}

func TestProjectService_Create(t *testing.T) {
	defer os.Remove(dbName)

	err := createTestData()
	if err != nil {
		assert.FailNow(t, err.Error())
	}

	db, err := gorm.Open("sqlite3", dbName)
	if err != nil {
		assert.FailNow(t, err.Error())
	}
	defer db.Close()

	dao := ProjectService{DB: db}

	t.Run("test new", func(t *testing.T) {
		p := model.Project{
			Name:        "TestProj111 ",
			Description: "Test Description Asdf ",
		}
		err := dao.Create(&p)

		assert.NoError(t, err)

		assert.NotZero(t, p.ID)
		assert.Equal(t, "testproj111", p.Name)
		assert.Equal(t, "Test Description Asdf", p.Description)
		assert.NotNil(t, p.CreatedAt)
		assert.NotNil(t, p.UpdatedAt)
		assert.Nil(t, p.DeletedAt)

		p2 := model.Project{}
		err = dao.FindByID(p.ID, &p2)

		assert.Equal(t, "testproj111", p2.Name)
		assert.Equal(t, "Test Description Asdf", p2.Description)
		assert.NotNil(t, p2.CreatedAt)
		assert.NotNil(t, p2.UpdatedAt)
		assert.Nil(t, p2.DeletedAt)
	})

	t.Run("test new with empty name", func(t *testing.T) {
		p := model.Project{
			Description: "Test Description Asdf 2",
		}
		err := dao.Create(&p)

		assert.NoError(t, err)

		assert.NotZero(t, p.ID)
		assert.NotEmpty(t, p.Name)
		assert.Equal(t, "Test Description Asdf 2", p.Description)
		assert.NotNil(t, p.CreatedAt)
		assert.NotNil(t, p.UpdatedAt)
		assert.Nil(t, p.DeletedAt)

		p2 := model.Project{}
		err = dao.FindByID(p.ID, &p2)

		assert.Equal(t, p.Name, p2.Name)
		assert.Equal(t, "Test Description Asdf 2", p2.Description)
		assert.NotNil(t, p2.CreatedAt)
		assert.NotNil(t, p2.UpdatedAt)
		assert.Nil(t, p2.DeletedAt)
	})

	t.Run("test new with ID", func(t *testing.T) {
		p := model.Project{
			Name:        " FooProject ",
			Description: " Bar Desc ",
		}
		p.ID = 123

		err := dao.Create(&p)

		assert.NoError(t, err)

		assert.Equal(t, uint(123), p.ID)
		assert.Equal(t, "fooproject", p.Name)
		assert.Equal(t, "Bar Desc", p.Description)
		assert.NotNil(t, p.CreatedAt)
		assert.NotNil(t, p.UpdatedAt)
		assert.Nil(t, p.DeletedAt)

		p2 := model.Project{}
		err = dao.FindByID(p.ID, &p2)

		assert.Equal(t, uint(123), p2.ID)
		assert.Equal(t, "fooproject", p2.Name)
		assert.Equal(t, "Bar Desc", p2.Description)
		assert.NotNil(t, p.CreatedAt)
		assert.NotNil(t, p.UpdatedAt)
		assert.Nil(t, p.DeletedAt)
	})

	t.Run("should fail with same ID", func(t *testing.T) {
		p := model.Project{
			Name:        "ACME",
			Description: "Lorem Ipsum",
		}
		p.ID = 123

		err := dao.Create(&p)

		assert.Error(t, err)
	})

	t.Run("should fail with same name", func(t *testing.T) {
		p := model.Project{
			Name:        "FooProject",
			Description: "Lorem Ipsum",
		}
		err := dao.Create(&p)

		assert.Error(t, err)
	})
}

func TestProjectService_Update(t *testing.T) {
	defer os.Remove(dbName)

	err := createTestData()
	if err != nil {
		assert.FailNow(t, err.Error())
	}

	db, err := gorm.Open("sqlite3", dbName)
	if err != nil {
		assert.FailNow(t, err.Error())
	}
	defer db.Close()

	dao := ProjectService{DB: db}

	t.Run("fail with new", func(t *testing.T) {
		p := model.Project{
			Name:        "testproject124",
			Description: "asdf",
		}
		p.ID = 4321

		err := dao.Update(&p)

		assert.Error(t, err)
	})

	t.Run("test update existing", func(t *testing.T) {
		p := model.Project{
			Name:        " New Name ",
			Description: "Baz",
		}
		p.ID = uint(1)

		err := dao.Update(&p)

		assert.NoError(t, err)

		assert.NotZero(t, p.ID)
		assert.Equal(t, "newname", p.Name)
		assert.Equal(t, "Baz", p.Description)
		assert.NotNil(t, p.CreatedAt)
		assert.NotNil(t, p.UpdatedAt)
		assert.Nil(t, p.DeletedAt)
	})

	t.Run("test update existing no name", func(t *testing.T) {
		p := model.Project{
			Description: "Foo Test Bar",
		}
		p.ID = uint(1)

		err := dao.Update(&p)

		assert.NoError(t, err)

		assert.NotZero(t, p.ID)
		assert.Equal(t, "newname", p.Name)
		assert.Equal(t, "Foo Test Bar", p.Description)
		assert.NotNil(t, p.CreatedAt)
		assert.NotNil(t, p.UpdatedAt)
		assert.Nil(t, p.DeletedAt)
	})
}
