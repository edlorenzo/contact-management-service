package response

import "github.com/gin-gonic/gin"

func EncodeJSONResp(c *gin.Context, data any, code int) {
	c.JSON(code, data)
}
