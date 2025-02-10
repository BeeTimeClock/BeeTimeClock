<script setup lang="ts">
import BeeTimeClock from 'src/service/BeeTimeClock';
import { onMounted, ref } from 'vue';
import { Team } from 'src/models/Team';
import TeamCreateDialog from 'components/dialog/TeamCreateDialog.vue';
import { QTableColumn } from 'quasar';
import { emptyPagination } from 'src/helper/objects';
import { useI18n } from 'vue-i18n';

const {t} = useI18n();
const teams = ref<Team[]>([]);
const showTeamCreateDialog = ref(false);
const columns: QTableColumn[] = [
  {
    name: 'id',
    field: 'ID',
    label: t('LABEL_ID'),
    align: 'left',
  },
  {
    name: 'teamName',
    field: 'Teamname',
    label: t('LABEL_TEAM_NAME'),
    align: 'left',
  },
  {
    name: 'teamLead',
    field: row => row.teamOwnerMapped.displayName,
    label: t('LABEL_TEAM_LEAD'),
    align: 'left',
  },
  {
    name: 'actions',
    field: '',
    label: t('LABEL_ACTION',2),
    align: 'left',
  }
];

function loadTeams() {
  BeeTimeClock.administrationGetTeams(true).then((result) => {
    if (result.status === 200) {
      teams.value = result.data.Data.map((s) => Team.fromApi(s));
    }
  });
}

onMounted(() => {
  loadTeams();
});
</script>

<template>
  <q-page padding>
    <q-btn
      :label="$t('LABEL_CREATE', {item: $t('LABEL_TEAM')})"
      class="full-width" color="positive" icon="add" @click="showTeamCreateDialog = true;"/>
    <q-table class="q-mt-md" :rows="teams" :columns="columns" hide-pagination :pagination="emptyPagination">
      <template v-slot:header="props">
        <q-tr :props="props" class="bg-primary text-white">
          <q-th
            v-for="col in props.cols"
            :key="col.name"
            :props="props"
          >
            {{ col.label }}
          </q-th>
        </q-tr>
      </template>

      <template v-slot:body="props">
        <q-tr :props="props">
          <q-td
            v-for="col in props.cols"
            :key="col.name"
            :props="props"
          >
            <div v-if="col.name == 'actions'">
              <q-btn icon="edit" color="primary" :to="{name: 'AdministrationTeamDetail', params: {teamId: props.row.ID}}"/>
            </div>
            <div v-else>{{ col.value }}</div>
          </q-td>
        </q-tr>
      </template>
    </q-table>
    <TeamCreateDialog v-model:show="showTeamCreateDialog" @created="loadTeams"/>
  </q-page>
</template>

<style scoped></style>
