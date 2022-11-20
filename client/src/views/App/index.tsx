import { Redirect, Route, Switch } from 'react-router-dom';
import Home from '../Home';
import Stat from '../Stat';
import { WishItemList } from '../WishItemList';

function App() {
  return (
    <Switch>
      <Route path="/stat/:userId/:gachaType">
        <Stat />
      </Route>
      <Route path="/list/:userId/">
        <WishItemList />
      </Route>
      <Route path="">
        <Home />
      </Route>
      <Redirect to="/" />
    </Switch>
  );
}

export default App;
