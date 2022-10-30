import React from 'react';
import { Route, Switch } from 'react-router-dom';
import Home from '../Home';
import Stat from '../Stat';
import { WishItemList } from '../WishItemList';
import styles from './App.module.less';

function App() {
  return (
    <Switch>
      <Route path="/stat/:userId/:gachaType">
        <Stat />
      </Route>
      <Route path="/list/:userId/">
        <WishItemList />
      </Route>
      <Route path="/">
        <Home />
      </Route>
    </Switch>
  );
}

export default App;
