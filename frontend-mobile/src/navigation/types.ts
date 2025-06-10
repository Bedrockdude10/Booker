// src/navigation/types.ts - Navigation type definitions
import { NavigatorScreenParams } from '@react-navigation/native';
import { Artist } from '../types';

export type RootStackParamList = {
  Auth: NavigatorScreenParams<AuthStackParamList>;
  Main: NavigatorScreenParams<MainStackParamList>;
};

export type AuthStackParamList = {
  Login: undefined;
  Signup: undefined;
  ForgotPassword: undefined;
};

export type MainStackParamList = {
  ArtistTabs: NavigatorScreenParams<ArtistTabParamList>;
  ArtistDetail: { artist: Artist };
};

export type ArtistTabParamList = {
  Discover: undefined;
  Profile: undefined;
};

// Helper hook for navigation with proper typing
import { useNavigation } from '@react-navigation/native';
import { StackNavigationProp } from '@react-navigation/stack';

export type NavigationProp = StackNavigationProp<MainStackParamList>;

export const useAppNavigation = () => {
  return useNavigation<NavigationProp>();
};