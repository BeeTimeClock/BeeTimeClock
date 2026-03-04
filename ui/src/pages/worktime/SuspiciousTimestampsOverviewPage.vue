<script setup lang="ts">
import { onMounted, ref } from 'vue';
import type {
  Timestamp,
  TimestampOvertimeReasonRequest,
} from 'src/models/Timestamp';
import BeeTimeClock from 'src/service/BeeTimeClock';
import { useI18n } from 'vue-i18n';
import { date, type QTableColumn, useQuasar } from 'quasar';
import { emptyPagination } from 'src/helper/objects';
import formatDate = date.formatDate;
import TimestampCorrectionDialog from 'components/TimestampCorrectionDialog.vue';
import type { ErrorResponse } from 'src/models/Base';
import { showErrorMessage } from 'src/helper/message';

const q = useQuasar();
const timestamps = ref<Timestamp[]>([]);
const selectedTimestamp = ref<Timestamp>();
const showTimestampCreateDialog = ref(false);
const { t } = useI18n();

const columns = [
  {
    name: 'expand',
  },
  {
    name: 'coming',
    field: 'ComingTimestamp',
    label: t('LABEL_COMING_TIMESTAMP'),
    format: (val: Date) => formatDate(val, 'DD.MM.YYYY HH:mm:ss'),
  },
  {
    name: 'going',
    field: 'GoingTimestamp',
    label: t('LABEL_GOING_TIMESTAMP'),
    format: (val: Date) => formatDate(val, 'DD.MM.YYYY HH:mm:ss'),
  },
  {
    name: 'suspiciousReason',
    field: 'SuspiciousReason',
    label: t('LABEL_SUSPICIOUS_REASON'),
    format: (val: string, row: Timestamp) =>
      getTimestampSuspiciousReasonLabel(row),
  },
  {
    name: 'hasCorrections',
    label: t('LABEL_HAS_CORRECTIONS'),
  },
  {
    name: 'actions',
    label: t('LABEL_ACTION', 2),
  },
] as QTableColumn[];

function loadTimestamps() {
  BeeTimeClock.timestampQuerySuspicious()
    .then((result) => {
      if (result.status === 200) {
        timestamps.value = result.data.Data;
      }
    })
    .catch((error: ErrorResponse) => {
      showErrorMessage(error.response?.data.Message);
    });
}

function editTimestamp(timestamp: Timestamp) {
  console.log('Timestamp: ', timestamp);
  selectedTimestamp.value = timestamp;
  showTimestampCreateDialog.value = true;
}

function setOvertimeReason(timestamp: Timestamp) {
  q.dialog({
    title: t('LABEL_SET_OVERTIME_REASON'),
      message: t('LABEL_SET_OVERTIME_REASON_DESC'),
    prompt: {
      model: '',
      type: 'text', // optional
    },
    cancel: true,
    persistent: true,
  }).onOk((reason) => {
    const overtimeReasonUpdateRequest = {
      OvertimeReason: reason,
    } as TimestampOvertimeReasonRequest;

    BeeTimeClock.timestampSetOvertimeReason(timestamp.ID, overtimeReasonUpdateRequest)
      .then((result) => {
        if (result.status === 200) {
          loadTimestamps();
        }
      })
      .catch((error: ErrorResponse) => {
        showErrorMessage(error.response?.data.Message);
      });
  });
}

function getTimestampSuspiciousReasonLabel(timestamp: Timestamp) {
  switch (timestamp.SuspiciousReason) {
    case 'broken_timestamp':
      return t('LABEL_BROKEN_TIMESTAMP');
    case 'correction':
      return t('LABEL_NEEDS_CORRECTION');
    case 'overtime':
      return t('LABEL_NEEDS_OVERTIME_REASON');
  }

  return '';
}

onMounted(() => {
  loadTimestamps();
});
</script>

<template>
  <q-page padding>
    <q-table
      :rows="timestamps"
      :columns="columns"
      :pagination="emptyPagination"
      hide-pagination
      :title="t('LABEL_SUSPICIOUS_TIMESTAMPS')"
      row-key="ID"
    >
      <template v-slot:body="props">
        <q-tr :props="props">
          <q-td v-for="col in props.cols" :key="col.name" :props="props">
            <q-td auto-width v-if="col.name == 'expand'">
              <q-btn
                v-if="props.row.Corrections.length > 0"
                size="sm"
                color="accent"
                round
                dense
                @click="props.expand = !props.expand"
                :icon="props.expand ? 'remove' : 'add'"
              />
            </q-td>
            <div v-else-if="col.name == 'actions'">
              <q-btn
                v-if="props.row.SuspiciousReason == 'overtime'"
                icon="add_comment"
                color="primary"
                class="q-mr-sm"
                @click="setOvertimeReason(props.row)"
              />
              <q-btn
                icon="edit"
                color="primary"
                @click="editTimestamp(props.row)"
              />
            </div>
            <div v-else-if="col.name == 'hasCorrections'">
              <q-icon
                size="large"
                :name="
                  props.row.Corrections.length > 0 ? 'check_circle' : 'cancel'
                "
                :color="props.row.Corrections.length > 0 ? 'positive' : ''"
              />
            </div>
            <div v-else>
              {{ col.value }}
            </div>
          </q-td>
        </q-tr>
        <q-tr v-show="props.expand" :props="props">
          <q-td colspan="100%">
            <q-list>
              <q-item
                v-for="correction in props.row.Corrections"
                :key="correction.ID"
              >
                <q-item-section>
                  <q-item-label caption>{{ t('LABEL_REASON') }}</q-item-label>
                  <q-item-label>{{ correction.ChangeReason }}</q-item-label>
                </q-item-section>
                <q-item-section>
                  <q-item-label caption>{{
                    t('LABEL_OLD_COMING')
                  }}</q-item-label>
                  <q-item-label
                    >{{
                      date.formatDate(
                        correction.OldComingTimestamp,
                        'DD.MM.YYYY HH:mm:ss',
                      )
                    }}
                  </q-item-label>
                </q-item-section>
                <q-item-section>
                  <q-item-label caption>{{
                    t('LABEL_OLD_GOING')
                  }}</q-item-label>
                  <q-item-label
                    >{{
                      date.formatDate(
                        correction.OldGoingTimestamp,
                        'DD.MM.YYYY HH:mm:ss',
                      )
                    }}
                  </q-item-label>
                </q-item-section>
              </q-item>
            </q-list>
          </q-td>
        </q-tr>
      </template>
    </q-table>
    <TimestampCorrectionDialog
      v-if="selectedTimestamp"
      v-model:show="showTimestampCreateDialog"
      v-model="selectedTimestamp"
      @refresh="loadTimestamps"
    />
  </q-page>
</template>

<style scoped></style>
