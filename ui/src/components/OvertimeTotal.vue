<script setup lang="ts">
import {ref} from 'vue';
import BeeTimeClock from 'src/service/BeeTimeClock';
import { formatIndustryHourMinutes } from 'src/helper/formatter';
import type { ErrorResponse } from 'src/models/Base';
import { showErrorMessage } from 'src/helper/message';
import { useI18n } from 'vue-i18n';

const { t } = useI18n();
const value = ref(0);

const props = defineProps({
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

BeeTimeClock.overtimeTotal().then(result => {
  if (result.status === 200) value.value = result.data.Data.Total;
}).catch((error: ErrorResponse) => showErrorMessage(error.message));
</script>

<template>
  <q-item v-if="props.item">
    <q-item-section avatar>
      <q-icon name="account_balance_wallet" color="grey-5" />
    </q-item-section>
    <q-item-section>
      <q-item-label caption>{{ t('LABEL_OVERTIME_TOTAL') }}</q-item-label>
    </q-item-section>
    <q-item-section side>
      <span class="text-body1 text-weight-bold" :class="value >= 0 ? 'text-positive' : 'text-negative'">
        {{ toHHMM(value) }}
      </span>
    </q-item-section>
  </q-item>

  <q-card v-else flat bordered class="full-height">
    <q-card-section :class="props.dense ? 'q-pa-sm' : 'q-pa-md'">
      <div class="row items-center no-wrap">
        <div class="col">
          <div class="text-overline text-grey-6">{{ t('LABEL_OVERTIME_TOTAL') }}</div>
          <div
            :class="[props.dense ? 'text-h5' : 'text-h4', 'text-weight-bold q-mt-xs', value >= 0 ? 'text-positive' : 'text-negative']"
          >
            {{ toHHMM(value) }}
          </div>
          <div class="text-caption text-grey-5 q-mt-xs">
            {{ formatIndustryHourMinutes(value) }}
          </div>
        </div>
        <q-icon v-if="!props.dense" name="account_balance_wallet" size="48px" color="grey-4" />
      </div>
    </q-card-section>
  </q-card>
</template>
