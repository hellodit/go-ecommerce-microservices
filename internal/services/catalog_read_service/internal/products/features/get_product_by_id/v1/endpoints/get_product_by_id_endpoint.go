package endpoints

import (
	"fmt"
	"net/http"

	"emperror.dev/errors"
	"github.com/labstack/echo/v4"
	customErrors "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/http/http_errors/custom_errors"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/web/route"
	"github.com/mehdihadeli/go-mediatr"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogreadservice/internal/products/contracts/params"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogreadservice/internal/products/features/get_product_by_id/v1/dtos"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogreadservice/internal/products/features/get_product_by_id/v1/queries"
)

type getProductByIdEndpoint struct {
	params.ProductRouteParams
}

func NewGetProductByIdEndpoint(
	params params.ProductRouteParams,
) route.Endpoint {
	return &getProductByIdEndpoint{
		ProductRouteParams: params,
	}
}

func (ep *getProductByIdEndpoint) MapEndpoint() {
	ep.ProductsGroup.GET("/:id", ep.handler())
}

// GetProductByID
// @Tags Products
// @Summary Get product
// @Description Get product by id
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} dtos.GetProductByIdResponseDto
// @Router /api/v1/products/{id} [get]
func (ep *getProductByIdEndpoint) handler() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		ep.CatalogsMetrics.GetProductByIdHttpRequests.Add(ctx, 1)

		request := &dtos.GetProductByIdRequestDto{}
		if err := c.Bind(request); err != nil {
			badRequestErr := customErrors.NewBadRequestErrorWrap(
				err,
				"[getProductByIdEndpoint_handler.Bind] error in the binding request",
			)
			ep.Logger.Errorf(
				fmt.Sprintf("[getProductByIdEndpoint_handler.Bind] err: %v", badRequestErr),
			)
			return badRequestErr
		}

		query := queries.NewGetProductById(request.Id)

		if err := ep.Validator.StructCtx(ctx, query); err != nil {
			validationErr := customErrors.NewValidationErrorWrap(
				err,
				"[getProductByIdEndpoint_handler.StructCtx]  query validation failed",
			)
			ep.Logger.Errorf("[getProductByIdEndpoint_handler.StructCtx] err: {%v}", validationErr)
			return validationErr
		}

		queryResult, err := mediatr.Send[*queries.GetProductById, *dtos.GetProductByIdResponseDto](
			ctx,
			query,
		)
		if err != nil {
			err = errors.WithMessage(
				err,
				"[getProductByIdEndpoint_handler.Send] error in sending GetProductById",
			)
			ep.Logger.Errorw(
				fmt.Sprintf(
					"[getProductByIdEndpoint_handler.Send] id: {%s}, err: {%v}",
					query.Id,
					err,
				),
				logger.Fields{"ProductId": query.Id},
			)
			return err
		}

		return c.JSON(http.StatusOK, queryResult)
	}
}
