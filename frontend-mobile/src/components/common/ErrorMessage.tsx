// src/components/common/ErrorMessage.tsx
import React from 'react';
import { View, Text, TouchableOpacity, StyleSheet } from 'react-native';
import { Ionicons } from '@expo/vector-icons';
import { COLORS, SPACING } from '../../utils/constants';

interface ErrorMessageProps {
  message: string;
  onRetry?: () => void;
  style?: object;
}

export const ErrorMessage: React.FC<ErrorMessageProps> = ({ 
  message, 
  onRetry,
  style = {} 
}) => {
  return (
    <View style={[styles.container, style]}>
      <Ionicons name="alert-circle-outline" size={48} color={COLORS.error} />
      <Text style={styles.message}>{message}</Text>
      {onRetry && (
        <TouchableOpacity style={styles.retryButton} onPress={onRetry}>
          <Ionicons name="refresh" size={20} color={COLORS.surface} />
          <Text style={styles.retryText}>Try Again</Text>
        </TouchableOpacity>
      )}
    </View>
  );
};

const styles = StyleSheet.create({
  container: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    padding: SPACING.xl,
    backgroundColor: COLORS.background,
  },
  message: {
    fontSize: 16,
    color: COLORS.text,
    textAlign: 'center',
    marginVertical: SPACING.lg,
    lineHeight: 22,
  },
  retryButton: {
    flexDirection: 'row',
    alignItems: 'center',
    backgroundColor: COLORS.primary,
    paddingHorizontal: SPACING.lg,
    paddingVertical: SPACING.md,
    borderRadius: 8,
    marginTop: SPACING.md,
  },
  retryText: {
    color: COLORS.surface,
    fontSize: 16,
    fontWeight: '600',
    marginLeft: SPACING.sm,
  },
});