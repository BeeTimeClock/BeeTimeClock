<template>
  <q-page>
    <q-page class="flex flex-center bg-grey-2">
      <q-card class="q-pa-md shadow-2 my_card" bordered>
        <q-card-section class="text-center">
          <img width="256" :src="'logo.png'" alt="btc-logo"/>
          <div class="text-grey-9 text-h3 text-weight-bold">Bee Time Clock</div>
          <div class="text-grey-9 text-h5 text-weight-bold">{{ $t('LABEL_SIGN_IN') }}</div>
        </q-card-section>
        <q-card-section>
          <q-input dense outlined v-model="email" :label="$t('LABEL_USERNAME')"></q-input>
          <q-input dense outlined class="q-mt-md" v-model="password" type="password"
                   :label="$t('LABEL_PASSWORD')"></q-input>
        </q-card-section>
        <q-card-section>
          <q-btn color="primary" size="md" :label="$t('LABEL_SIGN_IN')" no-caps class="full-width q-mb-lg"
                 @click="loginLocal">
          </q-btn>
          <q-btn v-if="authProviders?.Microsoft" color="primary" size="md" :label="$t('LABEL_SIGN_IN_MICROSOFT')"
                 no-caps class="full-width"
                 @click="loginWithMicrosoft">
            <q-avatar square class="q-ml-md" size="24px">
              <img :src="'ms_login.png'" alt="ms-logo"/>
            </q-avatar>
          </q-btn>
        </q-card-section>
      </q-card>
    </q-page>
  </q-page>
</template>

<script lang="ts">
import {AuthProviders} from 'src/models/Base';
import BeeTimeClock from 'src/service/BeeTimeClock';
import {useAuthStore} from 'stores/microsoft-auth';
import {defineComponent, ref} from 'vue';
import {showErrorMessage} from 'src/helper/message';

export default defineComponent({
  // name: 'PageName'
  data() {
    return {
      email: '',
      password: '',
      authProviders: ref(null as AuthProviders | null),
    }
  },
  methods: {
    async loadAuthProviders() {
      BeeTimeClock.getAuthProviders().then((result) => {
        if (result.status === 200) {
          this.authProviders = result.data.Data;
        }
      });
    },
    async loginLocal() {
      BeeTimeClock.login(this.email, this.password).then((result) => {
        if (result.status === 200) {
          useAuthStore().setAccessToken(result.data.Data.Token);
          useAuthStore().setAuthProvider('local');
          this.$router.push({name: 'Dashboard'});
        }
      }).catch((error) => {
        showErrorMessage(error)
        useAuthStore().logout();
      })
    },
    async loginWithMicrosoft() {
      await this.$msalProvider.msalInstance.loginPopup().then((result) => {
        useAuthStore().setAccessToken(result.idToken);
        useAuthStore().setAuthProvider('microsoft');
        this.$msalProvider.msalInstance.setActiveAccount(this.$msalProvider.msalInstance.getAllAccounts()[0])
        this.$router.push({name: 'Dashboard'})
      }).catch((error) => {
        showErrorMessage(error)
        useAuthStore().logout();
      })
    }
  },
  mounted() {
    this.loadAuthProviders();
  }
})
</script>
