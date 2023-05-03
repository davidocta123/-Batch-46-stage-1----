package main

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"text/template"
	"webku/connect"

	"github.com/labstack/echo/v4"
)

type AddProject struct {
	Image      string
	Diff       int
	Id		   int
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
}

var DataProject = []AddProject{
	//{
	//	Title:      "Dumbways Versi 1",
	//	Content:    "Selama Kursus saya dah mulai ngerti Golang",
	//	Author:     "David Octavyanto",
	//	DateMonth:  "Durasi 3 Bulan",
	//	StartDate:  "01 Februari 2023",
	//	EndDate:    "01 Juni 2023",
	//	DateDay:    "3 Bulan",
	//	TechJS:     "JavaScript",
	//	TechGolang: "Golang",
	//	TechGithub: "Github",
	//},
	//{
	//	Title:      "Dumbways Versi 2",
	//	Content:    "Demi Apapun Harus percaya ama dumbaways karena ada penyaluran kerja",
	//	Author:     "David Octavyanto",
	//	DateMonth:  "Durasi 4 Bulan",
	//	StartDate:  "01 Februari 2023",
	//	EndDate:    "01 Juli 2023",
	//	DateDay:    "4 Bulan",
	//	TechGolang: "Golang",
	//	TechGithub: "Github",
	//	TechNodeJs: "NodeJs",
	//},
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
	//e.POST("/addProject", formAddProject)       //localhost:5000/formAddProject
	e.GET("/editProject/:id", editProject)      // localhost:5000/editProject/:id
	//e.POST("/editProject/:id", editProjectDone) //localhost:5000/editProject/:id
	e.GET("/deleteProject/:id", deleteProject)  // localhost:5000/deleteProject/:id

	e.Logger.Fatal(e.Start("localhost:5000"))
}

// Func GET home
func home(c echo.Context) error {
	template, err := template.ParseFiles("views/index.html")

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"Message": err.Error()})
	}

	data, _ := connect.Conn.Query(context.Background(), "SELECT Id, Title, (End_Date - Start_Date) / 30 as Diff, Content, Author, Image FROM tb_projectweb;")

	var result []AddProject

	for data.Next() {

		var each = AddProject{}
		err := data.Scan(&each.Id, &each.Title, &each.Diff, &each.Content, &each.Author, &each.Image)
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
	template, err := template.ParseFiles("views/project-detail.html")

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"Message": err.Error()})
	}

	id, _ := strconv.Atoi(c.Param("id"))

	var detailProject = AddProject{}

	for index, data := range DataProject {
		if id == index {
			detailProject = AddProject{
				Title:      data.Title,
				Content:    data.Content,
				StartDate:  data.StartDate,
				EndDate:    data.EndDate,
				DateDay:    data.DateDay,
				TechJS:     data.TechJS,
				TechGolang: data.TechGolang,
				TechGithub: data.TechGithub,
				TechNodeJs: data.TechNodeJs,
			}
		}
	}

	data := map[string]interface{}{
		"projectDetail": detailProject,
	}

	return template.Execute(c.Response(), data)
}

// Func POST addProject
//func formAddProject(c echo.Context) error {

	// Untuk menghitung waktu
	// var Date1, Date2 string
	//var diffMonth, diffDay int

	//Date1 = c.FormValue("startDate")
	// Date2 = c.FormValue("endDate")

	//Merubah string date ke time
	// date1, _ := time.Parse("2006-01-02", Date1)
	//date2, _ := time.Parse("2006-01-02", Date2)

	//Format RFC822
	//date1NewForm := date1.Format(time.RFC822)
	//date2NewForm := date2.Format(time.RFC822)

	//untuk Bulan
	//diffMonth = int(date2.Sub(date1).Hours() / 24 / 30)
	//Ubah integer ke string untuk Bulan
	//diffValueMonth := strconv.Itoa(diffMonth)

	//untuk Day
	//diffDay = int(date2.Sub(date1).Hours() / 24)
	//Ubah integer ke string untuk Day
	//diffValueDay := strconv.Itoa(diffDay)

	//var newProject = AddProject{
		//Title:      c.FormValue("titleProject"),
		//Content:    c.FormValue("contentProject"),
		//Author:     "David Octavyanto",
		//DateMonth:  ("Durasi " + diffValueMonth + " Bulan"),
		//DateDay:    (diffValueDay + " Hari"),
		//StartDate:  date1NewForm,
		//EndDate:    date2NewForm,
		//TechJS:     c.FormValue("JavaScript"),
		//TechGolang: c.FormValue("Golang"),
		//TechGithub: c.FormValue("Github"),
		//TechNodeJs: c.FormValue("NodeJs"),
	//}

	//DataProject = append(DataProject, newProject)

	//return c.Redirect(http.StatusMovedPermanently, "/")
//}

// Func GET editProject
func editProject(c echo.Context) error {

	template, err := template.ParseFiles("views/edit-project.html")

	id, _ := strconv.Atoi(c.Param("id"))

	var editProject = AddProject{}

	for index, data := range DataProject {
		if id == index {
			editProject = AddProject{
				Title:   data.Title,
				Content: data.Content,
			}
		}
	}

	data := map[string]interface{}{
		"projectDetail": editProject,
	}

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"Message": err.Error()})
	}

	return template.Execute(c.Response(), data)
}

// Func POST editProject
//func editProjectDone(c echo.Context) error {

	// Untuk menghitung waktu
	//var Date1, Date2 string
	//var diffMonth, diffDay int

	//Date1 = c.FormValue("startDate")
	//Date2 = c.FormValue("endDate")

	//Merubah string date ke time
	//date1, _ := time.Parse("2006-01-02", Date1)
	//date2, _ := time.Parse("2006-01-02", Date2)

	//Format RFC822
	//date1NewForm := date1.Format(time.RFC822)
	//date2NewForm := date2.Format(time.RFC822)

	//untuk Bulan
	//diffMonth = int(date2.Sub(date1).Hours() / 24 / 30)
	//Ubah integer ke string untuk Bulan
	//diffValueMonth := strconv.Itoa(diffMonth)

	//untuk Day
	//diffDay = int(date2.Sub(date1).Hours() / 24)
	//Ubah integer ke string untuk Day
	//diffValueDay := strconv.Itoa(diffDay)

	//id, _ := strconv.Atoi(c.Param("id"))

	//var editProject = AddProject{
		//Title:      c.FormValue("titleProject"),
		//Content:    c.FormValue("contentProject"),
		//Author:     "David Octavyanto",
		//DateMonth:  ("Durasi " + diffValueMonth + " Bulan"),
		//DateDay:    (diffValueDay + " Hari"),
		//StartDate:  date1NewForm,
		//EndDate:    date2NewForm,
		//TechJS:     c.FormValue("JavaScript"),
		//TechGolang: c.FormValue("Golang"),
		//TechGithub: c.FormValue("Github"),
		//TechNodeJs: c.FormValue("NodeJs"),
	//}

	//DataProject = append(DataProject[:id+1], editProject)

	//return c.Redirect(http.StatusMovedPermanently, "/")
//}

// Func GET deleteProject
func deleteProject(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	DataProject = append(DataProject[:id], DataProject[id+1:]...)

	return c.Redirect(http.StatusMovedPermanently, "/")
}
