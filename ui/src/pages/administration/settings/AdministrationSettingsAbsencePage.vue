<script setup lang="ts">
import AbsenceReasonAdministrationTable from 'components/AbsenceReasonAdministrationTable.vue';
</script>

<template>
  <q-page padding>
    <q-card>
      <q-card-section class="bg-primary text-white text-subtitle2">
        {{ $t('LABEL_ABSENCE_REASON', 2) }}
      </q-card-section>
      <q-card-section>
        <AbsenceReasonAdministrationTable />
      </q-card-section>
    </q-card>
  </q-page>
</template>

<style scoped></style>
