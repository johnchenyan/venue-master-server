package controller

import (
	"github.com/deatil/go-datebin/datebin"
	"github.com/deatil/go-events/events"

	"github.com/deatil/lakego-doak/lakego/facade"
	"github.com/deatil/lakego-doak/lakego/router"

	"github.com/deatil/lakego-doak-admin/admin/auth/auth"
	"github.com/deatil/lakego-doak-admin/admin/model"
	auth_password "github.com/deatil/lakego-doak-admin/admin/password"
	"github.com/deatil/lakego-doak-admin/admin/support/http/code"
	"github.com/deatil/lakego-doak-admin/admin/support/utils"
	passport_validate "github.com/deatil/lakego-doak-admin/admin/validate/passport"
)

/**
 * 登陆相关
 *
 * @create 2021-9-2
 * @author deatil
 */
type Passport struct {
	Base
}

// 验证码
// @Summary 登陆验证码
// @Description 登陆验证码
// @Tags 登陆相关
// @Accept application/json
// @Produce application/json
// @Success 200 {string} json "{"success": true, "code": 0, "message": "string", "data": ""}"
// @Header 200 {string} string "Lakego-Admin-Captcha-Id"
// @Router /passport/captcha [get]
// @x-lakego {"slug": "lakego-admin.passport.captcha"}
func (this *Passport) Captcha(ctx *router.Context) {
	id, b64s, err := facade.Captcha.Make()
	if err != nil {
		this.Error(ctx, "error", code.StatusError)
		return
	}

	key := facade.Config("auth").GetString("passport.header-captcha-key")

	this.SetHeader(ctx, key, id)
	this.SuccessWithData(ctx, "获取成功", router.H{
		"captcha":   b64s,
		"captchaId": id,
	})
}

// 账号登陆
// @Summary 账号登陆
// @Description 账号登陆
// @Tags 登陆相关
// @Accept application/json
// @Produce application/json
// @Param Lakego-Admin-Captcha-Id header string true "验证码字段"
// @Param name formData string true "账号"
// @Param password formData string true "密码"
// @Param captcha formData string true "验证码"
// @Success 200 {string} json "{"success": true, "code": 0, "message": "string", "data": ""}"
// @Router /passport/login [post]
// @x-lakego {"slug": "lakego-admin.passport.login"}
func (this *Passport) Login(ctx *router.Context) {
	// 接收数据
	post := make(map[string]any)
	this.ShouldBindJSON(ctx, &post)

	validateErr := passport_validate.Login(post)
	if validateErr != "" {
		this.Error(ctx, validateErr, code.LoginError)
		return
	}
	name := post["name"].(string)
	password := post["password"].(string)
	captchaCode := post["captcha"].(string)
	captchaId := post["captchaId"].(string)

	// 验证码检测
	ok := facade.Captcha.Verify(captchaId, captchaCode, true)
	if !ok {
		this.Error(ctx, "验证码错误", code.LoginError)
		return
	}

	// 用户信息
	admin := map[string]any{}
	err := model.NewAdmin().
		Where(&model.Admin{Name: name}).
		First(&admin).
		Error
	if err != nil {
		this.Error(ctx, "账号或者密码错误", code.LoginError)
		return
	}

	// 验证密码
	checkStatus := auth_password.CheckPassword(admin["password"].(string), password, admin["password_salt"].(string))
	if !checkStatus {
		events.DoAction("admin.passport-login.password-error", name)

		this.Error(ctx, "账号或者密码错误", code.LoginError)
		return
	}

	// 生成 token
	aud := auth.GetJwtAud(ctx)
	jwter := auth.NewWithAud(aud)

	// 账号ID
	adminid := admin["id"].(string)

	// token 数据
	tokenData := map[string]string{
		"id": adminid,
	}

	// 授权 token
	accessToken, err := jwter.MakeAccessToken(tokenData)
	if err != nil {
		//events.DoAction("admin.passport-login.make-accesstoken-fail", err.Error())

		this.Error(ctx, "授权token生成失败", code.LoginError)
		return
	}

	// 刷新 token
	refreshToken, err := jwter.MakeRefreshToken(tokenData)
	if err != nil {
		//events.DoAction("admin.passport-login.make-refreshtoken-fail", err.Error())

		this.Error(ctx, "刷新token生成失败", code.LoginError)
		return
	}

	// 授权 token 过期时间
	expiresIn := jwter.GetAccessExpiresIn()

	// 更新登录时间
	model.NewAdmin().
		Where("id = ?", adminid).
		Updates(map[string]any{
			"last_active": int(datebin.NowTimestamp()),
			"last_ip":     router.GetRequestIp(ctx),
		})

	// 数据输出
	this.SuccessWithData(ctx, "登录成功", router.H{
		"access_token":  accessToken,
		"expires_in":    expiresIn,
		"refresh_token": refreshToken,
	})
}

// 刷新 token
// @Summary 刷新 token
// @Description 刷新 token
// @Tags 登陆相关
// @Accept application/json
// @Produce application/json
// @Param refresh_token formData string true "刷新 Token"
// @Success 200 {string} json "{"success": true, "code": 0, "message": "string", "data": ""}"
// @Router /passport/refresh-token [post]
// @x-lakego {"slug": "lakego-admin.passport.refresh-token"}
func (this *Passport) RefreshToken(ctx *router.Context) {
	// 接收数据
	post := make(map[string]any)
	this.ShouldBindJSON(ctx, &post)

	//events.DoAction("admin.passport-refreshtoken.start", post)

	var refreshToken any
	var ok bool

	if refreshToken, ok = post["refresh_token"]; !ok {
		this.Error(ctx, "refreshToken不能为空", code.JwtRefreshTokenFail)
		return
	}

	refreshTokenPutTime, _ := facade.Cache.Get(utils.MD5(refreshToken.(string)))
	refreshTokenPutTime = refreshTokenPutTime.(string)
	if refreshTokenPutTime != "" {
		this.Error(ctx, "refreshToken已失效", code.JwtRefreshTokenFail)
		return
	}

	// jwt
	aud := auth.GetJwtAud(ctx)
	jwter := auth.NewWithAud(aud)

	// 拿取数据
	adminId := jwter.GetRefreshTokenData(refreshToken.(string), "id")
	if adminId == "" {
		this.Error(ctx, "刷新Token失败", code.JwtRefreshTokenFail)
		return
	}

	// token 数据
	tokenData := map[string]string{
		"id": adminId,
	}

	// 授权 token
	accessToken, err := jwter.MakeAccessToken(tokenData)
	if err != nil {
		//events.DoAction("admin.passport-refreshtoken.make-accesstoken-fail", err.Error())

		this.Error(ctx, "生成 access_token 失败", code.JwtRefreshTokenFail)
		return
	}

	// 授权 token 过期时间
	expiresIn := jwter.GetAccessExpiresIn()

	events.DoAction("admin.passport-refreshtoken.end", adminId)

	// 数据输出
	this.SuccessWithData(ctx, "获取成功", router.H{
		"access_token": accessToken,
		"expires_in":   expiresIn,
	})
}

// 账号退出
// @Summary 当前账号退出
// @Description 当前账号退出
// @Tags 登陆相关
// @Accept  application/json
// @Produce application/json
// @Param Authorization header   string false "Bearer 用户令牌"
// @Param refresh_token formData string true  "刷新 Token"
// @Success 200 {string} json "{"success": true, "code": 0, "message": "string", "data": ""}"
// @Router /passport/logout [delete]
// @Security Bearer
// @x-lakego {"slug": "lakego-admin.passport.logout"}
func (this *Passport) Logout(ctx *router.Context) {
	// 接收数据
	post := make(map[string]any)
	this.ShouldBindJSON(ctx, &post)

	var refreshToken any
	var ok bool

	if refreshToken, ok = post["refresh_token"]; !ok {
		this.Error(ctx, "refreshToken 不能为空", code.JwtRefreshTokenFail)
		return
	}

	c := facade.Cache

	refreshTokenPutString, _ := c.Get(utils.MD5(refreshToken.(string)))
	refreshTokenPutString = refreshTokenPutString.(string)
	if refreshTokenPutString != "" {
		this.Error(ctx, "refreshToken 已失效", code.JwtRefreshTokenFail)
		return
	}

	// jwt
	aud := auth.GetJwtAud(ctx)
	jwter := auth.NewWithAud(aud)

	// 拿取数据
	claims, claimsErr := jwter.GetRefreshTokenClaims(refreshToken.(string))
	if claimsErr != nil {
		this.Error(ctx, "refreshToken 已失效", code.JwtRefreshTokenFail)
		return
	}

	// 当前账号ID
	adminId := jwter.GetDataFromTokenClaims(claims, "id")

	// 过期时间
	exp := jwter.GetFromTokenClaims(claims, "exp")
	iat := jwter.GetFromTokenClaims(claims, "iat")
	refreshTokenExpiresIn := exp.(float64) - iat.(float64)

	nowAdminId, _ := ctx.Get("admin_id")
	if adminId != nowAdminId.(string) {
		this.Error(ctx, "退出失败", code.JwtRefreshTokenFail)
		return
	}

	// 当前 accessToken
	accessToken, _ := ctx.Get("access_token")

	// 加入黑名单

	c.Put(utils.MD5(accessToken.(string)), "no", int64(refreshTokenExpiresIn))
	c.Put(utils.MD5(refreshToken.(string)), "no", int64(refreshTokenExpiresIn))
	//events.DoAction("admin.passport-logout.end", adminId)

	// 数据输出
	this.Success(ctx, "退出成功")
}
