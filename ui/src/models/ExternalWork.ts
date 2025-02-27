import { autoImplement } from 'src/helper/functions';

export interface ApiExternalWorkExpanse {
  ID: number;
  ExternalWorkID: number;
  Date: Date;
  DepartureTime: Date;
  ArrivalTime: Date;
  TravelDurationHours: number;
  PauseDurationHours: number;
  OnSiteFrom: Date;
  OnSiteTill: Date;
  Place: string;
}

export class ExternalWorkExpanse extends autoImplement<ApiExternalWorkExpanse>() {
  static fromApi(apiItem: ApiExternalWorkExpanse) : ExternalWorkExpanse {
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
}

export interface ApiExternalWorkCreateRequest {
  From: Date;
  Till: Date;
  Description: string;
}

export class ExternalWork extends autoImplement<ApiExternalWork>() {
  private mappedExpanses: ExternalWorkExpanse[] = [];

  static fromApi(apiItem: ApiExternalWork) : ExternalWork {
    const externalWork = new ExternalWork(apiItem);
    if (externalWork.WorkExpanses != null) {
      externalWork.mappedExpanses = externalWork.WorkExpanses.map(s => ExternalWorkExpanse.fromApi(s))
    }
    return externalWork;
  }

  get WorkExpansesMapped() {
    return this.mappedExpanses;
  }
}
