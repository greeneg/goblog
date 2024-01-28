package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/greeneg/goblog/model"
	pbkdf2auth "github.com/greeneg/goblog/pbkdf2_auth"
)

type Config struct {
	apiAuthMech string
	tcpPort     string
	webAuthMech string
}

type BlogApp struct {
	appPath    string
	configPath string
	authHeader string
	confStruct Config
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func parseAuthMechStr(mechStr string) (string, string, error) {
	authType := strings.Split(mechStr, "|")
	return authType[0], authType[1], nil // returns the type of auth, and the "target" string
}

func (b *BlogApp) postBlog(c *gin.Context) {
	// first, this needs authed, so check that front-end send us a token we can use
	authTok := c.GetHeader(b.authHeader)
	if authTok != "" {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized action!"})
		return
	}
	// what type of auth should be used?
	authType, tokFile, err := parseAuthMechStr(b.confStruct.apiAuthMech)
	if authType == "file" {
		authorization := pbkdf2auth.ValidateViaFile(strings.Replace(tokFile, "%CONFIG_DIR%", b.configPath, -1),
			[]byte(authTok))
		if authorization != true {
			c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized action!"})
			return
		}
	}

	var json model.Blog
	if err := c.ShouldBindJSON(&json); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	s, err := model.AddBlog(json)
	if s {
		c.IndentedJSON(http.StatusOK, gin.H{"message": "Successfully injected new blog record"})
		return
	} else {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err})
	}
}

func (b *BlogApp) readBlog(c *gin.Context) {
	bEnts, err := model.GetBlogs()
	checkErr(err)

	if bEnts == nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "no records found!"})
		return
	} else {
		c.IndentedJSON(http.StatusOK, gin.H{"data": bEnts})
	}
}

func (b *BlogApp) readBlogById(c *gin.Context) {
	id := c.Param("id")
	bEnt, err := model.GetBlogById(id)
	checkErr(err)

	if bEnt.Author == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "no records found for " + id})
		return
	} else {
		c.IndentedJSON(http.StatusOK, gin.H{"data": bEnt})
	}
}

func (b *BlogApp) updateBlog(c *gin.Context) {
	// first, this needs authed, so check that front-end send us a token we can use
	authTok := c.GetHeader(b.authHeader)
	if authTok != "" {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized action!"})
		return
	}
	// what type of auth should be used?
	authType, tokFile, err := parseAuthMechStr(b.confStruct.apiAuthMech)
	if authType == "file" {
		var authorization bool = pbkdf2auth.ValidateViaFile(strings.Replace(tokFile, "%CONFIG_DIR%", b.configPath, -1),
			[]byte(authTok))
		if authorization != true {
			c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized action!"})
			return
		}
	}

	var json model.Blog
	if err := c.ShouldBindJSON(&json); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	bEntId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
	}

	s, err := model.UpdateBlog(json, bEntId)
	if s {
		c.IndentedJSON(http.StatusOK, gin.H{"message": "successfully updated " + c.Param("id")})
	} else {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err})
	}
}

func (b *BlogApp) deleteBlog(c *gin.Context) {
	// first, this needs authed, so check that front-end send us a token we can use
	authTok := c.GetHeader(b.authHeader)
	if authTok != "" {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized action!"})
		return
	}
	// what type of auth should be used?
	authType, tokFile, err := parseAuthMechStr(b.confStruct.apiAuthMech)
	if authType == "file" {
		var authorization bool = pbkdf2auth.ValidateViaFile(strings.Replace(tokFile, "%CONFIG_DIR%", b.configPath, -1),
			[]byte(authTok))
		if authorization != true {
			c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized action!"})
			return
		}
	}

	bEntId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid ID: " + c.Param("id")})
		return
	}

	s, err := model.DeleteBlog(bEntId)
	if s {
		c.IndentedJSON(http.StatusOK, gin.H{"message": "successfully deleted record " + c.Param("id")})
	} else {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err})
	}
}

func (b *BlogApp) options(c *gin.Context) {
	ourOptions := "HTTP/1.1 200 OK\n" +
		"Allow: GET,POST,PUT,DELETE,OPTIONS\n" +
		"Access-Control-Allow-Origin: http://locahost:8000\n" +
		"Access-Control-Allow-Methods: GET,POST,PUT,DELETE,OPTIONS\n" +
		"Access-Control-Allow-Headers: Content-Type, X-Auth-Token\n"

	c.String(http.StatusOK, ourOptions)
}

func main() {
	err := model.ConnectDatabase()
	checkErr(err)

	r := gin.Default()

	// lets get our working directory
	appdir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	checkErr(err)

	// config path is derived from app working directory
	configDir := filepath.Join(appdir, "config")

	// now that we have our appdir and configDir, lets read in our app config
	// and marshall it to the Config struct
	config := Config{}
	jsonContent, err := os.ReadFile(filepath.Join(configDir, "config.json"))
	checkErr(err)
	err = json.Unmarshal(jsonContent, &config)
	checkErr(err)

	// create an app object that contains our routes and the configuration
	Blog := new(BlogApp)
	Blog.appPath = appdir
	Blog.configPath = configDir
	Blog.authHeader = "X-Auth-Token"
	Blog.confStruct = config

	// This part is the web UI
	webview := r.Group("/")
	{
		webview.GET("/")
	}

	// This part is for the API
	router := r.Group("/api/v1")
	{
		router.POST("/", Blog.postBlog)
		router.GET("/", Blog.readBlog)
		router.GET("/:id", Blog.readBlogById)
		router.PUT("/:id", Blog.updateBlog)
		router.DELETE("/:id", Blog.deleteBlog)
		router.OPTIONS("/", Blog.options)
	}

	r.Run(Blog.confStruct.tcpPort)
}
