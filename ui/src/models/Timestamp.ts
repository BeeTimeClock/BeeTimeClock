export interface Timestamp {
  ID: number;
  UserID: number;
  ComingTimestamp: Date;
  GoingTimestamp: Date;
  IsHomeoffice: boolean;
}

export interface TimestampGroup {
  Date: Date,
  IsHomeoffice: boolean,
  Timestamps: Timestamp[],
  WorkingHours: number,
  SubtractedHours: number,
}

export interface TimestampCreateRequest {
  IsHomeoffice: boolean;
}
