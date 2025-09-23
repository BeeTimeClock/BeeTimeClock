<script setup lang="ts">
import {Team} from 'src/models/Team';
import BeeTimeClock from 'src/service/BeeTimeClock';
import {onMounted, ref, watch} from 'vue';
import {useI18n} from 'vue-i18n';
import {showErrorMessage} from 'src/helper/message';
import {type AbsenceSummaryItem} from 'src/models/Absence';
import AbsenceSummaryTableComponent from 'components/AbsenceSummaryTableComponent.vue';

const {t} = useI18n();
const selectedTeam = ref<Team>();
const teams = ref<Team[]>([]);
const teamAbsenceSummaries = ref<AbsenceSummaryItem[]>([])

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

watch(selectedTeam, () => {
  loadTeamAbensces();
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
      <q-table class="q-mb-md" :rows="selectedTeam.Members"/>

      <AbsenceSummaryTableComponent v-model="teamAbsenceSummaries"/>
    </div>
  </q-page>
</template>

<style scoped>

</style>
