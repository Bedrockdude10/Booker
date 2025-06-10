// src/components/common/FilterModal.tsx
import React, { useState, useEffect } from 'react';
import {
  View,
  Text,
  Modal,
  TouchableOpacity,
  StyleSheet,
  ScrollView,
  TextInput,
  SafeAreaView,
} from 'react-native';
import { Ionicons } from '@expo/vector-icons';
import { FilterOptions } from '../../types';
import { COLORS, SPACING, GENRES } from '../../utils/constants';

interface FilterModalProps {
  visible: boolean;
  filters: FilterOptions;
  onApply: (filters: FilterOptions) => void;
  onClose: () => void;
  onClear: () => void;
}

export const FilterModal: React.FC<FilterModalProps> = ({
  visible,
  filters,
  onApply,
  onClose,
  onClear,
}) => {
  const [selectedGenres, setSelectedGenres] = useState<string[]>(filters.genre || []);
  const [location, setLocation] = useState(filters.location || '');
  const [minRating, setMinRating] = useState(filters.minRating?.toString() || '');

  useEffect(() => {
    // Reset form when filters change externally
    setSelectedGenres(filters.genre || []);
    setLocation(filters.location || '');
    setMinRating(filters.minRating?.toString() || '');
  }, [filters]);

  const toggleGenre = (genre: string) => {
    setSelectedGenres(prev => 
      prev.includes(genre) 
        ? prev.filter(g => g !== genre)
        : [...prev, genre]
    );
  };

  const handleApply = () => {
    const newFilters: FilterOptions = {};
    
    if (selectedGenres.length > 0) {
      newFilters.genre = selectedGenres;
    }
    
    if (location.trim()) {
      newFilters.location = location.trim();
    }
    
    if (minRating && !isNaN(parseFloat(minRating))) {
      newFilters.minRating = parseFloat(minRating);
    }

    onApply(newFilters);
  };

  const handleClear = () => {
    setSelectedGenres([]);
    setLocation('');
    setMinRating('');
    onClear();
  };

  const hasActiveFilters = selectedGenres.length > 0 || location.trim() || minRating;

  return (
    <Modal
      visible={visible}
      animationType="slide"
      presentationStyle="pageSheet"
      onRequestClose={onClose}
    >
      <SafeAreaView style={styles.container}>
        <View style={styles.header}>
          <TouchableOpacity onPress={onClose} style={styles.closeButton}>
            <Ionicons name="close" size={24} color={COLORS.text} />
          </TouchableOpacity>
          <Text style={styles.title}>Filter Artists</Text>
          <TouchableOpacity 
            onPress={handleClear} 
            style={styles.clearButton}
            disabled={!hasActiveFilters}
          >
            <Text style={[
              styles.clearText, 
              !hasActiveFilters && styles.clearTextDisabled
            ]}>
              Clear
            </Text>
          </TouchableOpacity>
        </View>

        <ScrollView style={styles.content} showsVerticalScrollIndicator={false}>
          {/* Genres Section */}
          <View style={styles.section}>
            <Text style={styles.sectionTitle}>Genres</Text>
            <View style={styles.genreGrid}>
              {GENRES.map((genre) => (
                <TouchableOpacity
                  key={genre}
                  style={[
                    styles.genreChip,
                    selectedGenres.includes(genre) && styles.genreChipSelected,
                  ]}
                  onPress={() => toggleGenre(genre)}
                >
                  <Text
                    style={[
                      styles.genreChipText,
                      selectedGenres.includes(genre) && styles.genreChipTextSelected,
                    ]}
                  >
                    {genre}
                  </Text>
                </TouchableOpacity>
              ))}
            </View>
          </View>

          {/* Location Section */}
          <View style={styles.section}>
            <Text style={styles.sectionTitle}>Location</Text>
            <TextInput
              style={styles.textInput}
              value={location}
              onChangeText={setLocation}
              placeholder="Enter city or region"
              placeholderTextColor={COLORS.textSecondary}
            />
          </View>

          {/* Rating Section */}
          <View style={styles.section}>
            <Text style={styles.sectionTitle}>Minimum Rating</Text>
            <TextInput
              style={styles.textInput}
              value={minRating}
              onChangeText={setMinRating}
              placeholder="e.g., 4.0"
              placeholderTextColor={COLORS.textSecondary}
              keyboardType="decimal-pad"
            />
          </View>
        </ScrollView>

        <View style={styles.footer}>
          <TouchableOpacity
            style={styles.applyButton}
            onPress={handleApply}
          >
            <Text style={styles.applyButtonText}>Apply Filters</Text>
          </TouchableOpacity>
        </View>
      </SafeAreaView>
    </Modal>
  );
};

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: COLORS.background,
  },
  header: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    paddingHorizontal: SPACING.md,
    paddingVertical: SPACING.md,
    borderBottomWidth: 1,
    borderBottomColor: COLORS.border,
    backgroundColor: COLORS.surface,
  },
  closeButton: {
    padding: SPACING.sm,
  },
  title: {
    fontSize: 18,
    fontWeight: 'bold',
    color: COLORS.text,
  },
  clearButton: {
    padding: SPACING.sm,
  },
  clearText: {
    fontSize: 16,
    color: COLORS.primary,
    fontWeight: '600',
  },
  clearTextDisabled: {
    color: COLORS.textSecondary,
  },
  content: {
    flex: 1,
    paddingHorizontal: SPACING.md,
  },
  section: {
    marginBottom: SPACING.xl,
  },
  sectionTitle: {
    fontSize: 18,
    fontWeight: '600',
    color: COLORS.text,
    marginBottom: SPACING.md,
    marginTop: SPACING.md,
  },
  genreGrid: {
    flexDirection: 'row',
    flexWrap: 'wrap',
    gap: SPACING.sm,
  },
  genreChip: {
    paddingHorizontal: SPACING.md,
    paddingVertical: SPACING.sm,
    borderRadius: 20,
    borderWidth: 1,
    borderColor: COLORS.border,
    backgroundColor: COLORS.surface,
  },
  genreChipSelected: {
    backgroundColor: COLORS.primary,
    borderColor: COLORS.primary,
  },
  genreChipText: {
    fontSize: 14,
    color: COLORS.text,
    fontWeight: '500',
  },
  genreChipTextSelected: {
    color: COLORS.surface,
  },
  textInput: {
    borderWidth: 1,
    borderColor: COLORS.border,
    borderRadius: 8,
    paddingHorizontal: SPACING.md,
    paddingVertical: SPACING.md,
    fontSize: 16,
    color: COLORS.text,
    backgroundColor: COLORS.surface,
  },
  footer: {
    padding: SPACING.md,
    backgroundColor: COLORS.surface,
    borderTopWidth: 1,
    borderTopColor: COLORS.border,
  },
  applyButton: {
    backgroundColor: COLORS.primary,
    paddingVertical: SPACING.md,
    borderRadius: 12,
    alignItems: 'center',
  },
  applyButtonText: {
    color: COLORS.surface,
    fontSize: 18,
    fontWeight: '600',
  },
});