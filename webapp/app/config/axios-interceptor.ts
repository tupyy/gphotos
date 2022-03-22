import axios from 'axios';

const TIMEOUT = 1 * 60 * 1000;
axios.defaults.timeout = TIMEOUT;
axios.defaults.baseURL = SERVER_API_URL;

const setupAxiosInterceptors = () => {
  const onResponseSuccess = response => response;

  axios.interceptors.response.use(onResponseSuccess);
};

export default setupAxiosInterceptors;
