//src/services/storage.ts
import * as SecureStore from 'expo-secure-store';

const AUTH_TOKEN_KEY = 'auth_token';
const USER_DATA_KEY = 'user_data';

export const storeAuthToken = async (token: string): Promise<void> => {
  await SecureStore.setItemAsync(AUTH_TOKEN_KEY, token);
};

export const getAuthToken = async (): Promise<string | null> => {
  return await SecureStore.getItemAsync(AUTH_TOKEN_KEY);
};

export const removeAuthToken = async (): Promise<void> => {
  await SecureStore.deleteItemAsync(AUTH_TOKEN_KEY);
};

export const storeUserData = async (userData: string): Promise<void> => {
  await SecureStore.setItemAsync(USER_DATA_KEY, userData);
};

export const getUserData = async (): Promise<string | null> => {
  return await SecureStore.getItemAsync(USER_DATA_KEY);
};

export const removeUserData = async (): Promise<void> => {
  await SecureStore.deleteItemAsync(USER_DATA_KEY);
};