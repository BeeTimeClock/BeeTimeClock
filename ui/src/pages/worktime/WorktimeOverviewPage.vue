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
      <q-table
        flat bordered
        :rows="timestampCurrentMonthGrouped"
        :columns="columns"
        row-key="Date"
        :pagination="pagination"
        hide-pagination
        v-model:expanded="expanded"
      >
        <template v-slot:header="props">
          <q-tr :props="props">
            <q-th auto-width/>
            <q-th
              v-for="col in props.cols"
              :key="col.name"
              :props="props"
            >
              {{ col.label }}
            </q-th>
          </q-tr>
        </template>
        <template v-slot:body="props">
          <q-tr :props="props">
            <q-td auto-width>
              <q-btn size="sm" color="accent" round dense @click="props.expand = !props.expand"
                     :icon="props.expand ? 'remove' : 'add'"/>
            </q-td>
            <q-td
              v-for="col in props.cols"
              :key="col.name"
              :props="props"
            >
              <div v-if="col.name == 'IsHomeoffice'">
                <q-icon size="large" v-if="col.value" name="done"/>
              </div>
              <div v-else>
                {{ col.value }}
              </div>
            </q-td>
          </q-tr>
          <q-tr v-show="props.expand" :props="props">
            <q-td colspan="100%">
              <q-markup-table>
                <thead>
                <q-tr>
                  <q-th class="text-left">{{ $t('LABEL_COMING_TIMESTAMP') }}</q-th>
                  <q-th class="text-left">{{ $t('LABEL_GOING_TIMESTAMP') }}</q-th>
                  <q-th></q-th>
                </q-tr>
                </thead>
                <tbody>
                <tr v-for="timestamp in props.row.Timestamps" :key="timestamp.ID">
                  <td>{{ formatDateTemplate(timestamp.ComingTimestamp, 'HH:mm') }}</td>
                  <td>
                    <div v-if="timestamp.GoingTimestamp && new Date(timestamp.GoingTimestamp).getFullYear() != 1">{{
                        formatDateTemplate(timestamp.GoingTimestamp, 'HH:mm')
                      }}
                    </div>
                  </td>
                  <td>
                    <q-btn color="primary" class="q-mr-md" :disable="timestamp.Corrections.length == 0"
                           icon="pending_actions" @click="promptTimestampCorrectionView(timestamp)"/>
                    <q-btn color="primary" icon="edit" @click="promptTimestampCorrectionCreate(timestamp)"/>
                  </td>
                </tr>
                </tbody>
              </q-markup-table>
            </q-td>
          </q-tr>
        </template>

      </q-table>
    </div>
  </q-page>
  <q-dialog v-model="prompt.timestampCorrectionCreate" persistent>
    <q-card class="q-dialog-plugin full-width">
      <q-card-section>
        <div class="text-h6">{{ $t('LABEL_TIMESTAMP_CORRECTION_CREATE') }}</div>
      </q-card-section>
      <q-card-section>
        <DateTimePickerComponent v-model="timestampCorrection.ComingTimestamp" :label="$t('LABEL_COMING_TIMESTAMP')"/>
        <DateTimePickerComponent class="q-mt-md" v-model="timestampCorrection.GoingTimestamp"
                                 :label="$t('LABEL_GOING_TIMESTAMP')"/>
        <q-checkbox v-model="timestampCorrection.IsHomeoffice" :label="$t('LABEL_HOMEOFFICE')"/>
        <q-input v-model="timestampCorrection.Reason" type="textarea" :label="$t('LABEL_REASON')"/>
      </q-card-section>
      <q-card-actions>
        <q-btn v-close-popup :label="$t('BTN_CANCEL')" color="negative"/>
        <q-btn v-close-popup :label="$t('BTN_CREATE')" color="positive" @click="timestampCorrectionCreate"/>
      </q-card-actions>
    </q-card>
  </q-dialog>
  <q-dialog v-model="prompt.timestampCorrectionView" persistent>
    <q-card class="q-dialog-plugin full-width">
      <q-card-section>
        <div class="text-h6">{{ $t('LABEL_TIMESTAMP_CORRECTION_VIEW') }}</div>
      </q-card-section>
      <q-card-section>
        <q-markup-table flat>
          <thead>
          <tr>
            <th>{{ $t('LABEL_REASON') }}</th>
            <th>{{ $t('LABEL_OLD_COMING') }}</th>
            <th>{{ $t('LABEL_OLD_GOING') }}</th>
            <th>{{ $t('LABEL_CREATED_AT') }}</th>
          </tr>
          </thead>
          <tbody>
          <tr v-for="correction in selectedTimestamp?.Corrections" :key="correction.ID">
            <td>{{ correction.ChangeReason }}</td>
            <td>{{ date.formatDate(correction.OldComingTimestamp, 'DD.MM.YYYY HH:mm') }}</td>
            <td>{{ date.formatDate(correction.OldGoingTimestamp, 'DD.MM.YYYY HH:mm') }}</td>
            <td>{{ date.formatDate(correction.CreatedAt, 'DD.MM.YYYY HH:mm') }}</td>
          </tr>
          </tbody>
        </q-markup-table>
      </q-card-section>
      <q-card-actions>
        <q-btn v-close-popup :label="$t('BTN_CLOSE')" color="primary"/>
      </q-card-actions>
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
import formatDate = date.formatDate;
import {TimestampCreateRequest} from 'src/models/Timestamp';
import {ErrorResponse} from 'src/models/Base';
import OvertimeCurrentMonth from "components/OvertimeCurrentMonth.vue";
import OvertimeTotal from "components/OvertimeTotal.vue";

export default defineComponent({
  computed: {
    date() {
      return date;
    }
  },
  data() {
    const columns = [
      {
        name: 'Date',
        required: true,
        label: this.$t('LABEL_DATE'),
        align: 'left',
        field: 'Date',
        format: (val: Date) => `${formatDate(val, 'ddd. DD.MM.YYYY')}`,
        sortable: true
      },
      {
        name: 'WorkingHours',
        align: 'left',
        label: this.$t('LABEL_WORKING_HOURS'),
        field: 'WorkingHours',
        format: (val: number) => val.toFixed(2)
      },
      {
        name: 'SubtractedHours',
        align: 'left',
        label: this.$t('LABEL_SUBTRACTED_HOURS'),
        field: 'SubtractedHours',
        format: (val: number) => val.toFixed(2)
      },
      {
        name: 'OvertimeHours',
        align: 'left',
        label: this.$t('LABEL_OVERTIME_HOURS'),
        field: 'OvertimeHours',
        format: (val: number) => val.toFixed(2)
      },
      {
        name: 'IsHomeoffice',
        align: 'left',
        label: this.$t('LABEL_HOMEOFFICE'),
        field: 'IsHomeoffice',
      },
    ];
    return {
      pagination: {
        page: 1,
        rowsPerPage: 0,
      },
      columns: columns,
      timestampCurrentMonthGrouped: ref([] as TimestampGroup[]),
      prompt: {
        timestampCorrectionCreate: false,
        timestampCorrectionView: false,
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
    promptTimestampCorrectionView(timestamp: Timestamp) {
      this.selectedTimestamp = timestamp;
      this.prompt.timestampCorrectionView = true;
    },
    formatDateTemplate(date: Date, format: string): string {
      return formatDate(date, format);
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
  components: {OvertimeTotal, OvertimeCurrentMonth, DateTimePickerComponent}
})
</script>
