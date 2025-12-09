package handlers

import (
	"errors"
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

func (handler *SubscriptionHandler) Create(c *fiber.Ctx) error {
	req := new(dto.SubscriptionDTO)

	if err := c.BodyParser(req); err != nil {
		handler.logger.Warn().Err(err).Msg("failed to parse subscription request")
		return httpext.Error(c, http.StatusBadRequest, "bad request")
	}

	resp, err := handler.service.Create(c.Context(), *req)
	if err != nil {
		return handler.error(c, err)
	}

	handler.logger.Info().
		Int("id", resp.ID).
		Str("service_name", resp.ServiceName).
		Str("user_id", resp.UserID).
		Msg("subscription created")
	return c.Status(http.StatusCreated).JSON(*resp)
}

func (handler *SubscriptionHandler) Update(c *fiber.Ctx) error {
	req := new(dto.SubscriptionDTO)

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
		return handler.error(c, err)
	}

	handler.logger.Info().
		Int("id", resp.ID).
		Str("service_name", resp.ServiceName).
		Str("user_id", resp.UserID).
		Msg("subscription updated")
	return c.Status(http.StatusOK).JSON(*resp)
}

func (handler *SubscriptionHandler) Delete(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return httpext.Error(c, http.StatusBadRequest, "bad request")
	}

	err = handler.service.Delete(c.Context(), id)
	if err != nil {
		return handler.error(c, err)
	}

	handler.logger.Info().
		Int("id", id).
		Msg("subscription deleted")
	return c.SendStatus(http.StatusNoContent)
}

func (handler *SubscriptionHandler) Index(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return httpext.Error(c, http.StatusBadRequest, "bad request")
	}

	resp, err := handler.service.Index(c.Context(), id)
	if err != nil {
		return handler.error(c, err)
	}

	return c.Status(http.StatusOK).JSON(*resp)
}

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

	subscriptions, total, err := handler.service.List(c.Context(), filters)
	if err != nil {
		return handler.error(c, err)
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"page":  filters.Page,
		"limit": filters.Limit,
		"total": total,
		"data":  subscriptions,
	})
}

func (handler *SubscriptionHandler) error(c *fiber.Ctx, err error) error {
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
		handler.logger.Error().Err(err).Msg("failed to create subscription")
		return httpext.Error(c, http.StatusInternalServerError, "internal server error")
	}
}
