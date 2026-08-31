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
  AbsenceSummaryItem,
} from 'src/models/Absence';
import AbsenceSummaryTableComponent from 'components/AbsenceSummaryTableComponent.vue';
import { type User } from 'src/models/Authentication';
import { date, type QTableColumn, useQuasar } from 'quasar';
import { emptyPagination } from 'src/helper/objects';
import { useAuthStore } from 'stores/microsoft-auth';
import AbsenceCreateDialog from 'components/absence/AbsenceCreateDialog.vue';
import UserAvatar from 'components/UserAvatar.vue';

const { t } = useI18n();
const $q = useQuasar();
const auth = useAuthStore();
const selectedTeam = ref<Team>();
const teams = ref<Team[]>([]);
const teamAbsenceSummaries = ref<AbsenceSummaryItem[]>([]);
const neededApprovals = ref<Absence[]>([]);
const isLoading = ref(true);

const approvalColumns = [
  {
    name: 'user',
    field: (row: Absence) => row.userMapped,
    label: t('LABEL_USER'),
    align: 'left',
    format: (val: User) => val?.displayName ?? '-',
  },
  {
    name: 'absenceFrom',
    label: t('LABEL_FROM'),
    field: 'AbsenceFrom',
    align: 'left',
    format: (val: Date) => date.formatDate(val, 'DD. MMM. YYYY'),
  },
  {
    name: 'absenceTill',
    label: t('LABEL_TILL'),
    field: 'AbsenceTill',
    align: 'left',
    format: (val: Date) => date.formatDate(val, 'DD. MMM. YYYY'),
  },
  {
    name: 'absenceNettoDays',
    label: t('LABEL_NETTO_DAYS'),
    field: 'NettoDays',
    align: 'left',
  },
  {
    name: 'actions',
    label: '',
    align: 'right',
  },
] as QTableColumn[];

const isLead = computed(() =>
  selectedTeam.value?.Members.find(
    (s) =>
      s.UserID === auth.getSession()?.ID &&
      (s.Level === TeamLevel.LeadSurrogate || s.Level === TeamLevel.Lead),
  ) != null,
);

const selectedUser = ref<User>();
const showAbsenceCreateDialog = ref(false);

function levelColor(level: TeamLevel) {
  if (level === TeamLevel.Lead) return 'primary';
  if (level === TeamLevel.LeadSurrogate) return 'secondary';
  return 'grey-5';
}

function levelLabel(level: TeamLevel) {
  if (level === TeamLevel.Lead) return t('LABEL_TEAM_LEAD');
  if (level === TeamLevel.LeadSurrogate) return t('LABEL_TEAM_LEAD_SURROGATE');
  return t('LABEL_TEAM_MEMBER', 1);
}

function loadMyTeams() {
  BeeTimeClock.getTeams()
    .then((result) => {
      if (result.status === 200) {
        teams.value = result.data.Data.map((s) => Team.fromApi(s));
        if (teams.value.length >= 1) selectedTeam.value = teams.value[0]!;
      }
    })
    .catch((error) => showErrorMessage(error));
}

function loadTeamAbensces() {
  if (!selectedTeam.value) return;
  isLoading.value = true;
  BeeTimeClock.queryTeamAbsenceSummary(selectedTeam.value.ID)
    .then((result) => {
      if (result.status === 200)
        teamAbsenceSummaries.value = result.data.Data.map((s) => AbsenceSummaryItem.fromApi(s));
    })
    .catch((error) => showErrorMessage(error))
    .finally(() => { isLoading.value = false; });
}

function loadNeededApprovals() {
  if (!selectedTeam.value || !isLead.value) return;
  BeeTimeClock.absenceTeamOpen(selectedTeam.value.ID)
    .then((result) => {
      if (result.status === 200)
        neededApprovals.value = result.data.Data.map((s) => Absence.fromApi(s));
    })
    .catch((error) => showErrorMessage(error));
}

function createAbsenceForTeamMember(teamMember: TeamMember) {
  selectedUser.value = teamMember.userMapped;
  showAbsenceCreateDialog.value = true;
}

function signAbsence(absence: Absence, status: AbsenceSignedStatus, message?: string) {
  if (!selectedTeam.value) return;
  BeeTimeClock.absenceTeamSign(
    selectedTeam.value.ID,
    absence.ID,
    { Status: status, Messages: message } as AbsenceSignRequest,
  )
    .then((result) => {
      if (result.status === 200) {
        showInfoMessage(t('MSG_CREATE_SUCCESS'));
        loadTeamAbensces();
        loadNeededApprovals();
      }
    })
    .catch((error) => showErrorMessage(error));
}

function declineAbsence(absence: Absence) {
  $q.dialog({
    title: t('TITLE_DECLINE'),
    message: t('MSG_ABSENCE_DECLINE'),
    prompt: { model: '', type: 'text' },
    cancel: true,
    persistent: true,
  }).onOk((data) => signAbsence(absence, AbsenceSignedStatus.Declined, data));
}

function acceptAbsence(absence: Absence) {
  $q.dialog({
    title: t('TITLE_ACCEPT'),
    message: t('MSG_ARE_YOU_SURE'),
    cancel: true,
    persistent: true,
  }).onOk(() => signAbsence(absence, AbsenceSignedStatus.Accepted));
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
      v-if="teams.length > 1"
      v-model="selectedTeam"
      :options="teams"
      :label="t('LABEL_TEAM')"
      option-label="Teamname"
      emit-value
      map-options
      outlined
      dense
      class="q-mb-lg"
    />

    <q-inner-loading :showing="isLoading" />

    <div v-if="selectedTeam && !isLoading" class="q-gutter-md">

      <q-card v-if="isLead && neededApprovals.length > 0" flat bordered>
        <q-card-section class="row items-center q-pb-none">
          <q-icon name="pending_actions" color="warning" size="20px" class="q-mr-sm" />
          <span class="text-subtitle2">{{ t('LABEL_ABSENCE') }} — {{ t('LABEL_IN_REVIEW') }}</span>
          <q-badge color="warning" :label="neededApprovals.length" class="q-ml-sm" />
        </q-card-section>
        <q-table
          :columns="approvalColumns"
          :rows="neededApprovals"
          :pagination="emptyPagination"
          flat
          hide-pagination
          hide-header
        >
          <template v-slot:body="props">
            <q-tr :props="props">
              <q-td v-for="col in props.cols" :key="col.name" :props="props">
                <template v-if="col.name === 'actions'">
                  <q-btn
                    flat
                    round
                    color="positive"
                    icon="check_circle"
                    @click="acceptAbsence(props.row)"
                  >
                    <q-tooltip>{{ t('TITLE_ACCEPT') }}</q-tooltip>
                  </q-btn>
                  <q-btn
                    flat
                    round
                    color="negative"
                    icon="cancel"
                    @click="declineAbsence(props.row)"
                  >
                    <q-tooltip>{{ t('TITLE_DECLINE') }}</q-tooltip>
                  </q-btn>
                </template>
                <template v-else>{{ col.value }}</template>
              </q-td>
            </q-tr>
          </template>
        </q-table>
      </q-card>

      <q-card flat bordered>
        <q-card-section class="row items-center q-pb-none">
          <q-icon name="groups" color="primary" size="20px" class="q-mr-sm" />
          <span class="text-subtitle2">{{ selectedTeam.Teamname }}</span>
          <q-badge outline color="primary" :label="selectedTeam.membersMapped.length" class="q-ml-sm" />
        </q-card-section>
        <q-list separator>
          <q-item v-for="member in selectedTeam.membersMapped" :key="member.ID">
            <q-item-section avatar>
              <UserAvatar :user="member.User" :size="36" />
            </q-item-section>
            <q-item-section>
              <q-item-label>{{ member.userMapped?.displayName ?? '-' }}</q-item-label>
              <q-item-label caption>
                <q-chip
                  dense
                  :color="levelColor(member.Level)"
                  text-color="white"
                  :label="levelLabel(member.Level)"
                  size="sm"
                />
              </q-item-label>
            </q-item-section>
            <q-item-section side v-if="isLead">
              <div class="row q-gutter-xs">
                <q-btn
                  flat
                  round
                  icon="visibility"
                  color="primary"
                  :to="{ name: 'TeamUserDetail', params: { teamId: selectedTeam.ID, userId: member.UserID } }"
                >
                  <q-tooltip>{{ t('MENU_OVERVIEW') }}</q-tooltip>
                </q-btn>
                <q-btn
                  flat
                  round
                  icon="event_busy"
                  color="primary"
                  @click="createAbsenceForTeamMember(member)"
                >
                  <q-tooltip>{{ t('MENU_ABSENCE') }}</q-tooltip>
                </q-btn>
              </div>
            </q-item-section>
          </q-item>
        </q-list>
      </q-card>

      <q-card flat bordered>
        <AbsenceSummaryTableComponent v-model="teamAbsenceSummaries" :show-reason="isLead" />
      </q-card>

    </div>

    <AbsenceCreateDialog
      v-if="selectedUser && selectedTeam"
      v-model:user="selectedUser"
      v-model:team="selectedTeam"
      v-model:show="showAbsenceCreateDialog"
      @create="loadTeamAbensces"
    />
  </q-page>
</template>
