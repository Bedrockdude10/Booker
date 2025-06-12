// src/types/index.ts
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

// src/types/index.ts - Updated Artist interface
export interface Artist {
  id: string;
  name: string;
  genre: string[];
  location: string;
  bio: string;
  imageUrl?: string;
  musicSamples?: string[];
  contactInfo: {
    email: string;
    phone?: string;
    website?: string;
    social?: {
      instagram?: string;
      spotify?: string;
      youtube?: string;
      // Add new streaming services
      bandcamp?: string;
      appleMusic?: string;
    };
  };
  bookingRate?: {
    min: number;
    max: number;
    currency: string;
  };
  tags?: string[];
  rating?: number;
  createdAt: string;
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
  genre?: string[];
  location?: string;
  minRating?: number;
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

export interface BackendArtist {
  _id: string;
  name: string;
  genres: string[];
  manager?: string;
  cities: string[];
  spotifyId?: string;
}

export interface RecommendationResult {
  artist: BackendArtist;
  score: number;
}

export interface RecommendationResponse {
  data: RecommendationResult[];
  total: number;
  requestedBy: string;
  metadata?: {
    userId?: string;
    basedOn?: string;
    reason?: string;
    [key: string]: any;
  };
}