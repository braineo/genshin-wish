import React from 'react';
import ReactDOM from 'react-dom';
import { BrowserRouter as Router, Route } from 'react-router-dom';
import './index.less';
import App from './App';
import reportWebVitals from './reportWebVitals';
import { setOptions } from 'genshin-db';
setOptions({
  matchAliases: false,
  matchCategories: false,
  verboseCategories: false,
  queryLanguages: ['English'],
  resultLanguage: 'CHS',
});

ReactDOM.render(
  <React.StrictMode>
    <Router>
      <Route path="/:configKey?">
        <App />
      </Route>
    </Router>
  </React.StrictMode>,
  document.getElementById('root'),
);

// If you want to start measuring performance in your app, pass a function
// to log results (for example: reportWebVitals(console.log))
// or send to an analytics endpoint. Learn more: https://bit.ly/CRA-vitals
reportWebVitals();
