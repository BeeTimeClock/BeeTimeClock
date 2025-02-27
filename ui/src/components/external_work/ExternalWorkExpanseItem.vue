<script setup lang="ts">
import { date } from 'quasar';
import { ExternalWork, ExternalWorkExpanse } from 'src/models/ExternalWork';
import { computed, ref } from 'vue';
import BeeTimeClock from 'src/service/BeeTimeClock';
import ExternalWorkExpanseItemTravel from 'components/external_work/ExternalWorkExpanseItemTravel.vue';
import ExternalWorkExpanseItemOnSite from 'components/external_work/ExternalWorkExpanseItemOnSite.vue';
import ExternalWorkExpanseItemSummaryItem from 'components/external_work/ExternalWorkExpanseItemSummaryItem.vue';
import { useI18n } from 'vue-i18n';
import ExternalWorkExpanseItemOptions from 'components/external_work/ExternalWorkExpanseItemOptions.vue';
import { all } from 'axios';

const { t } = useI18n();

const externalWorkExpanse = defineModel('externalworkexpanse', {
  type: ExternalWorkExpanse,
  required: true,
});
const allowEdit = defineModel('allow_edit', { type: Boolean, default: true });
const editMode = ref(false);
const emits = defineEmits(['updated']);
const isNew = computed(() => {
  return externalWorkExpanse.value.ID == 0;
});

const showTravelInput = ref(false);
const showOnSiteInput = ref(false);

const externalWork = computed(() => {
  return externalWorkExpanse.value.ExternalWork;
});

const hasTravelData = computed(() => {
  return (
    externalWorkExpanse.value.DepartureTime ||
    externalWorkExpanse.value.ArrivalTime
  );
});

const hasOnSiteData = computed(() => {
  return (
    externalWorkExpanse.value.OnSiteFrom || externalWorkExpanse.value.OnSiteTill
  );
});

const hasTravelTime = computed({
  get() {
    if (hasTravelData.value) return true;
    return showTravelInput.value;
  },
  set(newValue: boolean) {
    showTravelInput.value = newValue;
  },
});

const hasOnSiteTime = computed({
  get() {
    if (hasOnSiteData.value) return true;
    return showOnSiteInput.value;
  },
  set(newValue: boolean) {
    showOnSiteInput.value = newValue;
  },
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
      editMode.value = false;
      emits('updated');
    }
  });
}

function createExternalWorkExpanse() {
  BeeTimeClock.createExternalWorkExpanse(
    externalWorkExpanse.value.ExternalWorkID,
    externalWorkExpanse.value
  ).then((result) => {
    if (result.status == 201) {
      editMode.value = false;
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
      </div>
      <div class="col-auto">
        <template v-if="allow_edit">
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
        </template>
      </div>
    </q-card-section>
    <q-card-section horizontal class="row">
      <q-card-section class="col">
        <q-card-section class="row" v-if="editMode">
          <q-toggle
            class="col"
            v-model="hasTravelTime"
            :disable="hasTravelData"
            :label="$t('LABEL_TRAVEL_TIME')"
          />
          <q-toggle
            class="col"
            v-model="hasOnSiteTime"
            :disable="hasOnSiteData"
            :label="$t('LABEL_ON_SITE_TIME')"
          />
        </q-card-section>
        <q-card-section>
          <q-input
            class="q-mb-lg"
            v-model="externalWorkExpanse.Place"
            :label="$t('LABEL_PLACE')"
            :readonly="!editMode"
          />
          <ExternalWorkExpanseItemTravel
            v-if="hasTravelTime"
            v-model="externalWorkExpanse"
            v-model:editmode="editMode"
          />
          <ExternalWorkExpanseItemOnSite
            v-if="hasOnSiteTime"
            v-model="externalWorkExpanse"
            v-model:editmode="editMode"
          />
          <ExternalWorkExpanseItemOptions
            v-model="externalWorkExpanse"
            v-model:editmode="editMode"
          />
        </q-card-section>
      </q-card-section>
      <q-separator vertical />
      <q-card-section class="col-3">
        <q-list>
          <q-item-label header>{{ $t('LABEL_TIME', 2) }}</q-item-label>
          <ExternalWorkExpanseItemSummaryItem
            :caption="$t('LABEL_TOTAL_AWAY_HOURS')"
            :label="`${externalworkexpanse.TotalAwayHours} ${t(
              'LABEL_HOUR',
              externalworkexpanse.TotalAwayHours
            )}`"
          />
          <ExternalWorkExpanseItemSummaryItem
            :caption="$t('LABEL_TOTAL_WORKING_HOURS')"
            :label="`${externalworkexpanse.TotalWorkingHours} ${t(
              'LABEL_HOUR',
              externalworkexpanse.TotalAwayHours
            )}`"
          />
          <ExternalWorkExpanseItemSummaryItem
            :caption="$t('LABEL_TOTAL_OPERATION_HOURS')"
            :label="`${externalworkexpanse.TotalOperationHours} ${t(
              'LABEL_HOUR',
              externalworkexpanse.TotalAwayHours
            )}`"
          />
          <ExternalWorkExpanseItemSummaryItem
            :caption="$t('LABEL_TOTAL_OVERTIME_HOURS')"
            :label="`${externalworkexpanse.TotalOvertimeHours} ${t(
              'LABEL_HOUR',
              externalworkexpanse.TotalAwayHours
            )}`"
          />
          <q-item-label header>{{ $t('LABEL_EXPENSE', 2) }}</q-item-label>
          <ExternalWorkExpanseItemSummaryItem
            :caption="$t('LABEL_EXPENSES_WITHOUT_SOCIAL_INSURANCE')"
            :label="`${externalworkexpanse.ExpensesWithoutSocialInsurance} Euro`"
          />
          <ExternalWorkExpanseItemSummaryItem
            :caption="$t('LABEL_EXPENSES_WITH_SOCIAL_INSURANCE')"
            :label="`${externalworkexpanse.ExpensesWithSocialInsurance} Euro`"
          />
        </q-list>
      </q-card-section>
    </q-card-section>
  </q-card>
</template>

<style scoped></style>
