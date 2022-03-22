function isIsoDateString(value: any): boolean {
  return value && typeof value === "string";
}

export function handleDates(body: any) {
  // if (body === null || body === undefined || typeof body !== "object")
  //   return body;

  // for (const key of Object.keys(body)) {
  //   const value = body[key];
  //   if (isIsoDateString(value)) body[key] = parseISO(value);
  //   else if (typeof value === "object") handleDates(value);
  // }
}
