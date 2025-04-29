package controller

import (
	"fmt"
	"github.com/deatil/lakego-doak-admin/admin/model"
	"github.com/gin-gonic/gin"
)

type Income struct {
	Base
}

func (this *Income) List(ctx *gin.Context) {
	var linkInfos []model.LinkInfo
	err := model.NewLinkInfoModel().Find(&linkInfos).Error
	if err != nil {
		this.Error(ctx, fmt.Sprintf("取观察者链接失败: %s", err.Error()))
		return
	}
	
	this.SuccessWithData(ctx, "获取成功", linkInfos)
}
