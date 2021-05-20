import React, { useEffect, useState } from 'react';
import { useAxios } from '../../utils/axios';
import { WishLog } from 'genshin-wish';
import ItemList from '../../components/ItemList';

type RecentItemsProps = {
  userId: string;
  rarity: string;
  gachaType: string;
  size: number;
  title: string;
};

const Recents: React.FC<RecentItemsProps> = props => {
  const [recentLogs, setRecentLogs] = useState<WishLog[]>([]);
  const client = useAxios();
  useEffect(() => {
    const fetchLog = async () => {
      const recent = await client.get<{ data: WishLog[] }>(
        `log/${props.userId}`,
        {
          params: {
            rarity: props.rarity,
            gachaType: props.gachaType,
            size: props.size,
          },
        },
      );
      setRecentLogs(recent.data.data);
    };
    fetchLog();
  }, []);

  return <ItemList wishLogs={recentLogs} />;
};

export default Recents;
