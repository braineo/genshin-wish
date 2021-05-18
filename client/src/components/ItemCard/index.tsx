import React from 'react';
import classNames from 'classnames';
import { Character, character, weapon, Weapon } from 'genshin-db';
import styles from './index.module.less';

type ItemCardProps = {
  itemId: string;
  pityStar4: string;
  pityStar5: string;
  time: string;
  itemType: 'weapon' | 'character';
  rarity: string;
};

const ItemCard: React.FC<ItemCardProps> = props => {
  let itemInfo: Character | Weapon | null;
  if (props.itemType === 'character') {
    itemInfo = character(props.itemId);
  } else {
    itemInfo = weapon(props.itemId);
  }
  if (!itemInfo) {
    return <div />;
  }

  return (
    <li className={styles.item}>
      <div
        className={classNames(
          styles.iconWrapper,
          styles[`star${props.rarity}Bg`],
        )}
      >
        <img
          className={styles.icon}
          src={itemInfo.images.icon}
          alt="item-icon"
        />
      </div>
      <div>{itemInfo.name}</div>
    </li>
  );
};

export default ItemCard;
