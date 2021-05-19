import React, { useEffect, useState } from 'react';
import { useAxios } from '../../utils/axios';
import styles from './App.module.less';
import { Tabs } from 'antd';
import { useHistory, useParams } from 'react-router';
import Stat from '../Stat';
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
          <Stat />
        </TabPane>
      ))}
    </Tabs>
  );
}

export default App;
