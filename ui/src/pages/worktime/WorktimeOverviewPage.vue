<template>
  <q-page padding>
    <div class="row">
      <div class="col-6 q-pa-sm">
        <OvertimeCurrentMonth class="full-height"/>
      </div>
      <div class="col-6 q-pa-sm">
        <OvertimeTotal class="full-height"/>
      </div>
    </div>
    <div class="row">
      <q-card flat>
        <q-card-section>
          <q-card-actions>
            <q-btn class="q-ma-sm" color="positive" :label="$t('BTN_CHECK_IN')"
                   @click="actionCheckInOffice"/>
            <q-btn class="q-ma-sm" color="positive" :label="$t('BTN_CHECK_IN') + ' (Homeoffice)'"
                   @click="actionCheckInHomeoffice"/>
            <q-btn class="q-ma-sm" color="negative" :label="$t('BTN_CHECK_OUT')" @click="actionCheckOut"/>
            <q-btn class="q-ma-sm" color="primary" :label="$t('BTN_TIMESTAMP_ADD')"
                   @click="promptTimestampCorrectionCreate(null)"/>
          </q-card-actions>
        </q-card-section>
      </q-card>
    </div>
    <div class="q-pt-lg">
      <WorktimeOverviewTable v-model="timestampCurrentMonthGrouped"/>
    </div>
  </q-page>
  <q-dialog v-model="prompt.timestampCorrectionCreate" persistent>
    <q-card class="q-dialog-plugin full-width">
      <q-card-section>
        <div class="text-h6">{{ $t('LABEL_TIMESTAMP_CORRECTION_CREATE') }}</div>
      </q-card-section>
      <q-form @submit.prevent.stop="timestampCorrectionCreate">
        <q-card-section>
          <DateTimePickerComponent v-model="timestampCorrection.ComingTimestamp" :label="$t('LABEL_COMING_TIMESTAMP')"/>
          <DateTimePickerComponent class="q-mt-md" v-model="timestampCorrection.GoingTimestamp"
                                   :label="$t('LABEL_GOING_TIMESTAMP')"/>
          <q-checkbox v-model="timestampCorrection.IsHomeoffice" :label="$t('LABEL_HOMEOFFICE')"/>
          <q-input lazy-rules :rules="[val => !!val || $t('LABEL_FIELD_REQUIRED')]" v-model="timestampCorrection.Reason"
                   type="textarea" :label="$t('LABEL_REASON')"/>
        </q-card-section>
        <q-card-actions>
          <q-btn v-close-popup :label="$t('BTN_CANCEL')" color="negative" type="reset"/>
          <q-btn :label="$t('BTN_CREATE')" color="positive" type="submit"/>
        </q-card-actions>
      </q-form>
    </q-card>

  </q-dialog>
</template>

<script lang="ts">
import {date} from 'quasar';
import DateTimePickerComponent from 'src/components/DateTimePickerComponent.vue';
import {showErrorMessage, showInfoMessage} from 'src/helper/message';
import {Timestamp, TimestampCorrectionCreateRequest, TimestampGroup} from 'src/models/Timestamp';
import BeeTimeClock from 'src/service/BeeTimeClock';
import {defineComponent, ref} from 'vue';
import {TimestampCreateRequest} from 'src/models/Timestamp';
import {ErrorResponse} from 'src/models/Base';
import OvertimeCurrentMonth from 'components/OvertimeCurrentMonth.vue';
import OvertimeTotal from 'components/OvertimeTotal.vue';
import WorktimeOverviewTable from 'components/WorktimeOverviewTable.vue';
import formatDate = date.formatDate;

export default defineComponent({
  computed: {
    date() {
      return date;
    }
  },
  data() {
    return {
      timestampCurrentMonthGrouped: ref([] as TimestampGroup[]),
      prompt: {
        timestampCorrectionCreate: false,
      },
      timestampCorrection: {
        ComingTimestamp: ref(''),
        GoingTimestamp: ref(null as string | null),
        IsHomeoffice: false,
        Reason: ref(''),
      },
      selectedTimestamp: null as Timestamp | null,
      expanded: ref(['']),
    };
  },
  methods: {
    formatDate,
    promptTimestampCorrectionCreate(timestamp: Timestamp | null) {
      this.timestampCorrection = {
        ComingTimestamp: formatDate(new Date(), 'DD.MM.YYYY HH:mm'),
        GoingTimestamp: null,
        IsHomeoffice: false,
        Reason: '',
      }

      if (timestamp != null) {
        this.timestampCorrection.ComingTimestamp = formatDate(new Date(timestamp.ComingTimestamp), 'DD.MM.YYYY HH:mm');

        if (timestamp.GoingTimestamp != null && new Date(timestamp.GoingTimestamp).getFullYear() != 1) {
          this.timestampCorrection.GoingTimestamp = formatDate(new Date(timestamp.GoingTimestamp), 'DD.MM.YYYY HH:mm');
        }
      }

      this.selectedTimestamp = timestamp;
      this.prompt.timestampCorrectionCreate = true;
    },
    actionCheckInHomeoffice() {
      this.actionCheckIn(true);
    },
    actionCheckInOffice() {
      this.actionCheckIn(false);
    },
    actionCheckIn(isHomeoffice = false) {
      BeeTimeClock.timestampActionCheckin(isHomeoffice).then((result) => {
        if (result.status === 201) {
          showInfoMessage(this.$t('MSG_CHECK_IN_SUCCESS'));
          this.loadTimestampCurrentMonthGrouped();
        }
      }).catch((error: ErrorResponse) => {
        showErrorMessage(error.response?.data.Message);
      });
    },
    actionCheckOut() {
      BeeTimeClock.timestampActionCheckout().then((result) => {
        if (result.status === 200) {
          showInfoMessage(this.$t('MSG_CHECK_OUT_SUCCESS'));
          this.loadTimestampCurrentMonthGrouped();
        }
      }).catch((error: ErrorResponse) => {
        showErrorMessage(error.response?.data.Message);
      });
    },
    loadTimestampCurrentMonthGrouped() {
      BeeTimeClock.timestampQueryCurrentMonthGrouped().then((result) => {
        if (result.status === 200) {
          this.timestampCurrentMonthGrouped = result.data.Data.sort((a, b) => new Date(b.Date).getTime() - new Date(a.Date).getTime());
          if (this.timestampCurrentMonthGrouped.length > 0) {
            this.expanded = [this.timestampCurrentMonthGrouped[0].Date.toString()];
          }
        }
      });
    },
    timestampCorrectionCreate() {
      const comingTimestamp = date.extractDate(this.timestampCorrection.ComingTimestamp, 'DD.MM.YYYY HH:mm');

      let goingTimestamp = null;
      if (this.timestampCorrection.GoingTimestamp) {
        goingTimestamp = date.extractDate(this.timestampCorrection.GoingTimestamp, 'DD.MM.YYYY HH:mm');
      }

      this.prompt.timestampCorrectionCreate = false;

      if (this.selectedTimestamp != null) {
        const timestampCorrectionRequest = {
          NewComingTimestamp: comingTimestamp,
          NewGoingTimestamp: goingTimestamp,
          ChangeReason: this.timestampCorrection.Reason,
          IsHomeoffice: this.timestampCorrection.IsHomeoffice,
        } as TimestampCorrectionCreateRequest;

        BeeTimeClock.timestampCorrectionCreate(this.selectedTimestamp, timestampCorrectionRequest).then((result) => {
          if (result.status === 200) {
            showInfoMessage(this.$t('MSG_CREATE_SUCCESS'));
            this.loadTimestampCurrentMonthGrouped();
          }
        })
      } else {
        const timestampCreateRequest = {
          ComingTimestamp: comingTimestamp,
          GoingTimestamp: goingTimestamp,
          ChangeReason: this.timestampCorrection.Reason,
          IsHomeoffice: this.timestampCorrection.IsHomeoffice,
        } as TimestampCreateRequest;

        BeeTimeClock.timestampCreate(timestampCreateRequest).then((result) => {
          if (result.status === 201) {
            showInfoMessage(this.$t('MSG_CREATE_SUCCESS'));
            this.loadTimestampCurrentMonthGrouped();
          }
        });
      }
    }
  },
  mounted() {
    this.loadTimestampCurrentMonthGrouped();
  },
  components: { WorktimeOverviewTable, OvertimeTotal, OvertimeCurrentMonth, DateTimePickerComponent}
})
</script>
