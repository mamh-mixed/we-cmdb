package view

import (
	"fmt"

	"github.com/WeBankPartners/we-cmdb/cmdb-server/api/middleware"
	"github.com/WeBankPartners/we-cmdb/cmdb-server/models"
	"github.com/WeBankPartners/we-cmdb/cmdb-server/services/db"
	"github.com/gin-gonic/gin"
)

func GetViewMessage(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		middleware.ReturnParamEmptyError(c, "id")
		return
	}
	result, err := db.GetViewMessage(id)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	middleware.ReturnData(c, result)
}

func ViewCreate(c *gin.Context) {
	var param []*models.UpdateViewParam
	if err := c.ShouldBindJSON(&param); err != nil {
		middleware.ReturnParamValidateError(c, err)
		return
	}
	err := db.ViewCreate(param, middleware.GetRequestUser(c))
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	middleware.ReturnSuccess(c)
}

func ViewUpdate(c *gin.Context) {
	var param []*models.UpdateViewParam
	if err := c.ShouldBindJSON(&param); err != nil {
		middleware.ReturnParamValidateError(c, err)
		return
	}
	var viewIds []string
	for _, view := range param {
		viewIds = append(viewIds, view.Id)
	}
	if len(viewIds) == 0 {
		middleware.ReturnSuccess(c)
		return
	}
	permissionLegal, err := db.ValidateViewPermission(viewIds, middleware.GetRequestRoles(c))
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	if !permissionLegal {
		middleware.ReturnDataPermissionDenyError(c)
		return
	}
	err = db.ViewUpdate(param, middleware.GetRequestUser(c))
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	middleware.ReturnData(c, param[0])
}

func ViewDelete(c *gin.Context) {
	var param []string
	if err := c.ShouldBindJSON(&param); err != nil {
		middleware.ReturnParamValidateError(c, err)
		return
	}
	if len(param) == 0 {
		middleware.ReturnSuccess(c)
		return
	}
	permissionLegal, err := db.ValidateViewPermission(param, middleware.GetRequestRoles(c))
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	if !permissionLegal {
		middleware.ReturnDataPermissionDenyError(c)
		return
	}
	err = db.ViewDelete(param)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	middleware.ReturnSuccess(c)
}

func GraphList(c *gin.Context) {
	var param models.QueryRequestParam
	if err := c.ShouldBindJSON(&param); err != nil {
		middleware.ReturnParamValidateError(c, err)
		return
	}
	pageInfo, rowData, err := db.GraphList(&param)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	middleware.ReturnPageData(c, pageInfo, rowData)
}

func GraphCreate(c *gin.Context) {
	var param []*models.SysGraphTable
	if err := c.ShouldBindJSON(&param); err != nil {
		middleware.ReturnParamValidateError(c, err)
		return
	}
	err := db.GraphCreate(param)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	middleware.ReturnSuccess(c)
}

func GraphUpdate(c *gin.Context) {
	var param []*models.SysGraphTable
	if err := c.ShouldBindJSON(&param); err != nil {
		middleware.ReturnParamValidateError(c, err)
		return
	}
	var graphIds []string
	for _, graph := range param {
		graphIds = append(graphIds, graph.Id)
	}
	if len(graphIds) == 0 {
		middleware.ReturnSuccess(c)
		return
	}
	permissionLegal, err := db.ValidateGraphPermission(graphIds, middleware.GetRequestRoles(c))
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	if !permissionLegal {
		middleware.ReturnDataPermissionDenyError(c)
		return
	}
	err = db.GraphUpdate(param)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	middleware.ReturnSuccess(c)
}

func GraphDelete(c *gin.Context) {
	var param []string
	if err := c.ShouldBindJSON(&param); err != nil {
		middleware.ReturnParamValidateError(c, err)
		return
	}
	if len(param) == 0 {
		middleware.ReturnSuccess(c)
		return
	}
	permissionLegal, err := db.ValidateGraphPermission(param, middleware.GetRequestRoles(c))
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	if !permissionLegal {
		middleware.ReturnDataPermissionDenyError(c)
		return
	}
	err = db.GraphDelete(param)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	middleware.ReturnSuccess(c)
}

func ElementList(c *gin.Context) {
	var param models.QueryRequestParam
	if err := c.ShouldBindJSON(&param); err != nil {
		middleware.ReturnParamValidateError(c, err)
		return
	}
	pageInfo, rowData, err := db.ElementList(&param)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	middleware.ReturnPageData(c, pageInfo, rowData)
}

func ElementCreate(c *gin.Context) {
	var param []*models.SysGraphElementTable
	if err := c.ShouldBindJSON(&param); err != nil {
		middleware.ReturnParamValidateError(c, err)
		return
	}
	err := db.ElementCreate(param)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	middleware.ReturnSuccess(c)
}

func ElementUpdate(c *gin.Context) {
	var param []*models.SysGraphElementTable
	if err := c.ShouldBindJSON(&param); err != nil {
		middleware.ReturnParamValidateError(c, err)
		return
	}
	var elementIds []string
	for _, element := range param {
		elementIds = append(elementIds, element.Id)
	}
	if len(elementIds) == 0 {
		middleware.ReturnSuccess(c)
		return
	}
	permissionLegal, err := db.ValidateElementPermission(elementIds, middleware.GetRequestRoles(c))
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	if !permissionLegal {
		middleware.ReturnDataPermissionDenyError(c)
		return
	}
	err = db.ElementUpdate(param)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	middleware.ReturnSuccess(c)
}

func ElementDelete(c *gin.Context) {
	var param []string
	if err := c.ShouldBindJSON(&param); err != nil {
		middleware.ReturnParamValidateError(c, err)
		return
	}
	if len(param) == 0 {
		middleware.ReturnSuccess(c)
		return
	}
	permissionLegal, err := db.ValidateElementPermission(param, middleware.GetRequestRoles(c))
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	if !permissionLegal {
		middleware.ReturnDataPermissionDenyError(c)
		return
	}
	err = db.ElementDelete(param)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	middleware.ReturnSuccess(c)
}

func GetElementRefAttr(c *gin.Context) {
	parentElement := c.Query("parent_element")
	reportObject := c.Query("report_object")
	if parentElement == "" || reportObject == "" {
		middleware.ReturnParamValidateError(c, fmt.Errorf("Param can not empty "))
		return
	}
	result, err := db.GetElementRefAttrs(parentElement, reportObject)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	middleware.ReturnData(c, result)
}
