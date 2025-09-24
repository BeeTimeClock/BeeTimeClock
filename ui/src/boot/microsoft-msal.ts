import { boot } from 'quasar/wrappers';
import type { IPublicClientApplication } from '@azure/msal-browser';
import {
  InteractionRequiredAuthError,
  PublicClientApplication,
} from '@azure/msal-browser';
import { useAuthStore } from 'stores/microsoft-auth';
import BeeTimeClock from 'src/service/BeeTimeClock';

declare module '@vue/runtime-core' {
  interface ComponentCustomProperties {
    $msalProvider: MsalProvider;
  }
}

class MsalProvider {
  _msalInstance?: IPublicClientApplication;

  set msalInstance(msalInstance: IPublicClientApplication) {
    this._msalInstance = msalInstance;
  }
  get msalInstance() {
    console.log('Msal Instance: ', this._msalInstance);
    return this._msalInstance!;
  }

  refresh() {
    const accessTokenRequest = {
      scopes: ['openid', 'profile', 'email', 'Calendars.ReadWrite'],
    };

    const msalInstance = this.msalInstance;

    msalInstance
      .acquireTokenSilent(accessTokenRequest)
      .then(function (accessTokenResponse) {
        useAuthStore().setAccessToken(accessTokenResponse.idToken);
      })
      .catch(function (error) {
        console.log(error);
        if (error instanceof InteractionRequiredAuthError) {
          void msalInstance.acquireTokenRedirect(accessTokenRequest);
        }
      });
  }
}

const msalProvider = new MsalProvider();

export default boot(async({ app }) => {
  const authStore = useAuthStore();
  const microsoftSettings = await BeeTimeClock.getMicrosoftAuthSettings();

  authStore.setMicrosoftClientId(microsoftSettings.data.Data.ClientID);
  authStore.setMicrosoftAuthority(microsoftSettings.data.Data.TenantID);

  console.log('msal settings: ', authStore.msalConfig)

  msalProvider.msalInstance = await PublicClientApplication.createPublicClientApplication(authStore.msalConfig)

  app.config.globalProperties.$msalProvider = msalProvider;
});

export { msalProvider };
