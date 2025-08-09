import React, { useState } from 'react';
import { ColumnDef } from '@tanstack/react-table';
import { Layout } from '../components/layout/Layout';
import { DatePicker } from '../components/ui/DatePicker';
import { Select } from '../components/ui/Select';
import { Button } from '../components/ui/Button';
import { DataTable } from '../components/tables/DataTable';
import { BarChart } from '../components/charts/BarChart';
import { useContentHealth } from '../hooks/useReports';
import { formatNumber, formatDisplayDate } from '../utils/formatUtils';
import { getDateRangePresets } from '../utils/dateUtils';
import { ContentHealth as ContentHealthType } from '../types/reports';
import { Download, RefreshCw, FileText } from 'lucide-react';

const platformOptions = [
  { value: 'CTV', label: 'CTV' },
  { value: 'Audio', label: 'Audio' },
];

export function ContentHealth() {
  const datePresets = getDateRangePresets();
  const [platform, setPlatform] = useState('CTV');
  const [startDate, setStartDate] = useState(datePresets.last7Days.start);
  const [endDate, setEndDate] = useState(datePresets.last7Days.end);

  const { data, isLoading, error, refetch } = useContentHealth(platform, startDate, endDate);

  const columns: ColumnDef<ContentHealthType>[] = [
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
      accessorKey: 'title',
      header: 'Title',
      cell: ({ row }) => formatNumber(row.getValue('title')),
    },
    {
      accessorKey: 'series',
      header: 'Series',
      cell: ({ row }) => formatNumber(row.getValue('series')),
    },
    {
      accessorKey: 'episode',
      header: 'Episode',
      cell: ({ row }) => formatNumber(row.getValue('episode')),
    },
    {
      accessorKey: 'genre',
      header: 'Genre',
      cell: ({ row }) => formatNumber(row.getValue('genre')),
    },
    {
      accessorKey: 'language',
      header: 'Language',
      cell: ({ row }) => formatNumber(row.getValue('language')),
    },
    {
      accessorKey: 'length',
      header: 'Length',
      cell: ({ row }) => formatNumber(row.getValue('length')),
    },
    {
      accessorKey: 'livestream',
      header: 'Live Stream',
      cell: ({ row }) => formatNumber(row.getValue('livestream')),
    },
  ];

  const handleExportCSV = () => {
    if (!data?.data) return;

    const csvContent = [
      columns.map(col => col.header).join(','),
      ...data.data.map(row =>
        columns.map(col => {
          const value = row[col.accessorKey as keyof ContentHealthType];
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
    a.download = `content-health-${platform}-${startDate}-to-${endDate}.csv`;
    a.click();
    window.URL.revokeObjectURL(url);
  };

  // Prepare chart data for content fields
  const chartData = data?.data.map(item => ({
    date: item.date,
    title: item.title,
    series: item.series,
    episode: item.episode,
    genre: item.genre,
    language: item.language,
  })) || [];

  return (
    <Layout title="Content Health">
      <div className="space-y-6">
        {/* Filters */}
        <div className="bg-white p-6 rounded-lg shadow">
          <div className="flex items-center space-x-4 mb-4">
            <FileText className="w-5 h-5 text-gray-400" />
            <h3 className="text-lg font-medium text-gray-900">Content Health Filters</h3>
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
          <div className="space-y-6">
            <div className="bg-white p-6 rounded-lg shadow">
              <BarChart
                data={chartData}
                xKey="date"
                yKeys={['title', 'series', 'episode', 'genre']}
                title={`Content Fields Distribution - ${platform}`}
                colors={['#3b82f6', '#10b981', '#f59e0b', '#ef4444']}
              />
            </div>
          </div>
        )}

        {/* Statistics Cards */}
        {data?.data && data.data.length > 0 && (
          <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
            {[
              { key: 'title', label: 'Titles', color: 'blue' },
              { key: 'series', label: 'Series', color: 'green' },
              { key: 'episode', label: 'Episodes', color: 'yellow' },
              { key: 'genre', label: 'Genres', color: 'red' },
            ].map(({ key, label, color }) => {
              const total = data.data.reduce((sum, item) => sum + (item[key as keyof ContentHealthType] as number), 0);
              
              return (
                <div key={key} className="bg-white p-4 rounded-lg shadow">
                  <div className={`text-2xl font-bold text-${color}-600 mb-1`}>
                    {formatNumber(total)}
                  </div>
                  <div className="text-sm text-gray-600">{label}</div>
                </div>
              );
            })}
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
              searchPlaceholder="Search content health data..."
            />
          )}
        </div>
      </div>
    </Layout>
  );
}