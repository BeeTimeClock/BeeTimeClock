<script setup lang="ts">
import {onMounted, ref} from 'vue';
import {User} from 'src/models/Authentication';
import BeeTimeClock from 'src/service/BeeTimeClock';

const users = ref([] as User[]);

function loadUsers() {
  BeeTimeClock.administrationGetUsers(true).then(result => {
    if (result.status === 200) {
      users.value = result.data.Data;
    }
  })
}

onMounted(() => {
  loadUsers();
})
</script>

<template>
  <q-page padding>
    <q-markup-table>
      <thead>
      <tr>
        <th class="text-left">{{ $t('LABEL_USER') }}</th>
        <th class="text-left">{{ $t('LABEL_FIRST_NAME') }}</th>
        <th class="text-left">{{ $t('LABEL_LAST_NAME') }}</th>
        <th></th>
      </tr>
      </thead>
      <tbody>
      <tr v-for="user in users" :key="user.Username">
        <td>{{ user.Username }}</td>
        <td>{{ user.FirstName }}</td>
        <td>{{ user.LastName }}</td>
        <td>
          <q-btn color="primary" icon="edit" :to="{name: 'AdministrationUserDetail', params: { userId: user.ID }}"/>
        </td>
      </tr>
      </tbody>
    </q-markup-table>
  </q-page>
</template>

<style scoped>

</style>
