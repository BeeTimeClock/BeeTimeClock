<script setup lang="ts">
import {date} from 'quasar';
import {AbsenceSummaryItem} from 'src/models/Absence';
import {User} from 'src/models/Authentication';
import {computed, PropType} from 'vue';
import {useI18n} from 'vue-i18n';
import {useAuthStore} from 'stores/microsoft-auth';

const {t} = useI18n();
const auth = useAuthStore();

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

const rows = computed(() => {
  if (!props.modelValue) return [];

  const data = props.modelValue;
  return data.sort((a, b) => new Date(a.AbsenceFrom).getTime() - new Date(b.AbsenceFrom).getTime());
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
    name: 'absenceNettoDays',
    label: t('LABEL_NETTO_DAYS'),
    field: 'NettoDays',
  },
  {
    name: 'user',
    label: t('LABEL_USER'),
    field: 'User',
    format: (val: User) => `${val.FirstName} ${val.LastName}`
  },
  ]

  const pagination = {
    rowsPerPage: 10,
  }

if (auth.isAdministrator()) {
  columns.push({
    name: 'absenceReason',
    label: t('LABEL_REASON'),
    field: 'Reason',
    format: (val: string) => val,
  })
}
</script>

<template>
  <q-table :title="getTitle" :rows="rows" :columns="columns" :flat="flat" :pagination="pagination"/>
</template>

<style scoped>

</style>
