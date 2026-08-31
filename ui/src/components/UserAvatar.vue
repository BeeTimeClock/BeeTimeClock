<script setup lang="ts">
import { onMounted, ref, watch } from 'vue';
import type { ApiUser } from 'src/models/Authentication';
import { gravatarUrl } from 'src/helper/gravatar';

const props = defineProps<{ user: ApiUser; size?: number }>();
const url = ref<string | null>(null);

async function load() {
  if (props.user.AllowGravatar) {
    url.value = await gravatarUrl(props.user.Username, props.size ?? 40);
  } else {
    url.value = null;
  }
}

watch(() => props.user, load);
onMounted(load);
</script>

<template>
  <q-avatar :size="`${size ?? 40}px`">
    <img v-if="url" :src="url" referrerpolicy="no-referrer" />
    <q-icon v-else name="person" />
  </q-avatar>
</template>
