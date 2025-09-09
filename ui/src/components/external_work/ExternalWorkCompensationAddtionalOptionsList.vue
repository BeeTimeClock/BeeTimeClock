<script setup lang="ts">
import { ref } from 'vue';
import { ExternalWorkCompensation } from 'src/models/ExternalWork';
import { useQuasar } from 'quasar';
import { useI18n } from 'vue-i18n';

const q = useQuasar();
const {t} = useI18n();
const workCompensation = defineModel({
  type: ExternalWorkCompensation,
  required: true,
});
const showDialog = ref(false);
const newAdditionalOptionKey = ref<string>();
const newAdditionalOptionValue = ref<number>();

function submit() {
  if (!newAdditionalOptionValue.value) return;
  if (newAdditionalOptionKey.value == undefined) return;

  workCompensation.value.AdditionalOptions[newAdditionalOptionKey.value] =
    newAdditionalOptionValue.value;
  showDialog.value = false;
}

function deleteEntry(key: string|number) {
  q.dialog({
    title: t('LABEL_DELETE'),
    message: t('MSG_DELETE', {
      item: t('LABEL_OPTION'),
      identifier: key,
    }),
    cancel: true,
    persistent: true,
  }).onOk(() => {
    delete workCompensation.value.AdditionalOptions[key];
  })
}
</script>

<template>
  <q-list>
    <q-item
      v-for="(value, key) in workCompensation.AdditionalOptions"
      :key="key"
    >
      <q-item-section>
        <q-input readonly :model-value="key" :label="t('LABEL_NAME')" />
      </q-item-section>
      <q-item-section>
        <q-input
          v-model.number="workCompensation.AdditionalOptions[key]"
          :label="t('LABEL_COMPENSATION_IN_EURO')"
        />
      </q-item-section>
      <q-item-section side>
        <q-btn icon="delete" color="negative" @click="deleteEntry(key)"/>
      </q-item-section>
    </q-item>
    <q-item>
      <q-item-section>
        <q-btn
          class="full-width"
          color="positive"
          icon="add"
          :label="t('LABEL_ADD')"
          @click="showDialog = true"
        />
      </q-item-section>
    </q-item>
  </q-list>
  <q-dialog v-model="showDialog" persistent>
    <q-card>
      <q-card-section class="bg-primary text-white text-h6">
        {{ t('LABEL_ADD') }}
      </q-card-section>
      <q-form @submit="submit">
        <q-card-section>
          <q-input v-model="newAdditionalOptionKey" :label="t('LABEL_NAME')" />
          <q-input
            v-model.number="newAdditionalOptionValue"
            :label="t('LABEL_COMPENSATION_IN_EURO')"
          />
        </q-card-section>
        <q-card-section>
          <q-card-actions>
            <q-btn
              :label="t('LABEL_CANCEL')"
              color="negative"
              v-close-popup
              type="reset"
            />
            <q-btn
              :label="t('LABEL_CREATE')"
              color="positive"
              type="submit"
            />
          </q-card-actions>
        </q-card-section>
      </q-form>
    </q-card>
  </q-dialog>
</template>

<style scoped></style>
