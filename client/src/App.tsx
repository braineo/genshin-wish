import React, { useEffect, useState } from 'react';
import { useAxios } from './axios';
import styles from './App.module.less';
import { Tabs } from 'antd';
import { useHistory, useParams } from 'react-router';
import ItemCard from './components/ItemCard';
import ItemList from './components/ItemList';
const { TabPane } = Tabs;

type GachaConfig = {
  name: string;
  key: string;
};

function App() {
  const client = useAxios();
  const [gachaConfigs, setGachaConfigs] = useState<GachaConfig[]>([]);
  const params = useParams<{ configKey: string }>();
  const history = useHistory();
  useEffect(() => {
    const fetchConfigs = async () => {
      const gachaLog = await client.get<{ data: GachaConfig[] }>('gacha');
      setGachaConfigs(gachaLog.data.data);
    };
    fetchConfigs();
  }, []);

  const handleTabChange = (key: string) => {
    history.push(key);
  };

  const tabs = [{ name: '全部', key: 'all' }, ...gachaConfigs];

  return (
    <Tabs onChange={handleTabChange} type="card" activeKey={params.configKey}>
      {tabs.map(config => (
        <TabPane tab={config.name} key={config.key}>
          <ItemList />
        </TabPane>
      ))}
    </Tabs>
  );
}

export default App;
