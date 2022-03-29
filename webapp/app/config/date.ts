import dayjs from 'dayjs';
import { DATE_FORMAT } from './constants';

export function handleDates(body: any) {
  if (body === null || body === undefined || typeof body !== "object")
    return body;

  for (const key of Object.keys(body)) {
    const value = body[key];
    if (key == "date") body[key] = dayjs(body[key], DATE_FORMAT);
    else if (typeof value === "object") handleDates(value);
  }

  return body;
}
