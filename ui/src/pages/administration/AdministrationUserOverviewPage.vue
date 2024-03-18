<script setup lang="ts">
import { onMounted, ref } from 'vue';
import { User } from 'src/models/Authentication';
import BeeTimeClock from 'src/service/BeeTimeClock';
import UserAbsenceSummary from 'components/UserAbsenceSummary.vue';
import { AbsenceReason } from 'src/models/Absence';
import { useQuasar } from 'quasar';
import { useI18n } from 'vue-i18n';
import { showInfoMessage } from 'src/helper/message';

const q = useQuasar();
const { t } = useI18n();

const users = ref([] as User[]);
const absenceReasons = ref([] as AbsenceReason[]);

function loadUsers() {
  BeeTimeClock.administrationGetUsers(true).then(result => {
    if (result.status === 200) {
      users.value = result.data.Data;
    }
  });
}

function loadAbsenceReasons() {
  BeeTimeClock.absenceReasons().then(result => {
    absenceReasons.value = result.data.Data;
  });
}

function deleteUser(user: User) {
  q.dialog({
    message: t('MSG_DELETE', { item: t('LABEL_USER'), identifier: user.Username }),
    cancel: true,
    persistent: true
  }).onOk(() => {
    BeeTimeClock.administrationDeleteUser(user).then(result => {
      if (result.status === 204) {
        loadUsers();
        showInfoMessage(t('MSG_DELETE_SUCCESS', { item: t('LABEL_USER'), identifier: user.Username }));
      }
    });
  });
}

onMounted(() => {
  loadAbsenceReasons();
  loadUsers();
});
</script>

<template>
  <q-page padding>
    <q-markup-table>
      <thead>
      <tr class="bg-primary text-white">
        <th class="text-left">{{ $t('LABEL_USER') }}</th>
        <th class="text-left">{{ $t('LABEL_FIRST_NAME') }}</th>
        <th class="text-left">{{ $t('LABEL_LAST_NAME') }}</th>
        <th class="text-left">{{ $t('LABEL_ABSENCES') }}</th>
        <th></th>
      </tr>
      </thead>
      <tbody>
      <tr v-for="user in users" :key="user.Username">
        <td>{{ user.Username }}</td>
        <td>{{ user.FirstName }}</td>
        <td>{{ user.LastName }}</td>
        <td>
          <UserAbsenceSummary :user-id="user.ID" :absence-reasons="absenceReasons" />
        </td>
        <td>
          <q-btn class="q-ml-md" color="primary" icon="visibility"
                 :to="{name: 'AdministrationUserDetail', params: { userId: user.ID }}" />
          <q-btn class="q-ml-md" color="negative" icon="delete" @click="deleteUser(user)" />        </td>
      </tr>
      </tbody>
    </q-markup-table>
  </q-page>
</template>

<style scoped>

</style>
