import {AuthRequest, AuthResponse, User} from 'src/models/Authentication';
import {AuthProviders, BackendStatus, BaseResponse, MicrosoftAuthSettings, UserApikey, UserApikeyCreateRequest} from 'src/models/Base';
import {api} from 'boot/axios';
import {AxiosResponse} from 'axios';
import {
  SumResponse,
  Timestamp,
  TimestampCorrectionCreateRequest,
  TimestampCreateRequest,
  TimestampGroup
} from 'src/models/Timestamp';
import {Absence, AbsenceCreateRequest, AbsenceReason, AbsenceSummaryItem, AbsenceUserSummary} from 'src/models/Absence';

class BeeTimeClock {
  login(username: string, password: string): Promise<AxiosResponse<BaseResponse<AuthResponse>>> {
    const authRequest = {
      Username: username,
      Password: password,
    } as AuthRequest;

    return api.get('/api/v1/auth', { params: authRequest });
  }

  timestampQueryCurrentMonthGrouped() : Promise<AxiosResponse<BaseResponse<TimestampGroup[]>>> {
    return api.get('/api/v1/timestamp/query/current_month/grouped');
  }

  timestampOvertime() : Promise<AxiosResponse<BaseResponse<SumResponse>>> {
    return api.get('/api/v1/timestamp/overtime');

  }

  timestampQueryCurrentMonthOvertime() : Promise<AxiosResponse<BaseResponse<SumResponse>>> {
    return api.get('/api/v1/timestamp/query/current_month/overtime');
  }

  timestampActionCheckin(isHomeoffice = false) : Promise<AxiosResponse<BaseResponse<Timestamp>>> {
    const timestampCreateRequest = {
      IsHomeoffice: isHomeoffice,
    } as TimestampCreateRequest;

    return api.post('/api/v1/timestamp/action/checkin',  timestampCreateRequest);
  }

  timestampCorrectionCreate(timestamp: Timestamp, timestampCorrectionCreateRequest: TimestampCorrectionCreateRequest) : Promise<AxiosResponse<BaseResponse<Timestamp>>> {
    return api.post(`/api/v1/timestamp/${timestamp.ID}/correction`, timestampCorrectionCreateRequest)
  }

  timestampCreate(timestampCreateRequest: TimestampCreateRequest) : Promise<AxiosResponse<BaseResponse<Timestamp>>> {
    return api.post('/api/v1/timestamp', timestampCreateRequest);
  }

  timestampActionCheckout() : Promise<AxiosResponse<BaseResponse<Timestamp>>> {
    return api.post('/api/v1/timestamp/action/checkout', {});
  }

  absenceReasons() : Promise<AxiosResponse<BaseResponse<AbsenceReason[]>>> {
    return api.get('/api/v1/absence/reasons');
  }

  createAbsence(absenceCreateRequest: AbsenceCreateRequest) : Promise<AxiosResponse<BaseResponse<Absence>>> {
    absenceCreateRequest.AbsenceFrom = new Date(absenceCreateRequest.AbsenceFrom)
    absenceCreateRequest.AbsenceTill = new Date(absenceCreateRequest.AbsenceTill)

    return api.post('/api/v1/absence', absenceCreateRequest);
  }

  getAbsences() : Promise<AxiosResponse<BaseResponse<Absence[]>>> {
    return api.get('/api/v1/absence')
  }

  getMeUser() : Promise<AxiosResponse<BaseResponse<User>>> {
    return api.get('/api/v1/user/me');
  }

  queryAbsenceSummary() : Promise<AxiosResponse<BaseResponse<AbsenceSummaryItem[]>>> {
    return api.get('/api/v1/absence/query/users/summary')
  }

  queryMyAbsenceSummary() : Promise<AxiosResponse<BaseResponse<AbsenceUserSummary>>> {
    return api.get('/api/v1/absence/query/me/summary')
  }

  administrationGetUsers(withData?: boolean) : Promise<AxiosResponse<BaseResponse<User[]>>> {
    const params = {
      with_data: withData,
    } ;
    
    return api.get('/api/v1/administration/user', { params: params });
  }

  administrationGetUserById(userId: string) : Promise<AxiosResponse<BaseResponse<User>>> {
    return api.get(`/api/v1/administration/user/${userId}`);
  }

  administrationUpdateUser(user: User) : Promise<AxiosResponse<BaseResponse<User>>> {
    return api.put(`/api/v1/administration/user/${user.ID}`, user);
  }

  getStatus() : Promise<AxiosResponse<BaseResponse<BackendStatus>>> {
    return api.get('/api/v1/status');
  }

  getAuthProviders() : Promise<AxiosResponse<BaseResponse<AuthProviders>>> {
    return api.get('/api/v1/auth/providers');
  }

  getMicrosoftAuthSettings() : Promise<AxiosResponse<BaseResponse<MicrosoftAuthSettings>>> {
    return api.get('/api/v1/auth/microsoft');
  }

  getUserApikey() : Promise<AxiosResponse<BaseResponse<UserApikey[]>>> {
    return api.get('/api/v1/user/me/apikey');
  }

  createUserApikey(userApikeyCreateRequest: UserApikeyCreateRequest) : Promise<AxiosResponse<BaseResponse<UserApikey>>> {
    return api.post('/api/v1/user/me/apikey', userApikeyCreateRequest);
  }
}

export default new BeeTimeClock();
