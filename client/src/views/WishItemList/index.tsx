import { Table } from 'antd';
import { ColumnsType } from 'antd/lib/table';
import axios from 'axios';
import {
  Character,
  characters,
  QueryOptions,
  Weapon,
  weapons,
} from 'genshin-db';
import { WishLog } from 'genshin-wish';
import { useEffect, useState } from 'react';
import { useParams } from 'react-router';
import { Avatar } from '../../components/ItemCard';

const columns = [
  {
    title: '',
    dataIndex: ['Item'],
    render: (value: Character | Weapon) => (
      <Avatar itemInfo={value} rarity={value.rarity} />
    ),
  },
  { title: '物品名', dataIndex: ['Item', 'name'] },
  { title: '分类', dataIndex: ['Item', 'type'] },
  { title: '稀有度', dataIndex: ['Item', 'rarity'] },
  { title: '抽数', dataIndex: 'pityStar5' },
  { title: '日期', dataIndex: 'time' },
];

const queryOption = {
  matchAliases: false,
  matchCategories: false,
  verboseCategories: false,
  queryLanguages: ['English'],
  resultLanguage: 'CHS',
};

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
    <Table
      columns={columns}
      dataSource={wishLogs.map(wishLog => {
        if (wishLog.Item.type === 'character') {
          return {
            ...wishLog,
            Item: {
              ...wishLog.Item,
              ...characters(wishLog.Item.id, queryOption as QueryOptions),
            },
          };
        } else {
          return {
            ...wishLog,
            Item: {
              ...wishLog.Item,
              ...weapons(wishLog.Item.id, queryOption as QueryOptions),
            },
          };
        }
      })}
    />
  );
};
