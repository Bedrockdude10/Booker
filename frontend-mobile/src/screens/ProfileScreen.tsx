// src/screens/profile/ProfileScreen.tsx
import React, { useState, useEffect } from 'react';
import {
  View,
  Text,
  ScrollView,
  TouchableOpacity,
  StyleSheet,
  TextInput,
  Alert,
  KeyboardAvoidingView,
  Platform,
} from 'react-native';
import { Ionicons } from '@expo/vector-icons';
import { useAuth } from '../contexts/AuthContext';
import { apiService } from '../services/api';
import { User } from '../types';
import { LoadingSpinner } from '../components/common/LoadingSpinner';
import { globalStyles, theme } from '../styles';

// Only screen-specific styles that can't be achieved with global styles
const screenStyles = StyleSheet.create({
  avatarContainer: {
    width: 100,
    height: 100,
    borderRadius: 50,
    backgroundColor: theme.colors.border,
    justifyContent: 'center',
    alignItems: 'center',
    marginBottom: theme.spacing.md,
  },
  header: {
    alignItems: 'center',
    paddingVertical: theme.spacing.xl,
    backgroundColor: theme.colors.surface,
    borderBottomWidth: 1,
    borderBottomColor: theme.colors.border,
  },
  userRole: {
    fontSize: theme.fontSize.md,
    color: theme.colors.textSecondary,
    fontWeight: theme.fontWeight.medium,
  },
  editingButtons: {
    flexDirection: 'row',
    gap: theme.spacing.md,
    marginBottom: theme.spacing.md,
  },
  footer: {
    alignItems: 'center',
    paddingVertical: theme.spacing.xl,
  },
});

export const ProfileScreen = () => {
  const { user, logout } = useAuth();
  const [editing, setEditing] = useState(false);
  const [loading, setLoading] = useState(false);
  const [formData, setFormData] = useState<Partial<User>>({
    name: user?.name || '',
    email: user?.email || '',
    location: user?.location || '',
    phone: user?.phone || '',
    bio: user?.bio || '',
  });

  const handleSave = async () => {
    if (!user) return;

    try {
      setLoading(true);
      await apiService.updateProfile(user.id, formData);
      setEditing(false);
      Alert.alert('Success', 'Profile updated successfully');
    } catch (error) {
      const errorMessage = error instanceof Error ? error.message : 'Failed to update profile';
      Alert.alert('Error', errorMessage);
    } finally {
      setLoading(false);
    }
  };

  const handleCancel = () => {
    setFormData({
      name: user?.name || '',
      email: user?.email || '',
      location: user?.location || '',
      phone: user?.phone || '',
      bio: user?.bio || '',
    });
    setEditing(false);
  };

  const handleLogout = () => {
    Alert.alert(
      'Sign Out',
      'Are you sure you want to sign out?',
      [
        { text: 'Cancel', style: 'cancel' },
        { 
          text: 'Sign Out', 
          style: 'destructive',
          onPress: logout 
        },
      ]
    );
  };

  const renderField = (
    label: string,
    value: string,
    onChangeText: (text: string) => void,
    placeholder?: string,
    multiline = false,
    keyboardType: 'default' | 'email-address' | 'phone-pad' = 'default'
  ) => (
    <View style={globalStyles.inputContainer}>
      <Text style={globalStyles.inputLabel}>{label}</Text>
      {editing ? (
        <TextInput
          style={[
            globalStyles.textInput,
            globalStyles.inputWrapper,
            multiline && { minHeight: 100, textAlignVertical: 'top' }
          ]}
          value={value}
          onChangeText={onChangeText}
          placeholder={placeholder}
          placeholderTextColor={theme.colors.textSecondary}
          multiline={multiline}
          numberOfLines={multiline ? 4 : 1}
          keyboardType={keyboardType}
        />
      ) : (
        <Text style={[
          globalStyles.body,
          globalStyles.py_sm,
          !value && { color: theme.colors.textSecondary, fontStyle: 'italic' }
        ]}>
          {value || placeholder || 'Not provided'}
        </Text>
      )}
    </View>
  );

  if (!user) {
    return <LoadingSpinner />;
  }

  return (
    <KeyboardAvoidingView
      style={globalStyles.container}
      behavior={Platform.OS === 'ios' ? 'padding' : 'height'}
    >
      <ScrollView showsVerticalScrollIndicator={false}>
        <View style={screenStyles.header}>
          <View style={screenStyles.avatarContainer}>
            <Ionicons name="person" size={60} color={theme.colors.textSecondary} />
          </View>
          <Text style={globalStyles.h2}>{user.name}</Text>
          <Text style={screenStyles.userRole}>Promoter</Text>
        </View>

        <View style={globalStyles.p_md}>
          {renderField(
            'Full Name',
            formData.name || '',
            (text) => setFormData(prev => ({ ...prev, name: text })),
            'Enter your full name'
          )}

          {renderField(
            'Email',
            formData.email || '',
            (text) => setFormData(prev => ({ ...prev, email: text })),
            'your@email.com',
            false,
            'email-address'
          )}

          {renderField(
            'Location',
            formData.location || '',
            (text) => setFormData(prev => ({ ...prev, location: text })),
            'City, State'
          )}

          {renderField(
            'Phone',
            formData.phone || '',
            (text) => setFormData(prev => ({ ...prev, phone: text })),
            '+1 (555) 123-4567',
            false,
            'phone-pad'
          )}

          {renderField(
            'Bio',
            formData.bio || '',
            (text) => setFormData(prev => ({ ...prev, bio: text })),
            'Tell us about yourself and your experience as a promoter...',
            true
          )}
        </View>

        <View style={globalStyles.p_md}>
          {editing ? (
            <View style={screenStyles.editingButtons}>
              <TouchableOpacity
                style={[
                  globalStyles.button,
                  globalStyles.buttonSecondary,
                  globalStyles.flex1
                ]}
                onPress={handleCancel}
                disabled={loading}
              >
                <Text style={globalStyles.buttonTextSecondary}>Cancel</Text>
              </TouchableOpacity>
              
              <TouchableOpacity
                style={[
                  globalStyles.button,
                  globalStyles.buttonPrimary,
                  globalStyles.flex1,
                  loading && globalStyles.buttonDisabled
                ]}
                onPress={handleSave}
                disabled={loading}
              >
                {loading ? (
                  <LoadingSpinner size="small" color={theme.colors.textInverse} />
                ) : (
                  <View style={[globalStyles.row, globalStyles.alignCenter]}>
                    <Ionicons name="checkmark" size={20} color={theme.colors.textInverse} />
                    <Text style={[globalStyles.buttonTextPrimary, { marginLeft: theme.spacing.sm }]}>
                      Save
                    </Text>
                  </View>
                )}
              </TouchableOpacity>
            </View>
          ) : (
            <TouchableOpacity
              style={[globalStyles.button, globalStyles.buttonOutline, globalStyles.mb_md]}
              onPress={() => setEditing(true)}
            >
              <View style={[globalStyles.row, globalStyles.alignCenter]}>
                <Ionicons name="pencil" size={20} color={theme.colors.primary} />
                <Text style={[globalStyles.buttonTextOutline, { marginLeft: theme.spacing.sm }]}>
                  Edit Profile
                </Text>
              </View>
            </TouchableOpacity>
          )}

          <TouchableOpacity
            style={[
              globalStyles.button,
              globalStyles.buttonSecondary,
              { borderColor: theme.colors.error, borderWidth: 2 }
            ]}
            onPress={handleLogout}
          >
            <View style={[globalStyles.row, globalStyles.alignCenter]}>
              <Ionicons name="log-out-outline" size={20} color={theme.colors.error} />
              <Text style={[globalStyles.buttonTextSecondary, { 
                color: theme.colors.error,
                marginLeft: theme.spacing.sm 
              }]}>
                Sign Out
              </Text>
            </View>
          </TouchableOpacity>
        </View>

        <View style={screenStyles.footer}>
          <Text style={globalStyles.bodySmall}>
            Member since {new Date(user.createdAt).toLocaleDateString()}
          </Text>
        </View>
      </ScrollView>
    </KeyboardAvoidingView>
  );
};