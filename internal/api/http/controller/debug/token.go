package debug

import (
	"github.com/dysodeng/app/internal/pkg/api/token"
	"github.com/gin-gonic/gin"
)

func Token(ctx *gin.Context) {
	t, _ := token.GenerateToken("user", map[string]interface{}{
		"user_id": 1,
	}, nil)

	ctx.JSON(200, gin.H{
		"token": t,
	})
}
