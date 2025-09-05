<script setup lang="ts">
import { onMounted, ref } from 'vue';
import BeeTimeClock from 'src/service/BeeTimeClock';
import { ExternalWorkCompensation } from 'src/models/ExternalWork';
import { useI18n } from 'vue-i18n';
import { date } from 'quasar';
import { emptyPagination } from 'src/helper/objects';
import ExternalWorkCompensationAddtionalOptionsList
  from 'components/external_work/ExternalWorkCompensationAddtionalOptionsList.vue';
import type { ErrorResponse } from 'src/models/Base';
import { showErrorMessage } from 'src/helper/message';

const { t } = useI18n();
const externalWorkCompensations = ref<ExternalWorkCompensation[]>([]);
const isLoading = ref(true);
const selectedWorkCompensation = ref<ExternalWorkCompensation>();

const columns = [
  {
    name: 'isoCountryCodeA2',
    label: t('LABEL_COUNTRY'),
    field: 'IsoCountryCodeA2',
  },
  {
    name: 'validFrom',
    label: t('LABEL_VALID_FROM'),
    field: 'ValidFrom',
    format: (val: Date) => date.formatDate(val, 'DD. MMM. YYYY'),
  },
  {
    name: 'validTill',
    label: t('LABEL_VALID_TILL'),
    field: 'ValidTill',
    format: (val: Date) => date.formatDate(val, 'DD. MMM. YYYY'),
  },
];

function loadExternalWorkCompensation() {
  isLoading.value = true;
  BeeTimeClock.administrationExternalWorkCompensation()
    .then((result) => {
      if (result.status === 200) {
        externalWorkCompensations.value = result.data.Data.map((s) =>
          ExternalWorkCompensation.fromApi(s)
        );
      }
    }).catch((error: ErrorResponse) => {
    showErrorMessage(error.message);
  })
    .finally(() => {
      isLoading.value = false;
    });
}

function selectWorkCompensation(workCompensation: ExternalWorkCompensation) {
  selectedWorkCompensation.value = workCompensation;
}

function save() {
  if (!selectedWorkCompensation.value) return;
  BeeTimeClock.administrationExternalWorkCompensationUpdate(selectedWorkCompensation.value.ID, selectedWorkCompensation.value).then(result => {
    if (result.status === 200) {
      loadExternalWorkCompensation()
    }
  }).catch((error: ErrorResponse) => {
    showErrorMessage(error.message);
  })
}

onMounted(() => {
  loadExternalWorkCompensation();
});
</script>

<template>
  <q-page padding>
    <div v-if="!isLoading">
      <q-table
        :rows="externalWorkCompensations"
        :columns="columns"
        :pagination="emptyPagination"
        hide-pagination
      >
        <template v-slot:header="props">
          <q-tr :props="props">
            <q-th v-for="col in props.cols" :key="col.name" :props="props">
              {{ col.label }}
            </q-th>
            <q-th auto-width />
          </q-tr>
        </template>
        <template v-slot:body="props">
          <q-tr :props="props">
            <q-td v-for="col in props.cols" :key="col.name" :props="props">
              {{ col.value }}
            </q-td>
            <q-td auto-width>
              <q-btn
                icon="edit"
                color="primary"
                @click="selectWorkCompensation(props.row)"
              />
            </q-td>
          </q-tr>
        </template>
      </q-table>
    </div>
    <q-list v-if="selectedWorkCompensation">
      <q-item>
        <q-item-section>
          <q-input
            readonly
            v-model="selectedWorkCompensation.IsoCountryCodeA2"
          />
        </q-item-section>
      </q-item>
      <q-item>
        <q-item-section>
          <q-input
            :label="t('LABEL_PRIVATE_CAR_COMPENSATION_EURO')"
            type="number"
            v-model.number="selectedWorkCompensation.PrivateCarKmCompensation"
          />
        </q-item-section>
      </q-item>
      <q-expansion-item :label="t('LABEL_WITHOUT_SOCIAL_INSURANCE')" :content-inset-level="0.5">
        <q-list>
          <q-item v-for="(item, index) in selectedWorkCompensation.WithoutSocialInsuranceSlots" :key="index">
            <q-item-section>
              <q-input v-model.number="item.Hours" :label="t('LABEL_HOUR', 2)"/>
            </q-item-section>
            <q-item-section>
              <q-input v-model.number="item.Compensation" :label="t('LABEL_COMPENSATION_IN_EURO')"/>
            </q-item-section>
          </q-item>
          <q-item>
            <q-item-section>
              <q-btn class="full-width" color="positive" icon="add" :label="t('LABEL_ADD')"/>
            </q-item-section>
          </q-item>
        </q-list>
      </q-expansion-item>
      <q-expansion-item :label="t('LABEL_WITH_SOCIAL_INSURANCE')"  :content-inset-level="0.5">
        <q-list>
          <q-item v-for="(item, index) in selectedWorkCompensation.WithSocialInsuranceSlots" :key="index">
            <q-item-section>
              <q-input v-model.number="item.Hours" :label="t('LABEL_HOUR', 2)"/>
            </q-item-section>
            <q-item-section>
              <q-input v-model.number="item.Compensation" :label="t('LABEL_COMPENSATION_IN_EURO')"/>
            </q-item-section>
          </q-item>
          <q-item>
            <q-item-section>
              <q-btn class="full-width" color="positive" icon="add" :label="t('LABEL_ADD')"/>
            </q-item-section>
          </q-item>
        </q-list>
      </q-expansion-item>
      <q-expansion-item :label="t('LABEL_OPTION', 2)" :content-inset-level="0.5">
       <ExternalWorkCompensationAddtionalOptionsList v-model="selectedWorkCompensation"/>
      </q-expansion-item>
      <q-item>
        <q-item-section>
          <q-btn class="full-width" icon="save" :label="t('LABEL_SAVE')" color="positive" @click="save"/>
        </q-item-section>
      </q-item>
    </q-list>
    <q-inner-loading :showing="isLoading" />
  </q-page>
</template>

<style scoped></style>
