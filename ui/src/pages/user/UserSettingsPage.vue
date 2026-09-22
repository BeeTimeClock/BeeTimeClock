<script setup lang="ts">
import { onMounted, ref } from 'vue';
import {type User } from 'src/models/Authentication';
import BeeTimeClock from 'src/service/BeeTimeClock';
import { showErrorMessage, showInfoMessage } from 'src/helper/message';
import { useI18n } from 'vue-i18n';
import { type ErrorResponse } from 'src/models/Base';

const user = ref<User>();
const { t } = useI18n();

const MAX_RINGTONE_TONES = 24;

function ringtoneRule(val: string) {
  if (!val) return true;
  const notes = val.split(':').pop() ?? '';
  const count = notes.split(',').filter((n) => n.trim() !== '').length;
  return (
    count <= MAX_RINGTONE_TONES ||
    t('MSG_RINGTONE_TOO_LONG', { max: MAX_RINGTONE_TONES })
  );
}

function loadUser() {
  BeeTimeClock.getMeUser().then((result) => {
    if (result.status === 200) {
      user.value = result.data.Data;
    }
  }).catch((error: ErrorResponse) => {
    showErrorMessage(error.response?.data.Message);
  });
}

function save() {
  if (!user.value) return;
  BeeTimeClock.updateMeUser(user.value).then((result) => {
    if (result.status === 200) {
      user.value = result.data.Data;
      showInfoMessage(t('MSG_UPDATE_SUCCESS'));
    }
  }).catch((error: ErrorResponse) => {
    showErrorMessage(error.response?.data.Message);
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
        :label="t('LABEL_STAFF_NUMBER')"
        v-model.number="user.StaffNumber"
      />
      <q-toggle
        :label="t('LABEL_ALLOW_GRAVATAR')"
        v-model="user.AllowGravatar"
        class="q-mt-md"
      />
      <q-input
        :label="t('LABEL_COMING_RINGTONE')"
        :hint="t('HINT_RINGTONE_RTTTL')"
        v-model="user.ComingRingtone"
        :rules="[ringtoneRule]"
        class="q-mt-md"
      />
      <q-input
        :label="t('LABEL_GOING_RINGTONE')"
        :hint="t('HINT_RINGTONE_RTTTL')"
        v-model="user.GoingRingtone"
        :rules="[ringtoneRule]"
        class="q-mt-md"
      />
      <q-btn
        class="full-width q-mt-lg"
        :label="t('LABEL_SAVE')"
        icon="save"
        color="positive"
        type="submit"
      />
    </q-form>
  </q-page>
</template>

<style scoped></style>
