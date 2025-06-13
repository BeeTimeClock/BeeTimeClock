<script setup lang="ts">
import { onMounted, ref } from 'vue';
import { Timestamp } from 'src/models/Timestamp';
import BeeTimeClock from 'src/service/BeeTimeClock';
import { useI18n } from 'vue-i18n';
import { date } from 'quasar';
import formatDate = date.formatDate;

const timestamps = ref<Timestamp[]>([]);
const { t } = useI18n();

const columns = [
  {
    name: 'expand'
  },
  {
    name: 'coming',
    field: 'ComingTimestamp',
    label: t('LABEL_COMING_TIMESTAMP'),
    format: (val: Date) => formatDate(val, 'DD.MM.YYYY HH:mm:ss'),
  },
  {
    name: 'going',
    field: 'GoingTimestamp',
    label: t('LABEL_GOING_TIMESTAMP'),
    format: (val: Date) => formatDate(val, 'DD.MM.YYYY HH:mm:ss'),
  },
  {
    name: 'hasCorrections',
    label: t('LABEL_HAS_CORRECTIONS'),
  },
  {
    name: 'actions'
  },
];

const columnsCorrections = [
  {
    name: 'changeReason',
    field: 'ChangeReason',
    label: t('LABEL_REASON'),
  },
  {
    name: 'oldComingTimestamp',
    field: 'OldComingTimestamp',
    label: t('LABEL_COMING_TIMESTAMP'),
    format: (val: Date) => formatDate(val, 'DD.MM.YYYY HH:mm:ss'),
  },
  {
    name: 'oldGoingTimestamp',
    field: 'OldGoingTimestamp',
    label: t('LABEL_GOING_TIMESTAMP'),
    format: (val: Date) => formatDate(val, 'DD.MM.YYYY HH:mm:ss'),
  },
];

function loadTimestamps() {
  BeeTimeClock.timestampQuerySuspicious().then((result) => {
    if (result.status === 200) {
      timestamps.value = result.data.Data;
    }
  });
}

onMounted(() => {
  loadTimestamps();
});
</script>

<template>
  <q-page padding>
    <q-table :rows="timestamps" :columns="columns">
      <template v-slot:body="props">
        <q-tr :props="props">
          <q-td v-for="col in props.cols" :key="col.name" :props="props">
            <q-td auto-width v-if="col.name == 'expand'">
              <q-btn
                v-if="props.row.Corrections.length > 0"
                size="sm"
                color="accent"
                round
                dense
                @click="props.expand = !props.expand"
                :icon="props.expand ? 'remove' : 'add'"
              />
            </q-td>
            <div v-else-if="col.name == 'actions'">
              
            </div>
            <div v-else-if="col.name == 'hasCorrections'">
              <q-icon
                size="large"
                :name="
                  props.row.Corrections.length > 0 ? 'check_circle' : 'cancel'
                "
                :color="props.row.Corrections.length > 0 ? 'positive' : ''"
              />
            </div>
            <div v-else>
              {{ col.value }}
            </div>
          </q-td>
        </q-tr>
        <q-tr v-show="props.expand" :props="props">
          <q-td colspan="100%">
            <q-table :rows="props.row.Corrections" :columns="columnsCorrections"> </q-table>
          </q-td>
        </q-tr>
      </template>
    </q-table>
  </q-page>
</template>

<style scoped></style>
