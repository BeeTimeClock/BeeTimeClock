<script setup lang="ts">
import { Team, TeamLevel, type TeamMember } from 'src/models/Team';
import BeeTimeClock from 'src/service/BeeTimeClock';
import { computed, onMounted, ref, watch } from 'vue';
import { useI18n } from 'vue-i18n';
import { showErrorMessage, showInfoMessage } from 'src/helper/message';
import {
  Absence,
  AbsenceSignedStatus,
  type AbsenceSignRequest,
  type AbsenceSummaryItem,
} from 'src/models/Absence';
import AbsenceSummaryTableComponent from 'components/AbsenceSummaryTableComponent.vue';
import { type User } from 'src/models/Authentication';
import { date, type QTableColumn, useQuasar } from 'quasar';
import { emptyPagination } from 'src/helper/objects';
import { useAuthStore } from 'stores/microsoft-auth';
import AbsenceCreateDialog from 'components/absence/AbsenceCreateDialog.vue';

const { t } = useI18n();
const $q = useQuasar();
const auth = useAuthStore();
const selectedTeam = ref<Team>();
const teams = ref<Team[]>([]);
const teamAbsenceSummaries = ref<AbsenceSummaryItem[]>([]);
const neededApprovals = ref<Absence[]>([]);
const isLoading = ref(true);

const columns = [
  {
    name: 'user',
    field: (row: TeamMember) => row.userMapped,
    label: t('LABEL_USER'),
    align: 'left',
    format: (val: User) => (val ? val.displayName : '-'),
  },
  {
    name: 'user',
    field: 'Level',
    label: t('LABEL_USER'),
    align: 'left',
  },
  {
    name: 'actions',
    label: t('LABEL_ACTION', 2),
    align: 'right',
  },
] as QTableColumn[];

const approvalColumns = [
  {
    name: 'user',
    field: (row: TeamMember) => row.userMapped,
    label: t('LABEL_USER'),
    align: 'left',
    format: (val: User) => (val ? val.displayName : '-'),
  },
  {
    name: 'absenceFrom',
    label: t('LABEL_FROM'),
    field: 'AbsenceFrom',
    format: (val: Date) => date.formatDate(val, 'DD. MMM. YYYY'),
  },
  {
    name: 'absenceTill',
    label: t('LABEL_TILL'),
    field: 'AbsenceTill',
    format: (val: Date) => date.formatDate(val, 'DD. MMM. YYYY'),
  },
  {
    name: 'absenceNettoDays',
    label: t('LABEL_NETTO_DAYS'),
    field: 'NettoDays',
  },
  {
    name: 'reason',
    field: 'Reason',
    label: t('LABEL_REASON'),
    align: 'left',
  },
  {
    name: 'actions',
    label: t('LABEL_ACTION', 2),
    align: 'right',
  },
] as QTableColumn[];

const isLead = computed(() => {
  return (
    selectedTeam.value?.Members.find(
      (s) =>
        s.UserID === auth.getSession()?.ID &&
        (s.Level === TeamLevel.LeadSurrogate || s.Level === TeamLevel.Lead),
    ) != null
  );
});

const selectedUser = ref<User>();
const showAbsenceCreateDialog = ref(false);

function loadMyTeams() {
  BeeTimeClock.getTeams()
    .then((result) => {
      if (result.status === 200) {
        teams.value = result.data.Data.map((s) => Team.fromApi(s));
        if (teams.value.length >= 1) {
          selectedTeam.value = teams.value[0]!;
        }
      }
    })
    .catch((error) => {
      showErrorMessage(error);
    });
}

function loadTeamAbensces() {
  if (!selectedTeam.value) return;

  isLoading.value = true;
  BeeTimeClock.queryTeamAbsenceSummary(selectedTeam.value.ID)
    .then((result) => {
      if (result.status === 200) {
        teamAbsenceSummaries.value = result.data.Data;
      }
    })
    .catch((error) => {
      showErrorMessage(error);
    })
    .finally(() => {
      isLoading.value = false;
    });
}

function loadNeededApprovals() {
  if (!selectedTeam.value) return;
  if (!isLead.value) return;
  BeeTimeClock.absenceTeamOpen(selectedTeam.value.ID)
    .then((result) => {
      if (result.status === 200) {
        neededApprovals.value = result.data.Data.map((s) => Absence.fromApi(s));
      }
    })
    .catch((error) => {
      showErrorMessage(error);
    });
}

function createAbsenceForTeamMember(teamMember: TeamMember) {
  selectedUser.value = teamMember.userMapped;
  showAbsenceCreateDialog.value = true;
}

function signAbsence(
  absence: Absence,
  status: AbsenceSignedStatus,
  message?: string,
) {
  if (!selectedTeam.value) return;

  const payload = {
    Status: status,
    Messages: message,
  } as AbsenceSignRequest;

  BeeTimeClock.absenceTeamSign(selectedTeam.value.ID, absence.ID, payload)
    .then((result) => {
      if (result.status === 200) {
        showInfoMessage(t('MSG_CREATE_SUCCESS'));
        loadTeamAbensces();
        loadNeededApprovals();
      }
    })
    .catch((error) => {
      showErrorMessage(error);
    });
}

function declineAbsence(absence: Absence) {
  $q.dialog({
    title: t('TITLE_DECLINE'),
    message: t('MSG_ABSENCE_DECLINE'),
    prompt: {
      model: '',
      type: 'text', // optional
    },
    cancel: true,
    persistent: true,
  }).onOk((data) => {
    signAbsence(absence, AbsenceSignedStatus.Declined, data);
  });
}

function acceptAbsence(absence: Absence) {
  $q.dialog({
    title: t('TITLE_ACCEPT'),
    message: t('MSG_ARE_YOU_SURE'),
    cancel: true,
    persistent: true,
  }).onOk(() => {
    signAbsence(absence, AbsenceSignedStatus.Accepted);
  });
}

watch(selectedTeam, () => {
  loadTeamAbensces();
  loadNeededApprovals();
});

onMounted(() => {
  loadMyTeams();
  loadTeamAbensces();
});
</script>

<template>
  <q-page padding>
    <q-select
      :label="t('LABEL_TEAM')"
      v-model="selectedTeam"
      :options="teams"
      :readonly="teams.length == 1"
      emit-value
      map-options
      option-label="Teamname"
      class="q-mb-lg"
    />

    <div v-if="selectedTeam && !isLoading">
      <q-table
        class="q-mb-md"
        :columns="columns"
        :rows="selectedTeam.membersMapped"
        :pagination="emptyPagination"
        hide-pagination
      >
        <template v-slot:header="props">
          <q-tr :props="props">
            <q-th v-for="col in props.cols" :key="col.name" :props="props">
              {{ col.label }}
            </q-th>
          </q-tr>
        </template>

        <template v-slot:body="props">
          <q-tr :props="props" :key="`m_${props.row.index}`">
            <q-td v-for="col in props.cols" :key="col.name" :props="props">
              <template v-if="col.name == 'actions'">
                <q-btn
                  v-if="isLead"
                  color="primary"
                  icon="add"
                  @click="createAbsenceForTeamMember(props.row)"
                />
              </template>
              <template v-else>
                {{ col.value }}
              </template>
            </q-td>
          </q-tr>
        </template>
      </q-table>

      <q-table
        v-if="isLead && neededApprovals.length > 0"
        :columns="approvalColumns"
        :rows="neededApprovals"
      >
        <template v-slot:header="props">
          <q-tr :props="props">
            <q-th v-for="col in props.cols" :key="col.name" :props="props">
              {{ col.label }}
            </q-th>
          </q-tr>
        </template>

        <template v-slot:body="props">
          <q-tr :props="props" :key="`m_${props.row.index}`">
            <q-td v-for="col in props.cols" :key="col.name" :props="props">
              <template v-if="col.name == 'actions'">
                <q-btn
                  color="positive"
                  icon="check"
                  @click="acceptAbsence(props.row)"
                />
                <q-btn
                  class="q-ml-md"
                  color="negative"
                  icon="cancel"
                  @click="declineAbsence(props.row)"
                />
              </template>
              <template v-else>
                {{ col.value }}
              </template>
            </q-td>
          </q-tr>
        </template>
      </q-table>

      <AbsenceSummaryTableComponent
        class="q-mt-lg"
        v-model="teamAbsenceSummaries"
        :show-reason="isLead"
      />
      <AbsenceCreateDialog
        v-if="selectedUser"
        v-model:user="selectedUser"
        v-model:team="selectedTeam"
        v-model:show="showAbsenceCreateDialog"
        @create="
          () => {
            loadTeamAbensces();
          }
        "
      />
    </div>
    <q-inner-loading :showing="isLoading" />
  </q-page>
</template>

<style scoped></style>
