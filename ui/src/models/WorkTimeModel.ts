import { autoImplement } from 'src/helper/functions';

export type OvertimeSubtractionModel = 'hours' | 'percentage';

// Weekday values: 0=Sunday, 1=Monday, ..., 6=Saturday
export type WeekdayExceptionMap = Record<number, number>;

export interface ApiWorkTimeModel {
  ID: number;
  CreatedAt: Date;
  Name: string;
  WorkingHoursPerWeekday: number;
  OvertimeSubtractionModel: OvertimeSubtractionModel;
  OvertimeSubtractionAmount: number;
  HoursPerWeekdayException: WeekdayExceptionMap | null;
  HolidayDaysPerYear: number;
  OvertimeWarningThreshold: number;
}

export interface ApiWorkTimeModelCreateRequest {
  Name: string;
  WorkingHoursPerWeekday: number;
  OvertimeSubtractionModel: OvertimeSubtractionModel;
  OvertimeSubtractionAmount: number;
  HoursPerWeekdayException: WeekdayExceptionMap | null;
  HolidayDaysPerYear: number;
  OvertimeWarningThreshold: number;
}

export interface ApiWorkTimeModelUpdateRequest {
  WorkingHoursPerWeekday: number;
  OvertimeSubtractionModel: OvertimeSubtractionModel;
  OvertimeSubtractionAmount: number;
  HoursPerWeekdayException: WeekdayExceptionMap | null;
  HolidayDaysPerYear: number;
  OvertimeWarningThreshold: number;
}

export class WorkTimeModel extends autoImplement<ApiWorkTimeModel>() {
  static fromApi(apiItem: ApiWorkTimeModel) {
    return new WorkTimeModel(apiItem);
  }
}

export interface ApiUserWorkTime {
  ID: number;
  CreatedAt: Date;
  UserID: number;
  WorkTimeModelID: number;
  WorkTimeModel: ApiWorkTimeModel;
  ValidFrom: Date;
  ValidTill: Date | null;
}

export interface ApiUserWorkTimeCreateRequest {
  WorkTimeModelID: number;
  ValidFrom: Date;
  ValidTill: Date | null;
}

export interface ApiUserWorkTimeUpdateRequest {
  ValidFrom: Date;
  ValidTill: Date | null;
}

export class UserWorkTime extends autoImplement<ApiUserWorkTime>() {
  static fromApi(apiItem: ApiUserWorkTime) {
    return new UserWorkTime(apiItem);
  }
}
