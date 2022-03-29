import React from 'react';
import ReactDOM from 'react-dom';
import ErrorBoundary from './shared/error/error-boundary';
import AppComponent from './app';
import getStore from './config/store';
import { Provider } from 'react-redux';
import setupAxiosInterceptors from './config/axios-interceptor';

const store = getStore();

const rootEl = document.getElementById('root');
setupAxiosInterceptors();

const render = Component =>
  // eslint-disable-next-line react/no-render-return-value
  ReactDOM.render(
    <ErrorBoundary>
      <Provider store={store}>
        <div>
          <Component />
        </div>
      </Provider>
    </ErrorBoundary>,
    rootEl
  );

render(AppComponent);
