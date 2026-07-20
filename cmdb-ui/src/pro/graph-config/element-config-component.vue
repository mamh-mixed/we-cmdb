<template>
  <div>
    <div style="margin-bottom:8px">
      <Button type="primary" @click="openAddElement">{{ $t('add_element') }}</Button>
      <Button style="margin: 0 5px" @click="editGraph">{{ $t('edit_graph') }}</Button>
      <Button type="error" @click="removeGraph">{{ $t('delete_graph') }}</Button>
    </div>
    <Table border :columns="columns" :data="data"></Table>
    <!-- add element -->
    <Modal
      v-model="newElement.showNewElement"
      width="700"
      :mask-closable="false"
      :title="(newElement.isAdd ? $t('new') : $t('edit')) + ' ' + $t('view_element')"
    >
      <div slot="footer">
        <Button @click="cancelNewElementHandler">{{ $t('cancel') }}</Button>
        <Button type="primary" @click="elementManagementHandler">{{ $t('confirm') }}</Button>
      </div>
      <div :style="{ maxHeight: MODALHEIGHT + 'px', overflow: 'auto' }">
        <Form :label-width="160" label-colon>
          <FormItem label="ID" v-if="showFormField('graphElementId')">
            <Input v-model="newElement.newElementForm.graphElementId" :disabled="!newElement.isAdd"></Input>
          </FormItem>
          <FormItem :label="$t('parent_element')" v-if="showFormField('parentElement')">
            <Select
              v-model="newElement.newElementForm.parentElement"
              @on-clear="clearParentElement"
              filterable
              clearable
            >
              <Option v-for="item in data" :value="item.graphElementId" :key="item.graphElementId">{{
                item.graphElementId
              }}</Option>
            </Select>
          </FormItem>
          <FormItem :label="$t('report_object')" v-if="showFormField('reportObject')">
            <Select
              v-model="newElement.newElementForm.reportObject"
              @on-open-change="getReportObject"
              @on-change="getAttrByReportObject"
              @on-clear="clearReportObject"
              filterable
              clearable
            >
              <Option
                v-for="item in newElement.reportObjectOptions"
                :value="item.reportObjectId"
                :key="item.reportObjectId"
                >{{ item.dataName }}</Option
              >
            </Select>
          </FormItem>
          <FormItem :label="$t('is_displayed_table')" v-if="showFormField('showTable')">
            <Checkbox v-model="newElement.newElementForm.showTable" true-value="yes" false-value="no"></Checkbox>
          </FormItem>
          <FormItem :label="$t('is_editable')" v-if="showFormField('editable')">
            <Checkbox v-model="newElement.newElementForm.editable" true-value="yes" false-value="no"></Checkbox>
          </FormItem>
          <FormItem :label="$t('graph_type')" v-if="showFormField('graphType')">
            <Select
              v-if="(graph || {}).graphType === 'sequence'"
              v-model="newElement.newElementForm.graphType"
              filterable
              clearable
              @on-change="handleChangeGraphType"
            >
              <Option v-for="item in newElement.graphTypeSeqOptions" :value="item.id" :key="item.id">{{
                item.name
              }}</Option>
            </Select>
            <Select
              v-else
              v-model="newElement.newElementForm.graphType"
              filterable
              clearable
              @on-change="handleChangeGraphType"
            >
              <Option v-for="item in newElement.graphTypeOptions" :value="item.id" :key="item.id">{{
                item.name
              }}</Option>
            </Select>
          </FormItem>
          <FormItem :label="$t('display_expression')" v-if="showFormField('displayExpression')">
            <Input v-model="newElement.newElementForm.displayExpression"></Input>
          </FormItem>
          <FormItem :label="$t('edit_ref_attr')" v-if="showFormField('editRefAttr')">
            <Select v-model="newElement.newElementForm.editRefAttr" filterable clearable>
              <Option v-for="item in newElement.refAttrs" :value="item.ciTypeAttrId" :key="item.ciTypeAttrId">{{
                item.ciTypeAttrId
              }}</Option>
            </Select>
          </FormItem>
          <FormItem :label="$t('node_group_name')" v-if="showFormField('nodeGroupName')">
            <Select v-model="newElement.newElementForm.nodeGroupName" filterable clearable>
              <Option v-for="item in newElement.newElementForm.nodeGroupsItem" :value="item" :key="item">{{
                item
              }}</Option>
            </Select>
          </FormItem>
          <FormItem :label="$t('line_start_data')" v-if="showFormField('lineStartData')">
            <Select v-model="newElement.newElementForm.lineStartData" filterable clearable>
              <Option v-for="item in newElement.attrs" :value="item.dataName" :key="item.dataName">{{
                item.dataName
              }}</Option>
            </Select>
          </FormItem>
          <FormItem :label="$t('line_end_data')" v-if="showFormField('lineEndData')">
            <Select v-model="newElement.newElementForm.lineEndData" filterable clearable>
              <Option v-for="item in newElement.attrs" :value="item.dataName" :key="item.dataName">{{
                item.dataName
              }}</Option>
            </Select>
          </FormItem>
          <FormItem :label="$t('line_display_position')" v-if="showFormField('lineDisplayPosition')">
            <Select v-model="newElement.newElementForm.lineDisplayPosition" filterable clearable>
              <Option v-for="item in newElement.lineDisplayPositionOptions" :value="item.id" :key="item.id">{{
                item.name
              }}</Option>
            </Select>
          </FormItem>
          <FormItem :label="$t('graph_shape_data')" v-if="showFormField('graphShapeData')">
            <Select
              v-model="newElement.newElementForm.graphShapeData"
              filterable
              clearable
              allow-create
              @on-create="handleCreateDataOptions"
              @on-change="handleChangeGraphShapeData"
            >
              <Option v-for="item in newElement.attrs" :value="item.dataName" :key="item.dataName">{{
                item.dataName
              }}</Option>
            </Select>
          </FormItem>
          <FormItem :label="$t('graph_shapes')" v-if="showFormField('graphShapes')">
            <Input
              v-if="newElement.newElementForm.graphShapeData"
              v-model="newElement.newElementForm.graphShapes"
              type="textarea"
              :autosize="{ minRows: 3 }"
            >
            </Input>
            <Select v-else v-model="newElement.newElementForm.graphShapes" filterable clearable>
              <Option v-for="item in newElement.graphShapesOptions" :value="item.id" :key="item.id">{{
                item.name
              }}</Option>
            </Select>
          </FormItem>
          <!-- <FormItem :label="$t('graph_config_data')" v-if="showFormField('graphConfigData')">
            <Select v-model="newElement.newElementForm.graphConfigData" filterable clearable>
              <Option v-for="item in newElement.attrs" :value="item.dataName" :key="item.dataName">{{
                item.dataName
              }}</Option>
            </Select>
          </FormItem>
          <FormItem :label="$t('graph_configs')" v-if="showFormField('graphConfigs')">
            <Input v-model="newElement.newElementForm.graphConfigs" type="textarea" :autosize="{ minRows: 3 }"> </Input>
          </FormItem> -->
          <template v-if="showFormField('graphConfigData')">
            <FormItem :label="$t('graph_config_data')">
              <Row>
                <Col span="10">
                  <span style="color:red">*</span>
                  {{ $t('form_name') }}
                </Col>
                <Col span="10">
                  <span style="color:red">*</span>
                  {{ $t('cmdb_supported_version') }}
                </Col>
              </Row>
              <Row v-for="(item, itemIndex) in graphConfigData_tmp" :key="itemIndex" style="margin:6px 0;">
                <Col span="10">
                  <Select
                    v-model="item.name"
                    filterable
                    @on-change="changeGraphConfigDataKey(item.name)"
                    style="width: 92%;"
                  >
                    <Option
                      v-for="attr in newElement.attrs"
                      :value="attr.dataName"
                      :label="attr.dataName"
                      :key="attr.dataName"
                      v-if="
                        item.name === attr.dataName ||
                          graphConfigData_tmp.findIndex(tmp => tmp.name === attr.dataName) === -1
                      "
                    >
                      {{ attr.dataName }}
                    </Option>
                  </Select>
                </Col>
                <Col span="10">
                  <Select v-model="item.suport_version" style="width: 92%;">
                    <Option value="yes" key="yes">yes</Option>
                    <Option value="no" key="no">no</Option>
                  </Select>
                </Col>
                <Col span="2" offset="1">
                  <Button
                    type="error"
                    ghost
                    @click="deleteItem('graphConfigData_tmp', itemIndex)"
                    size="small"
                    style="cursor: pointer"
                    icon="md-trash"
                  ></Button>
                </Col>
              </Row>
              <Row>
                <Col span="2" offset="21">
                  <div style="cursor: pointer">
                    <Button
                      type="success"
                      ghost
                      :disabled="newElement.attrs.length === graphConfigData_tmp.length"
                      @click="addItem('graphConfigData_tmp', { name: '', suport_version: 'yes' })"
                      size="small"
                      icon="md-add"
                    ></Button>
                  </div>
                </Col>
              </Row>
            </FormItem>
            <FormItem v-if="Object.keys(graphConfigs_tmp).length > 0">
              <Row>
                <Col span="6">
                  <span style="visibility: hidden;">hide</span>
                </Col>
                <Col span="8">
                  <span style="color:red">*</span>
                  {{ $t('cmdb_key') }}
                </Col>
                <Col span="8">
                  <span style="color:red">*</span>
                  {{ $t('cmdb_value') }}
                </Col>
              </Row>
              <div v-for="graphConfigsKey in Object.keys(graphConfigs_tmp)" :key="graphConfigsKey">
                <Row>
                  <Col span="6">
                    {{ graphConfigsKey }}
                  </Col>
                  <Col span="18">
                    <Row
                      v-for="(item, itemIndex) in graphConfigs_tmp[graphConfigsKey]"
                      :key="itemIndex"
                      style="margin: 4px 0;"
                    >
                      <Col span="10">
                        <Input v-model="item.key" style="width: 90%;"></Input>
                      </Col>
                      <Col span="10">
                        <Input v-model="item.value" style="width: 90%;" type="textarea"></Input>
                      </Col>
                      <Col span="2">
                        <Button
                          type="error"
                          ghost
                          @click="deleteItem('graphConfigs_tmp', itemIndex, graphConfigsKey)"
                          size="small"
                          style="cursor: pointer"
                          icon="md-trash"
                        ></Button>
                      </Col>
                    </Row>
                    <Row>
                      <Col span="2" offset="20">
                        <div style="cursor: pointer">
                          <Button
                            type="success"
                            ghost
                            @click="addItem('graphConfigs_tmp', { key: '', value: '' }, graphConfigsKey)"
                            size="small"
                            icon="md-add"
                          ></Button>
                        </div>
                      </Col>
                    </Row>
                  </Col>
                </Row>
              </div>
              <!-- <Input v-model="newElement.newElementForm.graphFilterData"></Input> -->
            </FormItem>
          </template>
          <FormItem :label="$t('graph_filter_data')" v-if="showFormField('graphFilterData')">
            <Input v-model="newElement.newElementForm.graphFilterData"></Input>
          </FormItem>
          <FormItem :label="$t('graph_filter_values')" v-if="showFormField('graphFilterValues')">
            <Input v-model="newElement.newElementForm.graphFilterValues"></Input>
          </FormItem>
          <FormItem :label="$t('graph_order_data')" v-if="showFormField('orderData')">
            <Select v-model="newElement.newElementForm.orderData" filterable clearable>
              <Option v-for="item in newElement.attrs" :value="item.dataName" :key="item.dataName">{{
                item.dataName
              }}</Option>
            </Select>
          </FormItem>
          <FormItem :label="$t('graph_update_operation')" v-if="showFormField('updateOperation')">
            <Input v-model="newElement.newElementForm.updateOperation"></Input>
          </FormItem>
        </Form>
      </div>
    </Modal>
    <!-- 编辑graph模态框 -->
    <Modal
      v-model="newGraph.showNewGraph"
      width="700"
      @on-ok="editGraphHandler"
      @on-cancel="cancelGraphHandler"
      :title="$t('edit_graph')"
    >
      <Form :label-width="160" label-colon>
        <FormItem label="ID">
          <Input disabled v-model="newGraph.newGraphForm.graphId"></Input>
        </FormItem>
        <FormItem :label="$t('display_name')">
          <Input v-model="newGraph.newGraphForm.name"></Input>
        </FormItem>
        <FormItem :label="$t('graph_type')">
          <Select v-model="newGraph.newGraphForm.graphType" filterable clearable>
            <Option v-for="item in newGraph.graphTypeOptions" :value="item.id" :key="item.id">{{ item.name }}</Option>
          </Select>
        </FormItem>
        <FormItem :label="$t('node_group')">
          <Input v-model="newGraph.newGraphForm.nodeGroups" placeholder="eg: system->subSystem->unit->instance"></Input>
          <span style="color:#c5c5c5">eg: system->subSystem->unit->instance</span>
        </FormItem>
        <FormItem :label="$t('graph_dir')">
          <Select v-model="newGraph.newGraphForm.graphDir" filterable clearable>
            <Option v-for="item in newGraph.graphDirOptions" :value="item.id" :key="item.id">{{ item.name }}</Option>
          </Select>
        </FormItem>
        <FormItem :label="$t('cmdb_node_config')">
          <Input v-model="newGraph.newGraphForm.graphNodeConfig" type="textarea"></Input>
        </FormItem>
        <FormItem :label="$t('cmdb_edge_config')">
          <Input v-model="newGraph.newGraphForm.graphEdgeConfig" type="textarea"></Input>
        </FormItem>
      </Form>
    </Modal>
  </div>
</template>
<script>
import {
  addElementForGraph,
  deleteElementForGraph,
  deleteGraph,
  editElementForGraph,
  editGraph,
  getAttrByReportObject,
  getCiTypeAttributes,
  getReportObject
} from '@/api/server'
export default {
  data () {
    return {
      MODALHEIGHT: 600,
      columns: [
        {
          title: 'ID',
          key: 'graphElementId',
          width: 120,
          resizable: true,
          ellipsis: true,
          tooltip: true
        },
        // {
        //   title: this.$t('view_graph'),
        //   key: 'graph',
        //   width: 120,
        //   ellipsis: true,
        //   tooltip: true
        // },
        {
          title: this.$t('parent_element'),
          key: 'parentElement',
          width: 120,
          resizable: true,
          ellipsis: true,
          tooltip: true
        },
        {
          title: this.$t('report_object'),
          key: 'reportObject',
          width: 120,
          resizable: true,
          ellipsis: true,
          tooltip: true
        },
        {
          title: this.$t('is_displayed_table'),
          key: 'showTable',
          width: 120,
          resizable: true,
          ellipsis: true,
          tooltip: true
        },
        {
          title: this.$t('is_editable'),
          key: 'editable',
          width: 120,
          resizable: true,
          ellipsis: true,
          tooltip: true
        },
        {
          title: this.$t('graph_type'),
          key: 'graphType',
          width: 120,
          resizable: true,
          ellipsis: true,
          tooltip: true
        },
        {
          title: this.$t('display_expression'),
          key: 'displayExpression',
          width: 120,
          resizable: true,
          ellipsis: true,
          tooltip: true
        },
        {
          title: this.$t('edit_ref_attr'),
          key: 'editRefAttr',
          width: 120,
          resizable: true,
          ellipsis: true,
          tooltip: true
        },
        {
          title: this.$t('node_group_name'),
          key: 'nodeGroupName',
          width: 120,
          resizable: true,
          ellipsis: true,
          tooltip: true
        },
        {
          title: this.$t('line_start_data'),
          key: 'lineStartData',
          width: 120,
          resizable: true,
          ellipsis: true,
          tooltip: true
        },

        {
          title: this.$t('line_end_data'),
          key: 'lineEndData',
          width: 120,
          resizable: true,
          ellipsis: true,
          tooltip: true
        },
        {
          title: this.$t('line_display_position'),
          key: 'lineDisplayPosition',
          width: 120,
          resizable: true,
          ellipsis: true,
          tooltip: true
        },
        {
          title: this.$t('graph_shape_data'),
          key: 'graphShapeData',
          width: 120,
          resizable: true,
          ellipsis: true,
          tooltip: true
        },
        {
          title: this.$t('graph_shapes'),
          key: 'graphShapes',
          width: 120,
          resizable: true,
          ellipsis: true,
          tooltip: true
        },
        {
          title: this.$t('graph_config_data'),
          key: 'graphConfigData',
          width: 120,
          resizable: true,
          ellipsis: true,
          tooltip: true
        },
        {
          title: this.$t('graph_configs'),
          key: 'graphConfigs',
          width: 120,
          resizable: true,
          ellipsis: true,
          tooltip: true
        },
        {
          title: this.$t('graph_filter_values'),
          key: 'graphFilterData',
          width: 120,
          resizable: true,
          ellipsis: true,
          tooltip: true
        },
        {
          title: this.$t('graph_filter_data'),
          key: 'graphFilterValues',
          width: 120,
          resizable: true,
          ellipsis: true,
          tooltip: true
        },
        {
          title: this.$t('graph_order_data'),
          key: 'orderData',
          width: 120,
          resizable: true,
          ellipsis: true,
          tooltip: true
        },
        {
          title: this.$t('graph_update_operation'),
          key: 'updateOperation',
          width: 120,
          resizable: true,
          ellipsis: true,
          tooltip: true
        },
        {
          title: this.$t('seq_no'),
          key: 'seqNo',
          width: 120,
          resizable: true,
          ellipsis: true,
          tooltip: true
        },
        {
          title: this.$t('actions'),
          key: 'action',
          width: 150,
          align: 'center',
          fixed: 'right',
          render: (h, params) => {
            return h('div', [
              h(
                'Button',
                {
                  props: {
                    type: 'primary',
                    size: 'small'
                  },
                  style: {
                    marginRight: '5px'
                  },
                  on: {
                    click: () => {
                      this.editElement(params.row)
                    }
                  }
                },
                this.$t('edit')
              ),
              h(
                'Button',
                {
                  props: {
                    type: 'error',
                    size: 'small'
                  },
                  style: {
                    marginRight: '5px'
                  },
                  on: {
                    click: () => {
                      this.removeElemment(params.row)
                    }
                  }
                },
                this.$t('delete')
              )
            ])
          }
        }
      ],
      data: [],
      report: '',
      graph: null,
      newElement: {
        showNewElement: false,
        isAdd: true,
        newElementForm: {
          graphElementId: '',
          graph: '',
          parentElement: '',
          reportObject: '',
          showTable: 'yes',
          displayExpression: '',
          nodeGroupName: '',
          nodeGroupsItem: [],
          lineStartData: '',
          lineEndData: '',
          lineDisplayPosition: '',
          graphType: '',
          graphShapeData: '',
          graphShapes: '',
          graphConfigData: '[]',
          graphConfigs: '',
          editRefAttr: '',
          graphFilterData: '',
          graphFilterValues: '',
          editable: 'yes',
          orderData: '',
          updateOperation: '',
          seqNo: ''
        },
        attrs: [],
        refAttrs: [],
        reportObjectOptions: [],
        lineDisplayPositionOptions: [
          { id: 'head', name: 'head' },
          { id: 'tail', name: 'tail' }
        ],
        graphTypeOptions: [
          { id: 'subgraph', name: 'subgraph' },
          { id: 'line', name: 'line' },
          { id: 'node', name: 'node' },
          { id: 'image', name: 'image' },
          { id: 'none', name: 'none' }
        ],
        graphTypeSeqOptions: [
          { id: 'sequence_diagram', name: 'sequence_diagram' },
          { id: 'assist_item', name: 'assist_item' },
          { id: 'service_invoke_item', name: 'service_invoke_item' },
          { id: 'service_invoke', name: 'service_invoke' },
          { id: 'node', name: 'node' }
        ],
        graphShapesOptions: [
          { id: 'box', name: 'box' },
          { id: 'ellipse', name: 'ellipse' },
          { id: 'diamond', name: 'diamond' },
          { id: 'hexagon', name: 'hexagon' },
          { id: 'circle', name: 'circle' }
        ]
      },
      graphConfigData_tmp: [],
      graphConfigs_tmp: {},
      lastSeqNo: 0,
      newGraph: {
        showNewGraph: false,
        newGraphForm: {
          graphId: '',
          name: '',
          view: '',
          graphType: 'subgraph',
          nodeGroups: '',
          graphDir: 'LR',
          graphNodeConfig: '',
          graphEdgeConfig: ''
        },
        graphTypeOptions: [
          { id: 'subgraph', name: 'subgraph' },
          { id: 'group', name: 'group' },
          { id: 'sequence', name: 'sequence' }
        ],
        graphDirOptions: [
          { id: 'LR', name: 'LR' },
          { id: 'TB', name: 'TB' }
        ]
      }
    }
  },
  props: ['graphElement'],
  methods: {
    // #region 图标中图型配置数据
    addItem (key, value, scendKey) {
      if (key === 'graphConfigData_tmp') {
        this[key].push(value)
      }
      if (key === 'graphConfigs_tmp') {
        // 创建一个新副本，并对指定的键进行添加操作
        const updatedConfigs = { ...this[key] }
        updatedConfigs[scendKey] = [...updatedConfigs[scendKey], value]
        // 重新分配整个 graphConfigs_tmp 对象
        this.$set(this, key, updatedConfigs)
      }
    },
    deleteItem (key, index, scendKey) {
      if (key === 'graphConfigData_tmp') {
        this[key].splice(index, 1)
        this.changeGraphConfigDataKey('')
      }
      if (key === 'graphConfigs_tmp') {
        const updatedArray = [...this[key][scendKey]]
        updatedArray.splice(index, 1) // 从副本中移除元素
        this.$set(this[key], scendKey, updatedArray)
      }
    },
    changeGraphConfigDataKey (key) {
      const keys = Object.keys(this.graphConfigs_tmp)
      keys.forEach(key => {
        const index = this.graphConfigData_tmp.findIndex(item => item.name === key)
        if (index === -1) {
          delete this.graphConfigs_tmp[key]
        }
      })
      if (key) {
        this.graphConfigs_tmp[key] = [{ key: '', value: '' }]
      }
    },
    // #endregion
    showFormField (name) {
      let fieldExcludes = {
        subgraph: ['graph', 'nodeGroupName', 'seqNo', 'orderData', 'updateOperation'],
        group: ['graph', 'seqNo', 'orderData', 'updateOperation'],
        sequence: ['graph', 'nodeGroupName', 'lineDisplayPosition', 'seqNo']
      }
      let graphType = (this.graph || {}).graphType || 'subgraph'
      if (fieldExcludes[graphType].includes(name)) {
        return false
      }
      return true
    },
    handleChangeGraphType (value) {
      let graphShapesNodeOptions = [
        { id: 'box', name: 'box' },
        { id: 'ellipse', name: 'ellipse' },
        { id: 'diamond', name: 'diamond' },
        { id: 'hexagon', name: 'hexagon' },
        { id: 'circle', name: 'circle' }
      ]
      let graphShapesLineOptions = [
        { id: 'normal', name: 'normal' },
        { id: 'none', name: 'none' },
        { id: 'diamond', name: 'diamond' },
        { id: 'icurve', name: 'icurve' },
        { id: 'inv', name: 'inv' }
      ]
      switch (value) {
        case 'line':
          this.newElement.graphShapesOptions = graphShapesLineOptions
          break
        case 'subgraph':
        case 'image':
        default:
          // as 'node'
          this.newElement.graphShapesOptions = graphShapesNodeOptions
      }
    },
    handleChangeGraphShapeData (value) {
      if (!value) {
        this.newElement.newElementForm.graphShapes = ''
      }
    },
    handleCreateDataOptions (value) {
      this.newElement.attrs.push({ dataName: value })
    },
    initData (graphElement, graph, report) {
      this.data = graphElement
      if (this.data.length > 0) {
        this.lastSeqNo = Math.max(...this.data.map(item => Number(item.seqNo)))
      }
      this.graph = graph
      this.report = report
    },
    async openAddElement () {
      this.cancelNewElementHandler()
      this.newElement.isAdd = true
      this.newElement.showNewElement = true
    },
    async elementManagementHandler () {
      this.newElement.newElementForm.graph = this.graph.graphId
      if (this.newElement.isAdd) {
        this.newElement.newElementForm.seqNo = this.lastSeqNo + 1 + ''
      }
      // 排除空数据
      this.graphConfigData_tmp = this.graphConfigData_tmp.filter(item => item.name !== '')
      this.newElement.newElementForm.graphConfigData = JSON.stringify(this.graphConfigData_tmp)
      let tmp = {}
      const keys = Object.keys(this.graphConfigs_tmp)
      keys.forEach(key => {
        let obj = {}
        const val = this.graphConfigs_tmp[key].filter(item => item.key !== '' && item.key !== '')
        val.forEach(item => {
          obj[item.key] = item.value
        })
        tmp[key] = obj
      })
      this.newElement.newElementForm.graphConfigs = JSON.stringify(tmp)
      if (!this.newElement.newElementForm.reportObject) {
        this.$Message.warning(this.$t('report_object') + this.$t('is_required'))
        return
      }
      let method = this.newElement.isAdd ? addElementForGraph : editElementForGraph
      const { statusCode } = await method([this.newElement.newElementForm])
      if (statusCode === 'OK') {
        this.$Notice.success({
          title: 'Success',
          desc: 'Success !'
        })
        this.newElement.showNewElement = false
        this.$parent.$parent.$parent.getTabGraphData(this.graph.name)
      }
    },
    cancelNewElementHandler () {
      this.newElement.newElementForm = {
        graphElementId: '',
        graph: '',
        parentElement: '',
        reportObject: '',
        showTable: 'yes',
        displayExpression: '',
        nodeGroupName: '',
        nodeGroupsItem: [],
        lineStartData: '',
        lineEndData: '',
        lineDisplayPosition: '',
        graphType: '',
        graphShapeData: '',
        graphShapes: '',
        graphConfigData: '[]',
        graphConfigs: '',
        editRefAttr: '',
        graphFilterData: '',
        graphFilterValues: '',
        editable: 'yes',
        seqNo: ''
      }
      this.graphConfigData_tmp = []
      this.graphConfigs_tmp = {}
      this.clearParentElement()
      this.newElement.showNewElement = false
    },
    clearParentElement () {
      this.newElement.newElementForm.reportObject = ''
      this.newElement.reportObjectOptions = []
      this.newElement.newElementForm.lineStartData = ''
      this.newElement.newElementForm.lineEndData = ''
      this.newElement.attrs = []
      this.newElement.refAttrs = []
      this.newElement.editRefAttr = ''
      this.clearReportObject()
    },
    clearReportObject () {
      this.newElement.newElementForm.graphConfigData = '[]'
      this.newElement.newElementForm.graphConfigs = ''
      this.newElement.attrs = []
    },
    async editElement (row) {
      this.newElement.isAdd = false
      this.newElement.newElementForm = { ...row }
      this.newElement.newElementForm.nodeGroupsItem = this.graph.nodeGroups.split('->')
      // 兼容老数据
      if (
        this.newElement.newElementForm.graphConfigData !== '' &&
        !(
          this.newElement.newElementForm.graphConfigData.startsWith('[') &&
          this.newElement.newElementForm.graphConfigData.endsWith(']')
        )
      ) {
        const tmp = this.newElement.newElementForm.graphConfigData
        this.newElement.newElementForm.graphConfigData = JSON.stringify([
          {
            name: tmp,
            suport_version: 'yes'
          }
        ])
        this.newElement.newElementForm.graphConfigs = `{"${tmp}":${this.newElement.newElementForm.graphConfigs}}`
      }
      try {
        const graphConfigData = JSON.parse(this.newElement.newElementForm.graphConfigData)
        this.graphConfigData_tmp = graphConfigData
      } catch (err) {
        this.graphConfigData_tmp = []
      }

      try {
        const graphConfigs = JSON.parse(this.newElement.newElementForm.graphConfigs)
        const keys = Object.keys(graphConfigs)
        keys.forEach(key => {
          let obj = []
          const val = graphConfigs[key]
          const keysTmp = Object.keys(val)
          keysTmp.forEach(key => {
            obj.push({ key: key, value: val[key] })
          })
          // 使用 $set 确保响应式更新
          this.$set(this.graphConfigs_tmp, key, obj)
        })
      } catch (err) {
        this.graphConfigs_tmp = {}
      }
      this.handleChangeGraphType(this.newElement.newElementForm.graphType)
      await this.getReportObject()
      await this.getAttrByReportObject()
      this.newElement.showNewElement = true
    },
    async removeElemment (item) {
      this.$Modal.confirm({
        title: this.$t('delete_confirm'),
        'z-index': 1000000,
        onOk: async () => {
          const { statusCode } = await deleteElementForGraph([item.graphElementId])
          if (statusCode === 'OK') {
            this.$Notice.success({
              title: 'Delete',
              desc: 'Success !'
            })
            this.$parent.$parent.$parent.getTabGraphData(this.graph.name)
          }
        },
        onCancel: () => {}
      })
    },
    editGraph () {
      this.newGraph.newGraphForm = { ...this.graph }
      this.newGraph.showNewGraph = true
    },
    async editGraphHandler () {
      const { statusCode } = await editGraph([this.newGraph.newGraphForm])
      if (statusCode === 'OK') {
        this.$Notice.success({
          title: this.$t('edit_graph'),
          desc: 'Success !'
        })
        this.$parent.$parent.$parent.getGraphByView()
      }
    },
    cancelGraphHandler () {
      this.newGraph.newGraphForm = {
        graphId: '',
        name: '',
        view: '',
        graphType: 'subgraph',
        nodeGroups: '',
        graphDir: 'LR',
        graphNodeConfig: '',
        graphEdgeConfig: ''
      }
      this.newGraph.showNewGraph = false
    },
    async removeGraph () {
      this.$Modal.confirm({
        title: this.$t('delete_confirm'),
        'z-index': 1000000,
        onOk: async () => {
          const { statusCode } = await deleteGraph([this.graph.graphId])
          if (statusCode === 'OK') {
            this.$Notice.success({
              title: this.$t('delete_view'),
              desc: 'Success !'
            })
          }
          this.$parent.$parent.$parent.getGraphByView()
        },
        onCancel: () => {}
      })
    },
    async getReportObject () {
      // eslint-disable-next-line no-unused-vars
      let op = 'eq'
      // eslint-disable-next-line no-unused-vars
      if (!this.newElement.newElementForm.parentElement) {
        op = 'null'
      }
      const find = this.data.find(item => item.graphElementId === this.newElement.newElementForm.parentElement)
      const params = {
        paging: false,
        filters: [
          { name: 'parentObject', operator: op, value: (find && find.reportObject) || '' },
          { name: 'report', operator: 'eq', value: this.report.report }
        ]
      }
      const { statusCode, data } = await getReportObject(params)
      if (statusCode === 'OK') {
        this.newElement.reportObjectOptions = data.contents
      }
    },
    async getAttrByReportObject () {
      const find = this.newElement.reportObjectOptions.find(
        item => item.reportObjectId === this.newElement.newElementForm.reportObject
      )
      if (find) {
        const params = {
          paging: false,
          filters: [{ name: 'reportObject', operator: 'eq', value: this.newElement.newElementForm.reportObject }]
        }
        const { statusCode, data } = await getAttrByReportObject(params)
        if (statusCode === 'OK') {
          this.newElement.attrs = data.contents
        }
        // if (this.newElement.newElementForm.graphShapeData) {
        //   this.newElement.attrs.push({ dataName: this.newElement.newElementForm.graphShapeData })
        // }
        // if (this.newElement.newElementForm.graphConfigData) {
        //   this.newElement.attrs.push({ dataName: this.newElement.newElementForm.graphConfigData })
        // }
        const attrRes = await getCiTypeAttributes(find.ciType)
        if (attrRes.statusCode === 'OK') {
          this.newElement.refAttrs = attrRes.data.filter(item => item.inputType === 'ref')
        }
      }
    }
  },
  mounted () {
    this.MODALHEIGHT = window.MODALHEIGHT
  }
}
</script>
<style scoped lang="scss"></style>
