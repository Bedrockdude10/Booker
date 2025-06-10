// src/constants/index.ts - App constants (non-style related)
export const API_BASE_URL = __DEV__ 
  ? 'http://localhost:8080/api' 
  : 'https://your-production-api.com/api';

export const GENRES = [
  'Rock', 'Pop', 'Hip Hop', 'Electronic', 'Jazz', 'Blues', 
  'Country', 'Folk', 'R&B', 'Reggae', 'Classical', 'Alternative'
];