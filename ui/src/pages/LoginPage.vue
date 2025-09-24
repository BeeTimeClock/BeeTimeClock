<template>
  <q-layout view="hHh Lpr fFf"> <!-- Be sure to play with the Layout demo on docs -->
    <q-page-container>
      <q-page class="flex flex-center bg-grey-2">
        <q-card class="q-pa-md shadow-2 my_card" bordered>
          <q-card-section class="text-center">
            <img width="256" :src="'logo.svg'" alt="btc-logo"/>
            <div class="text-grey-9 text-h3 text-weight-bold">Bee Time Clock</div>
            <div class="text-grey-9 text-h5 text-weight-bold">{{ t('LABEL_SIGN_IN') }}</div>
          </q-card-section>
          <q-card-section>
            <div class="text-center text-weight-thin q-mb-sm">{{ t('LABEL_LOCAL_SIGN_IN') }}</div>
            <q-input dense outlined v-model="email" :label="t('LABEL_USERNAME')"></q-input>
            <q-input dense outlined class="q-mt-md" v-model="password" type="password"
                     :label="t('LABEL_PASSWORD')"></q-input>
            <q-btn color="primary" size="md" :label="t('LABEL_SIGN_IN')" no-caps class="full-width q-mt-lg"
                   @click="loginLocal">
            </q-btn>
          </q-card-section>
          <q-separator/>
          <q-card-section>
            <div class="text-center text-weight-thin q-mb-sm">{{ t('LABEL_SSO_SIGN_IN') }}</div>
            <q-btn v-if="authProviders?.Microsoft" color="primary" size="md" :label="t('LABEL_SIGN_IN_MICROSOFT')"
                   no-caps class="full-width"
                   @click="loginWithMicrosoft">
              <q-avatar square class="q-ml-md" size="24px">
                <img :src="'ms_login.png'" alt="ms-logo"/>
              </q-avatar>
            </q-btn>
          </q-card-section>
        </q-card>
      </q-page>
    </q-page-container>
  </q-layout>
</template>

<script lang="ts" setup>
import type {AuthProviders, ErrorResponse} from 'src/models/Base';
import BeeTimeClock from 'src/service/BeeTimeClock';
import {useAuthStore} from 'stores/microsoft-auth';
import {onMounted, ref} from 'vue';
import {showErrorMessage} from 'src/helper/message';
import {useRouter} from 'vue-router';
import {msalProvider} from 'boot/microsoft-msal';
import { useI18n } from 'vue-i18n';

const {t} = useI18n();

const router = useRouter();
const auth = useAuthStore();

const email = ref('');
const password = ref('');

const authProviders = ref(null as AuthProviders | null);

async function loadAuthProviders() {
  await BeeTimeClock.getAuthProviders().then((result) => {
    if (result.status === 200) {
      authProviders.value = result.data.Data;
    }
  }).catch((error: ErrorResponse) => {
    showErrorMessage(error.response?.data.Message);
  });
}

function loginLocal() {
  BeeTimeClock.login(email.value, password.value).then((result) => {
    if (result.status === 200) {
      auth.setAccessToken(result.data.Data.Token);
      auth.setAuthProvider('local');

      gotoDashboard();
    }
  }).catch((error) => {
    showErrorMessage(error)
    useAuthStore().logout();
  })
}

function loginWithMicrosoft() {
  if (!msalProvider) {
    showErrorMessage('hier ist was faul, meldet euch bei sebastian')
    return
  }

  msalProvider.msalInstance.loginPopup().then((result) => {
    try {

      auth.setAccessToken(result.idToken);
      auth.setAuthProvider('microsoft');


      const account = msalProvider.msalInstance.getAllAccounts()[0]
      if (account) {
        msalProvider.msalInstance.setActiveAccount(account)
      }
    } catch (e) {
      let message = 'Unknown error';
      if (e instanceof Error) {
        message = e.message;
      }
      showErrorMessage(`Cant set active account: ${message}`);
    } finally {
      gotoDashboard();
    }
  }).catch((error) => {
    showErrorMessage('Something went wrong with microsoft popup: ' + error)
    auth.logout();
  })
}

function gotoDashboard() {
  void router.push({name: 'Dashboard'})
}

onMounted(async () => {
  await loadAuthProviders();
})
</script>
