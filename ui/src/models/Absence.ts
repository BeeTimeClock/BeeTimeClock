import {User} from 'src/models/Authentication';

export interface AbsenceCreateRequest {
  AbsenceFrom: Date;
  AbsenceTill: Date;
  AbsenceReasonID: number;
}

export interface AbsenceReason {
  ID: number;
  Description: string;
}

export interface Absence {
  ID: number;
  AbsenceFrom: Date;
  AbsenceTill: Date;
  AbsenceReasonID: number;
  CreatedAt: Date;
  UpdatedAt: Date;
  NettoDays: number;
}

export interface AbsenceSummaryItem {
  ID: number;
  AbsenceFrom: Date;
  AbsenceTill: Date;
  NettoDays: number;
  User: User;
}

type AbsenceUserSummaryYearReason = {
  Upcoming: number;
  Past: number;
}

type AbsenceUserSummaryYearReasonMap = {
  [key: number]: AbsenceUserSummaryYearReason;
}

export type AbsenceUserSummaryYear = {
  ByAbsenceReason: AbsenceUserSummaryYearReasonMap
}

export type AbsenceUserSummaryYearMap = {
  [key: number]: AbsenceUserSummaryYear;
}

export type AbsenceUserSummary = {
  ByYear: AbsenceUserSummaryYearMap;
  HolidayDaysPerYear: number;
}
