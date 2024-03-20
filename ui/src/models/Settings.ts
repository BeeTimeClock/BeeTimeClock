export interface OfficeIPAddress {
  ID: number;
  IPAddress: string;
  Description: string;
}
export interface Settings {
  CheckinDetectionByIPAddress: boolean;
  OfficeIPAddresses: OfficeIPAddress[];
}
