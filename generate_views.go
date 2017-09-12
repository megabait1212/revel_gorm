package main

import (
	"github.com/yosssi/gohtml"
	"os"
	"path"
	"strings"
)

// generate controller
func generateViews(mname, fields, crupath string) {
	// get controller name and package
	_, f := path.Split(mname)

	// set controller name to uppercase
	modelName := strings.Title(f)

	updateForm, err := GetUpdateFormAttributes(modelName, fields)
	if err != nil {
		ColorLog("[ERRO] Could not generate views: %s\n", err)
		os.Exit(2)
	}

	ColorLog("[INFO] Using '%s' is generated in views path\n", modelName)

	// create controller folders
	filePath := path.Join(crupath, "app", "views", modelName)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// create controller directory
		if err := os.MkdirAll(filePath, 0777); err != nil {
			ColorLog("[ERRO] Could not create views directory: %s\n", err)
			os.Exit(2)
		}
	}

	viewFiles := []string{"Index", "Edit"}
	for _, filename := range viewFiles {
		currentFile := path.Join(crupath, "app", "views", modelName, filename+".html")
		if _, err := os.Stat(currentFile); os.IsNotExist(err) {
			if cf, err := os.OpenFile(currentFile, os.O_CREATE|os.O_EXCL|os.O_RDWR, 0666); err == nil {
				defer cf.Close()
				switch filename {
				case "Index":
					content := strings.Replace(indexTpl, "{{pageTitle}}", modelName, -1)
					content = strings.Replace(content, "{{pageHeader}}", modelName+"s", -1)
					content = strings.Replace(content, "{{gridMethod}}", "Grid", -1)
					content = strings.Replace(content, "{{gridDeleteMethod}}", "GridDelete", -1)
					content = strings.Replace(content, "{{newMethod}}", "New", -1)
					content = strings.Replace(content, "{{modelName}}", modelName, -1)
					gohtml.FormatWithLineNo(content)
					cf.WriteString(content)
				case "Edit":
					content := strings.Replace(editTpl, "{{pageTitle}}", modelName, -1)
					content = strings.Replace(content, "{{pageHeader}}", "Edit "+modelName, -1)
					content = strings.Replace(content, "{{method}}", "POST", -1)
					content = strings.Replace(content, "{{action}}", strings.ToLower(modelName)+"s/update", -1)
					content = strings.Replace(content, "{{modelName}}", modelName, -1)
					content = strings.Replace(content, "{{buttonName}}", "Save "+strings.ToLower(modelName), -1)
					content = strings.Replace(content, "{{formAttributes}}", updateForm, -1)
					content = strings.Replace(content, "{{Action}}", "Update", -1)
					gohtml.FormatWithLineNo(content)
					cf.WriteString(content)
				}
				// gofmt generated source code
				// FormatSourceCode(currentFile)

				ColorLog("[INFO] '%s' view file are generated as: %s\n", filename, currentFile)
			} else {
				// error creating file
				ColorLog("[ERRO] Could not create views for '%s': %s\n", modelName, err)
				os.Exit(2)
			}
		}

	}
}

// remove existing views file
func deleteViews(mname, crupath string) {
	_, f := path.Split(mname)
	modelName := strings.Title(f)
	folderPath := path.Join(crupath, "app", "views", modelName)
	if _, err := os.Stat(folderPath); !os.IsNotExist(err) {
		viewFiles := []string{"Index", "Show", "New", "Edit"}
		for _, filename := range viewFiles {
			currentFile := path.Join(crupath, "app", "views", modelName, filename+".html")
			if _, err := os.Stat(currentFile); !os.IsNotExist(err) {
				err = os.Remove(currentFile)
				if err != nil {
					ColorLog("[ERRO] Could not delete %s: %s\n", filename, err)
				}
			}
		}
		err = os.Remove(folderPath)
		if err != nil {
			ColorLog("[ERRO] Could not delete views folder: %s\n", err)
			os.Exit(2)
		}
		ColorLog("[INFO] views files are deleted: %s\n", folderPath)
	}

}

var indexTpl = `
{{set . "title" "{{pageTitle}}"}}
{{template "header.html" .}}

<header class="hero-unit">
  <div class="container">
    <div class="row">
      <div class="hero-text">
        <h1>{{pageHeader}}</h1>
      </div>
    </div>
  </div>
</header>

<div class="container">
  <div class="row">
    <div class="col-md-12">
      {{template "flash.html" .}}
    </div>
    <div class="col-md-12">
        <div class="panel-header">
	  		<a href="{{url "{{modelName}}.{{newMethod}}"}}" class="btn btn-default">New</a>
	    </div>
    	<div class="panel panel-default">
		  <div class="panel-body">
		  	  <div id="grid"></div>
          </div>
		</div>
    </div>
  </div>
</div>
<script type="text/javascript">
    $(function () {
        var dataManager2 = ej.DataManager({
            url: "{{url "{{modelName}}.{{gridMethod}}"}}",
            crossDomain: true,
            removeUrl: "{{url "{{modelName}}.{{gridDeleteMethod}}"}}",
            adaptor: new ej.UrlAdaptor()
        });

        $("#grid").ejGrid({
            dataSource: dataManager2,
            allowPaging: true,
            allowSorting: true,
            allowFiltering: true,
            isResponsive: true,
            pageSettings: {pageSize: 10},
            filterSettings: {
                filterType: "menu"
            },
            toolbarSettings: {
                showToolbar: true,
                toolbarItems: [ej.Grid.ToolBarItems.Delete]
            },
            editSettings: {
                allowEditing: true,
                allowDeleting: true
            },
            columns: [
                {type: "checkbox", width: 50},
                {
                    headerText: "",
                    template: "<a href='{{url "{{modelName}}"}}/{^{:ID}}'>Edit</a>"
                },
            ]
        });
    });
</script>

{{template "footer.html" .}}
`
var editTpl = `
{{set . "title" "{{pageTitle}}"}}
{{template "header.html" .}}

<header class="hero-unit">
  <div class="container">
    <div class="row">
      <div class="hero-text">
        <h1>{{pageHeader}}</h1>
      </div>
    </div>
  </div>
</header>

<div class="container">
  <div class="row">
    <div class="col-md-12">
      {{template "flash.html" .}}
    </div>
    <div class="col-md-12">
    	<form method="POST" action="{{if .attributeset.ID}}{{url "{{modelName}}.Update"}}{{else}}{{url "{{modelName}}.Create"}}{{end}}"}}">
			{{formAttributes}}
    		<input type="submit" class="btn btn-success" value="{{buttonName}}" />
		</form>
    </div>
  </div>
</div>
{{template "footer.html" .}}`
