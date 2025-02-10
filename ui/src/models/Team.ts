import { ApiUser, User } from 'src/models/Authentication';
import { autoImplement } from 'src/helper/functions';

export interface ApiTeam {
  ID: number;
  Teamname: string;
  TeamOwnerID: number;
  TeamOwner: ApiUser;
  CreatedAt: Date;
}

export interface ApiTeamMember {
  ID: number;
  TeamID: number;
  Team: ApiTeam;
  UserID: number;
  User: ApiUser;
}

export interface ApiTeamCreateRequest {
  Teamname: string;
  TeamOwnerID: number;
}

export interface ApiTeamMemberCreateRequest {
  UserID: number;
}

export class Team extends autoImplement<ApiTeam>() {
  static fromApi(apiItem: ApiTeam) {
    return new Team(apiItem)
  }

  get teamOwnerMapped() {
    return User.fromApi(this.TeamOwner)
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
