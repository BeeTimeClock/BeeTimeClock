<script setup lang="ts">
import AbsenceSummaryTableComponent from 'components/AbsenceSummaryTableComponent.vue';
import type {AbsenceCreateRequest, AbsenceReason} from 'src/models/Absence';
import {type AbsenceSummaryItem} from 'src/models/Absence';
import {computed, ref} from 'vue';
import BeeTimeClock from 'src/service/BeeTimeClock';
import type {ErrorResponse} from 'src/models/Base';
import {showErrorMessage, showInfoMessage} from 'src/helper/message';
import {useI18n} from 'vue-i18n';
import {date} from 'quasar';

const {t} = useI18n();
const show = defineModel('show', {default: false});
const emit = defineEmits<{
  create: [],
}>()
const initDate = defineModel<Date>('init-date', {
  default: new Date(),
});

const absenceCreateRequest = ref<AbsenceCreateRequest>(
  {} as AbsenceCreateRequest);
const absenceSummaryItems = ref([] as AbsenceSummaryItem[]);
const absenceReasons = ref<AbsenceReason[]>([]);

function loadAbsenceSummary() {
  BeeTimeClock.queryAbsenceSummary()
    .then((result) => {
      if (result.status === 200) {
        absenceSummaryItems.value = result.data.Data;
      }
    })
    .catch((error: ErrorResponse) => {
      showErrorMessage(error.message);
    });
}

const conflictingAbsences = computed(() => {
  return absenceSummaryItems.value.filter(
    (s) =>
      (s.AbsenceFrom > new Date(absenceCreateRequest.value.AbsenceFrom) &&
        s.AbsenceFrom < new Date(absenceCreateRequest.value.AbsenceTill)) ||
      (s.AbsenceTill < new Date(absenceCreateRequest.value.AbsenceTill) &&
        s.AbsenceTill > new Date(absenceCreateRequest.value.AbsenceFrom)),
  );
});

function loadAbsenceReasons() {
  BeeTimeClock.absenceReasons()
    .then((result) => {
      if (result.status === 200) {
        absenceReasons.value = result.data.Data;
      }
    })
    .catch((error: ErrorResponse) => {
      showErrorMessage(error.message);
    });
}

function createAbsence() {
  if (!absenceCreateRequest.value) return;
  BeeTimeClock.createAbsence(absenceCreateRequest.value)
    .then((result) => {
      if (result.status === 201) {
        emit('create')
        absenceCreateRequest.value = {} as AbsenceCreateRequest;
        showInfoMessage(t('MSG_CREATE_SUCCESS'));
      }
    })
    .catch((error: ErrorResponse) => {
      showErrorMessage(error.message);
    });
}

function onBeforeShow() {
  loadAbsenceReasons();
  loadAbsenceSummary();

  console.log(initDate.value.toDateString());
  absenceCreateRequest.value.AbsenceFrom = date.formatDate(initDate.value, 'YYYY-MM-DD');
  console.log(absenceCreateRequest.value);
}
</script>

<template>
  <q-dialog persistent v-model="show" @before-show="onBeforeShow">
    <q-card class="q-dialog-plugin full-width">
      <q-card-section>
        <q-input
          type="date"
          :label="t('LABEL_FROM')"
          v-model="absenceCreateRequest.AbsenceFrom"
        />
        <q-input
          type="date"
          :label="t('LABEL_TILL')"
          v-model="absenceCreateRequest.AbsenceTill"
        />
        <q-select
          :label="t('LABEL_REASON')"
          emit-value
          :options="absenceReasons"
          map-options
          option-value="ID"
          option-label="Description"
          v-model="absenceCreateRequest.AbsenceReasonID"
        />
      </q-card-section>
      <q-card-section>
        <AbsenceSummaryTableComponent
          v-if="conflictingAbsences.length > 0"
          v-model="conflictingAbsences"
          flat
          :title="t('LABEL_ABSENCE_CONFLICTING')"
        />
      </q-card-section>
      <q-card-actions>
        <q-btn v-close-popup :label="t('BTN_CANCEL')" color="negative"/>
        <q-btn
          v-close-popup
          :label="t('BTN_CREATE')"
          color="positive"
          @click="createAbsence"
        />
      </q-card-actions>
    </q-card>
  </q-dialog>
</template>

<style scoped></style>
