package helpers

import (
	"jxb-eprocurement/models"

	"github.com/gin-gonic/gin"
)

func GetUserPayload(c *gin.Context, user *models.USR_User) {
	userPayload, exist := c.Get("user")
	if !exist {
		return
	}

	userData, ok := userPayload.(*models.USR_User)
	if !ok {
		return
	}

	*user = *userData
}

func GetRolePayload(c *gin.Context, role *models.USR_Role) {
	rolePayload, exist := c.Get("role")
	if !exist {
		return
	}

	roleData, ok := rolePayload.(*models.USR_Role)
	if !ok {
		return
	}

	*role = *roleData
}
