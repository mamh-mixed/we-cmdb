<template>
  <div>
    <!-- 搜索区 -->
    <div>
      <Row>
        <Col span="7">
          <span>{{ $t('state_machine') }}</span>
          <Select
            v-model="stateMachine"
            clearable
            @on-clear="showTab = false"
            filterable
            style="width:300px; margin-left: 5px"
            ref="stateMachine"
            @on-open-change="getStateMachineList"
          >
            <Button type="success" style="width:100%" @click="addStateMachine" size="small">
              <Icon type="ios-add" size="24"></Icon>
            </Button>
            <Option v-for="item in stateMachineOptions" :value="item.id" :key="item.id"
              >{{ item.description
              }}<span style="float:right">
                <Button
                  @click.stop.prevent="editStateMachine(item)"
                  icon="ios-create-outline"
                  style="margin-right: 3px"
                  type="primary"
                  size="small"
                ></Button>
                <Button
                  @click.stop.prevent="deleteStateMachine(item)"
                  icon="ios-trash"
                  type="error"
                  size="small"
                ></Button> </span
            ></Option>
          </Select>
        </Col>

        <Col span="5" offset="1">
          <Button @click="getStateMachineInfo" :disabled="!stateMachine" type="primary">{{ $t('query') }}</Button>
          <Button style="margin: 0 10px" @click="preview" :disabled="!stateMachine">{{ $t('graphic_preview') }}</Button>
          <Button 
            @click="exportStateMachine" 
            class="btn-upload"
            :disabled="!stateMachine">
            <img src="@/styles/icon/DownloadOutlined.png" class="upload-icon" alt="" />
            {{ $t('export') }}
          </Button>
        </Col>
        <Col span="11">
          <Upload
            style="float:right"
            :action="uploadUrl"
            :show-upload-list="false"
            :max-size="1000"
            with-credentials
            :headers="{ Authorization: token }"
            :on-success="uploadSucess"
            :on-error="uploadFailed"
          >
            <Button class="btn-upload">
              <img src="@/styles/icon/UploadOutlined.png" class="upload-icon" />
              {{ $t('import') }}
            </Button>
          </Upload>
        </Col>
      </Row>
    </div>
    <!-- 展示区 -->
    <Tabs :value="currentTab" v-show="showTab">
      <TabPane :label="$t('state')" name="state">
        <StateTable ref="StateTable"></StateTable>
      </TabPane>
      <TabPane :label="$t('state_transition')" name="stateTransition">
        <StateTransition ref="StateTransition"></StateTransition>
      </TabPane>
    </Tabs>
    <Modal
      v-model="newStateMachine.isShow"
      :title="(newStateMachine.isAdd ? $t('new') : $t('edit')) + $t('state_machine')"
      @on-ok="confirmStateMachine"
      @on-cancel="newStateMachine.isShow = false"
    >
      <Form inline :label-width="80">
        <FormItem label="ID" v-if="newStateMachine.isAdd">
          <Input type="text" v-model="newStateMachine.form.id" style="width:400px"></Input>
        </FormItem>
        <FormItem :label="$t('description')">
          <Input type="text" v-model="newStateMachine.form.description" style="width:400px"></Input>
        </FormItem>
      </Form>
    </Modal>

    <Modal
      v-model="importInfo.isShow"
      :title="$t('confirm_import')"
      @on-ok="okImport"
      @on-cancel="importInfo.isShow = false"
      width="800"
    >
      <div style="overflow: auto;max-height: 500px;">
        <pre>{{ importInfo.new_states }}</pre>
        <Divider plain>new & old</Divider>
        <pre>{{ importInfo.old_states }}</pre>
      </div>
    </Modal>
    <Modal v-model="graphicPreview" width="800" :fullscreen="fullscreen" footer-hide :title="$t('graphic_preview')">
      <p slot="header">
        <span>{{ $t('details') }}</span>
        <Icon
          v-if="!fullscreen"
          @click="changeView"
          style="float: right;margin-right: 30px;margin-top: 3px;"
          type="ios-expand"
        />
        <Icon v-else @click="changeView" style="float: right;margin-right: 30px;margin-top: 3px;" type="ios-contract" />
      </p>
      <div class="graph-container" style="height:100%;width:100%" id="graph"></div>
    </Modal>
  </div>
</template>

<script>
import * as d3 from 'd3-selection'
// eslint-disable-next-line
import * as d3Graphviz from 'd3-graphviz'
import { getCookie } from '@/pages/util/cookie'
import StateTable from './state-management'
import StateTransition from './state-transition'
import axios from 'axios'
import {
  getStateMachineList,
  getStateMachineFullInfo,
  confirmImportStateMachine,
  editState,
  addStateMachine,
  editStateMachine,
  deleteStateMachine,
  addState,
  deleteState,
  getState
} from '@/api/server'

export const custom_api_enum = [
  {
    url: '/wecmdb/api/v1/state-config/state-machine/query',
    method: 'get'
  },
  {
    url: '/wecmdb/api/v1/state-config/state-machine/import',
    method: 'post'
  }
]

export default {
  components: { StateTable, StateTransition },
  data() {
    return {
      fullscreen: false,
      graphicPreview: false,
      currentTab: 'state',
      showTab: false,
      newStateMachine: {
        isShow: false,
        isAdd: true,
        form: {
          id: '',
          description: ''
        }
      },
      importInfo: {
        isShow: false,
        new_states: null,
        old_states: null,
        data: null
      },
      token: '',
      uploadUrl: '/wecmdb/api/v1/state-config/state-machine/import',
      stateMachine: '',
      stateMachineOptions: [],

      tableColumnsState: [
        {
          title: this.$t('table_name'),
          render: (h, params) => {
            return (
              <div>
                <Input
                  disabled={params.row.id !== ''}
                  value={params.row.name}
                  onInput={v => {
                    params.row.name = v
                  }}
                />
              </div>
            )
          }
        },
        {
          title: this.$t('is_confirm'),
          render: (h, params) => {
            return (
              <div>
                <Select
                  value={params.row.isConfirm}
                  style="width:100%"
                  filterable
                  on-on-change={v => {
                    params.row.isConfirm = v
                  }}
                >
                  {this.optionSet.map(item => {
                    return (
                      <Option value={item.value} key={item.value}>
                        {item.label}
                      </Option>
                    )
                  })}
                </Select>
              </div>
            )
          }
        },
        {
          title: this.$t('unique_path_trigger'),
          render: (h, params) => {
            return (
              <div>
                <Select
                  value={params.row.uniquePathTrigger}
                  style="width:100%"
                  filterable
                  on-on-change={v => {
                    params.row.uniquePathTrigger = v
                  }}
                >
                  {this.optionSet.map(item => {
                    return (
                      <Option value={item.value} key={item.value}>
                        {item.label}
                      </Option>
                    )
                  })}
                </Select>
              </div>
            )
          }
        },
        {
          title: this.$t('description'),
          render: (h, params) => {
            return (
              <div>
                <Input
                  value={params.row.description}
                  onInput={v => {
                    params.row.description = v
                  }}
                />
              </div>
            )
          }
        },
        {
          title: this.$t('actions'),
          width: 150,
          align: 'center',
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
                      this.saveRow(params.row)
                      // this.show(params.index)
                    }
                  }
                },
                this.$t('save')
              ),
              h(
                'Button',
                {
                  props: {
                    type: 'error',
                    size: 'small'
                  },
                  on: {
                    click: () => {
                      this.removeRow(params.row, params.index)
                    }
                  }
                },
                this.$t('delete')
              )
            ])
          }
        }
      ],
      optionSet: [
        { label: 'yes', value: 'yes' },
        { label: 'no', value: 'no' }
      ],
      tableDataState: [],
      emptyState: {
        description: '',
        id: '',
        isConfirm: 'yes',
        name: '',
        stateMachine: '',
        uniquePathTrigger: 'yes'
      },
      tableColumnsStateTran: [],
      tableDataStateTran: [],
      emptyStateTransition: {
        action: '',
        currentState: '',
        guid: '',
        operation: '',
        operationFormType: '',
        operationMultiple: '',
        operation_en: '',
        permission: '',
        stateMachine: '',
        targetState: ''
      },
      graph: {}
    }
  },
  mounted() {
    this.token = getCookie('accessToken')
  },
  methods: {
    changeView() {
      this.fullscreen = !this.fullscreen
      this.preview()
    },
    async uploadSucess(val) {
      if (val.data.diff_flag) {
        this.importInfo.data = val.data
        this.importInfo.new_states = val.data.new_states
        this.importInfo.old_states = val.data.old_states
        this.importInfo.isShow = true
      } else {
        this.successTip()
      }
    },
    async okImport() {
      const { statusCode } = await confirmImportStateMachine(this.importInfo.data)
      if (statusCode === 'OK') {
        this.successTip()
      }
    },
    successTip() {
      this.$Notice.success({
        title: 'Successful',
        desc: 'Successful'
      })
    },
    uploadFailed(val) {
      this.$Message.warning(this.$t('tips.failed'))
    },
    addStateMachine() {
      this.$refs.stateMachine.visible = false
      this.newStateMachine = {
        isShow: true,
        isAdd: true,
        form: {
          id: '',
          description: ''
        }
      }
    },
    editStateMachine(item) {
      this.$refs.stateMachine.visible = false
      this.newStateMachine = {
        isShow: true,
        isAdd: false,
        form: {
          ...item
        }
      }
    },
    async confirmStateMachine() {
      const method = this.newStateMachine.isAdd ? addStateMachine : editStateMachine
      const { statusCode } = await method([this.newStateMachine.form])
      if (statusCode === 'OK') {
        this.successTip()
        await this.getStateMachineList()
        this.stateMachine = this.newStateMachine.form.id
        this.getStateMachineInfo()
      }
    },
    deleteStateMachine(item) {
      this.$refs.stateMachine.visible = false
      this.$Modal.confirm({
        title: this.$t('confirm_delete'),
        'z-index': 1000000,
        loading: true,
        onOk: async () => {
          let res = await deleteStateMachine(item.id)
          this.$Modal.remove()
          if (res.statusCode === 'OK') {
            this.successTip()
            this.stateMachine = ''
            this.showTab = false
            this.getStateMachineList()
          }
        },
        onCancel: () => {}
      })
    },
    removeRow(item, index) {
      this.$Modal.confirm({
        title: this.$t('confirm_delete'),
        'z-index': 1000000,
        loading: true,
        onOk: async () => {
          let res = await deleteState(item.id)
          this.$Modal.remove()
          if (res.statusCode === 'OK') {
            this.successTip()
            this.tableDataState.splice(index, 1)
          }
        },
        onCancel: () => {}
      })
    },
    newState() {
      let tmp = JSON.parse(JSON.stringify(this.emptyState))
      tmp.stateMachine = this.stateMachine
      this.tableDataState.push(tmp)
    },
    async saveRow(item) {
      const method = item.id === '' ? addState : editState
      const { statusCode } = await method([item])
      if (statusCode === 'OK') {
        this.successTip()
        const params = {
          filters: [
            {
              name: 'stateMachine',
              operator: 'eq',
              value: this.stateMachine
            }
          ]
        }
        const stateRes = await getState(params)
        if (stateRes.statusCode === 'OK') {
          this.tableDataState = stateRes.data.contents
        }
      }
    },
    async getStateMachineList() {
      const params = {
        filters: []
      }
      const { statusCode, data } = await getStateMachineList(params)
      if (statusCode === 'OK') {
        this.stateMachineOptions = data.contents
      }
    },
    exportStateMachine() {
      axios({
        method: 'GET',
        url: `/wecmdb/api/v1/state-config/state-machine/query?export=yes&id=${this.stateMachine}`,
        headers: {
          Authorization: getCookie('accessToken')
        }
      })
        .then(response => {
          if (response.status < 400) {
            let content = JSON.stringify(response.data)
            let fileName = `state_machine_${new Date().getFullYear() +
              '-' +
              new Date().getMonth() +
              '-' +
              new Date().getDay() +
              '_' +
              new Date().getHours() +
              ':' +
              new Date().getMinutes() +
              ':' +
              new Date().getSeconds()}.json`
            let blob = new Blob([content])
            if ('msSaveOrOpenBlob' in navigator) {
              window.navigator.msSaveOrOpenBlob(blob, fileName)
            } else {
              if ('download' in document.createElement('a')) {
                // 非IE下载
                let elink = document.createElement('a')
                elink.download = fileName
                elink.style.display = 'none'
                elink.href = URL.createObjectURL(blob)
                document.body.appendChild(elink)
                elink.click()
                URL.revokeObjectURL(elink.href) // 释放URL 对象
                document.body.removeChild(elink)
              } else {
                // IE10+下载
                navigator.msSaveOrOpenBlob(blob, fileName)
              }
            }
          }
        })
        .catch(() => {
          this.$Message.warning('Error')
        })
    },
    async getStateMachineInfo() {
      this.showTab = true
      this.$refs.StateTable.initData(this.stateMachine)
      this.$refs.StateTransition.initData(this.stateMachine)
    },
    async preview() {
      this.graphicPreview = true
      const { statusCode, data } = await getStateMachineFullInfo(this.stateMachine)
      if (statusCode === 'OK') {
        const Dot = this.generateDot(data.states, data.transitions)
        this.newInitGraph(Dot)
      }
    },
    async newInitGraph(Dot) {
      let graph
      const initEvent = () => {
        graph = d3.select('#graph')
        graph
          .on('dblclick.zoom', null)
          .on('wheel.zoom', null)
          .on('mousewheel.zoom', null)

        this.graph.graphviz = graph
          .graphviz()
          .zoom(true)
          .fit(true)
          .height('100%')
          .width('100%')
          .attributer(function(d) {
            if (d.attributes.class === 'edge') {
              const keys = d.key.split('->')
              const from = keys[0].trim()
              const to = keys[1].trim()
              d.attributes.from = from
              d.attributes.to = to
            }

            if (d.tag === 'text') {
              const key = d.children[0].text
              d3.select(this).attr('text-key', key)
            }
          })
      }
      initEvent()
      this.renderGraph(Dot)
    },
    renderGraph(Dot) {
      this.graph.graphviz
        .transition()
        .renderDot(Dot)
        .on('end', () => {})
    },
    generateDot(states, transitions) {
      let dot = 'digraph{bgcolor="transparent"; ranksep=1.1; nodesep=.7; size="11,8"; rankdir=TB\n'
      dot += 'Node [fontname=Arial; style=filled; fixedsize="true"; width="1.5"; height="0.6"; fontsize=15];\n'
      dot += 'Edge [fontname=Arial; minlen="1"; color="gray"; fontsize=10];\n'

      states.forEach(state => {
        dot +=
          '"' + state.id + '" [id="' + state.id + '"; label="' + state.name + '"; tooltip="' + state.description + '"; '
        if (state.name === 'null_0' || state.name === 'null_1') {
          dot += 'shape="ellipse"; '
        } else {
          dot += 'shape="box"; '
        }
        if (state.isConfirm === 'no') {
          dot += 'color=mediumseagreen;]\n'
        } else {
          dot += 'color=dodgerblue2;]\n'
        }
      })
      transitions.forEach(transition => {
        dot +=
          '"' +
          transition.currentState +
          '"->"' +
          transition.targetState +
          '" [label="O:' +
          transition.operation_en +
          '\nA:' +
          transition.action +
          '\nP:' +
          transition.permission +
          '";]\n'
      })
      dot += '}'
      return dot
    }
  }
}
</script>
