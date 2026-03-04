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

const { t } = useI18n();
const route = useRoute();
const teamId = computed(() => {
  return parseInt(route.params.teamId as string);
});
const userId = computed(() => {
  return parseInt(route.params.userId as string);
});

const user = ref<User>();
const selectedTab = ref('worktime');
const timestampYearMonths = ref<TimestampYearMonthGrouped>({});
const timestampCurrentMonthGrouped = ref<TimestampGroup[]>([]);
const expanded = ref(['']);

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
  const result = await BeeTimeClock.administrationTimestampUserMonths(
    userId.value,
  );

  if (result.status === 200) {
    timestampYearMonths.value = result.data.Data;
  }
}

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

onMounted(async () => {
  await loadUser();
  await loadTimestampMonths();
  loadTimestampGrouped();
});
</script>

<template>
  <q-page>
    <div v-if="user">
      <q-tabs v-model="selectedTab" inline-label class="bg-primary text-white">
        <q-tab name="worktime" icon="alarms" :label="t('LABEL_WORKTIME')" />
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
            />
          </div>
        </q-tab-panel>
      </q-tab-panels>
    </div>
    <q-inner-loading :showing="!user" />
  </q-page>
</template>

<style scoped></style>
