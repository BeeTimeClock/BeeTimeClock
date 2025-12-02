<script setup lang="ts">
import type { User } from 'src/models/Authentication';
import { computed, onMounted, ref, watch } from 'vue';
import BeeTimeClock from 'src/service/BeeTimeClock';
import { useRoute } from 'vue-router';
import { useI18n } from 'vue-i18n';
import { showErrorMessage, showInfoMessage } from 'src/helper/message';
import WorktimeOverviewTable from 'components/WorktimeOverviewTable.vue';
import OvertimeMonth from 'components/OvertimeMonth.vue';
import type {
  TimestampGroup,
  TimestampYearMonthGrouped,
} from 'src/models/Timestamp';
import type { AbsenceReason } from 'src/models/Absence';
import { Absence } from 'src/models/Absence';
import type { QTableColumn} from 'quasar';
import { useQuasar } from 'quasar';
import type { ErrorResponse } from 'src/models/Base';
import { emptyPagination } from 'src/helper/objects';
import { OvertimeMonthQuota } from 'src/models/Overtime';
import { formatIndustryHourMinutes } from 'src/helper/formatter';

const { t } = useI18n();

const route = useRoute();
const userId = computed(() => {
  return parseInt(route.params.userId as string);
})
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
const overtimeTotal = ref<number>();
const q = useQuasar();

const accessLevelOptions = [
  {
    value: 'admin',
    label: t('LABEL_ADMINISTRATOR'),
  },
  {
    value: 'user',
    label: t('LABEL_USER'),
  },
];

const overtimeSubtractions = [
  {
    value: 'percentage',
    label: t('LABEL_PERCENTAGE'),
  },
  {
    value: 'hours',
    label: t('LABEL_HOUR', 2),
  },
];

const overtimeColumns = [
  {
    name: 'Date',
    required: true,
    label: t('LABEL_DATE'),
    align: "left",
    field: 'identifier',
    sortable: true,
  },
  {
    name: 'Hours',
    required: true,
    label: t('LABEL_HOUR', 2),
    field: 'Hours',
    align: "right",
    sortable: true,
    format: (val: number) => `${val.toFixed(2)}`,
  },
  {
    name: 'actions',
    label: t('LABEL_ACTION', 2),
  },
] as QTableColumn[];

function loadUser() {
  BeeTimeClock.administrationGetUserById(userId.value)
    .then((result) => {
      if (result.status === 200) {
        user.value = result.data.Data;
      }
    })
    .catch((error: ErrorResponse) => {
      showErrorMessage(error.message);
    });
}

function saveUser() {
  BeeTimeClock.administrationUpdateUser(user.value as User)
    .then((result) => {
      if (result.status === 200) {
        user.value = result.data.Data;
        showInfoMessage(t('MSG_UPDATE_SUCCESS'));
      }
    })
    .catch((error: ErrorResponse) => {
      showErrorMessage(error.message);
    });
}

const timestampYears = computed(() => {
  if (!timestampYearMonths.value) return [];
  const years = Object.keys(timestampYearMonths.value);
  return years.sort();
});

const timestampMonths = computed(() => {
  if (!timestampYearMonths.value) return [];
  const months = timestampYearMonths.value[selectedYear.value]!;
  return months.sort();
});

async function loadTimestampMonths() {
  const result = await BeeTimeClock.administrationTimestampUserMonths(userId.value);

  if (result.status === 200) {
    timestampYearMonths.value = result.data.Data;
  }
}

function loadTimestampGrouped() {
  BeeTimeClock.administrationTimestampQueryMonthGrouped(
    userId.value,
    selectedYear.value,
    selectedMonth.value,
  )
    .then((result) => {
      if (result.status === 200) {
        timestampCurrentMonthGrouped.value = result.data.Data.sort(
          (a, b) => new Date(b.Date).getTime() - new Date(a.Date).getTime(),
        );
        if (timestampCurrentMonthGrouped.value.length > 0) {
          expanded.value = [
            timestampCurrentMonthGrouped.value[0]!.Date.toString(),
          ];
        }
      }
    })
    .catch((error: ErrorResponse) => {
      showErrorMessage(error.message);
    });
}

function loadAbsenceYears() {
  BeeTimeClock.administrationAbsenceYears(userId.value)
    .then((result) => {
      if (result.status === 200) {
        absenceYears.value = result.data.Data;
      }
    })
    .catch((error: ErrorResponse) => {
      showErrorMessage(error.message);
    });
}

function loadAbsences() {
  BeeTimeClock.administrationAbsencesByYear(userId.value, selectedAbsenceYear.value)
    .then((result) => {
      if (result.status === 200) {
        absences.value = result.data.Data.map((s) => Absence.fromApi(s));
      }
    })
    .catch((error: ErrorResponse) => {
      showErrorMessage(error.message);
    });
}

function loadAbsenceReasons() {
  BeeTimeClock.absenceReasons()
    .then((result) => {
      if (result.status === 200) {
        absenceReasons.value = result.data.Data;
      }
    })
    .catch((error: ErrorResponse) => {
      showErrorMessage(error.message);
    });
}


const overtimeQuotas = ref<OvertimeMonthQuota[]>([]);


function loadOvertimeQuotas() {
  if (!user.value) return;
  BeeTimeClock.administrationOvertimeMonthQuotas(user.value.ID).then((result) => {
    if (result.status === 200) {
      overtimeQuotas.value = result.data.Data.sort((a, b) => b.Year - a.Year || b.Month - a.Month).map((s) =>
        OvertimeMonthQuota.fromApi(s)
      );
    }
  }).catch((error: ErrorResponse) => {
    showErrorMessage(error.response?.data.Message);
  });
}

function calculateOvertimeMonth(overtimeMonthQuota: OvertimeMonthQuota) {
  if (!user.value) return;
  BeeTimeClock.administrationCalculateOvertimeMonthQuota(
    user.value.ID,
    overtimeMonthQuota.Year,
    overtimeMonthQuota.Month
  ).then((result) => {
    if (result.status === 200) {
      loadOvertimeQuotas();
      loadOvertimeTotal();
    }
  }).catch((error: ErrorResponse) => {
    showErrorMessage(error.response?.data.Message);
  });
}

function loadOvertimeTotal() {
  if (!user.value) return;

  BeeTimeClock.administrationOvertimeTotal(user.value.ID).then(result => {
    if (result.status === 200) {
      overtimeTotal.value = result.data.Data.Total;
    }
  }).catch((error: ErrorResponse) => {
    showErrorMessage(error.response?.data.Message);
  });
}

onMounted(async () => {
  loadUser();
  loadAbsenceReasons();
  loadAbsenceYears();
  loadAbsences();
  await loadTimestampMonths();
  loadTimestampGrouped();
  loadOvertimeTotal();
  loadOvertimeQuotas();
});

watch(selectedYear, () => {
  console.log('year changed');
  if (
    timestampYearMonths.value[selectedYear.value]!.includes(selectedMonth.value)
  ) {
    loadTimestampGrouped();
    return;
  } else {
    selectedMonth.value = timestampYearMonths.value[selectedYear.value]![0]!;
  }
});

watch(selectedMonth, () => {
  console.log('month changed');
  loadTimestampGrouped();
});

watch(selectedAbsenceYear, () => {
  loadAbsences();
});

function getAbsenceReasonDescriptionById(id: number): string {
  const res = absenceReasons.value.filter((s) => s.ID == id);
  if (res.length == 0) return '';

  return res[0]!.Description;
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
    field: 'formatFromFull',
    sortable: true,
  },
  {
    name: 'absenceTill',
    label: t('LABEL_TILL'),
    field: 'formatTillFull',
    sortable: true,
  },
  {
    name: 'absenceNettoDays',
    label: t('LABEL_NETTO_DAYS'),
    field: 'NettoDays',
  },
] as QTableColumn[];

const pagination = {
  rowsPerPage: 10,
  sortBy: 'absenceFrom',
};

function deleteUserAbsence(absence: Absence) {
  q.dialog({
    message: t('MSG_DELETE', {
      item: t('LABEL_ABSENCE'),
      identifier: `${absence.formatFrom} - ${absence.formatTill}`,
    }),
    cancel: true,
    persistent: true,
  }).onOk(() => {
    BeeTimeClock.deleteAbsence(absence.ID)
      .then((result) => {
        if (result.status === 204) {
          loadAbsences();
          showInfoMessage(t('MSG_DELETE_SUCCESS'));
        }
      })
      .catch((error: ErrorResponse) => {
        showErrorMessage(error.message);
      });
  });
}
</script>

<template>
  <q-page>
    <div v-if="user">
      <q-tabs v-model="selectedTab" inline-label class="bg-primary text-white">
        <q-tab name="common" icon="account_circle" :label="t('LABEL_COMMON')" />
        <q-tab name="worktime" icon="alarms" :label="t('LABEL_WORKTIME')" />
        <q-tab name="overtime" icon="more_time" :label="t('LABEL_OVERTIME')" />
        <q-tab
          name="absence"
          icon="event_busy"
          :label="t('LABEL_ABSENCE', 2)"
        />
      </q-tabs>
      <q-tab-panels v-model="selectedTab">
        <q-tab-panel name="common">
          <q-card>
            <q-card-section>
              <q-input
                readonly
                :label="t('LABEL_USERNAME')"
                v-model="user.Username"
              />
              <q-input
                :label="t('LABEL_FIRST_NAME')"
                v-model="user.FirstName"
              />
              <q-input :label="t('LABEL_LAST_NAME')" v-model="user.LastName" />
              <q-select
                :label="t('LABEL_ACCESS_LEVEL')"
                :options="accessLevelOptions"
                v-model="user.AccessLevel"
                map-options
                emit-value
              />
              <q-select
                :label="t('LABEL_OVERTIME_SUBTRACTION_MODEL')"
                :options="overtimeSubtractions"
                v-model="user.OvertimeSubtractionModel"
                map-options
                emit-value
              />
              <q-input
                :label="t('LABEL_OVERTIME_SUBTRACTION_AMOUNT')"
                v-model.number="user.OvertimeSubtractionAmount"
                type="number"
              />
              <q-input
                :label="t('LABEL_STAFF_NUMBER')"
                v-model.number="user.StaffNumber"
                type="number"
              />
            </q-card-section>
            <q-card-actions>
              <q-btn
                :label="t('BTN_SAVE')"
                color="primary"
                @click="saveUser"
              />
            </q-card-actions>
          </q-card>
        </q-tab-panel>
        <q-tab-panel name="worktime">
          <div class="row">
            <div class="col">
              <q-select
                v-model="selectedYear"
                :options="timestampYears"
                :label="t('LABEL_YEAR')"
              />
            </div>
            <div class="col">
              <q-select
                class="q-ml-md"
                v-model="selectedMonth"
                :options="timestampMonths"
                :label="t('LABEL_MONTH')"
              />
            </div>
          </div>
          <div class="row q-mt-md">
            <div class="col">
              <OvertimeMonth
                v-if="user"
                v-model:model-user-id="user.ID"
                v-model:model-month="selectedMonth"
                v-model:model-year="selectedYear"
                class="full-width"
              />
            </div>
          </div>
          <div class="q-pt-lg">
            <WorktimeOverviewTable
              v-model="timestampCurrentMonthGrouped"
              @create="loadTimestampGrouped()"
            />
          </div>
        </q-tab-panel>
        <q-tab-panel name="overtime">
          <q-card v-if="overtimeTotal" class="q-mb-lg">
            <q-card-section class="bg-primary text-h6 text-white">
              {{t('LABEL_OVERTIME_TOTAL')}}
            </q-card-section>
            <q-card-section class="text-h6 text-center">
              {{ formatIndustryHourMinutes(overtimeTotal) }}
            </q-card-section>
          </q-card>
          <q-table
            flat
            bordered
            :rows="overtimeQuotas"
            :columns="overtimeColumns"
            row-key="identifier"
            :pagination="emptyPagination"
            hide-pagination
          >
            <template v-slot:header="props">
              <q-tr :props="props">
                <q-th auto-width />
                <q-th v-for="col in props.cols" :key="col.name" :props="props">
                  {{ col.label }}
                </q-th>
              </q-tr>
            </template>

            <template v-slot:body="props">
              <q-tr :props="props">
                <q-td auto-width>
                  <q-btn
                    size="xs"
                    color="secondary"
                    round
                    dense
                    @click="props.expand = !props.expand"
                    :icon="props.expand ? 'remove' : 'add'"
                  />
                </q-td>
                <q-td v-for="col in props.cols" :key="col.name" :props="props">
                  <template v-if="col.name === 'actions'">
                    <q-btn icon="refresh" @click="calculateOvertimeMonth(props.row)" color="primary"/>
                  </template>
                  <template v-else>
                    {{ col.value }}
                  </template>
                </q-td>
              </q-tr>
              <q-tr v-show="props.expand" :props="props">
                <q-td colspan="100%">
                  <q-list separator>
                    <q-item
                      v-for="entry in props.row.Summary"
                      :key="entry.Identifier"
                    >
                      <q-item-section>
                        <q-item-label caption>{{ entry.Source }}</q-item-label>
                      </q-item-section>
                      <q-item-section>
                        <q-item-section>{{ entry.Value.toFixed(2) }}</q-item-section>
                      </q-item-section>
                    </q-item>
                  </q-list>
                </q-td>
              </q-tr>
            </template>
          </q-table>
        </q-tab-panel>
        <q-tab-panel name="absence">
          <q-select
            v-model="selectedAbsenceYear"
            :label="t('LABEL_YEAR')"
            :options="absenceYears"
            class="full-width"
          />
          <q-table
            class="q-mt-lg"
            :rows="absences"
            :columns="columns"
            flat
            :pagination="pagination"
          >
            <template v-slot:header="props">
              <q-tr :props="props">
                <q-th v-for="col in props.cols" :key="col.name" :props="props">
                  {{ col.label }}
                </q-th>
                <q-th auto-width />
              </q-tr>
            </template>
            <template v-slot:body="props">
              <q-tr :props="props" :key="`m_${props.row.index}`">
                <q-td v-for="col in props.cols" :key="col.name" :props="props">
                  {{ col.value }}
                </q-td>
                <q-td auto-width>
                  <q-btn
                    icon="delete"
                    color="negative"
                    @click="deleteUserAbsence(props.row)"
                  />
                </q-td>
              </q-tr>
            </template>
          </q-table>
        </q-tab-panel>
      </q-tab-panels>
    </div>
  </q-page>
</template>

<style scoped></style>
