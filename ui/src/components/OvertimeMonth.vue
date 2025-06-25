<script setup lang="ts">
import { onMounted, ref, watch } from 'vue';
import BeeTimeClock from 'src/service/BeeTimeClock';
import { formatIndustryHourMinutes } from 'src/helper/formatter';
import { OvertimeResponse } from 'src/models/Timestamp';

const overtimeResponse = ref<OvertimeResponse | null>(null);

const props = defineProps({
  modelYear: {
    type: Number,
    default: new Date().getFullYear()
  },
  modelMonth: {
    type: Number,
    default: new Date().getMonth() + 1
  },
  modelUserId: {
    type: String,
  }
});

function loadOvertime() {
  if (props.modelUserId && props.modelUserId != '') {
    BeeTimeClock.administrationTimestampQueryMonthOvertime(props.modelUserId, props.modelYear, props.modelMonth).then(result => {
      if (result.status === 200) {
        overtimeResponse.value = result.data.Data;
      }
    });
  } else {
    BeeTimeClock.timestampQueryMonthOvertime(props.modelYear, props.modelMonth).then(result => {
      if (result.status === 200) {
        overtimeResponse.value = result.data.Data;
      }
    });
  }
}


watch(props, () => {
  loadOvertime();
});

onMounted(() => {
  loadOvertime();
});
</script>

<template>
  <q-card>
    <q-card-section class="bg-primary text-white text-subtitle">
      {{ $t('LABEL_OVERTIME_MONTH', {
      year: modelYear,
      month: new Date(modelYear, modelMonth - 1).toLocaleString('default', { month: 'long' })
    }) }}
    </q-card-section>
    <q-card-section v-if="overtimeResponse" class="text-h6 text-center">
      <div class="row">
        Geleistet: {{ formatIndustryHourMinutes(overtimeResponse.Total) }}
      </div>
      <div class="row">
        Abzug: {{ formatIndustryHourMinutes(overtimeResponse.Subtracted) }}
      </div>
    </q-card-section>
  </q-card>
</template>

<style scoped>

</style>
