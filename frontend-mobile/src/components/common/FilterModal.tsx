// src/components/common/FilterModal.tsx - Fixed to use backend filter property names
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
  // Use backend property names
  const [selectedGenres, setSelectedGenres] = useState<string[]>(filters.genres || []);
  const [selectedCities, setSelectedCities] = useState<string[]>(filters.cities || []);
  const [minRating, setMinRating] = useState(filters.minRating?.toString() || '');
  const [maxRating, setMaxRating] = useState(filters.maxRating?.toString() || '');
  const [hasManager, setHasManager] = useState<boolean | undefined>(filters.hasManager);
  const [hasSpotify, setHasSpotify] = useState<boolean | undefined>(filters.hasSpotify);

  useEffect(() => {
    // Reset form when filters change externally
    setSelectedGenres(filters.genres || []);
    setSelectedCities(filters.cities || []);
    setMinRating(filters.minRating?.toString() || '');
    setMaxRating(filters.maxRating?.toString() || '');
    setHasManager(filters.hasManager);
    setHasSpotify(filters.hasSpotify);
  }, [filters]);

  const toggleGenre = (genre: string) => {
    setSelectedGenres(prev => 
      prev.includes(genre) 
        ? prev.filter(g => g !== genre)
        : [...prev, genre]
    );
  };



  const removeCity = (cityToRemove: string) => {
    setSelectedCities(prev => prev.filter(city => city !== cityToRemove));
  };

  const handleApply = () => {
    const newFilters: FilterOptions = {};
    
    if (selectedGenres.length > 0) {
      newFilters.genres = selectedGenres;
    }
    
    if (selectedCities.length > 0) {
      newFilters.cities = selectedCities;
    }
    
    if (minRating && !isNaN(parseFloat(minRating))) {
      newFilters.minRating = parseFloat(minRating);
    }

    if (maxRating && !isNaN(parseFloat(maxRating))) {
      newFilters.maxRating = parseFloat(maxRating);
    }

    if (hasManager !== undefined) {
      newFilters.hasManager = hasManager;
    }

    if (hasSpotify !== undefined) {
      newFilters.hasSpotify = hasSpotify;
    }

    onApply(newFilters);
  };

  const handleClear = () => {
    setSelectedGenres([]);
    setSelectedCities([]);
    setMinRating('');
    setMaxRating('');
    setHasManager(undefined);
    setHasSpotify(undefined);
    onClear();
  };

  const hasActiveFilters = selectedGenres.length > 0 || 
                          selectedCities.length > 0 || 
                          minRating || 
                          maxRating ||
                          hasManager !== undefined ||
                          hasSpotify !== undefined;

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

          {/* Cities Section */}
          <View style={styles.section}>
            <Text style={styles.sectionTitle}>Cities</Text>
            {selectedCities.length > 0 && (
              <View style={styles.selectedCitiesContainer}>
                {selectedCities.map((city, index) => (
                  <View key={index} style={styles.cityTag}>
                    <Text style={styles.cityTagText}>{city}</Text>
                    <TouchableOpacity onPress={() => removeCity(city)}>
                      <Ionicons name="close-circle" size={16} color={COLORS.surface} />
                    </TouchableOpacity>
                  </View>
                ))}
              </View>
            )}
            <TextInput
              style={styles.textInput}
              placeholder="Enter cities separated by commas"
              placeholderTextColor={COLORS.textSecondary}
              value={selectedCities.join(', ')}
              onChangeText={(text) => {
                const cities = text.split(',').map(city => city.trim()).filter(city => city.length > 0);
                setSelectedCities(cities);
              }}
            />
          </View>

          {/* Rating Section */}
          <View style={styles.section}>
            <Text style={styles.sectionTitle}>Rating Range</Text>
            <View style={styles.ratingRow}>
              <TextInput
                style={[styles.textInput, styles.ratingInput]}
                value={minRating}
                onChangeText={setMinRating}
                placeholder="Min (e.g., 3.0)"
                placeholderTextColor={COLORS.textSecondary}
                keyboardType="decimal-pad"
              />
              <Text style={styles.ratingDash}>—</Text>
              <TextInput
                style={[styles.textInput, styles.ratingInput]}
                value={maxRating}
                onChangeText={setMaxRating}
                placeholder="Max (e.g., 5.0)"
                placeholderTextColor={COLORS.textSecondary}
                keyboardType="decimal-pad"
              />
            </View>
          </View>

          {/* Boolean Filters */}
          <View style={styles.section}>
            <Text style={styles.sectionTitle}>Additional Filters</Text>
            
            <TouchableOpacity 
              style={styles.booleanFilter}
              onPress={() => setHasManager(hasManager === true ? undefined : true)}
            >
              <Text style={styles.booleanFilterText}>Has Manager</Text>
              <Ionicons 
                name={hasManager === true ? "checkmark-circle" : "ellipse-outline"} 
                size={24} 
                color={hasManager === true ? COLORS.primary : COLORS.textSecondary} 
              />
            </TouchableOpacity>

            <TouchableOpacity 
              style={styles.booleanFilter}
              onPress={() => setHasSpotify(hasSpotify === true ? undefined : true)}
            >
              <Text style={styles.booleanFilterText}>Has Spotify</Text>
              <Ionicons 
                name={hasSpotify === true ? "checkmark-circle" : "ellipse-outline"} 
                size={24} 
                color={hasSpotify === true ? COLORS.primary : COLORS.textSecondary} 
              />
            </TouchableOpacity>
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
  selectedCitiesContainer: {
    flexDirection: 'row',
    flexWrap: 'wrap',
    gap: SPACING.sm,
    marginBottom: SPACING.sm,
  },
  cityTag: {
    flexDirection: 'row',
    alignItems: 'center',
    backgroundColor: COLORS.primary,
    paddingHorizontal: SPACING.sm,
    paddingVertical: SPACING.xs,
    borderRadius: 16,
    gap: SPACING.xs,
  },
  cityTagText: {
    color: COLORS.surface,
    fontSize: 14,
    fontWeight: '500',
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
  ratingRow: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: SPACING.md,
  },
  ratingInput: {
    flex: 1,
  },
  ratingDash: {
    fontSize: 18,
    color: COLORS.textSecondary,
  },
  booleanFilter: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    paddingVertical: SPACING.md,
    borderBottomWidth: 1,
    borderBottomColor: COLORS.border,
  },
  booleanFilterText: {
    fontSize: 16,
    color: COLORS.text,
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