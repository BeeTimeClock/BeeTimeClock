<script setup lang="ts">
import { showInfoMessage } from 'src/helper/message';
import { UserApikey, UserApikeyCreateRequest } from 'src/models/Base';
import BeeTimeClock from 'src/service/BeeTimeClock';
import { onMounted, ref } from 'vue';
import { useI18n } from 'vue-i18n';

  const apikeys = ref([] as UserApikey[])
  const promptApikeyAdd = ref(false)
  const promptApikeySuccess = ref(false)
  const apikeyCreateRequest = ref(null as UserApikeyCreateRequest|null)
  const apikeySuccessResponse = ref(null as UserApikey|null)

  const { t } = useI18n()

  function loadUserApikeys() {
    BeeTimeClock.getUserApikey().then(result => {
      if (result.status === 200) {
        apikeys.value = result.data.Data;
      }
    })
  }

  function showApikeyAdd() {
    apikeyCreateRequest.value = {} as UserApikeyCreateRequest
    promptApikeyAdd.value = true;
  }

  function showApikeySuccess(apikeySuccess: UserApikey) {
    apikeySuccessResponse.value = apikeySuccess;
    promptApikeySuccess.value = true
    console.log('show success');
  }

  function closeApikey() {
    apikeySuccessResponse.value = null;
    promptApikeySuccess.value = false;
  }

  function createApikey() {
    BeeTimeClock.createUserApikey(apikeyCreateRequest.value).then(result => {
      if (result.status === 201) {
        showInfoMessage(t('MSG_CREATE_SUCCESS'))
        showApikeySuccess(result.data.Data)
        loadUserApikeys()
      }
    })
  }

  onMounted(() => {
    loadUserApikeys();
  })
</script>

<template>
  <q-page padding>
    <q-markup-table>
      <thead>
        <tr>
          <th>{{ $t('LABEL_DESCRIPTION') }}</th>
          <th>{{ $t('LABEL_VALID_TILL') }}</th>
          <th>
            <q-btn color="primary" icon="add" @click="showApikeyAdd"/>
          </th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="apikey in apikeys" :key="apikey.ID">
          <td>{{ apikey.Description }}</td>
          <td>{{ apikey.ValidTill }}</td>
          <td></td>
        </tr>
      </tbody>
    </q-markup-table>
  </q-page>
  <q-dialog persistent v-model="promptApikeyAdd">
    <q-card>
      <q-card-section class="bg-primary text-white text-h6">
        {{ $t('LABEL_CREATE', {item: $t('LABEL_APIKEY')}) }}
      </q-card-section>
      <q-card-section>
        <q-input v-model="apikeyCreateRequest.Description" :label="$t('LABEL_DESCRIPTION')"/>
        <q-input type="date" v-model="apikeyCreateRequest.ValidTill" :label="$t('LABEL_VALID_TILL') + '(optional)'"/>
      </q-card-section>
      <q-card-section>
        <q-card-actions>
          <q-btn color="negative" v-close-popup :label="$t('BTN_CANCEL')"/>
          <q-btn color="positive" v-close-popup :label="$t('BTN_CREATE')" @click="createApikey"/>
        </q-card-actions>
      </q-card-section>
    </q-card>
  </q-dialog>
  <q-dialog persistent v-model="promptApikeySuccess">
    <q-card>
      <q-card-section class="bg-primary text-white text-h6">
        {{ $t('MSG_CREATE_SUCCESS', {item: $t('LABEL_APIKEY')}) }}
      </q-card-section>
      <q-card-section>
        {{ $t('MSG_APIKEY_SHOW_WARNING') }}
        <pre>{{ apikeySuccessResponse.Apikey }}</pre>
      </q-card-section>
      <q-card-section>
        <q-card-actions>
          <q-btn color="negative" v-close-popup :label="$t('BTN_CLOSE')" @click="closeApikey"/>
        </q-card-actions>
      </q-card-section>
    </q-card>
  </q-dialog>
</template>

<style scoped>

</style>
