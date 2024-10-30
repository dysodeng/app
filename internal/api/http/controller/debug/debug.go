package debug

import (
	"github.com/dysodeng/app/internal/dal/model/common"
	"github.com/dysodeng/app/internal/pkg/api"
	"github.com/dysodeng/app/internal/pkg/api/token"
	"github.com/dysodeng/app/internal/pkg/db"
	"github.com/gin-gonic/gin"
)

func Token(ctx *gin.Context) {
	t, _ := token.GenerateToken("user", map[string]interface{}{
		"user_id": 1,
	}, nil)
	ctx.JSON(200, api.Success(ctx, t))
}

func GormLogger(ctx *gin.Context) {
	var mailConfig common.MailConfig
	db.DB().WithContext(ctx).Where("a=?", "b").First(&mailConfig)
	db.DB().WithContext(ctx).Debug().First(&mailConfig)
	ctx.JSON(200, api.Success(ctx, mailConfig))
}
