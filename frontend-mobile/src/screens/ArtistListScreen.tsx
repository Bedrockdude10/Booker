// src/screens/artists/ArtistListScreen.tsx
import React, { useState, useEffect, useCallback } from 'react';
import {
  View,
  Text,
  FlatList,
  TouchableOpacity,
  StyleSheet,
  RefreshControl,
} from 'react-native';
import { Ionicons } from '@expo/vector-icons';
import { useNavigation } from '@react-navigation/native';
import { StackNavigationProp } from '@react-navigation/stack';
import { useAuth } from '../contexts/AuthContext';
import { apiService } from '../services/api';
import { Artist, FilterOptions } from '../types';
import { ArtistCard } from '../components/artists/ArtistCard';
import { FilterModal } from '../components/common/FilterModal';
import { LoadingSpinner } from '../components/common/LoadingSpinner';
import { ErrorMessage } from '../components/common/ErrorMessage';
import { COLORS, SPACING } from '../utils/constants';
import { MainStackParamList } from '../navigation/types';

type NavigationProp = StackNavigationProp<MainStackParamList>;

export const ArtistListScreen: React.FC = () => {
  const navigation = useNavigation<NavigationProp>();
  const { user } = useAuth();
  
  const [artists, setArtists] = useState<Artist[]>([]);
  const [loading, setLoading] = useState(true);
  const [refreshing, setRefreshing] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [filterVisible, setFilterVisible] = useState(false);
  const [filters, setFilters] = useState<FilterOptions>({});
  const [hasMore, setHasMore] = useState(true);
  const [offset, setOffset] = useState(0);

  const loadArtists = useCallback(async (reset = false) => {
    if (!user) return;

    try {
      if (reset) {
        setLoading(true);
        setOffset(0);
      }

      const currentOffset = reset ? 0 : offset;
      const response = await apiService.getPersonalizedRecommendations({
        userId: user.id,
        limit: 20,
        offset: currentOffset,
        filters,
      });

      if (reset) {
        setArtists(response);
      } else {
        setArtists(prev => [...prev, ...response]);
      }

      setHasMore(response.length === 20);
      setOffset(currentOffset + response.length);
      setError(null);
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Failed to load artists';
      setError(errorMessage);
      
      if (reset) {
        setArtists([]);
      }
    } finally {
      setLoading(false);
      setRefreshing(false);
    }
  }, [user, filters, offset]);

  useEffect(() => {
    loadArtists(true);
  }, [filters]);

  const handleRefresh = () => {
    setRefreshing(true);
    loadArtists(true);
  };

  const handleLoadMore = () => {
    if (!loading && hasMore) {
      loadArtists(false);
    }
  };

  const handleArtistPress = async (artist: Artist) => {
    try {
      // Track interaction
      if (user) {
        await apiService.trackInteraction({
          userId: user.id,
          artistId: artist.id,
          type: 'view',
        });
      }
      
      navigation.navigate('ArtistDetail', { artist });
    } catch (error) {
      console.error('Error tracking interaction:', error);
      // Still navigate even if tracking fails
      navigation.navigate('ArtistDetail', { artist });
    }
  };

  const handleApplyFilters = (newFilters: FilterOptions) => {
    setFilters(newFilters);
    setFilterVisible(false);
  };

  const clearFilters = () => {
    setFilters({});
    setFilterVisible(false);
  };

  const renderArtist = ({ item }: { item: Artist }) => (
    <ArtistCard
      artist={item}
      onPress={() => handleArtistPress(item)}
    />
  );

  const renderFooter = () => {
    if (!loading || artists.length === 0) return null;
    return (
      <View style={styles.footerLoader}>
        <LoadingSpinner size="small" />
      </View>
    );
  };

  const renderEmpty = () => (
    <View style={styles.emptyContainer}>
      <Ionicons name="musical-notes-outline" size={64} color={COLORS.textSecondary} />
      <Text style={styles.emptyText}>No artists found</Text>
      <Text style={styles.emptySubtext}>
        Try adjusting your filters or check back later for new recommendations
      </Text>
      {Object.keys(filters).length > 0 && (
        <TouchableOpacity style={styles.clearFiltersButton} onPress={clearFilters}>
          <Text style={styles.clearFiltersText}>Clear Filters</Text>
        </TouchableOpacity>
      )}
    </View>
  );

  if (loading && artists.length === 0) {
    return <LoadingSpinner />;
  }

  return (
    <View style={styles.container}>
      <View style={styles.header}>
        <Text style={styles.title}>Recommended for You</Text>
        <TouchableOpacity
          style={styles.filterButton}
          onPress={() => setFilterVisible(true)}
        >
          <Ionicons 
            name="filter" 
            size={24} 
            color={Object.keys(filters).length > 0 ? COLORS.primary : COLORS.textSecondary} 
          />
        </TouchableOpacity>
      </View>

      {error ? (
        <ErrorMessage 
          message={error} 
          onRetry={() => loadArtists(true)} 
        />
      ) : (
        <FlatList
          data={artists}
          renderItem={renderArtist}
          keyExtractor={(item) => item.id}
          contentContainerStyle={styles.listContainer}
          refreshControl={
            <RefreshControl
              refreshing={refreshing}
              onRefresh={handleRefresh}
              colors={[COLORS.primary]}
            />
          }
          onEndReached={handleLoadMore}
          onEndReachedThreshold={0.1}
          ListFooterComponent={renderFooter}
          ListEmptyComponent={renderEmpty}
        />
      )}

      <FilterModal
        visible={filterVisible}
        filters={filters}
        onApply={handleApplyFilters}
        onClose={() => setFilterVisible(false)}
        onClear={clearFilters}
      />
    </View>
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
    paddingVertical: SPACING.sm,
    backgroundColor: COLORS.surface,
    borderBottomWidth: 1,
    borderBottomColor: COLORS.border,
  },
  title: {
    fontSize: 18,
    fontWeight: 'bold',
    color: COLORS.text,
  },
  filterButton: {
    padding: SPACING.sm,
  },
  listContainer: {
    paddingVertical: SPACING.sm,
  },
  footerLoader: {
    paddingVertical: SPACING.md,
    alignItems: 'center',
  },
  emptyContainer: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    paddingHorizontal: SPACING.xl,
    paddingVertical: SPACING.xl * 2,
  },
  emptyText: {
    fontSize: 18,
    fontWeight: '600',
    color: COLORS.text,
    marginTop: SPACING.md,
    textAlign: 'center',
  },
  emptySubtext: {
    fontSize: 14,
    color: COLORS.textSecondary,
    marginTop: SPACING.sm,
    textAlign: 'center',
    lineHeight: 20,
  },
  clearFiltersButton: {
    marginTop: SPACING.lg,
    paddingHorizontal: SPACING.lg,
    paddingVertical: SPACING.sm,
    backgroundColor: COLORS.primary,
    borderRadius: 8,
  },
  clearFiltersText: {
    color: COLORS.surface,
    fontWeight: '600',
  },
});

