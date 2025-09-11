<script setup lang="ts">
import BeeTimeClock from 'src/service/BeeTimeClock';
import { onMounted, ref } from 'vue';
import { showErrorMessage } from 'src/helper/message';
import { emptyPagination } from 'src/helper/objects';
import { date, type QTableColumn } from 'quasar';
import { useI18n } from 'vue-i18n';
import formatDate = date.formatDate;
import TimestampCorrectionDialog from 'components/TimestampCorrectionDialog.vue';
import AbsenceCreateDialog from 'components/absence/AbsenceCreateDialog.vue';

interface MissingDay {
  date: string;
}

const { t } = useI18n();
const loading = ref(true);
const missingDays = ref<MissingDay[]>([]);
const promptTimestampCorrectionCreate = ref(false);
const promptAbsenceCreationDialog = ref(false);
const selectedDate = ref<string>();
const columns = ref<QTableColumn[]>([
  {
    name: 'date',
    label: `${t('LABEL_DATE')}`,
    field: 'date',
    sortable: true,
    format: (val) => formatDate(val, 'ddd. DD.MM.YYYY'),
    align: 'left',
  },
  {
    name: 'actions',
    label: t('LABEL_ACTION', 2),
    field: '',
    align: 'right',
  },
]);

function loadMissingDays() {
  BeeTimeClock.getMissingDays()
    .then((result) => {
      if (result.status === 200) {
        missingDays.value = result.data.Data.map((s) => {
          return { date: s } as MissingDay;
        }).sort(
          (a, b) => new Date(b.date).getTime() - new Date(a.date).getTime(),
        );
      }
    })
    .catch((error) => {
      showErrorMessage(error);
    })
    .finally(() => {
      loading.value = false;
    });
}

function openTimestampCorrectionDialog(input: string) {
  selectedDate.value = input;
  promptTimestampCorrectionCreate.value = true;
}

function openAbsenceCreationDialog(input: string) {
  selectedDate.value = input;
  promptAbsenceCreationDialog.value = true;
}

onMounted(() => {
  loadMissingDays();
});
</script>

<template>
  <q-page padding>
    <div v-if="!loading">
      <q-table
        :columns="columns"
        :rows="missingDays"
        hide-pagination
        :pagination="emptyPagination"
      >
        <template v-slot:header="props">
          <q-tr :props="props">
            <q-th v-for="col in props.cols" :key="col.name" :props="props">
              {{ col.label }}
              <template v-if="col.name == 'date'">
                <q-chip :label="missingDays.length" dense />
              </template>
            </q-th>
          </q-tr>
        </template>
        <template v-slot:body="props">
          <q-tr :props="props">
            <q-td v-for="col in props.cols" :key="col.name" :props="props">
              {{ col.value }}
              <template v-if="col.name == 'actions'">
                <q-btn
                  class="q-ml-md"
                  color="primary"
                  icon="timer"
                  :label="t('LABEL_TIMESTAMP_CORRECTION_CREATE')"
                  @click="openTimestampCorrectionDialog(props.row.date)"
                />
                <q-btn
                  class="q-ml-md"
                  color="secondary"
                  icon="event"
                  :label="t('LABEL_CREATE', { item: t('LABEL_ABSENCE') })"
                  @click="openAbsenceCreationDialog(props.row.date)"
                />
              </template>
            </q-td>
          </q-tr>
        </template>
      </q-table>
      <TimestampCorrectionDialog
        v-if="selectedDate"
        v-model:show="promptTimestampCorrectionCreate"
        :init-date="new Date(selectedDate)"
      />
      <AbsenceCreateDialog
        v-if="selectedDate"
        v-model:show="promptAbsenceCreationDialog"
        :init-date="new Date(selectedDate)"
        @create="loadMissingDays()"
      />
    </div>
    <q-inner-loading :showing="loading" />
  </q-page>
</template>

<style scoped></style>
