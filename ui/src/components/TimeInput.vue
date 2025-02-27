<script setup lang="ts">

import { computed, ref } from 'vue';

const label = defineModel('label', {type: String})
const readonly = defineModel('readonly', {type: Boolean, default: false})
const value = defineModel<Date|null>();
const date = defineModel('date', {type:String, required: true})

const timeIntern = ref();

const time = computed({
  get() {
    return timeIntern.value;
  },
  set(val: string|undefined|null) {
    if (val == '--:--') {
      timeIntern.value = null;
    } else {
      timeIntern.value = val;
    }

    if (isValid.value) {
      if (!timeIntern.value) {
        value.value = null;
        return;
      }
      const currentDate = new Date(date.value);
      console.log(currentDate.toISOString());
      const parts = timeIntern.value.split(':')
      const hours =parseInt(parts[0]);
      const minutes = parseInt(parts[1]);
      console.log(parts, hours, minutes)
      currentDate.setHours(hours, minutes)
      console.log(currentDate.toISOString());
      value.value = currentDate;
    }
  }
})
const isValid = computed(() => {
  if (!time.value || time.value == undefined || time.value == '') return true;
  console.log('current val: ', time.value)
  const timeRegex = /^([0-9]|0[0-9]|1[0-9]|2[0-3]):[0-5][0-9]$/;
  return timeRegex.test(time.value);
})
</script>

<template>
  <q-input
    mask="##:##"
    fill-mask="-"
    v-model="time"
    :label="label"
    :readonly="readonly"
    :clearable="!readonly"
    :error="!isValid"
    :error-message="$t('RULE_VALID_HOUR')"
  />
</template>

<style scoped>

</style>
