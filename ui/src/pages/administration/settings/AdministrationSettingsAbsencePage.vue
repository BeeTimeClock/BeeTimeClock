<script setup lang="ts">
import AbsenceReasonAdministrationTable from 'components/AbsenceReasonAdministrationTable.vue';
import { useI18n } from 'vue-i18n';
import BeeTimeClock from 'src/service/BeeTimeClock';
import { showErrorMessage, showInfoMessage } from 'src/helper/message';
import type { ErrorResponse } from 'src/models/Base';

const { t } = useI18n();

function recalculateAbsences() {
  BeeTimeClock.administrationAbsenceRecalculate()
    .then((result) => {
      if (result.status === 200) {
        showInfoMessage(t('MSG_CREATE_SUCCESS'));
      }
    })
    .catch((error: ErrorResponse) => {
      showErrorMessage(error.message);
    });
}
</script>

<template>
  <q-page padding>
    <q-card>
      <q-card-section class="bg-primary text-white text-subtitle2">
        {{ t('LABEL_ABSENCE_REASON', 2) }}
      </q-card-section>
      <q-card-section>
        <AbsenceReasonAdministrationTable />
      </q-card-section>
    </q-card>

    <q-card>
      <q-card-section class="bg-negative text-white text-subtitle2 q-mt-lg">
        {{ t('LABEL_DANGER_ZONE') }}
      </q-card-section>
      <q-card-section>
        <q-btn
          :label="t('LABEL_RECALCULATE_ABSENCES')"
          color="primary"
          class="full-width"
          @click="recalculateAbsences"
        />
      </q-card-section>
    </q-card>
  </q-page>
</template>

<style scoped></style>
