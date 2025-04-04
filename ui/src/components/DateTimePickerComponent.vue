<script setup lang="ts">
import { computed, ref, watch } from 'vue';
import { date } from 'quasar';

const props = defineProps({
  modelValue: {
    type: Date,
    required: true
  },
  label: {
    type: String
  }
});
const emit = defineEmits(['update:modelValue']);

const dateFormat = 'DD.MM.YYYY HH:mm';
const value = computed({
  get(): Date {
    return props.modelValue;
  },
  set(value: Date) {
    emit('update:modelValue', value);
  }
});

const formattedDate = ref(date.formatDate(value.value, dateFormat));

watch(formattedDate, (newDateValue) => {
  value.value = date.extractDate(newDateValue, dateFormat);
})
</script>

<template>
  <q-input filled v-model="formattedDate" :label="label" mask="">
    <template v-slot:prepend>
      <q-icon name="event" class="cursor-pointer">
        <q-popup-proxy cover transition-show="scale" transition-hide="scale">
          <q-date v-model="formattedDate" :mask="dateFormat" today-btn>
            <div class="row items-center justify-end">
              <q-btn v-close-popup :label="$t('BTN_CLOSE')" color="primary" flat />
            </div>
          </q-date>
        </q-popup-proxy>
      </q-icon>
    </template>

    <template v-slot:append>
      <q-icon name="access_time" class="cursor-pointer">
        <q-popup-proxy cover transition-show="scale" transition-hide="scale">
          <q-time v-model="formattedDate" :mask="dateFormat" format24h now-btn>
            <div class="row items-center justify-end">
              <q-btn v-close-popup :label="$t('BTN_CLOSE')" color="primary" flat />
            </div>
          </q-time>
        </q-popup-proxy>
      </q-icon>
    </template>
  </q-input>
</template>

<style scoped>

</style>
