package rest

import (
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"stitchdata.com/dbt/runner/dbt"
	"stitchdata.com/dbt/runner/rest/models"
)

var d *dbt.Dbt

func CreateProject(c echo.Context) error {
	return nil
}

func GetCompiledSQL(c echo.Context) error {
	project := c.QueryParam("project")
	path := c.QueryParam("path")
	name := c.QueryParam("name")
	accountid := c.Param("account_id")
	d = dbt.New(accountid)

	response := d.Execute("--project-dir", "/"+accountid+"/dbt_test", "--profiles-dir", "/"+accountid)
	if response.Status != "Success" {
		return c.String(http.StatusBadRequest, response.Message)
	}
	contents, err := ioutil.ReadFile("/" + accountid + "/" + project + "/target/compiled/my_new_project/models/" + strings.TrimRight(strings.TrimLeft(path, "/"), "/") + "/" + name)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	sql := string(contents)
	// r, err := json.Marshal(&models.ResponseCompiledSQL{
	// 	Name: name,
	// 	SQL:  sql,
	// })

	r := &models.ResponseCompiledSQL{
		Name: name,
		SQL:  sql,
	}
	// //rawIn := json.RawMessage(r)
	// if err != nil {
	// 	return c.String(http.StatusBadRequest, err.Error())
	// }

	// out, _ := strconv.Unquote(string(r))
	return c.JSON(http.StatusOK, r)
}
