// src/components/artists/ArtistDetailHeader.tsx
import React from 'react';
import { View, Text, StyleSheet, Image, TouchableOpacity, Linking, Alert } from 'react-native';
import { Ionicons } from '@expo/vector-icons';
import { Artist } from '../../types';
import { theme } from '../../styles';

interface ArtistDetailHeaderProps {
  artist: Artist;
}

type StreamingService = {
  name: 'bandcamp' | 'spotify' | 'appleMusic' | 'instagram';
  url: string;
  icon: keyof typeof Ionicons.glyphMap;
  color: string;
  displayName: string;
};

export const ArtistDetailHeader = ({ artist }: ArtistDetailHeaderProps) => {
  const renderGenres = () => (
    <View style={styles.genreContainer}>
      {artist.genre.map((genre, index) => (
        <View key={index} style={styles.genreTag}>
          <Text style={styles.genreText}>{genre}</Text>
        </View>
      ))}
    </View>
  );

  const handleStreamingServicePress = async (service: StreamingService) => {
    try {
      const supported = await Linking.canOpenURL(service.url);
      if (supported) {
        await Linking.openURL(service.url);
        
        // Optional: Track the interaction (won't break if tracking fails)
        // You can uncomment this when you want to add analytics
        // trackStreamingClick?.(artist.id, service.name);
      } else {
        Alert.alert(
          'Cannot Open Link',
          `Unable to open ${service.displayName}. Please check if the app is installed.`
        );
      }
    } catch (error) {
      console.error('Error opening URL:', error);
      Alert.alert(
        'Error',
        `Sorry, there was an error opening ${service.displayName}.`
      );
    }
  };

  const getAvailableStreamingServices = (): StreamingService[] => {
    const services: StreamingService[] = [];
    const social = artist.contactInfo.social;

    if (!social) return services;

    // Bandcamp - prioritized first
    if (social.bandcamp) {
      services.push({
        name: 'bandcamp',
        url: social.bandcamp,
        icon: 'musical-notes',
        color: '#1DA0C3',
        displayName: 'Bandcamp',
      });
    }

    // Spotify
    if (social.spotify) {
      services.push({
        name: 'spotify',
        url: social.spotify,
        icon: 'play-circle',
        color: '#1DB954',
        displayName: 'Spotify',
      });
    }

    // Apple Music
    if (social.appleMusic) {
      services.push({
        name: 'appleMusic',
        url: social.appleMusic,
        icon: 'play',
        color: '#FA243C',
        displayName: 'Apple Music',
      });
    }

    // Instagram
    if (social.instagram) {
      services.push({
        name: 'instagram',
        url: social.instagram,
        icon: 'logo-instagram',
        color: '#E4405F',
        displayName: 'Instagram',
      });
    }

    return services;
  };

  const renderStreamingServices = () => {
    const availableServices = getAvailableStreamingServices();

    if (availableServices.length === 0) {
      return null;
    }

    return (
      <View style={styles.streamingServicesContainer}>
        <Text style={styles.streamingServicesTitle}>Listen & Follow</Text>
        <View style={styles.streamingServices}>
          {availableServices.map((service) => (
            <TouchableOpacity
              key={service.name}
              style={[styles.serviceButton, { borderColor: service.color }]}
              onPress={() => handleStreamingServicePress(service)}
              activeOpacity={0.7}
            >
              <Ionicons 
                name={service.icon} 
                size={24} 
                color={service.color} 
              />
            </TouchableOpacity>
          ))}
        </View>
      </View>
    );
  };

  return (
    <View style={styles.container}>
      <View style={styles.imageContainer}>
        {artist.imageUrl ? (
          <Image source={{ uri: artist.imageUrl }} style={styles.artistImage} />
        ) : (
          <View style={[styles.artistImage, styles.placeholderImage]}>
            <Ionicons name="person" size={80} color={theme.colors.textSecondary} />
          </View>
        )}
      </View>

      <View style={styles.infoContainer}>
        <Text style={styles.artistName}>{artist.name}</Text>
        
        <View style={styles.locationContainer}>
          <Ionicons name="location-outline" size={16} color={theme.colors.textSecondary} />
          <Text style={styles.locationText}>{artist.location}</Text>
        </View>

        {artist.rating && (
          <View style={styles.ratingContainer}>
            <Ionicons name="star" size={16} color={theme.colors.warning} />
            <Text style={styles.ratingText}>{artist.rating.toFixed(1)} rating</Text>
          </View>
        )}

        {renderGenres()}
        {renderStreamingServices()}
      </View>
    </View>
  );
};

const styles = StyleSheet.create({
  container: {
    backgroundColor: theme.colors.surface,
    padding: theme.spacing.md,
    alignItems: 'center',
  },
  imageContainer: {
    marginBottom: theme.spacing.md,
  },
  artistImage: {
    width: 120,
    height: 120,
    borderRadius: 60,
  },
  placeholderImage: {
    backgroundColor: theme.colors.border,
    justifyContent: 'center',
    alignItems: 'center',
  },
  infoContainer: {
    alignItems: 'center',
  },
  artistName: {
    fontSize: 24,
    fontWeight: 'bold',
    color: theme.colors.text,
    marginBottom: theme.spacing.sm,
    textAlign: 'center',
  },
  locationContainer: {
    flexDirection: 'row',
    alignItems: 'center',
    marginBottom: theme.spacing.sm,
  },
  locationText: {
    fontSize: 16,
    color: theme.colors.textSecondary,
    marginLeft: theme.spacing.xs,
  },
  ratingContainer: {
    flexDirection: 'row',
    alignItems: 'center',
    marginBottom: theme.spacing.md,
  },
  ratingText: {
    fontSize: 14,
    color: theme.colors.text,
    marginLeft: theme.spacing.xs,
    fontWeight: '500',
  },
  genreContainer: {
    flexDirection: 'row',
    flexWrap: 'wrap',
    justifyContent: 'center',
    marginBottom: theme.spacing.md,
  },
  genreTag: {
    backgroundColor: theme.colors.primary,
    paddingHorizontal: theme.spacing.sm,
    paddingVertical: theme.spacing.xs,
    borderRadius: 16,
    marginRight: theme.spacing.xs,
    marginBottom: theme.spacing.xs,
  },
  genreText: {
    color: theme.colors.surface,
    fontSize: 12,
    fontWeight: '600',
  },
  streamingServicesContainer: {
    alignItems: 'center',
    width: '100%',
  },
  streamingServicesTitle: {
    fontSize: 16,
    fontWeight: '600',
    color: theme.colors.text,
    marginBottom: theme.spacing.sm,
  },
  streamingServices: {
    flexDirection: 'row',
    justifyContent: 'center',
    flexWrap: 'wrap',
    gap: theme.spacing.md,
  },
  serviceButton: {
    width: 56,
    height: 56,
    borderRadius: 28,
    borderWidth: 2,
    backgroundColor: theme.colors.surface,
    justifyContent: 'center',
    alignItems: 'center',
    ...theme.shadows.sm,
  },
});