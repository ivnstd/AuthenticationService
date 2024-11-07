package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func newErrorResponse(c *gin.Context, statusCode int, message string, description ...string) {
	response := gin.H{
		"status":  "error",
		"message": message,
	}

	logrus.Errorf("%v error.\tmessage: %v", statusCode, message)
	if len(description) > 0 {
		logrus.Errorf("Error description: %s", description[0])
	}
	c.AbortWithStatusJSON(statusCode, response)
}

func newSuccessResponse(c *gin.Context, statusCode int, dataKey string, dataValue interface{}) {
	response := gin.H{
		"status": "success",
		dataKey:  dataValue,
	}

	logrus.Infof("%v success.\t%v: %+v", statusCode, dataKey, dataValue)
	c.JSON(statusCode, response)
}
