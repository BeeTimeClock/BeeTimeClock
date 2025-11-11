<script setup lang="ts">
import {Team, TeamLevel, type TeamMember} from 'src/models/Team';
import BeeTimeClock from 'src/service/BeeTimeClock';
import {computed, onMounted, ref, watch} from 'vue';
import {useI18n} from 'vue-i18n';
import {showErrorMessage} from 'src/helper/message';
import {Absence, type AbsenceSummaryItem} from 'src/models/Absence';
import AbsenceSummaryTableComponent from 'components/AbsenceSummaryTableComponent.vue';
import {type User} from 'src/models/Authentication';
import {type QTableColumn} from 'quasar';
import {emptyPagination} from 'src/helper/objects';
import {useAuthStore} from 'stores/microsoft-auth';

const {t} = useI18n();
const auth = useAuthStore();
const selectedTeam = ref<Team>();
const teams = ref<Team[]>([]);
const teamAbsenceSummaries = ref<AbsenceSummaryItem[]>([])
const neededApprovals = ref<Absence[]>([])
const columns = [
  {
    name: 'user',
    field: (row: TeamMember) => row.userMapped,
    label: t('LABEL_USER'),
    align: 'left',
    format: (val: User) => val ? val.displayName : '-',
  },
  {
    name: 'user',
    field: 'Level',
    label: t('LABEL_USER'),
    align: 'left',
  }
] as QTableColumn[];

const isLead = computed(() => {
  return selectedTeam.value?.Members.find(s => s.UserID === auth.getSession()?.ID && (s.Level === TeamLevel.LeadSurrogate || s.Level === TeamLevel.Lead)) != null
})

function loadMyTeams() {
  BeeTimeClock.getTeams().then(result => {
    if (result.status === 200) {
      teams.value = result.data.Data.map(s => Team.fromApi(s))
      if (teams.value.length >= 1) {
        selectedTeam.value = teams.value[0]!;
      }
    }
  }).catch((error) => {
    showErrorMessage(error)
  })
}

function loadTeamAbensces() {
  if (!selectedTeam.value) return;

  BeeTimeClock.queryTeamAbsenceSummary(selectedTeam.value.ID).then(result => {
    if (result.status === 200) {
      teamAbsenceSummaries.value = result.data.Data;
    }
  }).catch((error) => {
    showErrorMessage(error);
  })
}

function loadNeededApprovals() {
  if (!selectedTeam.value) return
  BeeTimeClock.absenceApprovalTeamOpen(selectedTeam.value.ID).then(result => {
    if (result.status === 200) {
      neededApprovals.value = result.data.Data.map(s => Absence.fromApi(s));
    }
  }).catch((error) => {
    showErrorMessage(error);
  })
}

watch(selectedTeam, () => {
  loadTeamAbensces();
  loadNeededApprovals();
})

onMounted(() => {
  loadMyTeams()
  loadTeamAbensces();
})
</script>

<template>
  <q-page padding>
    <q-select :label="t('LABEL_TEAM')" v-model="selectedTeam" :options="teams" :readonly="teams.length == 1" emit-value
              map-options option-label="Teamname" class="q-mb-lg"/>

    <div v-if="selectedTeam">
      <q-table class="q-mb-md" :columns="columns" :rows="selectedTeam.membersMapped" :pagination="emptyPagination"
               hide-pagination/>

      <q-table v-if="isLead" :rows="neededApprovals"/>

      <AbsenceSummaryTableComponent v-model="teamAbsenceSummaries"/>
    </div>
  </q-page>
</template>

<style scoped>

</style>
