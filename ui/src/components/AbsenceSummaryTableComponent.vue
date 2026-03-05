<script setup lang="ts">
import { date, type QTableColumn } from 'quasar';
import { QCalendarMonth, today } from '@quasar/quasar-ui-qcalendar';
import type { AbsenceSummaryItem } from 'src/models/Absence';
import type { User } from 'src/models/Authentication';
import { type PropType, ref } from 'vue';
import { computed } from 'vue';
import { useI18n } from 'vue-i18n';
import { useAuthStore } from 'stores/microsoft-auth';

const { t } = useI18n();
const auth = useAuthStore();
const selectedTab = ref('calendar');

const showOnlyFuture = ref(true);

const calendarYear = computed({
  get() {
    return new Date(selectedDate.value).getFullYear();
  },
  set(value: number) {
    const month = new Date(selectedDate.value).getMonth();
    selectedDate.value = date.formatDate(
      new Date(value, month, 1),
      'YYYY-MM-DD',
    );
  },
});

const calendarMonth = computed({
  get() {
    return new Date(selectedDate.value).getMonth() + 1;
  },
  set(value: number) {
    const year = new Date(selectedDate.value).getFullYear();
    selectedDate.value = date.formatDate(
      new Date(year, value - 1, 1),
      'YYYY-MM-DD',
    );
  },
});

const calerndarYears = computed(() => {
  if (!props.modelValue) return [];

  const years = props.modelValue.map((item) =>
    new Date(item.AbsenceFrom).getFullYear(),
  );

  return [...new Set(years)].sort((a, b) => a - b);
});

const props = defineProps({
  modelValue: {
    type: Array as PropType<AbsenceSummaryItem[]>,
    required: true,
  },
  flat: {
    type: Boolean,
    default: false,
  },
  title: {
    type: String,
  },
  showReason: {
    type: Boolean,
    default: false,
  },
});

const calendar = ref<QCalendarMonth>();
const selectedDate = ref(today());

function getEventsForDate(date: string) {
  if (!props.modelValue) return [];
  const values = props.modelValue.filter((item) => {
    const from = new Date(item.AbsenceFrom);
    const till = new Date(item.AbsenceTill);
    const current = new Date(date);
    return current >= from && current <= till;
  });
  return values;
}

const getTitle = computed(() => {
  if (props.title) {
    return props.title;
  }
  return t('LABEL_EMPLOYEE_ABSENCES');
});

const rows = computed(() => {
  if (!props.modelValue) return [];

  const data = props.modelValue;
  return data
    .filter((a) => {
      if (showOnlyFuture.value) {
        const today = new Date();
        return new Date(a.AbsenceTill) >= today;
      }
      return true;
    })
    .sort(
      (a, b) =>
        new Date(a.AbsenceFrom).getTime() - new Date(b.AbsenceFrom).getTime(),
    );
});

const columns = [
  {
    name: 'absenceFrom',
    label: t('LABEL_FROM'),
    field: 'AbsenceFrom',
    format: (val: Date) => date.formatDate(val, 'DD. MMM. YYYY'),
  },
  {
    name: 'absenceTill',
    label: t('LABEL_TILL'),
    field: 'AbsenceTill',
    format: (val: Date) => date.formatDate(val, 'DD. MMM. YYYY'),
  },
  {
    name: 'absenceNettoDays',
    label: t('LABEL_NETTO_DAYS'),
    field: 'NettoDays',
  },
  {
    name: 'user',
    label: t('LABEL_USER'),
    field: 'User',
    format: (val: User) => `${val.FirstName} ${val.LastName}`,
  },
] as QTableColumn[];

const pagination = {
  rowsPerPage: 10,
};

if (auth.isAdministrator() || props.showReason) {
  columns.push({
    name: 'absenceReason',
    label: t('LABEL_REASON'),
    field: 'Reason',
  } as QTableColumn);
}
</script>

<template>
  <q-tabs v-model="selectedTab" inline-label class="bg-primary text-white">
    <q-tab
      name="calendar"
      icon="calendar_month"
      :label="t('LABEL_CALENDAR_VIEW')"
    />
    <q-tab name="table" icon="table_rows" :label="t('LABEL_TABLE_VIEW')" />
  </q-tabs>
  <q-tab-panels v-model="selectedTab">
    <q-tab-panel name="calendar">
      <div class="row q-mb-md">
        <div class="col q-pa-md">
          <q-select
            :options="calerndarYears"
            :label="t('LABEL_YEAR')"
            v-model="calendarYear"
          />
        </div>
        <div class="col q-pa-md">
          <q-select
            :options="[1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12]"
            :label="t('LABEL_MONTH')"
            v-model="calendarMonth"
          />
        </div>
      </div>
      <div class="row q-mb-md justify-center" v-if="calendar">
        <q-btn-group push>
          <q-btn
            push
            color="primary"
            icon="keyboard_double_arrow_left"
            @click="calendar.prev()"
          />
          <q-btn
            push
            color="primary"
            :label="t('LABEL_TODAY')"
            icon="today"
            @click="selectedDate = today()"
          />
          <q-btn
            push
            color="primary"
            icon="keyboard_double_arrow_right"
            @click="calendar.next()"
          />
        </q-btn-group>
      </div>
      <q-calendar-month
        ref="calendar"
        :weekdays="[1, 2, 3, 4, 5]"
        v-model="selectedDate"
        :day-min-height="80"
      >
        <template #day="{ scope: { timestamp } }">
          <template
            v-for="event in getEventsForDate(timestamp.date)"
            :key="event.id"
          >
            <div class="q-calendar__ellipsis bg-primary q-mt-sm">
              <div class="title q-calendar__ellipsis">
                {{ event.calendarCaption }}
              </div>
            </div>
          </template>
        </template>
      </q-calendar-month>
    </q-tab-panel>
    <q-tab-panel name="table">
      <q-table
        :rows="rows"
        :columns="columns"
        :flat="flat"
        :pagination="pagination"
      >
        <template v-slot:top>
          <div class="col-2 q-table__title">{{ getTitle }}</div>
          <q-space />
          <q-checkbox
            v-model="showOnlyFuture"
            :label="t('LABEL_ONLY_FUTURE_ABSENCES')"
          />
        </template>
      </q-table>
    </q-tab-panel>
  </q-tab-panels>
</template>

<style scoped></style>
