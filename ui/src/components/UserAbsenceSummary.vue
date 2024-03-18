<script setup lang="ts">
import { computed, onMounted, PropType, ref } from 'vue';
import BeeTimeClock from 'src/service/BeeTimeClock';
import { AbsenceReason, AbsenceUserSummary } from 'src/models/Absence';

const props = defineProps({
  userId: {
    type: Number,
    required: true
  },
  absenceReasons: {
    type: Array as PropType<AbsenceReason[]>,
    required: true,
  }
});

const userSummary = ref<AbsenceUserSummary | null>(null);

const byAbsenceReason = computed(() => {
  if (!userSummary.value) return null;
  return Object.keys(userSummary.value.ByYear).map((year) => {
    return userSummary.value?.ByYear[parseInt(year)].ByAbsenceReason
  })[0];
})

function loadUserSummary() {
  BeeTimeClock.administrationSummaryUserCurrentYear(props.userId, new Date().getFullYear()).then(response => {
    if (response.status === 200) {
      userSummary.value = response.data.Data;
    }
  });
}

function getAbsenceReason(absenceReasonId: number) : string|undefined {
  if (!props.absenceReasons) return undefined;
  return props.absenceReasons.find(x => x.ID == absenceReasonId)?.Description;
}

onMounted(() => {
  loadUserSummary();
});
</script>

<template>
  <q-markup-table v-if="userSummary">
    <thead>
    <tr>
      <th class="text-left">Typ</th>
      <th class="text-left">genommen</th>
      <th class="text-left">geplant</th>
    </tr>
    </thead>
    <tbody>
          <tr v-for="(days, reasonId) in byAbsenceReason" :key="reasonId">
            <td>{{ getAbsenceReason(reasonId) }}</td>
            <td>{{ days.Past }}</td>
            <td>{{ days.Upcoming }}</td>
          </tr>
    </tbody>
  </q-markup-table>

</template>

<style scoped>

</style>
