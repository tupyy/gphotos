import { loadingBarReducer as loadingBar } from 'react-redux-loading-bar';

import userManagement from './user-management';
import albumManagement from './album-management';

const rootReducer = {
  userManagement,
  albumManagement,
  loadingBar,
};

export default rootReducer;
