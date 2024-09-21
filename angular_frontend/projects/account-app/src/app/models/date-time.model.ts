export type Day = string;   // e.g., "4"
export type Month = string; // e.g., "August"
export type Year = string;  // e.g., "2024"
export type Time = string;  // e.g., "09:00 PM"

export interface DateTime {
  day: Day;
  month: Month;
  year: Year;
  time: Time;
  location: string;
}
