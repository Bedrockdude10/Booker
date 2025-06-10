// src/navigation/AppNavigator.tsx
import React from 'react';
import { NavigationContainer } from '@react-navigation/native';
import { createStackNavigator } from '@react-navigation/stack';
import { View, ActivityIndicator, StyleSheet } from 'react-native';
import { useAuth } from '../contexts/AuthContext';
import { AuthNavigator } from './AuthNavigator';
import { TabNavigator } from './TabNavigator';
import { RootStackParamList } from './types';
import { theme } from '../styles';

const RootStack = createStackNavigator<RootStackParamList>();

const LoadingScreen = () => (
  <View style={styles.loadingContainer}>
    <ActivityIndicator size="large" color={theme.colors.primary} />
  </View>
);

export const AppNavigator = () => {
    const { isAuthenticated, isLoading } = useAuth();
  
    if (isLoading) {
      return <LoadingScreen />;
    }
  
    return (
      <NavigationContainer>
        <RootStack.Navigator
          id={"RootStack" as any} 
          screenOptions={{ headerShown: false }}
          initialRouteName={isAuthenticated ? "Main" : "Auth"}
        >
          <RootStack.Screen name="Auth" component={AuthNavigator} />
          <RootStack.Screen name="Main" component={TabNavigator} />
        </RootStack.Navigator>
      </NavigationContainer>
    );
  };

const styles = StyleSheet.create({
  loadingContainer: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    backgroundColor: theme.colors.background,
  },
});