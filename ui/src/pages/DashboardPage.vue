<script setup lang="ts">
import {computed, onMounted, ref} from 'vue';
import BeeTimeClock from 'src/service/BeeTimeClock';
import type {User} from 'src/models/Authentication';
import {
  Absence,
  AbsenceSignedStatus,
  type AbsenceSignRequest,
  AbsenceSummaryItem,
  type AbsenceReason,
  type ApiAbsenceSummaryItem,
} from 'src/models/Absence';
import { Team, TeamLevel } from 'src/models/Team';
import AbsenceSummaryTableComponent from 'components/AbsenceSummaryTableComponent.vue';
import OvertimeTotal from 'components/OvertimeTotal.vue';
import OvertimeCurrentMonth from 'components/OvertimeMonth.vue';
import CheckInButton from 'components/CheckInButton.vue';
import UserAvatar from 'components/UserAvatar.vue';
import type { ErrorResponse } from 'src/models/Base';
import { showErrorMessage, showInfoMessage } from 'src/helper/message';
import { useI18n } from 'vue-i18n';
import { useAuthStore } from 'stores/microsoft-auth';
import { date, useQuasar } from 'quasar';

const { t, locale } = useI18n({ useScope: 'global' });
const $q = useQuasar();
const auth = useAuthStore();

const user = ref(null as User | null);
const absenceSummaryItems = ref([] as AbsenceSummaryItem[]);

interface PendingApproval {
  teamId: number;
  absence: Absence;
  conflicts: string[];
}
const pendingApprovals = ref<PendingApproval[]>([]);
const isTeamLead = ref(false);
const myPendingAbsences = ref<Absence[]>([]);
const myUpcomingAbsences = ref<Absence[]>([]);
const absenceReasons = ref<AbsenceReason[]>([]);

function absenceReasonLabel(id: number): string {
  return absenceReasons.value.find(r => r.ID === id)?.Description ?? '';
}

function loadMyUpcomingAbsences() {
  const today = new Date();
  today.setHours(0, 0, 0, 0);
  BeeTimeClock.getAbsences().then(result => {
    if (result.status !== 200) return;
    myUpcomingAbsences.value = result.data.Data
      .map(a => Absence.fromApi(a))
      .filter(a => {
        if (new Date(a.AbsenceTill) < today) return false;
        const needsApproval = absenceReasons.value.find(r => r.ID === a.AbsenceReasonID)?.NeedsApproval ?? false;
        return needsApproval ? a.SignedStatus === AbsenceSignedStatus.Accepted : true;
      })
      .sort((a, b) => new Date(a.AbsenceFrom).getTime() - new Date(b.AbsenceFrom).getTime())
      .slice(0, 5);
  }).catch(() => {});
}

function loadMyPendingAbsences() {
  BeeTimeClock.getOpenAbsences().then(result => {
    if (result.status !== 200) return;
    myPendingAbsences.value = result.data.Data
      .filter(a => absenceReasons.value.find(r => r.ID === a.AbsenceReasonID)?.NeedsApproval)
      .map(a => Absence.fromApi(a));
  }).catch(() => {});
}

const formattedDate = computed(() =>
  new Date().toLocaleDateString(locale.value, {
    weekday: 'long',
    year: 'numeric',
    month: 'long',
    day: 'numeric',
  })
);

function loadUser() {
  BeeTimeClock.getMeUser().then(result => {
    if (result.status === 200) user.value = result.data.Data;
  }).catch((error: ErrorResponse) => showErrorMessage(error.message));
}

function loadAbsenceSummary() {
  BeeTimeClock.queryAbsenceSummary().then(result => {
    if (result.status === 200)
      absenceSummaryItems.value = result.data.Data.map(s => AbsenceSummaryItem.fromApi(s));
  }).catch((error: ErrorResponse) => showErrorMessage(error.message));
}

function loadPendingApprovals() {
  const sessionId = auth.getSession()?.ID;
  if (!sessionId) return;

  BeeTimeClock.getTeams().then(result => {
    if (result.status !== 200) return;
    const leadTeams = result.data.Data
      .map(t => Team.fromApi(t))
      .filter(team => team.Members.some(m =>
        m.UserID === sessionId &&
        (m.Level === TeamLevel.Lead || m.Level === TeamLevel.LeadSurrogate)
      ));
    isTeamLead.value = leadTeams.length > 0;

    const promises = leadTeams.map(team =>
      Promise.all([
        BeeTimeClock.absenceTeamOpen(team.ID),
        BeeTimeClock.queryTeamAbsenceSummary(team.ID),
      ]).then(([openR, summaryR]) => {
        const openAbsences = openR.status === 200 ? openR.data.Data.map(a => Absence.fromApi(a)) : [];
        const teamSummary: ApiAbsenceSummaryItem[] = summaryR.status === 200 ? summaryR.data.Data : [];

        return openAbsences.map(absence => {
          const from = new Date(absence.AbsenceFrom).getTime();
          const till = new Date(absence.AbsenceTill).getTime();
          const conflicts = teamSummary
            .filter(s => {
              if (Number(s.User.ID) === Number(absence.User.ID)) return false;
              const sFrom = new Date(s.AbsenceFrom).getTime();
              const sTill = new Date(s.AbsenceTill).getTime();
              return from <= sTill && sFrom <= till;
            })
            .map(s => `${s.User.FirstName} ${s.User.LastName}`);
          return { teamId: team.ID, absence, conflicts };
        });
      }).catch(() => [])
    );

    void Promise.all(promises).then(results => {
      pendingApprovals.value = results.flat();
    });
  }).catch(() => {});
}

function signAbsence(teamId: number, absence: Absence, status: AbsenceSignedStatus, message?: string) {
  BeeTimeClock.absenceTeamSign(teamId, absence.ID, { Status: status, Messages: message } as AbsenceSignRequest)
    .then(result => {
      if (result.status === 200) {
        showInfoMessage(t('MSG_UPDATE_SUCCESS'));
        loadPendingApprovals();
      }
    })
    .catch((error: ErrorResponse) => showErrorMessage(error.response?.data.Message));
}

function acceptAbsence(item: PendingApproval) {
  $q.dialog({
    title: t('TITLE_ACCEPT'),
    message: t('MSG_ARE_YOU_SURE'),
    cancel: true,
    persistent: true,
  }).onOk(() => signAbsence(item.teamId, item.absence, AbsenceSignedStatus.Accepted));
}

function declineAbsence(item: PendingApproval) {
  $q.dialog({
    title: t('TITLE_DECLINE'),
    message: t('MSG_ABSENCE_DECLINE'),
    prompt: { model: '', type: 'text' },
    cancel: true,
    persistent: true,
  }).onOk((msg: string) => signAbsence(item.teamId, item.absence, AbsenceSignedStatus.Declined, msg));
}

onMounted(() => {
  loadUser();
  loadAbsenceSummary();
  loadPendingApprovals();
  BeeTimeClock.absenceReasons()
    .then(r => {
      if (r.status === 200) absenceReasons.value = r.data.Data;
      loadMyPendingAbsences();
      loadMyUpcomingAbsences();
    })
    .catch(() => {});
});
</script>

<template>
  <q-page>
    <div class="dashboard-header q-px-lg q-py-md row items-center">
      <div class="col">
        <div class="text-h6 text-weight-bold text-white" v-if="user">
          {{ t('LABEL_HELLO', { name: user.FirstName }) }}
        </div>
        <div class="text-caption text-white q-mt-xs" style="opacity: 0.75">
          {{ formattedDate }}
        </div>
      </div>
      <CheckInButton>
        <template #elapsed="{ elapsed }">
          <span class="text-white text-body1 text-weight-bold">{{ elapsed }}</span>
        </template>
      </CheckInButton>
    </div>

    <div class="q-pa-lg">
      <div class="row q-col-gutter-md q-mb-lg">
        <div class="col-12 col-sm-6">
          <OvertimeCurrentMonth class="full-height" />
        </div>
        <div class="col-12 col-sm-6">
          <OvertimeTotal class="full-height" />
        </div>
      </div>

      <q-card v-if="myUpcomingAbsences.length > 0" flat bordered class="q-mb-lg">
        <q-card-section class="row items-center q-pb-none">
          <q-icon name="event" color="primary" size="20px" class="q-mr-sm" />
          <span class="text-subtitle2">{{ t('LABEL_MY_UPCOMING_ABSENCES') }}</span>
        </q-card-section>
        <q-list separator>
          <q-item v-for="absence in myUpcomingAbsences" :key="absence.ID">
            <q-item-section>
              <q-item-label>{{ absenceReasonLabel(absence.AbsenceReasonID) }}</q-item-label>
              <q-item-label caption>
                {{ date.formatDate(absence.AbsenceFrom, 'DD. MMM') }} –
                {{ date.formatDate(absence.AbsenceTill, 'DD. MMM YYYY') }}
                ({{ absence.NettoDays }} {{ t('LABEL_DAY') }})
              </q-item-label>
            </q-item-section>
          </q-item>
        </q-list>
      </q-card>

      <q-card v-if="!isTeamLead && myPendingAbsences.length > 0" flat bordered class="q-mb-lg">
        <q-card-section class="row items-center q-pb-none">
          <q-icon name="hourglass_empty" color="warning" size="20px" class="q-mr-sm" />
          <span class="text-subtitle2">{{ t('LABEL_MY_PENDING_ABSENCES') }}</span>
          <q-badge color="warning" :label="myPendingAbsences.length" class="q-ml-sm" />
        </q-card-section>
        <q-list separator>
          <q-item v-for="absence in myPendingAbsences" :key="absence.ID">
            <q-item-section>
              <q-item-label>{{ absenceReasonLabel(absence.AbsenceReasonID) }}</q-item-label>
              <q-item-label caption>
                {{ date.formatDate(absence.AbsenceFrom, 'DD. MMM') }} –
                {{ date.formatDate(absence.AbsenceTill, 'DD. MMM YYYY') }}
                ({{ absence.NettoDays }} {{ t('LABEL_DAY') }})
              </q-item-label>
            </q-item-section>
            <q-item-section side>
              <q-badge color="warning" :label="t('LABEL_PENDING')" />
            </q-item-section>
          </q-item>
        </q-list>
      </q-card>

      <q-card v-if="pendingApprovals.length > 0" flat bordered class="q-mb-lg">
        <q-card-section class="row items-center q-pb-none">
          <q-icon name="pending_actions" color="warning" size="20px" class="q-mr-sm" />
          <span class="text-subtitle2">{{ t('LABEL_IN_REVIEW') }}</span>
          <q-badge color="warning" :label="pendingApprovals.length" class="q-ml-sm" />
        </q-card-section>
        <q-list separator>
          <q-item v-for="item in pendingApprovals" :key="item.absence.ID">
            <q-item-section avatar>
              <UserAvatar :user="item.absence.User" :size="36" />
            </q-item-section>
            <q-item-section>
              <q-item-label>{{ item.absence.userMapped.displayName }}</q-item-label>
              <q-item-label caption>
                {{ absenceReasonLabel(item.absence.AbsenceReasonID) }} ·
                {{ date.formatDate(item.absence.AbsenceFrom, 'DD. MMM') }} –
                {{ date.formatDate(item.absence.AbsenceTill, 'DD. MMM YYYY') }}
                ({{ item.absence.NettoDays }} {{ t('LABEL_DAY') }})
              </q-item-label>
              <q-item-label v-if="item.conflicts.length > 0" caption class="text-warning">
                <q-icon name="warning" size="12px" class="q-mr-xs" />
                {{ t('LABEL_CONFLICT_WITH') }}: {{ item.conflicts.join(', ') }}
              </q-item-label>
            </q-item-section>
            <q-item-section side>
              <div class="row q-gutter-xs">
                <q-btn flat round color="positive" icon="check_circle" @click="acceptAbsence(item)">
                  <q-tooltip>{{ t('TITLE_ACCEPT') }}</q-tooltip>
                </q-btn>
                <q-btn flat round color="negative" icon="cancel" @click="declineAbsence(item)">
                  <q-tooltip>{{ t('TITLE_DECLINE') }}</q-tooltip>
                </q-btn>
              </div>
            </q-item-section>
          </q-item>
        </q-list>
      </q-card>

      <q-card flat bordered>
        <AbsenceSummaryTableComponent v-model="absenceSummaryItems" />
      </q-card>
    </div>
  </q-page>
</template>

<style scoped>
.dashboard-header {
  background: linear-gradient(135deg, var(--q-primary) 0%, #1565c0 100%);
}
</style>
