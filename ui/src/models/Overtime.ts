import { autoImplement } from 'src/helper/functions';

export interface ApiOvertimeSummaryEntry {
  Source: string;
  Identifier: number;
  Value: number;
}

export interface ApiOvertimeMonthQuota {
  ID: number;
  UserID: number;
  Year: number;
  Month: number;
  Hours: number;
  Summary: ApiOvertimeSummaryEntry[];
}

export class OvertimeMonthQuota extends autoImplement<ApiOvertimeMonthQuota>() {
  static fromApi(item: ApiOvertimeMonthQuota) : OvertimeMonthQuota {
    return new OvertimeMonthQuota(item);
  }
}
