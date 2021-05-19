import React from 'react';
import classNames from 'classnames';
import { Character, character, weapon, Weapon } from 'genshin-db';
import styles from './index.module.less';
import RarityIndicator from '../RarityIndicator';
import { toEnElement } from '../../utils';

type ItemCardProps = {
  itemId: string;
  pityStar4: number;
  pityStar5: number;
  time: string;
  itemType: 'weapon' | 'character';
  rarity: string;
};

const isCharacter = (item: Character | Weapon): item is Character => {
  if (Object.keys(item).includes('gender')) {
    return true;
  }
  return false;
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
      <div className={styles.avatar}>
        <span
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
        </span>
        <RarityIndicator rarity={parseInt(props.rarity)} />
      </div>
      <span></span>
      <span
        className={classNames(
          styles.name,
          isCharacter(itemInfo)
            ? styles[`${toEnElement(itemInfo.element)}Text`]
            : '',
        )}
      >{`${itemInfo.name}(${
        isCharacter(itemInfo) ? itemInfo.element : itemInfo.weapontype
      })`}</span>
      {props.rarity === '4' ? (
        <span className={styles.pity}>{props.pityStar4}</span>
      ) : null}
      {props.rarity === '5' ? (
        <span className={styles.pity}>{props.pityStar5}</span>
      ) : null}
    </li>
  );
};

export default ItemCard;
