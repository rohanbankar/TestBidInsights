export interface PlatformStats {
  date: string;
  totalRequests: number;
  multiImpression: number;
  bigGuidance: number;
  addressable: number;
  complianceStrings: number;
  deals: number;
  tmax: number;
  invalidRequests: number;
  timeoutRate: number;
  bidRate: number;
  createdAt: string;
}

export interface ContentHealth {
  date: string;
  platform: string;
  totalRequests: number;
  album: number;
  artist: number;
  cat: number;
  context: number;
  data: number;
  embeddable: number;
  episode: number;
  genre: number;
  id: number;
  kwarray: number;
  keywords: number;
  length: number;
  language: number;
  livestream: number;
  season: number;
  series: number;
  title: number;
  url: number;
  videoquality: number;
  createdAt: string;
}

export interface VideoHealth {
  date: string;
  platform: string;
  percentCtv: number;
  api: number;
  boxingAllowed: number;
  delivery: number;
  h: number;
  linearity: number;
  maxBitrate: number;
  maxDuration: number;
  mimes: number;
  minBitrate: number;
  minCpmPerSec: number;
  minDuration: number;
  placement: number;
  playBackend: number;
  podDur: number;
  podId: number;
  pos: number;
  protocols: number;
  rqdDurs: number;
  skip: number;
  skipAfter: number;
  skipMin: number;
  slotInPod: number;
  startDelay: number;
  w: number;
  maxSeq: number;
  companionAd: number;
  companionType: number;
  protocol: number;
  placementType: number;
  createdAt: string;
}

export interface DashboardSummary {
  latestStats: PlatformStats;
  contentSummary: Record<string, number>;
  videoSummary: Record<string, number>;
  lastUpdated: string;
}

export interface ReportsResponse<T> {
  data: T[];
  count: number;
  query: {
    startDate?: string;
    endDate?: string;
    platform?: string;
  };
}