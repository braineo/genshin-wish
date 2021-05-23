import React from 'react';
import ItemCard from '../ItemCard';
import { WishLog } from 'genshin-wish';
import styles from './index.module.less';

type ItemListProps = {
  wishLogs: WishLog[];
};

const ItemList: React.FC<ItemListProps> = props => {
  return (
    <ul className={styles.itemList}>
      {props.wishLogs.map((log, index) => (
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
