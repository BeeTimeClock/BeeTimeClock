<script setup lang="ts">
import { computed, onMounted, ref } from 'vue';
import { Team, TeamMember } from 'src/models/Team';
import { useRoute } from 'vue-router';
import BeeTimeClock from 'src/service/BeeTimeClock';
import { emptyPagination } from 'src/helper/objects';
import { QTableColumn, useQuasar } from 'quasar';
import { useI18n } from 'vue-i18n';
import TeamMemberAddDialog from 'components/dialog/TeamMemberAddDialog.vue';
import { showInfoMessage } from 'src/helper/message';

const route = useRoute();
const team = ref<Team>();
const members = ref<TeamMember[]>();
const { t } = useI18n();
const showTeamMemberAdd = ref(false);
const q = useQuasar();

const columns: QTableColumn[] = [
  {
    name: 'displayName',
    field: (row) => row.userMapped.displayName,
    label: t('LABEL_NAME'),
    align: 'left',
  },
  {
    name: 'actions',
    field: '',
    label: t('LABEL_ACTION', 2),
    align: 'right',
  },
];

const teamId = computed(() => {
  return parseInt(route.params.teamId as string);
});

function loadTeam() {
  BeeTimeClock.administrationGetTeam(teamId.value, true).then((result) => {
    if (result.status === 200) {
      team.value = Team.fromApi(result.data.Data);
    }
  });
}

function loadMembers() {
  BeeTimeClock.administrationGetTeamMembers(teamId.value, true).then(
    (result) => {
      if (result.status === 200) {
        members.value = result.data.Data.map((s) => TeamMember.fromApi(s));
      }
    }
  );
}

function deleteTeamMember(teamMember: TeamMember) {
  if (!team.value) return;

  q.dialog({
    message: t('MSG_DELETE', {
      item: t('LABEL_TEAM_MEMBER'),
      identifier: teamMember.userMapped.displayName,
    }),
    cancel: true,
    persistent: true,
  }).onOk(() => {
    if (!team.value) return;
    BeeTimeClock.administrationDeleteTeamMember(
      team.value.ID,
      teamMember.ID
    ).then((result) => {
      if (result.status === 204) {
        loadMembers();
        showInfoMessage(
          t('MSG_DELETE_SUCCESS', {
            item: t('LABEL_TEAM_MEMBER'),
            identifier: teamMember.userMapped.displayName,
          })
        );
      }
    });
  });
}

onMounted(() => {
  loadTeam();
  loadMembers();
});
</script>

<template>
  <q-page padding>
    <q-card v-if="team">
      <q-card-section class="bg-primary text-white text-h6">
        {{ $t('LABEL_INFORMATION') }}
      </q-card-section>
      <q-card-section>
        <q-list>
          <q-item>
            <q-item-section side>
              <q-item-label caption>{{ $t('LABEL_ID') }}</q-item-label>
              <q-item-label>{{ team.ID }}</q-item-label>
            </q-item-section>
            <q-item-section>
              <q-item-label caption>{{ $t('LABEL_CREATED_AT') }}</q-item-label>
              <q-item-label>{{ team.CreatedAt }}</q-item-label>
            </q-item-section>
          </q-item>
          <q-item>
            <q-item-section>
              <q-item-label caption>{{ $t('LABEL_TEAM_LEAD') }}</q-item-label>
              <q-item-label
                >{{ team.teamOwnerMapped.displayName }}
              </q-item-label>
            </q-item-section>
          </q-item>
        </q-list>
      </q-card-section>
    </q-card>
    <q-card class="q-mt-md q-pa-none">
      <q-card-section class="bg-primary text-white text-h6">
        {{ $t('LABEL_TEAM_MEMBER', 2) }}
      </q-card-section>
      <q-card-section class="q-pa-none">
        <q-table
          :pagination="emptyPagination"
          hide-pagination
          :columns="columns"
          :rows="members"
          flat
          square
        >
          <template v-slot:header="props">
            <q-tr :props="props" class="bg-primary text-white">
              <q-th v-for="col in props.cols" :key="col.name" :props="props">
                <div v-if="col.name == 'actions'">
                  <q-btn
                    icon="add"
                    color="positive"
                    @click="showTeamMemberAdd = true"
                  />
                </div>
                <div v-else>
                  {{ col.label }}
                </div>
              </q-th>
            </q-tr>
          </template>

          <template v-slot:body="props">
            <q-tr :props="props">
              <q-td v-for="col in props.cols" :key="col.name" :props="props">
                <div v-if="col.name == 'actions'">
                  <q-btn icon="delete" color="negative" @click="deleteTeamMember(props.row)"/>
                </div>
                <div v-else>{{ col.value }}</div>
              </q-td>
            </q-tr>
          </template>
        </q-table>
      </q-card-section>
    </q-card>
    <TeamMemberAddDialog
      v-if="team"
      v-model:show="showTeamMemberAdd"
      :team="team"
      @created="loadMembers"
    />
  </q-page>
</template>

<style scoped></style>
