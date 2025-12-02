<script setup lang="ts">
import { computed, onMounted, ref } from 'vue';
import {
  User,
  UserWithAbsenceSummaryAndOvertime,
} from 'src/models/Authentication';
import BeeTimeClock from 'src/service/BeeTimeClock';
import type {
  AbsenceReason,
  AbsenceUserSummaryYearReason,
} from 'src/models/Absence';
import { type QTableColumn, useQuasar } from 'quasar';
import { useI18n } from 'vue-i18n';
import { showErrorMessage, showInfoMessage } from 'src/helper/message';
import type { ErrorResponse } from 'src/models/Base';
import { emptyPagination } from 'src/helper/objects';
import { formatIndustryHourMinutes } from 'src/helper/formatter';

const q = useQuasar();
const { t } = useI18n();

const users = ref([] as UserWithAbsenceSummaryAndOvertime[]);
const absenceReasons = ref([] as AbsenceReason[]);
const needle = ref('');

const selectedAbsenceReasons = ref<number[]>([]);

function loadUsers() {
  BeeTimeClock.administrationGetUsers(true)
    .then((result) => {
      if (result.status === 200) {
        users.value = result.data.Data.map(
          (s) => new UserWithAbsenceSummaryAndOvertime(User.fromApi(s)),
        );

        users.value.forEach((s) => {
          const userIndex = users.value.indexOf(s);

          void BeeTimeClock.administrationSummaryUserCurrentYear(
            s.ID,
            new Date().getFullYear(),
          ).then((absenceResult) => {
            users.value[userIndex]!.absenceSummary = absenceResult.data.Data;
          });

          void BeeTimeClock.administrationOvertimeTotal(s.ID).then(
            (overtimeResult) => {
              users.value[userIndex]!.overtime = overtimeResult.data.Data.Total;
            },
          );
        });
      }
    })
    .catch((error: ErrorResponse) => {
      showErrorMessage(error.message);
    });
}

function loadAbsenceReasons() {
  BeeTimeClock.absenceReasons()
    .then((result) => {
      absenceReasons.value = result.data.Data;
    })
    .catch((error: ErrorResponse) => {
      showErrorMessage(error.message);
    });
}

const columns = computed(() => {
  const result = [
    {
      name: 'username',
      label: t('LABEL_USERNAME'),
      field: 'Username',
      align: 'left',
    },
    {
      name: 'firstName',
      label: t('LABEL_FIRST_NAME'),
      field: 'FirstName',
      align: 'left',
    },
    {
      name: 'lastName',
      label: t('LABEL_LAST_NAME'),
      field: 'LastName',
      align: 'left',
    },
    {
      name: 'overtime',
      label: t('LABEL_OVERTIME'),
      field: 'overtime',
      format: (val: number) => formatIndustryHourMinutes(val ?? 0),
    },
  ] as QTableColumn[];

  result.push(
    ...absenceReasons.value
      .filter((s) => selectedAbsenceReasons.value.includes(s.ID))
      .map((s) => {
        return {
          name: s.ID.toString(),
          label: s.Description,
          field: (row: UserWithAbsenceSummaryAndOvertime) => {
            const year = row.absenceSummary?.ByYear[2025];
            if (year == undefined) return null;
            return year.ByAbsenceReason[s.ID] ?? null;
          },
          align: 'right',
          format: (val: AbsenceUserSummaryYearReason | null) =>
            `${val?.Past ?? 0} ${(val?.Upcoming ?? 0) > 0 ? `/ ${val?.Upcoming ?? 0}` : ''}`,
        } as QTableColumn;
      }),
  );
  return result;
});

const sortedFilteredUsers = computed(() => {
  const search = needle.value.toLowerCase();
  const data = users.value.filter((s) => {
    if (s.LastName.toLowerCase().indexOf(search) >= 0) return true;
    if (s.FirstName.toLowerCase().indexOf(search) >= 0) return true;
    if (s.Username.toLowerCase().indexOf(search) >= 0) return true;
  });
  return data.sort((a, b) => a.LastName.localeCompare(b.LastName));
});

function deleteUser(user: User) {
  q.dialog({
    message: t('MSG_DELETE', {
      item: t('LABEL_USER'),
      identifier: user.Username,
    }),
    cancel: true,
    persistent: true,
  }).onOk(() => {
    BeeTimeClock.administrationDeleteUser(user)
      .then((result) => {
        if (result.status === 204) {
          loadUsers();
          showInfoMessage(
            t('MSG_DELETE_SUCCESS', {
              item: t('LABEL_USER'),
              identifier: user.Username,
            }),
          );
        }
      })
      .catch((error: ErrorResponse) => {
        showErrorMessage(error.message);
      });
  });
}

onMounted(() => {
  loadAbsenceReasons();
  loadUsers();
});
</script>

<template>
  <q-page padding>
    <q-select
      filled
      v-model="selectedAbsenceReasons"
      :options="absenceReasons"
      :label="t('LABEL_ABSENCE_REASON', 2)"
      multiple
      emit-value
      map-options
      option-value="ID"
      option-label="Description"
      use-chips
    >
      <template v-slot:option="{ itemProps, opt, selected, toggleOption }">
        <q-item v-bind="itemProps">
          <q-item-section>
            <q-item-label>{{ opt.Description }}</q-item-label>
          </q-item-section>
          <q-item-section side>
            <q-toggle
              :model-value="selected"
              @update:model-value="toggleOption(opt)"
            />
          </q-item-section>
        </q-item>
      </template>
    </q-select>
    <q-input
      :label="t('LABEL_SEARCH')"
      v-model="needle"
      clearable
      @clear="needle = ''"
    />
    <q-table
      class="q-mt-md"
      :rows="sortedFilteredUsers"
      :pagination="emptyPagination"
      :columns="columns"
      hide-pagination
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
        <q-tr :props="props">
          <q-td v-for="col in props.cols" :key="col.name" :props="props">
            {{ col.value }}
          </q-td>
          <q-td auto-width>
            <q-btn
              class="q-ml-md"
              color="primary"
              icon="visibility"
              :to="{
                name: 'AdministrationUserDetail',
                params: { userId: props.row.ID },
              }"
            />
            <q-btn
              class="q-ml-md"
              color="negative"
              icon="delete"
              @click="deleteUser(props.row)"
            />
          </q-td>
        </q-tr>
      </template>
    </q-table>
  </q-page>
</template>

<style scoped></style>
