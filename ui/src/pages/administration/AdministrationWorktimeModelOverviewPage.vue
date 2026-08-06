<script setup lang="ts">
import BeeTimeClock from 'src/service/BeeTimeClock';
import { onMounted, ref } from 'vue';
import { WorkTimeModel, type WeekdayExceptionMap, type ApiWorkTimeModel } from 'src/models/WorkTimeModel';
import WorkTimeModelCreateDialog from 'components/dialog/WorkTimeModelCreateDialog.vue';
import type { QTableColumn } from 'quasar';
import { emptyPagination } from 'src/helper/objects';
import { useI18n } from 'vue-i18n';
import type { ErrorResponse } from 'src/models/Base';
import { showErrorMessage } from 'src/helper/message';

const { t } = useI18n();
const workTimeModels = ref<WorkTimeModel[]>([]);
const showDialog = ref(false);
const editingModel = ref<ApiWorkTimeModel | null>(null);

const weekdayNames: Record<number, string> = {
  0: t('LABEL_WEEKDAY_SUNDAY_SHORT'),
  1: t('LABEL_WEEKDAY_MONDAY_SHORT'),
  2: t('LABEL_WEEKDAY_TUESDAY_SHORT'),
  3: t('LABEL_WEEKDAY_WEDNESDAY_SHORT'),
  4: t('LABEL_WEEKDAY_THURSDAY_SHORT'),
  5: t('LABEL_WEEKDAY_FRIDAY_SHORT'),
  6: t('LABEL_WEEKDAY_SATURDAY_SHORT'),
};

function formatWeekdayExceptions(exceptions: WeekdayExceptionMap | null): string {
  if (!exceptions || Object.keys(exceptions).length === 0) {
    return '-';
  }
  return Object.entries(exceptions)
    .map(([day, hours]) => `${weekdayNames[parseInt(day)]}: ${hours}h`)
    .join(', ');
}

const columns: QTableColumn[] = [
  {
    name: 'id',
    field: 'ID',
    label: t('LABEL_ID'),
    align: 'left',
  },
  {
    name: 'name',
    field: 'Name',
    label: t('LABEL_NAME'),
    align: 'left',
  },
  {
    name: 'workingHoursPerWeekday',
    field: 'WorkingHoursPerWeekday',
    label: t('LABEL_WORKING_HOURS_PER_WEEKDAY'),
    align: 'left',
  },
  {
    name: 'hoursPerWeekdayException',
    field: 'HoursPerWeekdayException',
    label: t('LABEL_HOURS_PER_WEEKDAY_EXCEPTION'),
    align: 'left',
    format: (val: WeekdayExceptionMap | null) => formatWeekdayExceptions(val),
  },
  {
    name: 'overtimeSubtractionModel',
    field: 'OvertimeSubtractionModel',
    label: t('LABEL_OVERTIME_SUBTRACTION_MODEL'),
    align: 'left',
  },
  {
    name: 'overtimeSubtractionAmount',
    field: 'OvertimeSubtractionAmount',
    label: t('LABEL_OVERTIME_SUBTRACTION_AMOUNT'),
    align: 'left',
  },
  {
    name: 'holidayDaysPerYear',
    field: 'HolidayDaysPerYear',
    label: t('LABEL_HOLIDAY_DAYS_PER_YEAR'),
    align: 'left',
  },
  {
    name: 'overtimeWarningThreshold',
    field: 'OvertimeWarningThreshold',
    label: t('LABEL_OVERTIME_WARNING_THRESHOLD'),
    align: 'left',
  },
  {
    name: 'actions',
    field: '',
    label: t('LABEL_ACTION', 2),
    align: 'left',
  },
];

function loadWorkTimeModels() {
  BeeTimeClock.administrationGetWorkTimeModels()
    .then((result) => {
      if (result.status === 200) {
        workTimeModels.value = result.data.Data.map((s) =>
          WorkTimeModel.fromApi(s)
        );
      }
    })
    .catch((error: ErrorResponse) => {
      showErrorMessage(error.message);
    });
}

function openCreateDialog() {
  editingModel.value = null;
  showDialog.value = true;
}

function openEditDialog(model: WorkTimeModel) {
  editingModel.value = model;
  showDialog.value = true;
}

onMounted(() => {
  loadWorkTimeModels();
});
</script>

<template>
  <q-page padding>
    <q-btn
      :label="t('LABEL_CREATE', { item: t('LABEL_WORKTIME_MODEL') })"
      class="full-width"
      color="positive"
      icon="add"
      @click="openCreateDialog"
    />
    <q-table
      class="q-mt-md"
      :rows="workTimeModels"
      :columns="columns"
      hide-pagination
      :pagination="emptyPagination"
    >
      <template v-slot:header="props">
        <q-tr :props="props" class="bg-primary text-white">
          <q-th v-for="col in props.cols" :key="col.name" :props="props">
            {{ col.label }}
          </q-th>
        </q-tr>
      </template>

      <template v-slot:body="props">
        <q-tr :props="props">
          <q-td v-for="col in props.cols" :key="col.name" :props="props">
            <div v-if="col.name === 'actions'">
              <q-btn
                icon="edit"
                color="primary"
                dense
                flat
                @click="openEditDialog(props.row)"
              />
            </div>
            <div v-else>{{ col.value }}</div>
          </q-td>
        </q-tr>
      </template>
    </q-table>
    <WorkTimeModelCreateDialog
      v-model:show="showDialog"
      :work-time-model="editingModel"
      @created="loadWorkTimeModels"
      @updated="loadWorkTimeModels"
    />
  </q-page>
</template>

<style scoped></style>
