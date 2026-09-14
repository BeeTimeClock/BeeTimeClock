<script setup lang="ts">
import type { User } from 'src/models/Authentication';
import { computed, ref, watch } from 'vue';
import { useI18n } from 'vue-i18n';

const { t } = useI18n();
const userId = defineModel<number | undefined>();
const users = defineModel<User[]>('users', { default: [], required: true });

const sortedUsers = computed(() => {
  return users.value.toSorted((a, b) =>
    a.displayName.toLowerCase() < b.displayName.toLowerCase() ? -1 : 1,
  );
});

const options = ref<User[]>([]);

watch(
  sortedUsers,
  (value) => {
    options.value = value;
  },
  { immediate: true },
);

function filterUsers(val: string, update: (arg0: () => void) => void) {
  update(() => {
    const needle = val.toLowerCase();
    options.value = sortedUsers.value.filter((v) =>
      v.displayName.toLowerCase().includes(needle),
    );
  });
}
</script>

<template>
  <q-select
    v-if="users"
    v-model="userId"
    :options="options"
    emit-value
    option-label="displayName"
    option-value="ID"
    map-options
    @filter="filterUsers"
    use-input
    :label="t('LABEL_USER')"
    input-debounce="0"
  />
</template>

<style scoped></style>
