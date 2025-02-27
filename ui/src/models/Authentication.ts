import { autoImplement } from 'src/helper/functions';

export interface AuthRequest {
  Username: string;
  Password: string;
}

export interface AuthResponse {
  Token: string;
}

export interface ApiUser {
  ID: number;
  Username: string;
  FirstName: string;
  LastName: string;
  AccessLevel: string;
  OvertimeSubtractionModel: string;
  OvertimeSubtractionAmount: number;
  StaffNumber: number;
}

export class User extends autoImplement<ApiUser>() {
  static fromApi(apiItem: ApiUser) : User {
    return new User(apiItem);
  }

  get displayName() {
    return `${this.FirstName} ${this.LastName} (${this.Username})`
  }
}

