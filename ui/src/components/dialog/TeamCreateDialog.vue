<script setup lang="ts">
import { ref } from 'vue';
import { ApiTeamCreateRequest } from 'src/models/Team';
import BeeTimeClock from 'src/service/BeeTimeClock';
import { User } from 'src/models/Authentication';
import { showInfoMessage } from 'src/helper/message';
import { useI18n } from 'vue-i18n';
import { useRouter } from 'vue-router';

const router = useRouter();
const teamCreateRequest = ref<ApiTeamCreateRequest>({} as ApiTeamCreateRequest);
const users = ref<User[]>([]);
const show = defineModel('show', {type: Boolean, default: false})
const gotoDetailAfterCreation = defineModel('goto-detail', {type: Boolean, default: false})
const {t} = useI18n();

const emits = defineEmits(['created'])

function createTeam() {
  BeeTimeClock.administrationCreateTeam(teamCreateRequest.value).then(result => {
    if (result.status == 201) {
      show.value = false;
      emits('created', result.data.Data)
      showInfoMessage(t('MSG_CREATE_SUCCESS'))
      if (gotoDetailAfterCreation.value) {
        router.push({name: 'AdministrationTeamDetail', params: {teamId: result.data.Data.ID}})
      }
    }
  })
}

function loadData() {
  loadUsers();
}

function loadUsers() {
  BeeTimeClock.administrationGetUsers().then(result => {
    if (result.status === 200) {
      users.value = result.data.Data.map(s => User.fromApi(s));
    }
  })
}
</script>

<template>
  <q-dialog v-model="show" @before-show="loadData">
    <q-card>
      <q-card-section class="bg-primary text-white text-h6">
        {{ $t('LABEL_CREATE', { item: $t('LABEL_TEAM') }) }}
      </q-card-section>
      <q-form @submit="createTeam">
        <q-card-section>
          <q-input v-model="teamCreateRequest.Teamname" :label="$t('LABEL_TEAM_NAME')"/>
          <q-select v-model="teamCreateRequest.TeamOwnerID" :label="$t('LABEL_TEAM_LEAD')"
                    :options="users"
                    emit-value
                    map-options
                    option-value="ID"
                    option-label="displayName"
          />
        </q-card-section>
        <q-card-section>
          <q-card-actions>
            <q-btn color="negative" v-close-popup :label="$t('LABEL_CANCEL')"/>
            <q-btn color="positive" v-close-popup type="submit" :label="$t('LABEL_CREATE')"/>
          </q-card-actions>
        </q-card-section>
      </q-form>
    </q-card>
  </q-dialog>
</template>

<style scoped></style>
