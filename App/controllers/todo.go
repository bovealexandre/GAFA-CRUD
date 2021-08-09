package controllers

import (
	"context"
	"fmt"
	"housecms/App/models"
	"housecms/config"
	"log"

	"github.com/arangodb/go-driver"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// var todos = []*models.Todo{
// 	{
// 		Id:        1,
// 		Title:     "Walk the dog 🦮",
// 		Completed: false,
// 	},
// 	{
// 		Id:        2,
// 		Title:     "Walk the cat 🐈",
// 		Completed: false,
// 	},
// }

func checkCollections() {
	ctx := context.Background()
	found, err := config.DB.CollectionExists(ctx, "todo")
	if err != nil {
		// handle error
		log.Fatalf("Failed to get database info: %v", err)
	}

	fmt.Printf("Collection found: %v\n", found)

	if !found {
		options := &driver.CreateCollectionOptions{ReplicationFactor: 1, WaitForSync: true}
		col, err := config.DB.CreateCollection(ctx, "todo", options)
		if err != nil {
			// handle error
			log.Fatalf("Failed to create collection: %v", err)
		}

		fmt.Printf("collection created: %v\n", col)
	}
}

// get all todos
func GetTodos(c *fiber.Ctx) error {
	ctx := context.Background()
	query := `FOR t IN todo
	RETURN t`
	checkCollections()
	var todos = []*models.Todo{}
	cursor, err := config.DB.Query(ctx, query, nil)
	if err != nil {
		log.Fatalf("collection error: %v", err)
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": false,
			"data": fiber.Map{
				"todos": todos,
			},
		})
	}

	defer cursor.Close()
	results := models.Todos{}

	for cursor.HasMore() {
		result := models.Todo{}

		if _, err := cursor.ReadDocument(ctx, &result); err != nil {
			return err
		}

		results.Todos = append(results.Todos, result)
	}

	if results.Todos == nil {
		results.Todos = []models.Todo{}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    results,
	})
}

// Create a todo
func CreateTodo(c *fiber.Ctx) error {
	var body models.Todo

	err := c.BodyParser(&body)

	if err != nil {
		return c.Status(fiber.StatusCreated).JSON(fiber.Map{})
	}

	if body.Title == "" {
		return c.Status(fiber.StatusCreated).JSON(fiber.Map{})
	}

	todo := models.Todo{
		Id:        uuid.New().String(),
		Title:     body.Title,
		Completed: false,
	}

	ctx := context.Background()

	col, err := config.DB.Collection(ctx, "todo")
	if err != nil {
		// handle error
	}

	meta, err := col.CreateDocument(nil, todo)
	if err != nil {
		log.Fatalf("Failed to create document: %v", err)
	}
	fmt.Printf("Created document in collection '%s' in database '%s'\n", col.Name(), config.DB.Name())

	// Read the document back
	var result models.Todo
	if _, err := col.ReadDocument(nil, meta.Key, &result); err != nil {
		log.Fatalf("Failed to read document: %v", err)
	}
	fmt.Printf("Read book '%+v'\n", result)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"todo": todo,
		},
	})
}

// get a single todo
// PARAM: id
func GetTodo(c *fiber.Ctx) error {
	// get parameter value
	paramId := c.Params("id")

	if c.Params("id") == "" {
		return c.Status(fiber.StatusCreated).JSON(fiber.Map{})
	}

	ctx := context.Background()
	col, err := config.DB.Collection(ctx, "todo")
	if err != nil {
		// handle error
	}

	var result models.Todo
	if _, err := col.ReadDocument(ctx, paramId, &result); err != nil {
		fmt.Printf("Failed to read document: %v", err)
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": false,
			"data": fiber.Map{
				"todo": result,
			},
		})
	}

	// find todo and return
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"todo": result,
		},
	})
}

// Update a todo
// PARAM: id
func UpdateTodo(c *fiber.Ctx) error {
	// find parameter
	paramId := c.Params("id")

	if paramId == "" {
		return c.Status(fiber.StatusCreated).JSON(fiber.Map{})
	}

	var body models.Todo

	c.BodyParser(&body)

	todo := models.Todo{
		Id:        paramId,
		Title:     body.Title,
		Completed: body.Completed,
	}

	fmt.Printf(todo.Title)
	fmt.Print(todo.Completed)

	ctx := context.Background()

	col, err := config.DB.Collection(ctx, "todo")
	if err != nil {
		// handle error
	}

	col.UpdateDocument(ctx, paramId, todo)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    fiber.Map{
			// "todo": meta,
		},
	})
}

// Delete a todo
// PARAM: id
func DeleteTodo(c *fiber.Ctx) error {
	paramId := c.Params("id")

	if c.Params("id") == "" {
		return c.Status(fiber.StatusCreated).JSON(fiber.Map{})
	}

	ctx := context.Background()
	col, err := config.DB.Collection(ctx, "todo")
	if err != nil {
		// handle error
	}

	if _, err := col.RemoveDocument(ctx, paramId); err != nil {
		fmt.Printf("Failed to read document: %v", err)
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": false,
		})
	}

	// find todo and return
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
	})
}
