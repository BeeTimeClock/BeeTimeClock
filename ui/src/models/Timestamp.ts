export interface Timestamp {
  ID: number;
  UserID: number;
  ComingTimestamp: Date;
  GoingTimestamp: Date|null;
  IsHomeoffice: boolean;
  IsHomeofficeGoing: boolean;x
  Corrections: TimestampCorrection[];
}

export interface TimestampGroup {
  Date: Date,
  IsHomeoffice: boolean,
  Timestamps: Timestamp[],
  WorkingHours: number,
  SubtractedHours: number,
  OvertimeHours: number,
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
  NewGoingTimestamp: Date|null;
  ChangeReason: string;
  IsHomeoffice: boolean;
}

export interface SumResponse {
  Total: number;
}

export interface OvertimeResponse {
  Total: number;
  Subtracted: number;
  Needed: number;
}

export interface TimestampYearMonthGrouped {
  [details: number]: number[];
}
