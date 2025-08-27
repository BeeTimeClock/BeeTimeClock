<template>
  <q-layout view="lHh Lpr lFf">
    <q-header elevated>
      <q-toolbar>
        <q-btn
          flat
          dense
          round
          icon="menu"
          aria-label="Menu"
          @click="toggleLeftDrawer"
        />

        <q-toolbar-title>
          <q-avatar square>
            <img src="logo.svg" alt="btc-logo" />
          </q-avatar>
          Bee Time Clock
        </q-toolbar-title>

        <q-select
          v-model="locale"
          :options="localeOptions"
          dense
          borderless
          emit-value
          map-options
          options-dense
        />
        <q-btn
          class="q-ml-lg"
          color="secondary"
          :label="$t('BTN_LOG_OUT')"
          @click="logout"
        />
      </q-toolbar>
    </q-header>

    <q-drawer v-model="leftDrawerOpen" show-if-above bordered>
      <q-scroll-area
        style="
          height: calc(100% - 80px);
          margin-top: 80px;
          border-right: 1px solid #ddd;
        "
      >
        <q-list>
          <q-item-label header>
            {{ $t('LABEL_MENU') }}
          </q-item-label>
          <q-item clickable v-ripple :to="{ name: 'Dashboard' }">
            {{ $t('MENU_DASHBOARD') }}
          </q-item>
          <q-item clickable v-ripple :to="{ name: 'WorktimeOverview' }">
            {{ $t('MENU_WORKTIME') }}
          </q-item>
          <q-item clickable v-ripple :to="{ name: 'AbsenceOverview' }">
            {{ $t('MENU_ABSENCE') }}
          </q-item>
          <q-item clickable v-ripple :to="{ name: 'ExternalWorkOverview' }">
            {{ $t('MENU_EXTERNAL_WORK') }}
          </q-item>
          <q-item clickable v-ripple :to="{ name: 'OvertimeOverview' }">
            {{ $t('MENU_OVERTIME') }}
          </q-item>
          <q-item
            clickable
            v-ripple
            :to="{ name: 'SuspiciousTimestampsOverview' }"
          >
            {{ $t('MENU_SUSPICIOUS_TIMESTAMPS') }}
          </q-item>
          <q-item-label header> {{ $t('LABEL_ME') }}</q-item-label>
          <q-item clickable v-ripple :to="{ name: 'UserApikeyOverview' }">
            {{ $t('MENU_APIKEY') }}
          </q-item>
          <q-item clickable v-ripple :to="{ name: 'UserSettings' }">
            {{ $t('MENU_SETTINGS') }}
          </q-item>
          <div v-if="isAdministrator">
            <q-item-label header>
              {{ $t('LABEL_ADMINISTRATION') }}
            </q-item-label>
            <q-item
              clickable
              v-ripple
              :to="{ name: 'AdministrationUserOverview' }"
            >
              {{ $t('MENU_USERS') }}
            </q-item>
            <q-item
              clickable
              v-ripple
              :to="{ name: 'AdministrationTeamOverview' }"
            >
              {{ $t('MENU_TEAMS') }}
            </q-item>
            <q-expansion-item
              content-inset-level="0.5"
              :label="$t('MENU_SETTINGS')"
            >
              <q-list>
                <q-item
                  clickable
                  v-ripple
                  :to="{ name: 'AdministrationSettingsCommon' }"
                >
                  {{ $t('MENU_COMMON') }}
                </q-item>
                <q-item
                  clickable
                  v-ripple
                  :to="{ name: 'AdministrationSettingsTimestamp' }"
                >
                  {{ $t('MENU_TIMESTAMP') }}
                </q-item>
                <q-item
                  clickable
                  v-ripple
                  :to="{ name: 'AdministrationSettingsAbsence' }"
                >
                  {{ $t('MENU_ABSENCE') }}
                </q-item>
                <q-item
                  clickable
                  v-ripple
                  :to="{ name: 'AdministrationSettingsNotify' }"
                >
                  {{ $t('MENU_NOTIFY') }}
                </q-item>
                <q-item
                  clickable
                  v-ripple
                  :to="{ name: 'AdministrationSettingsExternalWork' }"
                >
                  {{ $t('MENU_EXTERNAL_WORK') }}
                </q-item>
              </q-list>
            </q-expansion-item>
          </div>
        </q-list>
      </q-scroll-area>
      <div
        class="absolute-top bg-primary q-pa-md text-white"
        style="height: 80px"
      >
        <div v-if="session">
          <div class="text-weight-bold">
            {{ session.FirstName }} {{ session.LastName }}
          </div>
          <div>{{ session.Username }}</div>
        </div>
      </div>
      <div class="absolute-bottom">
        <q-list>
          <q-item-label header>
            <div>Version: {{ commit }}<br /></div>
          </q-item-label>
        </q-list>
      </div>
    </q-drawer>

    <q-page-container>
      <router-view />
    </q-page-container>
  </q-layout>
</template>

<script lang="ts">
import { defineComponent, ref } from 'vue';
import { useAuthStore } from 'stores/microsoft-auth';
import { useI18n } from 'vue-i18n';
import type { User } from 'src/models/Authentication';
import BeeTimeClock from 'src/service/BeeTimeClock';
import type { BackendStatus} from 'src/models/Base';
import { type ErrorResponse } from 'src/models/Base';
import { showErrorMessage } from 'src/helper/message';

export default defineComponent({
  name: 'MainLayout',
  data() {
    return {
      session: ref(null as User | null),
      status: ref(null as BackendStatus | null),
    };
  },
  methods: {
    useAuthStore,
    logout() {
      useAuthStore().logout();
      this.$router.push({ name: 'Login' });
    },
    async refresh() {
      if (useAuthStore().getAuthProvider === 'microsoft') {
        this.$msalProvider.refresh();
      }
      const isLoggedIn = await useAuthStore().loadSession();
      if (!isLoggedIn) {
        console.log('unauth');
        this.logout();
        return;
      }

      this.session = useAuthStore().getSession();
      BeeTimeClock.getStatus()
        .then((result) => {
          if (result.status === 200) {
            this.status = result.data.Data;
          }
        })
        .catch((error: ErrorResponse) => {
          showErrorMessage(error.message);
        });
      this.isAdministrator = useAuthStore().isAdministrator();
    },
  },
  computed: {
    commit() {
      return process.env.VUE_APP_COMMIT;
    },
  },
  async mounted() {
    await this.refresh();
  },
  setup() {
    const leftDrawerOpen = ref(false);
    const { locale } = useI18n({ useScope: 'global' });
    const isAdministrator = ref(false);

    return {
      leftDrawerOpen,
      locale,
      localeOptions: [
        { value: 'en-US', label: 'English' },
        { value: 'de', label: 'Deutsch' },
      ],
      toggleLeftDrawer() {
        leftDrawerOpen.value = !leftDrawerOpen.value;
      },
      isAdministrator,
    };
  },
});
</script>
