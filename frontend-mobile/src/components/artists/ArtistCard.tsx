// src/components/artists/ArtistCard.tsx - Fixed to use helper functions
import React from 'react';
import {
  View,
  Text,
  TouchableOpacity,
  Image,
  StyleSheet,
} from 'react-native';
import { Ionicons } from '@expo/vector-icons';
import { Artist, getArtistLocation, getArtistBio } from '../../types';
import { globalStyles, theme } from '../../styles';

interface ArtistCardProps {
  artist: Artist;
  onPress: () => void;
}

export const ArtistCard: React.FC<ArtistCardProps> = ({ artist, onPress }) => {
  const renderGenres = () => {
    if (!artist.genres || artist.genres.length === 0) return null;
    
    const displayGenres = artist.genres.slice(0, 3);
    const remainingCount = artist.genres.length - 3;

    return (
      <View style={styles.genreContainer}>
        {displayGenres.map((genre, index) => (
          <View key={index} style={styles.genreTag}>
            <Text style={styles.genreText}>{genre}</Text>
          </View>
        ))}
        {remainingCount > 0 && (
          <View style={styles.genreTag}>
            <Text style={styles.genreText}>+{remainingCount}</Text>
          </View>
        )}
      </View>
    );
  };

  const renderRating = () => {
    if (!artist.rating) return null;
    return (
      <View style={[globalStyles.row, globalStyles.alignCenter]}>
        <Ionicons name="star" size={14} color={theme.colors.warning} />
        <Text style={[globalStyles.bodySmall, { marginLeft: theme.spacing.xs, fontWeight: theme.fontWeight.medium }]}>
          {artist.rating.toFixed(1)}
        </Text>
        {artist.ratingCount && (
          <Text style={[globalStyles.bodySmall, { marginLeft: theme.spacing.xs, color: theme.colors.textSecondary }]}>
            ({artist.ratingCount})
          </Text>
        )}
      </View>
    );
  };

  const renderContactInfo = () => {
    // Show manager info if available, otherwise show label info
    const managerName = artist.contactInfo?.manager || artist.manager;
    if (managerName) {
      return (
        <Text style={[globalStyles.bodySmall, { color: theme.colors.textSecondary, fontStyle: 'italic' }]}>
          Managed by {managerName}
        </Text>
      );
    }
    
    if (artist.contactInfo?.labelName) {
      return (
        <Text style={[globalStyles.bodySmall, { color: theme.colors.textSecondary, fontStyle: 'italic' }]}>
          {artist.contactInfo.labelName}
        </Text>
      );
    }
    
    return null;
  };

  return (
    <TouchableOpacity style={styles.container} onPress={onPress} activeOpacity={0.7}>
      <View style={styles.imageContainer}>
        {/* No imageUrl in API yet, so always show placeholder */}
        <View style={[styles.artistImage, styles.placeholderImage]}>
          <Ionicons name="person" size={40} color={theme.colors.textSecondary} />
        </View>
      </View>

      <View style={globalStyles.flex1}>
        <View style={[globalStyles.row, globalStyles.spaceBetween, globalStyles.alignCenter, globalStyles.mb_sm]}>
          <Text style={[globalStyles.h4, globalStyles.flex1]} numberOfLines={1}>
            {artist.name}
          </Text>
          {renderRating()}
        </View>

        <View style={[globalStyles.row, globalStyles.alignCenter, globalStyles.mb_sm]}>
          <Ionicons name="location-outline" size={14} color={theme.colors.textSecondary} />
          <Text style={[globalStyles.bodySmall, { marginLeft: theme.spacing.xs, flex: 1 }]} numberOfLines={1}>
            {getArtistLocation(artist)}
          </Text>
        </View>

        {renderGenres()}

        <Text style={[globalStyles.body, globalStyles.mb_sm]} numberOfLines={2}>
          {getArtistBio(artist)}
        </Text>

        <View style={[globalStyles.row, globalStyles.spaceBetween, globalStyles.alignCenter]}>
          {renderContactInfo()}
          <Ionicons name="chevron-forward" size={20} color={theme.colors.textSecondary} />
        </View>
      </View>
    </TouchableOpacity>
  );
};

const styles = StyleSheet.create({
  container: {
    flexDirection: 'row',
    backgroundColor: theme.colors.surface,
    marginHorizontal: theme.spacing.md,
    marginVertical: theme.spacing.xs,
    borderRadius: theme.borderRadius.lg,
    padding: theme.spacing.md,
    ...theme.shadows.md,
  },
  imageContainer: {
    marginRight: theme.spacing.md,
  },
  artistImage: {
    width: 80,
    height: 80,
    borderRadius: 40,
  },
  placeholderImage: {
    backgroundColor: theme.colors.borderLight,
    justifyContent: 'center',
    alignItems: 'center',
  },
  genreContainer: {
    flexDirection: 'row',
    flexWrap: 'wrap',
    marginBottom: theme.spacing.sm,
  },
  genreTag: {
    backgroundColor: theme.colors.primary,
    paddingHorizontal: theme.spacing.sm,
    paddingVertical: theme.spacing.xs / 2,
    borderRadius: theme.borderRadius.full,
    marginRight: theme.spacing.xs,
    marginBottom: theme.spacing.xs / 2,
  },
  genreText: {
    color: theme.colors.textInverse,
    fontSize: theme.fontSize.xs,
    fontWeight: theme.fontWeight.semibold,
  },
});