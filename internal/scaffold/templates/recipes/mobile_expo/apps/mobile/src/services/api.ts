import axios from 'axios';
import { Platform } from 'react-native';

// In Android Emulator, localhost is 10.0.2.2. In iOS Simulator, it is localhost.
const BASE_URL = Platform.select({
  android: 'http://10.0.2.2:8080/api',
  ios: 'http://localhost:8080/api',
  default: 'http://localhost:8080/api',
});

export const api = axios.create({
  baseURL: BASE_URL,
  timeout: 5000,
});
