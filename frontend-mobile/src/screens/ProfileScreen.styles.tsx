import { StyleSheet } from 'react-native';
import { COLORS, SPACING } from '../utils/constants';

export const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: COLORS.background,
  },
  header: {
    alignItems: 'center',
    paddingVertical: SPACING.xl,
    backgroundColor: COLORS.surface,
    borderBottomWidth: 1,
    borderBottomColor: COLORS.border,
  },
  avatarContainer: {
    width: 100,
    height: 100,
    borderRadius: 50,
    backgroundColor: COLORS.border,
    justifyContent: 'center',
    alignItems: 'center',
    marginBottom: SPACING.md,
  },
  userName: {
    fontSize: 24,
    fontWeight: 'bold',
    color: COLORS.text,
    marginBottom: SPACING.xs,
  },
  userRole: {
    fontSize: 16,
    color: COLORS.textSecondary,
    fontWeight: '500',
  },
  formContainer: {
    padding: SPACING.md,
  },
  fieldContainer: {
    marginBottom: SPACING.lg,
  },
  fieldLabel: {
    fontSize: 16,
    fontWeight: '600',
    color: COLORS.text,
    marginBottom: SPACING.sm,
  },
  fieldValue: {
    fontSize: 16,
    color: COLORS.text,
    paddingVertical: SPACING.sm,
  },
  fieldValueEmpty: {
    color: COLORS.textSecondary,
    fontStyle: 'italic',
  },
  textInput: {
    borderWidth: 1,
    borderColor: COLORS.border,
    borderRadius: 8,
    paddingHorizontal: SPACING.md,
    paddingVertical: SPACING.sm,
    fontSize: 16,
    color: COLORS.text,
    backgroundColor: COLORS.surface,
  },
  textInputMultiline: {
    minHeight: 100,
    textAlignVertical: 'top',
  },
  buttonContainer: {
    padding: SPACING.md,
  },
  editingButtons: {
    flexDirection: 'row',
    gap: SPACING.md,
    marginBottom: SPACING.md,
  },
  button: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'center',
    paddingVertical: SPACING.md,
    borderRadius: 12,
    marginBottom: SPACING.sm,
  },
  editButton: {
    backgroundColor: COLORS.surface,
    borderWidth: 2,
    borderColor: COLORS.primary,
  },
  editButtonText: {
    color: COLORS.primary,
    fontSize: 16,
    fontWeight: '600',
    marginLeft: SPACING.sm,
  },
  saveButton: {
    backgroundColor: COLORS.primary,
    flex: 1,
  },
  saveButtonText: {
    color: COLORS.surface,
    fontSize: 16,
    fontWeight: '600',
    marginLeft: SPACING.sm,
  },
  cancelButton: {
    backgroundColor: COLORS.surface,
    borderWidth: 1,
    borderColor: COLORS.border,
    flex: 1,
  },
  cancelButtonText: {
    color: COLORS.text,
    fontSize: 16,
    fontWeight: '600',
  },
  logoutButton: {
    backgroundColor: COLORS.surface,
    borderWidth: 2,
    borderColor: COLORS.error,
  },
  logoutButtonText: {
    color: COLORS.error,
    fontSize: 16,
    fontWeight: '600',
    marginLeft: SPACING.sm,
  },
  buttonDisabled: {
    opacity: 0.6,
  },
  footer: {
    alignItems: 'center',
    paddingVertical: SPACING.xl,
  },
  footerText: {
    fontSize: 14,
    color: COLORS.textSecondary,
  },
});