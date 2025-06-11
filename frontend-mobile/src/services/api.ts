import { API_BASE_URL } from '../utils/constants';
import { getAuthToken } from './storage';
import { 
  Artist, 
  LoginCredentials, 
  SignupCredentials, 
  User, 
  RecommendationParams,
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
    return this.request('/accounts/password/reset', {
      method: 'POST',
      body: JSON.stringify({ email }),
    });
  }

  // User endpoints
  async getProfile(userId: string): Promise<User> {
    return this.request(`/accounts/${userId}`);
  }

  async updateProfile(userId: string, data: Partial<User>): Promise<User> {
    return this.request(`/accounts/${userId}`, {
      method: 'PUT',
      body: JSON.stringify(data),
    });
  }

  // Artist recommendation endpoints
  async getPersonalizedRecommendations(params: RecommendationParams): Promise<Artist[]> {
    const queryParams = new URLSearchParams({
      limit: params.limit?.toString() || '20',
      offset: params.offset?.toString() || '0',
    });

    if (params.filters?.genre?.length) {
      queryParams.append('genres', params.filters.genre.join(','));
    }
    if (params.filters?.location) {
      queryParams.append('locations', params.filters.location);
    }

    const response: any = await this.request(`/recommendations/user/${params.userId}?${queryParams}`);
    
    // Extract and transform the artist data
    return response.data?.map((item: any) => ({
      id: item.artist._id,
      name: item.artist.name,
      genre: item.artist.genres,
      location: item.artist.cities[0] || '',
      bio: item.artist.manager ? `Managed by ${item.artist.manager}` : 'No bio available',
      imageUrl: undefined,
      rating: undefined,
      bookingRate: undefined,
      contactInfo: {
        email: '',
        phone: '',
        website: '',
        social: {}
      },
      createdAt: new Date().toISOString(),
    })) || [];
}

  async getArtistsByGenre(genre: string): Promise<Artist[]> {
    return this.request(`/recommendations/genre/${encodeURIComponent(genre)}`);
  }

  async getArtistsByCity(city: string): Promise<Artist[]> {
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
}

export const apiService = new ApiService();