<template>
  <q-page padding>
    <div v-if="!isLoading">

      <div class="row">
        <div class="col">
          <OvertimeTotal class="full-height" />
        </div>
      </div>
      <div class="row">
        <div class="col q-pa-md">
          <q-btn class="full-width" color="positive" :label="$t('BTN_CHECK_IN')"
                 @click="actionCheckIn()" />
        </div>
        <div class="col q-pa-md">
          <q-btn class="full-width" color="negative" :label="$t('BTN_CHECK_OUT')" @click="actionCheckOut" />
        </div>
        <div class="col q-pa-md">
          <q-btn class="full-width" color="primary" :label="$t('BTN_ADD', {item: $t('LABEL_TIMESTAMP')})" @click="promptTimestampCorrectionCreate = true"/>
        </div>
      </div>
      <div class="row">
        <div class="col">
          <q-select v-model="selectedYear" :options="timestampYears" :label="$t('LABEL_YEAR')" />
        </div>
        <div class="col">
          <q-select class="q-ml-md" v-model="selectedMonth" :options="timestampMonths" :label="$t('LABEL_MONTH')" />
        </div>
      </div>
      <div class="row q-mt-md">
        <div class="col">
          <OvertimeMonth v-model:model-month="selectedMonth" v-model:model-year="selectedYear" class="full-width" />
        </div>
      </div>
      <div class="q-pt-lg">
        <WorktimeOverviewTable v-model="timestampCurrentMonthGrouped" @create="loadTimestampGrouped()"/>
      </div>
      <TimestampCorrectionDialog v-model:model-show="promptTimestampCorrectionCreate" @refresh="loadTimestampGrouped()"/>
    </div>
    <q-inner-loading :showing="isLoading" />
  </q-page>
</template>

<script lang="ts" setup>
import { showErrorMessage, showInfoMessage } from 'src/helper/message';
import {
  TimestampGroup,
  TimestampYearMonthGrouped
} from 'src/models/Timestamp';
import BeeTimeClock from 'src/service/BeeTimeClock';
import { computed, onMounted, ref, watch } from 'vue';
import { ErrorResponse } from 'src/models/Base';
import OvertimeTotal from 'components/OvertimeTotal.vue';
import WorktimeOverviewTable from 'components/WorktimeOverviewTable.vue';
import OvertimeMonth from 'components/OvertimeMonth.vue';
import TimestampCorrectionDialog from 'components/TimestampCorrectionDialog.vue';
import { useI18n } from 'vue-i18n';

const { t } = useI18n();

const promptTimestampCorrectionCreate = ref(false);
const timestampCurrentMonthGrouped = ref<TimestampGroup[]>([]);
const expanded = ref(['']);
const timestampYearMonths = ref<TimestampYearMonthGrouped>({});
const selectedYear = ref<number>(new Date().getFullYear());
const selectedMonth = ref<number>(new Date().getMonth() + 1);
const isLoading = ref(true);

const timestampYears = computed(() => {
  if (!timestampYearMonths.value) return [];
  const years = Object.keys(timestampYearMonths.value);
  return years.sort();
});

const timestampMonths = computed(() => {
  if (!timestampYearMonths.value) return [];
  const months = timestampYearMonths.value[selectedYear.value];
  return months.sort();
});

function actionCheckIn(isHomeoffice?: boolean) {
  BeeTimeClock.timestampActionCheckin(isHomeoffice).then((result) => {
    if (result.status === 201) {
      showInfoMessage(t('MSG_CHECK_IN_SUCCESS'));
      loadTimestampGrouped();
    }
  }).catch((error: ErrorResponse) => {
    showErrorMessage(error.response?.data.Message);
  });
}

function actionCheckOut(isHomeoffice?: boolean) {
  BeeTimeClock.timestampActionCheckout(isHomeoffice).then((result) => {
    if (result.status === 200) {
      showInfoMessage(t('MSG_CHECK_OUT_SUCCESS'));
      loadTimestampGrouped();
    }
  }).catch((error: ErrorResponse) => {
    showErrorMessage(error.response?.data.Message);
  });
}

function loadTimestampGrouped() {
  BeeTimeClock.timestampQueryMonthGrouped(selectedYear.value, selectedMonth.value).then((result) => {
    if (result.status === 200) {
      timestampCurrentMonthGrouped.value = result.data.Data.sort((a, b) => new Date(b.Date).getTime() - new Date(a.Date).getTime());
      if (timestampCurrentMonthGrouped.value.length > 0) {
        expanded.value = [timestampCurrentMonthGrouped.value[0].Date.toString()];
      }
    }
  });
}

async function loadTimestampMonths() {
  isLoading.value = true;

  const result = await BeeTimeClock.timestampQueryMonths();

  if (result.status === 200) {
    timestampYearMonths.value = result.data.Data;
  }

  isLoading.value = false;
}

onMounted(async () => {
  await loadTimestampMonths();
  loadTimestampGrouped();
});

watch(selectedYear, () => {
  if (timestampYearMonths.value[selectedYear.value].includes(selectedMonth.value)) {
    loadTimestampGrouped();
    return;
  } else {
    selectedMonth.value = timestampYearMonths.value[selectedYear.value][0];
  }
});

watch(selectedMonth, () => {
  console.log('month changed');
  loadTimestampGrouped();
});


</script>
