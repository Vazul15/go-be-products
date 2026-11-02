package endpoint

import (
	"errors"

	"github.com/gin-gonic/gin"

	"products/internal/logic/product"
	"products/internal/model/apierror"
	"github.com/Vazul15/go-lib-logger"
	"github.com/Vazul15/go-lib-restserver"
	"github.com/Vazul15/go-lib-restserver/model/resterror"
)

// respondError logs err and writes the corresponding HTTP error response.
func respondError(c *gin.Context, err error) {
	logger.FromContext(c.Request.Context()).Error("Exception", "err", err)

	apiErr := apierror.Classify(toAPIError(err))
	restserver.RespondError(c, resterror.NewError(apiErr.Status, apiErr.Code, apiErr.Message, err))
}

// respondBadRequest reports a malformed client request (body, query or path) as a 400.
func respondBadRequest(c *gin.Context, err error) {
	respondError(c, apierror.Generic(err.Error()))
}

func toAPIError(err error) error {
	switch {
	case errors.Is(err, product.ErrNotFound):
		return apierror.NotFound(err.Error())
	case errors.Is(err, product.ErrInvalidProduct):
		return apierror.Generic(err.Error())
	default:
		return err
	}
}
