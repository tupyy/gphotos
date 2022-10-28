import React from 'react';
import ReactDOM from 'react-dom';
import AppComponent from './app';
import getStore from './config/store';
import { Provider } from 'react-redux';
import '@elastic/eui/dist/eui_theme_dark.css';

import { EuiProvider } from '@elastic/eui';

const store = getStore();

const rootEl = document.getElementById('root');

const render = Component =>
  // eslint-disable-next-line react/no-render-return-value
  ReactDOM.render(
      <EuiProvider colorMode='dark'>
      <Provider store={store}>
        <div>
          <Component />
        </div>
      </Provider>
    </EuiProvider>,
    rootEl
  );

render(AppComponent);
