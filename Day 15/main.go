package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"personal-webku/connect"
	"strconv"
	"text/template"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
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

type User struct {
	Name     string
	Email    string
	Password string
	Id       int
}

type SessionData struct {
	IsLogin bool
	Name    string
}

var DataProject = []AddProject{}
var userData = SessionData{}

// Func ketika server dimulai
func main() {

	//connect to database
	connect.DatabaseConnect()

	// Inisiasi echo di variabel e
	e := echo.New()

	//Inisiasi session
	e.Use(session.Middleware(sessions.NewCookieStore([]byte("session"))))

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

	e.GET("/formRegister", register) //localhost:5000/register
	e.GET("/formLogin", login)
	e.POST("/login", loginUser)
	e.POST("/register", registerUser)
	e.GET("/logout", logout)

	e.Logger.Fatal(e.Start("localhost:5000"))
}

// Func GET home
func home(c echo.Context) error {

	template, err := template.ParseFiles("views/index.html")

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"Message": err.Error()})
	}

	// if session.Values["isLogin"] != true {
	// 	userData.IsLogin = false
	// } else {
	// 	userData.IsLogin = session.Values["isLogin"].(bool)
	// 	userData.Name = session.Values["name"].(string)
	// }

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

	session, _ := session.Get("session", c)

	dataQuery := map[string]interface{}{
		"dataProject":  result,
		"dataSession":  userData,
		"FlashStatus":  session.Values["status"],
		"FlashMessage": session.Values["message"],
		"FlashName":    session.Values["name"],
	}

	// delete(session.Values, "status")
	delete(session.Values, "message")
	session.Save(c.Request(), c.Response())

	return template.Execute(c.Response(), dataQuery)
}

// Func GET contactMe
func contactMe(c echo.Context) error {
	template, err := template.ParseFiles("views/contact-me.html")

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"Message": err.Error()})
	}

	session, _ := session.Get("session", c)

	if session.Values["isLogin"] != true {
		userData.IsLogin = false
	} else {
		userData.IsLogin = session.Values["isLogin"].(bool)
		userData.Name = session.Values["name"].(string)
	}

	dataSession := map[string]interface{}{
		"dataSession": userData,
	}

	return template.Execute(c.Response(), dataSession)
}

// Func GET addProject
func addProject(c echo.Context) error {
	template, err := template.ParseFiles("views/add-project.html")

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"Message": err.Error()})
	}

	session, _ := session.Get("session", c)

	if session.Values["isLogin"] != true {
		userData.IsLogin = false
	} else {
		userData.IsLogin = session.Values["isLogin"].(bool)
		userData.Name = session.Values["name"].(string)
	}

	dataSession := map[string]interface{}{
		"dataSession": userData,
	}

	return template.Execute(c.Response(), dataSession)
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

	session, _ := session.Get("session", c)

	if session.Values["isLogin"] != true {
		userData.IsLogin = false
	} else {
		userData.IsLogin = session.Values["isLogin"].(bool)
		userData.Name = session.Values["name"].(string)
	}

	dataQuery := map[string]interface{}{
		"projectDetail": detailProject,
		"dataSession":   userData,
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

	session, _ := session.Get("session", c)

	if session.Values["isLogin"] != true {
		userData.IsLogin = false
	} else {
		userData.IsLogin = session.Values["isLogin"].(bool)
		userData.Name = session.Values["name"].(string)
	}

	dataQuery := map[string]interface{}{
		"projectDetail": editProject,
		"dataSession":   userData,
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

// Fungsi GET untuk menampilkan halaman register
func register(c echo.Context) error {
	template, err := template.ParseFiles("views/register.html")

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"Message": err.Error()})
	}

	return template.Execute(c.Response(), nil)
}

// Fungsi GET untuk menampilkan halaman login
func login(c echo.Context) error {
	session, _ := session.Get("session", c)

	messageFlash := map[string]interface{}{
		"FlashStatus":  session.Values["status"],
		"FlashMessage": session.Values["message"],
	}

	template, err := template.ParseFiles("views/login.html")

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"Message ": err.Error()})
	}

	delete(session.Values, "status")
	delete(session.Values, "message")
	session.Save(c.Request(), c.Response())

	return template.Execute(c.Response(), messageFlash)
}

// Fungsi POST form login
func loginUser(c echo.Context) error {
	err := c.Request().ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	email := c.FormValue("email")
	password := c.FormValue("password")

	user := User{}
	err = connect.Conn.QueryRow(context.Background(), "SELECT * FROM tb_user WHERE email=$1", email).Scan(&user.Name, &user.Email, &user.Password, &user.Id)
	if err != nil {
		fmt.Println("wow")
	}
	
	fmt.Println(err)
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return redirectWithMessage(c, "Password Salah !", false, "/formLogin")
	}
	session, _ := session.Get("session", c)
	session.Options.MaxAge = 3600
	session.Values["message"] = "Login Success"
	session.Values["status"] = true // show alert
	session.Values["name"] = user.Name
	session.Values["id"] = user.Id
	session.Values["isLogin"] = true // access login
	session.Save(c.Request(), c.Response())

	return redirectWithMessage(c, "Login success", true, "/")
	// return c.Redirect(http.StatusMovedPermanently, "/")
}

// Fungsi POST form register
func registerUser(c echo.Context) error {
	err := c.Request().ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	name := c.FormValue("name")
	email := c.FormValue("email")
	password := c.FormValue("password")

	fmt.Println(name)

	// generate password
	passwordHash, _ := bcrypt.GenerateFromPassword([]byte(password), 10)

	_, err = connect.Conn.Exec(context.Background(), "INSERT INTO tb_user (name, email, password) VALUES ($1, $2, $3)", name, email, passwordHash)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return redirectWithMessage(c, "Register success", true, "/formLogin")
}

func redirectWithMessage(c echo.Context, message string, status bool, path string) error {
	session, _ := session.Get("session", c)
	session.Values["message"] = message
	session.Values["status"] = status
	session.Save(c.Request(), c.Response())

	return c.Redirect(http.StatusSeeOther, path)
}

func logout(c echo.Context) error {
	session, _ := session.Get("session", c)
	session.Options.MaxAge = -1
	session.Values["isLogin"] = false
	session.Save(c.Request(), c.Response())

	return c.Redirect(http.StatusTemporaryRedirect, "/")
}
