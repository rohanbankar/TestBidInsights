import React, { useState } from 'react';
import { ColumnDef } from '@tanstack/react-table';
import { Layout } from '../components/layout/Layout';
import { DatePicker } from '../components/ui/DatePicker';
import { Select } from '../components/ui/Select';
import { Button } from '../components/ui/Button';
import { DataTable } from '../components/tables/DataTable';
import { BarChart } from '../components/charts/BarChart';
import { HeatMap } from '../components/charts/HeatMap';
import { useVideoHealth } from '../hooks/useReports';
import { formatNumber, formatPercentage, formatDisplayDate } from '../utils/formatUtils';
import { getDateRangePresets } from '../utils/dateUtils';
import { VideoHealth as VideoHealthType } from '../types/reports';
import { Download, RefreshCw, Play } from 'lucide-react';

const platformOptions = [
  { value: 'CTV', label: 'CTV' },
  { value: 'Display', label: 'Display' },
  { value: 'App', label: 'App' },
];

export function VideoHealth() {
  const datePresets = getDateRangePresets();
  const [platform, setPlatform] = useState('CTV');
  const [startDate, setStartDate] = useState(datePresets.last7Days.start);
  const [endDate, setEndDate] = useState(datePresets.last7Days.end);

  const { data, isLoading, error, refetch } = useVideoHealth(platform, startDate, endDate);

  const columns: ColumnDef<VideoHealthType>[] = [
    {
      accessorKey: 'date',
      header: 'Date',
      cell: ({ row }) => formatDisplayDate(row.getValue('date')),
    },
    {
      accessorKey: 'percentCtv',
      header: 'CTV %',
      cell: ({ row }) => formatPercentage(row.getValue('percentCtv')),
    },
    {
      accessorKey: 'placement',
      header: 'Placement',
      cell: ({ row }) => formatNumber(row.getValue('placement')),
    },
    {
      accessorKey: 'protocols',
      header: 'Protocols',
      cell: ({ row }) => formatNumber(row.getValue('protocols')),
    },
    {
      accessorKey: 'linearity',
      header: 'Linearity',
      cell: ({ row }) => formatNumber(row.getValue('linearity')),
    },
    {
      accessorKey: 'skip',
      header: 'Skip',
      cell: ({ row }) => formatNumber(row.getValue('skip')),
    },
    {
      accessorKey: 'startDelay',
      header: 'Start Delay',
      cell: ({ row }) => formatNumber(row.getValue('startDelay')),
    },
    {
      accessorKey: 'minDuration',
      header: 'Min Duration',
      cell: ({ row }) => formatNumber(row.getValue('minDuration')),
    },
    {
      accessorKey: 'maxDuration',
      header: 'Max Duration',
      cell: ({ row }) => formatNumber(row.getValue('maxDuration')),
    },
  ];

  const handleExportCSV = () => {
    if (!data?.data) return;

    const csvContent = [
      columns.map(col => col.header).join(','),
      ...data.data.map(row =>
        columns.map(col => {
          const value = row[col.accessorKey as keyof VideoHealthType];
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
    a.download = `video-health-${platform}-${startDate}-to-${endDate}.csv`;
    a.click();
    window.URL.revokeObjectURL(url);
  };

  // Prepare chart data
  const chartData = data?.data.map(item => ({
    date: item.date,
    placement: item.placement,
    protocols: item.protocols,
    linearity: item.linearity,
    skip: item.skip,
  })) || [];

  // Prepare heatmap data for protocol/placement correlation
  const heatMapData = data?.data.flatMap((item, index) => [
    { x: 'Protocols', y: `Day ${index + 1}`, value: item.protocols },
    { x: 'Placement', y: `Day ${index + 1}`, value: item.placement },
    { x: 'Skip', y: `Day ${index + 1}`, value: item.skip },
    { x: 'Linearity', y: `Day ${index + 1}`, value: item.linearity },
  ]) || [];

  return (
    <Layout title="Video Health">
      <div className="space-y-6">
        {/* Filters */}
        <div className="bg-white p-6 rounded-lg shadow">
          <div className="flex items-center space-x-4 mb-4">
            <Play className="w-5 h-5 text-gray-400" />
            <h3 className="text-lg font-medium text-gray-900">Video Health Filters</h3>
          </div>
          
          <div className="grid grid-cols-1 md:grid-cols-5 gap-4 items-end">
            <Select
              label="Platform"
              value={platform}
              onChange={(e) => setPlatform(e.target.value)}
              options={platformOptions}
            />
            
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
        {chartData.length > 0 && (
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
            <div className="bg-white p-6 rounded-lg shadow">
              <BarChart
                data={chartData}
                xKey="date"
                yKeys={['placement', 'protocols', 'skip', 'linearity']}
                title={`Video Properties - ${platform}`}
                colors={['#3b82f6', '#10b981', '#f59e0b', '#ef4444']}
              />
            </div>

            {heatMapData.length > 0 && (
              <div className="bg-white p-6 rounded-lg shadow">
                <HeatMap
                  data={heatMapData}
                  title="Video Properties Heatmap"
                  colorScheme="blue"
                />
              </div>
            )}
          </div>
        )}

        {/* Key Metrics */}
        {data?.data && data.data.length > 0 && (
          <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
            <div className="bg-white p-4 rounded-lg shadow">
              <div className="text-2xl font-bold text-primary-600 mb-1">
                {formatPercentage(
                  data.data.reduce((sum, item) => sum + item.percentCtv, 0) / data.data.length
                )}
              </div>
              <div className="text-sm text-gray-600">Avg CTV %</div>
            </div>
            
            <div className="bg-white p-4 rounded-lg shadow">
              <div className="text-2xl font-bold text-green-600 mb-1">
                {formatNumber(
                  data.data.reduce((sum, item) => sum + item.protocols, 0)
                )}
              </div>
              <div className="text-sm text-gray-600">Total Protocols</div>
            </div>
            
            <div className="bg-white p-4 rounded-lg shadow">
              <div className="text-2xl font-bold text-yellow-600 mb-1">
                {formatNumber(
                  data.data.reduce((sum, item) => sum + item.skip, 0)
                )}
              </div>
              <div className="text-sm text-gray-600">Skip Events</div>
            </div>
            
            <div className="bg-white p-4 rounded-lg shadow">
              <div className="text-2xl font-bold text-red-600 mb-1">
                {Math.round(
                  data.data.reduce((sum, item) => sum + item.maxDuration, 0) / data.data.length
                )}s
              </div>
              <div className="text-sm text-gray-600">Avg Max Duration</div>
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
              searchPlaceholder="Search video health data..."
            />
          )}
        </div>
      </div>
    </Layout>
  );
}