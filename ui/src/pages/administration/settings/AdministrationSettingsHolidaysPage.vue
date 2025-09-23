<script setup lang="ts">
import { useI18n } from 'vue-i18n';
import BeeTimeClock from 'src/service/BeeTimeClock';
import { onMounted, ref } from 'vue';
import { HolidayCustom } from 'src/models/Holiday';
import type { ErrorResponse } from 'src/models/Base';
import { showErrorMessage } from 'src/helper/message';

const {t} = useI18n();
const holidaysCustom = ref<HolidayCustom[]>([]);


function loadHolidaysCustom() {
  BeeTimeClock.administrationHolidaysCustom().then(result => {
    if (result.status === 200) {
      holidaysCustom.value = result.data.Data.map(s => HolidayCustom.fromApi(s));
    }
  }).catch((error: ErrorResponse) => {
    showErrorMessage(error.message);
  });
}

onMounted(() => {
  loadHolidaysCustom();
})
</script>

<template>
  <q-page padding>
    <q-card>
      <q-card-section class="bg-primary text-white text-subtitle2">
        {{ t('LABEL_HOLIDAY', 2) }}
      </q-card-section>
      <q-card-section>
        <q-table :rows="holidaysCustom"/>
      </q-card-section>
    </q-card>
  </q-page>
</template>

<style scoped></style>
