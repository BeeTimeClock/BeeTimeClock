<script setup lang="ts">
import BeeTimeClock from 'src/service/BeeTimeClock';
import { showInfoMessage } from 'src/helper/message';

function notifyCurrentWeek() {
  BeeTimeClock.administrationNotifyAbsenceWeek().then((result) => {
    if (result.status === 204) {
      showInfoMessage('notify sent');
    }
  });
}

</script>

<template>
  <q-page padding>
    <q-btn color="primary" class="full-width" label="Notify Absence Current Week" @click="notifyCurrentWeek" />
  </q-page>
</template>

<style scoped></style>
