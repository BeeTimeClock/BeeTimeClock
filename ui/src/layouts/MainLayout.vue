<script lang="ts" setup>
import {computed, onMounted, ref, watch} from 'vue';
import {useAuthStore} from 'stores/microsoft-auth';
import {useI18n} from 'vue-i18n';
import type {User} from 'src/models/Authentication';
import BeeTimeClock from 'src/service/BeeTimeClock';
import type {BackendStatus} from 'src/models/Base';
import {type ErrorResponse} from 'src/models/Base';
import {showErrorMessage} from 'src/helper/message';
import {useRouter} from 'vue-router';
import {msalProvider} from 'boot/microsoft-msal';
import {useSettingsStore} from 'stores/settings';
import UserAvatar from 'components/UserAvatar.vue';

const {t} = useI18n();

const authStore = useAuthStore();
const settingsStore = useSettingsStore();
const session = ref(null as User | null);
const status = ref(null as BackendStatus | null);
const router = useRouter();
const leftDrawerOpen = ref(true);
const miniMode = ref(false);
const {locale} = useI18n({useScope: 'global'});
const isAdministrator = ref(false);
const missingDaysCount = ref(0);
const suspiciousCount = ref(0);

const localeOptions = [
  {value: 'en-US', label: 'English'},
  {value: 'de', label: 'Deutsch'},
];

const commit = computed(() => {
  return process.env.VUE_APP_COMMIT;
});

function logout() {
  authStore.logout();
  void router.push({name: 'Login'});
}

function toggleLeftDrawer() {
  miniMode.value = !miniMode.value;
}

function loadMissingDaysCount() {
  BeeTimeClock.getMissingDaysCount()
    .then((result) => {
      if (result.status === 200) {
        missingDaysCount.value = result.data.Data.Count;
      }
    })
    .catch((error) => {
      showErrorMessage(error);
    });
}

function loadSuspiciousCount() {
  BeeTimeClock.timestampQuerySuspiciousCount()
    .then((result) => {
      if (result.status === 200) {
        suspiciousCount.value = result.data.Data.Count;
      }
    })
    .catch((error) => {
      showErrorMessage(error);
    });
}

watch(locale, () => {
  settingsStore.setLanguage(locale.value);
})


onMounted(async () => {
  settingsStore.load();
  const lang = settingsStore.language
  if (lang != null) {
    locale.value = lang
  }

  await refresh();
  loadMissingDaysCount();
  loadSuspiciousCount();
});

async function refresh() {
  if (authStore.getAuthProvider === 'microsoft') {
    msalProvider.refresh();
  }
  const isLoggedIn = await authStore.loadSession();
  if (!isLoggedIn) {
    console.log('unauth');
    logout();
    return;
  }

  session.value = authStore.getSession();

  BeeTimeClock.getStatus()
    .then((result) => {
      if (result.status === 200) {
        status.value = result.data.Data;
      }
    })
    .catch((error: ErrorResponse) => {
      showErrorMessage(error.message);
    });
  isAdministrator.value = authStore.isAdministrator();
}
</script>
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
            <img src="logo.svg" alt="btc-logo"/>
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
        <q-btn-dropdown
          v-if="session"
          flat
          color="white"
          class="q-ml-sm"
        >
          <template v-slot:label>
            <UserAvatar :user="session" :size="24" class="q-mr-sm" />
            {{ session.FirstName }} {{ session.LastName }}
          </template>
          <q-list>
            <q-item-label header class="text-grey-7">
              {{ session.Username }}
            </q-item-label>
            <q-separator />
            <q-item clickable v-close-popup :to="{ name: 'UserSettings' }">
              <q-item-section avatar><q-icon name="manage_accounts" /></q-item-section>
              <q-item-section>{{ t('MENU_SETTINGS') }}</q-item-section>
            </q-item>
            <q-item clickable v-close-popup :to="{ name: 'UserApikeyOverview' }">
              <q-item-section avatar><q-icon name="key" /></q-item-section>
              <q-item-section>{{ t('MENU_APIKEY') }}</q-item-section>
            </q-item>
            <q-separator />
            <q-item clickable v-close-popup @click="logout">
              <q-item-section avatar><q-icon name="logout" color="negative" /></q-item-section>
              <q-item-section class="text-negative">{{ t('BTN_LOG_OUT') }}</q-item-section>
            </q-item>
          </q-list>
        </q-btn-dropdown>
      </q-toolbar>
    </q-header>

    <q-drawer v-model="leftDrawerOpen" show-if-above bordered :width="280" :mini="miniMode" :mini-width="57">
      <q-scroll-area style="height: calc(100% - 50px); border-right: 1px solid #ddd">
        <q-list>
          <q-item-label v-if="!miniMode" header>{{ t('LABEL_MENU') }}</q-item-label>
          <q-separator v-else class="q-mt-sm q-mb-xs" />

          <q-item clickable v-ripple :to="{ name: 'Dashboard' }">
            <q-item-section avatar><q-icon name="dashboard" /></q-item-section>
            <q-item-section>{{ t('MENU_DASHBOARD') }}</q-item-section>
            <q-tooltip v-if="miniMode" anchor="center right" self="center left">{{ t('MENU_DASHBOARD') }}</q-tooltip>
          </q-item>
          <q-item clickable v-ripple :to="{ name: 'WorktimeOverview' }">
            <q-item-section avatar><q-icon name="schedule" /></q-item-section>
            <q-item-section>{{ t('MENU_WORKTIME') }}</q-item-section>
            <q-tooltip v-if="miniMode" anchor="center right" self="center left">{{ t('MENU_WORKTIME') }}</q-tooltip>
          </q-item>
          <q-item clickable v-ripple :to="{ name: 'AbsenceOverview' }">
            <q-item-section avatar><q-icon name="event_busy" /></q-item-section>
            <q-item-section>{{ t('MENU_ABSENCE') }}</q-item-section>
            <q-tooltip v-if="miniMode" anchor="center right" self="center left">{{ t('MENU_ABSENCE') }}</q-tooltip>
          </q-item>
          <q-item clickable v-ripple :to="{ name: 'ExternalWorkOverview' }">
            <q-item-section avatar><q-icon name="work_outline" /></q-item-section>
            <q-item-section>{{ t('MENU_EXTERNAL_WORK') }}</q-item-section>
            <q-tooltip v-if="miniMode" anchor="center right" self="center left">{{ t('MENU_EXTERNAL_WORK') }}</q-tooltip>
          </q-item>
          <q-item clickable v-ripple :to="{ name: 'OvertimeOverview' }">
            <q-item-section avatar><q-icon name="timer" /></q-item-section>
            <q-item-section>{{ t('MENU_OVERTIME') }}</q-item-section>
            <q-tooltip v-if="miniMode" anchor="center right" self="center left">{{ t('MENU_OVERTIME') }}</q-tooltip>
          </q-item>

          <q-item-label v-if="!miniMode" header>{{ t('LABEL_MONITORING') }}</q-item-label>
          <q-separator v-else class="q-mt-sm q-mb-xs" />

          <q-item clickable v-ripple :to="{ name: 'SuspiciousTimestampsOverview' }">
            <q-item-section avatar>
              <q-icon name="warning_amber">
                <q-badge v-if="miniMode && suspiciousCount > 0" color="negative" floating :label="suspiciousCount" />
              </q-icon>
            </q-item-section>
            <q-item-section>{{ t('MENU_SUSPICIOUS_TIMESTAMPS') }}</q-item-section>
            <q-item-section v-if="!miniMode" side>
              <q-chip v-if="suspiciousCount > 0" :label="suspiciousCount" dense color="negative" />
            </q-item-section>
            <q-tooltip v-if="miniMode" anchor="center right" self="center left">{{ t('MENU_SUSPICIOUS_TIMESTAMPS') }}</q-tooltip>
          </q-item>
          <q-item clickable v-ripple :to="{ name: 'MissingDaysOverview' }">
            <q-item-section avatar>
              <q-icon name="event_available">
                <q-badge v-if="miniMode && missingDaysCount > 0" color="negative" floating :label="missingDaysCount" />
              </q-icon>
            </q-item-section>
            <q-item-section>{{ t('MENU_MISSING_DAYS') }}</q-item-section>
            <q-item-section v-if="!miniMode" side>
              <q-chip v-if="missingDaysCount > 0" :label="missingDaysCount" dense color="negative" />
            </q-item-section>
            <q-tooltip v-if="miniMode" anchor="center right" self="center left">{{ t('MENU_MISSING_DAYS') }}</q-tooltip>
          </q-item>

          <q-item-label v-if="!miniMode" header>{{ t('LABEL_TEAM', 2) }}</q-item-label>
          <q-separator v-else class="q-mt-sm q-mb-xs" />

          <q-item clickable v-ripple :to="{ name: 'TeamOverview' }">
            <q-item-section avatar><q-icon name="groups" /></q-item-section>
            <q-item-section>{{ t('MENU_OVERVIEW') }}</q-item-section>
            <q-tooltip v-if="miniMode" anchor="center right" self="center left">{{ t('LABEL_TEAM', 2) }}</q-tooltip>
          </q-item>

          <template v-if="isAdministrator">
            <q-item-label v-if="!miniMode" header>{{ t('LABEL_ADMINISTRATION') }}</q-item-label>
            <q-separator v-else class="q-mt-sm q-mb-xs" />

            <q-expansion-item v-if="!miniMode" icon="people" :label="t('LABEL_PEOPLE')" :content-inset-level="0.5">
              <q-list>
                <q-item clickable v-ripple :to="{ name: 'AdministrationUserOverview' }">
                  <q-item-section avatar><q-icon name="person" /></q-item-section>
                  <q-item-section>{{ t('MENU_USERS') }}</q-item-section>
                </q-item>
                <q-item clickable v-ripple :to="{ name: 'AdministrationTeamOverview' }">
                  <q-item-section avatar><q-icon name="groups" /></q-item-section>
                  <q-item-section>{{ t('MENU_TEAMS') }}</q-item-section>
                </q-item>
              </q-list>
            </q-expansion-item>
            <q-item v-else clickable v-ripple>
              <q-item-section avatar><q-icon name="people" /></q-item-section>
              <q-menu anchor="top end" self="top start" auto-close>
                <q-list style="min-width: 160px">
                  <q-item-label header class="text-grey-7 q-py-xs">{{ t('LABEL_PEOPLE') }}</q-item-label>
                  <q-separator />
                  <q-item clickable v-ripple :to="{ name: 'AdministrationUserOverview' }">
                    <q-item-section avatar><q-icon name="person" /></q-item-section>
                    <q-item-section>{{ t('MENU_USERS') }}</q-item-section>
                  </q-item>
                  <q-item clickable v-ripple :to="{ name: 'AdministrationTeamOverview' }">
                    <q-item-section avatar><q-icon name="groups" /></q-item-section>
                    <q-item-section>{{ t('MENU_TEAMS') }}</q-item-section>
                  </q-item>
                </q-list>
              </q-menu>
            </q-item>

            <q-item clickable v-ripple :to="{ name: 'AdministrationWorktimeModelOverview' }">
              <q-item-section avatar><q-icon name="access_time" /></q-item-section>
              <q-item-section>{{ t('MENU_WORKTIME_MODELS') }}</q-item-section>
              <q-tooltip v-if="miniMode" anchor="center right" self="center left">{{ t('MENU_WORKTIME_MODELS') }}</q-tooltip>
            </q-item>
            <q-item clickable v-ripple :to="{ name: 'AdministrationSettingsCommon' }">
              <q-item-section avatar><q-icon name="settings" /></q-item-section>
              <q-item-section>{{ t('MENU_COMMON') }}</q-item-section>
              <q-tooltip v-if="miniMode" anchor="center right" self="center left">{{ t('MENU_COMMON') }}</q-tooltip>
            </q-item>
            <q-item clickable v-ripple :to="{ name: 'AdministrationSettingsHolidays' }">
              <q-item-section avatar><q-icon name="celebration" /></q-item-section>
              <q-item-section>{{ t('MENU_HOLIDAYS') }}</q-item-section>
              <q-tooltip v-if="miniMode" anchor="center right" self="center left">{{ t('MENU_HOLIDAYS') }}</q-tooltip>
            </q-item>
            <q-item clickable v-ripple :to="{ name: 'AdministrationSettingsTimestamp' }">
              <q-item-section avatar><q-icon name="fingerprint" /></q-item-section>
              <q-item-section>{{ t('MENU_TIMESTAMP') }}</q-item-section>
              <q-tooltip v-if="miniMode" anchor="center right" self="center left">{{ t('MENU_TIMESTAMP') }}</q-tooltip>
            </q-item>
            <q-item clickable v-ripple :to="{ name: 'AdministrationSettingsAbsence' }">
              <q-item-section avatar><q-icon name="event_busy" /></q-item-section>
              <q-item-section>{{ t('MENU_ABSENCE') }}</q-item-section>
              <q-tooltip v-if="miniMode" anchor="center right" self="center left">{{ t('MENU_ABSENCE') }}</q-tooltip>
            </q-item>
            <q-item clickable v-ripple :to="{ name: 'AdministrationSettingsExternalWork' }">
              <q-item-section avatar><q-icon name="work_outline" /></q-item-section>
              <q-item-section>{{ t('MENU_EXTERNAL_WORK') }}</q-item-section>
              <q-tooltip v-if="miniMode" anchor="center right" self="center left">{{ t('MENU_EXTERNAL_WORK') }}</q-tooltip>
            </q-item>
            <q-item clickable v-ripple :to="{ name: 'AdministrationSettingsNotify' }">
              <q-item-section avatar><q-icon name="notifications" /></q-item-section>
              <q-item-section>{{ t('MENU_NOTIFY') }}</q-item-section>
              <q-tooltip v-if="miniMode" anchor="center right" self="center left">{{ t('MENU_NOTIFY') }}</q-tooltip>
            </q-item>
            <q-item clickable v-ripple :to="{ name: 'AdministrationSettingsTerminal' }">
              <q-item-section avatar><q-icon name="point_of_sale" /></q-item-section>
              <q-item-section>{{ t('MENU_TERMINAL') }}</q-item-section>
              <q-tooltip v-if="miniMode" anchor="center right" self="center left">{{ t('MENU_TERMINAL') }}</q-tooltip>
            </q-item>
          </template>
        </q-list>
      </q-scroll-area>
      <div v-if="!miniMode" class="absolute-bottom row items-center" style="height: 50px; border-top: 1px solid #ddd">
        <q-item-label header>
          <div>Version: {{ commit }}</div>
        </q-item-label>
      </div>
    </q-drawer>

    <q-page-container>
      <router-view/>
    </q-page-container>
  </q-layout>
</template>
