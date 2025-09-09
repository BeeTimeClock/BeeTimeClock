import { autoImplement } from 'src/helper/functions';

export enum ApiExternalWorkStatus {
  Planned = 'planned',
  MissingInfo = 'missing_info',
  InReview = 'in_review',
  Accepted = 'accepted',
  Declined = 'declined',
  Invoiced = 'invoiced'
}

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
  PrivateCarKmCompensation: number;
}

export class ExternalWorkCompensation extends autoImplement<ApiExternalWorkCompensation>() {
  static fromApi(
    apiItem: ApiExternalWorkCompensation
  ): ExternalWorkCompensation {
    return new ExternalWorkCompensation(apiItem);
  }
}

export interface ApiExternalWorkExpanse {
  ID: number;
  ExternalWorkID: number;
  ExternalWork: ApiExternalWork;
  Date: Date;
  DepartureTime?: Date | undefined;
  ArrivalTime?: Date | undefined;
  TravelDurationHours: number;
  PauseDurationHours: number;
  RestDurationHours: number;
  TravelWithPrivateCarKm: number;
  OnSiteFrom?: Date|undefined;
  OnSiteTill?: Date|undefined;
  OnSiteDuration: number;
  Place: string;
  TotalWorkingHours: number;
  TotalOperationHours: number;
  TotalOvertimeHours: number;
  TotalAwayHours: number;
  ExpensesWithoutSocialInsurance: number;
  ExpensesWithSocialInsurance: number;
  AdditionalOptions: string[];
}

export class ExternalWorkExpanse extends autoImplement<ApiExternalWorkExpanse>() {
  static fromApi(apiItem: ApiExternalWorkExpanse): ExternalWorkExpanse {
    return new ExternalWorkExpanse(apiItem);
  }

  public get OnSiteFromDate() {
    return this.OnSiteFrom ? new Date(this.OnSiteFrom) : undefined;
  }

  public set OnSiteFromDate(input: Date|undefined) {
    this.OnSiteFrom = input;
  }

  public get OnSiteTillDate() {
    return this.OnSiteTill ? new Date(this.OnSiteTill) : undefined;
  }

  public set OnSiteTillDate(input: Date|undefined) {
    this.OnSiteTill = input;
  }

  public get ArrivalTimeDate() {
    return this.ArrivalTime ? new Date(this.ArrivalTime) : undefined;
  }

  public set ArrivalTimeDate(input: Date|undefined) {
    this.ArrivalTime = input;
  }

  public get DepartureTimeDate() {
    return this.DepartureTime ? new Date(this.DepartureTime) : undefined;
  }

  public set DepartureTimeDate(input: Date|undefined) {
    this.DepartureTime = input;
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
  TotalOptions: ExternalWorkCompensationAdditionalOptions;
  ExternalWorkCompensationID: number;
  ExternalWorkCompensation: ApiExternalWorkCompensation;
  Status: ApiExternalWorkStatus;
  InvoiceDate: Date;
  IsLocked: boolean;
}

export interface ApiExternalWorkCreateRequest {
  From: string;
  Till: string;
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

  get StatusLabel() {
    switch (this.Status) {
      case ApiExternalWorkStatus.Planned:
        return 'LABEL_PLANNED'
      case ApiExternalWorkStatus.MissingInfo:
        return 'LABEL_MISSING_INFO'
      case ApiExternalWorkStatus.InReview:
        return 'LABEL_IN_REVIEW'
      case ApiExternalWorkStatus.Accepted:
        return 'LABEL_ACCEPTED'
      case ApiExternalWorkStatus.Declined:
        return 'LABEL_DECLINED'
      case ApiExternalWorkStatus.Invoiced:
        return 'LABEL_INVOICED'
    }
  }

  get NeedsUserInput() {
    return this.Status == ApiExternalWorkStatus.MissingInfo || this.Status == ApiExternalWorkStatus.Planned
  }
}

export interface ApiExternalWorkInvoicedInfo {
  Identifier: string;
  InvoiceDate: Date;
}
