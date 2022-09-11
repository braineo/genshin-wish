import { Button } from 'antd';
import { WishLog } from 'genshin-wish';
import React, { useCallback, useEffect, useState } from 'react';
import { useHistory } from 'react-router-dom';
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
  const history = useHistory();
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

  const onSeeAll = useCallback(() => {
    history.push(`/list/${props.userId}/`);
  }, [history, props.userId]);

  return (
    <div className={styles.recentItems}>
      <div className={styles.title}>
        <div className={styles.titleText}>{props.title}</div>
        <Button onClick={onSeeAll}>查看全部</Button>
      </div>
      <ItemList wishLogs={recentLogs} />
    </div>
  );
};

export default Recents;
