import React from 'react';
import { clsx } from 'clsx';
import { formatNumber } from '../../utils/formatUtils';

interface HeatMapData {
  x: string;
  y: string;
  value: number;
}

interface HeatMapProps {
  data: HeatMapData[];
  title?: string;
  height?: number;
  colorScheme?: 'blue' | 'green' | 'red' | 'purple';
}

const colorSchemes = {
  blue: {
    low: 'bg-blue-100',
    medium: 'bg-blue-300',
    high: 'bg-blue-500',
    veryHigh: 'bg-blue-700',
  },
  green: {
    low: 'bg-green-100',
    medium: 'bg-green-300',
    high: 'bg-green-500',
    veryHigh: 'bg-green-700',
  },
  red: {
    low: 'bg-red-100',
    medium: 'bg-red-300',
    high: 'bg-red-500',
    veryHigh: 'bg-red-700',
  },
  purple: {
    low: 'bg-purple-100',
    medium: 'bg-purple-300',
    high: 'bg-purple-500',
    veryHigh: 'bg-purple-700',
  },
};

export function HeatMap({ 
  data, 
  title, 
  height = 400, 
  colorScheme = 'blue' 
}: HeatMapProps) {
  if (!data.length) {
    return (
      <div className="w-full p-8 text-center text-gray-500">
        No data available
      </div>
    );
  }

  // Get unique x and y values
  const xValues = [...new Set(data.map(d => d.x))];
  const yValues = [...new Set(data.map(d => d.y))];
  
  // Find min and max values for color scaling
  const values = data.map(d => d.value);
  const minValue = Math.min(...values);
  const maxValue = Math.max(...values);
  const range = maxValue - minValue;

  const getColorIntensity = (value: number): string => {
    if (range === 0) return colorSchemes[colorScheme].medium;
    
    const normalized = (value - minValue) / range;
    
    if (normalized <= 0.25) return colorSchemes[colorScheme].low;
    if (normalized <= 0.5) return colorSchemes[colorScheme].medium;
    if (normalized <= 0.75) return colorSchemes[colorScheme].high;
    return colorSchemes[colorScheme].veryHigh;
  };

  const getCellValue = (x: string, y: string): number => {
    const cell = data.find(d => d.x === x && d.y === y);
    return cell ? cell.value : 0;
  };

  return (
    <div className="w-full">
      {title && (
        <h3 className="text-lg font-medium text-gray-900 mb-4">{title}</h3>
      )}
      
      <div className="overflow-x-auto">
        <div className="inline-block min-w-full">
          <div className="grid gap-1 p-4" style={{ height }}>
            {/* Header row */}
            <div className="grid gap-1" style={{ gridTemplateColumns: `80px repeat(${xValues.length}, 1fr)` }}>
              <div></div>
              {xValues.map(x => (
                <div key={x} className="text-xs text-gray-600 text-center font-medium p-1">
                  {x}
                </div>
              ))}
            </div>
            
            {/* Data rows */}
            {yValues.map(y => (
              <div 
                key={y} 
                className="grid gap-1" 
                style={{ gridTemplateColumns: `80px repeat(${xValues.length}, 1fr)` }}
              >
                <div className="text-xs text-gray-600 font-medium flex items-center px-1">
                  {y}
                </div>
                {xValues.map(x => {
                  const value = getCellValue(x, y);
                  return (
                    <div
                      key={`${x}-${y}`}
                      className={clsx(
                        'rounded text-xs font-medium flex items-center justify-center h-8 cursor-default',
                        getColorIntensity(value),
                        value > (minValue + range * 0.5) ? 'text-white' : 'text-gray-900'
                      )}
                      title={`${x} × ${y}: ${formatNumber(value)}`}
                    >
                      {formatNumber(value)}
                    </div>
                  );
                })}
              </div>
            ))}
          </div>
        </div>
      </div>

      {/* Legend */}
      <div className="flex items-center justify-center mt-4 space-x-4 text-xs text-gray-600">
        <span>Low</span>
        <div className="flex space-x-1">
          <div className={`w-4 h-4 rounded ${colorSchemes[colorScheme].low}`}></div>
          <div className={`w-4 h-4 rounded ${colorSchemes[colorScheme].medium}`}></div>
          <div className={`w-4 h-4 rounded ${colorSchemes[colorScheme].high}`}></div>
          <div className={`w-4 h-4 rounded ${colorSchemes[colorScheme].veryHigh}`}></div>
        </div>
        <span>High</span>
      </div>
    </div>
  );
}