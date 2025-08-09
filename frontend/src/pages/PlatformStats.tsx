import React, { useState } from 'react';
import { ColumnDef } from '@tanstack/react-table';
import { Layout } from '../components/layout/Layout';
import { DatePicker } from '../components/ui/DatePicker';
import { Button } from '../components/ui/Button';
import { DataTable } from '../components/tables/DataTable';
import { BarChart } from '../components/charts/BarChart';
import { LineChart } from '../components/charts/LineChart';
import { usePlatformStats } from '../hooks/useReports';
import { formatNumber, formatPercentage, formatDisplayDate } from '../utils/formatUtils';
import { getDateRangePresets } from '../utils/dateUtils';
import { PlatformStats as PlatformStatsType } from '../types/reports';
import { Download, RefreshCw, Calendar } from 'lucide-react';

export function PlatformStats() {
  const datePresets = getDateRangePresets();
  const [startDate, setStartDate] = useState(datePresets.last7Days.start);
  const [endDate, setEndDate] = useState(datePresets.last7Days.end);

  const { data, isLoading, error, refetch } = usePlatformStats(startDate, endDate);

  const columns: ColumnDef<PlatformStatsType>[] = [
    {
      accessorKey: 'date',
      header: 'Date',
      cell: ({ row }) => formatDisplayDate(row.getValue('date')),
    },
    {
      accessorKey: 'totalRequests',
      header: 'Total Requests',
      cell: ({ row }) => formatNumber(row.getValue('totalRequests')),
    },
    {
      accessorKey: 'multiImpression',
      header: 'Multi Impression',
      cell: ({ row }) => formatNumber(row.getValue('multiImpression')),
    },
    {
      accessorKey: 'bidRate',
      header: 'Bid Rate',
      cell: ({ row }) => formatPercentage(row.getValue('bidRate')),
    },
    {
      accessorKey: 'timeoutRate',
      header: 'Timeout Rate',
      cell: ({ row }) => formatPercentage(row.getValue('timeoutRate')),
    },
    {
      accessorKey: 'deals',
      header: 'Deals',
      cell: ({ row }) => formatNumber(row.getValue('deals')),
    },
    {
      accessorKey: 'invalidRequests',
      header: 'Invalid Requests',
      cell: ({ row }) => formatNumber(row.getValue('invalidRequests')),
    },
  ];

  const handleExportCSV = () => {
    if (!data?.data) return;

    const csvContent = [
      columns.map(col => col.header).join(','),
      ...data.data.map(row =>
        columns.map(col => {
          const value = row[col.accessorKey as keyof PlatformStatsType];
          return typeof value === 'string' && value.includes(',') 
            ? `"${value}"` 
            : value;
        }).join(',')
      )
    ].join('\n');

    const blob = new Blob([csvContent], { type: 'text/csv' });
    const url = window.URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = `platform-stats-${startDate}-to-${endDate}.csv`;
    a.click();
    window.URL.revokeObjectURL(url);
  };

  return (
    <Layout title="Platform Statistics">
      <div className="space-y-6">
        {/* Filters */}
        <div className="bg-white p-6 rounded-lg shadow">
          <div className="flex items-center space-x-4 mb-4">
            <Calendar className="w-5 h-5 text-gray-400" />
            <h3 className="text-lg font-medium text-gray-900">Date Range</h3>
          </div>
          
          <div className="grid grid-cols-1 md:grid-cols-4 gap-4 items-end">
            <DatePicker
              label="Start Date"
              value={startDate}
              onChange={setStartDate}
            />
            <DatePicker
              label="End Date"
              value={endDate}
              onChange={setEndDate}
              min={startDate}
            />
            
            <div className="flex space-x-2">
              <Button onClick={() => refetch()} size="sm" variant="outline">
                <RefreshCw className="w-4 h-4 mr-1" />
                Refresh
              </Button>
              <Button onClick={handleExportCSV} size="sm" variant="outline">
                <Download className="w-4 h-4 mr-1" />
                Export
              </Button>
            </div>
          </div>

          {/* Quick Presets */}
          <div className="mt-4 flex flex-wrap gap-2">
            {Object.entries(datePresets).map(([key, preset]) => (
              <Button
                key={key}
                variant="ghost"
                size="sm"
                onClick={() => {
                  setStartDate(preset.start);
                  setEndDate(preset.end);
                }}
              >
                {preset.label}
              </Button>
            ))}
          </div>
        </div>

        {/* Charts */}
        {data?.data && data.data.length > 0 && (
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
            <div className="bg-white p-6 rounded-lg shadow">
              <BarChart
                data={data.data}
                xKey="date"
                yKeys={['totalRequests', 'multiImpression']}
                title="Requests Over Time"
                colors={['#3b82f6', '#10b981']}
              />
            </div>
            
            <div className="bg-white p-6 rounded-lg shadow">
              <LineChart
                data={data.data}
                xKey="date"
                yKeys={['bidRate', 'timeoutRate']}
                title="Performance Rates"
                colors={['#10b981', '#ef4444']}
                formatAsPercentage={true}
              />
            </div>
          </div>
        )}

        {/* Data Table */}
        <div className="bg-white p-6 rounded-lg shadow">
          {isLoading ? (
            <div className="flex items-center justify-center h-32">
              <RefreshCw className="w-8 h-8 animate-spin text-primary-600" />
            </div>
          ) : error ? (
            <div className="text-center py-8">
              <p className="text-red-600 mb-4">
                {error instanceof Error ? error.message : 'Failed to load data'}
              </p>
              <Button onClick={() => refetch()} variant="outline">
                Try Again
              </Button>
            </div>
          ) : (
            <DataTable
              data={data?.data || []}
              columns={columns}
              searchPlaceholder="Search platform stats..."
            />
          )}
        </div>
      </div>
    </Layout>
  );
}