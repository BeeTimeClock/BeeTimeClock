export function formatIndustryHourMinutes(hourMinutes: number) : string {
  let mins = hourMinutes;
  if (hourMinutes < 0) {
    mins = mins * -1;
  }

  const realHours = Math.floor(mins);
  const realMinutes = (mins - realHours) * 100 * 0.6;

  return `${hourMinutes.toFixed(2)} (${hourMinutes < 0 ? '-' : ''}${realHours}:${(realMinutes).toFixed(0)})`
}
