import { Col } from 'antd';
import React from 'react';
import { useParams } from 'react-router';
import styles from './index.module.less';
import RecentItems from './RecentItems';

const Recents: React.FC = () => {
  const { userId, gachaType } = useParams<{
    userId: string;
    gachaType: string;
  }>();
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
          // FIXME: handle mix gacha pool
          gachaType={
            gachaType === 'all'
              ? ''
              : gachaType === '301'
              ? '301+400'
              : gachaType
          }
          size={5}
        />
      ))}
    </Col>
  );
};

export default Recents;
