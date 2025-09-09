<script setup lang="ts">
import TimeInput from 'components/TimeInput.vue';
import { ExternalWorkExpanse } from 'src/models/ExternalWork';
import { useI18n } from 'vue-i18n';

const {t} = useI18n()

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
  <q-card class="q-mb-lg" v-if="externalWorkExpanse">
    <q-card-section class="bg-secondary"
      >{{ t('LABEL_TRAVEL') }}
    </q-card-section>
    <q-card-section>
      <div class="row">
        <div class="col q-px-md">
          <TimeInput
            class="col"
            v-model="externalWorkExpanse.DepartureTimeDate"
            v-model:date="externalWorkExpanse.Date"
            :readonly="!editMode"
            :label="t('LABEL_DEPARTURE_TIME_HOME')"
          />
        </div>
        <div class="col q-px-md">
          <TimeInput
            class="col"
            v-model="externalWorkExpanse.ArrivalTimeDate"
            v-model:date="externalWorkExpanse.Date"
            :readonly="!editMode"
            :label="t('LABEL_ARRIVAL_TIME_HOME')"
          />
        </div>
      </div>
      <div class="row">
        <div class="col q-px-md">
          <q-input
            v-model.number="externalWorkExpanse.TravelDurationHours"
            :label="t('LABEL_TRAVEL_DURATION_HOURS')"
            :readonly="!editMode"
          />
        </div>
        <div class="col q-px-md">
          <q-input
            v-model.number="externalWorkExpanse.RestDurationHours"
            :label="t('LABEL_REST_DURATION_HOURS')"
            :readonly="!editMode"
          />
        </div>
        <div class="col q-px-md">
          <q-input
            v-model.number="externalWorkExpanse.TravelWithPrivateCarKm"
            :label="t('LABEL_TRAVEL_WITH_PRIVATE_CAR_KM')"
            :readonly="!editMode"
          />
        </div>
      </div>
    </q-card-section>
  </q-card>
</template>

<style scoped></style>
