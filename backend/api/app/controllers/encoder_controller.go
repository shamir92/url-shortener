package controllers

import (
	"strconv"

	"github.com/create-go-app/fiber-go-template/app/helpers"
	"github.com/create-go-app/fiber-go-template/app/models"
	"github.com/create-go-app/fiber-go-template/platform/database"
	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

var validate = validator.New()

func ValidateStruct(data models.ShortUrl) []*models.ErrorResponse {
	var errors []*models.ErrorResponse
	err := validate.Struct(data)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element models.ErrorResponse
			element.FailedField = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errors = append(errors, &element)
		}
	}
	return errors
}

func CreateShortUrl(c *fiber.Ctx) error {
	db := database.DBConn
	var err error
	tempData := new(models.ShortUrl)
	if err := c.BodyParser(tempData); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})

	}
	errors := ValidateStruct(*tempData)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors)

	}
	getCurrentCounter := helpers.GetCurrentCounter()
	if getCurrentCounter.Status == "error" {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": getCurrentCounter.Data,
		})
	}

	currentCounter, err := strconv.ParseUint(getCurrentCounter.Data, 10, 64)
	var x = helpers.ShortUrlEncoder(currentCounter)
	tx := models.ShortUrl{}
	tx.LongUrl = tempData.LongUrl
	tx.ShortUrl = x
	tx.Email = tempData.Email
	tx.ID, err = uuid.NewRandom()

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	if err := db.Create(&tx).Error; err != nil {
		// error handling...
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"error":    false,
		"msg":      nil,
		"shortUrl": x,
	})
}

func GetShortUrl(c *fiber.Ctx) error {
	db := database.DBConn
	shortUrl := c.Params("shortUrl")
	tx := models.ShortUrl{}
	if err := db.Where("short_url = ?", shortUrl).First(&tx).Error; err != nil {
		// error handling...
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"url": tx.LongUrl,
	})
}
