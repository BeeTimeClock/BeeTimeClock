<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref } from 'vue';
import BeeTimeClock from 'src/service/BeeTimeClock';
import type { Timestamp } from 'src/models/Timestamp';
import type { ErrorResponse } from 'src/models/Base';
import { showErrorMessage, showInfoMessage } from 'src/helper/message';
import { useI18n } from 'vue-i18n';

const emit = defineEmits<{ checkin: []; checkout: [] }>();

const { t } = useI18n();
const lastTimestamp = ref<Timestamp | null>(null);
const now = ref(new Date());
let ticker: ReturnType<typeof setInterval> | null = null;

const isCheckedIn = computed(() =>
  lastTimestamp.value !== null &&
  new Date(lastTimestamp.value.GoingTimestamp as unknown as string).getFullYear() < 1999
);

const elapsedTime = computed(() => {
  if (!isCheckedIn.value || !lastTimestamp.value) return null;
  const diffMs = now.value.getTime() - new Date(lastTimestamp.value.ComingTimestamp).getTime();
  const totalMinutes = Math.floor(diffMs / 60000);
  const h = Math.floor(totalMinutes / 60);
  const m = totalMinutes % 60;
  return `${h}:${String(m).padStart(2, '0')}`;
});

function loadLastTimestamp() {
  BeeTimeClock.timestampQueryLast().then(result => {
    if (result.status === 200) lastTimestamp.value = result.data.Data;
  }).catch(() => { lastTimestamp.value = null; });
}

function actionCheckIn() {
  BeeTimeClock.timestampActionCheckin().then(result => {
    if (result.status === 201) {
      showInfoMessage(t('MSG_CHECK_IN_SUCCESS'));
      loadLastTimestamp();
      emit('checkin');
    }
  }).catch((error: ErrorResponse) => showErrorMessage(error.response?.data.Message));
}

function actionCheckOut() {
  BeeTimeClock.timestampActionCheckout(false).then(result => {
    if (result.status === 200) {
      showInfoMessage(t('MSG_CHECK_OUT_SUCCESS'));
      loadLastTimestamp();
      emit('checkout');
    }
  }).catch((error: ErrorResponse) => showErrorMessage(error.response?.data.Message));
}

onMounted(() => {
  loadLastTimestamp();
  ticker = setInterval(() => { now.value = new Date(); }, 60000);
});

onUnmounted(() => {
  if (ticker) clearInterval(ticker);
});
</script>

<template>
  <div v-if="isCheckedIn" class="row items-center q-gutter-sm">
    <slot name="elapsed" :elapsed="elapsedTime">
      <span>{{ elapsedTime }}</span>
    </slot>
    <q-btn
      color="negative"
      :label="t('BTN_CHECK_OUT')"
      icon="logout"
      unelevated
      @click="actionCheckOut"
    />
  </div>
  <q-btn
    v-else
    color="positive"
    :label="t('BTN_CHECK_IN')"
    icon="login"
    unelevated
    @click="actionCheckIn"
  />
</template>
