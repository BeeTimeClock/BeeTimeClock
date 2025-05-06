<script setup lang="ts">
import { onMounted, ref } from 'vue';
import { OvertimeMonthQuota } from 'src/models/Overtime';
import BeeTimeClock from 'src/service/BeeTimeClock';

const overtimeQuotas = ref<OvertimeMonthQuota[]>([]);

function loadOvertimeQuotas() {
  BeeTimeClock.overtimeMonthQuotas().then(result => {
    if (result.status === 200) {
      overtimeQuotas.value = result.data.Data.map(s => OvertimeMonthQuota.fromApi(s))
    }
  })
}

function calculateOvertimeMonth(overtimeMonthQuota: OvertimeMonthQuota) {
  BeeTimeClock.calculateOvertimeMonthQuota(overtimeMonthQuota.Year, overtimeMonthQuota.Month).then(result => {
    if (result.status === 200) {
      loadOvertimeQuotas();
    }
  })
}

onMounted(() => {
  loadOvertimeQuotas();
})
</script>

<template>
<q-page padding>
  <q-table :rows="overtimeQuotas">
    <template v-slot:header="props">
      <q-tr :props="props">
        <q-th
          v-for="col in props.cols"
          :key="col.name"
          :props="props"
        >
          {{ col.label }}
        </q-th>
        <q-th auto-width />
      </q-tr>
    </template>
    <template v-slot:body="props">
      <q-tr :props="props">
        <q-td
          v-for="col in props.cols"
          :key="col.name"
          :props="props"
        >
          <div v-if="col.name == 'IsHomeoffice'">
            <q-icon size="large" :name="props.row.IsHomeoffice ? 'check_circle' : 'cancel'"
                    :color="props.row.IsHomeoffice ? 'positive' : ''" />
          </div>
          <div v-else>
            {{ col.value }}
          </div>
        </q-td>
        <q-td auto-width>
          <q-btn icon="refresh" color="primary" @click="calculateOvertimeMonth(props.row as OvertimeMonthQuota)"/>
        </q-td>
      </q-tr>
    </template>
  </q-table>
</q-page>
</template>

<style scoped>

</style>
