package handlers

import (
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/noredis/subscriptions/internal/application/appservice"
	"github.com/noredis/subscriptions/internal/application/dto"
	"github.com/noredis/subscriptions/pkg/httpext"
	"github.com/rs/zerolog"
)

type CostHandler struct {
	logger  *zerolog.Logger
	service *appservice.CostService
}

func NewCostHandler(
	logger *zerolog.Logger,
	service *appservice.CostService,
) *CostHandler {
	return &CostHandler{
		logger:  logger,
		service: service,
	}
}

func (handler *CostHandler) Register(app *fiber.App) {
	app.Get("/costs/total", handler.Total)
}

// Total возвращает суммарную стоимость подписок.
//
// @Summary      Получить суммарную стоимость подписок
// @Description  Возвращает общую стоимость подписок с учётом фильтров.
// @Tags         cost
// @Produce      json
// @Param        service_name  query     string  false  "Фильтр по имени сервиса"
// @Param        user_id       query     string  false  "Фильтр по ID пользователя"
// @Param        start_date    query     string  false  "Дата начала (MM-YYYY)"
// @Param        end_date      query     string  false  "Дата окончания (MM-YYYY)"
// @Success      200  {object}  dto.TotalCostResponse  "Суммарная стоимость"
// @Failure      400  {object}  httpext.FiberError     "Некорректный запрос"
// @Failure      422  {object}  httpext.FiberError     "Ошибка валидации"
// @Failure      500  {object}  httpext.FiberError     "Внутренняя ошибка сервера"
// @Router       /cost/total [get]
func (handler *CostHandler) Total(c *fiber.Ctx) error {
	filters := dto.CostFilterDTO{
		ServiceName: c.Query("service_name"),
		UserID:      c.Query("user_id"),
		StartDate:   c.Query("start_date"),
		EndDate:     c.Query("end_date"),
	}

	cost, err := handler.service.Total(c.Context(), filters)
	if err != nil {
		return handler.error(c, err)
	}

	return c.Status(http.StatusOK).JSON(*cost)
}

func (handler *CostHandler) error(c *fiber.Ctx, err error) error {
	var vErrs validator.ValidationErrors

	switch {
	case errors.As(err, &vErrs):
		handler.logger.Info().Err(err).Msg("validation failed")
		return httpext.ValidationError(c, vErrs)
	default:
		handler.logger.Error().Err(err).Msg("failed to calculate total cost")
		return httpext.Error(c, http.StatusInternalServerError, "internal server error")
	}
}
