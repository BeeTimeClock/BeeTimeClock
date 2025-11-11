<script setup lang="ts">
import { useI18n } from 'vue-i18n';
import BeeTimeClock from 'src/service/BeeTimeClock';

const {t} = useI18n();
const holidays = BeeTimeClock.getHolidaysYear()
</script>

<template>
  <q-page padding>
    <q-card>
      <q-card-section class="bg-primary text-white text-subtitle2">
        {{ t('LABEL_HOLIDAY', 2) }}
      </q-card-section>
      <q-card-section>
      </q-card-section>
    </q-card>
  </q-page>
</template>

<style scoped></style>
