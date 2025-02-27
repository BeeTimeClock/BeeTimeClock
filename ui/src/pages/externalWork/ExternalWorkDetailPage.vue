<script setup lang="ts">
import { ExternalWork, ExternalWorkExpanse } from 'src/models/ExternalWork';
import { computed, onMounted, ref } from 'vue';
import BeeTimeClock from 'src/service/BeeTimeClock';
import { useRoute } from 'vue-router';
import ExternalWorkExpanseItem from 'components/external_work/ExternalWorkExpanseItem.vue';

const route = useRoute();
const externalWork = ref<ExternalWork>();
const loading = ref(true);
const expanses = ref<ExternalWorkExpanse[]>();

const externalWorkId = computed(() => {
  return parseInt(route.params.externalWorkId as string);
});

function loadExternalWork() {
  BeeTimeClock.getExternalWorkById(externalWorkId.value)
    .then((result) => {
      if (result.status === 200) {
        externalWork.value = ExternalWork.fromApi(result.data.Data);
        expanses.value = result.data.Data.WorkExpanses.map((s) =>
          ExternalWorkExpanse.fromApi(s)
        );
      }
    })
    .finally(() => {
      loading.value = false;
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
            </q-item>
            <q-item>
              <q-item-section>
                <q-item-label caption>{{ $t('LABEL_FROM') }}</q-item-label>
                <q-item-label>{{ externalWork.From }}</q-item-label>
              </q-item-section>
              <q-item-section>
                <q-item-label caption>{{ $t('LABEL_TILL') }}</q-item-label>
                <q-item-label>{{ externalWork.Till }}</q-item-label>
              </q-item-section>
            </q-item>
          </q-list>
        </q-card-section>
      </q-card>
      <div class="q-mt-lg">
        <ExternalWorkExpanseItem
          v-for="(expanse, index) in expanses"
          :key="index"
          v-model:externalworkexpanse="expanses[index]"
        />
      </div>
    </div>
    <q-inner-loading :showin="loading" />
  </q-page>
</template>

<style scoped></style>
