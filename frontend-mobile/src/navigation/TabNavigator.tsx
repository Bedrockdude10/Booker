// src/navigation/TabNavigator.tsx
import React from 'react';
import { createBottomTabNavigator, BottomTabScreenProps } from '@react-navigation/bottom-tabs';
import { createStackNavigator, StackScreenProps } from '@react-navigation/stack';
import { Ionicons } from '@expo/vector-icons';
import { RouteProp } from '@react-navigation/native';
import { ArtistListScreen } from '../screens/ArtistListScreen';
import { ArtistDetailScreen } from '../screens/ArtistDetailScreen';
import { ProfileScreen } from '../screens/ProfileScreen';
import { theme } from '../styles';
import { MainStackParamList, ArtistTabParamList } from './types';

const Tab = createBottomTabNavigator<ArtistTabParamList>();
const Stack = createStackNavigator<MainStackParamList>();

// Type for tab bar icon props
interface TabBarIconProps {
  focused: boolean;
  color: string;
  size: number;
}

const ArtistTabNavigator = () => {
  return (
    <Tab.Navigator
      screenOptions={({ route }: { route: RouteProp<ArtistTabParamList> }) => ({
        tabBarIcon: ({ focused, color, size }: TabBarIconProps) => {
          let iconName: keyof typeof Ionicons.glyphMap;

          if (route.name === 'Discover') {
            iconName = focused ? 'musical-notes' : 'musical-notes-outline';
          } else if (route.name === 'Profile') {
            iconName = focused ? 'person' : 'person-outline';
          } else {
            iconName = 'help-outline';
          }

          return <Ionicons name={iconName} size={size} color={color} />;
        },
        tabBarActiveTintColor: theme.colors.primary,
        tabBarInactiveTintColor: theme.colors.textSecondary,
        tabBarStyle: {
          backgroundColor: theme.colors.surface,
          borderTopColor: theme.colors.border,
        },
        headerStyle: {
          backgroundColor: theme.colors.surface,
        },
        headerTintColor: theme.colors.text,
        headerTitleStyle: {
          fontWeight: theme.fontWeight.bold,
        },
      })}
    >
      <Tab.Screen 
        name="Discover" 
        component={ArtistListScreen}
        options={{ 
          title: 'Discover Artists',
          headerTitle: 'Discover'
        }}
      />
      <Tab.Screen 
        name="Profile" 
        component={ProfileScreen}
        options={{ 
          title: 'My Profile',
          headerTitle: 'Profile'
        }}
      />
    </Tab.Navigator>
  );
};

export const TabNavigator = () => {
  return (
    <Stack.Navigator
      screenOptions={{
        headerStyle: {
          backgroundColor: theme.colors.surface,
        },
        headerTintColor: theme.colors.text,
        headerTitleStyle: {
          fontWeight: theme.fontWeight.bold,
        },
      }}
    >
      <Stack.Screen 
        name="ArtistTabs" 
        component={ArtistTabNavigator}
        options={{ headerShown: false }}
      />
      <Stack.Screen 
        name="ArtistDetail" 
        component={ArtistDetailScreen}
        options={({ route }: StackScreenProps<MainStackParamList, 'ArtistDetail'>) => ({
          title: route.params?.artist?.name || 'Artist Details',
          headerBackTitle: 'Back',
        })}
      />
    </Stack.Navigator>
  );
};