<template>
  <div class="">
    <!-- 选择view -->
    <Row>
      <Col span="8">
        <Form :label-width="100">
          <FormItem :label="$t('view_graph')">
            <Select
              v-model="currentView"
              @on-open-change="getViewList"
              @on-clear="clearSelectView"
              ref="viewSelect"
              filterable
              clearable
            >
              <Button type="success" style="width:100%" @click="openAddView" size="small">
                <Icon type="ios-add" size="24"></Icon>
              </Button>
              <Option v-for="item in viewOptions" :value="item.viewId" :key="item.viewId"
                >{{ item.name
                }}<span style="float:right">
                  <Button style="margin-right: 5px" @click="editView(item)" icon="ios-create-outline" type="primary" size="small"></Button>
                  <Button @click.stop.prevent="deleteView(item)" icon="ios-trash" type="error" size="small"></Button>
                </span>
              </Option>
            </Select>
          </FormItem>
        </Form>
      </Col>
      <Col span="2">
        <Button type="primary" style="margin-left:24px" @click="getGraphByView" :disabled="!currentView">
          {{ $t('search') }}
        </Button>
      </Col>
    </Row>
    <div></div>
    <!-- 配置区tab -->
    <div v-if="showNewViewBtn">
      <Tabs :value="currentGraph" @on-click="getTabGraphData" type="card">
        <TabPane
          v-for="(graph, graphIndex) in graphs"
          :key="graph.graphId"
          :label="graph.name"
          :name="graph.name"
          :index="graphIndex"
        >
          <ElementConfig :ref="graph.name" :graphElement="graph.graphElement"></ElementConfig>
        </TabPane>
        <Button @click="handleTabsAdd" type="primary" slot="extra">{{ $t('new_graph') }}</Button>
      </Tabs>
    </div>
    <!-- 新增view模态框 -->
    <Modal
      width="720"
      v-model="newView.showNewView"
      @on-ok="addNewViewHandler"
      @on-cancel="cancelNewViewHandler"
      :title="(newView.isAdd ? $t('new') : $t('edit')) + $t('view')"
    >
      <div class="modal-height">
        <Form :label-width="100" label-colon>
          <FormItem label="ID">
            <Input v-model="newView.newViewForm.viewId" :disabled="!newView.isAdd"></Input>
          </FormItem>
          <FormItem :label="$t('display_name')">
            <Input v-model="newView.newViewForm.name"></Input>
          </FormItem>
          <FormItem :label="$t('report')">
            <Select v-model="newView.newViewForm.report" filterable clearable :disabled="!newView.isAdd">
              <Option v-for="item in newView.reportList" :value="item.id" :key="item.id">{{ item.name }}</Option>
            </Select>
          </FormItem>
          <FormItem :label="$t('is_editable')">
            <Checkbox v-model="newView.newViewForm.editable" true-value="yes" false-value="no"></Checkbox>
          </FormItem>
          <FormItem :label="$t('suport_version')">
            <Checkbox v-model="newView.newViewForm.suportVersion" true-value="yes" false-value="no"></Checkbox>
          </FormItem>
          <FormItem :label="$t('multiple')">
            <Checkbox v-model="newView.newViewForm.multiple" true-value="yes" false-value="no"></Checkbox>
          </FormItem>
          <FormItem :label="$t('filter_attr')">
            <Select
              v-model="newView.newViewForm.filterAttr"
              :disabled="!newView.isAdd"
              @on-open-change="getAttrByReportCi"
              filterable
              clearable
            >
              <Option v-for="item in newView.ciTypeAttrs" :value="item.ciTypeAttrId" :key="item.ciTypeAttrId">{{
                item.name
              }}</Option>
            </Select>
          </FormItem>
          <FormItem :label="$t('filter_value')">
            <Input v-model="newView.newViewForm.filterValue" :disabled="!newView.isAdd"></Input>
          </FormItem>
        </Form>
        <RoleConfig ref="roleManagement"></RoleConfig>
      </div>
    </Modal>
    <!-- 新增graph模态框 -->
    <Modal
      v-model="newGraph.showNewGraph"
      @on-ok="addNewGraphHandler"
      @on-cancel="cancelNewGraphHandler"
      :title="$t('new_graph')"
    >
      <Form :label-width="100" label-colon>
        <FormItem label="ID">
          <Input v-model="newGraph.newGraphForm.graphId"></Input>
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
  addGraph,
  addView,
  deleteGraph,
  deleteView,
  editView,
  getCiTypeAttr,
  getElementByGraph,
  getGraphByView,
  getReportListByPermission,
  getViewById,
  graphViews
} from '@/api/server'
import ElementConfig from './graph-config/element-config-component.vue'
import RoleConfig from './role-config'
export default {
  name: '',
  data () {
    return {
      currentView: '',
      viewOptions: [],

      currentGraph: '',
      graphs: [],

      showNewViewBtn: false,
      newView: {
        isAdd: true,
        showNewView: false,
        newViewForm: {
          viewId: '',
          name: '',
          report: '',
          editable: 'yes',
          suportVersion: 'yes',
          multiple: 'yes',
          filterAttr: '',
          filterValue: '',
          USE: [],
          MGMT: []
        },
        reportList: [],
        ciTypeAttrs: []
      },
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
  methods: {
    async getAttrByReportCi () {
      const find = this.newView.reportList.find(item => item.id === this.newView.newViewForm.report)
      if (find) {
        let attrs = await getCiTypeAttr(find.ciType)
        if (attrs.statusCode === 'OK') {
          this.newView.ciTypeAttrs = attrs.data
        }
      }
    },
    async openAddView () {
      this.graphs = []
      this.currentGraph = ''
      this.$refs.viewSelect.visible = false
      this.newView.reportList = []
      this.newView.newViewForm.USE = []
      this.newView.newViewForm.MGMT = []
      this.$refs.roleManagement.initRoleData([], [])
      this.cancelNewViewHandler()
      const { statusCode, data } = await getReportListByPermission('MGMT', 'view')
      if (statusCode === 'OK') {
        this.newView.reportList = data
        this.newView.showNewView = true
      }
    },
    async removeGraph (name) {
      const find = this.graphs.find(item => item.name === name)
      this.$Modal.confirm({
        title: this.$t('delete_confirm'),
        'z-index': 1000000,
        onOk: async () => {
          const { statusCode } = await deleteGraph([find.graphId])
          if (statusCode === 'OK') {
            this.$Notice.success({
              title: this.$t('delete_view'),
              desc: 'Success !'
            })
          }
        },
        onCancel: () => {}
      })
    },
    async addNewViewHandler () {
      const method = this.newView.isAdd ? addView : editView
      const roleRef = this.$refs.roleManagement
      this.newView.newViewForm.USE = roleRef.USE
      this.newView.newViewForm.MGMT = roleRef.MGMT
      const { statusCode, data } = await method([this.newView.newViewForm])
      if (statusCode === 'OK') {
        this.$Notice.success({
          title: this.$t('new_view'),
          desc: 'Success !'
        })
        this.currentView = data.viewId
        this.getGraphByView()
      }
    },
    cancelNewViewHandler () {
      this.newView.newViewForm = {
        viewId: '',
        name: '',
        report: '',
        editable: 'yes',
        suportVersion: 'yes',
        multiple: 'yes',
        filterAttr: '',
        filterValue: '',
        USE: [],
        MGMT: []
      }
    },
    async editView (view) {
      this.$refs.viewSelect.visible = false
      this.newView.isAdd = false
      const { statusCode, data } = await getViewById(view.viewId)
      if (statusCode === 'OK') {
        this.newView.newViewForm = { ...data }
        const res = await getReportListByPermission('MGMT', 'view')
        if (res.statusCode === 'OK') {
          this.newView.reportList = res.data
          this.getAttrByReportCi()
        }

        this.$refs.roleManagement.initRoleData(data.MGMT, data.USE)
        this.newView.showNewView = true
      }
    },
    async deleteView (val) {
      this.$refs.viewSelect.visible = false
      this.$Modal.confirm({
        title: this.$t('delete_confirm'),
        'z-index': 1000000,
        onOk: async () => {
          const { statusCode } = await deleteView([val.viewId])
          if (statusCode === 'OK') {
            this.$Notice.success({
              title: this.$t('delete_view'),
              desc: 'Success !'
            })
          }
        },
        onCancel: () => {}
      })
    },
    async getGraphByView () {
      this.currentGraph = ''
      const params = {
        paging: false,
        filters: [{ name: 'view', operator: 'eq', value: this.currentView }]
      }
      const { statusCode, data } = await getGraphByView(params)
      if (statusCode === 'OK') {
        this.showNewViewBtn = true
        this.graphs = data.contents.map(graph => {
          graph.graphElement = []
          return graph
        })
        if (data.contents.length > 0) {
          this.currentGraph = data.contents[0].name
          await this.getTabGraphData(data.contents[0].name)
        }
      }
    },
    async getTabGraphData (graphName) {
      const find = this.graphs.find(item => item.name === graphName)
      const report = this.viewOptions.find(item => item.viewId === this.currentView)
      const params = {
        paging: false,
        filters: [{ name: 'graph', operator: 'eq', value: find.graphId }],
        sorting: { asc: false, field: 'seqNo' }
      }
      const { statusCode, data } = await getElementByGraph(params)
      if (statusCode === 'OK') {
        find.graphElement = data.contents
      }
      this.$refs[graphName][0].initData(find.graphElement, find, report)
    },
    handleTabsAdd () {
      this.cancelNewGraphHandler()
      this.newGraph.showNewGraph = true
    },
    async addNewGraphHandler () {
      const find = this.viewOptions.find(item => item.viewId === this.currentView)
      this.newGraph.newGraphForm.view = find.viewId
      const { statusCode } = await addGraph([this.newGraph.newGraphForm])
      if (statusCode === 'OK') {
        this.$Notice.success({
          title: this.$t('new_graph'),
          desc: 'Success !'
        })
        this.getGraphByView()
      }
    },
    cancelNewGraphHandler () {
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
    },
    clearSelectView () {
      this.showNewViewBtn = false
      this.currentGraph = ''
      this.graphs = []
    },
    async getViewList (val) {
      if (val) {
        const params = {
          permission: 'MGMT'
        }
        const { data, statusCode } = await graphViews(params)
        if (statusCode === 'OK') {
          this.viewOptions = data
        }
      }
    }
  },
  components: {
    RoleConfig,
    ElementConfig
  }
}
</script>

<style scoped lang="scss">
.role-transfer-title {
  text-align: center;
  font-size: 13px;
  font-weight: 700;
  background-color: rgb(226, 222, 222);
  margin-bottom: 5px;
}
.modal-height {
  max-height: calc(100vh - 300px);
  overflow-y: auto;
}
</style>
