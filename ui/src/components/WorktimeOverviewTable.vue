<script setup lang="ts">
import { formatIndustryHourMinutes } from 'src/helper/formatter';
import { ref } from 'vue';
import type { Timestamp, TimestampGroup } from 'src/models/Timestamp';
import formatDate = date.formatDate;
import { date, type QTableColumn } from 'quasar';
import { useI18n } from 'vue-i18n';
import TimestampCorrectionDialog from 'components/TimestampCorrectionDialog.vue';

const { t } = useI18n();

const value = defineModel<TimestampGroup[]>({ required: true });

const emits = defineEmits(['create']);

const columns = [
  {
    name: 'Date',
    required: true,
    label: t('LABEL_DATE'),
    align: 'left',
    field: 'Date',
    format: (val: Date) => `${formatDate(val, 'ddd. DD.MM.YYYY')}`,
    sortable: true,
  },
  {
    name: 'WorkingHours',
    align: 'left',
    label: t('LABEL_WORKING_HOURS'),
    field: 'WorkingHours',
    format: (val: number) => formatIndustryHourMinutes(val),
  },
  {
    name: 'SubtractedHours',
    align: 'left',
    label: t('LABEL_SUBTRACTED_HOURS'),
    field: 'SubtractedHours',
    format: (val: number) => formatIndustryHourMinutes(val),
  },
  {
    name: 'OvertimeHours',
    align: 'left',
    label: t('LABEL_OVERTIME_HOURS'),
    field: 'OvertimeHours',
    format: (val: number) => formatIndustryHourMinutes(val),
  },
  {
    name: 'IsHomeoffice',
    align: 'left',
    label: t('LABEL_HOMEOFFICE'),
    field: 'IsHomeoffice',
  },
] as QTableColumn[];

const pagination = {
  page: 1,
  rowsPerPage: 0,
};

function formatDateTemplate(date: Date, format: string): string {
  return formatDate(date, format);
}

const selectedTimestamp = ref<Timestamp | null>(null);
const isTimestampCorrentViewVisible = ref(false);
const promptNewTimestampCorrection = ref(false);

function promptTimestampCorrectionView(timestamp: Timestamp) {
  selectedTimestamp.value = timestamp;
  isTimestampCorrentViewVisible.value = true;
}
</script>

<template>
  <q-table
    flat
    bordered
    :rows="value"
    :columns="columns"
    row-key="Date"
    :pagination="pagination"
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
            size="sm"
            color="accent"
            round
            dense
            @click="props.expand = !props.expand"
            :icon="props.expand ? 'remove' : 'add'"
          />
        </q-td>
        <q-td v-for="col in props.cols" :key="col.name" :props="props">
          <div v-if="col.name == 'IsHomeoffice'">
            <q-icon
              size="large"
              :name="props.row.IsHomeoffice ? 'check_circle' : 'cancel'"
              :color="props.row.IsHomeoffice ? 'positive' : ''"
            />
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
                <q-th class="text-left">{{
                  t('LABEL_COMING_TIMESTAMP')
                }}</q-th>
                <q-th class="text-left">{{ t('LABEL_GOING_TIMESTAMP') }}</q-th>
                <q-th></q-th>
              </q-tr>
            </thead>
            <tbody>
              <tr v-for="timestamp in props.row.Timestamps" :key="timestamp.ID">
                <td>
                  {{ formatDateTemplate(timestamp.ComingTimestamp, 'HH:mm') }}
                  <q-icon
                    size="large"
                    :name="timestamp.IsHomeoffice ? 'check_circle' : 'cancel'"
                    :color="timestamp.IsHomeoffice ? 'positive' : ''"
                  >
                    <q-tooltip>Homeoffice</q-tooltip>
                  </q-icon>
                </td>
                <td>
                  <div
                    v-if="
                      timestamp.GoingTimestamp &&
                      new Date(timestamp.GoingTimestamp).getFullYear() != 1
                    "
                  >
                    {{ formatDateTemplate(timestamp.GoingTimestamp, 'HH:mm') }}
                    <q-icon
                      size="large"
                      :name="
                        timestamp.IsHomeofficeGoing ? 'check_circle' : 'cancel'
                      "
                      :color="timestamp.IsHomeofficeGoing ? 'positive' : ''"
                    >
                      <q-tooltip>Homeoffice</q-tooltip>
                    </q-icon>
                  </div>
                </td>
                <td>
                  <q-btn
                    color="primary"
                    class="q-mr-md"
                    :disable="timestamp.Corrections.length == 0"
                    icon="pending_actions"
                    @click="promptTimestampCorrectionView(timestamp)"
                  />
                  <q-btn
                    color="primary"
                    icon="edit"
                    @click="
                      selectedTimestamp = timestamp;
                      promptNewTimestampCorrection = true;
                    "
                  />
                </td>
              </tr>
            </tbody>
          </q-markup-table>
        </q-td>
      </q-tr>
    </template>
  </q-table>
  <TimestampCorrectionDialog
    v-if="selectedTimestamp"
    v-model:show="promptNewTimestampCorrection"
    v-model="selectedTimestamp"
    @refresh="emits('create')"
  />
  <q-dialog v-model="isTimestampCorrentViewVisible" persistent>
    <q-card class="q-dialog-plugin full-width">
      <q-card-section>
        <div class="text-h6">{{ t('LABEL_TIMESTAMP_CORRECTION_VIEW') }}</div>
      </q-card-section>
      <q-card-section>
        <q-markup-table flat>
          <thead>
            <tr>
              <th>{{ t('LABEL_REASON') }}</th>
              <th>{{ t('LABEL_OLD_COMING') }}</th>
              <th>{{ t('LABEL_OLD_GOING') }}</th>
              <th>{{ t('LABEL_CREATED_AT') }}</th>
            </tr>
          </thead>
          <tbody>
            <tr
              v-for="correction in selectedTimestamp?.Corrections"
              :key="correction.ID"
            >
              <td>{{ correction.ChangeReason }}</td>
              <td>
                {{
                  date.formatDate(
                    correction.OldComingTimestamp,
                    'DD.MM.YYYY HH:mm',
                  )
                }}
              </td>
              <td>
                {{
                  date.formatDate(
                    correction.OldGoingTimestamp,
                    'DD.MM.YYYY HH:mm',
                  )
                }}
              </td>
              <td>
                {{ date.formatDate(correction.CreatedAt, 'DD.MM.YYYY HH:mm') }}
              </td>
            </tr>
          </tbody>
        </q-markup-table>
      </q-card-section>
      <q-card-actions>
        <q-btn v-close-popup :label="t('BTN_CLOSE')" color="primary" />
      </q-card-actions>
    </q-card>
  </q-dialog>
</template>

<style scoped></style>
