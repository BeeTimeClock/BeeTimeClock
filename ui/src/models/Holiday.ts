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



