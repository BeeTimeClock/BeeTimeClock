<script setup lang="ts">
import { ref, computed } from 'vue';
import type {
  ApiWorkTimeModelCreateRequest,
  ApiWorkTimeModelUpdateRequest,
  WeekdayExceptionMap,
  ApiWorkTimeModel,
} from 'src/models/WorkTimeModel';
import BeeTimeClock from 'src/service/BeeTimeClock';
import { showErrorMessage, showInfoMessage } from 'src/helper/message';
import { useI18n } from 'vue-i18n';
import type { ErrorResponse } from 'src/models/Base';

const props = defineProps<{
  workTimeModel?: ApiWorkTimeModel | null;
}>();

const show = defineModel('show', { type: Boolean, default: false });
const { t } = useI18n();

const emits = defineEmits(['created', 'updated']);

const isEditMode = computed(() => !!props.workTimeModel);

const overtimeSubtractionModelOptions = [
  { value: 'hours', label: t('LABEL_HOUR', 2) },
  { value: 'percentage', label: t('LABEL_PERCENTAGE') },
];

// Weekday labels (0=Sunday, 1=Monday, ..., 6=Saturday)
const weekdays = [
  { value: 1, label: t('LABEL_WEEKDAY_MONDAY') },
  { value: 2, label: t('LABEL_WEEKDAY_TUESDAY') },
  { value: 3, label: t('LABEL_WEEKDAY_WEDNESDAY') },
  { value: 4, label: t('LABEL_WEEKDAY_THURSDAY') },
  { value: 5, label: t('LABEL_WEEKDAY_FRIDAY') },
  { value: 6, label: t('LABEL_WEEKDAY_SATURDAY') },
  { value: 0, label: t('LABEL_WEEKDAY_SUNDAY') },
];

// Store weekday exceptions as individual refs for easier binding
const weekdayExceptions = ref<Record<number, number | null>>({
  0: null,
  1: null,
  2: null,
  3: null,
  4: null,
  5: null,
  6: null,
});

const formData = ref({
  Name: '',
  WorkingHoursPerWeekday: 8,
  OvertimeSubtractionModel: 'percentage' as 'hours' | 'percentage',
  OvertimeSubtractionAmount: 10,
  HolidayDaysPerYear: 30,
  OvertimeWarningThreshold: 0,
});

// Build the exception map from non-null values
const hoursPerWeekdayException = computed<WeekdayExceptionMap | null>(() => {
  const exceptions: WeekdayExceptionMap = {};
  let hasExceptions = false;
  for (const [key, value] of Object.entries(weekdayExceptions.value)) {
    if (value !== null && value !== undefined && value !== 0) {
      exceptions[parseInt(key)] = value;
      hasExceptions = true;
    }
  }
  return hasExceptions ? exceptions : null;
});

function initForm() {
  if (props.workTimeModel) {
    // Edit mode - populate form with existing data
    formData.value = {
      Name: props.workTimeModel.Name,
      WorkingHoursPerWeekday: props.workTimeModel.WorkingHoursPerWeekday,
      OvertimeSubtractionModel: props.workTimeModel.OvertimeSubtractionModel,
      OvertimeSubtractionAmount: props.workTimeModel.OvertimeSubtractionAmount,
      HolidayDaysPerYear: props.workTimeModel.HolidayDaysPerYear,
      OvertimeWarningThreshold: props.workTimeModel.OvertimeWarningThreshold,
    };
    // Populate weekday exceptions
    const exceptions = props.workTimeModel.HoursPerWeekdayException || {};
    weekdayExceptions.value = {
      0: exceptions[0] ?? null,
      1: exceptions[1] ?? null,
      2: exceptions[2] ?? null,
      3: exceptions[3] ?? null,
      4: exceptions[4] ?? null,
      5: exceptions[5] ?? null,
      6: exceptions[6] ?? null,
    };
  } else {
    // Create mode - reset to defaults
    formData.value = {
      Name: '',
      WorkingHoursPerWeekday: 7.6,
      OvertimeSubtractionModel: 'percentage',
      OvertimeSubtractionAmount: 10,
      HolidayDaysPerYear: 30,
      OvertimeWarningThreshold: 0,
    };
    weekdayExceptions.value = {
      0: null,
      1: null,
      2: null,
      3: null,
      4: null,
      5: null,
      6: null,
    };
  }
}

function submitForm() {
  if (isEditMode.value) {
    updateWorkTimeModel();
  } else {
    createWorkTimeModel();
  }
}

function createWorkTimeModel() {
  const request: ApiWorkTimeModelCreateRequest = {
    ...formData.value,
    HoursPerWeekdayException: hoursPerWeekdayException.value,
  };

  BeeTimeClock.administrationCreateWorkTimeModel(request)
    .then((result) => {
      if (result.status == 201) {
        show.value = false;
        emits('created', result.data.Data);
        showInfoMessage(t('MSG_CREATE_SUCCESS', { item: t('LABEL_WORKTIME_MODEL') }));
      }
    })
    .catch((error: ErrorResponse) => {
      showErrorMessage(error.message);
    });
}

function updateWorkTimeModel() {
  if (!props.workTimeModel) return;

  const request: ApiWorkTimeModelUpdateRequest = {
    WorkingHoursPerWeekday: formData.value.WorkingHoursPerWeekday,
    OvertimeSubtractionModel: formData.value.OvertimeSubtractionModel,
    OvertimeSubtractionAmount: formData.value.OvertimeSubtractionAmount,
    HoursPerWeekdayException: hoursPerWeekdayException.value,
    HolidayDaysPerYear: formData.value.HolidayDaysPerYear,
    OvertimeWarningThreshold: formData.value.OvertimeWarningThreshold,
  };

  BeeTimeClock.administrationUpdateWorkTimeModel(props.workTimeModel.ID, request)
    .then((result) => {
      if (result.status == 200) {
        show.value = false;
        emits('updated', result.data.Data);
        showInfoMessage(t('MSG_UPDATE_SUCCESS', { item: t('LABEL_WORKTIME_MODEL') }));
      }
    })
    .catch((error: ErrorResponse) => {
      showErrorMessage(error.message);
    });
}
</script>

<template>
  <q-dialog v-model="show" @before-show="initForm">
    <q-card style="min-width: 500px">
      <q-card-section class="bg-primary text-white text-h6">
        <template v-if="isEditMode">
          {{ t('TITLE_UPDATE', { item: t('LABEL_WORKTIME_MODEL'), identifier: formData.Name }) }}
        </template>
        <template v-else>
          {{ t('LABEL_CREATE', { item: t('LABEL_WORKTIME_MODEL') }) }}
        </template>
      </q-card-section>
      <q-form @submit="submitForm">
        <q-card-section>
          <q-input
            v-model="formData.Name"
            :label="t('LABEL_NAME')"
            :rules="[(val) => !!val || t('LABEL_FIELD_REQUIRED')]"
            :readonly="isEditMode"
            :disable="isEditMode"
          />
          <q-input
            v-model.number="formData.WorkingHoursPerWeekday"
            type="number"
            step="0.1"
            :label="t('LABEL_WORKING_HOURS_PER_WEEKDAY')"
            :rules="[(val) => val > 0 || t('LABEL_FIELD_REQUIRED')]"
          />

          <div class="q-mt-md">
            <div class="text-subtitle2 q-mb-sm">{{ t('LABEL_HOURS_PER_WEEKDAY_EXCEPTION') }}</div>
            <div class="text-caption text-grey q-mb-sm">{{ t('LABEL_HOURS_PER_WEEKDAY_EXCEPTION_HINT') }}</div>
            <div class="row q-col-gutter-sm">
              <div v-for="weekday in weekdays" :key="weekday.value" class="col-6">
                <q-input
                  v-model.number="weekdayExceptions[weekday.value]"
                  type="number"
                  step="0.1"
                  :label="weekday.label"
                  dense
                  clearable
                />
              </div>
            </div>
          </div>

          <q-select
            v-model="formData.OvertimeSubtractionModel"
            :label="t('LABEL_OVERTIME_SUBTRACTION_MODEL')"
            :options="overtimeSubtractionModelOptions"
            emit-value
            map-options
            class="q-mt-md"
          />
          <q-input
            v-model.number="formData.OvertimeSubtractionAmount"
            type="number"
            step="0.1"
            :label="t('LABEL_OVERTIME_SUBTRACTION_AMOUNT')"
          />
          <q-input
            v-model.number="formData.HolidayDaysPerYear"
            type="number"
            :label="t('LABEL_HOLIDAY_DAYS_PER_YEAR')"
          />
          <q-input
            v-model.number="formData.OvertimeWarningThreshold"
            type="number"
            :label="t('LABEL_OVERTIME_WARNING_THRESHOLD')"
          />
        </q-card-section>
        <q-card-section>
          <q-card-actions>
            <q-btn color="negative" v-close-popup :label="t('LABEL_CANCEL')" />
            <q-btn
              color="positive"
              type="submit"
              :label="isEditMode ? t('LABEL_SAVE') : t('LABEL_CREATE')"
            />
          </q-card-actions>
        </q-card-section>
      </q-form>
    </q-card>
  </q-dialog>
</template>

<style scoped></style>
