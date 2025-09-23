<script setup lang="ts">
import { onMounted, ref } from 'vue';
import { AbsenceReason, AbsenceReasonImpact } from 'src/models/Absence';
import { useI18n } from 'vue-i18n';
import BeeTimeClock from 'src/service/BeeTimeClock';
import { showErrorMessage } from 'src/helper/message';
import type { ErrorResponse } from 'src/models/Base';
import AbsenceReasonUpdateDialog from 'components/dialog/AbsenceReasonUpdateDialog.vue';

const { t } = useI18n();
const absenceReasons = ref<AbsenceReason[]>([]);
const selectedAbsenceReason = ref<AbsenceReason>();
const showDialog = ref(false);


function loadAbsenceReasons() {
  BeeTimeClock.absenceReasons()
    .then((result) => {
      if (result.status === 200) {
        absenceReasons.value = result.data.Data.map((s) =>
          AbsenceReason.fromApi(s),
        );
      }
    })
    .catch((error: ErrorResponse) => {
      showErrorMessage(error.message);
    });
}

function createAbsenceReason() {
  selectedAbsenceReason.value = new AbsenceReason();
  showDialog.value = true;
}

function editAbsenceReason(absenceReason: AbsenceReason) {
  selectedAbsenceReason.value = absenceReason;
  showDialog.value = true;
}

onMounted(() => {
  loadAbsenceReasons();
});
</script>

<template>
  <q-list separator>
    <q-item v-for="absenceReason in absenceReasons" :key="absenceReason.ID">
      <q-item-section>
        <q-item-label>{{ absenceReason.Description }}</q-item-label>
        <q-item-label caption v-if="absenceReason.Impact">{{
          absenceReason.Impact
        }}</q-item-label>
        <q-item-label
          caption
          v-if="absenceReason.Impact == AbsenceReasonImpact.Hours"
          >{{ absenceReason.ImpactHours
          }}{{ t('LABEL_HOUR', absenceReason.ImpactHours) }}</q-item-label
        >
      </q-item-section>
      <q-item-section side>
        <q-btn
          icon="edit"
          color="primary"
          @click="editAbsenceReason(absenceReason)"
        />
      </q-item-section>
    </q-item>
    <q-item>
      <q-item-section>
        <q-btn
          class="full-width"
          icon="add"
          :label="t('BTN_ADD')"
          color="positive"
          @click="createAbsenceReason"
        />
      </q-item-section>
    </q-item>
  </q-list>
  <AbsenceReasonUpdateDialog v-if="selectedAbsenceReason" v-model="selectedAbsenceReason" v-model:show="showDialog"/>
</template>

<style scoped></style>
