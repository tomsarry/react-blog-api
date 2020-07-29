package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
)

// Article stores article information
type Article struct {
	ID      int64     `db:"id" json:"id" uri:"id"`
	Author  string    `db:"author" json:"author"`
	Content string    `db:"content" json:"content"`
	Title   string    `db:"title" json:"title"`
	Date    time.Time `db:"date" json:"date"`
}

// DB is the database
var DB *sqlx.DB

// Insert inserts a new entry in database
func (a Article) Insert() error {
	tx, err := DB.Beginx()

	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()
	_, err = tx.Exec("INSERT INTO articles (author, title, content) VALUES (?, ?, ?)", a.Author, a.Title, a.Content)

	return err
}

func createNewArticle() {
	newarticle := Article{
		Author:  "tom sarry",
		Title:   "An article from go",
		Content: "Ok guys this is epic, i just wrote an article from a go file.",
	}

	newarticle.Insert()
}

func retrieveArticle(id int) (Article, error) {

	var retrieved Article
	tx, err := DB.Beginx()
	if err != nil {
		fmt.Println(err)
		return retrieved, err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()
	tx.Get(&retrieved, "SELECT * FROM articles WHERE id = ?", id)
	return retrieved, err
}

func getArticles() ([]Article, error) {
	tx, err := DB.Beginx()
	if err != nil {
		return []Article{}, err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	var articles []Article
	err = DB.Select(&articles, "SELECT * from articles")
	return articles, err
}

func main() {
	godotenv.Load(".env")

	location := fmt.Sprintf("%v:%v@(%v)", os.Getenv("USER"), os.Getenv("DB_PSWD"), os.Getenv("HOST"))
	// var articles []Article
	var err error
	DB, err = sqlx.Connect("mysql", location+"/blog?parseTime=true")
	if err != nil {
		panic(err.Error())
	}
	defer DB.Close()

	if err != nil {
		panic(err.Error())
	}

	// createNewArticle()

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{os.Getenv("BLOG")},
		AllowMethods: []string{"GET", "PUT"},
	}))

	r.GET("/posts/:id", func(c *gin.Context) {
		//get id from get request
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.Status(500)
		}
		//make SQL selection based on ID
		article, err := retrieveArticle(id)
		c.JSON(200, article)
	})
	r.GET("/posts/", func(c *gin.Context) {
		//get all articles
		articles, err := getArticles()
		if err != nil {
			c.Status(500)
		}
		c.JSON(200, articles)
	})

	r.Run()

}
