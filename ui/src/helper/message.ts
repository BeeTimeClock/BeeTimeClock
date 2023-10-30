import {Notify} from 'quasar';

export function showInfoMessage(message: string) {
  Notify.create({
    message: message,
    type: 'info',
    timeout: 5000,
    closeBtn: true,
  })
}

export function showErrorMessage(message: string | undefined) {
  Notify.create({
    message: message ?? 'no message',
    type: 'error',
    closeBtn: true,
  })
}
