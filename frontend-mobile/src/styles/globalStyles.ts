// src/styles/globalStyles.ts
import { StyleSheet } from 'react-native';
import { theme } from './theme';

export const globalStyles = StyleSheet.create({
// Layout containers
container: {
    flex: 1,
    backgroundColor: theme.colors.background,
},

centerContainer: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    backgroundColor: theme.colors.background,
},

scrollContainer: {
    flexGrow: 1,
    paddingHorizontal: theme.spacing.lg,
},

// Flexbox utilities
row: {
    flexDirection: 'row',
},

column: {
    flexDirection: 'column',
},

justifyCenter: {
    justifyContent: 'center',
},

alignCenter: {
    alignItems: 'center',
},

spaceBetween: {
    justifyContent: 'space-between',
},

flex1: {
    flex: 1,
},

// Typography
h1: {
    fontSize: theme.fontSize['4xl'],
    fontWeight: theme.fontWeight.bold,
    color: theme.colors.text,
},

h2: {
    fontSize: theme.fontSize['3xl'],
    fontWeight: theme.fontWeight.bold,
    color: theme.colors.text,
},

h3: {
    fontSize: theme.fontSize['2xl'],
    fontWeight: theme.fontWeight.semibold,
    color: theme.colors.text,
},

h4: {
    fontSize: theme.fontSize.xl,
    fontWeight: theme.fontWeight.semibold,
    color: theme.colors.text,
},

bodyLarge: {
    fontSize: theme.fontSize.lg,
    color: theme.colors.text,
},

body: {
    fontSize: theme.fontSize.md,
    color: theme.colors.text,
},

bodySmall: {
    fontSize: theme.fontSize.sm,
    color: theme.colors.textSecondary,
},

caption: {
    fontSize: theme.fontSize.xs,
    color: theme.colors.textSecondary,
},

// Common input styles
inputContainer: {
    marginBottom: theme.spacing.lg,
},

inputLabel: {
    fontSize: theme.fontSize.md,
    fontWeight: theme.fontWeight.semibold,
    color: theme.colors.text,
    marginBottom: theme.spacing.sm,
},

inputWrapper: {
    flexDirection: 'row',
    alignItems: 'center',
    borderWidth: 1,
    borderColor: theme.colors.border,
    borderRadius: theme.borderRadius.lg,
    backgroundColor: theme.colors.surface,
    paddingHorizontal: theme.spacing.md,
},

inputWrapperError: {
    borderColor: theme.colors.error,
},

textInput: {
    flex: 1,
    paddingVertical: theme.spacing.md,
    fontSize: theme.fontSize.md,
    color: theme.colors.text,
},

inputIcon: {
    marginRight: theme.spacing.sm,
},

errorText: {
    color: theme.colors.error,
    fontSize: theme.fontSize.sm,
    marginTop: theme.spacing.xs,
},

// Common button styles
button: {
    paddingVertical: theme.spacing.md,
    paddingHorizontal: theme.spacing.lg,
    borderRadius: theme.borderRadius.lg,
    alignItems: 'center',
    justifyContent: 'center',
    flexDirection: 'row',
},

buttonPrimary: {
    backgroundColor: theme.colors.primary,
},

buttonSecondary: {
    backgroundColor: theme.colors.surface,
    borderWidth: 1,
    borderColor: theme.colors.border,
},

buttonOutline: {
    backgroundColor: 'transparent',
    borderWidth: 2,
    borderColor: theme.colors.primary,
},

buttonDisabled: {
    opacity: 0.6,
},

buttonText: {
    fontSize: theme.fontSize.md,
    fontWeight: theme.fontWeight.semibold,
},

buttonTextPrimary: {
    color: theme.colors.textInverse,
},

buttonTextSecondary: {
    color: theme.colors.text,
},

buttonTextOutline: {
    color: theme.colors.primary,
},

// Card styles
card: {
    backgroundColor: theme.colors.surface,
    borderRadius: theme.borderRadius.lg,
    padding: theme.spacing.md,
    ...theme.shadows.md,
},

// Spacing utilities (most commonly used)
mt_sm: { marginTop: theme.spacing.sm },
mt_md: { marginTop: theme.spacing.md },
mt_lg: { marginTop: theme.spacing.lg },

mb_sm: { marginBottom: theme.spacing.sm },
mb_md: { marginBottom: theme.spacing.md },
mb_lg: { marginBottom: theme.spacing.lg },

mx_md: { marginHorizontal: theme.spacing.md },
mx_lg: { marginHorizontal: theme.spacing.lg },

my_sm: { marginVertical: theme.spacing.sm },
my_md: { marginVertical: theme.spacing.md },

p_md: { padding: theme.spacing.md },
p_lg: { padding: theme.spacing.lg },

px_md: { paddingHorizontal: theme.spacing.md },
px_lg: { paddingHorizontal: theme.spacing.lg },

py_sm: { paddingVertical: theme.spacing.sm },
py_md: { paddingVertical: theme.spacing.md },
});