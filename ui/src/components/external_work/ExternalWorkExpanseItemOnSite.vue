<script setup lang="ts">
import TimeInput from 'components/TimeInput.vue';
import { ExternalWorkExpanse } from 'src/models/ExternalWork';
import { useI18n } from 'vue-i18n';

const {t} = useI18n();

const externalWorkExpanse = defineModel({
  type: ExternalWorkExpanse,
  required: true,
});

const editMode = defineModel('editmode', {
  type: Boolean,
  default: false,
});
</script>

<template>
  <q-card v-if="externalWorkExpanse" class="q-mb-lg">
    <q-card-section class="bg-secondary"
      >{{ t('LABEL_ON_SITE') }}
    </q-card-section>
    <q-card-section>
      <div class="row">
        <div class="col q-px-md">
          <TimeInput
            v-model="externalWorkExpanse.OnSiteFromDate"
            v-model:date="externalWorkExpanse.Date"
            :readonly="!editMode"
            :label="t('LABEL_ON_SITE_FROM')"
          />
        </div>
        <div class="col q-px-md">
          <TimeInput
            v-model="externalWorkExpanse.OnSiteTillDate"
            v-model:date="externalWorkExpanse.Date"
            :readonly="!editMode"
            :label="t('LABEL_ON_SITE_TILL')"
          />
        </div>
      </div>
      <div class="col q-px-md">
        <q-input
          v-model.number="externalWorkExpanse.PauseDurationHours"
          type="number"
          :label="t('LABEL_PAUSE_DURATION')"
        />
      </div>
      <div class="col q-px-md">
        <q-input
          v-if="!editMode"
          v-model.number="externalWorkExpanse.TotalWorkingHours"
          :label="t('LABEL_ON_SITE_DURATION')"
          readonly
        />
      </div>
    </q-card-section>
  </q-card>
</template>

<style scoped></style>
