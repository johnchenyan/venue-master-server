package controller

import (
	"fmt"
	"github.com/deatil/lakego-doak-admin/admin/model"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// VenueTemplate ******************************** 场地模版 ********************************* //

type Venue struct {
	Base
}

type fieldInfo struct {
	ID        uint
	FieldName string
}
type TemplateInfo struct {
	ID           uint
	TemplateName string
	FieldInfo    []fieldInfo
}

type TemplateIDList struct {
	ID           uint
	TemplateName string
}

// Index 获取所有场地模板
func (this *Venue) Index(ctx *gin.Context) {
	var templates []model.VenueTemplates
	err := model.NewVenueTemplate().Preload("Fields").Find(&templates).Error
	if err != nil {
		this.Error(ctx, fmt.Sprintf("取场地模板失败: %s", err.Error()))
		return
	}

	result := []TemplateInfo{}
	for _, template := range templates {
		fields := []fieldInfo{}
		for _, field := range template.Fields {
			fields = append(fields, fieldInfo{
				ID:        field.ID,
				FieldName: field.FieldName,
			})
		}
		result = append(result, TemplateInfo{
			ID:           template.ID,
			TemplateName: template.TemplateName,
			FieldInfo:    fields,
		})
	}

	this.SuccessWithData(ctx, "获取成功", result)
}

// List Index 获取所有场地模板名称
func (this *Venue) List(ctx *gin.Context) {
	var templates []model.VenueTemplates
	err := model.NewVenueTemplate().Find(&templates).Error
	if err != nil {
		this.Error(ctx, fmt.Sprintf("取场地模板失败: %s", err.Error()))
		return
	}

	result := []TemplateIDList{}
	for _, template := range templates {
		result = append(result, TemplateIDList{template.ID, template.TemplateName})
	}

	this.SuccessWithData(ctx, "获取成功", result)
}

// Detail 获取指定场地模板
func (this *Venue) Detail(ctx *gin.Context) {
	name := ctx.Param("name")
	if name == "" {
		this.Error(ctx, "模板名称不能为空")
		return
	}

	var template model.VenueTemplates
	err := model.NewVenueTemplate().Preload("Fields").Where("template_name = ?", name).First(&template).Error
	if err != nil {
		this.Error(ctx, fmt.Sprintf("模板不存在: %s", err.Error()))
		return
	}

	fields := []string{}
	for _, field := range template.Fields {
		fields = append(fields, field.FieldName)
	}

	this.SuccessWithData(ctx, "获取成功", gin.H{
		template.TemplateName: fields,
	})
}

// Create 创建新的场地模板
func (this *Venue) Create(ctx *gin.Context) {
	var data VenueTemplateNew
	if err := this.ShouldBindJSON(ctx, &data); err != nil {
		this.Error(ctx, "请求数据不正确")
		return
	}

	if data.TemplateName == "" {
		this.Error(ctx, "模版名称不能为空")
		return
	}

	ID, err := createVenueTemplate(data.TemplateName)
	if err != nil {
		this.Error(ctx, "创建模版失败")
	}

	for _, field := range data.Fields {
		if err := createTemplateField(ID, field.Value); err != nil {
			this.Error(ctx, "创建模版成功，新增模版字段失败")
		}
	}

	this.Success(ctx, "添加模版成功！")
}

// Update 更新场地模板以及场地属性
func (this *Venue) Update(ctx *gin.Context) {
	var data VenueTemplateChange
	if err := this.ShouldBindJSON(ctx, &data); err != nil {
		this.Error(ctx, "请求数据不正确")
		return
	}

	if data.TemplateNameAfter != data.TemplateNameBefore {
		// 更新模版名称
		if err := updateTemplateName(data.ID, data.TemplateNameAfter); err != nil {
			this.Error(ctx, "数据库更新错误")
			return
		}
	}

	for _, field := range data.Fields {
		if field.Status == VenueTemplateStatusNew { // 创建模版字段
			// 创建新的模板字段
			if err := createTemplateField(data.ID, field.Value); err != nil {
				this.Error(ctx, "创建模版字段失败")
				return
			}
		}
		if field.Status == VenueTemplateStatusModified { // 更新模版字段
			if err := updateTemplateField(uint(field.ID), field.Value); err != nil {
				this.Error(ctx, "更新模版字段错误")
				return
			}
		}

		if field.Status == VenueTemplateStatusDeleted { // 删除模版字段
			if err := deleteTemplateField(uint(field.ID)); err != nil {
				fmt.Println("删除模版字段错误")
				this.Error(ctx, "删除模版字段错误,先检查该字段是否有记录数据！")
				return
			}
		}
	}

	this.Success(ctx, "数据库更新成功！")
}

// Delete 删除指定场地模板
func (this *Venue) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "模板ID不能为空"})
		return
	}

	err := deleteTemplateByID(id)
	if err != nil {
		this.Error(ctx, "删除模版失败")
	}

	this.Success(ctx, "删除模版成功")
}

// GetFieldsByTemplateID 根据模版名称获取字段
func (this *Venue) GetFieldsByTemplateID(ctx *gin.Context) {
	templateId := ctx.Param("templateID")
	if templateId == "" {
		this.Error(ctx, "模板ID不能为空")
		return
	}

	// 解析模板ID
	id, err := strconv.Atoi(templateId)
	if err != nil {
		this.Error(ctx, "无效的模板 ID")
		return
	}

	field, err := getTemplateField(uint(id))
	if err != nil {
		this.Error(ctx, "获取字段失败")
	}

	this.SuccessWithData(ctx, "获取成功", field)
}

// GetVenueRecordAttributesByTemplateID 获取对应模版ID的fields 以及attributes
func (this *Venue) GetVenueRecordAttributesByTemplateID(ctx *gin.Context) {
	templateId := ctx.Param("templateID")
	if templateId == "" {
		this.Error(ctx, "模板ID不能为空")
		return
	}

	// 解析模板ID
	id, err := strconv.Atoi(templateId)
	if err != nil {
		this.Error(ctx, "无效的模板 ID")
		return
	}

	// 查询并联合获取相关记录和属性
	var records []struct {
		ID               uint   `gorm:"primaryKey"`
		AttributeFieldId uint   `gorm:"column:attribute_field_id"`
		AttributeName    string `gorm:"column:attribute_field_name"`
		AttributeValue   string `gorm:"column:attribute_field_value"`
	}

	venueRecordTable := (&model.VenueRecords{}).TableName()                  // 获取 VenueRecord 的表名
	venueRecordAttributeTable := (&model.VenueRecordAttribute{}).TableName() // 假设您也有这个表的常量

	// 获取模版对应的
	err = model.NewVenueRecord().Table(venueRecordTable+" AS r").
		Select("r.id, ra.field_id as attribute_field_id, ra.field_name AS attribute_field_name, ra.field_value AS attribute_field_value").
		Joins("LEFT JOIN "+venueRecordAttributeTable+" AS ra ON r.id = ra.record_id").
		Where("r.template_id = ?", id).
		Scan(&records).Error

	if err != nil {
		this.Error(ctx, "查询失败: "+err.Error())
		return
	}

	// 处理查询结果
	response := make(map[uint][]struct {
		FieldId    uint   `json:"field_id"`
		FieldName  string `json:"field_name"`
		FieldValue string `json:"field_value"`
	})

	for _, record := range records {
		if record.AttributeValue != "" {
			response[record.ID] = append(response[record.ID], struct {
				FieldId    uint   `json:"field_id"`
				FieldName  string `json:"field_name"`
				FieldValue string `json:"field_value"`
			}{
				FieldId:    record.AttributeFieldId,
				FieldName:  record.AttributeName,
				FieldValue: record.AttributeValue,
			})
		}
	}

	this.SuccessWithData(ctx, "获取成功", response)
}

// NewVenueRecord 新增场地记录属性
func (this *Venue) NewVenueRecord(ctx *gin.Context) {
	var data VenueRecordNew
	if err := this.ShouldBindJSON(ctx, &data); err != nil {
		this.Error(ctx, "请求数据不正确")
		return
	}

	if data.TemplateID == 0 {
		this.Error(ctx, "模版名称不能为空")
		return
	}

	recordId, err := createVenueRecord(data.TemplateID)
	if err != nil {
		println(err.Error())
		this.Error(ctx, "添加记录失败！")
		return
	}

	err = createVenueRecordAttributes(recordId, data.Fields)
	if err != nil {
		this.Error(ctx, "添加记录属性失败！")
		return
	}

	this.Success(ctx, "添加模版成功！")
}

// DeleteVenueRecord 删除场地记录属性
func (this *Venue) DeleteVenueRecord(ctx *gin.Context) {
	recordId := ctx.Param("recordID")
	if recordId == "" {
		this.Error(ctx, "记录ID不能为空")
		return
	}

	// 根据记录ID删除对应的删除记录数据
	if err := deleteVenueRecordAttributesByRecordId(recordId); err != nil {
		this.Error(ctx, "删除记录数据失败")
	}

	// 根据记录ID删除对应的删除记录
	if err := deleteVenueRecord(recordId); err != nil {
		this.Error(ctx, "删除记录失败")
	}

	this.Success(ctx, "删除记录成功！")
}

// UpdateVenueRecordAttributes 更新场地记录属性
func (this *Venue) UpdateVenueRecordAttributes(ctx *gin.Context) {
	var data VenueRecordUpdate
	if err := this.ShouldBindJSON(ctx, &data); err != nil {
		this.Error(ctx, "请求数据不正确")
		return
	}

	if data.RecordID == 0 {
		this.Error(ctx, "记录ID不能为空")
		return
	}

	for _, field := range data.Fields {
		if err := updateVenueRecordAttributes(data.RecordID, field.ID, field.Value); err != nil {
			this.Error(ctx, "更新记录失败")
		}
	}

	this.Success(ctx, "更新记录成功！")
}

// 数据库操作

// createVenueTemplate 创建新的场地模板
// 联表查询
func createVenueTemplate(templateName string) (uint, error) {
	var maxID uint
	// 查找数据库中最大的模版ID
	err := model.NewVenueTemplate().Select("MAX(id)").Scan(&maxID).Error
	if err != nil {
		return 0, err
	}

	newTemplate := model.VenueTemplates{
		ID:           maxID + 1, // 设置新模版的ID
		TemplateName: templateName,
	}

	if err := model.NewVenueTemplate().Create(&newTemplate).Error; err != nil {
		return 0, err
	}

	return newTemplate.ID, nil
}

// updateTemplateName 更新模版名称
func updateTemplateName(id uint, newName string) error {
	return model.NewVenueTemplate().Where("id = ?", id).
		Updates(map[string]interface{}{"template_name": newName}).Error
}

// deleteTemplateByID 根据 ID 删除模板
func deleteTemplateByID(id string) error {
	var template model.VenueTemplates
	err := model.NewVenueTemplate().Where("id = ?", id).First(&template).Error
	if err != nil {
		return model.ErrNotFound // 返回自定义错误，表示模板未找到
	}

	err = model.NewVenueTemplate().Delete(&template).Error
	if err != nil {
		return err // 返回删除过程中发生的错误
	}

	return nil // 删除成功
}

// createTemplateField 创建模版字段
func createTemplateField(templateID uint, fieldValue string) error {
	newField := model.TemplateFields{
		TemplateID: templateID,
		FieldName:  fieldValue,
		FieldType:  "VARCHAR",
		FieldOrder: 1,
	}
	return model.NewTemplateFields().Create(&newField).Error
}

// getTemplateField 根据模版ID获取模版字段
func getTemplateField(templateID uint) ([]fieldInfo, error) {
	var fields []fieldInfo
	// 执行查询，只获取 FieldName 字段
	err := model.NewTemplateFields().Select("id, field_name").Where("template_id = ?", templateID).Find(&fields).Error
	if err != nil {
		return nil, err
	}

	return fields, nil
}

// updateTemplateField 更新模版字段
func updateTemplateField(fieldID uint, fieldValue string) error {
	return model.NewTemplateFields().Where("id = ?", fieldID).
		Updates(map[string]interface{}{
			"field_name": fieldValue,
		}).Error
}

// deleteTemplateField 删除模版字段
func deleteTemplateField(fieldID uint) error {
	return model.NewTemplateFields().Where("id = ?", fieldID).Delete(&model.TemplateFields{}).Error
}

// createVenueRecord 插入场地记录的函数
func createVenueRecord(templateID uint) (uint, error) {
	venueRecord := model.VenueRecords{
		TemplateID: templateID,
	}
	if err := model.NewVenueRecord().Create(&venueRecord).Error; err != nil {
		return 0, err // 返回插入错误
	}
	return venueRecord.ID, nil // 返回新插入记录的 ID
}

// createVenueRecordAttributes 插入场地记录属性的函数
func createVenueRecordAttributes(recordID uint, fields []FieldsRecord) error {
	for _, field := range fields {
		venueRecordAttribute := model.VenueRecordAttribute{
			RecordID:   recordID,
			FieldName:  field.FieldName,
			FieldValue: field.Value,
			FieldID:    field.ID,
		}
		if err := model.NewVenueRecordAttribute().Create(&venueRecordAttribute).Error; err != nil {
			return err // 返回插入错误
		}
	}
	return nil
}

// deleteVenueRecord 删除记录
func deleteVenueRecord(recordID string) error {
	return model.NewVenueRecord().Where("id = ?", recordID).Delete(&model.VenueRecords{}).Error
}

// deleteVenueRecordAttributesByRecordId 删除记录数据
func deleteVenueRecordAttributesByRecordId(recordID string) error {
	return model.NewVenueRecordAttribute().Where("record_id = ?", recordID).Delete(&model.VenueRecordAttribute{}).Error
}

// updateVenueRecordAttributes 更新记录数据
func updateVenueRecordAttributes(recordID, fieldId uint, value string) error {
	return model.NewVenueRecordAttribute().Where("record_id = ? AND field_Id = ?", recordID, fieldId).
		Update("field_value", value).Error
}
