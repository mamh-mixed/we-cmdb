<template>
  <div>
    <!-- 展示区 -->
    <Button style="margin-bottom: 5px" type="primary" @click="newStateTran">{{ $t('new') }}</Button>
    <Table border :columns="tableColumnsStateTran" :data="tableDataStateTran" :max-height="MODALHEIGHT"></Table>
  </div>
</template>

<script>
import { getState, getStateTran, editStateTran, addStateTran, deleteStateTran } from '@/api/server'
export default {
  data () {
    return {
      MODALHEIGHT: 500,
      tableColumnsStateTran: [
        {
          title: this.$t('current_state'),
          key: 'currentState',
          render: (h, params) => {
            return (
              <div>
                <Select
                  value={params.row.currentState}
                  style="width:100%"
                  filterable
                  clearable
                  on-on-change={v => {
                    params.row.currentState = v
                  }}
                >
                  {this.state.map(item => {
                    return (
                      <Option value={item.id} key={item.id}>
                        {item.name}
                      </Option>
                    )
                  })}
                </Select>
              </div>
            )
          }
        },
        {
          title: this.$t('target_state'),
          key: 'targetState',
          render: (h, params) => {
            return (
              <div>
                <Select
                  value={params.row.targetState}
                  style="width:100%"
                  filterable
                  clearable
                  on-on-change={v => {
                    params.row.targetState = v
                  }}
                >
                  {this.state.map(item => {
                    return (
                      <Option value={item.id} key={item.id}>
                        {item.name}
                      </Option>
                    )
                  })}
                </Select>
              </div>
            )
          }
        },
        {
          title: this.$t('operation'),
          key: 'operation',
          render: (h, params) => {
            return (
              <div>
                <Input
                  value={params.row.operation}
                  clearable
                  onInput={v => {
                    params.row.operation = v
                  }}
                />
              </div>
            )
          }
        },
        {
          title: this.$t('operation_en'),
          key: 'operation_en',
          render: (h, params) => {
            return (
              <div>
                <Input
                  clearable
                  value={params.row.operation_en}
                  onInput={v => {
                    params.row.operation_en = v
                  }}
                />
              </div>
            )
          }
        },
        {
          title: this.$t('operation_form_type'),
          key: 'operationFormType',
          render: (h, params) => {
            return (
              <div>
                <Select
                  value={params.row.operationFormType}
                  style="width:100%"
                  filterable
                  clearable
                  on-on-change={v => {
                    params.row.operationFormType = v
                  }}
                >
                  {this.formTypes.map(item => {
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
          title: this.$t('permission'),
          key: 'permission',
          render: (h, params) => {
            return (
              <div>
                <Select
                  value={params.row.permission}
                  style="width:100%"
                  filterable
                  clearable
                  on-on-change={v => {
                    params.row.permission = v
                  }}
                >
                  {this.permissionOptions.map(item => {
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
          title: this.$t('action'),
          key: 'action',
          render: (h, params) => {
            return (
              <div>
                <Select
                  value={params.row.action}
                  style="width:100%"
                  filterable
                  clearable
                  on-on-change={v => {
                    params.row.action = v
                  }}
                >
                  {this.actionOptions.map(item => {
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
          title: this.$t('operation_multiple'),
          key: 'operationMultiple',
          render: (h, params) => {
            return (
              <div>
                <Select
                  value={params.row.operationMultiple}
                  style="width:100%"
                  filterable
                  clearable
                  on-on-change={v => {
                    params.row.operationMultiple = v
                  }}
                >
                  {this.YorN.map(item => {
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
          title: this.$t('operating_area'),
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
      formTypes: [
        { label: 'editable_form', value: 'editable_form' },
        { label: 'select_form', value: 'select_form' },
        { label: 'confirm_form', value: 'confirm_form' }
      ],
      permissionOptions: [
        { label: '新增 insert', value: 'insert' },
        { label: '执行 execute', value: 'execute' },
        { label: '修改 update', value: 'update' },
        { label: '删除 delete', value: 'delete' },
        { label: '确认 confirm', value: 'comfirm' }
      ],
      actionOptions: [
        { label: 'insert', value: 'insert' },
        { label: 'update', value: 'update' },
        { label: 'delete', value: 'delete' },
        { label: 'confirm', value: 'confirm' },
        { label: 'execute', value: 'execute' }
      ],
      YorN: [
        { label: 'yes', value: 'yes' },
        { label: 'no', value: 'no' }
      ],
      tableDataStateTran: [],
      emptyState: {
        action: 'insert',
        currentState: '',
        guid: '',
        operation: '',
        operationFormType: '',
        operationMultiple: 'no',
        operation_en: '',
        permission: 'insert',
        stateMachine: '',
        targetState: ''
      },
      state: []
    }
  },
  mounted () {
    this.MODALHEIGHT = document.body.scrollHeight - 300
  },
  methods: {
    async getState () {
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
        this.state = stateRes.data.contents
      }
    },
    removeRow (item, index) {
      this.$Modal.confirm({
        title: this.$t('confirm_delete'),
        'z-index': 1000000,
        loading: true,
        onOk: async () => {
          let res = await deleteStateTran(item.guid)
          this.$Modal.remove()
          if (res.statusCode === 'OK') {
            this.$Notice.success({
              title: 'Successful',
              desc: 'Successful'
            })
            this.tableDataStateTran.splice(index, 1)
          }
        },
        onCancel: () => {}
      })
    },
    newStateTran () {
      let tmp = JSON.parse(JSON.stringify(this.emptyState))
      tmp.stateMachine = this.stateMachine
      this.tableDataStateTran.push(tmp)
    },
    async initData (stateMachine) {
      this.stateMachine = stateMachine
      await this.getState()
      this.getTableData(this.stateMachine)
    },
    async getTableData (stateMachine) {
      const params = {
        filters: [
          {
            name: 'stateMachine',
            operator: 'eq',
            value: stateMachine
          }
        ]
      }
      const stateRes = await getStateTran(params)
      if (stateRes.statusCode === 'OK') {
        this.tableDataStateTran = stateRes.data.contents
      }
    },
    async saveRow (item) {
      const method = item.guid === '' ? addStateTran : editStateTran
      const { statusCode } = await method([item])
      if (statusCode === 'OK') {
        this.$Notice.success({
          title: 'Successful',
          desc: 'Successful'
        })
        this.getTableData(this.stateMachine)
      }
    }
  }
}
</script>
