package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	models "github.com/postgres/module"
	"github.com/postgres/storage"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

type Book struct {
	Author    string `json:"author"`
	Title     string `json:"title"`
	Publisher string `json:"publisher"`
}
type Repository struct {
	DB *gorm.DB
}

func (r *Repository) CreateBook(context *fiber.Ctx) error {

	book := Book{}

	err := context.BodyParser(&book)

	if err != nil {
		context.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{
				"message": "request failed",
			})
		return err
	}
	err = r.DB.Create(&book).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{
				"message": "NOT FOUNND",
			})
		return err
	}
	context.Status(http.StatusOK).JSON(&fiber.Map{
		"maeeage": "book",
	})
	return nil
}

func (r *Repository) DeletBook(context *fiber.Ctx) error {
	bookModels := &[]models.Books{}
	id := context.Params("id")
	if id == " " {
		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "id cannot be empty",
		})
		return nil
	}
	err := r.DB.Delete(bookModels, id)
	if err.Error != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "could not delete book",
		})
		return err.Error
	}
	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "books deleted successfully",
	})
	return nil

}

func (r *Repository) GetBooks(context *fiber.Ctx) error {
	bookModels := &[]models.Books{}

	err := r.DB.Find(bookModels).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not get books"})
		return err
	}
	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "books feteched successfully",
		"data":    bookModels,
	})
	return nil
}

func (r *Repository) GetBookById(context *fiber.Ctx) error {
	id := context.Params("id")
	bookModels := &models.Books{}
	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "id cannot be empty",
		})
		return nil
	}
	fmt.Println("the id is", id)
	err := r.DB.Where("id = ?", id).First(bookModels).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "could not get the book",
		})
		return err
	}
	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "book iss successfully fetched",
		"data":    bookModels,
	})
	return nil
}

// this makes this func a method
func (r *Repository) SetupRoutes(app *fiber.App) {
	api := app.Group("/api")
	// below gien statements are used to call methods for their specific needs
	api.Post("/create_books", r.CreateBook)
	api.Delete("delete_book/:id", r.DeletBook)
	api.Get("/get_books/id", r.GetBookById)
	api.Get("/books", r.GetBooks)
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)

	}

	config := &storage.Config{
		Host:     os.Getenv("DB_Host"),
		Port:     os.Getenv("DB_Port"),
		Password: os.Getenv("DB_Pass"),
		User:     os.Getenv("DB_User"),
		DBName:   os.Getenv("DB_Name"),
		SSLMode:  os.Getenv("DB_SSLMode"),
	}

	db, err := storage.NewConnection(config)
	if err != nil {
		log.Fatal("could not load the database")
	}
	err = models.MigrateBooks(db)
	if err != nil {
		log.Fatal("could nor migrate db")
	}

	r := Repository{
		DB: db,
	}

	app := fiber.New()
	r.SetupRoutes(app)
	app.Listen(":8080")
}
