import type { ApiUser} from 'src/models/Authentication';
import { User } from 'src/models/Authentication';
import { autoImplement } from 'src/helper/functions';

export enum TeamLevel {
  Lead = "lead",
  LeadSurrogate = "lead_surrogate",
  Member = "member"
}

export interface ApiTeam {
  ID: number;
  Teamname: string;
  CreatedAt: Date;
  Members: ApiTeamMember[];
}

export interface ApiTeamMember {
  ID: number;
  TeamID: number;
  Team: ApiTeam;
  UserID: number;
  User: ApiUser;
  Level: TeamLevel;
}

export interface ApiTeamCreateRequest {
  Teamname: string;
  TeamLeadID: number;
}

export interface ApiTeamMemberCreateRequest {
  UserID: number;
  Level: TeamLevel;
}

export class Team extends autoImplement<ApiTeam>() {
  static fromApi(apiItem: ApiTeam) {
    return new Team(apiItem)
  }

  get membersMapped() {
    return this.Members.map(s => TeamMember.fromApi(s))
  }
}

export class TeamMember extends autoImplement<ApiTeamMember>() {
  static fromApi(apiItem: ApiTeamMember) {
    return new TeamMember(apiItem);
  }

  get userMapped() {
    return User.fromApi(this.User)
  }
}
