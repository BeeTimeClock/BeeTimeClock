<script setup lang="ts">

import {User} from 'src/models/Authentication';
import {onMounted, ref} from 'vue';
import BeeTimeClock from 'src/service/BeeTimeClock';
import {useRoute} from 'vue-router';
import {useI18n} from 'vue-i18n';
import {showInfoMessage} from 'src/helper/message';

const {t} = useI18n();

const route = useRoute();
const userId = route.params.userId as string;
const user = ref(null as User | null);
const accessLevelOptions = [
  {
    value: 'admin',
    label: t('LABEL_ADMINISTRATOR'),
  },
  {
    value: 'user',
    label: t('LABEL_USER'),
  }
]

function loadUser() {
  BeeTimeClock.administrationGetUserById(userId).then(result => {
    if (result.status === 200) {
      user.value = result.data.Data;
    }
  })
}

function saveUser() {
  BeeTimeClock.administrationUpdateUser(user.value as User).then(result => {
    if (result.status === 200) {
      user.value = result.data.Data;
      showInfoMessage(t('MSG_UPDATE_SUCCESS'))
    }
  })
}

onMounted(() => {
  loadUser();
})

</script>

<template>
  <q-page padding>
    <div v-if="user">
      <q-card>
        <q-card-section>
          <q-input readonly :label="$t('LABEL_USERNAME')" v-model="user.Username"/>
          <q-input :label="$t('LABEL_FIRST_NAME')" v-model="user.FirstName"/>
          <q-input :label="$t('LABEL_LAST_NAME')" v-model="user.LastName"/>
          <q-select :label="$t('LABEL_ACCESS_LEVEL')" :options="accessLevelOptions" v-model="user.AccessLevel" map-options emit-value/>
        </q-card-section>
        <q-card-actions>
          <q-btn :label="$t('BTN_SAVE')" color="primary" @click="saveUser"/>
        </q-card-actions>
      </q-card>
    </div>
  </q-page>
</template>

<style scoped>

</style>
