import classNames from 'classnames';
import {
  Character,
  QueryOptions,
  Weapon,
  characters,
  weapons,
} from 'genshin-db';
import React from 'react';
import { toEnElement } from '../../utils';
import RarityIndicator from '../RarityIndicator';
import styles from './index.module.less';

const queryOption = {
  matchAliases: false,
  matchCategories: false,
  verboseCategories: false,
  queryLanguages: ['English'],
  resultLanguage: 'CHS',
};

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

const dateTimeFormatter = Intl.DateTimeFormat('zh', {
  timeZone: 'Asia/Shanghai',
  year: 'numeric',
  month: 'numeric',
  day: 'numeric',
});

const ItemName: React.FC<{ itemInfo: Character | Weapon | null }> = props => {
  const { itemInfo } = props;
  if (!itemInfo) {
    return <span></span>;
  }
  return (
    <div
      className={classNames(
        styles.name,
        isCharacter(itemInfo)
          ? styles[`${toEnElement(itemInfo.element)}Text`]
          : '',
      )}
    >{`${itemInfo.name}(${
      isCharacter(itemInfo) ? itemInfo.element : itemInfo.weapontype
    })`}</div>
  );
};

const ItemCard: React.FC<ItemCardProps> = props => {
  let itemInfo: Character | Weapon | undefined;
  if (props.itemType === 'character') {
    itemInfo = characters(props.itemId, queryOption as unknown as QueryOptions);
  } else {
    itemInfo = weapons(props.itemId, queryOption as unknown as QueryOptions);
  }
  if (!itemInfo) {
    return <div />;
  }

  const Avatar: React.FC = () => {
    return (
      <div className={styles.avatar}>
        <span
          className={classNames(
            styles.iconWrapper,
            styles[`star${props.rarity}Bg`],
          )}
        >
          <img
            className={styles.icon}
            src={itemInfo?.images.icon}
            alt="item-icon"
          />
        </span>
        <RarityIndicator rarity={parseInt(props.rarity)} />
      </div>
    );
  };

  return (
    <li className={styles.item}>
      <Avatar />
      <ItemName itemInfo={itemInfo} />
      <div className={styles.datetime}>
        {dateTimeFormatter.format(new Date(props.time))}
      </div>
      {props.rarity === '4' ? (
        <div className={styles.pity}>{props.pityStar4}</div>
      ) : null}
      {props.rarity === '5' ? (
        <div className={styles.pity}>{props.pityStar5}</div>
      ) : null}
    </li>
  );
};

export default ItemCard;
