<script setup lang="ts">
import { onMounted, ref } from 'vue';
import { OvertimeMonthQuota } from 'src/models/Overtime';
import BeeTimeClock from 'src/service/BeeTimeClock';
import { useI18n } from 'vue-i18n';

const { t } = useI18n();
const overtimeQuotas = ref<OvertimeMonthQuota[]>([]);

const columns = [
  {
    name: 'Date',
    required: true,
    label: t('LABEL_DATE'),
    align: 'left',
    field: 'identifier',
    sortable: true,
  },
  {
    name: 'Hours',
    required: true,
    label: t('LABEL_HOUR', 2),
    field: 'Hours',
    align: 'right',
    sortable: true,
    format: (val: number) => `${val.toFixed(2)}`,
  },
  {
    name: 'actions',
    label: t('LABEL_ACTION', 2),
  },
];

function loadOvertimeQuotas() {
  BeeTimeClock.overtimeMonthQuotas().then((result) => {
    if (result.status === 200) {
      overtimeQuotas.value = result.data.Data.sort((a, b) => b.Year - a.Year || b.Month - a.Month).map((s) =>
        OvertimeMonthQuota.fromApi(s)
      );
    }
  });
}

function calculateOvertimeMonth(overtimeMonthQuota: OvertimeMonthQuota) {
  BeeTimeClock.calculateOvertimeMonthQuota(
    overtimeMonthQuota.Year,
    overtimeMonthQuota.Month
  ).then((result) => {
    if (result.status === 200) {
      loadOvertimeQuotas();
    }
  });
}

onMounted(() => {
  loadOvertimeQuotas();
});
</script>

<template>
  <q-page padding>
    <q-table
      flat
      bordered
      :rows="overtimeQuotas"
      :columns="columns"
      row-key="identifier"
    >
      <template v-slot:header="props">
        <q-tr :props="props">
          <q-th auto-width />
          <q-th v-for="col in props.cols" :key="col.name" :props="props">
            {{ col.label }}
          </q-th>
        </q-tr>
      </template>

      <template v-slot:body="props">
        <q-tr :props="props">
          <q-td auto-width>
            <q-btn
              size="xs"
              color="secondary"
              round
              dense
              @click="props.expand = !props.expand"
              :icon="props.expand ? 'remove' : 'add'"
            />
          </q-td>
          <q-td v-for="col in props.cols" :key="col.name" :props="props">
            <template v-if="col.name === 'actions'">
              <q-btn icon="refresh" @click="calculateOvertimeMonth(props.row)" color="primary"/>
            </template>
            <template v-else>
              {{ col.value }}
            </template>
          </q-td>
        </q-tr>
        <q-tr v-show="props.expand" :props="props">
          <q-td colspan="100%">
            <q-list separator>
              <q-item
                v-for="entry in props.row.Summary"
                :key="entry.Identifier"
              >
                <q-item-section>
                  <q-item-label caption>{{ entry.Source }}</q-item-label>
                </q-item-section>
                <q-item-section>
                  <q-item-section>{{ entry.Value.toFixed(2) }}</q-item-section>
                </q-item-section>
              </q-item>
            </q-list>
          </q-td>
        </q-tr>
      </template>
    </q-table>
  </q-page>
</template>

<style scoped></style>
