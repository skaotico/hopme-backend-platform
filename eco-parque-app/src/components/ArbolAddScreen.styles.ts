import { StyleSheet, Platform, StatusBar } from 'react-native';

export const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#061114',
    paddingTop: Platform.OS === 'ios' ? 60 : (StatusBar.currentHeight || 24) + 20,
  },
  header: {
    flexDirection: 'row',
    alignItems: 'center',
    paddingHorizontal: 24,
    marginBottom: 32,
  },
  backButton: {
    width: 40,
    height: 40,
    borderRadius: 20,
    backgroundColor: '#112226',
    borderWidth: 1,
    borderColor: '#1D343B',
    justifyContent: 'center',
    alignItems: 'center',
    marginRight: 16,
  },
  title: {
    fontSize: 24,
    fontWeight: '800',
    color: '#FFFFFF',
    letterSpacing: -0.5,
  },
  formContent: {
    paddingHorizontal: 24,
  },
  inputContainer: {
    flexDirection: 'row',
    alignItems: 'center',
    backgroundColor: '#112226',
    borderWidth: 1,
    borderColor: '#1D343B',
    borderRadius: 16,
    marginBottom: 16,
    paddingHorizontal: 16,
    height: 60,
  },
  textAreaContainer: {
    height: 100,
    alignItems: 'flex-start',
    paddingTop: 16,
    paddingBottom: 12,
  },
  icon: {
    marginRight: 12,
  },
  textAreaIcon: {
    marginTop: 2,
  },
  input: {
    flex: 1,
    fontSize: 16,
    fontWeight: '600',
    color: '#FFFFFF',
  },
  textAreaInput: {
    height: '100%',
    textAlignVertical: 'top',
  },
  locationInfo: {
    flexDirection: 'row',
    alignItems: 'center',
    marginBottom: 24,
    paddingHorizontal: 4,
    backgroundColor: 'rgba(20, 184, 166, 0.1)',
    padding: 12,
    borderRadius: 12,
    borderWidth: 1,
    borderColor: 'rgba(20, 184, 166, 0.2)',
  },
  locationText: {
    color: '#14B8A6',
    fontSize: 14,
    marginLeft: 8,
    fontWeight: '500',
  },
  primaryButton: {
    backgroundColor: '#14B8A6',
    height: 60,
    borderRadius: 16,
    alignItems: 'center',
    justifyContent: 'center',
    marginTop: 8,
  },
  primaryButtonText: {
    color: '#061114',
    fontSize: 16,
    fontWeight: 'bold',
  },
  
  // Custom Selector Styles
  selectorButton: {
    flexDirection: 'row',
    alignItems: 'center',
    backgroundColor: '#112226',
    borderWidth: 1,
    borderColor: '#1D343B',
    borderRadius: 16,
    marginBottom: 16,
    paddingHorizontal: 16,
    height: 60,
    justifyContent: 'space-between',
  },
  selectorLeft: {
    flexDirection: 'row',
    alignItems: 'center',
  },
  selectorValue: {
    fontSize: 16,
    fontWeight: '600',
    color: '#FFFFFF',
    marginLeft: 12,
  },
  selectorPlaceholder: {
    fontSize: 16,
    fontWeight: '600',
    color: '#8FA3A9',
    marginLeft: 12,
  },
  
  // Modal layout
  modalOverlay: {
    flex: 1,
    backgroundColor: 'rgba(0, 0, 0, 0.7)',
    justifyContent: 'flex-end',
  },
  modalContent: {
    backgroundColor: '#061114',
    borderTopLeftRadius: 32,
    borderTopRightRadius: 32,
    maxHeight: '80%',
    paddingBottom: 40,
    borderWidth: 1,
    borderColor: '#1D343B',
  },
  modalHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    padding: 24,
    borderBottomWidth: 1,
    borderBottomColor: '#112226',
  },
  modalTitle: {
    fontSize: 20,
    fontWeight: '800',
    color: '#FFFFFF',
  },
  modalCloseButton: {
    padding: 4,
  },
  modalList: {
    padding: 24,
  },
  modalItem: {
    paddingVertical: 16,
    paddingHorizontal: 20,
    backgroundColor: '#112226',
    borderRadius: 16,
    marginBottom: 12,
    borderWidth: 1,
    borderColor: '#1D343B',
  },
  modalItemActive: {
    borderColor: '#14B8A6',
    backgroundColor: 'rgba(20, 184, 166, 0.1)',
  },
  modalItemText: {
    fontSize: 16,
    fontWeight: '700',
    color: '#FFFFFF',
  },
  modalItemSubtext: {
    fontSize: 13,
    color: '#8FA3A9',
    marginTop: 4,
    fontWeight: '500',
  },
});
