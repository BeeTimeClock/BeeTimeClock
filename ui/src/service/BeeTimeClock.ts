import type {
  ApiUser,
  AuthRequest,
  AuthResponse,
  User,
} from 'src/models/Authentication';
import type {
  ApiCountResult,
  AuthProviders,
  BackendStatus,
  BaseResponse,
  MicrosoftAuthSettings,
  UserApikey,
  UserApikeyCreateRequest,
} from 'src/models/Base';
import { api } from 'boot/axios';
import type { AxiosResponse } from 'axios';
import type {
  OvertimeResponse,
  SumResponse,
  Timestamp,
  TimestampCorrectionCreateRequest,
  TimestampCreateRequest,
  TimestampGroup,
  TimestampYearMonthGrouped,
} from 'src/models/Timestamp';
import type {
  AbsenceCreateRequest,
  AbsenceSignRequest,
  AbsenceSummaryItem,
  AbsenceUserSummary,
  ApiAbsence,
  ApiAbsenceReason,
} from 'src/models/Absence';
import type { Settings } from 'src/models/Settings';
import type {
  ApiTeam,
  ApiTeamCreateRequest,
  ApiTeamMember,
  ApiTeamMemberCreateRequest,
} from 'src/models/Team';
import type {
  ApiExternalWork,
  ApiExternalWorkCompensation,
  ApiExternalWorkCreateRequest,
  ApiExternalWorkExpanse,
  ApiExternalWorkInvoicedInfo,
  ExternalWorkCompensation,
} from 'src/models/ExternalWork';
import type { ApiOvertimeMonthQuota } from 'src/models/Overtime';
import type { ApiHoliday, ApiHolidayCustom } from 'src/models/Holiday';

class BeeTimeClock {
  login(
    username: string,
    password: string,
  ): Promise<AxiosResponse<BaseResponse<AuthResponse>>> {
    const authRequest = {
      Username: username,
      Password: password,
    } as AuthRequest;

    return api.get('/api/v1/auth', { params: authRequest });
  }

  timestampQueryCurrentMonthGrouped(): Promise<
    AxiosResponse<BaseResponse<TimestampGroup[]>>
  > {
    return api.get('/api/v1/timestamp/query/current_month/grouped');
  }

  timestampQueryMonthGrouped(
    year: number,
    month: number,
  ): Promise<AxiosResponse<BaseResponse<TimestampGroup[]>>> {
    return api.get(
      `/api/v1/timestamp/query/year/${year}/month/${month}/grouped`,
    );
  }

  timestampQueryCurrentMonthOvertime(): Promise<
    AxiosResponse<BaseResponse<SumResponse>>
  > {
    return api.get('/api/v1/timestamp/query/current_month/overtime');
  }

  timestampActionCheckin(
    isHomeoffice = false,
  ): Promise<AxiosResponse<BaseResponse<Timestamp>>> {
    const timestampCreateRequest = {
      IsHomeoffice: isHomeoffice,
    } as TimestampCreateRequest;

    return api.post('/api/v1/timestamp/action/checkin', timestampCreateRequest);
  }

  timestampCorrectionCreate(
    timestamp: Timestamp,
    timestampCorrectionCreateRequest: TimestampCorrectionCreateRequest,
  ): Promise<AxiosResponse<BaseResponse<Timestamp>>> {
    return api.post(
      `/api/v1/timestamp/${timestamp.ID}/correction`,
      timestampCorrectionCreateRequest,
    );
  }

  timestampCreate(
    timestampCreateRequest: TimestampCreateRequest,
  ): Promise<AxiosResponse<BaseResponse<Timestamp>>> {
    return api.post('/api/v1/timestamp', timestampCreateRequest);
  }

  timestampActionCheckout(
    isHomeoffice: boolean,
  ): Promise<AxiosResponse<BaseResponse<Timestamp>>> {
    const timestampCheckoutRequest = {
      IsHomeoffice: isHomeoffice,
    } as TimestampCreateRequest;

    return api.post('/api/v1/timestamp/action/checkout', {
      timestampCheckoutRequest,
    });
  }

  absenceReasons(): Promise<AxiosResponse<BaseResponse<ApiAbsenceReason[]>>> {
    return api.get('/api/v1/absence/reasons');
  }

  administrationUpdateAbsenceReason(
    absenceReasonId: number,
    absenceReason: ApiAbsenceReason,
  ): Promise<AxiosResponse<BaseResponse<ApiAbsenceReason>>> {
    return api.put(
      `/api/v1/administration/absence/reasons/${absenceReasonId}`,
      absenceReason,
    );
  }

  administrationCreateAbsenceReason(
    absenceReason: ApiAbsenceReason,
  ): Promise<AxiosResponse<BaseResponse<ApiAbsenceReason>>> {
    return api.post('/api/v1/administration/absence/reasons', absenceReason);
  }

  createAbsence(
    absenceCreateRequest: AbsenceCreateRequest,
  ): Promise<AxiosResponse<BaseResponse<ApiAbsence>>> {
    return api.post('/api/v1/absence', absenceCreateRequest);
  }

  createTeamUserAbsence(teamId: number, userId: number, absenceCreateRequest: AbsenceCreateRequest) : Promise<AxiosResponse<BaseResponse<ApiAbsence>>> {
    return api.post(`/api/v1/team/${teamId}/user/${userId}/absence`, absenceCreateRequest)
  }

  deleteAbsence(
    absenceId: number,
  ): Promise<AxiosResponse<BaseResponse<never>>> {
    return api.delete(`/api/v1/absence/${absenceId}`);
  }

  getAbsences(): Promise<AxiosResponse<BaseResponse<ApiAbsence[]>>> {
    return api.get('/api/v1/absence');
  }

  getMeUser(): Promise<AxiosResponse<BaseResponse<User>>> {
    return api.get('/api/v1/user/me');
  }

  updateMeUser(user: User): Promise<AxiosResponse<BaseResponse<User>>> {
    return api.put('/api/v1/user/me', user);
  }

  queryAbsenceSummary(): Promise<
    AxiosResponse<BaseResponse<AbsenceSummaryItem[]>>
  > {
    return api.get('/api/v1/absence/query/users/summary');
  }

  queryMyAbsenceSummary(): Promise<
    AxiosResponse<BaseResponse<AbsenceUserSummary>>
  > {
    return api.get('/api/v1/absence/query/me/summary');
  }

  getExternalWork(): Promise<AxiosResponse<BaseResponse<ApiExternalWork[]>>> {
    return api.get('/api/v1/external_work');
  }

  getExternalWorkInvoiced(): Promise<
    AxiosResponse<BaseResponse<ApiExternalWorkInvoicedInfo[]>>
  > {
    return api.get('/api/v1/external_work/invoiced');
  }

  getExternalWorkById(
    id: number,
  ): Promise<AxiosResponse<BaseResponse<ApiExternalWork>>> {
    return api.get(`/api/v1/external_work/${id}`);
  }

  deleteExternalWorkById(
    id: number,
  ): Promise<AxiosResponse<BaseResponse<ApiExternalWork>>> {
    return api.delete(`/api/v1/external_work/${id}`);
  }

  submitExternalWork(
    externalWorkId: number,
  ): Promise<AxiosResponse<BaseResponse<ApiExternalWork>>> {
    return api.post(`/api/v1/external_work/${externalWorkId}/action/submit`);
  }

  createExternalWork(
    externalWorkCreateRequest: ApiExternalWorkCreateRequest,
  ): Promise<AxiosResponse<BaseResponse<ApiExternalWork>>> {
    return api.post('/api/v1/external_work', externalWorkCreateRequest);
  }

  createExternalWorkExpanse(
    externalWorkId: number,
    externalWorkExpanse: ApiExternalWorkExpanse,
  ): Promise<AxiosResponse<BaseResponse<ApiExternalWorkExpanse>>> {
    return api.post(
      `/api/v1/external_work/${externalWorkId}/expanse`,
      externalWorkExpanse,
    );
  }

  updateExternalWorkExpanse(
    externalWorkId: number,
    externalWorkExpanseId: number,
    externalWorkExpanse: ApiExternalWorkExpanse,
  ): Promise<AxiosResponse<BaseResponse<ApiExternalWorkExpanse>>> {
    return api.put(
      `/api/v1/external_work/${externalWorkId}/expanse/${externalWorkExpanseId}`,
      externalWorkExpanse,
    );
  }

  administrationGetUsers(
    withData?: boolean,
  ): Promise<AxiosResponse<BaseResponse<ApiUser[]>>> {
    const params = {
      with_data: withData,
    };

    return api.get('/api/v1/administration/user', { params: params });
  }

  administrationGetTeams(
    withData?: boolean,
  ): Promise<AxiosResponse<BaseResponse<ApiTeam[]>>> {
    const params = {
      with_data: withData,
    };

    return api.get('/api/v1/administration/team', { params: params });
  }

  administrationCreateTeam(
    teamCreateRequest: ApiTeamCreateRequest,
  ): Promise<AxiosResponse<BaseResponse<ApiTeam>>> {
    return api.post('/api/v1/administration/team', teamCreateRequest);
  }

  administrationGetTeam(
    teamId: number,
    withData?: boolean,
  ): Promise<AxiosResponse<BaseResponse<ApiTeam>>> {
    const params = {
      with_data: withData,
    };
    return api.get(`/api/v1/administration/team/${teamId}`, { params: params });
  }

  administrationGetTeamMembers(
    teamId: number,
    withData?: boolean,
  ): Promise<AxiosResponse<BaseResponse<ApiTeamMember[]>>> {
    const params = {
      with_data: withData,
    };

    return api.get(`/api/v1/administration/team/${teamId}/member`, {
      params: params,
    });
  }

  administrationCreateTeamMember(
    teamId: number,
    teamMemberCreateRequest: ApiTeamMemberCreateRequest,
  ): Promise<AxiosResponse<BaseResponse<ApiTeamMember>>> {
    return api.post(
      `/api/v1/administration/team/${teamId}/member`,
      teamMemberCreateRequest,
    );
  }

  administrationDeleteTeamMember(
    teamId: number,
    teamMemberId: number,
  ): Promise<AxiosResponse<BaseResponse<never>>> {
    return api.delete(
      `/api/v1/administration/team/${teamId}/member/${teamMemberId}`,
    );
  }

  administrationGetUserById(
    userId: number,
  ): Promise<AxiosResponse<BaseResponse<User>>> {
    return api.get(`/api/v1/administration/user/${userId}`);
  }

  administrationUpdateUser(
    user: User,
  ): Promise<AxiosResponse<BaseResponse<User>>> {
    return api.put(`/api/v1/administration/user/${user.ID}`, user);
  }

  administrationDeleteUser(
    user: User,
  ): Promise<AxiosResponse<BaseResponse<never>>> {
    return api.delete(`/api/v1/administration/user/${user.ID}`);
  }

  administrationSummaryUserCurrentYear(
    userId: number,
    year: number,
  ): Promise<AxiosResponse<BaseResponse<AbsenceUserSummary>>> {
    return api.get(
      `/api/v1/administration/user/${userId}/absence/year/${year}/summary`,
    );
  }

  getStatus(): Promise<AxiosResponse<BaseResponse<BackendStatus>>> {
    return api.get('/api/v1/status');
  }

  getAuthProviders(): Promise<AxiosResponse<BaseResponse<AuthProviders>>> {
    return api.get('/api/v1/auth/providers');
  }

  getMicrosoftAuthSettings(): Promise<
    AxiosResponse<BaseResponse<MicrosoftAuthSettings>>
  > {
    return api.get('/api/v1/auth/microsoft');
  }

  getUserApikey(): Promise<AxiosResponse<BaseResponse<UserApikey[]>>> {
    return api.get('/api/v1/user/me/apikey');
  }

  createUserApikey(
    userApikeyCreateRequest: UserApikeyCreateRequest,
  ): Promise<AxiosResponse<BaseResponse<UserApikey>>> {
    return api.post('/api/v1/user/me/apikey', userApikeyCreateRequest);
  }

  timestampQueryMonths(): Promise<
    AxiosResponse<BaseResponse<TimestampYearMonthGrouped>>
  > {
    return api.get('/api/v1/timestamp/query/timestamp/months');
  }

  timestampQuerySuspicious(): Promise<
    AxiosResponse<BaseResponse<Timestamp[]>>
  > {
    return api.get('/api/v1/timestamp/query/suspicious');
  }

  administrationTimestampUserMonths(
    userId: number,
  ): Promise<AxiosResponse<BaseResponse<TimestampYearMonthGrouped>>> {
    return api.get(`/api/v1/administration/user/${userId}/timestamp/months`);
  }

  administrationTimestampQueryMonthGrouped(
    userId: number,
    year: number,
    month: number,
  ): Promise<AxiosResponse<BaseResponse<TimestampGroup[]>>> {
    return api.get(
      `/api/v1/administration/user/${userId}/timestamp/year/${year}/month/${month}/grouped`,
    );
  }

  administrationTimestampQueryMonthOvertime(
    userId: number,
    year: number,
    month: number,
  ): Promise<AxiosResponse<BaseResponse<OvertimeResponse>>> {
    return api.get(
      `/api/v1/administration/user/${userId}/timestamp/year/${year}/month/${month}/overtime`,
    );
  }

  administrationAbsenceYears(
    userId: number,
  ): Promise<AxiosResponse<BaseResponse<number[]>>> {
    return api.get(`/api/v1/administration/user/${userId}/absence/years`);
  }

  administrationAbsencesByYear(
    userId: number,
    year: number,
  ): Promise<AxiosResponse<BaseResponse<ApiAbsence[]>>> {
    return api.get(
      `/api/v1/administration/user/${userId}/absence/year/${year}`,
    );
  }

  timestampQueryMonthOvertime(
    year: number,
    month: number,
  ): Promise<AxiosResponse<BaseResponse<OvertimeResponse>>> {
    return api.get(
      `/api/v1/timestamp/query/year/${year}/month/${month}/overtime`,
    );
  }

  administrationSettings(): Promise<AxiosResponse<BaseResponse<Settings>>> {
    return api.get('/api/v1/administration/settings');
  }

  administrationSettingsSave(
    settings: Settings,
  ): Promise<AxiosResponse<BaseResponse<Settings>>> {
    return api.put('/api/v1/administration/settings', settings);
  }

  administrationNotifyAbsenceWeek(): Promise<
    AxiosResponse<BaseResponse<never>>
  > {
    return api.post('/api/v1/administration/notify/absence/week', {});
  }

  overtimeMonthQuotas(): Promise<
    AxiosResponse<BaseResponse<ApiOvertimeMonthQuota[]>>
  > {
    return api.get('/api/v1/overtime');
  }

  administrationOvertimeMonthQuotas(userId: number): Promise<
    AxiosResponse<BaseResponse<ApiOvertimeMonthQuota[]>>
  > {
    return api.get(`/api/v1/administration/user/${userId}/overtime`);
  }

  calculateOvertimeMonthQuota(
    year: number,
    month: number,
  ): Promise<AxiosResponse<BaseResponse<ApiOvertimeMonthQuota>>> {
    return api.post(`/api/v1/overtime/action/calculate/${year}/${month}`, {});
  }

  administrationCalculateOvertimeMonthQuota(userId: number,
    year: number,
    month: number,
  ): Promise<AxiosResponse<BaseResponse<ApiOvertimeMonthQuota>>> {
    return api.post(`/api/v1/administration/user/${userId}/overtime/action/calculate/${year}/${month}`, {});
  }

  overtimeTotal(): Promise<AxiosResponse<BaseResponse<SumResponse>>> {
    return api.get(`/api/v1/overtime/total`);
  }

  administrationOvertimeTotal(userId: number): Promise<AxiosResponse<BaseResponse<SumResponse>>> {
    return api.get(`/api/v1/administration/user/${userId}/overtime/total`);
  }

  externalWorkCompensation(): Promise<
    AxiosResponse<BaseResponse<ApiExternalWorkCompensation[]>>
  > {
    return api.get('/api/v1/external_work/compensation');
  }

  administrationExternalWorkCompensation(): Promise<
    AxiosResponse<BaseResponse<ApiExternalWorkCompensation[]>>
  > {
    return api.get('/api/v1/administration/external_work/compensation');
  }

  administrationExternalWorkCompensationUpdate(
    id: number,
    externalWorkCompensation: ExternalWorkCompensation,
  ): Promise<AxiosResponse<BaseResponse<ApiExternalWorkCompensation>>> {
    return api.put(
      `/api/v1/administration/external_work/compensation/${id}`,
      externalWorkCompensation,
    );
  }

  externalWorkDownloadPdf(): Promise<AxiosResponse<Blob>> {
    return api.get('/api/v1/external_work/action/export/pdf', {
      responseType: 'blob',
    });
  }

  externalWorkDownloadInvoicedPdf(
    identifier: string,
  ): Promise<AxiosResponse<Blob>> {
    return api.get(`/api/v1/external_work/action/export/pdf/${identifier}`, {
      responseType: 'blob',
    });
  }

  administrationUploadLogoFile(
    file: File | Blob,
  ): Promise<AxiosResponse<BaseResponse<never>>> {
    const fileData = new FormData();
    fileData.append('file', file);

    return api.post('/api/v1/administration/settings/logo', fileData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    });
  }

  administrationDebugHolidays() : Promise<AxiosResponse<BaseResponse<ApiHoliday[]>>> {
    return api.get('/api/v1/administration/debug/holidays')
  }

  getLogo(): Promise<AxiosResponse<Blob>> {
    return api.get('/api/v1/logo', {
      responseType: 'blob',
    });
  }

  getMissingDays(): Promise<AxiosResponse<BaseResponse<string[]>>> {
    return api.get(`/api/v1/timestamp/query/missing`)
  }

  getMissingDaysCount(): Promise<AxiosResponse<BaseResponse<ApiCountResult>>> {
    return api.get(`/api/v1/timestamp/query/missing/count`)
  }

  getMissingDaysMonth(year: number, month: number): Promise<AxiosResponse<BaseResponse<string[]>>> {
    return api.get(`/api/v1/timestamp/query/year/${year}/month/${month}/missing`)
  }

  getTeams(): Promise<AxiosResponse<BaseResponse<ApiTeam[]>>> {
    return api.get('/api/v1/team')
  }

  queryTeamAbsenceSummary(teamId: number): Promise<
    AxiosResponse<BaseResponse<AbsenceSummaryItem[]>>
  > {
    return api.get(`/api/v1/team/${teamId}/absence/query/users/summary`);
  }

  absenceTeamOpen(teamId: number) : Promise<AxiosResponse<BaseResponse<ApiAbsence[]>>> {
    return api.get(`/api/v1/team/${teamId}/absence/open`)
  }

  absenceTeamSign(teamId: number, absenceId: number, absenceSignRequest: AbsenceSignRequest) : Promise<AxiosResponse<BaseResponse<ApiAbsence>>> {
    return api.post(`/api/v1/team/${teamId}/absence/${absenceId}/sign`, absenceSignRequest)
  }

  administrationHolidaysCustom() : Promise<AxiosResponse<BaseResponse<ApiHolidayCustom[]>>> {
    return api.get(`/api/v1/administration/holidays/custom`)
  }

  holidaysYear(year: number) : Promise<AxiosResponse<BaseResponse<ApiHoliday[]>>> {
    return api.get(`/api/v1/holidays/year/${year}`)
  }

  administrationAbsenceRecalculate() : Promise<AxiosResponse<BaseResponse<never>>> {
    return api.post(`/api/v1/administration/absence/recalculate`)
  }
}

export default new BeeTimeClock();
