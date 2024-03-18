package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/dmitrii-a/image-previewer/internal/application"
	"github.com/dmitrii-a/image-previewer/internal/common"
	"github.com/gofiber/fiber/v2"
)

// ResizeImage is a handler for resizing images.
func ResizeImage(ctx *fiber.Ctx) error {
	width, err := strconv.Atoi(ctx.Params("width"))
	if err != nil {
		return ctx.Status(http.StatusBadRequest).SendString("Invalid width")
	}

	height, err := strconv.Atoi(ctx.Params("height"))
	if common.IsErr(err) {
		return ctx.Status(http.StatusBadRequest).SendString("Invalid height")
	}

	data, err := application.ImageApplicationService.ResizeImage(
		ctx.Params("*"),
		width,
		height,
		ctx.GetReqHeaders(),
	)

	var appErr *common.ApplicationError

	if errors.As(err, &appErr) {
		if appErr.Code == common.ImageDownloadErrorCode {
			return ctx.Status(http.StatusBadGateway).SendString(appErr.Error())
		}

		if appErr.Code == common.ImageDecodeErrorCode {
			return ctx.Status(http.StatusUnprocessableEntity).SendString(appErr.Error())
		}

		if appErr.Code == common.ImageNotFoundErrorCode {
			return ctx.Status(http.StatusNotFound).SendString(appErr.Error())
		}

		return ctx.Status(http.StatusInternalServerError).SendString(appErr.Error())
	} else if err != nil {
		return ctx.Status(http.StatusBadGateway).SendString(err.Error())
	}

	return ctx.Status(http.StatusOK).Send(data)
}

// HealthCheck is a handler for health checking.
func HealthCheck(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).SendString("ok")
}
