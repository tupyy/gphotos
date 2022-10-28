import React from 'react';
import Home from './modules/home/home';
import '@elastic/eui/dist/eui_theme_dark.css';
import {BrowserRouter as Router, Routes, Route } from "react-router-dom";

export const App = () => {
    return (
      <Router>
        <Routes>
          <Route path="/" element={<Home />} />
        </Routes>
      </Router>
    );
};

export default App;
