package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/malijoe/receipt-processor/application"
	"github.com/malijoe/receipt-processor/models"
	statuserrors "github.com/malijoe/receipt-processor/statusErrors"
)

func main() {
	router := gin.Default()

	app := application.NewApplication()
	// handler for POST /receipts/process endpoint
	router.POST("/receipts/process", func(ctx *gin.Context) {
		var receipt models.Receipt
		if err := ctx.ShouldBindBodyWithJSON(&receipt); err != nil {
			handleAppError(ctx, fmt.Errorf("%w: %s", statuserrors.ErrBadRequest, "The receipt is invalid."))
		}

		id, err := app.ProcessReceipt(ctx, receipt)
		if err != nil {
			handleAppError(ctx, err)
			return
		}
		ctx.JSON(http.StatusOK, map[string]any{"id": id})
	})
	// handler for GET /receipts/{id}/points
	router.GET("/receipts/:id/points", func(ctx *gin.Context) {
		id := ctx.Param("id")
		points, err := app.GetReceiptPoints(ctx, id)
		if err != nil {
			handleAppError(ctx, err)
			return
		}
		ctx.JSON(http.StatusOK, map[string]any{"points": points})
	})
	router.Run(":8080")
}

func handleAppError(ctx *gin.Context, err error) {
	switch e := err.(type) {
	case statuserrors.StatusError:
		ctx.JSON(e.Status(), e.Error())
	default:
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

}
