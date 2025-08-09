import { format, subDays, startOfWeek, endOfWeek, startOfMonth, endOfMonth } from 'date-fns';

export const formatDate = (date: Date): string => {
  return format(date, 'yyyy-MM-dd');
};

export const formatDisplayDate = (dateString: string): string => {
  return format(new Date(dateString), 'MMM dd, yyyy');
};

export const getDateRangePresets = () => {
  const today = new Date();
  
  return {
    today: {
      start: formatDate(today),
      end: formatDate(today),
      label: 'Today',
    },
    yesterday: {
      start: formatDate(subDays(today, 1)),
      end: formatDate(subDays(today, 1)),
      label: 'Yesterday',
    },
    last7Days: {
      start: formatDate(subDays(today, 7)),
      end: formatDate(today),
      label: 'Last 7 days',
    },
    last30Days: {
      start: formatDate(subDays(today, 30)),
      end: formatDate(today),
      label: 'Last 30 days',
    },
    thisWeek: {
      start: formatDate(startOfWeek(today)),
      end: formatDate(endOfWeek(today)),
      label: 'This week',
    },
    thisMonth: {
      start: formatDate(startOfMonth(today)),
      end: formatDate(endOfMonth(today)),
      label: 'This month',
    },
  };
};