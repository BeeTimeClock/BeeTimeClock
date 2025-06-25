<script setup lang="ts">
import {
  ExternalWorkExpanse,
} from 'src/models/ExternalWork';
import { computed } from 'vue';
import { formatCurrency } from 'src/helper/formatter';

const externalWorkExpanse = defineModel({
  type: ExternalWorkExpanse,
  required: true,
});

const group = computed({
  get() {
    if (externalWorkExpanse.value.AdditionalOptions == undefined) {
     return [];
    }

    return externalWorkExpanse.value.AdditionalOptions;
  },
  set(val: string[]) {
    externalWorkExpanse.value.AdditionalOptions = val;
  }
})
const editMode = defineModel('editmode', {
  type: Boolean,
  default: false,
});

const options = computed(() => {
  return Object.keys(externalWorkExpanse.value.ExternalWork.ExternalWorkCompensation.AdditionalOptions).map(s => {
    return {
      label: `${s} (${formatCurrency(externalWorkExpanse.value.ExternalWork.ExternalWorkCompensation.AdditionalOptions[s])})`,
      value: s
    }
  })
})
</script>

<template>
  <q-card v-if="externalWorkExpanse" class="q-mb-lg">
    <q-card-section class="bg-secondary">{{
      $t('LABEL_OPTION', 2)
    }}</q-card-section>
    <q-card-section>
      <q-option-group
        v-model="group"
        :options="options"
        type="toggle"
        :disable="!editMode"
      />
    </q-card-section>
  </q-card>
</template>

<style scoped></style>
