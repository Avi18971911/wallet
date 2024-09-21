import { Injectable } from '@angular/core';
import {DateTime} from "../../models/date-time.model";

@Injectable({
  providedIn: 'root'
})
export class DateFormatService {
  constructor() { }

  private monthNames = [
    'January', 'February', 'March', 'April', 'May', 'June',
    'July', 'August', 'September', 'October', 'November', 'December'
  ];

  private readonly DAY_REGEX = /^([1-9]|[12][0-9]|3[01])$/;
  private readonly MONTH_REGEX =
    /^(January|February|March|April|May|June|July|August|September|October|November|December)$/;
  private readonly YEAR_REGEX = /^\d{4}$/;
  private readonly TIME_REGEX = /^([1-9]|1[0-2]):[0-5][0-9]:[0-5][0-9] (AM|PM)$/;


  validateDateTime(dateTime: DateTime) {
    return this.DAY_REGEX.test(dateTime.day) &&
      this.MONTH_REGEX.test(dateTime.month) &&
      this.YEAR_REGEX.test(dateTime.year) &&
      this.TIME_REGEX.test(dateTime.time);
  }

  generateCurrentDateString(): DateTime {
    const currentDate = new Date();

    const formattedDate = currentDate.toLocaleDateString('en-US');
    const formattedTime = currentDate.toLocaleTimeString('en-US');
    const month = this.monthNames[currentDate.getMonth()];
    const year = currentDate.getFullYear().toString();
    const day = currentDate.getDate().toString();

    const generatedDateTime: DateTime = {
      day: day,
      month: month,
      year: year,
      time: formattedTime,
      location: 'Singapore'
    };

    if (!this.validateDateTime(generatedDateTime)) {
      console.log('Invalid date time generated: ', generatedDateTime);
    }

    return generatedDateTime;
  }

  getCurrentDate(): DateTime {
    return this.generateCurrentDateString();
  }
}
