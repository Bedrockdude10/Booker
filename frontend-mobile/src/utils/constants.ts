// src/utils/constants.ts - Updated to use centralized styles
import { theme } from '../styles/theme';

// Re-export for backward compatibility during migration
export const COLORS = theme.colors;
export const SPACING = theme.spacing;

// Re-export other constants
export { API_BASE_URL, GENRES } from '../constants';