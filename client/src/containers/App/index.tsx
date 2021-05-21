import React from 'react';
import styles from './App.module.less';
import { Switch, Route } from 'react-router';
import Stat from '../Stat';
import Home from '../Home';

function App() {
  return (
    <Switch>
      <Route path="/stat/:userId/:configKey">
        <Stat />
      </Route>
      <Route path="/">
        <Home />
      </Route>
    </Switch>
  );
}

export default App;
