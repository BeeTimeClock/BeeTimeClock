<script setup lang="ts">
import { ExternalWork, ExternalWorkExpanse } from 'src/models/ExternalWork';
import { computed, onMounted, ref } from 'vue';
import BeeTimeClock from 'src/service/BeeTimeClock';
import { useRoute, useRouter } from 'vue-router';
import ExternalWorkExpanseItem from 'components/external_work/ExternalWorkExpanseItem.vue';
import { formatCurrency } from 'src/helper/formatter';
import { showErrorMessage, showInfoMessage } from 'src/helper/message';
import { useI18n } from 'vue-i18n';
import type { AxiosError } from 'axios';
import type { BaseResponse, ErrorResponse } from 'src/models/Base';
import { date, useQuasar } from 'quasar';

const route = useRoute();
const router = useRouter();
const externalWork = ref<ExternalWork>();
const loading = ref(true);
const expanses = ref<ExternalWorkExpanse[]>();
const { t } = useI18n();
const q = useQuasar();

const externalWorkId = computed(() => {
  return parseInt(route.params.externalWorkId as string);
});

const externalWorkExpansesSorted = computed(() => {
  if (expanses.value == null) return [];
  const cloned = Object.assign([] as ExternalWorkExpanse[], expanses.value);
  return cloned.sort(
    (a, b) => new Date(a.Date).getTime() - new Date(b.Date).getTime(),
  );
});

function loadExternalWork() {
  BeeTimeClock.getExternalWorkById(externalWorkId.value)
    .then((result) => {
      if (result.status === 200) {
        externalWork.value = ExternalWork.fromApi(result.data.Data);
        expanses.value = result.data.Data.WorkExpanses.map((s) =>
          ExternalWorkExpanse.fromApi(s),
        );
      }
    })
    .catch((error: ErrorResponse) => {
      showErrorMessage(error.message);
    })
    .finally(() => {
      loading.value = false;
    });
}

function submit() {
  BeeTimeClock.submitExternalWork(externalWorkId.value)
    .then((result) => {
      if (result.status === 200) {
        showInfoMessage(t('MSG_SUBMIT'));
        loadExternalWork();
      }
    })
    .catch((error: ErrorResponse) => {
      showErrorMessage(error.response?.data.Message);
    });
}

function deleteExternalWork() {
  if (externalWork.value == null) return;

  q.dialog({
    title: t('LABEL_DELETE'),
    message: t('MSG_DELETE', {
      item: t('LABEL_EXTERNAL_WORK'),
      identifier: externalWork.value.Description,
    }),
    cancel: true,
    persistent: true,
  }).onOk(() => {
    BeeTimeClock.deleteExternalWorkById(externalWork.value!.ID)
      .then((result) => {
        if (result.status === 204) {
          showInfoMessage(
            t('MSG_DELETE_SUCCESS', {
              item: t('LABEL_EXTERNAL_WORK'),
              identifier: externalWork.value!.Description,
            }),
          );
          router
            .push({ name: 'ExternalWorkOverview' })
            .then(() => {})
            .catch((error: ErrorResponse) => {
              showErrorMessage(error.response?.data.Message);
            });
        }
      })
      .catch((error: AxiosError<BaseResponse<never>>) => {
        showErrorMessage(error.response?.data.Message);
      });
  });
}

onMounted(() => {
  loadExternalWork();
});
</script>

<template>
  <q-page padding>
    <div v-if="externalWork && !loading">
      <q-card>
        <q-card-section class="bg-primary text-h6 text-white">
          {{ $t('LABEL_INFORMATION') }}
          <div class="float-right">
            <q-btn
              v-if="externalWork.NeedsUserInput"
              :label="$t('LABEL_DELETE')"
              color="negative"
              class="q-mr-md"
              @click="deleteExternalWork"
            />
            <q-btn
              v-if="externalWork.NeedsUserInput"
              :label="$t('LABEL_SUBMIT')"
              color="secondary"
              @click="submit"
            />
          </div>
        </q-card-section>
        <q-card-section>
          <q-list>
            <q-item>
              <q-item-section>
                <q-item-label caption
                  >{{ $t('LABEL_DESCRIPTION') }}
                </q-item-label>
                <q-item-label>{{ externalWork.Description }}</q-item-label>
              </q-item-section>
              <q-item-section side>
                <q-chip color="secondary" square>{{
                  $t(externalWork.StatusLabel)
                }}</q-chip>
              </q-item-section>
            </q-item>
            <q-item>
              <q-item-section>
                <q-item-label caption>{{ $t('LABEL_FROM') }}</q-item-label>
                <q-item-label>{{
                  date.formatDate(externalWork.From, 'DD.MM.YYYY')
                }}</q-item-label>
              </q-item-section>
              <q-item-section>
                <q-item-label caption>{{ $t('LABEL_TILL') }}</q-item-label>
                <q-item-label>{{
                  date.formatDate(externalWork.Till, 'DD.MM.YYYY')
                }}</q-item-label>
              </q-item-section>
            </q-item>
            <q-item-label header>{{ $t('LABEL_CALCULATED') }}</q-item-label>
            <q-item>
              <q-item-section>
                <q-item-label caption>{{
                  $t('LABEL_TOTAL_OVERTIME_HOURS')
                }}</q-item-label>
                <q-item-label
                  >{{ externalWork.TotalOvertimeHours }}
                  {{ $t('LABEL_HOUR', externalWork.TotalOvertimeHours) }}
                </q-item-label>
              </q-item-section>
            </q-item>
            <q-item>
              <q-item-section>
                <q-item-label caption>{{
                  $t('LABEL_TOTAL_EXPENSE_WITH_SOCIAL_INSURANCE')
                }}</q-item-label>
                <q-item-label>{{
                  formatCurrency(externalWork.TotalExpensesWithSocialInsurance)
                }}</q-item-label>
              </q-item-section>
              <q-item-section>
                <q-item-label caption>{{
                  $t('LABEL_TOTAL_EXPENSE_WITHOUT_SOCIAL_INSURANCE')
                }}</q-item-label>
                <q-item-label>{{
                  formatCurrency(
                    externalWork.TotalExpensesWithoutSocialInsurance,
                  )
                }}</q-item-label>
              </q-item-section>
              <q-item-section>
                <q-item-label caption>{{
                  $t('LABEL_TOTAL_ADDITION_HAIRCUT')
                }}</q-item-label>
                <q-item-label
                  v-for="(value, key) in externalWork.TotalOptions"
                  :key="key"
                  >{{ key }}: {{ formatCurrency(value) }}</q-item-label
                >
              </q-item-section>
            </q-item>
          </q-list>
        </q-card-section>
      </q-card>
      <div class="q-mt-lg">
        <ExternalWorkExpanseItem
          v-for="(expanse, index) in externalWorkExpansesSorted"
          :key="index"
          v-model:externalworkexpanse="externalWorkExpansesSorted[index]!"
          @updated="loadExternalWork"
          :allow_edit="!externalWork.IsLocked"
        />
      </div>
    </div>
    <q-inner-loading :showing="loading" />
  </q-page>
</template>

<style scoped></style>
