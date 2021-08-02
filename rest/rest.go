package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"gopkg.in/yaml.v2"
	"stitchdata.com/dbt/runner/dbt"
	"stitchdata.com/dbt/runner/messages"
	"stitchdata.com/dbt/runner/rest/models"
)

var d *dbt.Dbt

func GetProfiles(c echo.Context) error {
	response := messages.Response{
		Status:  "Success",
		Message: "Hello World",
	}
	result, err := json.Marshal(response)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(c.Param("account_id"))
	fmt.Println(string(result))
	return c.JSON(http.StatusOK, string(result))
}

func PostCompile(c echo.Context) error {
	accountid := c.Param("account_id")
	d = dbt.New(accountid)
	//response := d.Compile("--project-dir", "/"+d.Path+"/dbt_test")
	response := d.Compile("--project-dir", "/"+accountid+"/dbt_test", "--profiles-dir", "/"+accountid)

	if response.Status == "Success" {
		return c.String(http.StatusOK, response.Message)
	} else {
		return c.String(http.StatusNotAcceptable, response.Message)
	}
}

func PostExecute(c echo.Context) error {
	accountid := c.Param("account_id")
	d = dbt.New(accountid)

	response := d.Execute("--project-dir", "/"+accountid+"/dbt_test", "--profiles-dir", "/"+accountid)

	if response.Status == "Success" {
		return c.String(http.StatusOK, response.Message)
	} else {
		return c.String(http.StatusNotAcceptable, response.Message)
	}
}

func CreateAccount(c echo.Context) error {
	var resp *messages.Response
	path := "/" + c.Param("account_id")
	fi, err := os.Lstat(path)
	if os.IsNotExist(err) {
		_ = os.Mkdir(path, fi.Mode().Perm())
		resp = &messages.Response{
			Status:  "Success",
			Message: "Directory create for Account " + c.Param("account_id"),
		}
	}
	
	if err == nil || os.IsExist(err) {
		resp = &messages.Response{
			Status:  "Success",
			Message: "Directory already exists for Account " + c.Param("account_id"),
		}
	}

	return c.JSON(http.StatusOK, resp)
}

func GetProjectVars(c echo.Context) error {
	accountid := c.Param("account_id")
	project := c.QueryParam("project")

	bytes, err := ioutil.ReadFile("/" + accountid + "/" + project + "/dbt_project.yml")
	if err != nil {
		c.JSON(http.StatusBadRequest, &messages.Response{
			Status:  "Failure",
			Message: err.Error(),
		})
	}
	var prj *models.DBTProject
	err = yaml.Unmarshal(bytes, &prj)
	if err != nil {
		c.JSON(http.StatusBadRequest, &messages.Response{
			Status:  "Failure",
			Message: err.Error(),
		})
	}
	var vars []*messages.Var
	for k, v := range prj.Vars {
		vars = append(vars, &messages.Var{
			Name:  k,
			Value: v,
			Type:  getDataType(v),
		})

	}
	return c.JSON(http.StatusOK, vars)

}

func CreateProject(c echo.Context) error {
	path := "/" + c.Param("account_id")
	reader := c.Request().Body
	var resp *messages.Response
	fi, err := os.Lstat(path)
	if os.IsNotExist(err) {
		resp = &messages.Response{
			Status:  "Failure",
			Message: "Directory for Account " + c.Param("account_id") + " not created. Please call POST /account/" + c.Param("account_id") + " to properly create directory",
		}
		return c.JSON(http.StatusBadRequest, resp)
	}
	if err != nil {
		return c.JSON(http.StatusBadRequest, &messages.Response{
			Status: "Failure",
			Message: err.Error(),
		})
	}
	if fi.IsDir() {
		buf := new(bytes.Buffer)
		defer reader.Close()
		buf.ReadFrom(reader)
		var clone *messages.Clone
		err = json.Unmarshal(buf.Bytes(), &clone)
		if err != nil {
			return c.JSON(http.StatusBadRequest, &messages.Response{
				Status: "Failure",
				Message: err.Error(),
			})
		}
		args := []string{"clone",clone.URL,path+"/"+clone.Name}
		cmd := exec.Command("git", args...)
		var out bytes.Buffer
		cmd.Stdout = &out
		err := cmd.Run()
		if err != nil {
			return c.JSON(http.StatusBadRequest,&messages.Response{
				Status:  "Error",
				Message: err.Error(),
			})
		}
	}

	return c.JSON(http.StatusOK, &messages.Response{
		Status: "Success",
		Message: "Project Cloned",
	})
}

func GetCompiledSQL(c echo.Context) error {
	project := c.QueryParam("project")
	path := c.QueryParam("path")
	name := c.QueryParam("name")
	accountid := c.Param("account_id")
	d = dbt.New(accountid)

	contents, err := ioutil.ReadFile("/" + accountid + "/" + project + "/target/compiled/my_new_project/models/" + strings.TrimRight(strings.TrimLeft(path, "/"), "/") + "/" + name)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error()+"\n The project was not compiled before hand.")
	}
	sql := string(contents)

	r := &models.ResponseCompiledSQL{
		Name: name,
		SQL:  sql,
	}
	return c.JSON(http.StatusOK, r)
}

func getDataType(val string) string {
	v := strings.Trim(val, " ")
	if v == "" {
		return "string"
	}
	_, err := strconv.ParseBool(v)
	if err == nil {
		return "bool"
	}
	_, err = strconv.ParseFloat(v, 32)
	if err == nil {
		return "float"
	}
	_, err = strconv.Atoi(v)
	if err == nil {
		return "int"
	}
	m, _ := regexp.Match("^\\d{4}\\-\\d{2}\\-\\d{2}$", []byte(v))
	if m {
		return "date"
	}
	return "string"
}
