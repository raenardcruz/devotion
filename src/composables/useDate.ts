/**
 * Composable for date-related utilities using local time.
 */
export function useDate() {
  /**
   * Returns the current local date as an ISO string (YYYY-MM-DD).
   * This avoids the UTC offset issue with new Date().toISOString().
   */
  const getLocalISOString = (date: Date = new Date()): string => {
    const year = date.getFullYear();
    const month = String(date.getMonth() + 1).padStart(2, '0');
    const day = String(date.getDate()).padStart(2, '0');
    return `${year}-${month}-${day}`;
  };

  /**
   * Returns the day of the week (0-6) in local time.
   */
  const getLocalDay = (date: Date = new Date()): number => {
    return date.getDay();
  };

  return {
    getLocalISOString,
    getLocalDay
  };
}
