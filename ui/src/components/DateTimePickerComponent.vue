<script setup lang="ts">
import {computed, ref, watch} from 'vue';
import {date} from 'quasar';

const modelValue = defineModel<Date|undefined>({required: true});
const label = defineModel('label', {type: String});
const emit = defineEmits(['update:modelValue']);

const dateFormat = 'DD.MM.YYYY HH:mm';
const value = computed({
  get(): Date|undefined {
    return modelValue.value;
  },
  set(value: Date) {
    emit('update:modelValue', value);
  }
});

const pickerOpen = ref(false)

const formattedDate = ref(date.formatDate(value.value, dateFormat));

watch(formattedDate, (newDateValue) => {
  value.value = date.extractDate(newDateValue, dateFormat);
})
</script>

<template>
  <q-input filled v-model="formattedDate" :label="label" mask="">
    <template v-slot:append>
      <q-icon name="event" class="cursor-pointer" @click="pickerOpen = true">
      </q-icon>
    </template>
  </q-input>

  <q-dialog v-model="pickerOpen" id="date-picker">
      <div class="row justify-center" style="max-width: 50vw">
        <q-date v-model="formattedDate" :mask="dateFormat" today-btn square flat/>
        <q-time v-model="formattedDate" :mask="dateFormat" format24h now-btn square flat>
          <div class="row items-center justify-end">
            <q-btn v-close-popup :label="$t('BTN_CLOSE')" color="primary" flat/>
          </div>
        </q-time>
      </div>
  </q-dialog>
</template>

<style scoped>
.q-dialog__inner--minimized > div {
  max-width: unset !important;
}
</style>
