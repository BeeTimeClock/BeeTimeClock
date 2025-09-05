<script setup lang="ts">
import type { ApiTeamMemberCreateRequest} from 'src/models/Team';
import { Team } from 'src/models/Team';
import { User } from 'src/models/Authentication';
import { ref } from 'vue';
import BeeTimeClock from 'src/service/BeeTimeClock';
import { showErrorMessage, showInfoMessage } from 'src/helper/message';
import { useI18n } from 'vue-i18n';
import type { ErrorResponse } from 'src/models/Base';

const { t } = useI18n();
const team = defineModel('team', { type: Team, required: true });
const show = defineModel('show', { type: Boolean, default: false });
const users = ref<User[]>([]);
const teamMemberCreateRequest = ref<ApiTeamMemberCreateRequest>(
  {} as ApiTeamMemberCreateRequest,
);
const emits = defineEmits(['created']);

function createMember() {
  BeeTimeClock.administrationCreateTeamMember(
    team.value.ID,
    teamMemberCreateRequest.value,
  )
    .then((result) => {
      if (result.status === 201) {
        emits('created');
        showInfoMessage(t('MSG_CREATE_SUCCESS'));
      }
    })
    .catch((error: ErrorResponse) => {
      showErrorMessage(error.message);
    });
}

function loadUsers() {
  BeeTimeClock.administrationGetUsers().then((result) => {
    if (result.status === 200) {
      users.value = result.data.Data.map((s) => User.fromApi(s));
    }
  }).catch((error: ErrorResponse) => {
    showErrorMessage(error.message);
  });;
}
</script>

<template>
  <q-dialog v-model="show" @before-show="loadUsers">
    <q-card>
      <q-card-section class="bg-primary text-white text-h6">
        {{ t('LABEL_CREATE', { item: t('LABEL_TEAM_MEMBER') }) }}
      </q-card-section>
      <q-form @submit="createMember">
        <q-card-section>
          <q-select
            v-model="teamMemberCreateRequest.UserID"
            :options="users"
            map-options
            emit-value
            option-label="displayName"
            option-value="ID"
            :label="t('LABEL_USER')"
          />
        </q-card-section>
        <q-card-section>
          <q-card-actions>
            <q-btn color="negative" v-close-popup :label="t('LABEL_CANCEL')" />
            <q-btn
              color="positive"
              v-close-popup
              type="submit"
              :label="t('LABEL_CREATE')"
            />
          </q-card-actions>
        </q-card-section>
      </q-form>
    </q-card>
  </q-dialog>
</template>

<style scoped></style>
