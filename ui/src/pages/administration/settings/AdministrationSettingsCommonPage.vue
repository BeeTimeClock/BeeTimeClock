<script setup lang="ts">
import { computed, onMounted, ref } from 'vue';
import BeeTimeClock from 'src/service/BeeTimeClock';
import { showInfoMessage } from 'src/helper/message';

const logoFile = ref<File|Blob>();

function loadLogo() {
  BeeTimeClock.getLogo().then(result => {
    if (result.status === 200) {
      logoFile.value = result.data;
    }
  })
}


const imageUrl = computed(() => {
  if (logoFile.value) {
    return URL.createObjectURL(logoFile.value);
  }
  return '';
})


function save() {
  if (logoFile.value) {
    BeeTimeClock.administrationUploadLogoFile(logoFile.value).then((result) => {
      if (result.status === 204) {
        showInfoMessage('Saved');
      }
    });
  }
}

onMounted(() => {
  loadLogo()
})
</script>

<template>
  <q-page padding>
    <div class="row">
      <div class="col-3">
        <q-file v-model="logoFile" :label="$t('LABEL_LOGO')" />
      </div>
      <div class="col q-pa-md">
        <q-img :src="imageUrl" height="100px" :fit="'contain'" />
      </div>
    </div>
    <q-btn :label="$t('LABEL_SAVE')" class="full-width" color="positive" icon="save" @click="save" />
  </q-page>
</template>

<style scoped></style>
