<script setup lang="ts">
import { date } from 'quasar';
import { ExternalWorkExpanse } from 'src/models/ExternalWork';
import { ref } from 'vue';
import TimeInput from 'components/TimeInput.vue';

const externalWorkExpanse = defineModel('externalworkexpanse', {
  type: ExternalWorkExpanse,
  required: true,
});
const editMode = ref(false);
</script>

<template>
  <q-card class="q-mb-lg">
    <q-card-section class="bg-primary text-white text-h6 row">
      <div class="col">
        {{ date.formatDate(externalWorkExpanse.Date, 'ddd. DD.MM.YYYY') }}
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
            @click="editMode = false"
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
                :label="$t('LABEL_DEPARTURE_TIME')"
              />
            </q-item-label>
          </q-item-section>
          <q-item-section>
            <q-item-label>
              <TimeInput
                v-model="externalWorkExpanse.ArrivalTime"
                v-model:date="externalWorkExpanse.Date"
                :readonly="!editMode"
                :label="$t('LABEL_ARRIVAL_TIME')"
              />
            </q-item-label>
          </q-item-section>
          <q-item-section></q-item-section>
          <q-item-section>
            <q-input
              v-model.number="externalWorkExpanse.TravelDurationHours"
              type="number"
              :label="$t('LABEL_TRAVEL_DURATION')"
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
              v-model="externalWorkExpanse.PauseDurationHours"
              :label="$t('LABEL_ON_SITE_DURATION')"
              readonly
            />
          </q-item-section>
        </q-item>
        <q-item>
          <q-item-section>
            <q-input v-model="externalWorkExpanse.Place" :label="$t('LABEL_PLACE')" :readonly="!editMode"/>
          </q-item-section>
        </q-item>
      </q-list>
    </q-card-section>
  </q-card>
</template>

<style scoped></style>
