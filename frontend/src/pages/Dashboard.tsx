import React from 'react';
import { Layout } from '../components/layout/Layout';
import { useDashboard } from '../hooks/useReports';
import { formatNumber, formatPercentage } from '../utils/formatUtils';
import { BarChart } from '../components/charts/BarChart';
import { LineChart } from '../components/charts/LineChart';
import { 
  TrendingUp, 
  Activity, 
  Users, 
  AlertCircle,
  RefreshCw
} from 'lucide-react';

interface StatCardProps {
  title: string;
  value: string | number;
  change?: string;
  icon: React.ReactNode;
  trend?: 'up' | 'down' | 'neutral';
}

function StatCard({ title, value, change, icon, trend = 'neutral' }: StatCardProps) {
  const trendColors = {
    up: 'text-green-600',
    down: 'text-red-600',
    neutral: 'text-gray-600',
  };

  return (
    <div className="bg-white overflow-hidden shadow rounded-lg">
      <div className="p-5">
        <div className="flex items-center">
          <div className="flex-shrink-0">
            <div className="text-gray-400">{icon}</div>
          </div>
          <div className="ml-5 w-0 flex-1">
            <dl>
              <dt className="text-sm font-medium text-gray-500 truncate">{title}</dt>
              <dd className="flex items-baseline">
                <div className="text-2xl font-semibold text-gray-900">
                  {typeof value === 'number' ? formatNumber(value) : value}
                </div>
                {change && (
                  <div className={`ml-2 flex items-baseline text-sm font-semibold ${trendColors[trend]}`}>
                    {change}
                  </div>
                )}
              </dd>
            </dl>
          </div>
        </div>
      </div>
    </div>
  );
}

export function Dashboard() {
  const { data, isLoading, error, refetch } = useDashboard();

  if (isLoading) {
    return (
      <Layout title="Dashboard">
        <div className="flex items-center justify-center h-64">
          <RefreshCw className="w-8 h-8 animate-spin text-primary-600" />
        </div>
      </Layout>
    );
  }

  if (error) {
    return (
      <Layout title="Dashboard">
        <div className="text-center py-12">
          <AlertCircle className="mx-auto h-12 w-12 text-red-400" />
          <h3 className="mt-2 text-sm font-medium text-gray-900">Error loading dashboard</h3>
          <p className="mt-1 text-sm text-gray-500">
            {error instanceof Error ? error.message : 'Something went wrong'}
          </p>
          <button
            onClick={() => refetch()}
            className="mt-4 bg-primary-600 text-white px-4 py-2 rounded-md text-sm hover:bg-primary-700"
          >
            Try again
          </button>
        </div>
      </Layout>
    );
  }

  const stats = data?.latestStats;
  const contentSummary = data?.contentSummary || {};
  const videoSummary = data?.videoSummary || {};

  // Prepare chart data
  const contentChartData = Object.entries(contentSummary).map(([platform, value]) => ({
    platform,
    requests: value,
  }));

  const videoChartData = Object.entries(videoSummary).map(([platform, value]) => ({
    platform,
    avgCTV: value,
  }));

  return (
    <Layout title="Dashboard">
      <div className="space-y-8">
        {/* Stats Cards */}
        {stats && (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
            <StatCard
              title="Total Requests"
              value={stats.totalRequests}
              icon={<Activity className="h-6 w-6" />}
              trend="up"
            />
            <StatCard
              title="Bid Rate"
              value={formatPercentage(stats.bidRate)}
              icon={<TrendingUp className="h-6 w-6" />}
              trend={stats.bidRate > 50 ? 'up' : 'down'}
            />
            <StatCard
              title="Multi Impression"
              value={stats.multiImpression}
              icon={<Users className="h-6 w-6" />}
              trend="up"
            />
            <StatCard
              title="Timeout Rate"
              value={formatPercentage(stats.timeoutRate)}
              icon={<AlertCircle className="h-6 w-6" />}
              trend={stats.timeoutRate < 5 ? 'up' : 'down'}
            />
          </div>
        )}

        {/* Charts Grid */}
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
          {/* Content Summary Chart */}
          {contentChartData.length > 0 && (
            <div className="bg-white p-6 rounded-lg shadow">
              <BarChart
                data={contentChartData}
                xKey="platform"
                yKeys={['requests']}
                title="Content Requests by Platform (Last 7 Days)"
                height={300}
                colors={['#3b82f6']}
              />
            </div>
          )}

          {/* Video Summary Chart */}
          {videoChartData.length > 0 && (
            <div className="bg-white p-6 rounded-lg shadow">
              <LineChart
                data={videoChartData}
                xKey="platform"
                yKeys={['avgCTV']}
                title="Average CTV Percentage by Platform"
                height={300}
                colors={['#10b981']}
                formatAsPercentage={true}
              />
            </div>
          )}
        </div>

        {/* Recent Activity */}
        <div className="bg-white shadow rounded-lg">
          <div className="px-4 py-5 sm:p-6">
            <h3 className="text-lg leading-6 font-medium text-gray-900 mb-4">
              System Overview
            </h3>
            <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
              <div className="text-center">
                <div className="text-2xl font-bold text-primary-600">
                  {Object.keys(contentSummary).length}
                </div>
                <div className="text-sm text-gray-500">Active Platforms</div>
              </div>
              <div className="text-center">
                <div className="text-2xl font-bold text-green-600">
                  {stats ? formatNumber(stats.deals) : '0'}
                </div>
                <div className="text-sm text-gray-500">Deals Processed</div>
              </div>
              <div className="text-center">
                <div className="text-2xl font-bold text-yellow-600">
                  {stats ? formatNumber(stats.invalidRequests) : '0'}
                </div>
                <div className="text-sm text-gray-500">Invalid Requests</div>
              </div>
            </div>
          </div>
        </div>

        {/* Last Updated */}
        {data?.lastUpdated && (
          <div className="text-center text-sm text-gray-500">
            Last updated: {new Date(data.lastUpdated).toLocaleString()}
          </div>
        )}
      </div>
    </Layout>
  );
}