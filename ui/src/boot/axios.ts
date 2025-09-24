import { boot } from 'quasar/wrappers';
import type { AxiosInstance } from 'axios';
import axios from 'axios';
import {ACCESS_TOKEN_STORE_KEY, AUTH_PROVIDER_STORE_KEY} from 'stores/microsoft-auth';

declare module '@vue/runtime-core' {
  interface ComponentCustomProperties {
    $axios: AxiosInstance;
  }
}

// Be careful when using SSR for cross-request state pollution
// due to creating a Singleton instance here;
// If any client changes this (global) instance, it might be a
// good idea to move this instance creation inside of the
// "export default () => {}" function below (which runs individually
// for each client)
const options = {
  baseURL: process.env.VUE_APP_BACKEND_ADDRESS || '',
  headers: {
    'Content-Type': 'application/json',
  }
};
const api = axios.create(options);

api.interceptors.request.use(request => {
  const token = localStorage.getItem(ACCESS_TOKEN_STORE_KEY);

  if (token != null && token !== '') {
    request.headers.setAuthorization(`Bearer ${token}`);
    request.headers.set('X-Auth-Provider', localStorage.getItem(AUTH_PROVIDER_STORE_KEY));
  } else {
    request.headers.setAuthorization(null);
  }

  return request;
});

export default boot(({ app }) => {
  // for use inside Vue files (Options API) through this.$axios and this.$api
  app.config.globalProperties.$axios = axios;
  // ^ ^ ^ this will allow you to use this.$axios (for Vue Options API form)
  //       so you won't necessarily have to import axios in each vue file

  app.config.globalProperties.$api = api;
  // ^ ^ ^ this will allow you to use this.$api (for Vue Options API form)
  //       so you can easily perform requests against your app's API
});

export { api };
