package route

import (
	"github.com/deatil/lakego-doak/lakego/router"

	"github.com/deatil/lakego-doak-admin/admin/controller"
)

/**
 * 后台路由
 */
func Route(engine router.IRouter) {
	// 登陆
	passportController := new(controller.Passport)
	engine.GET("/passport/captcha", passportController.Captcha)
	engine.POST("/passport/login", passportController.Login)
	engine.PUT("/passport/refresh-token", passportController.RefreshToken)
	engine.DELETE("/passport/logout", passportController.Logout)

	// 个人信息
	profileController := new(controller.Profile)
	engine.GET("/profile", profileController.Index)
	engine.PUT("/profile", profileController.Update)
	engine.PATCH("/profile/avatar", profileController.UpdateAvatar)
	engine.PATCH("/profile/password", profileController.UpdatePasssword)
	engine.GET("/profile/rules", profileController.Rules)

	// 上传
	uploadController := new(controller.Upload)
	engine.POST("/upload/file", uploadController.File)

	// 附件
	attachmentController := new(controller.Attachment)
	engine.GET("/attachment", attachmentController.Index)
	engine.GET("/attachment/:id", attachmentController.Detail)
	engine.PATCH("/attachment/:id/enable", attachmentController.Enable)
	engine.PATCH("/attachment/:id/disable", attachmentController.Disable)
	engine.DELETE("/attachment/:id", attachmentController.Delete)
	engine.GET("/attachment/downcode/:id", attachmentController.DownloadCode)
	engine.GET("/attachment/download/:code", attachmentController.Download)

	// 管理员
	adminController := new(controller.Admin)
	engine.GET("/admin", adminController.Index)
	engine.GET("/admin/groups", adminController.Groups)
	engine.GET("/admin/:id", adminController.Detail)
	engine.GET("/admin/:id/rules", adminController.Rules)
	engine.POST("/admin", adminController.Create)
	engine.PUT("/admin/:id", adminController.Update)
	engine.DELETE("/admin/:id", adminController.Delete)
	engine.PATCH("/admin/:id/enable", adminController.Enable)
	engine.PATCH("/admin/:id/disable", adminController.Disable)
	engine.PATCH("/admin/:id/avatar", adminController.UpdateAvatar)
	engine.PATCH("/admin/:id/password", adminController.UpdatePasssword)
	engine.PATCH("/admin/:id/access", adminController.Access)
	engine.DELETE("/admin/logout/:refreshToken", adminController.Logout)
	engine.PUT("/admin/reset-permission", adminController.ResetPermission)

	// 系统信息
	systemController := new(controller.System)
	engine.GET("/system/info", systemController.Info)
	engine.GET("/system/rules", systemController.Rules)

	// 设置场地相关路由
	VenueRoutes(engine)

	// 设置收益相关路由
	IncomeRoutes(engine)

	// 设置托管相关路由
	CustodyRoutes(engine)

	SettlementRoutes(engine)

	BTCMiningPoolRoutes(engine)
}

/**
 * 后台管理员路由
 */
func AdminRoute(engine router.IRouter) {
	// 权限菜单
	authRuleController := new(controller.AuthRule)
	engine.GET("/auth/rule", authRuleController.Index)
	engine.GET("/auth/rule/tree", authRuleController.IndexTree)
	engine.GET("/auth/rule/children", authRuleController.IndexChildren)
	engine.GET("/auth/rule/:id", authRuleController.Detail)
	engine.POST("/auth/rule", authRuleController.Create)
	engine.PUT("/auth/rule/:id", authRuleController.Update)
	engine.DELETE("/auth/rule/clear", authRuleController.Clear)
	engine.DELETE("/auth/rule/:id", authRuleController.Delete)
	engine.PATCH("/auth/rule/:id/sort", authRuleController.Listorder)
	engine.PATCH("/auth/rule/:id/enable", authRuleController.Enable)
	engine.PATCH("/auth/rule/:id/disable", authRuleController.Disable)

	// 权限分组
	authGroupController := new(controller.AuthGroup)
	engine.GET("/auth/group", authGroupController.Index)
	engine.GET("/auth/group/tree", authGroupController.IndexTree)
	engine.GET("/auth/group/children", authGroupController.IndexChildren)
	engine.GET("/auth/group/:id", authGroupController.Detail)
	engine.POST("/auth/group", authGroupController.Create)
	engine.PUT("/auth/group/:id", authGroupController.Update)
	engine.DELETE("/auth/group/:id", authGroupController.Delete)
	engine.PATCH("/auth/group/:id/sort", authGroupController.Listorder)
	engine.PATCH("/auth/group/:id/enable", authGroupController.Enable)
	engine.PATCH("/auth/group/:id/disable", authGroupController.Disable)
	engine.PATCH("/auth/group/:id/access", authGroupController.Access)
}

/**
 * 场地相关路由
 */
func VenueRoutes(engine router.IRouter) {
	// 场地模板相关路由
	venueController := new(controller.Venue)
	engine.GET("/venue-templates", venueController.Index)          // 获取所有场地模板
	engine.GET("/venue-templates/list", venueController.List)      // 获取所有场地模板名称
	engine.GET("/venue-templates/:name", venueController.Detail)   // 获取指定场地模板
	engine.POST("/venue-templates/new", venueController.Create)    // 创建新的场地模板
	engine.POST("/venue-templates/update", venueController.Update) // 更新新的场地模板
	engine.DELETE("/venue-templates/:id", venueController.Delete)  // 删除指定场地模板

	// 根据模版ID获取字段
	// 获取对应模版ID的fields 以及attributes
	engine.GET("/venue/:templateID/GetFieldsByTemplateID", venueController.GetFieldsByTemplateID)

	// 获取对应模版ID的fields 以及attributes
	engine.GET("/venue/:templateID/GetVenueRecordByTemplateID", venueController.GetVenueRecordAttributesByTemplateID)

	// 添加场地属性记录
	engine.POST("/venue/newVenueRecord", venueController.NewVenueRecord)

	// 删除场地属性记录
	engine.DELETE("/venue/deleteVenueRecord/:recordID", venueController.DeleteVenueRecord)

	// 更新场地属性记录
	engine.POST("/venue/updateVenueRecordAttributes", venueController.UpdateVenueRecordAttributes)

}

func IncomeRoutes(engine router.IRouter) {
	// 收益相关路由
	incomeController := new(controller.Income)
	//engine.GET("/venue-templates", venueController.Index)          // 获取所有场地模板
	//engine.GET("/venue-templates/list", venueController.List)      // 获取所有场地模板名称
	//engine.GET("/venue-templates/:name", venueController.Detail)   // 获取指定场地模板
	//engine.POST("/venue-templates/new", venueController.Create)    // 创建新的场地模板
	//engine.POST("/venue-templates/update", venueController.Update) // 创建新的场地模板
	//engine.DELETE("/venue-templates/:id", venueController.Delete)  // 删除指定场地模板

	// 根据模版ID获取字段
	// 获取对应模版ID的fields 以及attributes
	engine.GET("/link/list", incomeController.List)

	// 获取对应模版ID的fields 以及attributes
	//engine.GET("/venue/:templateID/GetVenueRecordByTemplateID", venueController.GetVenueRecordAttributesByTemplateID)
	//
	//// 添加场地属性记录
	//engine.POST("/venue/newVenueRecord", venueController.NewVenueRecord)
	//
	//// 删除场地属性记录
	//engine.DELETE("/venue/deleteVenueRecord/:recordID", venueController.DeleteVenueRecord)
	//
	//// 更新场地属性记录
	//engine.POST("/venue/updateVenueRecordAttributes", venueController.UpdateVenueRecordAttributes)

}

func CustodyRoutes(engine router.IRouter) {
	// 托管相关路由
	CustodyController := new(controller.Custody)

	engine.GET("/custody/custodyInfoList", CustodyController.ListCustodyInfo)

	// 新增托管信息
	engine.POST("/custody/newCustodyInfo", CustodyController.CreateCustodyInfo)

	// 删除托管信息
	engine.DELETE("/custody/deleteCustodyInfo/:custodyInfoId", CustodyController.DeleteCustodyInfo)

	// 更新托管信息
	engine.POST("/custody/updateCustodyInfo", CustodyController.UpdateCustodyInfo)

	// 获取托管统计信息
	engine.GET("/custody/custodyStatisticsList/:timeRange", CustodyController.ListCustodyStatistics)

	// 获取价格信息
	engine.GET("/custody/dailyAveragePriceList", CustodyController.ListDailyAveragePrice)

	// 获取托管费曲线图数据
	engine.GET("/custody/hostingFeeRatioList", CustodyController.ListHostingFeeRatio)
}

func SettlementRoutes(engine router.IRouter) {
	SettlementController := new(controller.Settlement)

	engine.POST("/settlement/findSettlementData", SettlementController.FindSettlementData)

	engine.POST("/settlement/findSettlementDataWithPagination", SettlementController.FindSettlementDataWithPagination)

	engine.POST("/settlement/findSettlementAverage", SettlementController.FindSettlementAverage)

	engine.GET("/settlement/settlementPointList/:type", SettlementController.SettlementPointList)

	engine.POST("/settlement/downloadSettlementData", SettlementController.DownLoadSettlementData)
}

func BTCMiningPoolRoutes(engine router.IRouter) {
	btcMiningPoolController := new(controller.BtcMiningPool)

	engine.GET("/miningPool/listBtcMiningPool/:poolType/:poolCategory", btcMiningPoolController.ListBtcMiningPool)

	engine.POST("/miningPool/createBtcMiningPool", btcMiningPoolController.CreateBtcMiningPool)

	engine.POST("/miningPool/updateBtcMiningPool", btcMiningPoolController.UpdateBtcMiningPool)

	engine.DELETE("/miningPool/deleteBtcMiningPool/:miningPoolId", btcMiningPoolController.DeleteBtcMiningPool)

	engine.GET("/miningPool/listBtcMiningPoolHashRate/:poolType/:poolCategory", btcMiningPoolController.ListBtcMiningPoolHashRate)

	// overview
	engine.GET("/miningPool/getTotalRealTimeStatus/:poolType", btcMiningPoolController.TotalRealTimeStatus)
	engine.GET("/miningPool/getTotalLastDayStatus/:poolType", btcMiningPoolController.TotalLastDayStatus)
	engine.GET("/miningPool/getTotalLastWeekStatus/:poolType", btcMiningPoolController.TotalLastWeekStatus)
	engine.GET("/miningPool/getHashRateEfficiency/:poolType/:day", btcMiningPoolController.LastestHashRateEfficiency)
	engine.GET("/miningPool/getLastestHashRate/:poolType/:day", btcMiningPoolController.LastestHashRate)
}
