import { log } from 'console';
import React, { useEffect, useState } from 'react';
import { useAxios } from '../../axios';
import ItemCard from '../ItemCard';
type GachaLog = {
  gachaType: string;
  id: string;
  Item: GachaItem;
  time: string;
  pityStar4: number;
  pityStar5: number;
};
type GachaItem = {
  id: string;
  name: string;
  type: 'weapon' | 'character';
  rarity: string;
};
const ItemList: React.FC = () => {
  const [gachaLogs, setGachaLogs] = useState<GachaLog[]>([]);
  const client = useAxios();
  useEffect(() => {
    const fetchLog = async () => {
      const gachaLog = await client.get<{ data: GachaLog[] }>('log/820575774', {
        params: {
          rarity: '5',
          type: 'character',
        },
      });
      setGachaLogs(gachaLog.data.data);
    };
    fetchLog();
  }, []);

  return (
    <ul>
      {gachaLogs.map(log => (
        <ItemCard
          itemType={log.Item.type}
          itemId={log.Item.id}
          rarity={log.Item.rarity}
          {...log}
        />
      ))}
    </ul>
  );
};

export default ItemList;
