package db

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/WeBankPartners/go-common-lib/guid"
	"github.com/WeBankPartners/we-cmdb/cmdb-server/models"
)

func CreateCiTemplate(inputParam *models.SysCiTemplateTable) error {
	if err := validateCiTemplateInput(inputParam); err != nil {
		return err
	}
	var actions []*execAction
	if inputParam.ImageFile != "" {
		actions = append(actions, &execAction{Sql: "insert into sys_ci_template(id,description,image_file,state_machine) value (?,?,?,?)", Param: []interface{}{inputParam.Id, inputParam.Description, inputParam.ImageFile, inputParam.StateMachine}})
	} else {
		actions = append(actions, &execAction{Sql: "insert into sys_ci_template(id,description,state_machine) value (?,?,?)", Param: []interface{}{inputParam.Id, inputParam.Description, inputParam.StateMachine}})
	}
	return transaction(actions)
}

func UpdateCiTemplate(inputParam *models.SysCiTemplateTable) error {
	var ciTemplateTable []*models.SysCiTemplateTable
	err := x.SQL("select * from sys_ci_template where id=?", inputParam.Id).Find(&ciTemplateTable)
	if err != nil {
		return fmt.Errorf("Try to get ci template data fail,%s ", err.Error())
	}
	if len(ciTemplateTable) == 0 {
		return fmt.Errorf("Can not find ci template with id:%s ", inputParam.Id)
	}
	if err = validateCiTemplateInput(inputParam); err != nil {
		return err
	}
	var actions []*execAction
	if inputParam.ImageFile != "" {
		actions = append(actions, &execAction{Sql: "update sys_ci_template set description=?,state_machine=?,image_file=? where id=?", Param: []interface{}{inputParam.Description, inputParam.StateMachine, inputParam.ImageFile, inputParam.Id}})
	} else {
		actions = append(actions, &execAction{Sql: "update sys_ci_template set description=?,state_machine=?,image_file=NULL where id=?", Param: []interface{}{inputParam.Description, inputParam.StateMachine, inputParam.Id}})
	}
	err = transaction(actions)
	if err != nil {
		return fmt.Errorf("Update database fail,%s ", err.Error())
	}
	if inputParam.ImageFile != "" {
		oldImageFileName, _ := getImageFileName(ciTemplateTable[0].ImageFile)
		CiTypesImageDelete(ciTemplateTable[0].ImageFile, oldImageFileName)
	}
	return nil
}

func validateCiTemplateInput(param *models.SysCiTemplateTable) error {
	if param.Id == "" {
		return fmt.Errorf("Param id can not empty ")
	}
	if param.StateMachine == "" {
		return fmt.Errorf("Param stateMachine can not empty ")
	}
	_, stateMachineList, queryMachineErr := ListStateMachine(&models.QueryRequestParam{Filters: []*models.QueryRequestFilterObj{{Name: "id", Operator: "eq", Value: param.StateMachine}}})
	if queryMachineErr != nil {
		return fmt.Errorf("Try to query depend state machine fail,%s ", queryMachineErr.Error())
	}
	if len(stateMachineList) == 0 {
		return fmt.Errorf("Can not find stateMachine:%s ", param.StateMachine)
	}
	return nil
}

func DeleteCiTemplate(ids []string) error {
	if len(ids) == 0 {
		return nil
	}
	var actions []*execAction
	for _, id := range ids {
		actions = append(actions, &execAction{Sql: "delete from sys_ci_template where id=?", Param: []interface{}{id}})
		actions = append(actions, &execAction{Sql: "delete from sys_ci_template_attr where ci_template=?", Param: []interface{}{id}})
		actions = append(actions, &execAction{Sql: "update sys_ci_type set ci_template=NULL where ci_template=?", Param: []interface{}{id}})
	}
	return transactionWithoutForeignCheck(actions)
}

func ListCiAttrTemplate(param *models.QueryRequestParam) (pageInfo models.PageInfo, rowData []*models.SysCiTemplateAttrTable, err error) {
	rowData = []*models.SysCiTemplateAttrTable{}
	filterSql, queryColumn, queryParam := transFiltersToSQL(param, &models.TransFiltersParam{IsStruct: true, StructObj: models.SysCiTemplateAttrTable{}, PrimaryKey: "id"})
	baseSql := fmt.Sprintf("SELECT %s FROM sys_ci_template_attr WHERE 1=1 %s ", queryColumn, filterSql)
	if param.Paging {
		pageInfo.StartIndex = param.Pageable.StartIndex
		pageInfo.PageSize = param.Pageable.PageSize
		pageInfo.TotalRows = queryCount(baseSql, queryParam...)
		pageSql, pageParam := transPageInfoToSQL(*param.Pageable)
		baseSql += pageSql
		queryParam = append(queryParam, pageParam...)
	}
	err = x.SQL(baseSql, queryParam...).Find(&rowData)
	return
}

func CreateCiAttrTemplate(param []*models.SysCiTemplateAttrTable) (result []string, err error) {
	if len(param) == 0 {
		return
	}
	var actions []*execAction
	for _, inputParam := range param {
		if err = validateCiTemplateAttrParam(inputParam); err != nil {
			break
		}
		inputParam.Id = inputParam.CiTemplate + models.SysTableIdConnector + inputParam.Name
		result = append(result, inputParam.Id)
		inputParam.Status = "notCreated"
		inputParam.Source = "template"
		inputParam.AutofillAble = "no"
		inputParam.EditGroupControl = "no"
		if strings.Contains(inputParam.DataType, "(") {
			tmpDataType := inputParam.DataType[:strings.Index(inputParam.DataType, "(")]
			inputParam.DataLength, _ = strconv.Atoi(inputParam.DataType[strings.Index(inputParam.DataType, "(")+1 : len(inputParam.DataType)-1])
			inputParam.DataType = tmpDataType
		}
		tmpAction := execAction{Sql: "insert into sys_ci_template_attr(id,ci_template,name,display_name,description,status,input_type,data_type,data_length,text_validate,ui_search_order,ui_form_order,unique_constraint,ui_nullable,nullable,editable,display_by_default,permission_usage,reset_on_edit,source,customizable,autofillable,edit_group_control,confirm_nullable) value (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"}
		tmpAction.Param = []interface{}{inputParam.Id, inputParam.CiTemplate, inputParam.Name, inputParam.DisplayName, inputParam.Description, inputParam.Status, inputParam.InputType, inputParam.DataType, inputParam.DataLength, inputParam.TextValidate, inputParam.UiSearchOrder, inputParam.UiFormOrder, inputParam.UniqueConstraint, inputParam.UiNullable, inputParam.Nullable, inputParam.Editable, inputParam.DisplayByDefault, inputParam.PermissionUsage, inputParam.ResetOnEdit, inputParam.Source, inputParam.Customizable, inputParam.AutofillAble, inputParam.EditGroupControl, inputParam.ConfirmNullable}
		actions = append(actions, &tmpAction)
	}
	if err != nil {
		return result, err
	}
	return result, transaction(actions)
}

func UpdateCiAttrTemplate(param []*models.SysCiTemplateAttrTable) error {
	if len(param) == 0 {
		return nil
	}
	var actions []*execAction
	for _, inputParam := range param {
		if strings.Contains(inputParam.DataType, "(") {
			tmpDataType := inputParam.DataType[:strings.Index(inputParam.DataType, "(")]
			inputParam.DataLength, _ = strconv.Atoi(inputParam.DataType[strings.Index(inputParam.DataType, "(")+1 : len(inputParam.DataType)-1])
			inputParam.DataType = tmpDataType
		}
		tmpAction := execAction{Sql: "update sys_ci_template_attr set display_name=?,description=?,input_type=?,data_type=?,data_length=?,text_validate=?,ui_search_order=?,ui_form_order=?,unique_constraint=?,ui_nullable=?,nullable=?,editable=?,display_by_default=?,permission_usage=?,reset_on_edit=?,source=?,customizable=?,confirm_nullable=? where id=?"}
		tmpAction.Param = []interface{}{inputParam.DisplayName, inputParam.Description, inputParam.InputType, inputParam.DataType, inputParam.DataLength, inputParam.TextValidate, inputParam.UiSearchOrder, inputParam.UiFormOrder, inputParam.UniqueConstraint, inputParam.UiNullable, inputParam.Nullable, inputParam.Editable, inputParam.DisplayByDefault, inputParam.PermissionUsage, inputParam.ResetOnEdit, inputParam.Source, inputParam.Customizable, inputParam.ConfirmNullable, inputParam.Id}
		actions = append(actions, &tmpAction)
	}
	return transaction(actions)
}

func validateCiTemplateAttrParam(param *models.SysCiTemplateAttrTable) error {
	if param.Name == "" {
		return fmt.Errorf("Param name can not empty ")
	}
	if param.CiTemplate == "" {
		return fmt.Errorf("Param ciTemplate can not empty ")
	}
	return nil
}

func DeleteCiAttrTemplate(ids []string) error {
	if len(ids) == 0 {
		return nil
	}
	var actions []*execAction
	for _, id := range ids {
		actions = append(actions, &execAction{Sql: "delete from sys_ci_template_attr where id=?", Param: []interface{}{id}})
	}
	return transaction(actions)
}

func ListStateMachine(param *models.QueryRequestParam) (pageInfo models.PageInfo, rowData []*models.SysStateMachineTable, err error) {
	rowData = []*models.SysStateMachineTable{}
	filterSql, queryColumn, queryParam := transFiltersToSQL(param, &models.TransFiltersParam{IsStruct: true, StructObj: models.SysStateMachineTable{}, PrimaryKey: "id"})
	baseSql := fmt.Sprintf("SELECT %s FROM sys_state_machine WHERE 1=1 %s ", queryColumn, filterSql)
	if param.Paging {
		pageInfo.StartIndex = param.Pageable.StartIndex
		pageInfo.PageSize = param.Pageable.PageSize
		pageInfo.TotalRows = queryCount(baseSql, queryParam...)
		pageSql, pageParam := transPageInfoToSQL(*param.Pageable)
		baseSql += pageSql
		queryParam = append(queryParam, pageParam...)
	}
	err = x.SQL(baseSql, queryParam...).Find(&rowData)
	return
}

func GetStateMachineDetail(id string) (result models.GetStateMachineList, err error) {
	_, stateMachineList, queryMachineErr := ListStateMachine(&models.QueryRequestParam{Filters: []*models.QueryRequestFilterObj{{Name: "id", Operator: "eq", Value: id}}})
	if queryMachineErr != nil {
		err = fmt.Errorf("Try to query state machine table fail,%s ", queryMachineErr.Error())
		return
	}
	if len(stateMachineList) == 0 {
		err = fmt.Errorf("Can not find stata machine with id=%s ", id)
		return
	}
	result.Id = stateMachineList[0].Id
	result.Description = stateMachineList[0].Description
	result.StartState = stateMachineList[0].StartState
	result.FinalState = stateMachineList[0].FinalState
	_, stateList, queryStateErr := ListState(&models.QueryRequestParam{Filters: []*models.QueryRequestFilterObj{{Name: "stateMachine", Operator: "eq", Value: result.Id}}})
	if queryStateErr != nil {
		err = fmt.Errorf("Try to query state table fail,%s ", queryStateErr.Error())
		return
	}
	result.States = stateList
	_, transitionList, queryTransErr := ListStateTransition(&models.QueryRequestParam{Filters: []*models.QueryRequestFilterObj{{Name: "stateMachine", Operator: "eq", Value: result.Id}}})
	if queryTransErr != nil {
		err = fmt.Errorf("Try to query transition table fail,%s ", queryTransErr.Error())
		return
	}
	result.Transitions = transitionList
	return
}

func CreateStateMachine(param []*models.SysStateMachineTable) error {
	if len(param) == 0 {
		return nil
	}
	var actions []*execAction
	for _, inputParam := range param {
		startStateId := fmt.Sprintf("%s_null_0", inputParam.Id)
		finalStateId := fmt.Sprintf("%s_null_1", inputParam.Id)
		actions = append(actions, &execAction{Sql: "insert into sys_state_machine(id,description,start_state,final_state) value (?,?,?,?)", Param: []interface{}{inputParam.Id, inputParam.Description, startStateId, finalStateId}})
		actions = append(actions, &execAction{Sql: "insert into sys_state(id,name,description,state_machine,unique_path_trigger,is_confirm) values (?,'null_0','初始态',?,'no','no'),(?,'null_1','终止态',?,'no','yes')", Param: []interface{}{startStateId, inputParam.Id, finalStateId, inputParam.Id}})
	}
	return transactionWithoutForeignCheck(actions)
}

func UpdateStateMachine(param []*models.SysStateMachineTable) error {
	if len(param) == 0 {
		return nil
	}
	var actions []*execAction
	for _, inputParam := range param {
		actions = append(actions, &execAction{Sql: "update sys_state_machine set description=? where id=?", Param: []interface{}{inputParam.Description, inputParam.Id}})
	}
	return transaction(actions)
}

func DeleteStateMachine(ids []string) error {
	if len(ids) == 0 {
		return nil
	}
	querySql, queryParams := createListParams(ids, "")
	queryParams = append(queryParams, queryParams...)
	queryParams = append([]interface{}{"select id,'ci' as name from sys_ci_type where state_machine in (" + querySql + ") union select id,'ci_template' as name from sys_ci_template where state_machine in (" + querySql + ")"}, queryParams...)
	queryDep, queryErr := x.QueryString(queryParams...)
	if queryErr != nil {
		return fmt.Errorf("Try to query depend ci fail,%s ", queryErr.Error())
	}
	if len(queryDep) > 0 {
		return fmt.Errorf("Find dependence %s id:%s ", queryDep[0]["name"], queryDep[0]["id"])
	}
	var actions []*execAction
	for _, id := range ids {
		actions = append(actions, &execAction{Sql: "delete from sys_state_machine where id=?", Param: []interface{}{id}})
		actions = append(actions, &execAction{Sql: "delete from sys_state where state_machine=?", Param: []interface{}{id}})
		actions = append(actions, &execAction{Sql: "delete from sys_state_transition where state_machine=?", Param: []interface{}{id}})
	}
	return transactionWithoutForeignCheck(actions)
}

func ListState(param *models.QueryRequestParam) (pageInfo models.PageInfo, rowData []*models.SysStateTable, err error) {
	rowData = []*models.SysStateTable{}
	filterSql, queryColumn, queryParam := transFiltersToSQL(param, &models.TransFiltersParam{IsStruct: true, StructObj: models.SysStateTable{}, PrimaryKey: "id"})
	baseSql := fmt.Sprintf("SELECT %s FROM sys_state WHERE 1=1 %s ", queryColumn, filterSql)
	if param.Paging {
		pageInfo.StartIndex = param.Pageable.StartIndex
		pageInfo.PageSize = param.Pageable.PageSize
		pageInfo.TotalRows = queryCount(baseSql, queryParam...)
		pageSql, pageParam := transPageInfoToSQL(*param.Pageable)
		baseSql += pageSql
		queryParam = append(queryParam, pageParam...)
	}
	err = x.SQL(baseSql, queryParam...).Find(&rowData)
	return
}

func CreateState(param []*models.SysStateTable) error {
	if len(param) == 0 {
		return nil
	}
	var err error
	var actions []*execAction
	for _, inputParam := range param {
		if inputParam.StateMachine == "" {
			err = fmt.Errorf("StateMachine can not empty ")
			break
		}
		if !strings.HasSuffix(inputParam.Name, "_0") && !strings.HasSuffix(inputParam.Name, "_1") {
			err = fmt.Errorf("Name:%s is illegal,must end with _0 or _1 ", inputParam.Name)
			break
		}
		if inputParam.UniquePathTrigger == "" {
			inputParam.UniquePathTrigger = "no"
		}
		if inputParam.IsConfirm == "" {
			inputParam.IsConfirm = "no"
		}
		tmpAction := execAction{Sql: "insert into sys_state(id,name,description,state_machine,unique_path_trigger,is_confirm) value (?,?,?,?,?,?)"}
		tmpAction.Param = []interface{}{inputParam.StateMachine + models.SysTableIdConnector + inputParam.Name, inputParam.Name, inputParam.Description, inputParam.StateMachine, inputParam.UniquePathTrigger, inputParam.IsConfirm}
		actions = append(actions, &tmpAction)
	}
	if err != nil {
		return err
	}
	return transaction(actions)
}

func UpdateState(param []*models.SysStateTable) error {
	if len(param) == 0 {
		return nil
	}
	var err error
	var actions []*execAction
	for _, inputParam := range param {
		if inputParam.Id == "" {
			err = fmt.Errorf("Param id can not empty ")
			break
		}
		if strings.HasSuffix(inputParam.Id, "null_0") || strings.HasSuffix(inputParam.Id, "null_1") {
			continue
		}
		if inputParam.UniquePathTrigger == "" {
			inputParam.UniquePathTrigger = "no"
		}
		if inputParam.IsConfirm == "" {
			inputParam.IsConfirm = "no"
		}
		actions = append(actions, &execAction{Sql: "update sys_state set description=?,unique_path_trigger=?,is_confirm=? where id=?", Param: []interface{}{inputParam.Description, inputParam.UniquePathTrigger, inputParam.IsConfirm, inputParam.Id}})
	}
	if err != nil {
		return err
	}
	return transaction(actions)
}

func DeleteState(ids []string) error {
	if len(ids) == 0 {
		return nil
	}
	var transitionTable []*models.SysStateTransitionTable
	querySql, queryParams := createListParams(ids, "")
	queryParams = append(queryParams, queryParams...)
	err := x.SQL("select guid,state_machine,operation_en from sys_state_transition where current_state in ("+querySql+") or target_state in ("+querySql+")", queryParams...).Find(&transitionTable)
	if err != nil {
		return fmt.Errorf("Try to query depend transitions fail,%s ", err.Error())
	}
	if len(transitionTable) > 0 {
		return fmt.Errorf("Find dependence with transtion: %s %s ", transitionTable[0].StateMachine, transitionTable[0].OperationEn)
	}
	var actions []*execAction
	for _, id := range ids {
		actions = append(actions, &execAction{Sql: "delete from sys_state where id=?", Param: []interface{}{id}})
	}
	return transaction(actions)
}

func ListStateTransition(param *models.QueryRequestParam) (pageInfo models.PageInfo, rowData []*models.SysStateTransitionTable, err error) {
	rowData = []*models.SysStateTransitionTable{}
	filterSql, queryColumn, queryParam := transFiltersToSQL(param, &models.TransFiltersParam{IsStruct: true, StructObj: models.SysStateTransitionTable{}, PrimaryKey: "guid"})
	baseSql := fmt.Sprintf("SELECT %s FROM sys_state_transition WHERE 1=1 %s ", queryColumn, filterSql)
	if param.Paging {
		pageInfo.StartIndex = param.Pageable.StartIndex
		pageInfo.PageSize = param.Pageable.PageSize
		pageInfo.TotalRows = queryCount(baseSql, queryParam...)
		pageSql, pageParam := transPageInfoToSQL(*param.Pageable)
		baseSql += pageSql
		queryParam = append(queryParam, pageParam...)
	}
	err = x.SQL(baseSql, queryParam...).Find(&rowData)
	return
}

func validateStateTransitionParam(param *models.SysStateTransitionTable) error {
	if param.CurrentState == "" {
		return fmt.Errorf("Param currentState can not empty ")
	}
	if param.TargetState == "" {
		return fmt.Errorf("Param targetState can not empty ")
	}
	if param.Operation == "" {
		return fmt.Errorf("Param operation can not emtpy ")
	}
	if param.OperationEn == "" {
		return fmt.Errorf("Param operation_en can not empty ")
	}
	if param.Permission != "insert" && param.Permission != "update" && param.Permission != "delete" && param.Permission != "execute" {
		return fmt.Errorf("Param permission is illegal ")
	}
	if param.Action != "insert" && param.Action != "update" && param.Action != "delete" && param.Action != "confirm" && param.Action != "callback" && param.Action != "execute" {
		return fmt.Errorf("Param action is illegal ")
	}
	if param.OperationFormType != "editable_form" && param.OperationFormType != "select_form" && param.OperationFormType != "confirm_form" {
		return fmt.Errorf("Param operationFormType is illegal ")
	}
	if param.OperationMultiple != "yes" && param.OperationMultiple != "no" {
		return fmt.Errorf("Param operationMultiple is illegal ")
	}
	if param.Action == "confirm" {
		if strings.HasSuffix(param.CurrentState, "1") {
			return fmt.Errorf("Confirm action should start with state 0 ")
		}
		if strings.HasSuffix(param.TargetState, "0") {
			return fmt.Errorf("Confirm action should end with state 1 ")
		}
	}
	if strings.HasSuffix(param.CurrentState, "null_0") {
		if param.Action != "insert" {
			return fmt.Errorf("CurrentState:%s must with insert action ", param.CurrentState)
		}
	}
	if strings.HasSuffix(param.TargetState, "null_1") {
		if param.Action != "delete" {
			return fmt.Errorf("TargetState:%s must with delete action ", param.TargetState)
		}
	}
	return nil
}

func CreateStateTransition(param []*models.SysStateTransitionTable) error {
	if len(param) == 0 {
		return nil
	}
	var err error
	var actions []*execAction
	guidList := guid.CreateGuidList(len(param))
	for i, inputParam := range param {
		if inputParam.StateMachine == "" {
			err = fmt.Errorf("StateMachine can not empty ")
			break
		}
		if err = validateStateTransitionParam(inputParam); err != nil {
			break
		}
		tmpAction := execAction{Sql: "insert into sys_state_transition(guid,state_machine,current_state,target_state,operation,operation_en,permission,`action`,operation_form_type,operation_multiple) value (?,?,?,?,?,?,?,?,?,?)"}
		tmpAction.Param = []interface{}{guidList[i], inputParam.StateMachine, inputParam.CurrentState, inputParam.TargetState, inputParam.Operation, inputParam.OperationEn, inputParam.Permission, inputParam.Action, inputParam.OperationFormType, inputParam.OperationMultiple}
		actions = append(actions, &tmpAction)
	}
	if err != nil {
		return err
	}
	return transaction(actions)
}

func UpdateStateTransition(param []*models.SysStateTransitionTable) error {
	if len(param) == 0 {
		return nil
	}
	var err error
	var actions []*execAction
	for _, inputParam := range param {
		if inputParam.Guid == "" {
			err = fmt.Errorf("Guid can not empty ")
			break
		}
		if err = validateStateTransitionParam(inputParam); err != nil {
			break
		}
		tmpAction := execAction{Sql: "update sys_state_transition set state_machine=?,current_state=?,target_state=?,operation=?,operation_en=?,permission=?,`action`=?,operation_form_type=?,operation_multiple=? where guid=?"}
		tmpAction.Param = []interface{}{inputParam.StateMachine, inputParam.CurrentState, inputParam.TargetState, inputParam.Operation, inputParam.OperationEn, inputParam.Permission, inputParam.Action, inputParam.OperationFormType, inputParam.OperationMultiple, inputParam.Guid}
		actions = append(actions, &tmpAction)
	}
	if err != nil {
		return err
	}
	return transaction(actions)
}

func DeleteStateTransition(ids []string) error {
	if len(ids) == 0 {
		return nil
	}
	var actions []*execAction
	for _, id := range ids {
		actions = append(actions, &execAction{Sql: "delete from sys_state_transition where guid=?", Param: []interface{}{id}})
	}
	return transaction(actions)
}

func ImportStateMachine(param models.GetStateMachineList) (result models.ImportStateMachineResult, err error) {
	result = models.ImportStateMachineResult{DiffFlag: false, StateMachine: &models.SysStateMachineTable{}}
	_, stateMachineList, machineQueryErr := ListStateMachine(&models.QueryRequestParam{Filters: []*models.QueryRequestFilterObj{{Name: "id", Operator: "eq", Value: param.Id}}})
	if machineQueryErr != nil {
		err = fmt.Errorf("Try to query state machine fail,%s ", machineQueryErr.Error())
		return
	}
	var actions []*execAction
	if len(stateMachineList) > 0 {
		if stateMachineList[0].Description != param.Description {
			result.DiffFlag = true
			result.StateMachine = stateMachineList[0]
		}
		_, stateList, queryStateErr := ListState(&models.QueryRequestParam{Filters: []*models.QueryRequestFilterObj{{Name: "stateMachine", Operator: "eq", Value: param.Id}}})
		if queryStateErr != nil {
			err = fmt.Errorf("Try to query state table fail,%s ", queryStateErr.Error())
			return
		}
		for _, newState := range param.States {
			existFlag := false
			for _, oldState := range stateList {
				if newState.Id == oldState.Id {
					existFlag = true
					if newState.Description != oldState.Description || newState.UniquePathTrigger != oldState.UniquePathTrigger || newState.IsConfirm != oldState.IsConfirm {
						result.DiffFlag = true
						result.NewStates = append(result.NewStates, newState)
						result.OldStates = append(result.OldStates, oldState)
					}
					break
				}
			}
			if !existFlag {
				tmpAction := execAction{Sql: "insert into sys_state(id,name,description,state_machine,unique_path_trigger,is_confirm) value (?,?,?,?,?,?)"}
				tmpAction.Param = []interface{}{newState.Id, newState.Name, newState.Description, newState.StateMachine, newState.UniquePathTrigger, newState.IsConfirm}
				actions = append(actions, &tmpAction)
			}
		}
		_, transitionList, queryTransErr := ListStateTransition(&models.QueryRequestParam{Filters: []*models.QueryRequestFilterObj{{Name: "stateMachine", Operator: "eq", Value: param.Id}}})
		if queryTransErr != nil {
			err = fmt.Errorf("Try to query transtition fail,%s ", queryTransErr.Error())
			return
		}
		for _, newTrans := range param.Transitions {
			existFlag := false
			for _, oldTrans := range transitionList {
				if newTrans.Guid == oldTrans.Guid {
					existFlag = true
					if newTrans.CurrentState != oldTrans.CurrentState || newTrans.TargetState != oldTrans.TargetState || newTrans.Operation != oldTrans.Operation || newTrans.OperationEn != oldTrans.OperationEn || newTrans.Action != oldTrans.Action || newTrans.Permission != oldTrans.Permission || newTrans.OperationFormType != oldTrans.OperationFormType || newTrans.OperationMultiple != oldTrans.OperationMultiple {
						result.DiffFlag = true
						result.NewTransitions = append(result.NewTransitions, newTrans)
						result.OldTransitions = append(result.OldTransitions, oldTrans)
					}
					break
				}
			}
			if !existFlag {
				tmpAction := execAction{Sql: "insert into sys_state_transition(guid,state_machine,current_state,target_state,operation,operation_en,permission,`action`,operation_form_type,operation_multiple) value (?,?,?,?,?,?,?,?,?,?)"}
				tmpAction.Param = []interface{}{newTrans.Guid, newTrans.StateMachine, newTrans.CurrentState, newTrans.TargetState, newTrans.Operation, newTrans.OperationEn, newTrans.Permission, newTrans.Action, newTrans.OperationFormType, newTrans.OperationMultiple}
				actions = append(actions, &tmpAction)
			}
		}
	} else {
		actions = append(actions, &execAction{Sql: "insert into sys_state_machine(id,description,start_state,final_state) value (?,?,?,?)", Param: []interface{}{param.Id, param.Description, param.StartState, param.FinalState}})
		for _, state := range param.States {
			tmpAction := execAction{Sql: "insert into sys_state(id,name,description,state_machine,unique_path_trigger,is_confirm) value (?,?,?,?,?,?)"}
			tmpAction.Param = []interface{}{state.Id, state.Name, state.Description, state.StateMachine, state.UniquePathTrigger, state.IsConfirm}
			actions = append(actions, &tmpAction)
		}
		for _, transition := range param.Transitions {
			tmpAction := execAction{Sql: "insert into sys_state_transition(guid,state_machine,current_state,target_state,operation,operation_en,permission,`action`,operation_form_type,operation_multiple) value (?,?,?,?,?,?,?,?,?,?)"}
			tmpAction.Param = []interface{}{transition.Guid, transition.StateMachine, transition.CurrentState, transition.TargetState, transition.Operation, transition.OperationEn, transition.Permission, transition.Action, transition.OperationFormType, transition.OperationMultiple}
			actions = append(actions, &tmpAction)
		}
	}
	if len(actions) == 0 {
		return
	}
	err = transactionWithoutForeignCheck(actions)
	return
}

func ConfirmStateImport(param models.ImportStateMachineResult) error {
	var actions []*execAction
	if param.StateMachine != nil && param.StateMachine.Id != "" {
		actions = append(actions, &execAction{Sql: "update sys_state_machine set description=? where id=?", Param: []interface{}{param.StateMachine.Description, param.StateMachine.Id}})
	}
	if len(param.NewStates) > 0 {
		for _, state := range param.NewStates {
			actions = append(actions, &execAction{Sql: "update sys_state set description=?,unique_path_trigger=?,is_confirm=? where id=?", Param: []interface{}{state.Description, state.UniquePathTrigger, state.IsConfirm, state.Id}})
		}
	}
	if len(param.NewTransitions) > 0 {
		for _, trans := range param.NewTransitions {
			tmpAction := execAction{Sql: "update sys_state_transition set state_machine=?,current_state=?,target_state=?,operation=?,operation_en=?,permission=?,`action`=?,operation_form_type=?,operation_multiple=? where guid=?"}
			tmpAction.Param = []interface{}{trans.StateMachine, trans.CurrentState, trans.TargetState, trans.Operation, trans.OperationEn, trans.Permission, trans.Action, trans.OperationFormType, trans.OperationMultiple, trans.Guid}
			actions = append(actions, &tmpAction)
		}
	}
	if len(actions) == 0 {
		return nil
	}
	return transactionWithoutForeignCheck(actions)
}
