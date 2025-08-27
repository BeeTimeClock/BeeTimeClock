<script setup lang="ts">
import { computed, onMounted, ref } from 'vue';
import { emptyPagination } from 'src/helper/objects';
import type {
  ApiExternalWorkCreateRequest,
  ApiExternalWorkInvoicedInfo,
} from 'src/models/ExternalWork';
import {
  ApiExternalWorkStatus,
  ExternalWork,
  ExternalWorkCompensation,
} from 'src/models/ExternalWork';
import type { QTableColumn } from 'quasar';
import { date } from 'quasar';
import formatDate = date.formatDate;
import { useI18n } from 'vue-i18n';
import BeeTimeClock from 'src/service/BeeTimeClock';
import { showErrorMessage, showInfoMessage } from 'src/helper/message';
import type { ErrorResponse } from 'src/models/Base';

const { t } = useI18n();
const externalWorkItems = ref<ExternalWork[]>();
const externalWorkItemsInvoiced = ref<ApiExternalWorkInvoicedInfo[]>();
const promptCreateExternalWork = ref(false);
const externalWorkCreateRequest = ref<ApiExternalWorkCreateRequest>();
const externalWorkCompensations = ref<ExternalWorkCompensation[]>([]);
const creating = ref(false);

const externalWorkItemsFiltered = computed(() => {
  if (externalWorkItems.value == null) return [];

  return externalWorkItems.value.filter(
    (s) => s.Status != ApiExternalWorkStatus.Invoiced,
  );
});

const columns = [
  {
    name: 'Description',
    required: true,
    label: t('LABEL_DESCRIPTION'),
    align: 'left',
    field: 'Description',
  },
  {
    name: 'From',
    required: true,
    label: t('LABEL_FROM'),
    align: 'left',
    field: 'From',
    format: (val: Date) => `${formatDate(val, 'ddd. DD.MM.YYYY')}`,
    sortable: true,
    sortOrder: 'ad',
  },
  {
    name: 'Till',
    required: true,
    label: t('LABEL_TILL'),
    align: 'left',
    field: 'Till',
    format: (val: Date) => `${formatDate(val, 'ddd. DD.MM.YYYY')}`,
    sortable: true,
  },
  {
    name: 'Status',
    label: t('LABEL_STATUS'),
    align: 'left',
    field: 'Status',
  },
] as QTableColumn[];

const columnsInvoiced = [
  {
    name: 'InvoiceDate',
    required: true,
    label: t('LABEL_DATE'),
    align: 'left',
    field: 'InvoiceDate',
    format: (val: Date) => `${formatDate(val, 'ddd. DD.MM.YYYY')}`,
  },
] as QTableColumn[];

function loadExternalWorkItems() {
  BeeTimeClock.getExternalWork()
    .then((result) => {
      if (result.status === 200) {
        externalWorkItems.value = result.data.Data.map((s) =>
          ExternalWork.fromApi(s),
        );
      }
    })
    .catch((error: ErrorResponse) => {
      showErrorMessage(error.response?.data.Message);
    });
}

function loadExternalWorkItemsInvoiced() {
  BeeTimeClock.getExternalWorkInvoiced()
    .then((result) => {
      if (result.status === 200) {
        externalWorkItemsInvoiced.value = result.data.Data;
      }
    })
    .catch((error: ErrorResponse) => {
      showErrorMessage(error.response?.data.Message);
    });
}

function loadExternalWorkCompensation() {
  BeeTimeClock.externalWorkCompensation()
    .then((result) => {
      if (result.status === 200) {
        externalWorkCompensations.value = result.data.Data.map((s) =>
          ExternalWorkCompensation.fromApi(s),
        );
      }
    })
    .catch((error: ErrorResponse) => {
      showErrorMessage(error.response?.data.Message);
    });
}

function openCreateExternalWorkDialog() {
  externalWorkCreateRequest.value = {} as ApiExternalWorkCreateRequest;
  promptCreateExternalWork.value = true;
}

function saveExternalWork() {
  if (!externalWorkCreateRequest.value) return;
  creating.value = true;
  BeeTimeClock.createExternalWork(externalWorkCreateRequest.value)
    .then((result) => {
      if (result.status === 201) {
        showInfoMessage(
          t('MSG_CREATE_SUCCESS', { item: t('LABEL_EXTERNAL_WORK') }),
        );
        promptCreateExternalWork.value = false;
        loadExternalWorkItems();
      }
    })
    .catch((error: ErrorResponse) => {
      showErrorMessage(error.response?.data.Message);
    })
    .finally(() => {
      creating.value = false;
    });
}

function downloadPdf() {
  BeeTimeClock.externalWorkDownloadPdf()
    .then((result) => {
      if (result.status === 200) {
        console.log(result.data);
        const blob = new Blob([result.data], { type: 'application/pdf' });
        const link = document.createElement('a');
        link.href = URL.createObjectURL(blob);
        link.download = 'Spesen_Test.pdf';
        link.click();
        URL.revokeObjectURL(link.href);

        loadExternalWorkItems();
      }
    })
    .catch((error: ErrorResponse) => {
      showErrorMessage(error.response?.data.Message);
    });
}

function downloadInvoicedPdf(identifier: string) {
  BeeTimeClock.externalWorkDownloadInvoicedPdf(identifier)
    .then((result) => {
      if (result.status === 200) {
        console.log(result.data);
        const blob = new Blob([result.data], { type: 'application/pdf' });
        const link = document.createElement('a');
        link.href = URL.createObjectURL(blob);
        link.download = 'Spesen_Test.pdf';
        link.click();
        URL.revokeObjectURL(link.href);
      }
    })
    .catch((error: ErrorResponse) => {
      showErrorMessage(error.response?.data.Message);
    });
}

onMounted(() => {
  loadExternalWorkCompensation();
  loadExternalWorkItems();
  loadExternalWorkItemsInvoiced();
});
</script>

<template>
  <q-page padding>
    <q-table
      :title="$t('LABEL_EXTERNAL_WORK', 2)"
      :rows="externalWorkItemsFiltered"
      :columns="columns"
      :pagination="emptyPagination"
    >
      <template v-slot:top>
        <div class="col-2 q-table__title">
          {{ $t('LABEL_EXTERNAL_WORK', 2) }}
        </div>
        <q-space />
        <q-btn
          color="positive"
          icon="add"
          class="q-mr-md"
          @click="openCreateExternalWorkDialog()"
        />
        <q-btn
          color="secondary"
          icon="post_add"
          :label="$t('LABEL_GENERATE_REPORTS_FROM_ACCEPTED')"
          @click="downloadPdf"
        />
      </template>
      <template v-slot:header="props">
        <q-tr :props="props">
          <q-th v-for="col in props.cols" :key="col.name" :props="props">
            {{ col.label }}
          </q-th>
          <q-th auto-width />
        </q-tr>
      </template>
      <template v-slot:body="props">
        <q-tr :props="props" :key="`m_${props.row.index}`">
          <q-td v-for="col in props.cols" :key="col.name" :props="props">
            {{ col.value }}
          </q-td>
          <q-td>
            <q-btn
              color="primary"
              icon="chevron_right"
              :to="{
                name: 'ExternalWorkDetail',
                params: { externalWorkId: props.row.ID },
              }"
            />
          </q-td>
        </q-tr>
      </template>
    </q-table>
    <q-table v-if="externalWorkItemsInvoiced"
      :title="$t('LABEL_EXTERNAL_WORK_INVOICED', 2)"
      :rows="externalWorkItemsInvoiced"
      :columns="columnsInvoiced"
      :pagination="emptyPagination"
      class="q-mt-lg"
    >
      <template v-slot:top>
        <div class="col-2 q-table__title">
          {{ $t('LABEL_EXTERNAL_WORK_INVOICED', 2) }}
        </div>
      </template>
      <template v-slot:header="props">
        <q-tr :props="props">
          <q-th v-for="col in props.cols" :key="col.name" :props="props">
            {{ col.label }}
          </q-th>
          <q-th auto-width />
        </q-tr>
      </template>
      <template v-slot:body="props">
        <q-tr :props="props" :key="`m_${props.row.index}`">
          <q-td v-for="col in props.cols" :key="col.name" :props="props">
            {{ col.value }}
          </q-td>
          <q-td>
            <q-btn
              color="primary"
              icon="chevron_right"
              @click="downloadInvoicedPdf(props.row.Identifier)"
            />
          </q-td>
        </q-tr>
      </template>
    </q-table>
    <q-dialog
      v-model="promptCreateExternalWork"
      v-if="externalWorkCreateRequest"
    >
      <q-card>
        <q-card-section class="text-h6 bg-primary text-white">
          {{ $t('LABEL_CREATE', { item: $t('LABEL_EXTERNAL_WORK') }) }}
        </q-card-section>
        <q-form @submit="saveExternalWork">
          <q-card-section>
            <q-select
              v-model="externalWorkCreateRequest.ExternalWorkCompensationID"
              :options="externalWorkCompensations"
              emit-value
              map-options
              option-value="ID"
              option-label="IsoCountryCodeA2"
              :label="$t('LABEL_EXTERNAL_WORK_COMPENSATION')"
            />
            <q-input
              v-model="externalWorkCreateRequest.Description"
              class="q-mb-md"
              :label="$t('LABEL_DESCRIPTION')"
            />
            <q-input
              type="date"
              v-model="externalWorkCreateRequest.From"
              class="q-mb-md"
              :label="$t('LABEL_FROM')"
            />
            <q-input
              type="date"
              v-model="externalWorkCreateRequest.Till"
              :label="$t('LABEL_TILL')"
              :has-time="false"
            />
          </q-card-section>
          <q-card-section>
            <q-card-actions>
              <q-btn
                :label="$t('BTN_SAVE')"
                icon="save"
                color="positive"
                type="submit"
                :loading="creating"
              />
            </q-card-actions>
          </q-card-section>
        </q-form>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<style scoped></style>
