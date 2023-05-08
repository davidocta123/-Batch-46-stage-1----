package main

import (
	"context"
	"fmt"
	"net/http"
	"webku/connect"
	"strconv"
	"text/template"

	"github.com/labstack/echo/v4"
)

type AddProject struct {
	Id         int
	Title      string
	Content    string
	Author     string
	DateMonth  string
	DateDay    string
	StartDate  string
	EndDate    string
	TechJS     string
	TechGolang string
	TechGithub string
	TechNodeJs string
	Image      string
	Diff       int
}

var DataProject = []AddProject{
	// {
	// 	Title:      "Dumbways Versi 1",
	// 	Content:    "Mampus Gw, Bingung begini",
	// 	Author:     "Yoga Jaim",
	// 	DateMonth:  "Durasi 1 Bulan",
	// 	StartDate:  "01 Maret 2023",
	// 	EndDate:    "01 April 2023",
	// 	DateDay:    "30 Hari",
	// 	TechJS:     "JavaScript",
	// 	TechGolang: "Golang",
	// 	TechGithub: "Github",
	// },
	// {
	// 	Title:      "Dumbways Versi 2",
	// 	Content:    "Mampus, Apa yang dipelajari ya",
	// 	Author:     "Yoga Gila",
	// 	DateMonth:  "Durasi 1 Bulan",
	// 	StartDate:  "01 Maret 2023",
	// 	EndDate:    "01 April 2023",
	// 	DateDay:    "30 Hari",
	// 	TechGolang: "Golang",
	// 	TechGithub: "Github",
	// 	TechNodeJs: "NodeJs",
	// },
}

// Func ketika server dimulai
func main() {

	//connect to database
	connect.DatabaseConnect()

	// Inisiasi echo di variabel e
	e := echo.New()

	//Static for Access Folder
	e.Static("assets/", "assets")

	//Routing
	e.GET("/", home)                            // localhost:5000/
	e.GET("/contactMe", contactMe)              // localhost:5000/contactMe
	e.GET("/addProject", addProject)            // localhost:5000/addProject
	e.GET("/projectDetail/:id", projectDetail)  // localhost:5000/projectDetail
	e.POST("/addProject", formAddProject)       //localhost:5000/formAddProject
	e.GET("/editProject/:id", editProject)      // localhost:5000/editProject/:id
	e.POST("/editProject/:id", editProjectDone) //localhost:5000/editProject/:id
	e.GET("/deleteProject/:id", deleteProject)  // localhost:5000/deleteProject/:id

	e.Logger.Fatal(e.Start("localhost:5000"))
}

// Func GET home
func home(c echo.Context) error {
	template, err := template.ParseFiles("views/index.html")

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"Message": err.Error()})
	}

	data, _ := connect.Conn.Query(context.Background(), "SELECT Id, Title, (End_Date - Start_Date) / 30 as Diff, Content, Author, Techno[1], Techno[2], Techno[3], Techno[4], Image FROM tb_projectweb;")

	var result []AddProject

	for data.Next() {

		var each = AddProject{}
		err := data.Scan(&each.Id, &each.Title, &each.Diff, &each.Content, &each.Author, &each.TechJS, &each.TechGolang, &each.TechGithub, &each.TechNodeJs, &each.Image)
		if err != nil {
			fmt.Println(err.Error())
			return c.JSON(http.StatusInternalServerError, map[string]string{"Message ": err.Error()})
		}

		result = append(result, each)
	}

	dataQuery := map[string]interface{}{
		"dataProject": result,
	}

	return template.Execute(c.Response(), dataQuery)
}

// Func GET contactMe
func contactMe(c echo.Context) error {
	template, err := template.ParseFiles("views/contact-me.html")

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"Message": err.Error()})
	}

	return template.Execute(c.Response(), nil)
}

// Func GET addProject
func addProject(c echo.Context) error {
	template, err := template.ParseFiles("views/add-project.html")

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"Message": err.Error()})
	}

	return template.Execute(c.Response(), nil)
}

// Func menampilkan detailProject
func projectDetail(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	template, err := template.ParseFiles("views/project-detail.html")

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"Message": err.Error()})
	}

	var detailProject = AddProject{}

	dataErr := connect.Conn.QueryRow(context.Background(), "SELECT Id, Title, End_Date-Start_Date as Diff, TO_CHAR(Start_Date, 'DD-Mon-YYYY') Start_Date, TO_CHAR(End_Date, 'DD-Mon-YYYY') End_Date, Content, Image, Author, Techno[1], Techno[2], Techno[3], Techno[4] FROM tb_projectweb WHERE Id = $1;", id).Scan(&detailProject.Id, &detailProject.Title, &detailProject.Diff, &detailProject.StartDate, &detailProject.EndDate, &detailProject.Content, &detailProject.Image, &detailProject.Author, &detailProject.TechJS, &detailProject.TechGolang, &detailProject.TechGithub, &detailProject.TechNodeJs)

	if dataErr != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"Message ": dataErr.Error()})
	}

	dataQuery := map[string]interface{}{
		"projectDetail": detailProject,
	}

	return template.Execute(c.Response(), dataQuery)
}

// Func POST addProject
func formAddProject(c echo.Context) error {

	Title := c.FormValue("titleProject")
	Content := c.FormValue("contentProject")
	Author := "David Octavyanto"
	StartDate := c.FormValue("startDate")
	EndDate := c.FormValue("endDate")
	TechJS := c.FormValue("JavaScript")
	TechGolang := c.FormValue("Golang")
	TechGithub := c.FormValue("Github")
	TechNodeJs := c.FormValue("NodeJs")
	Image := "image.png"

	_, err := connect.Conn.Exec(context.Background(), "INSERT INTO tb_projectweb (Title, Content, Author, Start_Date, End_Date, Techno[1], Techno[2], Techno[3], Techno[4], Image) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)", Title, Content, Author, StartDate, EndDate, TechJS, TechGolang, TechGithub, TechNodeJs, Image)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"Message ": err.Error()})
	}

	return c.Redirect(http.StatusMovedPermanently, "/")
}

// Func GET editProject
func editProject(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	template, err := template.ParseFiles("views/edit-project.html")

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"Message": err.Error()})
	}

	var editProject = AddProject{}

	dataErr := connect.Conn.QueryRow(context.Background(), "SELECT Id, Title, Content FROM tb_projectweb WHERE Id = $1;", id).Scan(&editProject.Id, &editProject.Title, &editProject.Content)

	if dataErr != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"Message ": dataErr.Error()})
	}

	dataQuery := map[string]interface{}{
		"projectDetail": editProject,
	}

	return template.Execute(c.Response(), dataQuery)
}

// Func POST editProject
func editProjectDone(c echo.Context) error {

	Id, _ := strconv.Atoi(c.Param("id"))

	Title := c.FormValue("titleProject")
	Content := c.FormValue("contentProject")
	Author := "David Octavyanto"
	StartDate := c.FormValue("startDate")
	EndDate := c.FormValue("endDate")
	TechJS := c.FormValue("JavaScript")
	TechGolang := c.FormValue("Golang")
	TechGithub := c.FormValue("Github")
	TechNodeJs := c.FormValue("NodeJs")
	Image := "image.png"

	_, err := connect.Conn.Exec(context.Background(), "UPDATE tb_projectweb SET Title=$2, Content=$3, Author=$4, Start_Date=$5, End_Date=$6, Techno[1]=$7, Techno[2]=$8, Techno[3]=$9, Techno[4]=$10, Image=$11 WHERE Id = $1", Id, Title, Content, Author, StartDate, EndDate, TechJS, TechGolang, TechGithub, TechNodeJs, Image)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"Message ": err.Error()})
	}

	return c.Redirect(http.StatusMovedPermanently, "/")
}

// Func GET deleteProject
func deleteProject(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	_, err := connect.Conn.Exec(context.Background(), "DELETE FROM tb_projectweb WHERE Id=$1", id)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"Message ": err.Error()})
	}

	return c.Redirect(http.StatusMovedPermanently, "/")
}