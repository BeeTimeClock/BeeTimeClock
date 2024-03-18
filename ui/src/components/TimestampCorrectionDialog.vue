<script setup lang="ts">

import DateTimePickerComponent from 'components/DateTimePickerComponent.vue';
import { Timestamp, TimestampCorrectionCreateRequest, TimestampCreateRequest } from 'src/models/Timestamp';
import { computed, PropType, ref } from 'vue';
import BeeTimeClock from 'src/service/BeeTimeClock';
import { showInfoMessage } from 'src/helper/message';
import { date } from 'quasar';
import formatDate = date.formatDate;
import { useI18n } from 'vue-i18n';

const {t} = useI18n();
const timestampCorrectionCreateRequest = ref<TimestampCorrectionCreateRequest|null>(null);

const props = defineProps({
  modelValue: {
    type: Object as PropType<Timestamp>,
  },
  modelShow: {
    type: Boolean,
    default: false,
  }
})

const emits = defineEmits(['update:modelShow', 'refresh'])

const show = computed({
  get() {
    return props.modelShow;
  },
  set(val: boolean) {
    emits('update:modelShow', val);
  }
})


function onShow() {
  timestampCorrectionCreateRequest.value = {} as TimestampCorrectionCreateRequest;

  if (props.modelValue != undefined) {
    timestampCorrectionCreateRequest.value.NewComingTimestamp = props.modelValue.ComingTimestamp;
    if (props.modelValue.GoingTimestamp) {
      timestampCorrectionCreateRequest.value.NewGoingTimestamp = props.modelValue.GoingTimestamp;
    }
    timestampCorrectionCreateRequest.value.IsHomeoffice = props.modelValue.IsHomeoffice
  } else {
    timestampCorrectionCreateRequest.value.NewComingTimestamp = new Date();
  }
}

function createTimestampCorrection() {
  if (!timestampCorrectionCreateRequest.value) return;

  if (props.modelValue != undefined) {
    BeeTimeClock.timestampCorrectionCreate(props.modelValue, timestampCorrectionCreateRequest.value).then((result) => {
      if (result.status === 200) {
        showInfoMessage(t('MSG_CREATE_SUCCESS'));
        emits('refresh');
        show.value = false
      }
    });
  } else {
    const timestampCreateRequest = {
      ComingTimestamp: new Date(timestampCorrectionCreateRequest.value.NewComingTimestamp),
      ChangeReason: timestampCorrectionCreateRequest.value.ChangeReason,
      IsHomeoffice: timestampCorrectionCreateRequest.value.IsHomeoffice,
    } as TimestampCreateRequest;

    if (timestampCorrectionCreateRequest.value.NewGoingTimestamp) {
      timestampCreateRequest.GoingTimestamp = new Date(timestampCorrectionCreateRequest.value.NewGoingTimestamp)
    }

    BeeTimeClock.timestampCreate(timestampCreateRequest).then((result) => {
      if (result.status === 201) {
        showInfoMessage(t('MSG_CREATE_SUCCESS'));
        emits('refresh');
        show.value = false
      }
    });
  }
}
</script>

<template>
  <q-dialog v-model="show" persistent @beforeShow="onShow()">
    <q-card class="q-dialog-plugin full-width" v-if="timestampCorrectionCreateRequest">
      <q-card-section>
        <div class="text-h6">{{ $t('LABEL_TIMESTAMP_CORRECTION_CREATE') }}</div>
      </q-card-section>
      <q-form @submit.prevent.stop="createTimestampCorrection">
        <q-card-section>
          <DateTimePickerComponent v-model="timestampCorrectionCreateRequest.NewComingTimestamp"
                                   :label="$t('LABEL_COMING_TIMESTAMP')" />
          <DateTimePickerComponent class="q-mt-md" v-model="timestampCorrectionCreateRequest.NewGoingTimestamp"
                                   :label="$t('LABEL_GOING_TIMESTAMP')" />
          <q-checkbox v-model="timestampCorrectionCreateRequest.IsHomeoffice" :label="$t('LABEL_HOMEOFFICE')" />
          <q-input :rules="[val => val.trim().length >= 20 || t('LABEL_FIELD_REQUIRED') + ' ' + t('LABEL_MIN_CHARS', {count: 20})]" v-model="timestampCorrectionCreateRequest.ChangeReason"
                   type="textarea" :label="`${t('LABEL_REASON')} (${t('LABEL_MIN_CHARS', {count: 20})})`" />
        </q-card-section>
        <q-card-actions>
          <q-btn v-close-popup :label="$t('BTN_CANCEL')" color="negative" type="reset" />
          <q-btn :label="$t('BTN_CREATE')" color="positive" type="submit" />
        </q-card-actions>
      </q-form>
    </q-card>
  </q-dialog>
</template>

<style scoped>

</style>
