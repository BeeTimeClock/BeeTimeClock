import {defineStore} from 'pinia';

const SETTINGS_STORE_KEY = 'settings';

interface SettingsState {
  language: string;
}

export const useSettingsStore = defineStore('settings', {
  state: () => {
    return {
       language: '',
    } as SettingsState;
  },
  getters: {
    language() {
      const data = localStorage.getItem(SETTINGS_STORE_KEY)

      if (data == null) {
        return null;
      }
      const settings = JSON.parse(data) as SettingsState;
      return settings.language;
    },
  },

  actions: {
    setLanguage(lang: string) {
      this.$state.language = lang;
      this.save();
    },
    save() {
      localStorage.setItem(SETTINGS_STORE_KEY, JSON.stringify(this.$state));
    },
    load() {

    }
  },
});
