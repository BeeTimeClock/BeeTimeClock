import { autoImplement } from 'src/helper/functions';

export interface ApiHoliday {
  ID: number;
  Name: string;
  Date:  string;
  State: string;
}

export class Holiday extends autoImplement<ApiHoliday>() {
  static fromApi(apiItem: ApiHoliday) : Holiday {
    return new Holiday(apiItem);
  }
}

export interface ApiHolidayCustom {
  ID: number;
  Name: string;
  Date?: string;
  Month?: number;
  Day?: number;
  Yearly?: boolean;
}

export class HolidayCustom extends autoImplement<ApiHolidayCustom>() {
  static fromApi(apiItem: ApiHolidayCustom) : HolidayCustom {
    return new HolidayCustom(apiItem);
  }
}

