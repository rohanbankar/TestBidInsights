import React from 'react';
import { Input } from './Input';

interface DatePickerProps {
  label?: string;
  value?: string;
  onChange: (value: string) => void;
  min?: string;
  max?: string;
  error?: string;
  className?: string;
}

export function DatePicker({
  label,
  value,
  onChange,
  min,
  max,
  error,
  className,
}: DatePickerProps) {
  return (
    <Input
      type="date"
      label={label}
      value={value}
      onChange={(e) => onChange(e.target.value)}
      min={min}
      max={max}
      error={error}
      className={className}
    />
  );
}