<script setup lang="ts">
import WorktimeOverviewTable from 'components/WorktimeOverviewTable.vue';
import OvertimeMonth from 'components/OvertimeMonth.vue';
import { computed, onMounted, ref, watch } from 'vue';
import type { User } from 'src/models/Authentication';
import type {
  Timestamp,
  TimestampGroup,
  TimestampYearMonthGrouped,
} from 'src/models/Timestamp';
import { useI18n } from 'vue-i18n';
import BeeTimeClock from 'src/service/BeeTimeClock';
import type { ErrorResponse } from 'src/models/Base';
import { showErrorMessage } from 'src/helper/message';
import { useRoute } from 'vue-router';
import { OvertimeMonthQuota } from 'src/models/Overtime';
import OvertimeTableComponent from 'components/overtime/OvertimeTableComponent.vue';
import { formatIndustryHourMinutes } from 'src/helper/formatter';

const { t } = useI18n();
const route = useRoute();
const teamId = computed(() => {
  return parseInt(route.params.teamId as string);
});
const userId = computed(() => {
  return parseInt(route.params.userId as string);
});

const overtimeQuotas = ref<OvertimeMonthQuota[]>([]);
const user = ref<User>();
const selectedTab = ref('worktime');
const timestampYearMonths = ref<TimestampYearMonthGrouped>({});
const timestampCurrentMonthGrouped = ref<TimestampGroup[]>([]);
const expanded = ref(['']);
const overtimeTotal = ref<number>();

const selectedYear = ref<number>(new Date().getFullYear());
const selectedMonth = ref<number>(new Date().getMonth() + 1);

const timestampYears = computed(() => {
  if (!timestampYearMonths.value) return [];
  const years = Object.keys(timestampYearMonths.value);
  return years.sort();
});

const timestampMonths = computed(() => {
  if (!timestampYearMonths.value) return [];
  const months = timestampYearMonths.value[selectedYear.value]!;
  return months != undefined ? months.sort() : months;
});

async function loadUser() {
  await BeeTimeClock.teamGetUserById(teamId.value, userId.value)
    .then((result) => {
      if (result.status === 200) {
        user.value = result.data.Data;
      }
    })
    .catch((error: ErrorResponse) => {
      showErrorMessage(error.message);
    });
}

function loadOvertimeTotal() {
  if (!user.value) return;

  BeeTimeClock.teamUserOvertimeTotal(teamId.value, user.value.ID)
    .then((result) => {
      if (result.status === 200) {
        overtimeTotal.value = result.data.Data.Total;
      }
    })
    .catch((error: ErrorResponse) => {
      showErrorMessage(error.response?.data.Message);
    });
}

function loadOvertimeQuotas() {
  if (!user.value) return;
  BeeTimeClock.teamUserOvertimeMonthQuotas(teamId.value, user.value.ID)
    .then((result) => {
      if (result.status === 200) {
        overtimeQuotas.value = result.data.Data.sort(
          (a, b) => b.Year - a.Year || b.Month - a.Month,
        ).map((s) => OvertimeMonthQuota.fromApi(s));
      }
    })
    .catch((error: ErrorResponse) => {
      showErrorMessage(error.response?.data.Message);
    });
}

function calculateOvertimeMonth(overtimeMonthQuota: OvertimeMonthQuota) {
  if (!user.value) return;
  BeeTimeClock.teamUserCalculateOvertimeMonthQuota(
    teamId.value,
    user.value.ID,
    overtimeMonthQuota.Year,
    overtimeMonthQuota.Month,
  )
    .then((result) => {
      if (result.status === 200) {
        loadOvertimeQuotas();
        loadOvertimeTotal();
      }
    })
    .catch((error: ErrorResponse) => {
      showErrorMessage(error.response?.data.Message);
    });
}

function loadTimestampGrouped() {
  BeeTimeClock.teamTimestampQueryMonthGrouped(
    teamId.value,
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

function deleteTimestamp(timestamp: Timestamp) {
  BeeTimeClock.teamTimestampDelete(teamId.value, userId.value, timestamp.ID)
    .then((result) => {
      if (result.status === 204) {
        loadTimestampGrouped();
      }
    })
    .catch((error: ErrorResponse) => {
      showErrorMessage(error.message);
    });
}

async function loadTimestampMonths() {
  const result = await BeeTimeClock.teamTimestampUserMonths(
    teamId.value,
    userId.value,
  );

  if (result.status === 200) {
    timestampYearMonths.value = result.data.Data;
  }
}

watch(selectedYear, () => {
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
  loadTimestampGrouped();
});

onMounted(async () => {
  await loadUser();
  await loadTimestampMonths();
  loadTimestampGrouped();
  loadOvertimeQuotas();
  loadOvertimeTotal();
});
</script>

<template>
  <q-page padding>
    <div v-if="user">
      <q-card>
        <q-card-section>
          <div class="row">
            <q-input
              class="col"
              readonly
              :label="t('LABEL_USERNAME')"
              v-model="user.Username"
            />
            <q-input
              class="col"
              :label="t('LABEL_STAFF_NUMBER')"
              v-model="user.StaffNumber"
              readonly
            />
          </div>

          <div class="row">
            <q-input
              class="col"
              :label="t('LABEL_FIRST_NAME')"
              v-model="user.FirstName"
              readonly
            />
            <q-input
              class="col"
              :label="t('LABEL_LAST_NAME')"
              v-model="user.LastName"
              readonly
            />
          </div>
          <div class="row">
            <q-input
              class="col"
              :label="t('LABEL_OVERTIME_TOTAL')"
              :model-value="formatIndustryHourMinutes(overtimeTotal ?? 0)"
              readonly
            />
          </div>
        </q-card-section>
      </q-card>
      <q-card class="q-mt-md">
        <q-tabs
          v-model="selectedTab"
          inline-label
          class="bg-primary text-white"
        >
          <q-tab name="worktime" icon="alarms" :label="t('LABEL_WORKTIME')" />
          <q-tab
            name="overtime"
            icon="more_time"
            :label="t('LABEL_OVERTIME')"
          />
        </q-tabs>
        <q-tab-panels v-model="selectedTab">
          <q-tab-panel name="worktime">
            <div class="row">
              <div class="col">
                <q-select
                  v-model.number="selectedYear"
                  :options="timestampYears"
                  :label="t('LABEL_YEAR')"
                />
              </div>
              <div class="col">
                <q-select
                  class="q-ml-md"
                  v-model.number="selectedMonth"
                  :options="timestampMonths"
                  :label="t('LABEL_MONTH')"
                />
              </div>
            </div>
            <div class="row q-mt-md">
              <div class="col">
                <OvertimeMonth
                  v-if="user"
                  v-model:model-team-id="teamId"
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
                @delete="deleteTimestamp"
                allow-delete
                disable-edit
              />
            </div>
          </q-tab-panel>
          <q-tab-panel name="overtime">
            <OvertimeTableComponent
              v-model="overtimeQuotas"
              @calculateOvertimeMonth="calculateOvertimeMonth"
            />
          </q-tab-panel>
        </q-tab-panels>
      </q-card>
    </div>
    <q-inner-loading :showing="!user" />
  </q-page>
</template>

<style scoped></style>
