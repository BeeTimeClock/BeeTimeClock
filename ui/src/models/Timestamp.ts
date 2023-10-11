export interface Timestamp {
  ID: number;
  UserID: number;
  ComingTimestamp: Date;
  GoingTimestamp: Date;
  IsHomeoffice: boolean;
  Corrections: TimestampCorrection[];
}

export interface TimestampGroup {
  Date: Date,
  IsHomeoffice: boolean,
  Timestamps: Timestamp[],
  WorkingHours: number,
  SubtractedHours: number,
}

export interface TimestampCreateRequest {
  ComingTimestamp: Date;
  GoingTimestamp: Date;
  IsHomeoffice: boolean;
	ChangeReason: string;
}

export interface TimestampCorrection {
  ID: number;
  ChangeReason: string;
  OldComingTimestamp: Date;
  OldGoingTimestamp: Date;
  CreatedAt: Date;
}

export interface TimestampCorrectionCreateRequest {
  NewComingTimestamp: Date;
  NewGoingTimestamp: Date;
  ChangeReason: string;
}
