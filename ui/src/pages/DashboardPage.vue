<script setup lang="ts">
import {computed, onMounted, ref} from 'vue';
import BeeTimeClock from 'src/service/BeeTimeClock';
import type {User} from 'src/models/Authentication';
import {AbsenceSummaryItem} from 'src/models/Absence';
import AbsenceSummaryTableComponent from 'components/AbsenceSummaryTableComponent.vue';
import OvertimeTotal from 'components/OvertimeTotal.vue';
import OvertimeCurrentMonth from 'components/OvertimeMonth.vue';
import CheckInButton from 'components/CheckInButton.vue';
import type { ErrorResponse } from 'src/models/Base';
import { showErrorMessage } from 'src/helper/message';
import { useI18n } from 'vue-i18n';

const { t, locale } = useI18n({ useScope: 'global' });
const user = ref(null as User | null);
const absenceSummaryItems = ref([] as AbsenceSummaryItem[]);

const formattedDate = computed(() =>
  new Date().toLocaleDateString(locale.value, {
    weekday: 'long',
    year: 'numeric',
    month: 'long',
    day: 'numeric',
  })
);

function loadUser() {
  BeeTimeClock.getMeUser().then(result => {
    if (result.status === 200) user.value = result.data.Data;
  }).catch((error: ErrorResponse) => showErrorMessage(error.message));
}

function loadAbsenceSummary() {
  BeeTimeClock.queryAbsenceSummary().then(result => {
    if (result.status === 200)
      absenceSummaryItems.value = result.data.Data.map(s => AbsenceSummaryItem.fromApi(s));
  }).catch((error: ErrorResponse) => showErrorMessage(error.message));
}

onMounted(() => {
  loadUser();
  loadAbsenceSummary();
});
</script>

<template>
  <q-page>
    <div class="dashboard-header q-px-lg q-py-md row items-center">
      <div class="col">
        <div class="text-h6 text-weight-bold text-white" v-if="user">
          {{ t('LABEL_HELLO', { name: user.FirstName }) }}
        </div>
        <div class="text-caption text-white q-mt-xs" style="opacity: 0.75">
          {{ formattedDate }}
        </div>
      </div>
      <CheckInButton>
        <template #elapsed="{ elapsed }">
          <span class="text-white text-body1 text-weight-bold">{{ elapsed }}</span>
        </template>
      </CheckInButton>
    </div>

    <div class="q-pa-lg">
      <div class="row q-col-gutter-md q-mb-lg">
        <div class="col-12 col-sm-6">
          <OvertimeCurrentMonth class="full-height" />
        </div>
        <div class="col-12 col-sm-6">
          <OvertimeTotal class="full-height" />
        </div>
      </div>

      <q-card flat bordered>
        <AbsenceSummaryTableComponent v-model="absenceSummaryItems" />
      </q-card>
    </div>
  </q-page>
</template>

<style scoped>
.dashboard-header {
  background: linear-gradient(135deg, var(--q-primary) 0%, #1565c0 100%);
}
</style>
