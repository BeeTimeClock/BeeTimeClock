<script setup lang="ts">

import { computed } from 'vue';
import { date } from 'quasar';

const label = defineModel('label', {type: String})
const readonly = defineModel('readonly', {type: Boolean, default: false})
const value = defineModel<Date|null>();
const dayDate = defineModel('date', {type:String, required: true})

const dateFormat = 'HH:mm';
const time = computed({
  get() {
    if (!value.value) return '';

    return date.formatDate(
      new Date(value.value),
      dateFormat
    );
  },
  set(newValue) {
    if (newValue == '') {
      value.value = undefined;
      return;
    }

    const day = date.formatDate(
      new Date(
        value.value ??
        dayDate.value
      ),
      'YYYY-MM-DD'
    );

    const full = `${day} ${newValue}`;

    value.value = new Date(full);
  },
});
</script>

<template>
  <q-input
    type="time"
    v-model="time"
    :label="label"
    :readonly="readonly"
    :clearable="!readonly"
  />
</template>

<style scoped>

</style>
