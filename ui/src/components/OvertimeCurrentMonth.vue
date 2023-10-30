<script setup lang="ts">
import {ref} from 'vue';
import BeeTimeClock from 'src/service/BeeTimeClock';

const value = ref(0);

BeeTimeClock.timestampQueryCurrentMonthOvertime().then(result => {
  if (result.status === 200) {
    value.value = result.data.Data.Total;
  }
})

</script>

<template>
  <q-card>
    <q-card-section class="bg-primary text-white text-subtitle">
      {{ $t('LABEL_OVERTIME_CURRENT_MONTH') }}
    </q-card-section>
    <q-card-section class="text-h6 text-center">
      {{ value }}
    </q-card-section>
  </q-card>
</template>

<style scoped>

</style>
