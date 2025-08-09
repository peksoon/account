export function formatDateToString(date) {
    const year = date.getFullYear();
    const month = String(date.getMonth() + 1).padStart(2, '0');
    const day = String(date.getDate()).padStart(2, '0');
    return `${year}-${month}-${day}`;
  }

export function formatDateToLocalISOString(date) {
  const offset = date.getTimezoneOffset() * 60000;
  const DateOffset = new Date(date.getTime() - offset);
  return DateOffset.toISOString().slice(0, 19).replace('T', ' ')
}