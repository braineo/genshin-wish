import { WishLog } from 'genshin-wish';
import React, { useEffect, useState } from 'react';
import ItemList from '../../components/ItemList';
import { useAxios } from '../../utils/axios';
import styles from './index.module.less';

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

  return (
    <div className={styles.recentItems}>
      {props.title}
      <ItemList wishLogs={recentLogs} />
    </div>
  );
};

export default Recents;
