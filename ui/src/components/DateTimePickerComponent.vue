<script setup lang="ts">
import { computed } from 'vue';

  const props = defineProps({
    modelValue: {
      type: String,
      required: true,
    },
    label: {
      type: String,
    }
  });
  const emit = defineEmits(['update:modelValue']);

  const value = computed({
    get() : string {
      return props.modelValue;
    },
    set(value: string) {
      emit('update:modelValue', value);
    }
  })
</script>

<template>
  <q-input filled v-model="value" :label="label">
    <template v-slot:prepend>
      <q-icon name="event" class="cursor-pointer">
        <q-popup-proxy cover transition-show="scale" transition-hide="scale">
          <q-date v-model="value" mask="DD.MM.YYYY HH:mm">
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
          <q-time v-model="value" mask="DD.MM.YYYY HH:mm" format24h>
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
