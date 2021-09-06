import React from 'react';
import { Col } from 'antd';
import { useParams } from 'react-router';
import RecentItems from './RecentItems';
import styles from './index.module.less';

const Recents: React.FC = () => {
  const { userId, gachaType } =
    useParams<{ userId: string; gachaType: string }>();
  const itemConfigs = [
    { title: '最近五星', rarity: '5' },
    { title: '最近四星', rarity: '4' },
  ];
  return (
    <Col>
      {itemConfigs.map((data, index) => (
        <RecentItems
          key={index}
          title={data.title}
          rarity={data.rarity}
          userId={userId}
          gachaType={gachaType === 'all' ? '' : gachaType}
          size={5}
        />
      ))}
    </Col>
  );
};

export default Recents;
