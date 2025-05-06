<script setup lang="ts">
import { onMounted, ref } from 'vue';
import { emptyPagination } from 'src/helper/objects';
import {
  ApiExternalWorkCreateRequest,
  ExternalWork,
} from 'src/models/ExternalWork';
import { date } from 'quasar';
import formatDate = date.formatDate;
import { useI18n } from 'vue-i18n';
import DateTimePickerComponent from 'components/DateTimePickerComponent.vue';
import BeeTimeClock from 'src/service/BeeTimeClock';
import { showInfoMessage } from 'src/helper/message';

const { t } = useI18n();
const externalWorkItems = ref<ExternalWork[]>();
const promptCreateExternalWork = ref(false);
const externalWorkCreateRequest = ref<ApiExternalWorkCreateRequest>();

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
  },
  {
    name: 'Till',
    required: true,
    label: t('LABEL_TILL'),
    align: 'left',
    field: 'Till',
    format: (val: Date) => `${formatDate(val, 'ddd. DD.MM.YYYY')}`,
  },
];

function loadExternalWorkItems() {
  BeeTimeClock.getExternalWork().then((result) => {
    if (result.status === 200) {
      externalWorkItems.value = result.data.Data.map((s) =>
        ExternalWork.fromApi(s)
      );
    }
  });
}

function openCreateExternalWorkDialog() {
  externalWorkCreateRequest.value = {} as ApiExternalWorkCreateRequest;
  promptCreateExternalWork.value = true;
}

function saveExternalWork() {
  if (!externalWorkCreateRequest.value) return;
  BeeTimeClock.createExternalWork(externalWorkCreateRequest.value).then(
    (result) => {
      if (result.status === 201) {
        showInfoMessage(
          t('MSG_CREATE_SUCCESS', { item: t('LABEL_EXTERNAL_WORK') })
        );
        promptCreateExternalWork.value = false;
        loadExternalWorkItems();
      }
    }
  );
}

onMounted(() => {
  loadExternalWorkItems();
});
</script>

<template>
  <q-page padding>
    <q-table
      :title="$t('LABEL_EXTERNAL_WORK', 2)"
      :rows="externalWorkItems"
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
          @click="openCreateExternalWorkDialog()"
        />
      </template>
      <template v-slot:header="props">
        <q-tr :props="props">
          <q-th
            v-for="col in props.cols"
            :key="col.name"
            :props="props"
          >
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
            <q-btn color="primary" icon="chevron_right" :to="{name: 'ExternalWorkDetail', params: {externalWorkId: props.row.ID}}"/>
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
            <q-input
              v-model="externalWorkCreateRequest.Description"
              class="q-mb-md"
              :label="$t('LABEL_DESCRIPTION')"
            />
            <DateTimePickerComponent
              v-model="externalWorkCreateRequest.From"
              class="q-mb-md"
              :label="$t('LABEL_FROM')"
            />
            <DateTimePickerComponent
              v-model="externalWorkCreateRequest.Till"
              :label="$t('LABEL_TILL')"
            />
          </q-card-section>
          <q-card-section>
            <q-card-actions>
              <q-btn
                :label="$t('BTN_SAVE')"
                icon="save"
                color="positive"
                type="submit"
              />
            </q-card-actions>
          </q-card-section>
        </q-form>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<style scoped></style>
