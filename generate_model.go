package main

import (
	"os"
	"path"
	"strings"
)

// generate model file
func generateModel(mname, fields, crupath string) {
	// get name and package
	p, f := path.Split(mname)

	// Title to model name
	modelName := strings.Title(f)

	// set default package
	packageName := "models"
	if p != "" {
		i := strings.LastIndex(p[:len(p)-1], "/")
		packageName = p[i+1 : len(p)-1]
	}

	// get Struct from fileds
	modelStruct, timePkg, err := GetStruct(modelName, fields)
	if err != nil {
		ColorLog("[ERRO] Could not generate models struct: %s\n", err)
		os.Exit(2)
	}

	ColorLog("[INFO] Using '%s' as model name\n", modelName)
	ColorLog("[INFO] Using '%s' as package name\n", packageName)

	// create models folder if not exist
	filePath := path.Join(crupath, "app", "models", p)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// create controller directory
		if err := os.MkdirAll(filePath, 0777); err != nil {
			ColorLog("[ERRO] Could not create models directory: %s\n", err)
			os.Exit(2)
		}
	}

	// create model file with template
	fpath := path.Join(filePath, strings.ToLower(modelName)+".go")
	if f, err := os.OpenFile(fpath, os.O_CREATE|os.O_EXCL|os.O_RDWR, 0666); err == nil {
		defer f.Close()

		// get current path
		paths := strings.Split(crupath, "/")

		// get app name
		projectName := paths[len(paths)-1:][0]

		// get mongodb pcakge path
		databasePkg := path.Join(projectName, "app", "models", "database")

		content := strings.Replace(modelTpl, "{{packageName}}", "models", -1)
		if timePkg == true {
			content = strings.Replace(content, "{{timePkg}}", `"time"`, -1)
		} else {
			content = strings.Replace(content, "{{timePkg}}", "", -1)
		}
		content = strings.Replace(content, "{{databasePkg}}", databasePkg, -1)
		content = strings.Replace(content, "{{modelStruct}}", modelStruct, -1)
		content = strings.Replace(content, "{{modelStructName}}", modelName, -1)
		content = strings.Replace(content, "{{modelObjectName}}", strings.ToLower(modelName), -1)
		f.WriteString(content)
		// gofmt generated source code
		FormatSourceCode(fpath)
		ColorLog("[INFO] model file generated: %s\n", fpath)
	} else {
		// error creating file
		ColorLog("[ERRO] Could not create model file: %s\n", err)
		os.Exit(2)
	}

	// create grid grid.go
	gridModelFp := path.Join(crupath, "app", "models", "grid.go")
	if _, err := os.Stat(gridModelFp); os.IsNotExist(err) {
		if cf, err := os.OpenFile(gridModelFp, os.O_CREATE|os.O_EXCL|os.O_RDWR, 0666); err == nil {
			defer cf.Close()
			content := strings.Replace(gridModelTpl, "{{packageName}}", packageName, -1)
			cf.WriteString(content)
			// gofmt generated source code
			FormatSourceCode(gridModelFp)
			ColorLog("[INFO] grid model file generated: %s\n", gridModelFp)
		} else {
			// error creating file
			ColorLog("[ERRO] Could not create grid model file: %s\n", err)
			os.Exit(2)
		}
	}
}

// remove existing model file
func deleteModel(mname, crupath string) {
	_, f := path.Split(mname)
	modelName := strings.Title(f)
	filePath := path.Join(crupath, "app", "models", modelName+".go")
	if _, err := os.Stat(filePath); !os.IsNotExist(err) {
		err = os.Remove(filePath)
		if err != nil {
			ColorLog("[ERRO] Could not delete model struct: %s\n", err)
			os.Exit(2)
		}
		ColorLog("[INFO] model file deleted: %s\n", filePath)

	}

}

var modelTpl = `package {{packageName}}

import (
	"{{databasePkg}}"
	"github.com/jinzhu/gorm"
	"errors"
	{{timePkg}}
	"github.com/revel/revel"
)

{{modelStruct}}


func (m *{{modelStructName}}) BeforeCreate(scope *gorm.Scope) error {
 // Do something before create
  return nil
}

func (m *{{modelStructName}}) BeforeUpdate(scope *gorm.Scope) error {
 // Do something before create
  return nil
}

func (m *{{modelStructName}}) AfterUpdate(scope *gorm.Scope) error {
 // Do something before create
  return nil
}

// Add{{modelStructName}} insert a new {{modelStructName}} into database and returns
// last inserted {{modelObjectName}} on success.
func (m {{modelStructName}}) Add{{modelStructName}}() ({{modelStructName}}, error) {
	var err error
	if !database.DB.NewRecord(m) {
		return m, errors.New("primary key should be blank")
	} 

	tx := database.DB.Begin()
	if err = tx.Create(&m).Error; err != nil {
     	tx.Rollback()
     	return m, err
  	}
	tx.Commit()

	return m, err
}

// Update{{modelStructName}} update a {{modelStructName}} into database and returns
// last nil on success.
func (m {{modelStructName}}) Update{{modelStructName}}() ({{modelStructName}}, error){
	var err error

	if database.DB.NewRecord(m) {
		return m, errors.New("primary key should not be blank")
	} 

	tx := database.DB.Begin()
	if err = tx.Save(&m).Error; err != nil {
     	tx.Rollback()
     	return m, err
  	}
	tx.Commit()
	return m, err
}

// Delete{{modelStructName}} Delete {{modelStructName}} from database and returns
// last nil on success.
func (m {{modelStructName}}) Delete{{modelStructName}}() error{
	var err error
	tx := database.DB.Begin()
	if err = tx.Delete(&m).Error; err != nil {
     	tx.Rollback()
     	return err
  	}
  	tx.Commit()
  	return err
}

// Get{{modelStructName}}s Get all {{modelStructName}} from database and returns
// list of {{modelStructName}} on success
func (m {{modelStructName}}) Get{{modelStructName}}s() ([]{{modelStructName}},error) {
	var (
		err error
		{{modelObjectName}}s []{{modelStructName}}
	)

	tx := database.DB.Begin()
	if err = tx.Find(&{{modelObjectName}}s).Error; err != nil {
     	tx.Rollback()
     	return {{modelObjectName}}s, err
  	}
	tx.Commit()
	return {{modelObjectName}}s, err
}

// Get{{modelStructName}} Get a {{modelStructName}} from database and returns
// a {{modelStructName}} on success
func (m {{modelStructName}}) Get{{modelStructName}}(id uint64) ({{modelStructName}}, error){
	var (
		{{modelObjectName}} {{modelStructName}}
		err error
	)
	tx := database.DB.Begin()
	if err = tx.Last(&{{modelObjectName}}, id).Error; err != nil {
     	tx.Rollback()
     	return {{modelObjectName}}, err
  	}
	tx.Commit()
	return {{modelObjectName}}, err
}

func Migrate{{modelStructName}}(){
	database.DB.AutoMigrate(&{{modelStructName}}{})
}

func ({{modelStructName}}) TableName() string {
  return "{{modelObjectName}}s"
}

func ({{modelObjectName}} *{{modelStructName}}) Validate(v *revel.Validation) {
	//Validation rules here
}

`
var gridModelTpl = `
package {{packageName}}

import (
	"fmt"
	"regexp"
	"strings"
)

func prepareOrder(o []string) string {
	return strings.Join(o, ", ")
}

func prepareWhere(w []map[string]string) (string, []interface{}) {
	var where string
	var args []interface{}
	for i, v := range w {
		col := ToSnakeCase(v["column"])
		val := v["value"]
		fmt.Println(v["operator"])

		switch v["operator"] {
		case "equal":
			where += col + " = ?"
		case "notequal":
			where += col + " <> ?"
		case "greaterthan":
			where += col + " > ?"
		case "lessthan":
			where += col + " < ?"
		case "greaterthanorequal":
			where += col + " >= ?"
		case "lessthanorequal":
			where += col + " <= ?"
		case "startswith":
			where += col + " LIKE ?"
			val = val + "%"
		case "endswith":
			where += col + " LIKE ?"
			val = "%" + val
		case "contains":
			where += col + " LIKE ?"
			val = "%" + val + "%"
		}

		if len(w)-1 > i {
			where += " AND "
		}
		args = append(args, val)
	}
	fmt.Println(where)
	return where, args
}

func ToSnakeCase(str string) string {
	var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
	var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}
`
