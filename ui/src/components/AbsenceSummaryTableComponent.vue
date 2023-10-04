<script setup lang="ts">
import {date} from 'quasar';
import {computed, PropType} from 'vue';
import {AbsenceSummaryItem} from 'src/models/Absence';
import {useI18n} from 'vue-i18n';
import {User} from "src/models/Authentication";

const {t} = useI18n();

const props = defineProps({
  modelValue: {
    type: Array as PropType<AbsenceSummaryItem[]>,
    required: true,
  },
  flat: {
    type: Boolean,
    default: false,
  },
  title: {
    type: String,
  }
})

const getTitle = computed(() => {
  if (props.title) {
    return props.title;
  }
  return t('LABEL_EMPLOYEE_ABSENCES')
})

const columns = [
  {
    name: 'absenceFrom',
    label: t('LABEL_FROM'),
    field: 'AbsenceFrom',
    format: (val: Date) => date.formatDate(val, 'DD. MMM. YYYY')
  },
  {
    name: 'absenceTill',
    label: t('LABEL_TILL'),
    field: 'AbsenceTill',
    format: (val: Date) => date.formatDate(val, 'DD. MMM. YYYY')
  },
  {
    name: 'user',
    label: t('LABEL_USER'),
    field: 'User',
    format: (val: User) => `${val.FirstName} ${val.LastName}`
  },
]
</script>

<template>
  <q-table :title="getTitle" :rows="modelValue" :columns="columns" hide-pagination :flat="flat"/>
</template>

<style scoped>

</style>
