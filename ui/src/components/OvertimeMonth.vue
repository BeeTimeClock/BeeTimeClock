<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue';
import BeeTimeClock from 'src/service/BeeTimeClock';
import { formatIndustryHourMinutes } from 'src/helper/formatter';
import type { OvertimeResponse } from 'src/models/Timestamp';
import type { ErrorResponse } from 'src/models/Base';
import { showErrorMessage } from 'src/helper/message';
import { useI18n } from 'vue-i18n';

const { t, locale } = useI18n();

const overtimeResponse = ref<OvertimeResponse | null>(null);

const props = defineProps({
  modelYear: { type: Number, default: new Date().getFullYear() },
  modelMonth: { type: Number, default: new Date().getMonth() + 1 },
  modelUserId: { type: Number },
  modelTeamId: { type: Number },
  dense: { type: Boolean, default: false },
  item: { type: Boolean, default: false },
});

function toHHMM(hours: number): string {
  const sign = hours < 0 ? '-' : '+';
  const abs = Math.abs(hours);
  const h = Math.floor(abs);
  const m = Math.round((abs - h) * 60);
  return `${sign}${h}:${String(m).padStart(2, '0')}`;
}

const monthName = computed(() =>
  new Date(props.modelYear, props.modelMonth - 1).toLocaleString(locale.value, { month: 'long', year: 'numeric' })
);

const netOvertime = computed(() => {
  if (!overtimeResponse.value) return 0;
  return overtimeResponse.value.Total - overtimeResponse.value.Subtracted;
});

function loadOvertime() {
  if (props.modelUserId && props.modelUserId != 0) {
    if (props.modelTeamId && props.modelTeamId != 0) {
      BeeTimeClock.teamTimestampQueryMonthOvertime(props.modelTeamId, props.modelUserId, props.modelYear, props.modelMonth)
        .then(result => { if (result.status === 200) overtimeResponse.value = result.data.Data; })
        .catch((error: ErrorResponse) => showErrorMessage(error.message));
    } else {
      BeeTimeClock.administrationTimestampQueryMonthOvertime(props.modelUserId, props.modelYear, props.modelMonth)
        .then(result => { if (result.status === 200) overtimeResponse.value = result.data.Data; })
        .catch((error: ErrorResponse) => showErrorMessage(error.message));
    }
  } else {
    BeeTimeClock.timestampQueryMonthOvertime(props.modelYear, props.modelMonth)
      .then(result => { if (result.status === 200) overtimeResponse.value = result.data.Data; })
      .catch((error: ErrorResponse) => showErrorMessage(error.message));
  }
}

watch(props, () => loadOvertime());
onMounted(() => loadOvertime());
</script>

<template>
  <q-item v-if="item">
    <q-item-section avatar>
      <q-icon name="timer" color="grey-5" />
    </q-item-section>
    <q-item-section>
      <q-item-label caption>{{ monthName }}</q-item-label>
    </q-item-section>
    <q-item-section side v-if="overtimeResponse">
      <span class="text-body1 text-weight-bold" :class="netOvertime >= 0 ? 'text-positive' : 'text-negative'">
        {{ toHHMM(netOvertime) }}
      </span>
    </q-item-section>
  </q-item>

  <q-card v-else flat bordered class="full-height">
    <q-card-section :class="dense ? 'q-pa-sm' : 'q-pa-md'">
      <div class="row items-center no-wrap q-mb-md">
        <div class="col">
          <div class="text-overline text-grey-6">{{ t('LABEL_OVERTIME') }}</div>
          <div class="text-h6 text-weight-bold">{{ monthName }}</div>
        </div>
        <q-icon v-if="!dense" name="timer" size="48px" color="grey-4" />
      </div>

      <template v-if="overtimeResponse">
        <template v-if="overtimeResponse.Subtracted !== 0">
          <div class="row items-center q-mb-xs">
            <div class="col text-body2 text-grey-7">{{ t('LABEL_WORKING_HOURS') }}</div>
            <div
              class="col-auto text-body1 text-weight-medium"
              :class="overtimeResponse.Total >= 0 ? 'text-positive' : 'text-negative'"
            >
              {{ toHHMM(overtimeResponse.Total) }}
            </div>
          </div>
          <div class="row items-center q-mb-sm">
            <div class="col text-body2 text-grey-7">{{ t('LABEL_SUBTRACTED_HOURS') }}</div>
            <div class="col-auto text-body1 text-weight-medium text-negative">
              {{ formatIndustryHourMinutes(overtimeResponse.Subtracted) }}
            </div>
          </div>
          <q-separator class="q-mb-sm" />
        </template>
        <div class="row items-center">
          <div class="col text-body2 text-grey-7 text-weight-bold">{{ t('LABEL_OVERTIME_HOURS') }}</div>
          <div
            class="col-auto text-h6 text-weight-bold"
            :class="netOvertime >= 0 ? 'text-positive' : 'text-negative'"
          >
            {{ toHHMM(netOvertime) }}
          </div>
        </div>
      </template>

      <q-inner-loading :showing="!overtimeResponse">
        <q-spinner-dots size="40px" color="primary" />
      </q-inner-loading>
    </q-card-section>
  </q-card>
</template>
