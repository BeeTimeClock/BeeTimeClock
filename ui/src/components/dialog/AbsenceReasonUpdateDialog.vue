<script setup lang="ts">
import type { AbsenceReason } from 'src/models/Absence';
import BeeTimeClock from 'src/service/BeeTimeClock';
import { showErrorMessage, showInfoMessage } from 'src/helper/message';
import type { ErrorResponse } from 'src/models/Base';
import { computed } from 'vue';
import { useI18n } from 'vue-i18n';

const {t} = useI18n();

const show = defineModel('show', { default: false });
const value = defineModel<AbsenceReason>({ required: true });
const overtimeImpactOptions = [
  {
    label: t('LABEL_OVERTIME_IMPACT_NONE'),
    value: 'none'
  },
  {
    label: t('LABEL_OVERTIME_IMPACT_DURATION'),
    value: 'duration'
  },
  {
    label: t('LABEL_OVERTIME_IMPACT_HOURS'),
    value: 'hours'
  },
  {
    label: t('LABEL_OVERTIME_IMPACT_DAYS'),
    value: 'days'
  },
]

const isNewReason = computed(() => {
  if (!value.value) return true;
  return !value.value.ID;
});

function saveAbsenceReason() {
  if (!value.value) return;

  if (isNewReason.value) {
    BeeTimeClock.administrationCreateAbsenceReason(value.value)
      .then((result) => {
        if (result.status === 201) {
          showInfoMessage(t('MSG_CREATE_SUCCESS', { item: t('LABEL_REASON') }));
          show.value = false;
        }
      })
      .catch((error: ErrorResponse) => {
        showErrorMessage(error.message);
      });
  } else {
    BeeTimeClock.administrationUpdateAbsenceReason(
      value.value.ID,
      value.value,
    )
      .then((result) => {
        if (result.status === 200) {
          showInfoMessage(t('MSG_UPDATE_SUCCESS'));
          show.value = false;
        }
      })
      .catch((error: ErrorResponse) => {
        showErrorMessage(error.message);
      });
  }
}
</script>

<template>
  <q-dialog v-model="show">
    <q-card>
      <q-card-section class="bg-primary text-h6 text-white">
        <div v-if="isNewReason">
          {{ t('TITLE_CREATE', { item: t('LABEL_REASON') }) }}
        </div>
        <div v-else>
          {{ t('TITLE_UPDATE', { item: t('LABEL_REASON') }) }}
        </div>
      </q-card-section>
      <q-form @submit="saveAbsenceReason">
        <q-card-section>
          <q-input
            v-model="value.Description"
            :label="t('LABEL_DESCRIPTION')"
          />
          <q-toggle
            v-model="value.NeedsApproval"
            :label="t('LABEL_NEEDS_APPROVAL')"
          />
          <q-select v-model="value.Impact" :options="overtimeImpactOptions" map-options emit-value :label="t('LABEL_OVERTIME_IMPACT')"/>
          <q-input type="number" v-model.number="value.ImpactHours" :label="t('LABEL_HOUR', 2)"/>
          <q-input type="number" v-model.number="value.ImpactDays" :label="t('LABEL_DAY', 2)"/>
        </q-card-section>
        <q-card-section>
          <q-card-actions>
            <q-btn color="negative" :label="t('BTN_CANCEL')" v-close-popup />
            <q-btn
              color="positive"
              :label="isNewReason ? t('BTN_CREATE') : t('BTN_SAVE')"
              type="submit"
            />
          </q-card-actions>
        </q-card-section>
      </q-form>
    </q-card>
  </q-dialog>
</template>

<style scoped></style>
