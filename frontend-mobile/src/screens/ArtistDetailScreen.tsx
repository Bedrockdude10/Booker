// src/screens/artists/ArtistDetailScreen.tsx - Fixed to use helper functions and correct property names
import React, { useState, useEffect } from 'react';
import {
  View,
  Text,
  ScrollView,
  TouchableOpacity,
  StyleSheet,
  Image,
  Linking,
  Alert,
  Dimensions,
} from 'react-native';
import { Ionicons } from '@expo/vector-icons';
import { RouteProp, useRoute } from '@react-navigation/native';
import { useAuth } from '../contexts/AuthContext';
import { apiService } from '../services/api';
import { COLORS, SPACING } from '../utils/constants';
import { 
  MainStackParamList, 
  getArtistId, 
  getArtistLocation, 
  getArtistBio, 
  getArtistEmail, 
  getArtistWebsite,
  getArtistSocialLinks 
} from '../types';

type ArtistDetailRouteProp = RouteProp<MainStackParamList, 'ArtistDetail'>;

const { width } = Dimensions.get('window');

export const ArtistDetailScreen: React.FC = () => {
  const route = useRoute<ArtistDetailRouteProp>();
  const { user } = useAuth();
  const { artist } = route.params;

  const [isSaved, setIsSaved] = useState(false);
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    // Track that user viewed this artist detail
    trackView();
  }, []);

  const trackView = async () => {
    if (!user) return;
    
    try {
      await apiService.trackInteraction({
        userId: user.id,
        artistId: getArtistId(artist),
        type: 'view',
      });
    } catch (error) {
      console.error('Error tracking view:', error);
    }
  };

  const handleContact = async () => {
    try {
      setLoading(true);
      
      // Track contact interaction
      if (user) {
        await apiService.trackInteraction({
          userId: user.id,
          artistId: getArtistId(artist),
          type: 'contact',
        });
      }

      const email = getArtistEmail(artist);
      if (!email) {
        Alert.alert(
          'Contact Information',
          'No email address available for this artist. Please check their social media or website for contact details.',
          [{ text: 'OK' }]
        );
        return;
      }

      const subject = `Booking Inquiry - ${artist.name}`;
      const body = `Hi ${artist.name},\n\nI'm interested in discussing a potential booking opportunity.\n\nBest regards,\n${user?.name}`;
      
      const mailtoUrl = `mailto:${email}?subject=${encodeURIComponent(subject)}&body=${encodeURIComponent(body)}`;
      
      const canOpen = await Linking.canOpenURL(mailtoUrl);
      if (canOpen) {
        await Linking.openURL(mailtoUrl);
      } else {
        Alert.alert(
          'Contact Artist',
          `Email: ${email}`,
          [
            { text: 'Copy Email', onPress: () => copyToClipboard(email) },
            { text: 'OK' },
          ]
        );
      }
    } catch (error) {
      Alert.alert('Error', 'Unable to open email client');
    } finally {
      setLoading(false);
    }
  };

  const copyToClipboard = (text: string) => {
    // Note: You'll need to install expo-clipboard for this to work
    // For now, just show an alert
    Alert.alert('Email Copied', text);
  };

  const handleSocialLink = async (url: string, platform: string) => {
    try {
      const canOpen = await Linking.canOpenURL(url);
      if (canOpen) {
        await Linking.openURL(url);
      } else {
        Alert.alert('Error', `Cannot open ${platform} link`);
      }
    } catch (error) {
      Alert.alert('Error', `Failed to open ${platform}`);
    }
  };

  const renderGenres = () => (
    <View style={styles.genreContainer}>
      {artist.genres.map((genre, index) => (
        <View key={index} style={styles.genreTag}>
          <Text style={styles.genreText}>{genre}</Text>
        </View>
      ))}
    </View>
  );

  const renderSocialLinks = () => {
    const socialLinks = getArtistSocialLinks(artist);
    if (socialLinks.length === 0) return null;

    return (
      <View style={styles.socialContainer}>
        <Text style={styles.sectionTitle}>Social Media & Streaming</Text>
        <View style={styles.socialLinks}>
          {socialLinks.map((link) => (
            <TouchableOpacity
              key={link.key}
              style={[styles.socialButton, { backgroundColor: link.color }]}
              onPress={() => handleSocialLink(link.url!, link.key)}
            >
              <Ionicons name={link.icon as any} size={24} color={COLORS.surface} />
            </TouchableOpacity>
          ))}
        </View>
      </View>
    );
  };

  const renderContactInfo = () => {
    const email = getArtistEmail(artist);
    const website = getArtistWebsite(artist);
    const managerName = artist.contactInfo?.manager || artist.manager;
    const labelName = artist.contactInfo?.labelName;
    
    const hasContactInfo = email || website || managerName || labelName;
    if (!hasContactInfo) return null;

    return (
      <View style={styles.section}>
        <Text style={styles.sectionTitle}>Contact Information</Text>
        
        {email && (
          <View style={styles.contactRow}>
            <Ionicons name="mail-outline" size={20} color={COLORS.textSecondary} />
            <Text style={styles.contactText}>{email}</Text>
          </View>
        )}
        
        {website && (
          <TouchableOpacity
            style={styles.contactRow}
            onPress={() => handleSocialLink(website, 'Website')}
          >
            <Ionicons name="globe-outline" size={20} color={COLORS.primary} />
            <Text style={[styles.contactText, { color: COLORS.primary }]}>
              {website}
            </Text>
          </TouchableOpacity>
        )}
        
        {managerName && (
          <View style={styles.contactRow}>
            <Ionicons name="person-outline" size={20} color={COLORS.textSecondary} />
            <Text style={styles.contactText}>
              Manager: {managerName}
            </Text>
          </View>
        )}
        
        {labelName && (
          <View style={styles.contactRow}>
            <Ionicons name="business-outline" size={20} color={COLORS.textSecondary} />
            <Text style={styles.contactText}>
              Label: {labelName}
            </Text>
          </View>
        )}
      </View>
    );
  };

  const renderRating = () => {
    if (!artist.rating) return null;
    
    return (
      <View style={styles.ratingContainer}>
        <View style={styles.ratingStars}>
          <Ionicons name="star" size={20} color="#FFD700" />
          <Text style={styles.ratingText}>{artist.rating.toFixed(1)}</Text>
          {artist.ratingCount && (
            <Text style={styles.ratingCount}>({artist.ratingCount} reviews)</Text>
          )}
        </View>
      </View>
    );
  };

  return (
    <ScrollView style={styles.container} showsVerticalScrollIndicator={false}>
      {/* Artist Image */}
      <View style={styles.imageContainer}>
        {/* No imageUrl in API yet, so always show placeholder */}
        <View style={[styles.artistImage, styles.placeholderImage]}>
          <Ionicons name="person" size={80} color={COLORS.textSecondary} />
        </View>
      </View>

      <View style={styles.contentContainer}>
        {/* Artist Name, Location & Rating */}
        <View style={styles.headerInfo}>
          <Text style={styles.artistName}>{artist.name}</Text>
          <View style={styles.locationContainer}>
            <Ionicons name="location-outline" size={16} color={COLORS.textSecondary} />
            <Text style={styles.locationText}>{getArtistLocation(artist)}</Text>
          </View>
          {renderRating()}
        </View>

        {/* Genres */}
        {renderGenres()}

        {/* Bio */}
        <View style={styles.section}>
          <Text style={styles.sectionTitle}>About</Text>
          <Text style={styles.bioText}>{getArtistBio(artist)}</Text>
        </View>

        {/* Contact Info */}
        {renderContactInfo()}

        {/* Social Media */}
        {renderSocialLinks()}

        {/* Contact Button */}
        {getArtistEmail(artist) && (
          <TouchableOpacity
            style={[styles.contactButton, loading && styles.contactButtonDisabled]}
            onPress={handleContact}
            disabled={loading}
          >
            <Ionicons 
              name="mail" 
              size={20} 
              color={COLORS.surface} 
              style={styles.contactButtonIcon}
            />
            <Text style={styles.contactButtonText}>
              {loading ? 'Opening...' : 'Contact Artist'}
            </Text>
          </TouchableOpacity>
        )}
      </View>
    </ScrollView>
  );
};

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: COLORS.background,
  },
  imageContainer: {
    width: '100%',
    height: 250,
    backgroundColor: COLORS.border,
  },
  artistImage: {
    width: '100%',
    height: '100%',
  },
  placeholderImage: {
    justifyContent: 'center',
    alignItems: 'center',
    backgroundColor: COLORS.border,
  },
  contentContainer: {
    padding: SPACING.md,
  },
  headerInfo: {
    marginBottom: SPACING.md,
  },
  artistName: {
    fontSize: 28,
    fontWeight: 'bold',
    color: COLORS.text,
    marginBottom: SPACING.xs,
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
    marginTop: SPACING.xs,
  },
  ratingStars: {
    flexDirection: 'row',
    alignItems: 'center',
  },
  ratingText: {
    fontSize: 18,
    fontWeight: '600',
    color: COLORS.text,
    marginLeft: SPACING.xs,
  },
  ratingCount: {
    fontSize: 14,
    color: COLORS.textSecondary,
    marginLeft: SPACING.xs,
  },
  genreContainer: {
    flexDirection: 'row',
    flexWrap: 'wrap',
    marginBottom: SPACING.lg,
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
  section: {
    marginBottom: SPACING.lg,
  },
  sectionTitle: {
    fontSize: 18,
    fontWeight: '600',
    color: COLORS.text,
    marginBottom: SPACING.sm,
  },
  bioText: {
    fontSize: 16,
    color: COLORS.text,
    lineHeight: 24,
  },
  contactRow: {
    flexDirection: 'row',
    alignItems: 'center',
    marginBottom: SPACING.sm,
  },
  contactText: {
    fontSize: 16,
    color: COLORS.text,
    marginLeft: SPACING.sm,
    flex: 1,
  },
  socialContainer: {
    marginBottom: SPACING.lg,
  },
  socialLinks: {
    flexDirection: 'row',
    flexWrap: 'wrap',
    gap: SPACING.sm,
  },
  socialButton: {
    width: 48,
    height: 48,
    borderRadius: 24,
    justifyContent: 'center',
    alignItems: 'center',
  },
  contactButton: {
    backgroundColor: COLORS.primary,
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'center',
    paddingVertical: SPACING.md,
    borderRadius: 12,
    marginTop: SPACING.md,
    marginBottom: SPACING.xl,
  },
  contactButtonDisabled: {
    opacity: 0.6,
  },
  contactButtonIcon: {
    marginRight: SPACING.sm,
  },
  contactButtonText: {
    color: COLORS.surface,
    fontSize: 18,
    fontWeight: '600',
  },
});