import React, { useEffect, useState } from 'react';
import { useAxios } from '../../utils/axios';
import ItemCard from '../ItemCard';
import { WishLog } from 'genshin-wish';

const ItemList: React.FC = () => {
  const [wishLogs, setWishLogs] = useState<WishLog[]>([]);
  const client = useAxios();
  useEffect(() => {
    const fetchLog = async () => {
      const gachaLog = await client.get<{ data: WishLog[] }>('log/815648055', {
        params: {
          rarity: '5',
          itemType: '',
        },
      });
      setWishLogs(gachaLog.data.data);
    };
    fetchLog();
  }, []);

  return (
    <ul>
      {wishLogs.map((log, index) => (
        <ItemCard
          key={index}
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
