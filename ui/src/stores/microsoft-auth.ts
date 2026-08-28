import { defineStore } from 'pinia';
import { type Configuration } from '@azure/msal-browser';
import BeeTimeClock from 'src/service/BeeTimeClock';
import { type User } from 'src/models/Authentication';

export const ACCESS_TOKEN_STORE_KEY = 'accessToken';
export const AUTH_PROVIDER_STORE_KEY = 'authProvider';
const SESSION_STORE_KEY = 'session';

export const useAuthStore = defineStore('auth', {
  state: () => {
    return {
      accessToken: localStorage.getItem(ACCESS_TOKEN_STORE_KEY) ?? '',
      authProvider: localStorage.getItem(AUTH_PROVIDER_STORE_KEY) ?? '',
      msalConfig: {
        auth: {
          clientId: '',
          authority: '',
          redirectUri: window.location.origin,
        },
        cache: {
          cacheLocation: 'localStorage',
        },
        system: {
          allowNativeBroker: false,
        },
      } as Configuration,
    };
  },
  getters: {
    getMsalConfig(state) {
      return state.msalConfig;
    },
    getAccessToken(state): undefined | string {
      return state.accessToken !== '' ? state.accessToken : undefined;
    },
    getAuthProvider(state): string {
      return state.authProvider;
    },
    loggedIn(state): boolean {
      return state.accessToken != '';
    },
  },

  actions: {
    getSession(): User | null {
      const session = localStorage.getItem(SESSION_STORE_KEY);

      if (!session) return null;
      return JSON.parse(session) as User;
    },
    isAdministrator(): boolean {
      const session = this.getSession();
      if (!session) return false;
      return session.AccessLevel == 'admin';
    },
    logout() {
      localStorage.clear();
      sessionStorage.clear();
      this.accessToken = '';
      this.authProvider = '';
    },
    setAccessToken(token: string) {
      this.accessToken = token;
      localStorage.setItem(ACCESS_TOKEN_STORE_KEY, token);
    },
    setAuthProvider(provider: string) {
      this.authProvider = provider;
      localStorage.setItem(AUTH_PROVIDER_STORE_KEY, provider);
    },
    setMicrosoftAuthority(tenantId: string) {
      this.$state.msalConfig.auth.authority = `https://login.microsoftonline.com/${tenantId}`;
    },
    setMicrosoftClientId(clientId: string) {
      this.$state.msalConfig.auth.clientId = clientId;
    },
    async loadSession(): Promise<boolean> {
      try {
        const result = await BeeTimeClock.getMeUser();

        localStorage.setItem(
          SESSION_STORE_KEY,
          JSON.stringify(result.data.Data),
        );
        return true;
      } catch (err) {
        console.log(err);
        return false;
      }
    },
  },
});
