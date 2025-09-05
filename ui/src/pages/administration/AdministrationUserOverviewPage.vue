<script setup lang="ts">
import { computed, onMounted, ref } from 'vue';
import { User } from 'src/models/Authentication';
import BeeTimeClock from 'src/service/BeeTimeClock';
import UserAbsenceSummary from 'components/UserAbsenceSummary.vue';
import type { AbsenceReason } from 'src/models/Absence';
import { useQuasar } from 'quasar';
import { useI18n } from 'vue-i18n';
import { showErrorMessage, showInfoMessage } from 'src/helper/message';
import type { ErrorResponse } from 'src/models/Base';

const q = useQuasar();
const { t } = useI18n();

const users = ref([] as User[]);
const absenceReasons = ref([] as AbsenceReason[]);
const needle = ref('');

function loadUsers() {
  BeeTimeClock.administrationGetUsers(true)
    .then((result) => {
      if (result.status === 200) {
        users.value = result.data.Data.map(s => User.fromApi(s));
      }
    })
    .catch((error: ErrorResponse) => {
      showErrorMessage(error.message);
    });
}

function loadAbsenceReasons() {
  BeeTimeClock.absenceReasons().then((result) => {
    absenceReasons.value = result.data.Data;
  }).catch((error: ErrorResponse) => {
    showErrorMessage(error.message);
  });
}

const sortedFilteredUsers = computed(() => {
  const search = needle.value.toLowerCase();
  const data = users.value.filter((s) => {
    if (s.LastName.toLowerCase().indexOf(search) >= 0) return true;
    if (s.FirstName.toLowerCase().indexOf(search) >= 0) return true;
    if (s.Username.toLowerCase().indexOf(search) >= 0) return true;
  });
  return data.sort((a, b) => a.LastName.localeCompare(b.LastName));
});

function deleteUser(user: User) {
  q.dialog({
    message: t('MSG_DELETE', {
      item: t('LABEL_USER'),
      identifier: user.Username,
    }),
    cancel: true,
    persistent: true,
  }).onOk(() => {
    BeeTimeClock.administrationDeleteUser(user).then((result) => {
      if (result.status === 204) {
        loadUsers();
        showInfoMessage(
          t('MSG_DELETE_SUCCESS', {
            item: t('LABEL_USER'),
            identifier: user.Username,
          }),
        );
      }
    }).catch((error: ErrorResponse) => {
      showErrorMessage(error.message);
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
    <q-input :label="t('LABEL_SEARCH')" v-model="needle" />
    <q-markup-table class="q-mt-lg">
      <thead>
        <tr class="bg-primary text-white">
          <th class="text-left">{{ t('LABEL_USER') }}</th>
          <th class="text-left">{{ t('LABEL_FIRST_NAME') }}</th>
          <th class="text-left">{{ t('LABEL_LAST_NAME') }}</th>
          <th class="text-left">{{ t('LABEL_ABSENCE', 2) }}</th>
          <th></th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="user in sortedFilteredUsers" :key="user.Username">
          <td>{{ user.Username }}</td>
          <td>{{ user.FirstName }}</td>
          <td>{{ user.LastName }}</td>
          <td>
            <UserAbsenceSummary
              :user-id="user.ID"
              :absence-reasons="absenceReasons"
            />
          </td>
          <td>
            <q-btn
              class="q-ml-md"
              color="primary"
              icon="visibility"
              :to="{
                name: 'AdministrationUserDetail',
                params: { userId: user.ID },
              }"
            />
            <q-btn
              class="q-ml-md"
              color="negative"
              icon="delete"
              @click="deleteUser(user)"
            />
          </td>
        </tr>
      </tbody>
    </q-markup-table>
  </q-page>
</template>

<style scoped></style>
