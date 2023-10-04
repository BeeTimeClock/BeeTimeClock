export interface BaseResponse<T> {
  Status: string;
  Timestamp: Date;
  Message: string|undefined;
  Error: string|undefined;
  Data: T;
}

export interface BackendStatus {
  Commit: string;
}

export interface MicrosoftAuthSettings {
  ClientID: string;
  TenantID: string;
}

export interface AuthProviders {
  Local: boolean;
  Microsoft: boolean;
}
