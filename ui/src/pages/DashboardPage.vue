<script setup lang="ts">
import {onMounted, ref} from 'vue';
import BeeTimeClock from 'src/service/BeeTimeClock';
import {User} from 'src/models/Authentication';
import {AbsenceSummaryItem} from "src/models/Absence";
import AbsenceSummaryTableComponent from "components/AbsenceSummaryTableComponent.vue";

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
      <AbsenceSummaryTableComponent v-model="absenceSummaryItems"/>
    </div>
  </q-page>
</template>

<style scoped>

</style>
