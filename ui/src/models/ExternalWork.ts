import { autoImplement } from 'src/helper/functions';

export interface ApiExternalWorkCompensationHourSlot {
  Hours: number;
  Compensation: number;
}

export interface ExternalWorkCompensationAdditionalOptions {
  [details: string]: number;
}

export interface ApiExternalWorkCompensation {
  ID: number;
  IsoCountryCodeA2: string;
  WithSocialInsuranceSlots: ApiExternalWorkCompensationHourSlot[];
  WithoutSocialInsuranceSlots: ApiExternalWorkCompensationHourSlot[];
  AdditionalOptions: ExternalWorkCompensationAdditionalOptions;
  ValidFrom: Date;
  ValidTill: Date;
}

export class ExternalWorkCompensation extends autoImplement<ApiExternalWorkCompensation>() {
  static fromApi(apiItem: ApiExternalWorkCompensation) : ExternalWorkCompensation {
    return new ExternalWorkCompensation(apiItem);
  }
}

export interface ApiExternalWorkExpanse {
  ID: number;
  ExternalWorkID: number;
  Date: Date;
  DepartureTime?: Date;
  ArrivalTime: Date;
  TravelDurationHours: number;
  PauseDurationHours: number;
  RestDurationHours: number;
  TravelWithPrivateCarKm: number;
  OnSiteFrom: Date;
  OnSiteTill: Date;
  OnSiteDuration: number;
  Place: string;
  TotalWorkingHours: number;
  TotalOperationHours: number;
  TotalOvertimeHours: number;
  TotalAwayHours: number;
  ExpensesWithoutSocialInsurance: number;
  ExpensesWithSocialInsurance: number;
  ExternalWorkCompensationID: number;
  ExternalWorkCompensation: ApiExternalWorkCompensation;
}

export class ExternalWorkExpanse extends autoImplement<ApiExternalWorkExpanse>() {
  static fromApi(apiItem: ApiExternalWorkExpanse): ExternalWorkExpanse {
    return new ExternalWorkExpanse(apiItem);
  }
}

export interface ApiExternalWork {
  ID: number;
  UserID: number;
  From: Date;
  Till: Date;
  Description: string;
  WorkExpanses: ApiExternalWorkExpanse[];
  TotalOvertimeHours: number;
  TotalExpensesWithoutSocialInsurance: number;
  TotalExpensesWithSocialInsurance: number;
}

export interface ApiExternalWorkCreateRequest {
  From: Date;
  Till: Date;
  Description: string;
  ExternalWorkCompensationID: number;
}

export class ExternalWork extends autoImplement<ApiExternalWork>() {
  private mappedExpanses: ExternalWorkExpanse[] = [];

  static fromApi(apiItem: ApiExternalWork): ExternalWork {
    const externalWork = new ExternalWork(apiItem);
    if (externalWork.WorkExpanses != null) {
      externalWork.mappedExpanses = externalWork.WorkExpanses.map((s) =>
        ExternalWorkExpanse.fromApi(s)
      );
    }
    return externalWork;
  }

  get WorkExpansesMapped() {
    return this.mappedExpanses;
  }
}
