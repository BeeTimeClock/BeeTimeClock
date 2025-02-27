<script setup lang="ts">
const value = defineModel('label',{type: String})
const caption = defineModel('caption', {type: String})
</script>

<template>
  <q-item>
    <q-item-section>
      <q-item-label v-if="caption" caption>{{caption}}</q-item-label>
      <q-item-label>{{value}}</q-item-label>
    </q-item-section>
  </q-item>
</template>

<style scoped>

</style>
