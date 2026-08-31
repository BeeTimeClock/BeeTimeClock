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
  Timestamp,
  TimestampGroup,
  TimestampYearMonthGrouped,
} from 'src/models/Timestamp';
import { Absence } from 'src/models/Absence';
import { date, type QTableColumn } from 'quasar';
import { useQuasar } from 'quasar';
import type { ErrorResponse } from 'src/models/Base';
import { emptyPagination } from 'src/helper/objects';
import { OvertimeMonthQuota } from 'src/models/Overtime';
import { formatIndustryHourMinutes } from 'src/helper/formatter';
import { type MissingDay } from 'src/models/MissingDays';
import formatDate = date.formatDate;
import OvertimeTableComponent from 'components/overtime/OvertimeTableComponent.vue';
import {
  UserWorkTime,
  WorkTimeModel,
} from 'src/models/WorkTimeModel';
import { UserToken } from 'src/models/Terminal';

const { t } = useI18n();

const route = useRoute();
const userId = computed(() => {
  return parseInt(route.params.userId as string);
});
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
const overtimeTotal = ref<number>();
const q = useQuasar();
const missingDays = ref<MissingDay[]>([]);
const userWorkTimes = ref<UserWorkTime[]>([]);
const workTimeModels = ref<WorkTimeModel[]>([]);
const showUserWorkTimeDialog = ref(false);
const newUserWorkTime = ref({
  WorkTimeModelID: null as number | null,
  ValidFrom: '',
  ValidTill: null as string | null,
});
const userTokens = ref<UserToken[]>([]);
const showUserTokenDialog = ref(false);
const newUserToken = ref({
  TokenType: 'chip',
  TokenIdentifier: '',
});

const tokenTypeOptions = [
  {
    value: 'chip',
    label: 'Chip',
  },
];

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

const missingDaysColumn = ref<QTableColumn[]>([
  {
    name: 'date',
    label: `${t('LABEL_DATE')}`,
    field: 'date',
    sortable: true,
    format: (val) => formatDate(val, 'ddd. DD.MM.YYYY'),
    align: 'left',
  },
]);

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
  const result = await BeeTimeClock.administrationTimestampUserMonths(
    userId.value,
  );

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
  BeeTimeClock.administrationAbsencesByYear(
    userId.value,
    selectedAbsenceYear.value,
  )
    .then((result) => {
      if (result.status === 200) {
        absences.value = result.data.Data.map((s) => Absence.fromApi(s));
      }
    })
    .catch((error: ErrorResponse) => {
      showErrorMessage(error.message);
    });
}

const overtimeQuotas = ref<OvertimeMonthQuota[]>([]);

function loadOvertimeQuotas() {
  if (!user.value) return;
  BeeTimeClock.administrationOvertimeMonthQuotas(user.value.ID)
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
  BeeTimeClock.administrationCalculateOvertimeMonthQuota(
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

function loadOvertimeTotal() {
  if (!user.value) return;

  BeeTimeClock.administrationOvertimeTotal(user.value.ID)
    .then((result) => {
      if (result.status === 200) {
        overtimeTotal.value = result.data.Data.Total;
      }
    })
    .catch((error: ErrorResponse) => {
      showErrorMessage(error.response?.data.Message);
    });
}

function loadMissingDays() {
  if (!user.value) return;

  BeeTimeClock.administrationGetMissingDays(userId.value)
    .then((result) => {
      if (result.status === 200) {
        missingDays.value = result.data.Data.map((s) => {
          return { date: s } as MissingDay;
        }).sort(
          (a, b) => new Date(b.date).getTime() - new Date(a.date).getTime(),
        );
      }
    })
    .catch((error) => {
      showErrorMessage(error);
    });
}

function deleteTimestamp(timestamp: Timestamp) {
  BeeTimeClock.administrationTimestampUserDelete(userId.value, timestamp.ID)
    .then((result) => {
      if (result.status === 204) {
        loadTimestampGrouped();
      }
    })
    .catch((error: ErrorResponse) => {
      showErrorMessage(error.message);
    });
}

function loadUserWorkTimes() {
  BeeTimeClock.administrationGetUserWorkTimes(userId.value)
    .then((result) => {
      if (result.status === 200) {
        userWorkTimes.value = result.data.Data.map((s) =>
          UserWorkTime.fromApi(s),
        );
      }
    })
    .catch((error: ErrorResponse) => {
      showErrorMessage(error.message);
    });
}

function loadWorkTimeModels() {
  BeeTimeClock.administrationGetWorkTimeModels()
    .then((result) => {
      if (result.status === 200) {
        workTimeModels.value = result.data.Data.map((s) =>
          WorkTimeModel.fromApi(s),
        );
      }
    })
    .catch((error: ErrorResponse) => {
      showErrorMessage(error.message);
    });
}

function openUserWorkTimeDialog() {
  newUserWorkTime.value = {
    WorkTimeModelID: null,
    ValidFrom: '',
    ValidTill: null,
  };
  showUserWorkTimeDialog.value = true;
}

function createUserWorktime() {
  if (
    !newUserWorkTime.value.WorkTimeModelID ||
    !newUserWorkTime.value.ValidFrom
  ) {
    return;
  }

  BeeTimeClock.administrationCreateUserWorkTime(userId.value, {
    WorkTimeModelID: newUserWorkTime.value.WorkTimeModelID,
    ValidFrom: new Date(newUserWorkTime.value.ValidFrom),
    ValidTill: newUserWorkTime.value.ValidTill
      ? new Date(newUserWorkTime.value.ValidTill)
      : null,
  })
    .then((result) => {
      if (result.status === 201) {
        showUserWorkTimeDialog.value = false;
        loadUserWorkTimes();
        showInfoMessage(
          t('MSG_CREATE_SUCCESS', { item: t('LABEL_USER_WORKTIME') }),
        );
      }
    })
    .catch((error: ErrorResponse) => {
      showErrorMessage(error.message);
    });
}

function deleteUserWorkTime(userWorktime: UserWorkTime) {
  q.dialog({
    message: t('MSG_DELETE', {
      item: t('LABEL_USER_WORKTIME'),
      identifier: userWorktime.WorkTimeModel.Name,
    }),
    cancel: true,
    persistent: true,
  }).onOk(() => {
    BeeTimeClock.administrationDeleteUserWorkTime(userId.value, userWorktime.ID)
      .then((result) => {
        if (result.status === 204) {
          loadUserWorkTimes();
          showInfoMessage(
            t('MSG_DELETE_SUCCESS', {
              item: t('LABEL_USER_WORKTIME'),
              identifier: '',
            }),
          );
        }
      })
      .catch((error: ErrorResponse) => {
        showErrorMessage(error.message);
      });
  });
}

function loadUserTokens() {
  BeeTimeClock.administrationUserTokenList(userId.value)
    .then((result) => {
      if (result.status === 200) {
        userTokens.value = result.data.Data.map((s) => UserToken.fromApi(s));
      }
    })
    .catch((error: ErrorResponse) => {
      showErrorMessage(error.message);
    });
}

function openUserTokenDialog() {
  newUserToken.value = {
    TokenType: 'chip',
    TokenIdentifier: '',
  };
  showUserTokenDialog.value = true;
}

function createUserToken() {
  if (!newUserToken.value.TokenIdentifier) {
    return;
  }

  BeeTimeClock.administrationUserTokenCreate(userId.value, {
    TokenType: newUserToken.value.TokenType,
    TokenIdentifier: newUserToken.value.TokenIdentifier,
  })
    .then((result) => {
      if (result.status === 201) {
        showUserTokenDialog.value = false;
        loadUserTokens();
        showInfoMessage(t('MSG_CREATE_SUCCESS', { item: t('LABEL_TOKEN') }));
      }
    })
    .catch((error: ErrorResponse) => {
      showErrorMessage(error.message);
    });
}

function deleteUserToken(userToken: UserToken) {
  q.dialog({
    message: t('MSG_DELETE', {
      item: t('LABEL_TOKEN'),
      identifier: userToken.TokenIdentifier,
    }),
    cancel: true,
    persistent: true,
  }).onOk(() => {
    BeeTimeClock.administrationUserTokenDelete(userId.value, userToken.ID)
      .then((result) => {
        if (result.status === 204) {
          loadUserTokens();
          showInfoMessage(
            t('MSG_DELETE_SUCCESS', {
              item: t('LABEL_TOKEN'),
              identifier: '',
            }),
          );
        }
      })
      .catch((error: ErrorResponse) => {
        showErrorMessage(error.message);
      });
  });
}

const userTokenColumns = ref<QTableColumn[]>([
  {
    name: 'tokenType',
    label: t('LABEL_TOKEN_TYPE'),
    field: 'TokenType',
    align: 'left',
  },
  {
    name: 'tokenIdentifier',
    label: t('LABEL_TOKEN_IDENTIFIER'),
    field: 'TokenIdentifier',
    align: 'left',
  },
]);

const userWorktimeColumns = ref<QTableColumn[]>([
  {
    name: 'workTimeModel',
    label: t('LABEL_WORKTIME_MODEL'),
    field: (row: UserWorkTime) => row.WorkTimeModel?.Name,
    align: 'left',
  },
  {
    name: 'validFrom',
    label: t('LABEL_VALID_FROM'),
    field: 'ValidFrom',
    format: (val: Date) => formatDate(val, 'DD.MM.YYYY'),
    align: 'left',
  },
  {
    name: 'validTill',
    label: t('LABEL_VALID_TILL'),
    field: 'ValidTill',
    format: (val: Date | null) =>
      val ? formatDate(val, 'DD.MM.YYYY') : t('LABEL_NO_END_DATE'),
    align: 'left',
  },
]);

onMounted(async () => {
  loadUser();
  loadAbsenceYears();
  loadAbsences();
  await loadTimestampMonths();
  loadTimestampGrouped();
  loadOvertimeTotal();
  loadOvertimeQuotas();
  loadMissingDays();
  loadUserWorkTimes();
  loadWorkTimeModels();
  loadUserTokens();
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

const columns = [
  {
    name: 'absenceReason',
    label: t('LABEL_REASON'),
    field: 'Reason',
    align: 'left',
  },
  {
    name: 'absenceFrom',
    label: t('LABEL_FROM'),
    field: 'AbsenceFrom',
    format: (val: Date) => date.formatDate(val, 'ddd DD. MMM. YYYY'),
    sortable: true,
  },
  {
    name: 'absenceTill',
    label: t('LABEL_TILL'),
    field: 'AbsenceTill',
    format: (val: Date) => date.formatDate(val, 'ddd DD. MMM. YYYY'),
    sortable: true,
  },
  {
    name: 'absenceNettoDays',
    label: t('LABEL_NETTO_DAYS'),
    field: 'NettoDays',
  },
  {
    name: 'signed',
    label: t('LABEL_STATUS'),
    field: 'SignedStatus',
  },
  {
    name: 'signedBy',
    label: t('LABEL_SIGNED_BY'),
    field: 'signedUserMapped',
    format: (val: User | null) => val?.displayName,
  },
  {
    name: 'createdAt',
    label: t('LABEL_CREATED_AT'),
    field: 'CreatedAt',
    format: (val: string) => date.formatDate(val, 'ddd DD. MMM. YYYY'),
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
        <q-tab name="absence" icon="event" :label="t('LABEL_ABSENCE', 2)" />
        <q-tab
          name="missing-days"
          icon="event_busy"
          :label="t('LABEL_MISSING_DAY', 2)"
        />
        <q-tab name="tokens" icon="key" :label="t('LABEL_TOKEN', 2)" />
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
              <q-input
                :label="t('LABEL_STAFF_NUMBER')"
                v-model.number="user.StaffNumber"
                type="number"
              />
            </q-card-section>
            <q-card-actions>
              <q-btn :label="t('BTN_SAVE')" color="primary" @click="saveUser" />
            </q-card-actions>
          </q-card>
          <q-card class="q-mt-md">
            <q-card-section class="bg-primary text-white text-h6">
              <div class="row">
                {{ t('LABEL_WORKTIME_ASSIGNMENTS') }}
                <q-space />
                <q-btn
                  icon="add"
                  color="secondary"
                  @click="openUserWorkTimeDialog"
                />
              </div>
            </q-card-section>
            <q-card-section>
              <q-table
                :rows="userWorkTimes"
                :columns="userWorktimeColumns"
                flat
              >
                <template v-slot:body="props">
                  <q-tr :props="props">
                    <q-td
                      v-for="col in props.cols"
                      :key="col.name"
                      :props="props"
                      :align="col.align || 'left'"
                    >
                      {{ col.value }}
                    </q-td>
                    <q-td auto-width>
                      <q-btn
                        icon="delete"
                        color="negative"
                        @click="deleteUserWorkTime(props.row)"
                      />
                    </q-td>
                  </q-tr>
                </template>
              </q-table>
            </q-card-section>
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
                dense
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
              disable-edit
            />
          </div>
        </q-tab-panel>
        <q-tab-panel name="overtime">
          <q-card v-if="overtimeTotal" class="q-mb-lg">
            <q-card-section class="bg-primary text-h6 text-white">
              {{ t('LABEL_OVERTIME_TOTAL') }}
            </q-card-section>
            <q-card-section class="text-h6 text-center">
              {{ formatIndustryHourMinutes(overtimeTotal) }}
            </q-card-section>
          </q-card>
          <OvertimeTableComponent
            v-model="overtimeQuotas"
            @calculate="calculateOvertimeMonth"
          />
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
        <q-tab-panel name="missing-days">
          <q-btn class="full-width" label="refresh" @click="loadMissingDays" />
          <q-table
            :columns="missingDaysColumn"
            :rows="missingDays"
            hide-pagination
            :pagination="emptyPagination"
          >
            <template v-slot:header="props">
              <q-tr :props="props">
                <q-th v-for="col in props.cols" :key="col.name" :props="props">
                  {{ col.label }}
                </q-th>
              </q-tr>
            </template>
            <template v-slot:body="props">
              <q-tr :props="props">
                <q-td v-for="col in props.cols" :key="col.name" :props="props">
                  {{ col.value }}
                </q-td>
              </q-tr>
            </template>
          </q-table>
        </q-tab-panel>
        <q-tab-panel name="tokens">
          <q-card>
            <q-card-section class="bg-primary text-white text-h6">
              <div class="row items-center">
                {{ t('LABEL_TOKEN', 2) }}
                <q-space />
                <q-btn
                  icon="add"
                  color="secondary"
                  @click="openUserTokenDialog"
                />
              </div>
            </q-card-section>
            <q-card-section>
              <q-table
                :rows="userTokens"
                :columns="userTokenColumns"
                row-key="ID"
                flat
              >
                <template v-slot:body="props">
                  <q-tr :props="props">
                    <q-td
                      v-for="col in props.cols"
                      :key="col.name"
                      :props="props"
                      :align="col.align || 'left'"
                    >
                      {{ col.value }}
                    </q-td>
                    <q-td auto-width>
                      <q-btn
                        icon="delete"
                        color="negative"
                        flat
                        dense
                        @click="deleteUserToken(props.row)"
                      />
                    </q-td>
                  </q-tr>
                </template>
              </q-table>
            </q-card-section>
          </q-card>
        </q-tab-panel>
      </q-tab-panels>
    </div>
    <q-dialog v-model="showUserWorkTimeDialog">
      <q-card>
        <q-card-section class="bg-primary text-white text-h6">
          {{ t('LABEL_CREATE', { item: t('LABEL_USER_WORKTIME') }) }}
        </q-card-section>
        <q-form @submit="createUserWorktime">
          <q-card-section>
            <q-select
              v-model="newUserWorkTime.WorkTimeModelID"
              :options="
                workTimeModels.map((s) => ({
                  label: s.Name,
                  value: s.ID,
                }))
              "
              :label="t('LABEL_WORKTIME_MODEL')"
              class="full-width"
              map-options
              emit-value
            />
            <q-input
              v-model="newUserWorkTime.ValidFrom"
              :label="t('LABEL_VALID_FROM')"
              type="date"
            />
            <q-input
              v-model="newUserWorkTime.ValidTill"
              :label="t('LABEL_VALID_TILL')"
              type="date"
            />
          </q-card-section>
          <q-card-actions align="right">
            <q-btn
              flat
              :label="t('BTN_CANCEL')"
              color="primary"
              v-close-popup
              type="reset"
            />
            <q-btn
              flat
              :label="t('BTN_CREATE')"
              color="primary"
              type="submit"
            />
          </q-card-actions>
        </q-form>
      </q-card>
    </q-dialog>
    <q-dialog v-model="showUserTokenDialog">
      <q-card>
        <q-card-section class="bg-primary text-white text-h6">
          {{ t('LABEL_CREATE', { item: t('LABEL_TOKEN') }) }}
        </q-card-section>
        <q-form @submit="createUserToken">
          <q-card-section>
            <q-select
              v-model="newUserToken.TokenType"
              :options="tokenTypeOptions"
              :label="t('LABEL_TOKEN_TYPE')"
              class="full-width"
              map-options
              emit-value
            />
            <q-input
              v-model="newUserToken.TokenIdentifier"
              :label="t('LABEL_TOKEN_IDENTIFIER')"
              :rules="[(val) => !!val || t('LABEL_TOKEN_IDENTIFIER')]"
            />
          </q-card-section>
          <q-card-actions align="right">
            <q-btn
              flat
              :label="t('BTN_CANCEL')"
              color="primary"
              v-close-popup
              type="reset"
            />
            <q-btn
              flat
              :label="t('BTN_CREATE')"
              color="primary"
              type="submit"
            />
          </q-card-actions>
        </q-form>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<style scoped></style>
