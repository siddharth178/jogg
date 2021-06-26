package web

import (
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type NopController struct {
	logger *zap.SugaredLogger
}

func (c *NopController) Nop(ctx echo.Context) error {
	c.logger.Info("nop called")
	return nil
}
