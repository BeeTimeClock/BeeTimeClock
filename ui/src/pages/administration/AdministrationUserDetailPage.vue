<script setup lang="ts">

import { User } from 'src/models/Authentication';
import { computed, onMounted, ref, watch } from 'vue';
import BeeTimeClock from 'src/service/BeeTimeClock';
import { useRoute } from 'vue-router';
import { useI18n } from 'vue-i18n';
import { showInfoMessage } from 'src/helper/message';
import WorktimeOverviewTable from 'components/WorktimeOverviewTable.vue';
import OvertimeMonth from 'components/OvertimeMonth.vue';
import TimestampCorrectionDialog from 'components/TimestampCorrectionDialog.vue';
import { TimestampGroup, TimestampYearMonthGrouped } from 'src/models/Timestamp';
import AbsenceSummaryTableComponent from 'components/AbsenceSummaryTableComponent.vue';
import { Absence, AbsenceReason, AbsenceSummaryItem } from 'src/models/Absence';
import { date } from 'quasar';

const { t } = useI18n();

const route = useRoute();
const userId = route.params.userId as string;
const user = ref(null as User | null);
const selectedTab = ref('common');
const timestampYearMonths = ref<TimestampYearMonthGrouped>({});
const timestampCurrentMonthGrouped = ref<TimestampGroup[]>([]);
const expanded = ref(['']);
const selectedYear = ref<number>(new Date().getFullYear());
const selectedMonth = ref<number>(new Date().getMonth() + 1);
const absenceYears = ref<number[]>([]);
const absences = ref<Absence[]>([]);
const selectedAbsenceYear = ref<number>(new Date().getFullYear());
const absenceReasons = ref<AbsenceReason[]>([]);

const accessLevelOptions = [
  {
    value: 'admin',
    label: t('LABEL_ADMINISTRATOR')
  },
  {
    value: 'user',
    label: t('LABEL_USER')
  }
];

const overtimeSubtractions = [
  {
    value: 'percentage',
    label: t('LABEL_PERCENTAGE')
  },
  {
    value: 'hours',
    label: t('LABEL_HOURS')
  }
];

function loadUser() {
  BeeTimeClock.administrationGetUserById(userId).then(result => {
    if (result.status === 200) {
      user.value = result.data.Data;
    }
  });
}

function saveUser() {
  BeeTimeClock.administrationUpdateUser(user.value as User).then(result => {
    if (result.status === 200) {
      user.value = result.data.Data;
      showInfoMessage(t('MSG_UPDATE_SUCCESS'));
    }
  });
}

const timestampYears = computed(() => {
  if (!timestampYearMonths.value) return [];
  const years = Object.keys(timestampYearMonths.value);
  return years.sort();
});

const timestampMonths = computed(() => {
  if (!timestampYearMonths.value) return [];
  const months = timestampYearMonths.value[selectedYear.value];
  return months.sort();
});

async function loadTimestampMonths() {
  const result = await BeeTimeClock.administrationTimestampUserMonths(userId);

  if (result.status === 200) {
    timestampYearMonths.value = result.data.Data;
  }
}

function loadTimestampGrouped() {
  BeeTimeClock.administrationTimestampQueryMonthGrouped(userId, selectedYear.value, selectedMonth.value).then((result) => {
    if (result.status === 200) {
      timestampCurrentMonthGrouped.value = result.data.Data.sort((a, b) => new Date(b.Date).getTime() - new Date(a.Date).getTime());
      if (timestampCurrentMonthGrouped.value.length > 0) {
        expanded.value = [timestampCurrentMonthGrouped.value[0].Date.toString()];
      }
    }
  });
}


function loadAbsenceYears() {
  BeeTimeClock.administrationAbsenceYears(userId).then(result => {
    if (result.status === 200) {
      absenceYears.value = result.data.Data
    }
  })
}

function loadAbsences() {
  BeeTimeClock.administrationAbsencesByYear(userId, selectedAbsenceYear.value).then(result => {
    if (result.status === 200) {
      absences.value = result.data.Data;
    }
  })
}

function loadAbsenceReasons() {
  BeeTimeClock.absenceReasons().then(result => {
    if (result.status === 200) {
      absenceReasons.value = result.data.Data
    }
  })
}

onMounted(async () => {
  loadUser();
  loadAbsenceReasons();
  loadAbsenceYears();
  loadAbsences();
  await loadTimestampMonths();
  loadTimestampGrouped()
});

watch(selectedYear, () => {
  console.log('year changed');
  if (timestampYearMonths.value[selectedYear.value].includes(selectedMonth.value)) {
    loadTimestampGrouped();
    return;
  } else {
    selectedMonth.value = timestampYearMonths.value[selectedYear.value][0];
  }
});

watch(selectedMonth, () => {
  console.log('month changed');
  loadTimestampGrouped();
});

watch(selectedAbsenceYear, () => {
  loadAbsences();
})

function getAbsenceReasonDescriptionById(id: number): string {
  const res = absenceReasons.value.filter(s => s.ID == id);
  if (res.length == 0) return '';

  return res[0].Description;
}

const columns = [
  {
    name: 'absenceReasonId',
    label: t('LABEL_REASON'),
    field: 'AbsenceReasonID',
    format: (val: number) => getAbsenceReasonDescriptionById(val),
    align: 'left',
  },
  {
    name: 'absenceFrom',
    label: t('LABEL_FROM'),
    field: 'AbsenceFrom',
    format: (val: Date) => date.formatDate(val, 'DD. MMM. YYYY')
  },
  {
    name: 'absenceTill',
    label: t('LABEL_TILL'),
    field: 'AbsenceTill',
    format: (val: Date) => date.formatDate(val, 'DD. MMM. YYYY')
  },
  {
    name: 'absenceNettoDays',
    label: t('LABEL_NETTO_DAYS'),
    field: 'NettoDays',
  }
]

const pagination = {
  rowsPerPage: 10,
}

</script>

<template>
  <q-page>
    <div v-if="user">
      <q-tabs
        v-model="selectedTab"
        inline-label
        class="bg-primary text-white"
      >
        <q-tab name="common" icon="account_circle" :label="t('LABEL_COMMON')" />
        <q-tab name="worktime" icon="alarms" :label="t('LABEL_WORKTIME')" />
        <q-tab name="absence" icon="event_busy" :label="t('LABEL_ABSENCES')" />
      </q-tabs>
      <q-tab-panels v-model="selectedTab">
        <q-tab-panel name="common">
          <q-card>
            <q-card-section>
              <q-input readonly :label="$t('LABEL_USERNAME')" v-model="user.Username" />
              <q-input :label="$t('LABEL_FIRST_NAME')" v-model="user.FirstName" />
              <q-input :label="$t('LABEL_LAST_NAME')" v-model="user.LastName" />
              <q-select :label="$t('LABEL_ACCESS_LEVEL')" :options="accessLevelOptions" v-model="user.AccessLevel"
                        map-options emit-value />
              <q-select :label="$t('LABEL_OVERTIME_SUBTRACTION_MODEL')" :options="overtimeSubtractions"
                        v-model="user.OvertimeSubtractionModel" map-options emit-value />
              <q-input :label="$t('LABEL_OVERTIME_SUBTRACTION_AMOUNT')" v-model.number="user.OvertimeSubtractionAmount"
                       type="number" />
            </q-card-section>
            <q-card-actions>
              <q-btn :label="$t('BTN_SAVE')" color="primary" @click="saveUser" />
            </q-card-actions>
          </q-card>
        </q-tab-panel>
        <q-tab-panel name="worktime">
          <div class="row">
            <div class="col">
              <q-select v-model="selectedYear" :options="timestampYears" :label="$t('LABEL_YEAR')" />
            </div>
            <div class="col">
              <q-select class="q-ml-md" v-model="selectedMonth" :options="timestampMonths"
                        :label="$t('LABEL_MONTH')" />
            </div>
          </div>
          <div class="row q-mt-md">
            <div class="col">
              <OvertimeMonth v-if="user" v-model:model-user-id="user.ID" v-model:model-month="selectedMonth" v-model:model-year="selectedYear"
                             class="full-width" />
            </div>
          </div>
          <div class="q-pt-lg">
            <WorktimeOverviewTable v-model="timestampCurrentMonthGrouped" @create="loadTimestampGrouped()" />
          </div>
        </q-tab-panel>
        <q-tab-panel name="absence">
          <q-select v-model="selectedAbsenceYear" :label="t('LABEL_YEAR')" :options="absenceYears" class="full-width"/>
          <q-table class="q-mt-lg" :rows="absences" :columns="columns" :flat="flat" :pagination="pagination"/>
        </q-tab-panel>
      </q-tab-panels>
    </div>
  </q-page>
</template>

<style scoped>

</style>
