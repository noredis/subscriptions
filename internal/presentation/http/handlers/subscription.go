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
	app.Post("/subscriptions", handler.CreateSubscription)
	app.Put("/subscriptions/:id", handler.UpdateSubscription)
}

func (handler *SubscriptionHandler) CreateSubscription(c *fiber.Ctx) error {
	req := new(dto.SubscriptionDTO)

	if err := c.BodyParser(req); err != nil {
		handler.logger.Warn().Err(err).Msg("failed to parse subscription request")
		return httpext.Error(c, http.StatusBadRequest, "bad request")
	}

	resp, err := handler.service.CreateSubscription(c.Context(), *req)
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

func (handler *SubscriptionHandler) UpdateSubscription(c *fiber.Ctx) error {
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

	resp, err := handler.service.UpdateSubscription(c.Context(), *req, id)
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
