import { PageHeader, Tabs } from 'antd';
import axios from 'axios';
import { useEffect, useState } from 'react';
import { useHistory, useParams } from 'react-router-dom';
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
  const { userId, gachaType } = useParams<{
    userId: string;
    gachaType: string;
  }>();
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
        title="祈愿概要"
        subTitle={`${userId}`}
      />
      <Tabs onChange={handleTabChange} type="card" activeKey={gachaType}>
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
