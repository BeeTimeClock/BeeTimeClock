<script setup lang="ts">
import {ref} from 'vue';
import BeeTimeClock from 'src/service/BeeTimeClock';
import { formatIndustryHourMinutes } from 'src/helper/formatter';

const value = ref(0);

BeeTimeClock.overtimeTotal().then(result => {
  if (result.status === 200) {
    value.value = result.data.Data.Total;
  }
})

</script>

<template>
  <q-card>
    <q-card-section class="bg-primary text-white text-subtitle">
      {{ $t('LABEL_OVERTIME_TOTAL') }}
    </q-card-section>
    <q-card-section class="text-h6 text-center">
      {{ formatIndustryHourMinutes(value) }}
    </q-card-section>
  </q-card>
</template>

<style scoped>

</style>
