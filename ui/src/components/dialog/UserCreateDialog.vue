<script setup lang="ts">
import { ref } from 'vue';
import BeeTimeClock from 'src/service/BeeTimeClock';
import type { UserCreateRequest } from 'src/models/Authentication';
import { useI18n } from 'vue-i18n';
import { showErrorMessage, showInfoMessage } from 'src/helper/message';
import type { ErrorResponse } from 'src/models/Base';

const { t } = useI18n();
const emit = defineEmits<{ created: []; 'update:show': [value: boolean] }>();
defineProps<{ show: boolean }>();

const emptyForm = (): UserCreateRequest => ({
  Username: '',
  Password: '',
  AccessLevel: 'user',
  FirstName: '',
  LastName: '',
});

const form = ref<UserCreateRequest>(emptyForm());

const accessLevelOptions = [
  { value: 'user', label: t('LABEL_USER', 1) },
  { value: 'admin', label: t('LABEL_ADMINISTRATOR', 1) },
];

function close() {
  emit('update:show', false);
}

function submit() {
  BeeTimeClock.administrationCreateUser(form.value)
    .then(result => {
      if (result.status === 201) {
        showInfoMessage(t('MSG_CREATE_SUCCESS', { item: t('LABEL_USER', 1) }));
        emit('created');
        close();
        form.value = emptyForm();
      }
    })
    .catch((error: ErrorResponse) => showErrorMessage(error.response?.data.Message));
}
</script>

<template>
  <q-dialog :model-value="show" @update:model-value="emit('update:show', $event)">
    <q-card style="min-width: 400px">
      <q-card-section>
        <div class="text-h6">{{ t('TITLE_CREATE', { item: t('LABEL_USER', 1) }) }}</div>
      </q-card-section>
      <q-card-section>
        <q-form @submit="submit" class="q-gutter-sm">
          <q-input
            :label="t('LABEL_USERNAME')"
            v-model="form.Username"
            :rules="[v => !!v || t('LABEL_FIELD_REQUIRED')]"
            autofocus
          />
          <q-input
            :label="t('LABEL_PASSWORD')"
            v-model="form.Password"
            type="password"
            :rules="[v => !!v || t('LABEL_FIELD_REQUIRED')]"
          />
          <q-input :label="t('LABEL_FIRST_NAME')" v-model="form.FirstName" />
          <q-input :label="t('LABEL_LAST_NAME')" v-model="form.LastName" />
          <q-select
            :label="t('LABEL_ACCESS_LEVEL')"
            v-model="form.AccessLevel"
            :options="accessLevelOptions"
            emit-value
            map-options
          />
          <q-card-actions align="right">
            <q-btn flat :label="t('BTN_CANCEL')" @click="close" />
            <q-btn color="primary" :label="t('BTN_CREATE')" type="submit" />
          </q-card-actions>
        </q-form>
      </q-card-section>
    </q-card>
  </q-dialog>
</template>
