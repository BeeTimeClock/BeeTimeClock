/** eslint-disable @typescript-eslint/consistent-type-imports */
<script setup lang="ts">
import { computed, onMounted, ref } from 'vue';
import type {
  AbsenceSummaryItem,
  AbsenceUserSummary,
  AbsenceUserSummaryYear} from 'src/models/Absence';
import  {
  AbsenceSignedStatus
} from 'src/models/Absence';
import { AbsenceReason } from 'src/models/Absence';
import { Absence } from 'src/models/Absence';
import BeeTimeClock from 'src/service/BeeTimeClock';
import type { QTableColumn } from 'quasar';
import { date, useQuasar } from 'quasar';
import AbsenceSummaryTableComponent from 'components/AbsenceSummaryTableComponent.vue';
import AbsenceCreateDialog from 'src/components/absence/AbsenceCreateDialog.vue';
import { useI18n } from 'vue-i18n';
import { showErrorMessage, showInfoMessage } from 'src/helper/message';
import type { AxiosError } from 'axios';
import type { BaseResponse } from 'src/models/Base';
import { type ErrorResponse } from 'src/models/Base';

const { t } = useI18n();
const q = useQuasar();

const absences = ref([] as Absence[]);
const absenceSummaryItems = ref([] as AbsenceSummaryItem[]);
const mySummary = ref(null as AbsenceUserSummary | null);
const absenceReasons = ref<AbsenceReason[]>([]);
const promptAbsenceCreationDialog = ref(false);

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
    name: 'absenceSigningStatus',
    label: t('LABEL_SIGN'),
    field: 'SignedStatus'
  },
  {
    name: 'absenceActions',
    label: t('LABEL_ACTION', 2),
  },
] as QTableColumn[];

const pagination = {
  rowsPerPage: 10,
};

const myAbsences = computed(() => {
  if (!absences.value) return [];

  const data = absences.value;
  return data.sort(
    (a, b) =>
      new Date(a.AbsenceFrom).getTime() - new Date(b.AbsenceFrom).getTime(),
  );
});

function loadAbsenceReasons() {
  BeeTimeClock.absenceReasons()
    .then((result) => {
      if (result.status === 200) {
        absenceReasons.value = result.data.Data.map((s) =>
          AbsenceReason.fromApi(s),
        );
      }
    })
    .catch((error: ErrorResponse) => {
      showErrorMessage(error.message);
    });
}

function openAbsenceCreationDialog() {
  promptAbsenceCreationDialog.value = true;
}

function loadAbsenceSummary() {
  BeeTimeClock.queryAbsenceSummary()
    .then((result) => {
      if (result.status === 200) {
        absenceSummaryItems.value = result.data.Data;
      }
    })
    .catch((error: ErrorResponse) => {
      showErrorMessage(error.message);
    });
}

function loadAbsences() {
  BeeTimeClock.getAbsences()
    .then((result) => {
      if (result.status === 200) {
        absences.value = result.data.Data.map((s) => Absence.fromApi(s));
      }
    })
    .catch((error: ErrorResponse) => {
      showErrorMessage(error.message);
    });
}

function loadMySummary() {
  BeeTimeClock.queryMyAbsenceSummary()
    .then((result) => {
      if (result.status === 200) {
        mySummary.value = result.data.Data;
      }
    })
    .catch((error: ErrorResponse) => {
      showErrorMessage(error.message);
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
  const res = absenceReasons.value.filter((s: AbsenceReason) => s.ID == id);
  if (res.length == 0) return '';

  return res[0]!.Description;
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
  loadAbsenceReasons();
  loadAbsences();
  loadAbsenceSummary();
  loadMySummary();
}

onMounted(() => {
  loadAbsenceReasons();
  refresh();
});
</script>

<template>
  <q-page padding>
    <div class="q-mb-lg" v-if="getCurrentYearSummary()">
      <div class="text-h4">Jahresübersicht (verbraucht)</div>
      <div class="row q-mt-sm">
        <div
          class="col bg-primary q-mr-xl text-white text-center rounded-borders"
          v-for="(absenceSummary, absenceReason) in getCurrentYearSummary()!
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
      :title="t('LABEL_MY_ABSENCES')"
      :rows="myAbsences"
      :columns="myAbsencesColumns"
      :pagination="pagination"
    >
      <template v-slot:top>
        <div class="col-2 q-table__title">{{ t('LABEL_MY_ABSENCES') }}</div>
        <q-space />
        <q-btn
          color="positive"
          icon="add"
          @click="openAbsenceCreationDialog()"
        />
      </template>
      <template v-slot:body="props">
        <q-tr :props="props" :key="`m_${props.row.index}`">
          <q-td v-for="col in props.cols" :key="col.name" :props="props">
            <div v-if="col.name == 'absenceSigningStatus'">
              <q-icon v-if="props.row.SignedStatus == AbsenceSignedStatus.Accepted" name="check" color="positive" size="sm"/>
              <q-icon v-else-if="props.row.SignedStatus == AbsenceSignedStatus.Declined" name="cancel" color="negative" size="sm"/>
            </div>
            <div v-else-if="col.name == 'absenceActions'">
              <q-btn
                v-if="props.row.Deletable"
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
    <AbsenceCreateDialog
      v-model:show="promptAbsenceCreationDialog"
      @create="refresh()"
    />
  </q-page>
</template>

<style scoped></style>
