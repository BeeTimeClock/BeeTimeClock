<template>
  <q-page padding>
    <div>
      <q-btn color="positive" :label="$t('BTN_CHECK_IN')" @click="actionCheckInOffice"/>
      <q-btn class="q-ml-lg" color="positive" :label="$t('BTN_CHECK_IN') + ' (Homeoffice)'" @click="actionCheckInHomeoffice"/>
      <q-btn class="q-ml-lg" color="negative" :label="$t('BTN_CHECK_OUT')" @click="actionCheckOut"/>
    </div>
    <div class="q-pt-lg">
      <q-table
        flat bordered
        :rows="timestampCurrentMonthGrouped"
        :columns="columns"
        row-key="Date"
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
              {{ col.value }}
            </q-td>
          </q-tr>
          <q-tr v-show="props.expand" :props="props">
            <q-td colspan="100%">
              <q-markup-table>
                <thead>
                <q-tr>
                  <q-th class="text-left">{{ $t('LABEL_COMING_TIMESTAMP') }}</q-th>
                  <q-th class="text-left">{{ $t('LABEL_GOING_TIMESTAMP') }}</q-th>
                </q-tr>
                </thead>
                <tbody>
                <tr v-for="timestamp in props.row.Timestamps" :key="timestamp.ID">
                  <td>{{ formatDateTemplate(timestamp.ComingTimestamp, 'HH:mm') }}</td>
                  <td>{{ formatDateTemplate(timestamp.GoingTimestamp, 'HH:mm') }}</td>
                </tr>
                </tbody>
              </q-markup-table>
            </q-td>
          </q-tr>
        </template>

      </q-table>
    </div>
  </q-page>
</template>

<script lang="ts">
import {defineComponent, ref} from 'vue'
import BeeTimeClock from 'src/service/BeeTimeClock';
import {TimestampGroup} from 'src/models/Timestamp';
import {date} from 'quasar';
import formatDate = date.formatDate;
import {showInfoMessage} from 'src/helper/message';

export default defineComponent({
  computed: {
    date() {
      return date
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
        format: (val: Date) => `${formatDate(val, 'DD.MM.YYYY')}`,
        sortable: true
      },
      {name: 'WorkingHours', align: 'left', label: this.$t('LABEL_WORKING_HOURS'), field: 'WorkingHours', format: (val: number) => val.toFixed(2)},
      {name: 'SubtractedHours', align: 'left', label: this.$t('LABEL_SUBTRACTED_HOURS'), field: 'SubtractedHours', format: (val: number) => val.toFixed(2)},
      {name: 'IsHomeoffice', align: 'left', label: this.$t('LABEL_HOMEOFFICE'), field: 'IsHomeoffice'},
    ]

    return {
      columns: columns,
      timestampCurrentMonthGrouped: ref([] as TimestampGroup[]),
    }
  },
  methods: {
    formatDate,
    formatDateTemplate(date: Date, format: string) : string {
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
        if (result.status === 200) {
          this.loadTimestampCurrentMonthGrouped();
          showInfoMessage(this.$t('CHECK_IN_SUCCESS'));
        }
      })
    },
    actionCheckOut() {
      BeeTimeClock.timestampActionCheckout().then((result) => {
        if (result.status === 200) {
          this.loadTimestampCurrentMonthGrouped();
        }
      })
    },
    loadTimestampCurrentMonthGrouped() {
      BeeTimeClock.timestampQueryCurrentMonthGrouped().then((result) => {
        if (result.status === 200) {
          this.timestampCurrentMonthGrouped = result.data.Data as TimestampGroup[];
        }
      });
    }
  },
  mounted() {
    this.loadTimestampCurrentMonthGrouped();
  }
})
</script>
