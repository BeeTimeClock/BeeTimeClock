<script setup lang="ts">
import { onMounted, ref } from 'vue';
import { User } from 'src/models/Authentication';
import BeeTimeClock from 'src/service/BeeTimeClock';
import { showInfoMessage } from 'src/helper/message';
import { useI18n } from 'vue-i18n';

const user = ref<User>();
const { t } = useI18n();

function loadUser() {
  BeeTimeClock.getMeUser().then((result) => {
    if (result.status === 200) {
      user.value = result.data.Data;
    }
  });
}

function save() {
  if (!user.value) return;
  BeeTimeClock.updateMeUser(user.value).then((result) => {
    if (result.status === 200) {
      user.value = result.data.Data;
      showInfoMessage(t('MSG_UPDATE_SUCCESS'));
    }
  });
}

onMounted(() => {
  loadUser();
});
</script>

<template>
  <q-page padding>
    <q-form v-if="user" @submit="save">
      <q-input
        :label="$t('LABEL_STAFF_NUMBER')"
        v-model.number="user.StaffNumber"
      />
      <q-btn
        class="full-width q-mt-lg"
        :label="$t('LABEL_SAVE')"
        icon="save"
        color="positive"
        type="submit"
      />
    </q-form>
  </q-page>
</template>

<style scoped></style>
