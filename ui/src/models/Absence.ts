import type {User} from 'src/models/Authentication';
import { autoImplement } from 'src/helper/functions';
import { date } from 'quasar';

export interface AbsenceCreateRequest {
  AbsenceFrom: string;
  AbsenceTill: string;
  AbsenceReasonID: number;
}

export enum AbsenceReasonImpact {
  Duration = 'duration',
  Hours = 'hours'
}

export interface ApiAbsenceReason {
  ID: number;
  Description: string;
  Impact?: AbsenceReasonImpact;
  ImpactHours: number;
}

export class AbsenceReason extends autoImplement<ApiAbsenceReason>() {
  static fromApi(apiItem: ApiAbsenceReason) : AbsenceReason {
    return new AbsenceReason(apiItem);
  }
}

export interface AbsenceReasonCreateRequest {
  Description: string;
}

export interface ApiAbsence {
  ID: number;
  AbsenceFrom: Date;
  AbsenceTill: Date;
  AbsenceReasonID: number;
  CreatedAt: Date;
  UpdatedAt: Date;
  NettoDays: number;
  Deletable: boolean;
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

export class Absence extends autoImplement<ApiAbsence>() {
  static fromApi(apiItem: ApiAbsence) : Absence {
    return new Absence(apiItem);
  }

  get formatTill() {
    return date.formatDate(this.AbsenceTill, 'DD. MMM. YYYY')
  }

  get formatFrom() {
    return date.formatDate(this.AbsenceFrom, 'DD. MMM. YYYY')
  }

  get formatTillFull() {
    return date.formatDate(this.AbsenceTill, 'ddd DD. MMM. YYYY')
  }

  get formatFromFull() {
    return date.formatDate(this.AbsenceFrom, 'ddd DD. MMM. YYYY')
  }
}
