<script setup lang="ts">
import { onMounted, ref } from 'vue';
import { AbsenceReason } from 'src/models/Absence';
import { useI18n } from 'vue-i18n';
import BeeTimeClock from 'src/service/BeeTimeClock';

const { t } = useI18n();
const absenceReasons = ref<AbsenceReason[]>([]);
const showDialog = ref(false);
const columns = [
  {
    name: 'ID',
    required: true,
    label: t('LABEL_ID'),
    align: 'left',
    field: 'ID',
    sortable: true,
  },
  {
    name: 'Description',
    required: true,
    label: t('LABEL_DESCRIPTION'),
    align: 'left',
    field: 'Description',
    sortable: true,
  },
  {
    name: 'Action',
    align: 'right',
  },
];

function loadAbsenceReasons() {
  BeeTimeClock.absenceReasons().then(result => {
    if (result.status === 200) {
      absenceReasons.value = result.data.Data.map(s => AbsenceReason.fromApi(s))
    }
  })
}

onMounted(() => {
  loadAbsenceReasons()
})
</script>

<template>
  <q-table :columns="columns" :rows="absenceReasons" />
  <q-dialog v-model="showDialog">
    <q-card>
      <q-card-section></q-card-section>
      <q-form>
        <q-card-section></q-card-section>
        <q-card-section>
          <q-card-actions>
            <q-btn color="negative"/>
            <q-btn color="positive"/>
          </q-card-actions>
        </q-card-section>
      </q-form>
    </q-card>
  </q-dialog>
</template>

<style scoped></style>
