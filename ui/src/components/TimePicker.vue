<script setup lang="ts">
import { computed } from 'vue';

const value = defineModel('value', {type: Date, required: true})
const label = defineModel('label', {type: String})
const readonly = defineModel('readonly', {type: Boolean, default: false})

const time = computed({
  get() {
     const current = new Date(value.value);
     return `${current.getHours()}:${current.getMinutes()}`
  },
  set(val) {
    const parts = val.split(':')
    const newDate = new Date(value.value)
    newDate.setHours(parseInt(parts[0]!))
    newDate.setMinutes(parseInt(parts[1]!))

    value.value = newDate.toISOString();
  }
})
</script>

<template>
  <q-input filled v-model="time" mask="time" :rules="['time']" :label="label" :readonly="readonly" >
    <template v-slot:append>
      <q-icon name="access_time" class="cursor-pointer">
        <q-popup-proxy cover transition-show="scale" transition-hide="scale">
          <q-time v-model="time">
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
