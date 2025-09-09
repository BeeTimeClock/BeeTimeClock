import { autoImplement } from 'src/helper/functions';

export interface Timestamp {
  ID: number;
  UserID: number;
  ComingTimestamp: Date;
  GoingTimestamp: Date|null;
  IsHomeoffice: boolean;
  IsHomeofficeGoing: boolean;
  Corrections: TimestampCorrection[];
}

export interface ApiTimestampGroup {
  Date: Date,
  IsHomeoffice: boolean,
  Timestamps: Timestamp[],
  WorkingHours: number,
  SubtractedHours: number,
  OvertimeHours: number,
}

export class TimestampGroup extends autoImplement<ApiTimestampGroup>() {
  public expanded: boolean = false;

  static fromApi(apiItem: ApiTimestampGroup) : TimestampGroup {
    return new TimestampGroup(apiItem);
  }
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
  NewGoingTimestamp?: Date;
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
