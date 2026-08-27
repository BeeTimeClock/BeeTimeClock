<script setup lang="ts">
import { useI18n } from 'vue-i18n';
import BeeTimeClock from 'src/service/BeeTimeClock';
import { showErrorMessage, showInfoMessage } from 'src/helper/message';
import type { ApiTerminalCreateRequest} from 'src/models/Terminal';
import { Terminal } from 'src/models/Terminal';
import { onMounted, ref } from 'vue';
import { copyToClipboard, useQuasar, type QTableColumn } from 'quasar';
import ClipboardDialog from '../../../components/ClipboardDialog.vue';
import { emptyPagination } from 'src/helper/objects';

const { t } = useI18n();
const q = useQuasar();
const terminals = ref<Terminal[]>([]);

const terminalCreateRequest = ref<ApiTerminalCreateRequest>({} as ApiTerminalCreateRequest);
const apikey = ref<string|null>();
const showDialog = ref<boolean>(false);

const columns = [
  {
    name: 'id',
    label: t('LABEL_ID'),
    field: 'ID',
  },
  {
    name: 'terminalName',
    label: t('LABEL_NAME'),
    field: 'TerminalName',
  },
  {
    name: 'clientId',
    label: t('LABEL_CLIENT_ID'),
    field: 'ClientId',
  },
  {
    name: 'actions',
    label: t('LABEL_ACTION', 2),
    field: 'ID',
    align: 'right',
  },
] as QTableColumn[];

function loadTerminals() {
  BeeTimeClock.administrationTerminalList()
    .then((result) => {
      if (result.status === 200) {
        terminals.value = result.data.Data.map((s) => Terminal.fromApi(s));
      }
    })
    .catch((exception) => {
      showErrorMessage(exception);
    });
}

function createTerminal() {
  BeeTimeClock.administrationTerminalCreate(terminalCreateRequest.value)
    .then((result) => {
      if (result.status === 201) {
        showInfoMessage(t('MSG_CREATE_SUCCESS'))
        apikey.value = result.data.Data.Apikey;
        showDialog.value = true;
        loadTerminals();
      }
    })
    .catch((exception) => {
      showErrorMessage(exception);
    });
}

function regenerateTerminal(terminal: Terminal) {
  q.dialog({
    message: t('MSG_TERMINAL_REGENERATE_CONFIRM', {
      identifier: terminal.TerminalName,
    }),
    cancel: true,
    persistent: true,
  }).onOk(() => {
    BeeTimeClock.administrationTerminalRegenerate(terminal.ID)
      .then((result) => {
        if (result.status === 200) {
          showInfoMessage(t('MSG_UPDATE_SUCCESS'));
          apikey.value = result.data.Data.Apikey;
          showDialog.value = true;
          loadTerminals();
        }
      })
      .catch((exception) => {
        showErrorMessage(exception);
      });
  });
}

function deleteTerminal(terminal: Terminal) {
  q.dialog({
    message: t('MSG_DELETE', {
      item: t('LABEL_TERMINAL'),
      identifier: terminal.TerminalName,
    }),
    cancel: true,
    persistent: true,
  }).onOk(() => {
    BeeTimeClock.administrationTerminalDelete(terminal.ID)
      .then((result) => {
        if (result.status === 204) {
          showInfoMessage(t('MSG_DELETE_SUCCESS'));
          loadTerminals();
        }
      })
      .catch((exception) => {
        showErrorMessage(exception);
      });
  });
}

function copyClientId(clientId: string) {
  void copyToClipboard(clientId);
  showInfoMessage(t('MSG_COPIED_TO_CLIPBOARD'));
}

function clear() {
  apikey.value = null;
}

onMounted(() => {
  loadTerminals();
});
</script>

<template>
  <q-page padding>
    <q-table
      :rows="terminals"
      :columns="columns"
      hide-pagination
      :pagination="emptyPagination"
    >
      <template v-slot:body-cell-clientId="props">
        <q-td :props="props">
          <code class="q-mr-md">{{ props.value }}</code>
          <q-btn
            round
            dense
            outline
            size="xs"
            color="grey-7"
            icon="content_copy"
            @click="copyClientId(props.value)"
          />
        </q-td>
      </template>
      <template v-slot:body-cell-actions="props">
        <q-td :props="props" class="text-right">
          <q-btn
            icon="autorenew"
            color="primary"
            flat
            dense
            @click="regenerateTerminal(props.row)"
          >
            <q-tooltip>{{ t('BTN_REGENERATE_APIKEY') }}</q-tooltip>
          </q-btn>
          <q-btn
            icon="delete"
            color="negative"
            flat
            dense
            class="q-ml-sm"
            @click="deleteTerminal(props.row)"
          />
        </q-td>
      </template>
    </q-table>
    <q-card class="q-mt-md">
      <q-card-section class="bg-primary text-white text-subtitle2">
        {{ t('LABEL_CREATE', { item: t('LABEL_TERMINAL') }) }}
      </q-card-section>
      <q-form @submit="createTerminal">
        <q-card-section>
          <q-input
            v-model="terminalCreateRequest.TerminalName"
            :label="t('LABEL_NAME')"
            :rules="[(val) => !!val || t('LABEL_NAME')]"
          />
        </q-card-section>
        <q-card-section>
          <q-card-actions>
            <q-btn :label="t('BTN_CREATE')" type="submit" color="primary" />
          </q-card-actions>
        </q-card-section>
      </q-form>
    </q-card>
    <ClipboardDialog v-if="apikey" v-model:value="apikey" v-model:show="showDialog" @closed="clear()"/>
  </q-page>
</template>

<style scoped></style>
