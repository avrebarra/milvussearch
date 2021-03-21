package server

import (
	"time"

	"github.com/labstack/echo/v4"
)

type HandlerCommon struct{ Config Config }

func (h *HandlerCommon) PingV1() echo.HandlerFunc {
	type RequestData struct {
	}

	type ResponseData struct {
		State          string    `json:"state"`
		StartTimestamp time.Time `json:"start_time"`
		Uptime         string    `json:"uptime"`
	}

	starttime := time.Now()

	return func(c echo.Context) error {
		// var ctx context.Context = context.Background()
		var response Resp = ListResp[ServerRespSuccess]

		response.Data = ResponseData{
			State:          "ready",
			StartTimestamp: starttime,
			Uptime:         time.Since(starttime).Round(time.Second).String(),
		}

		return c.JSON(200, response.Normalize())
	}
}
