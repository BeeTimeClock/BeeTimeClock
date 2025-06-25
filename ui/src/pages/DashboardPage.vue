<script setup lang="ts">
import {onMounted, ref} from 'vue';
import BeeTimeClock from 'src/service/BeeTimeClock';
import {User} from 'src/models/Authentication';
import {AbsenceSummaryItem} from 'src/models/Absence';
import AbsenceSummaryTableComponent from 'components/AbsenceSummaryTableComponent.vue';
import OvertimeTotal from 'components/OvertimeTotal.vue';
import OvertimeCurrentMonth from 'components/OvertimeMonth.vue';

const user = ref(null as User | null);
const absenceSummaryItems = ref([] as AbsenceSummaryItem[])

function loadUser() {
  BeeTimeClock.getMeUser().then(result => {
    if (result.status === 200) {
      user.value = result.data.Data;
    }
  })
}

function loadAbsenceSummary() {
  BeeTimeClock.queryAbsenceSummary().then(result => {
    if (result.status === 200) {
      absenceSummaryItems.value = result.data.Data;
    }
  })
}


onMounted(() => {
  loadUser();
  loadAbsenceSummary();
})
</script>

<template>
  <q-page padding>
    <div v-if="user">
      <h2>Hallo {{ user.FirstName }} {{ user.LastName }}</h2>
    </div>
    <div>
      <OvertimeCurrentMonth class="q-mt-sm"/>
      <OvertimeTotal class="q-mt-sm"/>
      <AbsenceSummaryTableComponent v-model="absenceSummaryItems" class="q-mt-sm"/>
    </div>
  </q-page>
</template>

<style scoped>

</style>
