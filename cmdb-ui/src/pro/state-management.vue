<template>
  <div>
    <!-- 展示区 -->
    <Button style="margin-bottom: 5px" type="primary" @click="newState">{{ $t('new') }}</Button>
    <Table border :columns="tableColumnsState" :data="tableDataState" :max-height="MODALHEIGHT"></Table>
  </div>
</template>

<script>
import { editState, addState, deleteState, getState } from '@/api/server'
export default {
  data () {
    return {
      MODALHEIGHT: 500,
      tableColumnsState: [
        {
          title: this.$t('table_name'),
          render: (h, params) => {
            return (
              <div>
                <Input
                  disabled={params.row.id !== '' || params.row.name.startsWith('null_')}
                  value={params.row.name}
                  clearable
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
                  disabled={params.row.name.startsWith('null_')}
                  value={params.row.isConfirm}
                  style="width:100%"
                  filterable
                  clearable
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
                  disabled={params.row.name.startsWith('null_')}
                  value={params.row.uniquePathTrigger}
                  style="width:100%"
                  filterable
                  clearable
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
                  disabled={params.row.name.startsWith('null_')}
                  value={params.row.description}
                  clearable
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
      }
    }
  },
  mounted () {
    this.MODALHEIGHT = document.body.scrollHeight - 300
  },
  methods: {
    removeRow (item, index) {
      this.$Modal.confirm({
        title: this.$t('confirm_delete'),
        'z-index': 1000000,
        loading: true,
        onOk: async () => {
          let res = await deleteState(item.id)
          this.$Modal.remove()
          if (res.statusCode === 'OK') {
            this.$Notice.success({
              title: 'Successful',
              desc: 'Successful'
            })
            this.tableDataState.splice(index, 1)
          }
        },
        onCancel: () => {}
      })
    },
    newState () {
      let tmp = JSON.parse(JSON.stringify(this.emptyState))
      tmp.stateMachine = this.stateMachine
      this.tableDataState.push(tmp)
    },
    initData (stateMachine) {
      this.stateMachine = stateMachine
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
      const stateRes = await getState(params)
      if (stateRes.statusCode === 'OK') {
        this.tableDataState = stateRes.data.contents
      }
    },
    async saveRow (item) {
      if (!item.name.endsWith('_0') && !item.name.endsWith('_1')) {
        this.$Notice.warning({
          title: 'Error',
          desc: "Name must be end width '_0' or '_1'"
        })
        return
      }
      const method = item.id === '' ? addState : editState
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
