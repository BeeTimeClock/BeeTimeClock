<script setup lang="ts">
import { useI18n } from 'vue-i18n';
import { copyToClipboard } from 'quasar';
import { showInfoMessage } from 'src/helper/message';

const { t } = useI18n();
const show = defineModel<boolean>('show', { default: false });
const value = defineModel<string>('value', { required: true });
const emits = defineEmits(['closed']);

function copy() {
  void copyToClipboard(value.value);
  showInfoMessage(t('MSG_COPIED_TO_CLIPBOARD'));
}

function close() {
  show.value = false;
  emits('closed')
}
</script>

<template>
  <q-dialog v-model="show">
    <q-card style="min-width: 400px">
      <q-card-section class="bg-primary text-white text-subtitle2">
        {{ t('LABEL_CLIPBOARD') }}
      </q-card-section>
      <q-card-section>
        <div class="q-mb-sm">{{ t('MSG_APIKEY_SHOW_WARNING') }}</div>
        <q-input v-model="value" readonly outlined>
          <template v-slot:after>
            <q-btn round dense flat icon="content_copy" @click="copy()" />
          </template>
        </q-input>
      </q-card-section>
      <q-card-actions align="right">
        <q-btn flat color="primary" :label="t('BTN_CLOSE')" @click="close" />
      </q-card-actions>
    </q-card>
  </q-dialog>
</template>

<style scoped></style>
