<script setup lang="ts">
import {OfficeIPAddress, Settings} from 'src/models/Settings';
import {onMounted, ref} from 'vue';
import BeeTimeClock from 'src/service/BeeTimeClock';
import {useI18n} from 'vue-i18n';
import AbsenceReasonAdministrationTable from 'components/AbsenceReasonAdministrationTable.vue';
import { showInfoMessage } from 'src/helper/message';

const {t} = useI18n();

const applicationSettings = ref<Settings | undefined>(undefined);
const newIpAddress = ref<OfficeIPAddress | null>();
const promptOfficeIpAdd = ref(false);
const tab = ref('timestamp')

function loadSettings() {
  BeeTimeClock.administrationSettings().then(result => {
    if (result.status === 200) {
      applicationSettings.value = result.data.Data;
    }
  })
}

function updateSettings() {
  if (!applicationSettings.value) return;

  BeeTimeClock.administrationSettingsSave(applicationSettings.value).then(result => {
    if (result.status === 200) {
      applicationSettings.value = result.data.Data;
    }
  })
}

function addNewIpAddress() {
  if (!newIpAddress.value) return;
  if (!applicationSettings.value) return;

  applicationSettings.value.OfficeIPAddresses.push(newIpAddress.value);
}

function notifyCurrentWeek() {
  BeeTimeClock.administrationNotifyAbsenceWeek().then(result => {
    if (result.status === 204) {
      showInfoMessage('notify sent')
    }
  })
}

onMounted(() => {
  loadSettings();
})
</script>

<template>
  <q-page>
    <div v-if="applicationSettings">
      <q-splitter
        v-model="splitterModel"
      >
        <template v-slot:before>
          <q-tabs
            v-model="tab"
            vertical
            class="text-teal"
          >
            <q-tab name="timestamp" icon="watch" :label="$t('LABEL_TIMESTAMP')" />
            <q-tab name="absence" icon="event" :label="$t('LABEL_ABSENCE')" />
            <q-tab name="notify" icon="alert" :label="$t('LABEL_NOTIFY')" />
          </q-tabs>
        </template>

        <template v-slot:after>
          <q-tab-panels
            v-model="tab"
            vertical
          >
            <q-tab-panel name="timestamp">
              <q-toggle :label="t('LABEL_CHECKIN_DETECTION_BY_IP_ADDRESS')"
                        v-model="applicationSettings.CheckinDetectionByIPAddress"/>
              <q-card class="q-mt-lg">
                <q-card-section class="bg-primary text-white text-h6">
                  Office IPs
                </q-card-section>
                <q-card-section>
                  <q-list>
                    <q-item v-for="officeIp in applicationSettings.OfficeIPAddresses" :key="officeIp.ID">
                      <q-item-section>
                        <q-item-label>{{ officeIp.IPAddress }}</q-item-label>
                        <q-item-label caption>{{ officeIp.Description }}</q-item-label>
                      </q-item-section>
                    </q-item>
                    <q-item>
                      <q-btn class="full-width" :label="t('BTN_ADD')" color="primary"
                             @click="newIpAddress = {} as OfficeIPAddress; promptOfficeIpAdd = true"/>
                    </q-item>
                  </q-list>
                </q-card-section>
              </q-card>
              <q-btn class="full-width q-mt-lg" color="primary" :label="t('BTN_SAVE')" @click="updateSettings"/>
            </q-tab-panel>

            <q-tab-panel name="absence">
              <q-card>
                <q-card-section class="bg-primary text-white text-h6">
                  {{$t('LABEL_ABSENCE_REASON', 2 )}}
                </q-card-section>
                <q-card-section>
                  <AbsenceReasonAdministrationTable/>
                </q-card-section>
              </q-card>
            </q-tab-panel>
            <q-tab-panel name="notify">
              <q-btn label="Notify Absence Current Week" @click="notifyCurrentWeek"/>
            </q-tab-panel>
          </q-tab-panels>
        </template>

      </q-splitter>
    </div>
    <q-dialog v-model="promptOfficeIpAdd">
      <q-card v-if="newIpAddress">
        <q-card-section>
          <q-input :label="t('LABEL_IP_ADDRESS')" v-model="newIpAddress.IPAddress"/>
          <q-input :label="t('LABEL_DESCRIPTION')" v-model="newIpAddress.Description"/>
        </q-card-section>
        <q-card-actions>
          <q-btn :label="t('BTN_CANCEL')" v-close-popup type="reset" color="negative"/>
          <q-btn :label="t('BTN_ADD')" v-close-popup @click="addNewIpAddress" color="positive"/>
        </q-card-actions>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<style scoped>

</style>
