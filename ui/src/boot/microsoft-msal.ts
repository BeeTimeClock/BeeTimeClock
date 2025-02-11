import {boot} from 'quasar/wrappers'
import {InteractionRequiredAuthError, IPublicClientApplication, PublicClientApplication} from '@azure/msal-browser';
import {useAuthStore} from 'stores/microsoft-auth';
import BeeTimeClock from 'src/service/BeeTimeClock';
import {getCurrentInstance} from 'vue';

declare module '@vue/runtime-core' {
  interface ComponentCustomProperties {
    $msalProvider: MsalProvider;
  }
}

class MsalProvider {
  msalInstance: IPublicClientApplication;

  constructor(msalInstance: IPublicClientApplication) {
    this.msalInstance = msalInstance;
  }

  refresh() {
    const accessTokenRequest = {
      scopes: [
        'openid',
        'profile',
        'email',
        'Calendars.ReadWrite',
      ],
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
          msalInstance.acquireTokenRedirect(accessTokenRequest);
        }
      });
  }
}

export default boot(async ({app}) => {
  const microsoftSettings = await BeeTimeClock.getMicrosoftAuthSettings();

  useAuthStore().setMicrosoftClientId(microsoftSettings.data.Data.ClientID);
  useAuthStore().setMicrosoftAuthority(microsoftSettings.data.Data.TenantID);

  app.config.globalProperties.$msalProvider = new MsalProvider(await PublicClientApplication.createPublicClientApplication(useAuthStore().getMsalConfig));
})

export function useMicrosoftAuth() {
  const app = getCurrentInstance();
  return app?.appContext.config.globalProperties.$msalProvider;
}
