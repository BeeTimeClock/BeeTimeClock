<script setup lang="ts">
import { date } from 'quasar';
import { ExternalWorkExpanse } from 'src/models/ExternalWork';
import { computed, ref } from 'vue';
import TimeInput from 'components/TimeInput.vue';
import BeeTimeClock from 'src/service/BeeTimeClock';

const externalWorkExpanse = defineModel('externalworkexpanse', {
  type: ExternalWorkExpanse,
  required: true,
});
const editMode = ref(false);
const emits = defineEmits(['updated']);
const isNew = computed(() => {
  return externalWorkExpanse.value.ID == 0;
});

function updateOrCreate() {
  if (isNew.value) {
    createExternalWorkExpanse();
  } else {
    updateExternalWorkExpanse();
  }
}

function updateExternalWorkExpanse() {
  BeeTimeClock.updateExternalWorkExpanse(
    externalWorkExpanse.value.ExternalWorkID,
    externalWorkExpanse.value.ID,
    externalWorkExpanse.value
  ).then((result) => {
    if (result.status == 200) {
      emits('updated');
    }
  });
}

function createExternalWorkExpanse() {
  BeeTimeClock.createExternalWorkExpanse(
    externalWorkExpanse.value.ExternalWorkID,
    externalWorkExpanse.value
  ).then((result) => {
    if (result.status == 200) {
      emits('updated');
    }
  });
}
</script>

<template>
  <q-card class="q-mb-lg" v-if="externalWorkExpanse">
    <q-card-section class="bg-primary text-white text-h6 row">
      <div class="col">
        {{ date.formatDate(externalWorkExpanse.Date, 'ddd. DD.MM.YYYY') }}
        {{ externalWorkExpanse.ID }}
        {{ externalWorkExpanse.ExternalWorkID }}
      </div>
      <div class="col-auto">
        <q-btn
          v-if="!editMode"
          color="secondary"
          icon="edit"
          @click="editMode = true"
        />
        <template v-else>
          <q-btn
            color="positive"
            icon="save"
            @click="updateOrCreate"
            class="q-mr-md"
          />
          <q-btn color="negative" icon="cancel" @click="editMode = false" />
        </template>
      </div>
    </q-card-section>
    <q-card-section>
      <q-list>
        <q-item>
          <q-item-section>
            <q-item-label>
              <TimeInput
                v-model="externalWorkExpanse.DepartureTime"
                v-model:date="externalWorkExpanse.Date"
                :readonly="!editMode"
                :label="$t('LABEL_DEPARTURE_TIME_HOME')"
              />
            </q-item-label>
          </q-item-section>
          <q-item-section>
            <q-item-label>
              <TimeInput
                v-model="externalWorkExpanse.ArrivalTime"
                v-model:date="externalWorkExpanse.Date"
                :readonly="!editMode"
                :label="$t('LABEL_ARRIVAL_TIME_HOME')"
              />
            </q-item-label>
          </q-item-section>
          <q-item-section>
            <q-input
              v-model.number="externalWorkExpanse.TravelDurationHours"
              :label="$t('LABEL_TRAVEL_DURATION_HOURS')"
              :readonly="!editMode"
            />
          </q-item-section>
          <q-item-section>
            <q-input
              v-model.number="externalWorkExpanse.RestDurationHours"
              :label="$t('LABEL_REST_DURATION_HOURS')"
              :readonly="!editMode"
            />
          </q-item-section>
        </q-item>
        <q-item>
          <q-item-section>
            <TimeInput
              v-model="externalWorkExpanse.OnSiteFrom"
              v-model:date="externalWorkExpanse.Date"
              :readonly="!editMode"
              :label="$t('LABEL_ON_SITE_FROM')"
            />
          </q-item-section>
          <q-item-section>
            <TimeInput
              v-model="externalWorkExpanse.OnSiteTill"
              v-model:date="externalWorkExpanse.Date"
              :readonly="!editMode"
              :label="$t('LABEL_ON_SITE_TILL')"
            />
          </q-item-section>
          <q-item-section>
            <q-input
              v-model.number="externalWorkExpanse.PauseDurationHours"
              type="number"
              :label="$t('LABEL_PAUSE_DURATION')"
            />
          </q-item-section>
          <q-item-section>
            <q-input
              v-if="!editMode"
              v-model.number="externalWorkExpanse.TotalWorkingHours"
              :label="$t('LABEL_ON_SITE_DURATION')"
              readonly
            />
          </q-item-section>
        </q-item>
        <q-item>
          <q-item-section>
            <q-input
              v-model="externalWorkExpanse.Place"
              :label="$t('LABEL_PLACE')"
              :readonly="!editMode"
            />
          </q-item-section>
        </q-item>
      </q-list>
    </q-card-section>
    <q-card-section>
      <q-list>
        <q-expansion-item :label="$t('LABEL_CALCULATION')">
          <q-list>
            <q-item>
              <q-item-section>
                <q-item-label caption>Spesen (soz frei)</q-item-label>
                <q-item-label>20 Euro</q-item-label>
              </q-item-section>
            </q-item>
          </q-list>
        </q-expansion-item>
      </q-list>
    </q-card-section>
  </q-card>
</template>

<style scoped></style>
