import React from 'react';
import { useParams } from 'react-router';
import RecentItems from './RecentItems';

const Recents: React.FC = () => {
  const { userId, configKey } =
    useParams<{ userId: string; configKey: string }>();
  const itemConfigs = [
    { title: '最近记录', rarity: '' },
    { title: '最近五星', rarity: '5' },
    { title: '最近四星', rarity: '4' },
  ];
  return (
    <div>
      {itemConfigs.map((data, index) => (
        <RecentItems
          key={index}
          title={data.title}
          rarity={data.rarity}
          userId={userId}
          gachaType={configKey === 'all' ? '' : configKey}
          size={3}
        />
      ))}
    </div>
  );
};

export default Recents;
