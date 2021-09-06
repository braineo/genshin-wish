import axios from 'axios';
import { WishLog } from 'genshin-wish';
import { useEffect, useState } from 'react';
import { useParams } from 'react-router';
import ItemList from '../../components/ItemList';

const client = axios.create({
  baseURL: '/api/v1',
});

export const WishItemList = () => {
  const [wishLogs, setWishLogs] = useState<WishLog[]>([]);
  const { userId, gachaType } =
    useParams<{ userId: string; gachaType: string }>();

  useEffect(() => {
    const fetchLog = async () => {
      const gachaLog = await client.get<{ data: WishLog[] }>(`log/${userId}`, {
        params: {
          gachaType: gachaType === 'all' ? '' : gachaType,
        },
      });
      setWishLogs(gachaLog.data.data);
    };
    fetchLog();
  }, []);

  return (
    <div>
      <ItemList wishLogs={wishLogs} />
    </div>
  );
};
