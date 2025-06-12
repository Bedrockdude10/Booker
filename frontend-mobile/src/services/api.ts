// src/services/api.ts - Updated to match backend API contract
import { API_BASE_URL } from '../utils/constants';
import { getAuthToken } from './storage';
import { 
  Artist, 
  LoginCredentials, 
  SignupCredentials, 
  User, 
  RecommendationParams,
  RecommendationResponse,
} from '../types';

class ApiService {
  private async request<T>(
    endpoint: string, 
    options: RequestInit = {}
  ): Promise<T> {
    const token = await getAuthToken();
    
    const config: RequestInit = {
      headers: {
        'Content-Type': 'application/json',
        ...(token && { Authorization: `Bearer ${token}` }),
        ...options.headers,
      },
      ...options,
    };

    const response = await fetch(`${API_BASE_URL}${endpoint}`, config);
    
    if (!response.ok) {
      const errorData = await response.json().catch(() => ({ message: 'Network error' }));
      throw new Error(errorData.message || `HTTP ${response.status}`);
    }

    return response.json();
  }

  // Auth endpoints
  async login(credentials: LoginCredentials): Promise<{ user: User; token: string }> {
    return this.request('/auth/login', {
      method: 'POST',
      body: JSON.stringify(credentials),
    });
  }

  async signup(credentials: SignupCredentials): Promise<{ user: User; token: string }> {
    return this.request('/auth/register', {
      method: 'POST',
      body: JSON.stringify(credentials),
    });
  }

  async requestPasswordReset(email: string): Promise<{ message: string }> {
    return this.request('/auth/password/reset', {
      method: 'POST',
      body: JSON.stringify({ email }),
    });
  }

  // User endpoints
  async getProfile(userId: string): Promise<User> {
    return this.request(`/accounts`);
  }

  async updateProfile(userId: string, data: Partial<User>): Promise<User> {
    return this.request(`/accounts`, {
      method: 'PUT',
      body: JSON.stringify(data),
    });
  }

  // Artist recommendation endpoints - Updated to match backend
  async getPersonalizedRecommendations(params: RecommendationParams): Promise<RecommendationResponse> {
    const queryParams = new URLSearchParams({
      limit: params.limit?.toString() || '20',
      offset: params.offset?.toString() || '0',
    });

    // Use backend filter property names
    if (params.filters?.genres?.length) {
      queryParams.append('genres', params.filters.genres.join(','));
    }
    if (params.filters?.cities?.length) {
      queryParams.append('cities', params.filters.cities.join(','));
    }
    if (params.filters?.minRating) {
      queryParams.append('minRating', params.filters.minRating.toString());
    }
    if (params.filters?.maxRating) {
      queryParams.append('maxRating', params.filters.maxRating.toString());
    }
    if (params.filters?.hasManager !== undefined) {
      queryParams.append('hasManager', params.filters.hasManager.toString());
    }
    if (params.filters?.hasSpotify !== undefined) {
      queryParams.append('hasSpotify', params.filters.hasSpotify.toString());
    }

    return this.request(`/recommendations/user/${params.userId}?${queryParams}`);
  }

  async getArtistsByGenre(genre: string): Promise<RecommendationResponse> {
    return this.request(`/recommendations/genre/${encodeURIComponent(genre)}`);
  }

  async getArtistsByCity(city: string): Promise<RecommendationResponse> {
    return this.request(`/recommendations/city/${encodeURIComponent(city)}`);
  }

  async getArtistDetail(artistId: string): Promise<Artist> {
    return this.request(`/artists/${artistId}`);
  }

  // Interaction tracking
  async trackInteraction(data: {
    userId: string;
    artistId: string;
    type: 'view' | 'contact' | 'save';
  }): Promise<void> {
    return this.request('/recommendations/interactions', {
      method: 'POST',
      body: JSON.stringify(data),
    });
  }

  // Track streaming service clicks - optional analytics
  async trackStreamingServiceClick(data: {
    userId?: string;
    artistId: string;
    service: 'bandcamp' | 'spotify' | 'appleMusic' | 'instagram';
    url?: string;
  }): Promise<void> {
    // Only track if we have the interactions endpoint available
    try {
      return this.request('/recommendations/interactions', {
        method: 'POST',
        body: JSON.stringify({
          userId: data.userId,
          artistId: data.artistId,
          type: 'streaming_service_click',
          metadata: {
            service: data.service,
            url: data.url,
          },
        }),
      });
    } catch (error) {
      // If tracking fails, don't break the user experience
      console.warn('Failed to track streaming service click:', error);
    }
  }
}

export const apiService = new ApiService();