import React, { useEffect, useState } from 'react';
import { Tabs, PageHeader } from 'antd';
import axios from 'axios';
import { useHistory, useParams } from 'react-router';
import BannerStatistics from './BannerStatistics';

const { TabPane } = Tabs;

const client = axios.create({
  baseURL: '/api/v1',
});

type GachaConfig = {
  name: string;
  key: string;
};

const Stat: React.FC = () => {
  const { userId, configKey } =
    useParams<{ userId: string; configKey: string }>();
  const [gachaConfigs, setGachaConfigs] = useState<GachaConfig[]>([]);
  const history = useHistory();
  const handleTabChange = (key: string) => {
    history.push(key);
  };

  useEffect(() => {
    const fetchConfigs = async () => {
      const gachaLog = await client.get<{ data: GachaConfig[] }>('gacha');
      setGachaConfigs(gachaLog.data.data);
    };
    fetchConfigs();
  }, []);

  const tabs = [{ name: '全部', key: 'all' }, ...gachaConfigs];

  // const chartOption

  return (
    <>
      <PageHeader
        onBack={() => history.push('/')}
        title={`${userId}`}
        subTitle={`${userId}`}
      />
      <Tabs onChange={handleTabChange} type="card" activeKey={configKey}>
        {tabs.map(config => (
          <TabPane tab={config.name} key={config.key}>
            <BannerStatistics />
          </TabPane>
        ))}
      </Tabs>
    </>
  );
};

export default Stat;
