<script setup lang="ts">
import { onMounted, ref } from 'vue';
import { Holiday } from 'src/models/Holiday';
import BeeTimeClock from 'src/service/BeeTimeClock';
import { showErrorMessage } from 'src/helper/message';

const tab = ref('holiday')
const holidays = ref<Holiday[]>([])

function loadHolidays() {
  BeeTimeClock.administrationDebugHolidays().then(result => {
    if (result.status === 200) {
      holidays.value = result.data.Data.map(s => Holiday.fromApi(s))
    }
  }).catch((error) => {
    showErrorMessage(error)
  })
}

onMounted(() => {
  loadHolidays()
})
</script>

<template>
  <q-page>
    <q-tabs v-model="tab">
      <q-tab name="holiday"/>
    </q-tabs>
    <q-tab-panels v-model="tab">
      <q-tab-panel name="holiday">
        <q-table :rows="holidays"/>
      </q-tab-panel>
    </q-tab-panels>
  </q-page>
</template>

<style scoped></style>
