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
import { COLORS, SPACING } from '../utils/constants';
import { styles } from './ProfileScreen.styles';

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
    <View style={styles.fieldContainer}>
      <Text style={styles.fieldLabel}>{label}</Text>
      {editing ? (
        <TextInput
          style={[styles.textInput, multiline && styles.textInputMultiline]}
          value={value}
          onChangeText={onChangeText}
          placeholder={placeholder}
          placeholderTextColor={COLORS.textSecondary}
          multiline={multiline}
          numberOfLines={multiline ? 4 : 1}
          keyboardType={keyboardType}
        />
      ) : (
        <Text style={[styles.fieldValue, !value && styles.fieldValueEmpty]}>
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
      style={styles.container}
      behavior={Platform.OS === 'ios' ? 'padding' : 'height'}
    >
      <ScrollView showsVerticalScrollIndicator={false}>
        <View style={styles.header}>
          <View style={styles.avatarContainer}>
            <Ionicons name="person" size={60} color={COLORS.textSecondary} />
          </View>
          <Text style={styles.userName}>{user.name}</Text>
          <Text style={styles.userRole}>Promoter</Text>
        </View>

        <View style={styles.formContainer}>
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

        <View style={styles.buttonContainer}>
          {editing ? (
            <View style={styles.editingButtons}>
              <TouchableOpacity
                style={[styles.button, styles.cancelButton]}
                onPress={handleCancel}
                disabled={loading}
              >
                <Text style={styles.cancelButtonText}>Cancel</Text>
              </TouchableOpacity>
              
              <TouchableOpacity
                style={[styles.button, styles.saveButton, loading && styles.buttonDisabled]}
                onPress={handleSave}
                disabled={loading}
              >
                {loading ? (
                  <LoadingSpinner size="small" color={COLORS.surface} />
                ) : (
                  <>
                    <Ionicons name="checkmark" size={20} color={COLORS.surface} />
                    <Text style={styles.saveButtonText}>Save</Text>
                  </>
                )}
              </TouchableOpacity>
            </View>
          ) : (
            <TouchableOpacity
              style={[styles.button, styles.editButton]}
              onPress={() => setEditing(true)}
            >
              <Ionicons name="pencil" size={20} color={COLORS.primary} />
              <Text style={styles.editButtonText}>Edit Profile</Text>
            </TouchableOpacity>
          )}

          <TouchableOpacity
            style={[styles.button, styles.logoutButton]}
            onPress={handleLogout}
          >
            <Ionicons name="log-out-outline" size={20} color={COLORS.error} />
            <Text style={styles.logoutButtonText}>Sign Out</Text>
          </TouchableOpacity>
        </View>

        <View style={styles.footer}>
          <Text style={styles.footerText}>
            Member since {new Date(user.createdAt).toLocaleDateString()}
          </Text>
        </View>
      </ScrollView>
    </KeyboardAvoidingView>
  );
};