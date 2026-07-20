package db

import (
	"fmt"
	"strings"
	"time"

	"github.com/WeBankPartners/we-cmdb/cmdb-server/models"
)

func GetViewMessage(id string) (result models.UpdateViewParam, err error) {
	result = models.UpdateViewParam{Id: id}
	var viewTable []*models.SysViewTable
	err = x.SQL("select * from sys_view where id=?", id).Find(&viewTable)
	if err != nil {
		return
	}
	if len(viewTable) == 0 {
		err = fmt.Errorf("Can not find view with id=%s ", id)
		return
	}
	result.Name = viewTable[0].Name
	result.Multiple = viewTable[0].Multiple
	result.Editable = viewTable[0].Editable
	result.SuportVersion = viewTable[0].SuportVersion
	result.FilterAttr = viewTable[0].FilterAttr
	result.FilterValue = viewTable[0].FilterValue
	result.Report = viewTable[0].Report
	result.USE = []string{}
	result.MGMT = []string{}
	var viewRoleTable []*models.SysRoleViewTable
	x.SQL("select * from sys_role_view where `view`=?", id).Find(&viewRoleTable)
	for _, v := range viewRoleTable {
		if v.Permission == "USE" {
			result.USE = append(result.USE, v.Role)
		}
		if v.Permission == "MGMT" {
			result.MGMT = append(result.MGMT, v.Role)
		}
	}
	return
}

func ViewCreate(param []*models.UpdateViewParam, operator string) error {
	if len(param) == 0 {
		return nil
	}
	nowTime := time.Now().Format(models.DateTimeFormat)
	var actions []*execAction
	var err error
	for _, inputParam := range param {
		if len(inputParam.MGMT) == 0 {
			err = fmt.Errorf("Param MGMT can not empty ")
			break
		}
		tmpAction := execAction{Param: []interface{}{}}
		if inputParam.FilterAttr == "" {
			tmpAction = execAction{Sql: "insert into sys_view(id,name,report,editable,suport_version,multiple,create_time,create_user,update_time,update_user) value (?,?,?,?,?,?,?,?,?,?)"}
			tmpAction.Param = append(tmpAction.Param, inputParam.Id, inputParam.Name, inputParam.Report, inputParam.Editable, inputParam.SuportVersion, inputParam.Multiple, nowTime, operator, nowTime, operator)
		} else {
			tmpAction = execAction{Sql: "insert into sys_view(id,name,report,editable,suport_version,multiple,filter_attr,filter_value,create_time,create_user,update_time,update_user) value (?,?,?,?,?,?,?,?,?,?,?,?)"}
			tmpAction.Param = append(tmpAction.Param, inputParam.Id, inputParam.Name, inputParam.Report, inputParam.Editable, inputParam.SuportVersion, inputParam.Multiple, inputParam.FilterAttr, inputParam.FilterValue, nowTime, operator, nowTime, operator)
		}
		actions = append(actions, &tmpAction)
		actions = append(actions, getPermissionActions(inputParam.MGMT, inputParam.USE, inputParam.Id)...)
	}
	if err != nil {
		return err
	}
	return transaction(actions)
}

func ViewUpdate(param []*models.UpdateViewParam, operator string) error {
	if len(param) == 0 {
		return nil
	}
	nowTime := time.Now().Format(models.DateTimeFormat)
	var actions []*execAction
	var err error
	for _, inputParam := range param {
		if len(inputParam.MGMT) == 0 {
			err = fmt.Errorf("Param MGMT can not empty ")
			break
		}
		tmpAction := execAction{Param: []interface{}{}}
		if inputParam.FilterAttr == "" {
			tmpAction = execAction{Sql: "update sys_view set name=?,report=?,editable=?,suport_version=?,multiple=?,filter_attr=NULL,filter_value='',update_time=?,update_user=? where id=?"}
			tmpAction.Param = append(tmpAction.Param, inputParam.Name, inputParam.Report, inputParam.Editable, inputParam.SuportVersion, inputParam.Multiple, nowTime, operator, inputParam.Id)
		} else {
			tmpAction = execAction{Sql: "update sys_view set name=?,report=?,editable=?,suport_version=?,multiple=?,filter_attr=?,filter_value=?,update_time=?,update_user=? where id=?"}
			tmpAction.Param = append(tmpAction.Param, inputParam.Name, inputParam.Report, inputParam.Editable, inputParam.SuportVersion, inputParam.Multiple, inputParam.FilterAttr, inputParam.FilterValue, nowTime, operator, inputParam.Id)
		}
		actions = append(actions, &tmpAction)
		actions = append(actions, &execAction{Sql: "delete from sys_role_view where `view`=?", Param: []interface{}{inputParam.Id}})
		actions = append(actions, getPermissionActions(inputParam.MGMT, inputParam.USE, inputParam.Id)...)
	}
	if err != nil {
		return err
	}
	return transaction(actions)
}

func ViewDelete(ids []string) error {
	if len(ids) == 0 {
		return nil
	}
	filterSql, filterParams := createListParams(ids, "")
	var viewTable []*models.SysViewTable
	err := x.SQL("select id from sys_view where id in ("+filterSql+")", filterParams...).Find(&viewTable)
	if err != nil {
		err = fmt.Errorf("Try to query view from database fail,%s ", err.Error())
		return err
	}
	var actions []*execAction
	for _, view := range viewTable {
		graphActions, tmpErr := getGraphDeleteActions([]string{}, view.Id)
		if tmpErr != nil {
			err = tmpErr
			break
		}
		actions = append(actions, graphActions...)
		actions = append(actions, &execAction{Sql: "delete from sys_view where id=?", Param: []interface{}{view.Id}})
	}
	if err != nil {
		return err
	}
	actions = append(actions, &execAction{Sql: "delete from sys_role_view where view in (" + filterSql + ")", Param: filterParams})
	if len(actions) == 0 {
		return nil
	}
	return transactionWithoutForeignCheck(actions)
}

func GraphList(param *models.QueryRequestParam) (pageInfo models.PageInfo, rowData []*models.SysGraphTable, err error) {
	rowData = []*models.SysGraphTable{}
	filterSql, queryColumn, queryParam := transFiltersToSQL(param, &models.TransFiltersParam{IsStruct: true, StructObj: models.SysGraphTable{}, PrimaryKey: "id"})
	baseSql := fmt.Sprintf("SELECT %s FROM sys_graph WHERE 1=1 %s ", queryColumn, filterSql)
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

func GraphCreate(param []*models.SysGraphTable) error {
	if len(param) == 0 {
		return nil
	}
	var actions []*execAction
	for _, inputParam := range param {
		if !strings.HasPrefix(inputParam.Id, inputParam.View) {
			inputParam.Id = inputParam.View + models.SysTableIdConnector + inputParam.Id
		}
		tmpAction := execAction{Sql: "insert into sys_graph(`id`,`name`,`view`,`graph_type`,`node_groups`,`graph_dir`,`graph_node_config`,`graph_edge_config`) value (?,?,?,?,?,?,?,?)", Param: []interface{}{}}
		tmpAction.Param = append(tmpAction.Param, inputParam.Id, inputParam.Name, inputParam.View, inputParam.GraphType, inputParam.NodeGroups, inputParam.GraphDir, inputParam.GraphNodeConfig, inputParam.GraphEdgeConfig)
		actions = append(actions, &tmpAction)
	}
	return transaction(actions)
}

func GraphUpdate(param []*models.SysGraphTable) error {
	if len(param) == 0 {
		return nil
	}
	var actions []*execAction
	for _, inputParam := range param {
		tmpAction := execAction{Sql: "update sys_graph set `name`=?,`view`=?,`graph_type`=?,`node_groups`=?,`graph_dir`=?,`graph_node_config`=?,`graph_edge_config`=? where id=?", Param: []interface{}{}}
		tmpAction.Param = append(tmpAction.Param, inputParam.Name, inputParam.View, inputParam.GraphType, inputParam.NodeGroups, inputParam.GraphDir, inputParam.GraphNodeConfig, inputParam.GraphEdgeConfig, inputParam.Id)
		actions = append(actions, &tmpAction)
	}
	return transaction(actions)
}

func GraphDelete(ids []string) error {
	if len(ids) == 0 {
		return nil
	}
	actions, err := getGraphDeleteActions(ids, "")
	if err != nil {
		return err
	}
	if len(actions) == 0 {
		return nil
	}
	return transactionWithoutForeignCheck(actions)
}

func ElementList(param *models.QueryRequestParam) (pageInfo models.PageInfo, rowData []*models.SysGraphElementTable, err error) {
	rowData = []*models.SysGraphElementTable{}
	filterSql, queryColumn, queryParam := transFiltersToSQL(param, &models.TransFiltersParam{IsStruct: true, StructObj: models.SysGraphElementTable{}, PrimaryKey: "id"})
	baseSql := fmt.Sprintf("SELECT %s FROM sys_graph_element WHERE 1=1 %s ", queryColumn, filterSql)
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

func ElementCreate(param []*models.SysGraphElementTable) error {
	if len(param) == 0 {
		return nil
	}
	var actions []*execAction
	for _, inputParam := range param {
		if !strings.HasPrefix(inputParam.Id, inputParam.Graph) {
			inputParam.Id = inputParam.Graph + models.SysTableIdConnector + inputParam.Id
		}
		tmpAction := execAction{Param: []interface{}{}}
		tmpAction.Param = append(tmpAction.Param, inputParam.Id, inputParam.Graph, inputParam.ReportObject, inputParam.ShowTable, inputParam.DisplayExpression, inputParam.NodeGroupName, inputParam.LineStartData, inputParam.LineEndData, inputParam.LineDisplayPosition, inputParam.GraphType, inputParam.GraphShapeData, inputParam.GraphShapes, inputParam.GraphConfigData, inputParam.GraphConfigs, inputParam.GraphFilterData, inputParam.GraphFilterValues, inputParam.Editable, inputParam.SeqNo, inputParam.OrderData, inputParam.UpdateOperation)
		tmpColumns := append([]string{"id"}, models.BaseElementColumns...)
		if inputParam.ParentElement != "" {
			tmpColumns = append(tmpColumns, "parent_element")
			tmpAction.Param = append(tmpAction.Param, inputParam.ParentElement)
		}
		if inputParam.EditRefAttr != "" {
			tmpColumns = append(tmpColumns, "edit_ref_attr")
			tmpAction.Param = append(tmpAction.Param, inputParam.EditRefAttr)
		}
		tmpSpecCharList := []string{}
		for i := 0; i < len(tmpColumns); i++ {
			tmpSpecCharList = append(tmpSpecCharList, "?")
		}
		tmpAction.Sql = "insert into sys_graph_element(" + strings.Join(tmpColumns, ",") + ") value (" + strings.Join(tmpSpecCharList, ",") + ")"
		actions = append(actions, &tmpAction)
	}
	return transaction(actions)
}

func ElementUpdate(param []*models.SysGraphElementTable) error {
	if len(param) == 0 {
		return nil
	}
	var actions []*execAction
	for _, inputParam := range param {
		tmpAction := execAction{Param: []interface{}{}}
		tmpAction.Param = append(tmpAction.Param, inputParam.Graph, inputParam.ReportObject, inputParam.ShowTable, inputParam.DisplayExpression, inputParam.NodeGroupName, inputParam.LineStartData, inputParam.LineEndData, inputParam.LineDisplayPosition, inputParam.GraphType, inputParam.GraphShapeData, inputParam.GraphShapes, inputParam.GraphConfigData, inputParam.GraphConfigs, inputParam.GraphFilterData, inputParam.GraphFilterValues, inputParam.Editable, inputParam.SeqNo, inputParam.OrderData, inputParam.UpdateOperation)
		tmpColumns := []string{}
		tmpColumns = append(tmpColumns, models.BaseElementColumns...)
		if inputParam.ParentElement != "" {
			tmpColumns = append(tmpColumns, "parent_element")
			tmpAction.Param = append(tmpAction.Param, inputParam.ParentElement)
		}
		if inputParam.EditRefAttr != "" {
			tmpColumns = append(tmpColumns, "edit_ref_attr")
			tmpAction.Param = append(tmpAction.Param, inputParam.EditRefAttr)
		}
		for i, v := range tmpColumns {
			tmpColumns[i] = v + "=?"
		}
		tmpAction.Sql = "update sys_graph_element set " + strings.Join(tmpColumns, ",") + " where id=?"
		tmpAction.Param = append(tmpAction.Param, inputParam.Id)
		actions = append(actions, &tmpAction)
	}
	return transaction(actions)
}

func ElementDelete(ids []string) error {
	if len(ids) == 0 {
		return nil
	}
	actions, err := getElementDeleteActions(ids, "")
	if err != nil {
		return err
	}
	if len(actions) == 0 {
		return nil
	}
	return transactionWithoutForeignCheck(actions)
}

func GetElementRefAttrs(parentElement, reportObject string) (result []*models.SysCiTypeAttrTable, err error) {
	var sourceReportObjects, targetReportObjects []*models.SysReportObjectTable
	result = []*models.SysCiTypeAttrTable{}
	err = x.SQL("select ci_type from sys_report_object where id=?", reportObject).Find(&sourceReportObjects)
	if err != nil {
		return
	}
	if len(sourceReportObjects) == 0 {
		err = fmt.Errorf("Can not find reportObject with id=%s ", reportObject)
		return
	}
	sourceCiType := sourceReportObjects[0].CiType
	err = x.SQL("select ci_type from sys_report_object where id in (select report_object from sys_graph_element where id=?)", parentElement).Find(&targetReportObjects)
	if err != nil {
		return
	}
	if len(targetReportObjects) == 0 {
		err = fmt.Errorf("Can not find element or related reportObject with element id=%s ", parentElement)
		return
	}
	targetCiType := targetReportObjects[0].CiType
	err = x.SQL("select * from sys_ci_type_attr where ci_type=? and ref_ci_type=? and status='created'", sourceCiType, targetCiType).Find(&result)
	return
}

func getGraphDeleteActions(ids []string, view string) (actions []*execAction, err error) {
	var filterSql string
	var filterParams []interface{}
	if len(ids) > 0 {
		filterSql, filterParams = createListParams(ids, "")
		filterSql = " where id in (" + filterSql + ")"
	} else if view != "" {
		filterSql = " where view=?"
		filterParams = []interface{}{view}
	}
	var fetchRows []*models.SysGraphTable
	err = x.SQL("select id from sys_graph "+filterSql, filterParams...).Find(&fetchRows)
	if err != nil {
		err = fmt.Errorf("Try to query graph from database fail,%s ", err.Error())
		return
	}
	for _, row := range fetchRows {
		elementActions, tmpErr := getElementDeleteActions([]string{}, row.Id)
		if tmpErr != nil {
			err = tmpErr
			break
		}
		actions = append(actions, elementActions...)
		actions = append(actions, &execAction{Sql: "delete from sys_graph where id=?", Param: []interface{}{row.Id}})
	}
	return
}

func getElementDeleteActions(ids []string, graph string) (actions []*execAction, err error) {
	var filterSql string
	var filterParams []interface{}
	if len(ids) > 0 {
		filterSql, filterParams = createListParams(ids, "")
		filterSql = " where id in (" + filterSql + ")"
	} else if graph != "" {
		filterSql = " where graph=?"
		filterParams = []interface{}{graph}
	}
	var fetchRows, childRows []*models.SysGraphElementTable
	err = x.SQL("select id,parent_element from sys_graph_element "+filterSql, filterParams...).Find(&fetchRows)
	if err != nil {
		err = fmt.Errorf("Try to query element from database fail,%s ", err.Error())
		return
	}
	var parentIds []string
	for _, row := range fetchRows {
		if row.ParentElement == "" {
			parentIds = append(parentIds, row.Id)
		}
	}
	if len(parentIds) > 0 {
		err = x.SQL("select id from sys_graph_element where parent_element in ('" + strings.Join(parentIds, "','") + "')").Find(&childRows)
		if err != nil {
			err = fmt.Errorf("Try to query child element from database fail,%s ", err.Error())
			return
		}
		if len(childRows) > 0 {
			fetchRows = append(fetchRows, childRows...)
		}
	}
	for _, row := range fetchRows {
		actions = append(actions, &execAction{Sql: "delete from sys_graph_element where id=?", Param: []interface{}{row.Id}})
	}
	return
}

func getPermissionActions(mgmt, use []string, viewId string) []*execAction {
	actions := []*execAction{}
	for _, v := range mgmt {
		actions = append(actions, &execAction{Sql: "insert into sys_role_view(`id`,`role`,`view`,`permission`) value (?,?,?,?)", Param: []interface{}{v + models.SysTableIdConnector + viewId + models.SysTableIdConnector + "MGMT", v, viewId, "MGMT"}})
	}
	for _, v := range use {
		actions = append(actions, &execAction{Sql: "insert into sys_role_view(`id`,`role`,`view`,`permission`) value (?,?,?,?)", Param: []interface{}{v + models.SysTableIdConnector + viewId + models.SysTableIdConnector + "USE", v, viewId, "USE"}})
	}
	return actions
}

func ValidateViewPermission(viewIds, roles []string) (legal bool, err error) {
	if len(viewIds) == 0 {
		return true, nil
	}
	if len(roles) == 0 {
		return false, nil
	}
	queryRows, queryErr := x.QueryString("select distinct `view` from sys_role_view where `role` in ('" + strings.Join(roles, "','") + "') and permission='MGMT'")
	if queryErr != nil {
		err = fmt.Errorf("Query permission role fail,%s ", queryErr.Error())
		return
	}
	legal = true
	for _, reqView := range viewIds {
		tmpLegal := false
		for _, legalView := range queryRows {
			if reqView == legalView["view"] {
				tmpLegal = true
				break
			}
		}
		if !tmpLegal {
			legal = false
			break
		}
	}
	return
}

func ValidateGraphPermission(graphIds, roles []string) (legal bool, err error) {
	if len(graphIds) == 0 {
		return true, nil
	}
	filterSql, filterParam := createListParams(graphIds, "")
	var graphTable []*models.SysGraphTable
	err = x.SQL("select distinct `view` from sys_graph where id in ("+filterSql+")", filterParam...).Find(&graphTable)
	if err != nil {
		err = fmt.Errorf("Try to get graph owner view rocord fail,%s ", err.Error())
		return
	}
	viewIds := []string{}
	for _, row := range graphTable {
		viewIds = append(viewIds, row.View)
	}
	return ValidateViewPermission(viewIds, roles)
}

func ValidateElementPermission(elementIds, roles []string) (legal bool, err error) {
	if len(elementIds) == 0 {
		return true, nil
	}
	filterSql, filterParam := createListParams(elementIds, "")
	var elementTable []*models.SysGraphElementTable
	err = x.SQL("select distinct `graph` from sys_graph_element where id in ("+filterSql+")", filterParam...).Find(&elementTable)
	if err != nil {
		err = fmt.Errorf("Try to get element owner graph rocord fail,%s ", err.Error())
		return
	}
	graphIds := []string{}
	for _, row := range elementTable {
		graphIds = append(graphIds, row.Graph)
	}
	return ValidateGraphPermission(graphIds, roles)
}
