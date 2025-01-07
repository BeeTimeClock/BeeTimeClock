<script setup lang="ts">
import { computed, onMounted, ref } from 'vue';
import {
  Absence,
  AbsenceCreateRequest,
  AbsenceReason,
  AbsenceSummaryItem,
  AbsenceUserSummary,
  AbsenceUserSummaryYear,
} from 'src/models/Absence';
import BeeTimeClock from 'src/service/BeeTimeClock';
import { date, useQuasar } from 'quasar';
import AbsenceSummaryTableComponent from 'components/AbsenceSummaryTableComponent.vue';
import { useI18n } from 'vue-i18n';
import { showErrorMessage, showInfoMessage } from 'src/helper/message';
import { AxiosError } from 'axios';
import { BaseResponse } from 'src/models/Base';

const { t } = useI18n();
const q = useQuasar();

const promptCreateAbsence = ref(false);

const absenceCreateRequest = ref({} as AbsenceCreateRequest);
const absences = ref([] as Absence[]);
const absenceSummaryItems = ref([] as AbsenceSummaryItem[]);
const mySummary = ref(null as AbsenceUserSummary | null);

const myAbsencesColumns = [
  {
    name: 'absenceFrom',
    label: t('LABEL_FROM'),
    field: 'formatFromFull',
  },
  {
    name: 'absenceTill',
    label: t('LABEL_TILL'),
    field: 'formatTillFull',
  },
  {
    name: 'absenceReason',
    label: t('LABEL_REASON'),
    field: 'AbsenceReasonID',
    format: (val: number) => getAbsenceReasonDescriptionById(val),
  },
  {
    name: 'absenceNettoDays',
    label: t('LABEL_NETTO_DAYS'),
    field: 'NettoDays',
  },
  {
    name: 'absenceCreatedAt',
    label: t('LABEL_CREATED_AT'),
    field: 'CreatedAt',
    format: (val: Date) => date.formatDate(val, 'DD.MM.YYYY HH:mm'),
  },
  {
    name: 'absenceActions',
    label: t('LABEL_ACTIONS'),
  },
];

const pagination = {
  rowsPerPage: 10,
};

let absenceReasons = [] as AbsenceReason[];

const conflictingAbsences = computed(() => {
  return absenceSummaryItems.value.filter(
    (s) =>
      (s.AbsenceFrom > absenceCreateRequest.value.AbsenceFrom &&
        s.AbsenceFrom < absenceCreateRequest.value.AbsenceTill) ||
      (s.AbsenceTill < absenceCreateRequest.value.AbsenceTill &&
        s.AbsenceTill > absenceCreateRequest.value.AbsenceFrom),
  );
});

const myAbsences = computed(() => {
  if (!absences.value) return [];

  const data = absences.value;
  return data.sort(
    (a, b) =>
      new Date(a.AbsenceFrom).getTime() - new Date(b.AbsenceFrom).getTime(),
  );
});

function loadAbsenceReasons() {
  BeeTimeClock.absenceReasons().then((result) => {
    if (result.status === 200) {
      absenceReasons = result.data.Data;
    }
  });
}

function createAbsence() {
  BeeTimeClock.createAbsence(absenceCreateRequest.value).then((result) => {
    if (result.status === 201) {
      refresh();
    }
  });
}

function loadAbsenceSummary() {
  BeeTimeClock.queryAbsenceSummary().then((result) => {
    if (result.status === 200) {
      absenceSummaryItems.value = result.data.Data;
    }
  });
}

function loadAbsences() {
  BeeTimeClock.getAbsences().then((result) => {
    if (result.status === 200) {
      absences.value = result.data.Data.map(s => Absence.fromApi(s));
    }
  });
}

function loadMySummary() {
  BeeTimeClock.queryMyAbsenceSummary().then((result) => {
    if (result.status === 200) {
      mySummary.value = result.data.Data;
    }
  });
}

function getCurrentYearSummary(): AbsenceUserSummaryYear | null {
  const currentYear = new Date().getFullYear();

  if (!mySummary.value) {
    return null;
  }

  if (
    Object.keys(mySummary.value?.ByYear).filter(
      (s) => s == currentYear.toString(),
    ).length == 0
  ) {
    return null;
  }

  return mySummary.value?.ByYear[currentYear] ?? null;
}

function getAbsenceReasonDescriptionById(id: number): string {
  const res = absenceReasons.filter((s) => s.ID == id);
  if (res.length == 0) return '';

  return res[0].Description;
}

function deleteAbsence(absence: Absence) {
  const spanMessage = `${date.formatDate(absence.AbsenceFrom, 'DD.MM.YYYY')} - ${date.formatDate(absence.AbsenceTill, 'DD.MM.YYYY')}`;
  q.dialog({
    title: t('LABEL_DELETE'),
    message: t('MSG_DELETE', {
      item: t('LABEL_ABSENCE'),
      identifier: spanMessage,
    }),
    cancel: true,
    persistent: true,
  }).onOk(() => {
    BeeTimeClock.deleteAbsence(absence.ID)
      .then((result) => {
        if (result.status === 204) {
          showInfoMessage(
            t('MSG_DELETE_SUCCESS', {
              item: t('LABEL_ABSENCE'),
              identifier: spanMessage,
            }),
          );
          refresh();
        }
      })
      .catch((error: AxiosError<BaseResponse<never>>) => {
        showErrorMessage(error.response?.data.Message);
      });
  });
}

function refresh() {
  loadAbsences();
  loadAbsenceSummary();
  loadMySummary();
}

onMounted(async () => {
  await loadAbsenceReasons();
  await refresh();
});
</script>

<template>
  <q-page padding>
    <div class="q-mb-lg" v-if="getCurrentYearSummary()">
      <div class="text-h4">Jahresübersicht (verbraucht)</div>
      <div class="row q-mt-sm">
        <div
          class="col bg-primary q-mr-xl text-white text-center rounded-borders"
          v-for="(absenceSummary, absenceReason) in getCurrentYearSummary()
            .ByAbsenceReason"
          :key="absenceReason"
        >
          <div class="text-h5">
            {{ absenceSummary.Past
            }}<span v-if="absenceSummary.Upcoming > 0" class="text-warning">
              / {{ absenceSummary.Upcoming }}</span
            >
          </div>
          <div class="text-h6">
            {{ getAbsenceReasonDescriptionById(absenceReason) }}
          </div>
        </div>
      </div>
    </div>
    <q-table
      :title="$t('LABEL_MY_ABSENCES')"
      :rows="myAbsences"
      :columns="myAbsencesColumns"
      :pagination="pagination"
    >
      <template v-slot:top>
        <div class="col-2 q-table__title">{{ $t('LABEL_MY_ABSENCES') }}</div>
        <q-space />
        <q-btn
          color="positive"
          icon="add"
          @click="promptCreateAbsence = true"
        />
      </template>
      <template v-slot:body="props">
        <q-tr :props="props" :key="`m_${props.row.index}`">
          <q-td v-for="col in props.cols" :key="col.name" :props="props">
            <div v-if="col.name == 'absenceActions'">
              <q-btn
                icon="delete"
                color="negative"
                @click="deleteAbsence(props.row)"
              />
            </div>
            <div v-else>
              {{ col.value }}
            </div>
          </q-td>
        </q-tr>
      </template>
    </q-table>
    <AbsenceSummaryTableComponent
      class="q-mt-lg"
      v-model="absenceSummaryItems"
    />
  </q-page>
  <q-dialog persistent v-model="promptCreateAbsence">
    <q-card class="q-dialog-plugin full-width">
      <q-card-section>
        <q-input
          type="date"
          :label="$t('LABEL_FROM')"
          v-model="absenceCreateRequest.AbsenceFrom"
        />
        <q-input
          type="date"
          :label="$t('LABEL_TILL')"
          v-model="absenceCreateRequest.AbsenceTill"
        />
        <q-select
          :label="$t('LABEL_REASON')"
          emit-value
          :options="absenceReasons"
          map-options
          option-value="ID"
          option-label="Description"
          v-model="absenceCreateRequest.AbsenceReasonID"
        />
      </q-card-section>
      <q-card-section>
        <AbsenceSummaryTableComponent
          v-if="conflictingAbsences.length > 0"
          v-model="conflictingAbsences"
          flat
          :title="$t('LABEL_ABSENCE_CONFLICTING')"
        />
      </q-card-section>
      <q-card-actions>
        <q-btn v-close-popup :label="$t('BTN_CANCEL')" color="negative" />
        <q-btn
          v-close-popup
          :label="$t('BTN_CREATE')"
          color="positive"
          @click="createAbsence"
        />
      </q-card-actions>
    </q-card>
  </q-dialog>
</template>

<style scoped></style>
