//src/contexts/AuthContext.tsx
import { createContext, useContext, useEffect, useReducer, ReactNode } from 'react';
import { AuthState, User, LoginCredentials, SignupCredentials, AuthContextType } from '../types';
import { apiService } from '../services/api';
import { 
  storeAuthToken, 
  getAuthToken, 
  removeAuthToken, 
  storeUserData, 
  getUserData, 
  removeUserData 
} from '../services/storage';

const AuthContext = createContext<AuthContextType | undefined>(undefined);

type AuthAction = 
  | { type: 'SET_LOADING'; payload: boolean }
  | { type: 'SET_USER'; payload: { user: User; token: string } }
  | { type: 'LOGOUT' };

const authReducer = (state: AuthState, action: AuthAction): AuthState => {
  switch (action.type) {
    case 'SET_LOADING':
      return { ...state, isLoading: action.payload };
    case 'SET_USER':
      return {
        ...state,
        user: action.payload.user,
        token: action.payload.token,
        isAuthenticated: true,
        isLoading: false,
      };
    case 'LOGOUT':
      return {
        user: null,
        token: null,
        isAuthenticated: false,
        isLoading: false,
      };
    default:
      return state;
  }
};

const initialState: AuthState = {
  user: null,
  token: null,
  isAuthenticated: false,
  isLoading: true,
};

export const AuthProvider = ({ children }: { children: ReactNode }) => {
  const [state, dispatch] = useReducer(authReducer, initialState);

  useEffect(() => {
    loadStoredAuth();
  }, []);

  const loadStoredAuth = async () => {
    try {
      const token = await getAuthToken();
      const userDataString = await getUserData();

      if (token && userDataString) {
        const user = JSON.parse(userDataString);
        dispatch({ type: 'SET_USER', payload: { user, token } });
      }
    } catch (error) {
      console.error('Error loading stored auth:', error);
    } finally {
      dispatch({ type: 'SET_LOADING', payload: false });
    }
  };

  const login = async (credentials: LoginCredentials): Promise<void> => {
    try {
      dispatch({ type: 'SET_LOADING', payload: true });
      
      const response = await apiService.login(credentials);
      
      // Ensure token is a string
      const tokenString = typeof response.token === 'string' ? response.token : String(response.token);
      
      await storeAuthToken(tokenString);
      await storeUserData(JSON.stringify(response.user));

      dispatch({ type: 'SET_USER', payload: { user: response.user, token: tokenString } });
    } catch (error) {
      dispatch({ type: 'SET_LOADING', payload: false });
      throw error;
    }
  };

  const signup = async (credentials: SignupCredentials): Promise<void> => {
    try {
      dispatch({ type: 'SET_LOADING', payload: true });
      
      const response = await apiService.signup(credentials);
      
      // Ensure token is a string and user data is properly serialized
      const tokenString = typeof response.token === 'string' ? response.token : String(response.token);
      const userDataString = JSON.stringify(response.user);
      
      await storeAuthToken(tokenString);
      await storeUserData(userDataString);

      dispatch({ type: 'SET_USER', payload: { user: response.user, token: tokenString } });
    } catch (error) {
      dispatch({ type: 'SET_LOADING', payload: false });
      throw error;
    }
  };

  const logout = async (): Promise<void> => {
    try {
      await removeAuthToken();
      await removeUserData();
      dispatch({ type: 'LOGOUT' });
    } catch (error) {
      console.error('Error during logout:', error);
      // Still dispatch logout even if cleanup fails
      dispatch({ type: 'LOGOUT' });
    }
  };

  const requestPasswordReset = async (email: string): Promise<void> => {
    await apiService.requestPasswordReset(email);
  };

  return (
    <AuthContext.Provider
      value={{
        ...state,
        login,
        signup,
        logout,
        requestPasswordReset,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
};

export const useAuth = (): AuthContextType => {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
};