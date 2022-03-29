import axios from 'axios';
import { handleDates } from './date';

const TIMEOUT = 1 * 60 * 1000;
axios.defaults.timeout = TIMEOUT;
axios.defaults.baseURL = SERVER_API_URL;

const setupAxiosInterceptors = () => {
  const onResponseSuccess = response => handleDates(response);

  axios.interceptors.response.use(onResponseSuccess);
};

export default setupAxiosInterceptors;
