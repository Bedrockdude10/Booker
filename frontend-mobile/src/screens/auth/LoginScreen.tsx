// src/screens/auth/LoginScreen.tsx
import React, { useState } from 'react';
import {
  View,
  Text,
  TextInput,
  TouchableOpacity,
  KeyboardAvoidingView,
  Platform,
  ScrollView,
  Alert,
} from 'react-native';
import { Ionicons } from '@expo/vector-icons';
import { useNavigation } from '@react-navigation/native';
import { StackNavigationProp } from '@react-navigation/stack';
import { useAuth } from '../../contexts/AuthContext';
import { LoadingSpinner } from '../../components/common/LoadingSpinner';
import { globalStyles, theme } from '../../styles';
import { AuthStackParamList } from '../../navigation/AuthNavigator';

type NavigationProp = StackNavigationProp<AuthStackParamList, 'Login'>;

export const LoginScreen: React.FC = () => {
  const navigation = useNavigation<NavigationProp>();
  const { login, isLoading } = useAuth();

  const [formData, setFormData] = useState({
    email: '',
    password: '',
  });
  const [showPassword, setShowPassword] = useState(false);
  const [errors, setErrors] = useState<{[key: string]: string}>({});

  // ... validation and handler functions remain the same ...

  const handleInputChange = (field: keyof typeof formData, value: string) => {
    setFormData(prev => ({ ...prev, [field]: value }));
    // Clear error for this field when user starts typing
    if (errors[field]) {
      setErrors(prev => ({ ...prev, [field]: '' }));
    }
  };

  const handleLogin = async () => {
    if (!validateForm()) return;
  
    try {
      await login({
        email: formData.email.toLowerCase().trim(),
        password: formData.password,
      });
    } catch (error) {
      const errorMessage = error instanceof Error ? error.message : 'Login failed';
      Alert.alert('Login Failed', errorMessage);
    }
  };

  const validateForm = () => {
    const newErrors: {[key: string]: string} = {};
  
    if (!formData.email.trim()) {
      newErrors.email = 'Email is required';
    } else if (!/\S+@\S+\.\S+/.test(formData.email)) {
      newErrors.email = 'Please enter a valid email';
    }
  
    if (!formData.password) {
      newErrors.password = 'Password is required';
    } else if (formData.password.length < 6) {
      newErrors.password = 'Password must be at least 6 characters';
    }
  
    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  return (
    <KeyboardAvoidingView
      style={globalStyles.container}
      behavior={Platform.OS === 'ios' ? 'padding' : 'height'}
    >
      <ScrollView 
        contentContainerStyle={[globalStyles.scrollContainer, globalStyles.justifyCenter]}
        keyboardShouldPersistTaps="handled"
      >
        <View style={[globalStyles.alignCenter, globalStyles.mb_lg]}>
          <Text style={[globalStyles.h1, globalStyles.mb_sm]}>Welcome to Booker</Text>
          <Text style={[globalStyles.body, { textAlign: 'center' }]}>Connect with amazing artists</Text>
        </View>

        <View>
          {/* Email Input */}
          <View style={globalStyles.inputContainer}>
            <Text style={globalStyles.inputLabel}>Email</Text>
            <View style={[globalStyles.inputWrapper, errors.email && globalStyles.inputWrapperError]}>
              <Ionicons name="mail-outline" size={20} color={theme.colors.textSecondary} style={globalStyles.inputIcon} />
              <TextInput
                style={globalStyles.textInput}
                value={formData.email}
                onChangeText={(value) => handleInputChange('email', value)}
                placeholder="your@email.com"
                placeholderTextColor={theme.colors.textSecondary}
                keyboardType="email-address"
                autoCapitalize="none"
                autoComplete="email"
                textContentType="emailAddress"
              />
            </View>
            {errors.email ? <Text style={globalStyles.errorText}>{errors.email}</Text> : null}
          </View>

          {/* Password Input */}
          <View style={globalStyles.inputContainer}>
            <Text style={globalStyles.inputLabel}>Password</Text>
            <View style={[globalStyles.inputWrapper, errors.password && globalStyles.inputWrapperError]}>
              <Ionicons name="lock-closed-outline" size={20} color={theme.colors.textSecondary} style={globalStyles.inputIcon} />
              <TextInput
                style={globalStyles.textInput}
                value={formData.password}
                onChangeText={(value) => handleInputChange('password', value)}
                placeholder="Enter your password"
                placeholderTextColor={theme.colors.textSecondary}
                secureTextEntry={!showPassword}
                autoComplete="password"
                textContentType="password"
              />
              <TouchableOpacity
                onPress={() => setShowPassword(!showPassword)}
                style={{ padding: theme.spacing.sm }}
              >
                <Ionicons 
                  name={showPassword ? "eye-outline" : "eye-off-outline"} 
                  size={20} 
                  color={theme.colors.textSecondary} 
                />
              </TouchableOpacity>
            </View>
            {errors.password ? <Text style={globalStyles.errorText}>{errors.password}</Text> : null}
          </View>

          {/* Forgot Password Link */}
          <TouchableOpacity
            style={[{ alignItems: 'flex-end' }, globalStyles.mb_lg]}
            onPress={() => navigation.navigate('ForgotPassword')}
          >
            <Text style={[globalStyles.bodySmall, { color: theme.colors.primary, fontWeight: theme.fontWeight.medium }]}>
              Forgot your password?
            </Text>
          </TouchableOpacity>

          {/* Login Button */}
          <TouchableOpacity
            style={[globalStyles.button, globalStyles.buttonPrimary, isLoading && globalStyles.buttonDisabled, globalStyles.mb_lg]}
            onPress={handleLogin}
            disabled={isLoading}
          >
            {isLoading ? (
              <LoadingSpinner size="small" color={theme.colors.textInverse} />
            ) : (
              <Text style={[globalStyles.buttonText, globalStyles.buttonTextPrimary]}>Sign In</Text>
            )}
          </TouchableOpacity>

          {/* Sign Up Link */}
          <View style={[globalStyles.row, globalStyles.justifyCenter, globalStyles.alignCenter]}>
            <Text style={globalStyles.body}>Don't have an account?</Text>
            <TouchableOpacity onPress={() => navigation.navigate('Signup')}>
              <Text style={[globalStyles.body, { color: theme.colors.primary, fontWeight: theme.fontWeight.semibold, marginLeft: theme.spacing.xs }]}>
                Sign Up
              </Text>
            </TouchableOpacity>
          </View>
        </View>
      </ScrollView>
    </KeyboardAvoidingView>
  );
};