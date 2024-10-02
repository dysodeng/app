package provider

import (
	"fmt"

	linkIface "github.com/dysodeng/app/internal/service/contracts/link"
	linkImpl "github.com/dysodeng/app/internal/service/link"

	areaImpl "github.com/dysodeng/app/internal/service/common/area"
	commonMessageImpl "github.com/dysodeng/app/internal/service/common/message"
	areaIface "github.com/dysodeng/app/internal/service/contracts/common/area"
	commonMessageIface "github.com/dysodeng/app/internal/service/contracts/common/message"

	"github.com/dysodeng/app/internal/config"
	"github.com/dysodeng/app/internal/service"

	"github.com/defval/di"
	"github.com/pkg/errors"
)

// ServiceProvider 服务容器初始化
func ServiceProvider() {
	if config.App.Env != config.Prod {
		di.SetTracer(&di.StdTracer{})
	}

	var err error
	service.Container, err = di.New(
		// 公共服务
		di.Provide(areaImpl.NewAreaService, di.As(new(areaIface.ServiceInterface))),
		di.Provide(commonMessageImpl.NewCodeMessageService, di.As(new(commonMessageIface.CodeMessageServiceInterface))),
		di.Provide(linkImpl.NewLinkService, di.As(new(linkIface.ServiceInterface))),
	)

	if err != nil {
		err = errors.Wrap(err, "service provider error.")
		fmt.Printf("%+v\n", err)
		panic(err)
	}
}
