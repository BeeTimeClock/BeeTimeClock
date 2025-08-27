<script setup lang="ts">
import BeeTimeClock from 'src/service/BeeTimeClock';
import { showErrorMessage, showInfoMessage } from 'src/helper/message';
import type { ErrorResponse } from 'src/models/Base';

function notifyCurrentWeek() {
  BeeTimeClock.administrationNotifyAbsenceWeek()
    .then((result) => {
      if (result.status === 204) {
        showInfoMessage('notify sent');
      }
    })
    .catch((error: ErrorResponse) => {
      showErrorMessage(error.message);
    });
}
</script>

<template>
  <q-page padding>
    <q-btn
      color="primary"
      class="full-width"
      label="Notify Absence Current Week"
      @click="notifyCurrentWeek"
    />
  </q-page>
</template>

<style scoped></style>
