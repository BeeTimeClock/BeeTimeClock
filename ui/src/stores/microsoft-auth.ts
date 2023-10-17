import {defineStore} from 'pinia';
import {Configuration} from "@azure/msal-browser";
import BeeTimeClock from "src/service/BeeTimeClock";
import {User} from "src/models/Authentication";
import { Cookies } from 'quasar';


export const ACCESS_TOKEN_STORE_KEY = 'accessToken';
export const AUTH_PROVIDER_STORE_KEY = 'authProvider';
const SESSION_STORE_KEY = 'session';

export const useAuthStore = defineStore('auth', {
  state: () => {
    return {
      msalConfig: {
        auth: {
          clientId: '',
          authority: '',
          redirectUri: process.env.APP_URL,
        },
        cache: {
          cacheLocation: 'localStorage',
        },
        system: {
          allowNativeBroker: false,
        }
      } as Configuration,
    }
  },
  getters: {
    getMsalConfig(state) {
      return state.msalConfig;
    },
    getAccessToken(): undefined | string {
      return localStorage.getItem(ACCESS_TOKEN_STORE_KEY) ?? undefined
    },
    getAuthProvider(): string {
      return localStorage.getItem(AUTH_PROVIDER_STORE_KEY) ?? ''
    },
    loggedIn(): boolean {
      const accessToken = this.getAccessToken;
      return accessToken != undefined && accessToken != ''
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
      console.log(session);
      if (!session) return false;
      return session.AccessLevel == 'admin';
    },
    logout() {
      this.setAccessToken('');
      this.setAuthProvider('');
      localStorage.clear();
      sessionStorage.clear();
    },
    setAccessToken(token: string) {
      localStorage.setItem(ACCESS_TOKEN_STORE_KEY, token);
    },
    setAuthProvider(provider: string) {
      localStorage.setItem(AUTH_PROVIDER_STORE_KEY, provider);
    },
    setMicrosoftAuthority(tenantId: string) {
      this.$state.msalConfig.auth.authority = `https://login.microsoftonline.com/${tenantId}`;
    },
    setMicrosoftClientId(clientId: string) {
      this.$state.msalConfig.auth.clientId = clientId;
    },
    async loadSession(): boolean {
      try {
        const result = await BeeTimeClock.getMeUser();

        localStorage.setItem(SESSION_STORE_KEY, JSON.stringify(result.data.Data));
        return true;

      } catch (err) {
        console.log(err);
        return false;
      }
    },
  }
});
