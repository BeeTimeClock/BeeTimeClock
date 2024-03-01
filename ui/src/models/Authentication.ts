export interface AuthRequest {
  Username: string;
  Password: string;
}

export interface AuthResponse {
  Token: string;
}

export interface User {
  ID: number;
  Username: string;
  FirstName: string;
  LastName: string;
  AccessLevel: string;
  OvertimeSubtractionModel: string;
  OvertimeSubtractionAmount: number;
}

