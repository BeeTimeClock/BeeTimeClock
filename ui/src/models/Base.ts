import {AxiosError} from 'axios';

export interface BaseResponse<T> {
  Status: string;
  Timestamp: Date;
  Message: string|undefined;
  Error: string|undefined;
  Data: T;
}
export type ErrorResponse = AxiosError<BaseResponse<undefined>>

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

export interface UserApikey {
  ID: number;
  Description: string;
  Apikey: string
  ValidTill: Date;
}

export interface UserApikeyCreateRequest {
  Description: string;
  ValidTill: Date;
}
