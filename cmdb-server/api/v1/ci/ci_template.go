package ci

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/WeBankPartners/we-cmdb/cmdb-server/api/middleware"
	"github.com/WeBankPartners/we-cmdb/cmdb-server/models"
	"github.com/WeBankPartners/we-cmdb/cmdb-server/services/db"
	"github.com/gin-gonic/gin"
)

func CreateCiTemplate(c *gin.Context) {
	var err error
	var param models.SysCiTemplateTable
	if err = c.ShouldBindJSON(&param); err != nil {
		middleware.ReturnParamValidateError(c, err)
		return
	}
	var imageGuid, imageFileName string
	if param.FileName != "" {
		imageGuid, imageFileName, err = saveImageFile(param.ImageFile, param.FileName)
		if err != nil {
			middleware.ReturnParamValidateError(c, err)
			return
		}
		param.ImageFile = imageGuid
	}
	err = db.CreateCiTemplate(&param)
	if err != nil {
		if imageGuid != "" {
			db.CiTypesImageDelete(imageGuid, imageFileName)
		}
		middleware.ReturnServerHandleError(c, err)
		return
	}
	middleware.ReturnSuccess(c)
}

func UpdateCiTemplate(c *gin.Context) {
	var err error
	var param models.SysCiTemplateTable
	if err = c.ShouldBindJSON(&param); err != nil {
		middleware.ReturnParamValidateError(c, err)
		return
	}
	var newImageGuid, imageFileName string
	if param.FileName != "" {
		newImageGuid, imageFileName, err = saveImageFile(param.ImageFile, param.FileName)
		if err != nil {
			middleware.ReturnParamValidateError(c, err)
			return
		}
		param.ImageFile = newImageGuid
	}
	err = db.UpdateCiTemplate(&param)
	if err != nil {
		if newImageGuid != "" {
			db.CiTypesImageDelete(newImageGuid, imageFileName)
		}
		middleware.ReturnServerHandleError(c, err)
		return
	}
	middleware.ReturnSuccess(c)
}

func DeleteCiTemplate(c *gin.Context) {
	ids := strings.Split(c.Query("ids"), ",")
	err := db.DeleteCiTemplate(ids)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	middleware.ReturnSuccess(c)
}

func ListCiAttrTemplate(c *gin.Context) {
	var param models.QueryRequestParam
	if err := c.ShouldBindJSON(&param); err != nil {
		middleware.ReturnParamValidateError(c, err)
		return
	}
	pageInfo, rowData, err := db.ListCiAttrTemplate(&param)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	middleware.ReturnPageData(c, pageInfo, rowData)
}

func CreateCiAttrTemplate(c *gin.Context) {
	var param []*models.SysCiTemplateAttrTable
	if err := c.ShouldBindJSON(&param); err != nil {
		middleware.ReturnParamValidateError(c, err)
		return
	}
	result, err := db.CreateCiAttrTemplate(param)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	middleware.ReturnData(c, result)
}

func UpdateCiAttrTemplate(c *gin.Context) {
	var param []*models.SysCiTemplateAttrTable
	if err := c.ShouldBindJSON(&param); err != nil {
		middleware.ReturnParamValidateError(c, err)
		return
	}
	err := db.UpdateCiAttrTemplate(param)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	middleware.ReturnSuccess(c)
}

func DeleteCiAttrTemplate(c *gin.Context) {
	ids := strings.Split(c.Query("ids"), ",")
	err := db.DeleteCiAttrTemplate(ids)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	middleware.ReturnSuccess(c)
}

func ListStateMachine(c *gin.Context) {
	var param models.QueryRequestParam
	if err := c.ShouldBindJSON(&param); err != nil {
		middleware.ReturnParamValidateError(c, err)
		return
	}
	pageInfo, rowData, err := db.ListStateMachine(&param)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	middleware.ReturnPageData(c, pageInfo, rowData)
}

func GetStateMachineDetail(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		middleware.ReturnParamEmptyError(c, "id")
		return
	}
	export := c.Query("export")
	result, err := db.GetStateMachineDetail(id)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	if export == "yes" {
		b, err := json.Marshal(result)
		if err != nil {
			middleware.ReturnServerHandleError(c, fmt.Errorf("Export state machine fail, json marshal object error:%s ", err.Error()))
			return
		}
		c.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s_%s.json", "state_machine", time.Now().Format("20060102150405")))
		c.Data(http.StatusOK, "application/octet-stream", b)
		return
	}
	middleware.ReturnData(c, result)
}

func ImportStateMachine(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ResponseErrorJson{StatusCode: "PARAM_HANDLE_ERROR", StatusMessage: "Http read upload file fail:" + err.Error(), Data: nil})
		return
	}
	f, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ResponseErrorJson{StatusCode: "PARAM_HANDLE_ERROR", StatusMessage: "File open error:" + err.Error(), Data: nil})
		return
	}
	var paramObj models.GetStateMachineList
	b, err := ioutil.ReadAll(f)
	defer f.Close()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ResponseErrorJson{StatusCode: "PARAM_HANDLE_ERROR", StatusMessage: "Read content fail error:" + err.Error(), Data: nil})
		return
	}
	err = json.Unmarshal(b, &paramObj)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ResponseErrorJson{StatusCode: "PARAM_HANDLE_ERROR", StatusMessage: "Json unmarshal fail error:" + err.Error(), Data: nil})
		return
	}
	result, err := db.ImportStateMachine(paramObj)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ResponseErrorJson{StatusCode: "PARAM_HANDLE_ERROR", StatusMessage: err.Error(), Data: nil})
		return
	}
	c.JSON(http.StatusOK, models.ResponseErrorJson{StatusCode: "OK", StatusMessage: "Success", Data: result})
}

func ConfirmStateImport(c *gin.Context) {
	var param models.ImportStateMachineResult
	if err := c.ShouldBindJSON(&param); err != nil {
		middleware.ReturnParamValidateError(c, err)
		return
	}
	err := db.ConfirmStateImport(param)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	middleware.ReturnSuccess(c)
}

func CreateStateMachine(c *gin.Context) {
	var param []*models.SysStateMachineTable
	if err := c.ShouldBindJSON(&param); err != nil {
		middleware.ReturnParamValidateError(c, err)
		return
	}
	err := db.CreateStateMachine(param)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	middleware.ReturnSuccess(c)
}

func UpdateStateMachine(c *gin.Context) {
	var param []*models.SysStateMachineTable
	if err := c.ShouldBindJSON(&param); err != nil {
		middleware.ReturnParamValidateError(c, err)
		return
	}
	err := db.UpdateStateMachine(param)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	middleware.ReturnSuccess(c)
}

func DeleteStateMachine(c *gin.Context) {
	ids := strings.Split(c.Query("ids"), ",")
	err := db.DeleteStateMachine(ids)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	middleware.ReturnSuccess(c)
}

func ListState(c *gin.Context) {
	var param models.QueryRequestParam
	if err := c.ShouldBindJSON(&param); err != nil {
		middleware.ReturnParamValidateError(c, err)
		return
	}
	pageInfo, rowData, err := db.ListState(&param)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	middleware.ReturnPageData(c, pageInfo, rowData)
}

func CreateState(c *gin.Context) {
	var param []*models.SysStateTable
	if err := c.ShouldBindJSON(&param); err != nil {
		middleware.ReturnParamValidateError(c, err)
		return
	}
	err := db.CreateState(param)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	middleware.ReturnSuccess(c)
}

func UpdateState(c *gin.Context) {
	var param []*models.SysStateTable
	if err := c.ShouldBindJSON(&param); err != nil {
		middleware.ReturnParamValidateError(c, err)
		return
	}
	err := db.UpdateState(param)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	middleware.ReturnSuccess(c)
}

func DeleteState(c *gin.Context) {
	ids := strings.Split(c.Query("ids"), ",")
	err := db.DeleteState(ids)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	middleware.ReturnSuccess(c)
}

func ListStateTransition(c *gin.Context) {
	var param models.QueryRequestParam
	if err := c.ShouldBindJSON(&param); err != nil {
		middleware.ReturnParamValidateError(c, err)
		return
	}
	pageInfo, rowData, err := db.ListStateTransition(&param)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	middleware.ReturnPageData(c, pageInfo, rowData)
}

func CreateStateTransition(c *gin.Context) {
	var param []*models.SysStateTransitionTable
	if err := c.ShouldBindJSON(&param); err != nil {
		middleware.ReturnParamValidateError(c, err)
		return
	}
	err := db.CreateStateTransition(param)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	middleware.ReturnSuccess(c)
}

func UpdateStateTransition(c *gin.Context) {
	var param []*models.SysStateTransitionTable
	if err := c.ShouldBindJSON(&param); err != nil {
		middleware.ReturnParamValidateError(c, err)
		return
	}
	err := db.UpdateStateTransition(param)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	middleware.ReturnSuccess(c)
}

func DeleteStateTransition(c *gin.Context) {
	ids := strings.Split(c.Query("ids"), ",")
	err := db.DeleteStateTransition(ids)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	middleware.ReturnSuccess(c)
}
