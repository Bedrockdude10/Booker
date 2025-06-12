// src/types/index.ts - Single source of truth for Artist data
export interface User {
  id: string;
  email: string;
  name: string;
  role: 'promoter' | 'artist';
  location?: string;
  phone?: string;
  bio?: string;
  createdAt: string;
  isActive: boolean;
}

// Single Artist interface that matches the API exactly
export interface Artist {
  _id: string;
  name: string;
  genres: string[];
  manager?: string;
  cities: string[];
  spotifyId?: string;
  rating?: number;
  ratingCount?: number;
  contactInfo?: {
    social?: {
      spotify?: string;
      appleMusic?: string;
      instagram?: string;
      youtube?: string;
      facebook?: string;
      twitter?: string;
      tiktok?: string;
      website?: string;
      soundcloud?: string;
      beatport?: string;
      bandcamp?: string;
      discogs?: string;
      email?: string;
    };
    manager?: string;
    managerInfo?: string;
    bookingInfo?: string;
    labelName?: string;
    labelURL?: string;
  };
}

export interface RecommendationResult {
  artist: Artist;
  score: number;
}

export interface RecommendationResponse {
  data: RecommendationResult[];
  total: number;
  requestedBy: string;
  metadata?: {
    filters?: Record<string, any>;
    type?: string;
    userId?: string;
    basedOn?: string;
    reason?: string;
    [key: string]: any;
  };
}

export interface AuthState {
  user: User | null;
  token: string | null;
  isLoading: boolean;
  isAuthenticated: boolean;
}

export interface LoginCredentials {
  email: string;
  password: string;
}

export interface SignupCredentials {
  email: string;
  password: string;
  name: string;
  role: 'promoter';
  location?: string;
  phone?: string;
}

export interface FilterOptions {
  // Map to backend filter names
  genres?: string[];
  cities?: string[];
  minRating?: number;
  maxRating?: number;
  hasManager?: boolean;
  hasSpotify?: boolean;
}

export interface RecommendationParams {
  userId: string;
  limit?: number;
  offset?: number;
  filters?: FilterOptions;
}

export interface AuthContextType extends AuthState {
  login: (credentials: LoginCredentials) => Promise<void>;
  signup: (credentials: SignupCredentials) => Promise<void>;
  logout: () => Promise<void>;
  requestPasswordReset: (email: string) => Promise<void>;
}

export type MainStackParamList = {
  ArtistTabs: undefined;
  ArtistDetail: { artist: Artist };
};

export type ArtistTabParamList = {
  Discover: undefined;
  Profile: undefined;
};

// Helper functions for UI display (computed properties)
export const getArtistId = (artist: Artist): string => artist._id;

export const getArtistLocation = (artist: Artist): string => {
  if (!artist.cities || artist.cities.length === 0) {
    return 'Location not specified';
  }
  return artist.cities.length === 1 
    ? artist.cities[0]
    : `${artist.cities[0]} +${artist.cities.length - 1}`;
};

export const getArtistBio = (artist: Artist): string => {
  const parts: string[] = [];
  
  // Try contactInfo.manager first, then fallback to manager field
  const managerName = artist.contactInfo?.manager || artist.manager;
  if (managerName) {
    parts.push(`Managed by ${managerName}`);
  }
  
  if (artist.contactInfo?.labelName) {
    parts.push(`Signed to ${artist.contactInfo.labelName}`);
  }
  
  if (parts.length > 0) {
    return parts.join('. ') + '.';
  }
  
  // Fallback bio
  return `${artist.genres.join(', ')} artist from ${getArtistLocation(artist)}`;
};

export const getArtistEmail = (artist: Artist): string | undefined => {
  if (artist.contactInfo?.social?.email) {
    return artist.contactInfo.social.email;
  }
  
  // Check if managerInfo contains an email
  if (artist.contactInfo?.managerInfo && 
      artist.contactInfo.managerInfo.includes('@')) {
    return artist.contactInfo.managerInfo;
  }
  
  return undefined;
};

export const getArtistWebsite = (artist: Artist): string | undefined => {
  return artist.contactInfo?.social?.website;
};

// Social media helper
export const getArtistSocialLinks = (artist: Artist) => {
  const social = artist.contactInfo?.social;
  if (!social) return [];

  return [
    { key: 'instagram', icon: 'logo-instagram', color: '#E1306C', url: social.instagram },
    { key: 'spotify', icon: 'musical-note', color: '#1DB954', url: social.spotify },
    { key: 'youtube', icon: 'logo-youtube', color: '#FF0000', url: social.youtube },
    { key: 'appleMusic', icon: 'logo-apple', color: '#000000', url: social.appleMusic },
    { key: 'facebook', icon: 'logo-facebook', color: '#1877F2', url: social.facebook },
    { key: 'twitter', icon: 'logo-twitter', color: '#1DA1F2', url: social.twitter },
    { key: 'tiktok', icon: 'logo-tiktok', color: '#000000', url: social.tiktok },
    { key: 'soundcloud', icon: 'cloud', color: '#FF5500', url: social.soundcloud },
    { key: 'bandcamp', icon: 'musical-notes', color: '#629AA0', url: social.bandcamp },
  ].filter(link => link.url);
};