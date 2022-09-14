import React from 'react';
import Home from './modules/home/home';
import Header from './shared/layout/header/header';
import {createTheme, ThemeProvider} from '@mui/material/styles';
import Container from '@mui/material/Container';

const darkTheme = createTheme({
  palette: {
    mode: 'dark',
  },
});

export const App = () => {
    return (
      <ThemeProvider theme={darkTheme}>
        <div>
          <Header/>
          <Container>
            <Home/>
          </Container>
        </div>
      </ThemeProvider>
    );
};

export default App;
