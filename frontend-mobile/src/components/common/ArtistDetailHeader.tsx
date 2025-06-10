// src/components/artists/ArtistDetailHeader.tsx
import React from 'react';
import { View, Text, StyleSheet, Image } from 'react-native';
import { Ionicons } from '@expo/vector-icons';
import { Artist } from '../../types';
import { COLORS, SPACING } from '../../utils/constants';

interface ArtistDetailHeaderProps {
  artist: Artist;
}

export const ArtistDetailHeader: React.FC<ArtistDetailHeaderProps> = ({ artist }) => {
  const renderGenres = () => (
    <View style={styles.genreContainer}>
      {artist.genre.map((genre, index) => (
        <View key={index} style={styles.genreTag}>
          <Text style={styles.genreText}>{genre}</Text>
        </View>
      ))}
    </View>
  );

  return (
    <View style={styles.container}>
      <View style={styles.imageContainer}>
        {artist.imageUrl ? (
          <Image source={{ uri: artist.imageUrl }} style={styles.artistImage} />
        ) : (
          <View style={[styles.artistImage, styles.placeholderImage]}>
            <Ionicons name="person" size={80} color={COLORS.textSecondary} />
          </View>
        )}
      </View>

      <View style={styles.infoContainer}>
        <Text style={styles.artistName}>{artist.name}</Text>
        
        <View style={styles.locationContainer}>
          <Ionicons name="location-outline" size={16} color={COLORS.textSecondary} />
          <Text style={styles.locationText}>{artist.location}</Text>
        </View>

        {artist.rating && (
          <View style={styles.ratingContainer}>
            <Ionicons name="star" size={16} color={COLORS.warning} />
            <Text style={styles.ratingText}>{artist.rating.toFixed(1)} rating</Text>
          </View>
        )}

        {renderGenres()}
      </View>
    </View>
  );
};

const styles = StyleSheet.create({
  container: {
    backgroundColor: COLORS.surface,
    padding: SPACING.md,
    alignItems: 'center',
  },
  imageContainer: {
    marginBottom: SPACING.md,
  },
  artistImage: {
    width: 120,
    height: 120,
    borderRadius: 60,
  },
  placeholderImage: {
    backgroundColor: COLORS.border,
    justifyContent: 'center',
    alignItems: 'center',
  },
  infoContainer: {
    alignItems: 'center',
  },
  artistName: {
    fontSize: 24,
    fontWeight: 'bold',
    color: COLORS.text,
    marginBottom: SPACING.sm,
    textAlign: 'center',
  },
  locationContainer: {
    flexDirection: 'row',
    alignItems: 'center',
    marginBottom: SPACING.sm,
  },
  locationText: {
    fontSize: 16,
    color: COLORS.textSecondary,
    marginLeft: SPACING.xs,
  },
  ratingContainer: {
    flexDirection: 'row',
    alignItems: 'center',
    marginBottom: SPACING.md,
  },
  ratingText: {
    fontSize: 14,
    color: COLORS.text,
    marginLeft: SPACING.xs,
    fontWeight: '500',
  },
  genreContainer: {
    flexDirection: 'row',
    flexWrap: 'wrap',
    justifyContent: 'center',
  },
  genreTag: {
    backgroundColor: COLORS.primary,
    paddingHorizontal: SPACING.sm,
    paddingVertical: SPACING.xs,
    borderRadius: 16,
    marginRight: SPACING.xs,
    marginBottom: SPACING.xs,
  },
  genreText: {
    color: COLORS.surface,
    fontSize: 12,
    fontWeight: '600',
  },
});