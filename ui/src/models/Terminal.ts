import { autoImplement } from 'src/helper/functions';

export interface ApiTerminal {
  ID: number;
  TerminalName: string;
  ClientId: string;
  Apikey: string;
}

export interface ApiTerminalCreateRequest {
  TerminalName: string;
}

export interface ApiUserToken {
  ID: number;
  UserID: number;
  TokenType: string;
  TokenIdentifier: string;
}

export interface ApiUserTokenCreateRequest {
  TokenType?: string;
  TokenIdentifier: string;
}

export class Terminal extends autoImplement<ApiTerminal>() {
  static fromApi(apiItem: ApiTerminal) : Terminal {
    return new Terminal(apiItem);
  }
}

export class UserToken extends autoImplement<ApiUserToken>() {
  static fromApi(apiItem: ApiUserToken) : UserToken {
    return new UserToken(apiItem);
  }
}
