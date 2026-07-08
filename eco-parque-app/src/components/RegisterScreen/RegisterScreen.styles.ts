import { StyleSheet } from 'react-native';

export const styles = StyleSheet.create({
  container: {
    flex: 1,
    justifyContent: 'center',
    backgroundColor: '#061114',
    padding: 24,
  },
  card: {
    backgroundColor: '#112226',
    borderRadius: 32,
    padding: 28,
    borderWidth: 1,
    borderColor: '#1D343B',
  },
  header: {
    alignItems: 'center',
    marginBottom: 36,
  },
  iconContainer: {
    backgroundColor: '#1B333A',
    width: 72,
    height: 72,
    borderRadius: 36,
    justifyContent: 'center',
    alignItems: 'center',
    marginBottom: 16,
    borderWidth: 1,
    borderColor: '#26464F',
  },
  title: {
    fontSize: 28,
    fontWeight: '800',
    color: '#FFFFFF',
    letterSpacing: -0.5,
  },
  subtitle: {
    fontSize: 15,
    color: '#8FA3A9',
    marginTop: 8,
  },
  inputContainer: {
    flexDirection: 'row',
    alignItems: 'center',
    backgroundColor: '#061114',
    borderWidth: 1,
    borderColor: '#1D343B',
    borderRadius: 16,
    marginBottom: 16,
    paddingHorizontal: 16,
    height: 60,
  },
  icon: {
    marginRight: 12,
  },
  input: {
    flex: 1,
    fontSize: 16,
    color: '#FFFFFF',
  },
  primaryButton: {
    backgroundColor: '#14B8A6',
    height: 60,
    borderRadius: 16,
    alignItems: 'center',
    justifyContent: 'center',
    marginTop: 12,
  },
  primaryButtonText: {
    color: '#061114',
    fontSize: 16,
    fontWeight: 'bold',
  },
  secondaryButton: {
    alignItems: 'center',
    marginTop: 28,
  },
  secondaryButtonText: {
    color: '#8FA3A9',
    fontSize: 14,
  },
  linkText: {
    color: '#14B8A6',
    fontWeight: '700',
  },
});
