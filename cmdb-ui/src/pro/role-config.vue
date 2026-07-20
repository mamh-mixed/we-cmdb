<template>
  <div>
    <div style="width: 100%; overflow-x: auto">
      <div style="min-width: 670px;">
        <div class="role-transfer-title">{{ $t('mgmt_role') }}</div>
        <Transfer
          :titles="transferTitles"
          :list-style="transferStyle"
          :data="allRoles"
          :target-keys="MGMT"
          @on-change="handleMgmtRoleTransferChange"
          filterable
        ></Transfer>
        <div style="margin-top: 30px">
          <div class="role-transfer-title">{{ $t('use_role') }}</div>
          <Transfer
            :titles="transferTitles"
            :list-style="transferStyle"
            :data="allRolesBackUp"
            :target-keys="USE"
            @on-change="handleUseRoleTransferChange"
            filterable
          ></Transfer>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { getRolesByCurrentUser, getAllRoles } from '@/api/server'
export default {
  name: '',
  data () {
    return {
      allRoles: [],
      MGMT: [],
      allRolesBackUp: [],
      USE: [],
      transferTitles: [this.$t('unselected_role'), this.$t('selected_role')],
      transferStyle: { width: '300px' }
    }
  },
  methods: {
    initRoleData (MGMT, USE) {
      this.USE = USE
      this.MGMT = MGMT
      this.getRolesByCurrentUser()
      this.getRoleList()
    },
    async getRolesByCurrentUser () {
      const { statusCode, data } = await getRolesByCurrentUser()
      if (statusCode === 'OK') {
        this.allRoles = data.map(_ => {
          return {
            ..._,
            key: _.roleName,
            label: _.description
          }
        })
      }
    },
    async getRoleList () {
      const { statusCode, data } = await getAllRoles()
      if (statusCode === 'OK') {
        this.allRolesBackUp = data.map(_ => {
          return {
            ..._,
            key: _.roleName,
            label: _.description
          }
        })
      }
    },
    async handleMgmtRoleTransferChange (newTargetKeys, direction, moveKeys) {
      this.MGMT = newTargetKeys
    },
    async handleUseRoleTransferChange (newTargetKeys, direction, moveKeys) {
      this.USE = newTargetKeys
    }
  },
  components: {}
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
</style>
