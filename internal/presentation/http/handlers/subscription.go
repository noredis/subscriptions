package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/noredis/subscriptions/internal/application/appservice"
	"github.com/noredis/subscriptions/internal/application/dto"
	"github.com/noredis/subscriptions/internal/domain/failure"
	"github.com/noredis/subscriptions/pkg/httpext"
	"github.com/rs/zerolog"
)

type SubscriptionHandler struct {
	service *appservice.SubscriptionService
	logger  *zerolog.Logger
}

func NewSubscriptionHandler(
	service *appservice.SubscriptionService,
	logger *zerolog.Logger,
) *SubscriptionHandler {
	return &SubscriptionHandler{
		service: service,
		logger:  logger,
	}
}

func (handler *SubscriptionHandler) Register(app *fiber.App) {
	app.Post("/subscriptions", handler.Create)
	app.Put("/subscriptions/:id", handler.Update)
	app.Delete("/subscriptions/:id", handler.Delete)
	app.Get("/subscriptions/:id", handler.Index)
	app.Get("/subscriptions", handler.List)
}

// Create создаёт новую подписку.
//
// @Summary      Создать подписку
// @Description  Создаёт новую подписку для пользователя и сервиса.
// @Tags         subscriptions
// @Accept       json
// @Produce      json
// @Param        request  body      dto.SubscriptionRequest   true  "Данные для создания подписки"
// @Success      201      {object}  dto.SubscriptionResponse  "Подписка успешно создана"
// @Failure      400      {object}  httpext.FiberError        "Некорректный запрос"
// @Failure      409      {object}  httpext.FiberError        "Подписка уже существует"
// @Failure      422      {object}  httpext.FiberError        "Ошибка валидации"
// @Failure      500      {object}  httpext.FiberError        "Внутренняя ошибка сервера"
// @Router       /subscriptions [post]
func (handler *SubscriptionHandler) Create(c *fiber.Ctx) error {
	req := new(dto.SubscriptionRequest)

	if err := c.BodyParser(req); err != nil {
		handler.logger.Warn().Err(err).Msg("failed to parse subscription request")
		return httpext.Error(c, http.StatusBadRequest, "bad request")
	}

	resp, err := handler.service.Create(c.Context(), *req)
	if err != nil {
		return handler.error(c, err, "failed to create subscription")
	}

	handler.logger.Info().
		Int("id", resp.ID).
		Str("service_name", resp.ServiceName).
		Str("user_id", resp.UserID).
		Msg("subscription created")
	c.Location(fmt.Sprintf("/subscriptions/%d", resp.ID))
	return c.Status(http.StatusCreated).JSON(*resp)
}

// Update обновляет данные подписки.
//
// @Summary      Обновить подписку
// @Description  Обновляет данные существующей подписки по её идентификатору.
// @Tags         subscriptions
// @Accept       json
// @Produce      json
// @Param        id       path      int                        true  "ID подписки"
// @Param        request  body      dto.SubscriptionRequest    true  "Данные для обновления подписки"
// @Success      200      {object}  dto.SubscriptionResponse   "Подписка успешно обновлена"
// @Failure      400      {object}  httpext.FiberError         "Некорректный запрос"
// @Failure      404      {object}  httpext.FiberError         "Подписка не найдена"
// @Failure      409      {object}  httpext.FiberError         "Подписка уже существует"
// @Failure      422      {object}  httpext.FiberError         "Ошибка валидации"
// @Failure      500      {object}  httpext.FiberError         "Внутренняя ошибка сервера"
// @Router       /subscriptions/{id} [put]
func (handler *SubscriptionHandler) Update(c *fiber.Ctx) error {
	req := new(dto.SubscriptionRequest)

	if err := c.BodyParser(req); err != nil {
		handler.logger.Warn().Err(err).Msg("failed to parse subscription request")
		return httpext.Error(c, http.StatusBadRequest, "bad request")
	}

	idParam := c.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return httpext.Error(c, http.StatusBadRequest, "bad request")
	}

	resp, err := handler.service.Update(c.Context(), *req, id)
	if err != nil {
		return handler.error(c, err, "failed to update subscription")
	}

	handler.logger.Info().
		Int("id", resp.ID).
		Str("service_name", resp.ServiceName).
		Str("user_id", resp.UserID).
		Msg("subscription updated")
	return c.Status(http.StatusOK).JSON(*resp)
}

// Delete удаляет подписку.
//
// @Summary      Удалить подписку
// @Description  Удаляет подписку по её идентификатору.
// @Tags         subscriptions
// @Param        id   path      int                   true  "ID подписки"
// @Success      204  "Подписка успешно удалена"
// @Failure      400  {object}  httpext.FiberError    "Некорректный запрос"
// @Failure      404  {object}  httpext.FiberError    "Подписка не найдена"
// @Failure      500  {object}  httpext.FiberError    "Внутренняя ошибка сервера"
// @Router       /subscriptions/{id} [delete]
func (handler *SubscriptionHandler) Delete(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return httpext.Error(c, http.StatusBadRequest, "bad request")
	}

	err = handler.service.Delete(c.Context(), id)
	if err != nil {
		return handler.error(c, err, "failed to delete subscription")
	}

	handler.logger.Info().
		Int("id", id).
		Msg("subscription deleted")
	return c.SendStatus(http.StatusNoContent)
}

// Index возвращает информацию о конкретной подписке.
//
// @Summary      Получить подписку по ID
// @Description  Возвращает данные подписки по её идентификатору.
// @Tags         subscriptions
// @Produce      json
// @Param        id   path      int                       true  "ID подписки"
// @Success      200  {object}  dto.SubscriptionResponse  "Данные подписки"
// @Failure      400  {object}  httpext.FiberError        "Некорректный идентификатор"
// @Failure      404  {object}  httpext.FiberError        "Подписка не найдена"
// @Failure      500  {object}  httpext.FiberError        "Внутренняя ошибка сервера"
// @Router       /subscriptions/{id} [get]
func (handler *SubscriptionHandler) Index(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return httpext.Error(c, http.StatusBadRequest, "bad request")
	}

	resp, err := handler.service.Index(c.Context(), id)
	if err != nil {
		return handler.error(c, err, "failed to index subscription")
	}

	return c.Status(http.StatusOK).JSON(*resp)
}

// List возвращает список подписок с пагинацией и фильтрами.
//
// @Summary      Получить список подписок
// @Description  Возвращает список подписок с поддержкой пагинации и фильтрации.
// @Tags         subscriptions
// @Produce      json
// @Param        page          query     int     false  "Номер страницы"         default(1)
// @Param        limit         query     int     false  "Количество элементов"   default(20)
// @Param        service_name  query     string  false  "Фильтр по имени сервиса"
// @Param        user_id       query     string  false  "Фильтр по ID пользователя"
// @Param        start_date    query     string  false  "Фильтр по дате начала (MM-YYYY)"
// @Param        end_date      query     string  false  "Фильтр по дате окончания (MM-YYYY)"
// @Success      200  {object}  dto.SubscriptionListResponse  "Список подписок"
// @Failure      400  {object}  httpext.FiberError            "Некорректный запрос"
// @Failure      422  {object}  httpext.FiberError            "Ошибка валидации"
// @Failure      500  {object}  httpext.FiberError            "Внутренняя ошибка сервера"
// @Router       /subscriptions [get]
func (handler *SubscriptionHandler) List(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 20)

	filters := dto.SubscriptionFilterDTO{
		Page:        page,
		Limit:       limit,
		ServiceName: c.Query("service_name"),
		UserID:      c.Query("user_id"),
		StartDate:   c.Query("start_date"),
		EndDate:     c.Query("end_date"),
	}

	resp, err := handler.service.List(c.Context(), filters)
	if err != nil {
		return handler.error(c, err, "failed to list subscriptions")
	}

	return c.Status(http.StatusOK).JSON(*resp)
}

func (handler *SubscriptionHandler) error(c *fiber.Ctx, err error, err500msg string) error {
	var vErrs validator.ValidationErrors

	switch {
	case errors.As(err, &vErrs):
		handler.logger.Info().Err(err).Msg("validation for subscription failed")
		return httpext.ValidationError(c, vErrs)
	case errors.Is(err, failure.ErrUserAlreadyHasThisSubscription):
		handler.logger.Info().Err(err).Msg("subscription already exists")
		return httpext.Error(c, http.StatusConflict, err.Error())
	case errors.Is(err, failure.ErrSubscriptionNotFound):
		handler.logger.Info().Err(err).Msg("subscription not found")
		return httpext.Error(c, http.StatusNotFound, err.Error())
	default:
		handler.logger.Error().Err(err).Msg(err500msg)
		return httpext.Error(c, http.StatusInternalServerError, "internal server error")
	}
}
