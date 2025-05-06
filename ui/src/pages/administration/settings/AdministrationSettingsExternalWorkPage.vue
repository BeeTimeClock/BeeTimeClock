<script setup lang="ts">
import {onMounted, ref} from 'vue';
import BeeTimeClock from 'src/service/BeeTimeClock';
import { ExternalWorkCompensation } from 'src/models/ExternalWork';

const externalWorkCompensations = ref<ExternalWorkCompensation[]>([]);
const isLoading = ref(true);

function loadExternalWorkCompensation() {
  isLoading.value = true;
  BeeTimeClock.administrationExternalWorkCompensation().then(result => {
    if (result.status === 200) {
      externalWorkCompensations.value = result.data.Data.map(s => ExternalWorkCompensation.fromApi(s))
    }
  }).finally(() => {
    isLoading.value = false;
  })
}

onMounted(() => {
  loadExternalWorkCompensation();
})
</script>

<template>
  <q-page padding>
    <div v-if="!isLoading">
      <q-table :rows="externalWorkCompensations"/>
    </div>
    <q-inner-loading :showing="isLoading"/>
  </q-page>
</template>

<style scoped>

</style>
