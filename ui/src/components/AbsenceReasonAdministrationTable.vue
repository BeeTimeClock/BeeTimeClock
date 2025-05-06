<script setup lang="ts">
import { computed, onMounted, ref } from 'vue';
import { AbsenceReason, AbsenceReasonImpact } from 'src/models/Absence';
import { useI18n } from 'vue-i18n';
import BeeTimeClock from 'src/service/BeeTimeClock';
import { showInfoMessage } from 'src/helper/message';

const { t } = useI18n();
const absenceReasons = ref<AbsenceReason[]>([]);
const selectedAbsenceReason = ref<AbsenceReason>();
const showDialog = ref(false);

const isNewReason = computed(() => {
  if (!selectedAbsenceReason.value) return true;
  return !selectedAbsenceReason.value.ID;
});

function loadAbsenceReasons() {
  BeeTimeClock.absenceReasons().then((result) => {
    if (result.status === 200) {
      absenceReasons.value = result.data.Data.map((s) =>
        AbsenceReason.fromApi(s)
      );
    }
  });
}

function saveAbsenceReason() {
  if (!selectedAbsenceReason.value) return;

  if (isNewReason.value) {
    BeeTimeClock.administrationCreateAbsenceReason(
      selectedAbsenceReason.value
    ).then((result) => {
      if (result.status === 201) {
        showInfoMessage(t('MSG_CREATE_SUCCESS', { item: t('LABEL_REASON') }));
        showDialog.value = false;
      }
    });
  } else {
    BeeTimeClock.administrationUpdateAbsenceReason(
      selectedAbsenceReason.value.ID,
      selectedAbsenceReason.value
    ).then((result) => {
      if (result.status === 200) {
        showInfoMessage(t('MSG_UPDATE_SUCCESS'));
        showDialog.value = false;
      }
    });
  }
}

function createAbsenceReason() {
  selectedAbsenceReason.value = new AbsenceReason();
  showDialog.value = true;
}

function editAbsenceReason(absenceReason: AbsenceReason) {
  selectedAbsenceReason.value = absenceReason;
  showDialog.value = true;
}

onMounted(() => {
  loadAbsenceReasons();
});
</script>

<template>
  <q-list separator>
    <q-item v-for="absenceReason in absenceReasons" :key="absenceReason.ID">
      <q-item-section>
        <q-item-label>{{ absenceReason.Description }}</q-item-label>
        <q-item-label caption v-if="absenceReason.Impact">{{
          absenceReason.Impact
        }}</q-item-label>
        <q-item-label
          caption
          v-if="absenceReason.Impact == AbsenceReasonImpact.Hours"
          >{{ absenceReason.ImpactHours }}{{$t('LABEL_HOUR', absenceReason.ImpactHours)}}</q-item-label
        >
      </q-item-section>
      <q-item-section side>
        <q-btn
          icon="edit"
          color="primary"
          @click="editAbsenceReason(absenceReason)"
        />
      </q-item-section>
    </q-item>
    <q-item>
      <q-item-section>
        <q-btn
          class="full-width"
          icon="add"
          :label="$t('BTN_ADD')"
          color="positive"
          @click="createAbsenceReason"
        />
      </q-item-section>
    </q-item>
  </q-list>
  <q-dialog v-model="showDialog" v-if="selectedAbsenceReason">
    <q-card>
      <q-card-section class="bg-primary text-h6 text-white">
        <div v-if="isNewReason">
          {{ $t('TITLE_CREATE', { item: $t('LABEL_REASON') }) }}
        </div>
        <div v-else>
          {{ $t('TITLE_UPDATE', { item: $t('LABEL_REASON') }) }}
        </div>
      </q-card-section>
      <q-form @submit="saveAbsenceReason">
        <q-card-section>
          <q-input
            v-model="selectedAbsenceReason.Description"
            :label="$t('LABEL_DESCRIPTION')"
          />
        </q-card-section>
        <q-card-section>
          <q-card-actions>
            <q-btn color="negative" :label="$t('BTN_CANCEL')" v-close-popup />
            <q-btn
              color="positive"
              :label="isNewReason ? $t('BTN_CREATE') : $t('BTN_SAVE')"
              type="submit"
            />
          </q-card-actions>
        </q-card-section>
      </q-form>
    </q-card>
  </q-dialog>
</template>

<style scoped></style>
